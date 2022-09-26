/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"sync"

	logrus "github.com/sirupsen/logrus"
	"k8s.io/test-infra/prow/bugzilla"
	"k8s.io/test-infra/prow/config"
	"k8s.io/test-infra/prow/git/v2"
	github "k8s.io/test-infra/prow/github"
	_ "k8s.io/test-infra/prow/hook/plugin-imports"
	"k8s.io/test-infra/prow/plugins"
	"k8s.io/test-infra/prow/plugins/ownersconfig"
	"k8s.io/test-infra/prow/repoowners"
	"k8s.io/test-infra/prow/slack"
)

const (
	// env var names
	// supplied by GH Action Runtime
	ghEventPath = "GITHUB_EVENT_PATH"
	ghEventName = "GITHUB_EVENT_NAME"
	ghRepo      = "GITHUB_REPOSITORY"

	pgaLocalRun = "PGA_LOCAL"

	// Project Admins, configure OAuth Tokens on repo as a secret
	// pga will pick this up as an env var in a Github Action with ${{secrets.git statoauth}}
	repoOauthToken = "REPO_OAUTH_TOKEN" // Stored as a secret on the repo (org level also??)

	failedCommentCoerceFmt = "Could not coerce %s event to a GenericCommentEvent. Unknown 'action': %q."
)

var (
	pluginsConfig      *plugins.ConfigAgent
	clientConfig       *plugins.ClientAgent
	configurationAgent *config.Agent
	ghClient           github.Client
	configDir          string
	wg                 sync.WaitGroup
)

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)

	setupConfigDirLocation()

	pluginsConfig = getProwPluginConfigAgent()
	configurationAgent = getConfigAgent()
}

// ghaWriteEnvVarsToFile writes runtime environment variables (env vars) to a file called env
//
// if and only if...
//   - logrus.DebugLevel is set
//   - the env var is NOT present in the []string returned by getUnexportedEnvVars() string[]
//
// When run as part of a Github Action the env file can be found by visiting the
// Github Repo > Actions > the Workflow Run > Summary >  Artifacts > environment variables and event payload
//
// Implemented for development of pga setting up tests
//
// For this to work the following yaml needs to be a the Github action Workflow
//
// ``` yaml
// - name: "Capture runtime artefacts for development on pga"
// uses: actions/upload-artifact@v3
// with:
//   name: environment variables and event payload
//   path: ${{ github.workspace }}
//
func ghaWriteEnvVarsToFile() {
	if logrus.GetLevel() == logrus.DebugLevel {
		env := os.Environ()
		var b []byte
		var ignoredEnvVars = getUnexportedEnvVars()
		var skip bool = false

		for _, s := range env {
			for _, excluded := range ignoredEnvVars {
				if s == excluded {
					skip = true
				}
				break
			}
			if !skip {
				envTuple := s + string('\n')
				b = append(b, envTuple...)
				skip = false
			}
		}
		storeDataAsArtefact("env", b)
	}
}

// Returns a list of env vars that should not be archived
func getUnexportedEnvVars() []string {
	return []string{"PATH", "UNSPLASH_ACCESS_KEY"}
}

// comments tagged #27150 refer to issue number on k8s/test-infra
func main() {
	// #27150 no Command Line Options, Github runtime supplied env vars only
	ghaWriteEnvVarsToFile()
	eventName := getMandatoryEnvVar(ghEventName)
	repo := getMandatoryEnvVar(ghRepo)

	eventPayload := getGithubEventPayload()
	clientConfig = getClientConfig(repo)
	err := processGithubAction(eventName, eventPayload, repo, ghClient)
	if err != nil {
		logrus.WithError(err).Errorf("Error demuxing event %s", eventName)
	}
	// Wait for all handlers to complete.
	wg.Wait()
}

// setupConfigDirLocation changes the dir used to search for config files
// so that pga can run from the command line in a dev environment (a "local
// run") or can run as part of the container image that implements the Custom
// Github Action, prow-github-action
// Cofiguration directories searched
// ./kodata     for local runs
// /var/run/ko/ for running in ko built container
func setupConfigDirLocation() {
	if thisIsALocalRun() {
		configDir = "./kodata/"
		logrus.Infof("Local run. Env var %s is set. Config Dir %s", pgaLocalRun, configDir)

	} else {
		configDir = "/var/run/ko/"
		logrus.Infof("GHA run. Env var %s is empty. Config Dir %s", pgaLocalRun, configDir)
	}
}

// thisIsALocalRun returns true if pgaLocalRun is set as runtime env var, false otherwise
// used to flip between CI runtime confguration settings and development runtime configuration settings.
func thisIsALocalRun() bool {
	localRun := false

	runLocally := getOptionalEnvVar(pgaLocalRun)

	if len(runLocally) > 0 {
		localRun = true
	}
	return localRun
}

