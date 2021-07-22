docker run \
--env CYPRESS_username="$USERNAME" \
--env CYPRESS_password="$USERPW" \
--env CYPRESS_consoleUrl="$CONSOLEURL" \
--env CYPRESS_projectName="newProject"  \
-it -v $PWD:/e2e -w /e2e cypress/included:8.0.0
