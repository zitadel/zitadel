import { test } from "@playwright/test";

test("username, password and u2f login", async ({ page }) => {
  test.skip();
  // Given u2f is enabled on the organization of the user
  // Given the user has only u2f configured as second factor
  // User enters username
  // User enters password
  // Popup for u2f is directly opened
  // User verifies u2f
  // User is redirected to the app (default redirect url)
});

test("username, password and u2f login, multiple mfa options", async ({ page }) => {
  test.skip();
  // Given u2f and semailms otp is enabled on the organization of the user
  // Given the user has u2f and email otp configured as second factor
  // User enters username
  // User enters password
  // Popup for u2f is directly opened
  // User aborts u2f verification
  // User clicks button to use email otp as second factor
  // User receives an email with a verification code
  // User enters code in ui
  // User is redirected to the app (default redirect url)
});
