import { test } from "@playwright/test";

test("username, password and email otp login, enter code manually", async ({ page }) => {
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User enters the code into the ui
  // User is redirected to the app (default redirect url)
});

test("username, password and email otp login, click link in email", async ({ page }) => {
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User clicks link in the email
  // User is redirected to the app (default redirect url)
});

test("username, password and email otp login, resend code", async ({ page }) => {
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User clicks resend code
  // User receives a new email with a verification code
  // User enters the new code in the ui
  // User is redirected to the app (default redirect url)
});

test("username, password and email otp login, wrong code", async ({ page }) => {
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User enters a wrond code
  // Error message - "Invalid code" is shown
});

test("username, password and email otp login, multiple mfa options", async ({ page }) => {
  // Given email otp and sms otp is enabled on the organization of the user
  // Given the user has email and sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User clicks button to use sms otp as second factor
  // User receives an sms with a verification code
  // User enters code in ui
  // User is redirected to the app (default redirect url)
});
