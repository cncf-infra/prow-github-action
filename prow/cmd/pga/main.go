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
	logrus "github.com/sirupsen/logrus"
	"io/ioutil"
	github "k8s.io/test-infra/prow/github"
	"os"
)

const (
	defaultWebhookPath = "/hook"
	ghEventPath        = "GITHUB_EVENT_PATH"
	ghEventName        = "GITHUB_EVENT_NAME"
	ghRepo             = "GITHUB_ACTION_REPOSITORY"
)

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.WarnLevel)
}

// comments tagged #27150 refer to issue on k8s/test-infra
func main() {
	// #27150 no Command Line Options, Github runtime supplied env vars only
	eventName := os.Getenv(ghEventName)
	eventPath := os.Getenv(ghEventPath)
	repo := os.Getenv(ghRepo)
	// TODO: event = FILE.READ(eventPath)
	eventBody, err := ioutil.ReadFile(eventPath)
	// if err != nil {
	// 	log.Fatalf("unable to read file: %v", err)
	// }
	err = processGithubAction(eventName, "GUID???", []byte(eventBody), repo)
	if err != nil {
		logrus.WithError(err).Errorf("Error demuxing event %s", eventName)
	}

}

// #27150 https://github.com/kubernetes/test-infra/blob/master/prow/hook/server.go#L91-L176
// Inspired by demuxEvent in above ref

func processGithubAction(eventType, eventGUID string, payload []byte, srcRepo string) error {
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
		var ic github.IssueCommentEvent
		if err := json.Unmarshal(payload, &ic); err != nil {
			return err
		}
		ic.GUID = eventGUID
		srcRepo = ic.Repo.FullName
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
