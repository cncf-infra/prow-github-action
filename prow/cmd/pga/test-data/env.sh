#!/bin/shA
# Description:
# A local enve file for local development env lifted from a ghithub action runner and hacked to work locally
export GITHUB_WORKSPACE=/home/runner/work/mock-project-repo/mock-project-repo
export GITHUB_PATH=/home/runner/work/_temp/_runner_file_commands/add_path_147af699-3596-4bf2-9136-035f8eba92ed
export GITHUB_ACTION=__mxschmitt_action-tmate
export GITHUB_RUN_NUMBER=34
export GITHUB_TRIGGERING_ACTOR=RobertKielty
export GITHUB_REF_TYPE=branch
export GITHUB_ACTIONS=true
export GITHUB_SHA=e6ed88e1afefe1c57249145cdb0442f92a7d6f56
export GITHUB_REF=refs/heads/main
export GITHUB_REF_PROTECTED=false
export GITHUB_API_URL=https://api.github.com
export GITHUB_ENV=/home/runner/work/_temp/_runner_file_commands/set_env_147af699-3596-4bf2-9136-035f8eba92ed
export GITHUB_EVENT_PATH=./test-data/issue_comment.json
export GITHUB_EVENT_NAME=issue_comment
export GITHUB_RUN_ID=2932302119
export GITHUB_STEP_SUMMARY=/home/runner/work/_temp/_runner_file_commands/step_summary_147af699-3596-4bf2-9136-035f8eba92ed
export GITHUB_ACTOR=RobertKielty
export GITHUB_RUN_ATTEMPT=1
export GITHUB_GRAPHQL_URL=https://api.github.com/graphql
export GITHUB_SERVER_URL=https://github.com
export GITHUB_REF_NAME=main
export GITHUB_JOB=prow_github_hook
export GITHUB_REPOSITORY=cncf-infra/mock-project-repo
export GITHUB_RETENTION_DAYS=90
export GITHUB_ACTION_REPOSITORY=mxschmitt/action-tmate
export GITHUB_BASE_REF=
export GITHUB_REPOSITORY_OWNER=cncf-infra
export GITHUB_HEAD_REF=
export GITHUB_ACTION_REF=v3
export GITHUB_WORKFLOW=Assign reviewers from OWNERS files
LOCAL_TOKEN_FILE=~/.github/rk.pat
if [ -f "${LOCAL_TOKEN_FILE}" ]; then
    token="$(cat ~/.github/rk.pat)"
    export REPO_OAUTH_TOKEN=$token
else 
    echo "Local token file not found, tried $(LOCAL_TOKEN_FILE) "
    echo "Set the variable LOCAL_TOKEN_FILE to your token" 
fi
# grep -v 

