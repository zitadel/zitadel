#!/usr/bin/env bash

set -a; source $1; set +a

shift

npx playwright test tests --config ./tests/e2e/playwright.config.ts "$@"
