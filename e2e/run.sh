#!/bin/bash

set -e

DO_BUILD=1
DO_DEPLOY=1
DO_TEST=1

while getopts 'bdt:' OPTION; do
  case "$OPTION" in
    b)
      echo "skipping build"
      DO_BUILD=0
      ;;
    d)
      echo "skipping deployment"
      DO_DEPLOY=0
      ;;
    s)
      echo "skipping tests"
      DO_TEST=0
      ;;
    ?)
      echo "script usage: $(basename \$0) [-b] [-d] [-t]" >&2
      echo "-b   skip build"
      echo "-d   skip deployment"
      echo "-t   skip tests"
      exit 1
      ;;
  esac
done
shift "$(($OPTIND -1))"

if [ "$DO_BUILD" -eq "1" ]; then
    # build the zitadel binary
    goreleaser build --snapshot --single-target --rm-dist
fi

# extract some metadata for building and tagging the docker image
function extract_metadata ()
{
    cat .artifacts/goreleaser/$1 | jq -r $2
}
BUILD_DATE="$(extract_metadata metadata.json '.date')"

# Use simple local date
BUILD_DATE="${BUILD_DATE%.*}"
# Replace colons and plus signs
export BUILD_DATE="${BUILD_DATE//:/_}"


if [ "$DO_BUILD" -eq "1" ]; then
    BUILD_PATH="$(dirname $(extract_metadata artifacts.json '.[0].path'))"
    BUILD_VERSION="$(extract_metadata metadata.json '.version')"

    # build the docker image
    docker build --file ./build/Dockerfile --tag zitadel:latest --tag zitadel:$BUILD_VERSION --tag zitadel:$BUILD_DATE $BUILD_PATH
fi

if [ "$DO_DEPLOY" -eq "1" ]; then
    # run cockroach and zitadel
    docker compose --file ./docs/docs/guides/installation/run/docker-compose.yaml --file ./e2e/docker-compose-overwrite.yaml up --detach
fi
