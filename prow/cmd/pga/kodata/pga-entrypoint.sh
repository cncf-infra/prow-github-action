#!/bin/sh
#
# Launch custom (no k8s infra) hook build as a GH Action to run a prow plugin
#
# DEBUG
set -x
set -e

env | grep INPUT_

PROW_CONFIGFILE=$HOME/config.yaml
PLUGIN_CONFIGFILE=$HOME/plugins.yaml
HMAC_FILE=$HOME/hmac
GITHUB_TOKENFILE=$HOME/github_token

# Goose plugin https://github.com/kubernetes/test-infra/blob/master/prow/plugins/goose/goose.go#L58
UNSPLASH_KEYFILE=/etc/unsplash-api/honk.txt

if [ "${INPUT_PROW_CONFIG}" != "" ]; then
    echo "${INPUT_PROW_CONFIG}" > "${PROW_CONFIGFILE}"
else
    cp "/var/run/ko/config.yaml" "${PROW_CONFIGFILE}"
fi

if [ "${INPUT_PLUGIN_CONFIG}" != "" ]; then
    echo "${INPUT_PLUGIN_CONFIG}" > "${PLUGIN_CONFIGFILE}"
else
    cp "/var/run/ko/plugins.yaml" "${PLUGIN_CONFIGFILE}"
fi

if [ "${INPUT_HMAC}" != "" ]; then
    echo "${INPUT_HMAC}" > "${HMAC_FILE}"
else
    cp "/var/run/ko/hmac" "${HMAC_FILE}"
fi

echo "${GITHUB_TOKEN}" > "${GITHUB_TOKENFILE}"
echo "${UNSPLASH_ACCESS_KEY}" > "${UNSPLASH_KEYFILE}"


set +e
apk add git

/ko-app/hook \
    --config-path "${PROW_CONFIGFILE}" \
    --plugin-config "${PLUGIN_CONFIGFILE}" \
    --hmac-secret-file "${HMAC_FILE}" \
    --github-token-path "${GITHUB_TOKENFILE}" \
    --dry-run=false

sleep 9999