# RUNNER_ARCH=X64
# RUNNER_TEMP=/home/runner/work/_temp
# ACTIONS_RUNTIME_URL=https://pipelines.actions.githubusercontent.com/bLixyNLvOFER2O8iwwVXW1FqaLVIRMasmHluEdlKaiDeVWFpgV/
# EDGEWEBDRIVER=/usr/local/share/edge_driver
# INVOCATION_ID=c010611c82af4bac86434a63a91b84f5
# JAVA_HOME_17_X64=/usr/lib/jvm/temurin-17-jdk-amd64
# ANDROID_NDK_HOME=/usr/local/lib/android/sdk/ndk/25.0.8775105
# HOMEBREW_NO_AUTO_UPDATE=1
# NVM_DIR=/home/runner/.nvm
# SGX_AESM_ADDR=1
# ANDROID_HOME=/usr/local/lib/android/sdk
# TERM=screen-256color
# ACCEPT_EULA=Y
# RUNNER_USER=runner
# ACTIONS_RUNTIME_TOKEN=eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6ImVCWl9jbjNzWFlBZDBjaDRUSEJLSElnT3dPRSJ9.eyJuYW1laWQiOiJkZGRkZGRkZC1kZGRkLWRkZGQtZGRkZC1kZGRkZGRkZGRkZGQiLCJzY3AiOiJBY3Rpb25zLkdlbmVyaWNSZWFkOjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCBBY3Rpb25zLlVwbG9hZEFydGlmYWN0czowMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAvMTpCdWlsZC9CdWlsZC8zNyBEaXN0cmlidXRlZFRhc2suR2VuZXJhdGVJZFRva2VuOjhiMWQxYTM0LWIzYzItNDMzOC05NjkwLWVmZjk0MTZmZjg2Mzo4NzQxYWM4NC0yOGQ3LTU5YzYtOGJlMy1jZjkwMjEzODU1NGUgTG9jYXRpb25TZXJ2aWNlLkNvbm5lY3QgUmVhZEFuZFVwZGF0ZUJ1aWxkQnlVcmk6MDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwLzE6QnVpbGQvQnVpbGQvMzciLCJJZGVudGl0eVR5cGVDbGFpbSI6IlN5c3RlbTpTZXJ2aWNlSWRlbnRpdHkiLCJodHRwOi8vc2NoZW1hcy54bWxzb2FwLm9yZy93cy8yMDA1LzA1L2lkZW50aXR5L2NsYWltcy9zaWQiOiJERERERERERC1ERERELUREREQtRERERC1EREREREREREREREQiLCJodHRwOi8vc2NoZW1hcy5taWNyb3NvZnQuY29tL3dzLzIwMDgvMDYvaWRlbnRpdHkvY2xhaW1zL3ByaW1hcnlzaWQiOiJkZGRkZGRkZC1kZGRkLWRkZGQtZGRkZC1kZGRkZGRkZGRkZGQiLCJhdWkiOiJmNGRmMjY3Mi0zMmQyLTRmODMtYmY0NC1mYzQ1ZGUxOTEyNTgiLCJzaWQiOiIzYTEwMGM3Zi04MTE1LTRlODQtYmE5Zi1hMDhiMTAyMGUxYTciLCJhYyI6Ilt7XCJTY29wZVwiOlwicmVmcy9oZWFkcy9tYWluXCIsXCJQZXJtaXNzaW9uXCI6M31dIiwiYWNzbCI6IjEwIiwib2lkY19zdWIiOiJyZXBvOmNuY2YtaW5mcmEvbW9jay1wcm9qZWN0LXJlcG86cmVmOnJlZnMvaGVhZHMvbWFpbiIsIm9pZGNfZXh0cmEiOiJ7XCJyZWZcIjpcInJlZnMvaGVhZHMvbWFpblwiLFwic2hhXCI6XCJlNmVkODhlMWFmZWZlMWM1NzI0OTE0NWNkYjA0NDJmOTJhN2Q2ZjU2XCIsXCJyZXBvc2l0b3J5XCI6XCJjbmNmLWluZnJhL21vY2stcHJvamVjdC1yZXBvXCIsXCJyZXBvc2l0b3J5X293bmVyXCI6XCJjbmNmLWluZnJhXCIsXCJyZXBvc2l0b3J5X293bmVyX2lkXCI6XCI2NDY1OTg3N1wiLFwicnVuX2lkX
# CI6XCIyOTMyMzAyMTE5XCIsXCJydW5fbnVtYmVyXCI6XCIzNFwiLFwicnVuX2F0dGVtcHRcIjpcIjFcIixcInJlcG9zaXRvcnlfdmlzaWJpbGl0eVwiOlwicHVibGljXCIsXCJyZXBvc2l0b3J5X2lkXCI6XCI
# 1MjUxNDcyNzNcIixcImFjdG9yX2lkXCI6XCIyMDgwNTlcIixcImFjdG9yXCI6XCJSb2JlcnRLaWVsdHlcIixcIndvcmtmbG93XCI6XCJBc3NpZ24gcmV2aWV3ZXJzIGZyb20gT1dORVJTIGZpbGVzXCIsXCJoZ
# WFkX3JlZlwiOlwiXCIsXCJiYXNlX3JlZlwiOlwiXCIsXCJldmVudF9uYW1lXCI6XCJwdXNoXCIsXCJyZWZfdHlwZVwiOlwiYnJhbmNoXCIsXCJqb2Jfd29ya2Zsb3dfcmVmXCI6XCJjbmNmLWluZnJhL21vY2s
# tcHJvamVjdC1yZXBvLy5naXRodWIvd29ya2Zsb3dzL3Byb3ctaG9vay55YW1sQHJlZnMvaGVhZHMvbWFpblwifSIsIm9yY2hpZCI6IjhiMWQxYTM0LWIzYzItNDMzOC05NjkwLWVmZjk0MTZmZjg2My5wcm93X
# 2dpdGh1Yl9ob29rLl9fZGVmYXVsdCIsImlzcyI6InZzdG9rZW4uYWN0aW9ucy5naXRodWJ1c2VyY29udGVudC5jb20iLCJhdWQiOiJ2c3Rva2VuLmFjdGlvbnMuZ2l0aHVidXNlcmNvbnRlbnQuY29tfHZzbzo
# 5OWY1MzMyOS02NTUyLTRjMTctOTk2Zi1lNjg1MjFhYWQxN2QiLCJuYmYiOjE2NjE1MDE0OTMsImV4cCI6MTY2MTUyNDI5M30.NrTdJDxuVRvJb2fXV8C9QkdU21rjYnPYTW-Ggm7yvOaBBUK02RO4xB025Ua4B
# bLMmd2Im8ErRjn2pP9U29OvO3WkiRDj3c2X8e4CbAQ-IMoE5JHPMIpTdZptgNpXL8XU-QgT80LiJ627ieBEZV_If9QN4-jkiMPo3XAv4vfP24dy9tL1xnmZqI74lUFQwxrrNuZfnc64V0FYbbkHUmwtiWm8UTY
# d0BlKvYuiIYAomyvaT8gcYq8-mLwnUo50qyAF4i27EcUKCGJ6r3kgnTFRkcq44CT8xfit-EfchUnygQ5FjjqQFbU698_FWUstSMBEr2s-_v5dc9oSQ4sWkDeg0g
# USER=runner
# TMUX_PANE=%0
# HOMEBREW_CELLAR=/home/linuxbrew/.linuxbrew/Cellar
# PIPX_HOME=/opt/pipx
# GECKOWEBDRIVER=/usr/local/share/gecko_driver
# CHROMEWEBDRIVER=/usr/local/share/chrome_driver
# SHLVL=0
# ANDROID_SDK_ROOT=/usr/local/lib/android/sdk
# VCPKG_INSTALLATION_ROOT=/usr/local/share/vcpkg
# HOMEBREW_REPOSITORY=/home/linuxbrew/.linuxbrew/Homebrew
# RUNNER_TOOL_CACHE=/opt/hostedtoolcache
# ImageVersion=20220821.1
# DOTNET_NOLOGO=1
# STATS_PFS=true
# GRAALVM_11_ROOT=/usr/local/graalvm/graalvm-ce-java11-22.2.0
# XDG_RUNTIME_DIR=/run/user/1001
# ACTIONS_ID_TOKEN_REQUEST_URL=https://pipelines.actions.githubusercontent.com/bLixyNLvOFER2O8iwwVXW1FqaLVIRMasmHluEdlKaiDeVWFpgV/00000000-0000-0000-0000-000000
# 000000/_apis/distributedtask/hubs/Actions/plans/8b1d1a34-b3c2-4338-9690-eff9416ff863/jobs/8741ac84-28d7-59c6-8be3-cf902138554e/idtoken?api-version=2.0
# AZURE_EXTENSION_DIR=/opt/az/azcliextensions
# PERFLOG_LOCATION_SETTING=RUNNER_PERFLOG
# ANDROID_NDK_ROOT=/usr/local/lib/android/sdk/ndk/25.0.8775105
# CHROME_BIN=/usr/bin/google-chrome
# GOROOT_1_18_X64=/opt/hostedtoolcache/go/1.18.5/x64
# JOURNAL_STREAM=8:22806
# RUNNER_WORKSPACE=/home/runner/work/mock-project-repo
# LEIN_HOME=/usr/local/lib/lein
# LEIN_JAR=/usr/local/lib/lein/self-installs/leiningen-2.9.10-standalone.jar
# PATH=/home/linuxbrew/.linuxbrew/bin:/home/linuxbrew/.linuxbrew/sbin:/home/runner/.local/bin:/opt/pipx_bin:/home/runner/.cargo/bin:/home/runner/.config/compose
# r/vendor/bin:/usr/local/.ghcup/bin:/home/runner/.dotnet/tools:/snap/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/gam
# es:/snap/bin
# RUNNER_PERFLOG=/home/runner/perflog
# CI=true
# SWIFT_PATH=/usr/share/swift/usr/bin
# ImageOS=ubuntu20
# GOROOT_1_19_X64=/opt/hostedtoolcache/go/1.19.0/x64
# DEBIAN_FRONTEND=noninteractive
# AGENT_TOOLSDIRECTORY=/opt/hostedtoolcache
# OLDPWD=/home/runner/work/mock-project-repo/mock-project-repo
# _=/usr/bin/env
# 
# 