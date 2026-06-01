// Note for all tests, in case Apple doesn't deliver all relevant information per default
// We should add an action in the needed cases

import test from "@playwright/test";

test("login with Apple IDP", async ({ page }) => {
  test.skip();
  // Given an Apple IDP is configured on the organization
  // Given the user has an Apple added as auth method
  // User authenticates with Apple
  // User is redirected back to login
  // User is redirected to the app
});

test("login with Apple IDP - error", async ({ page }) => {
  test.skip();
  // Given an Apple IDP is configured on the organization
  // Given the user has an Apple added as auth method
  // User is redirected to Apple
  // User authenticates with Apple and gets an error
  // User is redirect back to login
  // An error is shown to the user "Something went wrong in Apple Login"
});

test("login with Apple IDP, no user existing - auto register", async ({ page }) => {
  test.skip();
  // Given idp Apple is configure on the organization as only authencation method
  // Given idp Apple is configure with account creation alloweed, and automatic creation enabled
  // Given no user exists yet
  // User is automatically redirected to Apple
  // User authenticates in Apple
  // User is redirect to ZITADEL login
  // User is created in ZITADEL
  // User is redirected to the app (default redirect url)
});

test("login with Apple IDP, no user existing - auto register not possible", async ({ page }) => {
  test.skip();
  // Given idp Apple is configure on the organization as only authencation method
  // Given idp Apple is configure with account creation alloweed, and automatic creation enabled
  // Given no user exists yet
  // User is automatically redirected to Apple
  // User authenticates in Apple
  // User is redirect to ZITADEL login
  // Because of missing informaiton on the user auto creation is not possible
  // User will see the registration page with pre filled user information
  // User fills missing information
  // User clicks register button
  // User is created in ZITADEL
  // User is redirected to the app (default redirect url)
});

test("login with Apple IDP, no user existing - auto register enabled - manual creation disabled, creation not possible", async ({
  page,
}) => {
  test.skip();
  // Given idp Apple is configure on the organization as only authencation method
  // Given idp Apple is configure with account creation not allowed, and automatic creation enabled
  // Given no user exists yet
  // User is automatically redirected to Apple
  // User authenticates in Apple
  // User is redirect to ZITADEL login
  // Because of missing informaiton on the user auto creation is not possible
  // Error message is shown, that registration of the user was not possible due to missing information
});

test("login with Apple IDP, no user linked - auto link", async ({ page }) => {
  test.skip();
  // Given idp Apple is configure on the organization as only authencation method
  // Given idp Apple is configure with account linking allowed, and linking set to existing email
  // Given user with email address user@zitadel.com exists
  // User is automatically redirected to Apple
  // User authenticates in Apple with user@zitadel.com
  // User is redirect to ZITADEL login
  // User is linked with existing user in ZITADEL
  // User is redirected to the app (default redirect url)
});

test("login with Apple IDP, no user linked, linking not possible", async ({ page }) => {
  test.skip();
  // Given idp Apple is configure on the organization as only authencation method
  // Given idp Apple is configure with manually account linking  not allowed, and linking set to existing email
  // Given user with email address user@zitadel.com doesn't exists
  // User is automatically redirected to Apple
  // User authenticates in Apple with user@zitadel.com
  // User is redirect to ZITADEL login
  // User with email address user@zitadel.com can not be found
  // User will get an error message that account linking wasn't possible
});

test("login with Apple IDP, no user linked, user link successful", async ({ page }) => {
  test.skip();
  // Given idp Apple is configure on the organization as only authencation method
  // Given idp Apple is configure with manually account linking allowed, and linking set to existing email
  // Given user with email address user@zitadel.com doesn't exists
  // User is automatically redirected to Apple
  // User authenticates in Apple with user@zitadel.com
  // User is redirect to ZITADEL login
  // User with email address user@zitadel.com can not be found
  // User is prompted to link the account manually
  // User is redirected to the app (default redirect url)
});
