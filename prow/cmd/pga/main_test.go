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
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"k8s.io/test-infra/prow/plugins"
)

// Make sure that pga's plugins are valid.
func TestPlugins(t *testing.T) {
	pa := &plugins.ConfigAgent{}
	if err := pa.Load("./kodata/plugins.yaml", nil, "", true, false); err != nil {
		t.Fatalf("Could not load plugins: %v.", err)
	}
}

// pga is env var driven so we can make use of that fact to run local tests.
//
// for each test case
// create the scenario you want to test on Github Actions
//   capture data from Github Actions on the Job Run by going to the Summary Page of
//   the job and downloading the job artefacts, envFile and event.json
//   organise env vars in folders by plugin and test case setting up
//     ./test-data/PLUGIN_NAME/TEST_CASE/envFiles
//     ./test-data/PLUGIN_NAME/TEST_CASE/event.json
//
// In envFiles setting PGA_LOCAL to any string
// allow us to config data files from ./kodata
// instead of where ko build places those files
//
// This test then runs the relevent functions in pga to describe pga's behaviour
// GITHUB_EVENT_PATH to file w/ event data

func TestRunAllGutHubActionTests(t *testing.T) {
	cases := []struct {
		name              string
		envFile           string
		eventFile         string
		expectedComment   string
		expectedCommand   string
		expectedLabels    []string
		expectedAssignees []string
		err               bool
	}{
		{
			name:              "The yuk plugin tells a joke",
			envFile:           "./test-data/yuk/env",
			eventFile:         "./test-data/yuk/event",
			expectedCommand:   "/joke",
			expectedComment:   "",
			expectedLabels:    []string{},
			expectedAssignees: []string{},
			err:               false,
		},
		{
			name:              "The Blunderbuss plugins assigns two reviewers to this PR",
			envFile:           "./test-data/blunderbuss/auto-cc/env",
			eventFile:         "./test-data/blunderbuss/auto-cc/event.json",
			expectedCommand:   "/auto-cc",
			expectedComment:   "",
			expectedLabels:    []string{},
			expectedAssignees: []string{"@RobKielty", "@BobyMcBobs"},
			err:               false,
		},
	}
	// Range over the cases, read in env vars
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Add env vars to runtime
			err := setEnvVarsFromFile(tc.envFile, t)
			// For now temp set REPO_OAUTH_TOKEN on env thss making this not a unit test!
			// TODO Inject fake GH client for proper unit testing
			if err != nil {
				t.Logf("Skipping %s. Cause : %v ", tc.name, err)
				t.SkipNow()
			}

			eventName := getMandatoryEnvVar(ghEventName)
			repo := getMandatoryEnvVar(ghRepo)

			eventPayload := getGithubEventPayload()
			clientConfig = getClientConfig(repo)
			t.Logf("eventName: %v", eventName)
			t.Logf("repo: %v", repo)
			t.Logf("clientConfig: %v", clientConfig)
			t.Logf("eventPayload[0-180]: %v", string(eventPayload[0:180]))
			pgaErr := processGithubAction(eventName, eventPayload, repo, ghClient)
			// Check expected state
			if pgaErr == nil {
				fmt.Println("pga ran!?")
			} else {
				fmt.Printf("pga error %v", pgaErr)
			}

		})
	}
}

func setEnvVarsFromFile(envFile string, t *testing.T) error {
	data, err := ioutil.ReadFile(envFile)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if len(line) > 0 && line[0:1] != "#" { // ignore empty lines and #-commented out lines
			kv := strings.Split(string(line), "=")
			t.Setenv(kv[0], kv[1])
		}
	}
	return nil
}
