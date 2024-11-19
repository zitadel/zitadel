import { test } from "@playwright/test";

test("username, password and totp login", async ({ page }) => {
  // Given totp is enabled on the organizaiton of the user
  // Given the user has only totp configured as second factor
  // User enters username
  // User enters password
  // Screen for entering the code is shown directly
  // User enters the code into the ui
  // User is redirected to the app (default redirect url)
});

test("username, password and totp otp login, wrong code", async ({ page }) => {
  // Given totp is enabled on the organizaiton of the user
  // Given the user has only totp configured as second factor
  // User enters username
  // User enters password
  // Screen for entering the code is shown directly
  // User enters a wrond code
  // Error message - "Invalid code" is shown
});

test("username, password and totp login, multiple mfa options", async ({ page }) => {
  // Given totp and email otp is enabled on the organizaiton of the user
  // Given the user has totp and email otp configured as second factor
  // User enters username
  // User enters password
  // Screen for entering the code is shown directly
  // Button to switch to email otp is shown
  // User clicks button to use email otp instead
  // User receives an email with a verification code
  // User enters code in ui
  // User is redirected to the app (default redirect url)
});
