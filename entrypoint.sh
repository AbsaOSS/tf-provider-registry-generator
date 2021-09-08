#!/bin/bash
set -a
set -x
set -o pipefail
set -o errexit

TOKEN=$INPUT_TOKEN
ARTIFACTS_DIR=$INPUT_ARTIFACTS_DIR
REPOSITORY=$INPUT_REPOSITORY
BRANCH=$INPUT_BRANCH
GPG_KEYID=$INPUT_KEYID
GPG_ARMOR=$INPUT_GPG_ASCII_ARMOR
USERNAME=$INPUT_USERNAME
EMAIL=$INPUT_EMAIL
WEB_ROOT=$INPUT_WEBROOT
NAMESPACE=$INPUT_NAMESPACE

OWNER=$(cut -d '/' -f 1 <<< "$GITHUB_REPOSITORY")
if [[ -z "$REPOSITORY" ]]; then
    REPOSITORY=$(cut -d '/' -f 2 <<< "$GITHUB_REPOSITORY")
fi

if [[ -z "$REPO_URL" ]]; then
    REPO_URL="https://x-access-token:${TOKEN}@github.com/${OWNER}/${REPOSITORY}"
fi

if [[ -z "$INPUT_NAMESPACE" ]]; then
  INPUT_NAMESPACE=${REPOSITORY}
fi

if [[ -z "$INPUT_USERNAME" ]]; then
    INPUT__USERNAME="${GITHUB_ACTOR}"
fi

if [[ -z "$INPUT_EMAIL" ]]; then
    INPUT_EMAIL="${GITHUB_ACTOR}@users.noreply.github.com"
fi

/usr/bin/tfreg-golang