// getMandatoryEnvVar returns the value of envVar if set in the environment
// otherwise exits if envVar is not set.
func getMandatoryEnvVar(envVar string) string {
	value := os.Getenv(envVar)
	if value == "" {
		logrus.Fatalf("Env Var %v is not set. Exiting", envVar)
	}
	return value
}
func getOptionalEnvVar(envVar string) string {
	value := os.Getenv(envVar)
	return value
}

func getGithubClient() github.Client {
	oauthToken := getMandatoryEnvVar(repoOauthToken)

	_, userGenerator, ghClient, err := github.NewClientFromOptions(logrus.Fields{},
		github.ClientOptions{
			Censor:          func(b []byte) []byte { return b },
			GetToken:        func() []byte { return []byte(oauthToken) },
			AppID:           "",
			AppPrivateKey:   nil,
			GraphqlEndpoint: github.DefaultGraphQLEndpoint,
			Bases:           []string{github.DefaultAPIEndpoint},
			DryRun:          false,
		}.Default(),
	)
	if userGenerator != nil {
		user, err := userGenerator()
		if err != nil {
			logrus.WithError(err).Warn("pga.getGithubClient : Cannot generate user")
		}
		logrus.Infof("pga.getGithubClient : user %v", user)
	}
	if err != nil {
		logrus.WithError(err).Errorf("Error creating Github Client. Err: %v ", err)
		logrus.WithError(err).Debugf("oauthToken: %v ", oauthToken)
	}
	logrus.Debugf("GH Client created: %v ", ghClient)
	return ghClient
}

func getClientConfig(repo string) *plugins.ClientAgent {
	clientConfig = new(plugins.ClientAgent)
	clientConfig.GitClient = getGitClient(repo)
	clientConfig.GitHubClient = getGithubClient()
	if logrus.GetLevel() == logrus.DebugLevel && !thisIsALocalRun() {
		fmt.Printf("clientConfig.GitClient: %v\n", clientConfig.GitClient)
	}
	clientConfig.OwnersClient = getOwnersClient(repo)
	clientConfig.BugzillaClient = &bugzilla.Fake{}
	clientConfig.SlackClient = slack.NewFakeClient()
	return clientConfig
}
func getGitClient(repo string) git.ClientFactory {
	cf, err := git.NewClientFactory()
	if err != nil {
		log.Fatalf("Cannot initialise git client! Error is %v", err)
	}

	return cf
}
func getOwnersClient(repo string) repoowners.Interface {
	mdYAMLEnabled := func(org, repo string) bool {
		return pluginsConfig.Config().MDYAMLEnabled(org, repo)
	}
	skipCollaborators := func(org, repo string) bool {
		return pluginsConfig.Config().SkipCollaborators(org, repo)
	}
	ownersDirDenylist := func() *config.OwnersDirDenylist {
		// OwnersDirDenylist struct contains some defaults that's required by all
		// repos, so this function cannot return nil
		res := &config.OwnersDirDenylist{}
		if l := getConfigAgent().Config().OwnersDirDenylist; l != nil {
			res = l
		}
		return res
	}
	resolver := func(org, repo string) ownersconfig.Filenames {
		return pluginsConfig.Config().OwnersFilenames(org, repo)
	}

	ownersClient := repoowners.NewClient(
		clientConfig.GitClient,
		getGithubClient(),
		mdYAMLEnabled,
		skipCollaborators, ownersDirDenylist, resolver)
	return ownersClient
}

func getConfigAgent() *config.Agent {
	configAgent := &config.Agent{}
	configAgent.Start(configDir+"config.yaml", configDir+"emptyJobConfig.yaml", []string{}, "")
	return configAgent
}

func getProwPluginConfigAgent() *plugins.ConfigAgent {
	pluginConfigAgent := &plugins.ConfigAgent{}
	if err := pluginConfigAgent.Load(configDir+"plugins.yaml", nil, "", false, false); err != nil {
		// TODO if err := pluginConfigAgent.Load("./kodata/plugins.yaml", nil, "", false, true); err != nil {
		logrus.WithError(err).Fatalf("getProwPluginConfigAgent: pluginConfigAgent.Load %v", configDir+"plugins.yaml")
	}
	logrus.Debugf("IssueCommentHandlers %v", pluginConfigAgent.IssueCommentHandlers("cncf-infra", "mock-project-repo"))
	logrus.Debugf("GenericCommentHandlers %v", pluginConfigAgent.GenericCommentHandlers("cncf-infra", "mock-project-repo"))
	return pluginConfigAgent
}

