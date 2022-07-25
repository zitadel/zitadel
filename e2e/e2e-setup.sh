#!/bin/bash

set -ex

export projectRoot="."

ENVFILE=$1

if [ -z ${ENVFILE:x} ]; then
    echo "Not sourcing any env file"
else
    set -a; source $ENVFILE; set +a
fi

env

go run ./cmd/e2e-setup/*.go "$@"
