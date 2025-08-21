#!/bin/bash

if [ "$FAIL_COMMANDS_ON_ERRORS" == "true" ]; then
    set -e
fi

echo
echo
echo
echo -e "THANKS FOR CONTRIBUTING TO ZITADEL 🚀"
echo
echo "Your dev container is configured for fixing login integration tests."
echo "The login is running in a separate container with the same configuration."
echo "It calls the mock-zitadel container which provides a mocked Zitadel gRPC API."
echo
echo "Also the test suite is configured correctly."
echo "For example, run a single test file:"
echo "pnpm cypress run --spec integration/integration/login.cy.ts"
echo
echo "You can also run the test interactively."
echo "However, this is only possible from outside the dev container." 
echo "On your host machine, run:"
echo "cd apps/login"
echo "pnpm cypress open"
echo
echo "If you want to change the login code, you can replace the login container by a hot reloading dev server."
echo "docker stop login-integration"
echo "pnpm turbo dev"
echo "Navigate to the page you want to fix, for example:"
echo "http://localhost:3001/ui/v2/login/verify?userId=221394658884845598&code=abc"
echo "Change some code and reload the page for instant feedback."
echo
echo "When you are done, make sure all integration tests pass:"
echo "pnpm cypress run"
echo

if [ "$FAIL_COMMANDS_ON_ERRORS" != "true" ]; then
    exit 0
fi
