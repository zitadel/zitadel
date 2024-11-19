import { test } from "@playwright/test";

test("username, password and sms otp login", async ({ page }) => {
  // Given sms otp is enabled on the organization of the user
  // Given the user has only sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives a sms with a verification code
  // User enters the code into the ui
  // User is redirected to the app (default redirect url)
});

test("username, password and sms otp login, resend code", async ({ page }) => {
  // Given sms otp is enabled on the organization of the user
  // Given the user has only sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives a sms with a verification code
  // User clicks resend code
  // User receives a new sms with a verification code
  // User is redirected to the app (default redirect url)
});

test("username, password and sms otp login, wrong code", async ({ page }) => {
  // Given sms otp is enabled on the organization of the user
  // Given the user has only sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives a sms with a verification code
  // User enters a wrong code
  // Error message - "Invalid code" is shown
});
