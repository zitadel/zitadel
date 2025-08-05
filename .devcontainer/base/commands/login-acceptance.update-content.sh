#!/bin/bash

if [ "$FAIL_COMMANDS_ON_ERRORS" == "true" ]; then
    echo "Running in fail-on-errors mode" 
    set -e
fi

pnpm install --frozen-lockfile \
    --filter @zitadel/login \
    --filter @zitadel/client \
    --filter @zitadel/proto  \
    --filter zitadel-monorepo
pnpm exec playwright install --with-deps
PLAYWRIGHT_HTML_OPEN=never pnpm test:acceptance:login

if [ "$FAIL_COMMANDS_ON_ERRORS" != "true" ]; then
    exit 0
fi