func getGithubEventPayload() []byte {
	path := os.Getenv(ghEventPath)
	if path == "" {
		logrus.Fatalf("Env var %s is not set\n", ghEventPath)
	}
	eventPayload, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatalf("Error reading: %s Reason: %v", path, err)
	}
	return eventPayload
}

// #27150 https://github.com/kubernetes/test-infra/blob/master/prow/hook/server.go#L91-L176
// Inspired by demuxEvent in above ref

func processGithubAction(eventType string, payload []byte, srcRepo string, ghclient github.Client) error {
	l := logrus.WithFields(
		logrus.Fields{
			"eventType": eventType,
		},
	)
	switch eventType {
	case "issues":
		var i github.IssueEvent
		if err := json.Unmarshal(payload, &i); err != nil {
			return err
		}
	case "issue_comment":
		var event github.IssueCommentEvent
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}
		storeDataAsArtefact("issue_comment_payload.json", payload)
		handleIssueCommentEvent(event, l)
	case "pull_request":
		var pr github.PullRequestEvent
		if err := json.Unmarshal(payload, &pr); err != nil {
			return err
		}
		// pr.GUID = eventGUID
		// srcRepo = pr.Repo.FullName
		storeDataAsArtefact("pull_request_payload.json", payload)
		handlePullRequestEvent(l, pr)
	default:
		var ge github.GenericEvent
		if err := json.Unmarshal(payload, &ge); err != nil {
			return err
		}
		l.Debug("Ignoring unhandled event type. ( k8s/test-infra issue #27150 No external plugins for now.)")
	}
	return nil
}

func storeDataAsArtefact(fileName string, data []byte) {

	if logrus.GetLevel() == logrus.DebugLevel && !thisIsALocalRun() {
		path, err := os.Getwd()
		if err != nil {
			logrus.Debug("Could not get working directory")
		}

		logrus.Debugf("working directory %s\n", path)
		e := os.WriteFile(fileName, data, 0644)
		if e != nil {
			logrus.Error(err)
		}
	}
}

func handlePullRequestEvent(l *logrus.Entry, prEvent github.PullRequestEvent) {
	l.Infof("%s", prEvent.Action)
	prHandlerMap := pluginsConfig.PullRequestHandlers(
		prEvent.PullRequest.Base.Repo.Owner.Login,
		prEvent.PullRequest.Base.Repo.Name)
	if len(prHandlerMap) == 0 {
		l.Debugf("No PR comment handlers configured for %v ", prEvent.PullRequest.Base.Repo.Name)
	}
	for pluginName, handler := range prHandlerMap {
		wg.Add(1)
		go func(pluginName string, handler plugins.PullRequestHandler) {
			defer wg.Done()
			agent := plugins.NewAgent(configurationAgent, pluginsConfig, clientConfig, prEvent.Repo.Owner.Login, nil, l, pluginName)
			agent.InitializeCommentPruner(
				prEvent.Repo.Owner.Login,
				prEvent.Repo.Name,
				prEvent.PullRequest.Number,
			)
			err := errorOnPanic(func() error { return handler(agent, prEvent) })
			if err != nil {
				agent.Logger.WithError(err).Error("pga handlingPullRequestEvent")
			}
		}(pluginName, handler)
	}
	action := genericCommentAction(string(prEvent.Action))
	handleGenericComment(
		l,
		&github.GenericCommentEvent{
			ID:           prEvent.PullRequest.ID,
			NodeID:       prEvent.PullRequest.NodeID,
			GUID:         prEvent.GUID,
			IsPR:         true,
			Action:       action,
			Body:         prEvent.PullRequest.Body,
			HTMLURL:      prEvent.PullRequest.HTMLURL,
			Number:       prEvent.PullRequest.Number,
			Repo:         prEvent.Repo,
			User:         prEvent.PullRequest.User,
			IssueAuthor:  prEvent.PullRequest.User,
			Assignees:    prEvent.PullRequest.Assignees,
			IssueState:   prEvent.PullRequest.State,
			IssueTitle:   prEvent.PullRequest.Title,
			IssueBody:    prEvent.PullRequest.Body,
			IssueHTMLURL: prEvent.PullRequest.HTMLURL,
		},
	)
}
func handleIssueCommentEvent(event github.IssueCommentEvent, l *logrus.Entry) {

	l = l.WithFields(logrus.Fields{
		// 	github.OrgLogField:  event.Repo.Owner.Login,
		// 	github.RepoLogField: event.Repo.Name,
		// 	github.PrLogField:   event.Issue.Number,
		// 	"author":            event.Comment.User.Login,
		// 	"url":               event.Comment.HTMLURL,
	})
	l.Debugf("HANDLING %v on %v", event.Action, event.Issue.ID)

	commentHandlerMap := pluginsConfig.IssueCommentHandlers(event.Repo.Owner.Login, event.Repo.Name)

	l.Debugf("commentHandlerMap %v ", commentHandlerMap)

	i := 0
	for pluginName, handler := range commentHandlerMap {
		i++
		wg.Add(1)
		l := logrus.WithFields(
			logrus.Fields{
				"Prow Plugin": pluginName,
				"handler":     handler,
			},
		)

		l.Debugf("Plugin NUMBER %d", i)
		go func(pluginName string, handler plugins.IssueCommentHandler) {
			l.Debugf("IN ISSUE COMMENTHANDLER  %v", pluginName)
			defer wg.Done()
			agent := plugins.NewAgent(
				configurationAgent, pluginsConfig, clientConfig, event.Repo.Owner.Login, nil, l, pluginName)
			agent.InitializeCommentPruner(
				event.Repo.Owner.Login,
				event.Repo.Name,
				event.Issue.Number,
			)
			// start := time.Now()
			err := errorOnPanic(func() error { return handler(agent, event) })
			//labels := prometheus.Labels{"event_type": logrus.Data[eventTypeField].(string), "action": string(ic.Action), "plugin": p, "took_action": strconv.FormatBool(agent.TookAction())}
			if err != nil {
				agent.Logger.WithError(err).Error("Error handling IssueCommentEvent.")
				// s.Metrics.PluginHandleErrors.With(labels).Inc()
			}
			//  s.Metrics.PluginHandleDuration.With(labels).Observe(time.Since(start).Seconds())
		}(pluginName, handler)
	}
	action := genericCommentAction(string(event.Action))
	if action == "" {
		l.Errorf(failedCommentCoerceFmt, "issue_comment", string(event.Action))
		return
	}
	handleGenericComment(
		l,
		&github.GenericCommentEvent{
			ID:           event.Issue.ID,
			NodeID:       event.Issue.NodeID,
			CommentID:    &event.Comment.ID,
			GUID:         event.GUID,
			IsPR:         event.Issue.IsPullRequest(),
			Action:       action,
			Body:         event.Comment.Body,
			HTMLURL:      event.Comment.HTMLURL,
			Number:       event.Issue.Number,
			Repo:         event.Repo,
			User:         event.Comment.User,
			IssueAuthor:  event.Issue.User,
			Assignees:    event.Issue.Assignees,
			IssueState:   event.Issue.State,
			IssueTitle:   event.Issue.Title,
			IssueBody:    event.Issue.Body,
			IssueHTMLURL: event.Issue.HTMLURL,
		},
	)

}

