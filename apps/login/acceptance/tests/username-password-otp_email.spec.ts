import { test } from "./fixtures.js";
import { test as base } from "@playwright/test";
import { code, codeResend, emailOtpFromMockServer, smsOtpFromMockServer } from "./code.js";
import { codeScreenExpect } from "./code-screen.js";
import { loginScreenExpect, loginWithPassword, loginWithPasswordAndEmailOTP } from "./login.js";

test.skip("DOESN'T WORK: username, password and email otp login, enter code manually", async ({ registeredUser, userService, page }) => {
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User enters the code into the ui
  // User is redirected to the app (default redirect url)
  await registeredUser.create()
  await userService.native.addOTPEmail({
    userId: registeredUser.res?.id,
  })
  await loginWithPasswordAndEmailOTP(page, registeredUser.username, registeredUser.password, registeredUser.email);
  await loginScreenExpect(page, registeredUser.fullName);
});

test("username, password and email otp login, click link in email", async ({ page }) => {
  base.skip();
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User clicks link in the email
  // User is redirected to the app (default redirect url)
});

test.skip("DOESN'T WORK: username, password and email otp login, resend code", async ({ registeredUser, page }) => {
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User clicks resend code
  // User receives a new email with a verification code
  // User enters the new code in the ui
  // User is redirected to the app (default redirect url)
  await registeredUser.create()
  await loginWithPassword(page, registeredUser.username, registeredUser.password);
  await emailOtpFromMockServer(page, registeredUser.username);
  await codeResend(page);
  await emailOtpFromMockServer(page, registeredUser.username);
  await loginScreenExpect(page, registeredUser.fullName);
});

test("username, password and email otp login, wrong code", async ({ registeredUser, userService, page }) => {
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User enters a wrong code
  // Error message - "Invalid code" is shown
  await registeredUser.create()
  // Drain first code?
  await userService.native.addOTPEmail({
    userId: registeredUser.res?.id,
  })
  const c = "wrongcode";
  await loginWithPassword(page, registeredUser.username, registeredUser.password);
  await code(page, c);
  await codeScreenExpect(page, c);
});

test("username, password and email otp login, multiple mfa options", async ({ page }) => {
  base.skip();
  // Given email otp and sms otp is enabled on the organization of the user
  // Given the user has email and sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User clicks button to use sms otp as second factor
  // User receives a sms with a verification code
  // User enters code in ui
  // User is redirected to the app (default redirect url)
});
