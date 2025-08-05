#!/bin/bash

if [ "$FAIL_COMMANDS_ON_ERRORS" == "true" ]; then
    set -e
fi

pnpm install --frozen-lockfile --recursive
pnpm turbo lint test:unit

if [ "$FAIL_COMMANDS_ON_ERRORS" != "true" ]; then
    exit 0
fi
