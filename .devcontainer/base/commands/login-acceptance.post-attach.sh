#!/bin/bash

if [ "$FAIL_COMMANDS_ON_ERRORS" == "true" ]; then
    set -e
fi

echo
echo
echo
echo -e "THANKS FOR CONTRIBUTING TO ZITADEL 🚀"
echo
nohup bash -c "pnpm playwright show-report --host 0.0.0.0 &"
echo "View the Playwright report at http://localhost:9323"
echo
echo "Your dev container is configured for fixing login acceptance tests."
echo "The login is running in a separate container with the same configuration."
echo "It calls a local zitadel container with a fully implemented gRPC API."
echo
echo "Also the test suite is configured correctly."
echo "For example, rerun only failed tests:"
echo "pnpm playwright test --last-failed"
echo
echo "You can also run the test interactively."
echo "However, this is only possible from outside the dev container." 
echo "On your host machine, run:"
echo "cd apps/login"
echo "pnpm playwright open"
echo "Also consider using the VSCode extension for Playwright:"
echo "https://playwright.dev/docs/getting-started-vscode"
echo
echo "If you want to change the login code, you can replace the login container by a hot reloading dev server."
echo "docker stop login-acceptance"
echo "pnpm turbo dev"
echo "Navigate to the page you want to fix, for example:"
echo "http://localhost:3000/ui/v2/login/loginname"
echo "Change some code and reload the page for instant feedback."
echo
echo "When you are done, make sure all acceptance tests pass:"
echo "pnpm playwright test"
echo

if [ "$FAIL_COMMANDS_ON_ERRORS" != "true" ]; then
    exit 0
fi
