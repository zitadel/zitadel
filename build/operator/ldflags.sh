#!/bin/bash

set -e

VERSION=${1}
if [ "${VERSION}" == "" ]; then
  VERSION="$(git rev-parse --abbrev-ref HEAD | sed -e 's/heads\///')"
fi

echo -n "-extldflags -static -X main.Version=${VERSION} -X main.githubClientID=${GITHUBOAUTHCLIENTID} -X main.githubClientSecret=${GITHUBOAUTHCLIENTSECRET}"
