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
	"io/ioutil"
	"os"

	logrus "github.com/sirupsen/logrus"
	github "k8s.io/test-infra/prow/github"
)

const (
	defaultWebhookPath = "/hook"

	// env var names
	// supplied by GH Action Runtmie
	ghEventPath = "GITHUB_EVENT_PATH"
	ghEventName = "GITHUB_EVENT_NAME"
	ghRepo      = "GITHUB_ACTION_REPOSITORY"

	// Configured by project admistrators on the repo as secrets
	// ${{secrets.oauth}}
	repoOauthToken = "REPO_OAUTH_TOKEN" // Stored as a secret on the repo (org level also??)
	prowPlugin  = "PROW_PLUGIN"         // Just one for now, list of plugins later?
)

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
}

// comments tagged #27150 refer to issue on k8s/test-infra
func main() {
	// #27150 no Command Line Options, Github runtime supplied env vars only
	eventName := os.Getenv(ghEventName)
	repo := os.Getenv(ghRepo)

	eventPayload := getGithubEventPayload()
	ghClient := getGithubClient()
	//	plugin := getProwPlugin()

	err := processGithubAction(eventName, "GUID???", eventPayload, repo, ghClient)
	if err != nil {
		logrus.WithError(err).Errorf("Error demuxing event %s", eventName)
	}
}

func getGithubClient() github.Client {
	oauthToken:= os.Getenv(repoOauthToken)
	options := new(github.ClientOptions)
	options.GetToken = func() []byte { return []byte(oauthToken)}

	_, _, ghClient , err := github.NewClientFromOptions(logrus.Fields{}, (*options))
	if err != nil {
		logrus.WithError(err).Errorf("Error creating a GH Client", )
	}
	return ghClient
}

// have to select a prow plugin to use
func getProwPlugin() {
	prowPlugin := os.Getenv(prowPlugin)
	if prowPlugin == "" {
		logrus.Fatalf("Env var %s is not set\n", prowPlugin)
	}
	// TODO load plugin here
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

func processGithubAction(eventType, eventGUID string, payload []byte, srcRepo string, ghclient github.Client) error {
	l := logrus.WithFields(
		logrus.Fields{
			"eventType":        eventType,
			"github.EventGUID": eventGUID,
		},
	)
	switch eventType {
	case "issues":
		var i github.IssueEvent
		if err := json.Unmarshal(payload, &i); err != nil {
			return err
		}
		i.GUID = eventGUID
		srcRepo = i.Repo.FullName
	case "issue_comment":
		var ice github.IssueCommentEvent
		if err := json.Unmarshal(payload, &ice); err != nil {
			return err
		}
		// ice.GUID = eventGUID
		// srcRepo = ice.Repo.FullName
		handleIssueCommentEvent(ice)
	case "pull_request":
		var pr github.PullRequestEvent
		if err := json.Unmarshal(payload, &pr); err != nil {
			return err
		}
		pr.GUID = eventGUID
		srcRepo = pr.Repo.FullName
	default:
		var ge github.GenericEvent
		if err := json.Unmarshal(payload, &ge); err != nil {
			return err
		}
		srcRepo = ge.Repo.FullName
		l.Debug("Ignoring unhandled event type. ( k8s/test-infra issue #27150 No external plugins for now.)")
	}
	return nil
}

func handleIssueCommentEvent(ice github.IssueCommentEvent) {
	logrus.Infof("ice %v", ice)
	logrus.Infof("ice %v", ice)
}
