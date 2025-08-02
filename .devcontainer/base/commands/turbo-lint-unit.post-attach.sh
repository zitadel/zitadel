#!/bin/bash

if [ "$FAIL_COMMANDS_ON_ERRORS" == "true" ]; then
    set -e
fi

echo
echo
echo
echo -e "THANKS FOR CONTRIBUTING TO ZITADEL 🚀"
echo
echo "Your dev container is configured for fixing linting and unit tests."
echo "No other services are running alongside this container."
echo
echo "To fix all auto-fixable linting errors, run:"
echo "pnpm turbo lint:fix"
echo
echo "To watch console linting errors, run:"
echo "pnpm turbo watch lint --filter console"
echo
echo "To watch @zitadel/client unit test failures, run:"
echo "pnpm turbo watch test:unit --filter @zitadel/client"
echo
echo "To watch @zitadel/login relevant unit tests and linting failures, run:"
echo "pnpm turbo watch lint test:unit --filter @zitadel/login..."
echo

if [ "$FAIL_COMMANDS_ON_ERRORS" != "true" ]; then
    exit 0
fi
