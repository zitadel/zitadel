import test from "@playwright/test";

test("login with Google IDP", async ({ page }) => {
  // Given a Google IDP is configured on the organization
  // Given the user has Google IDP added as auth method
  // User authenticates with the Google IDP
  // User is redirected back to login
  // User is redirected to the app
});

test("login with Google IDP - error", async ({ page }) => {
  // Given the Google IDP is configured on the organization
  // Given the user has Google IDP added as auth method
  // User is redirected to the Google IDP
  // User authenticates with the Google IDP and gets an error
  // User is redirected back to login
  // An error is shown to the user "Something went wrong"
});

test("login with Google IDP, no user existing - auto register", async ({ page }) => {
  // Given idp Google is configure on the organization as only authencation method
  // Given idp Google is configure with account creation alloweed, and automatic creation enabled
  // Given no user exists yet
  // User is automatically redirected to Google
  // User authenticates in Google
  // User is redirect to ZITADEL login
  // User is created in ZITADEL
  // User is redirected to the app (default redirect url)
});

test("login with Google IDP, no user existing - auto register not possible", async ({ page }) => {
  // Given idp Google is configure on the organization as only authencation method
  // Given idp Google is configure with account creation alloweed, and automatic creation enabled
  // Given no user exists yet
  // User is automatically redirected to Google
  // User authenticates in Google
  // User is redirect to ZITADEL login
  // Because of missing informaiton on the user auto creation is not possible
  // User will see the registration page with pre filled user information
  // User fills missing information
  // User clicks register button
  // User is created in ZITADEL
  // User is redirected to the app (default redirect url)
});

test("login with Google IDP, no user existing - auto register enabled - manual creation disabled, creation not possible", async ({
  page,
}) => {
  // Given idp Google is configure on the organization as only authencation method
  // Given idp Google is configure with account creation not allowed, and automatic creation enabled
  // Given no user exists yet
  // User is automatically redirected to Google
  // User authenticates in Google
  // User is redirect to ZITADEL login
  // Because of missing informaiton on the user auto creation is not possible
  // Error message is shown, that registration of the user was not possible due to missing information
});

test("login with Google IDP, no user linked - auto link", async ({ page }) => {
  // Given idp Google is configure on the organization as only authencation method
  // Given idp Google is configure with account linking allowed, and linking set to existing email
  // Given user with email address user@zitadel.com exists
  // User is automatically redirected to Google
  // User authenticates in Google with user@zitadel.com
  // User is redirect to ZITADEL login
  // User is linked with existing user in ZITADEL
  // User is redirected to the app (default redirect url)
});

test("login with Google IDP, no user linked, linking not possible", async ({ page }) => {
  // Given idp Google is configure on the organization as only authencation method
  // Given idp Google is configure with manually account linking  not allowed, and linking set to existing email
  // Given user with email address user@zitadel.com doesn't exists
  // User is automatically redirected to Google
  // User authenticates in Google with user@zitadel.com
  // User is redirect to ZITADEL login
  // User with email address user@zitadel.com can not be found
  // User will get an error message that account linking wasn't possible
});

test("login with Google IDP, no user linked, linking successful", async ({ page }) => {
  // Given idp Google is configure on the organization as only authencation method
  // Given idp Google is configure with manually account linking allowed, and linking set to existing email
  // Given user with email address user@zitadel.com doesn't exists
  // User is automatically redirected to Google
  // User authenticates in Google with user@zitadel.com
  // User is redirect to ZITADEL login
  // User with email address user@zitadel.com can not be found
  // User is prompted to link the account manually
  // User is redirected to the app (default redirect url)
});