func handleGenericComment(l *logrus.Entry, ce *github.GenericCommentEvent) {
	gnrcCommentHandlerMap := getProwPluginConfigAgent().GenericCommentHandlers(ce.Repo.Owner.Login, ce.Repo.Name)
	for p, h := range gnrcCommentHandlerMap {
		wg.Add(1)
		go func(pluginName string, h plugins.GenericCommentHandler) {
			defer wg.Done()
			// agent := plugins.NewAgent(s.ConfigAgent, s.Plugins, s.ClientAgent, ce.Repo.Owner.Login, s.Metrics.Metrics, l, p)
			agent := plugins.NewAgent(configurationAgent, pluginsConfig, clientConfig, ce.Repo.Owner.Login, nil, l, pluginName)
			agent.InitializeCommentPruner(
				ce.Repo.Owner.Login,
				ce.Repo.Name,
				ce.Number,
			)
			// start := time.Now()
			err := errorOnPanic(func() error { return h(agent, *ce) })
			// labels := prometheus.Labels{"event_type": l.Data[eventTypeField].(string), "action": string(ce.Action), "plugin": p, "took_action": strconv.FormatBool(agent.TookAction())}
			if err != nil {
				agent.Logger.WithError(err).Errorf("handleGenericComment error for the %v plugin on its handler %v", pluginName, h)
			}
			// Metrics.PluginHandleDuration.With(labels).Observe(time.Since(start).Seconds())
		}(p, h)
	}
}

func errorOnPanic(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught: |%v|", r)
			debug.PrintStack()
		}
	}()
	return f()
}

// genericCommentAction normalizes the action string to a GenericCommentEventAction or returns ""
// if the action is unrelated to the comment text. (For example a PR 'label' action.)
func genericCommentAction(action string) github.GenericCommentEventAction {
	switch action {
	case "created", "opened", "submitted":
		return github.GenericCommentActionCreated
	case "edited":
		return github.GenericCommentActionEdited
	case "deleted", "dismissed":
		return github.GenericCommentActionDeleted
	}
	// The action is not related to the text body.
	return ""
}
