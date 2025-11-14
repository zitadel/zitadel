import { test } from "./fixtures.js";
import { resendCode, verifyEmailOTPFromMockServer, verifyOTPCode } from "./code.js";
import { codeScreenExpect } from "./code-screen.js";
import { loginScreenExpect, loginWithPassword } from "./login.js";
import { eventualEmailOTP } from "./mock.js";

test("username, password and email otp login, enter code manually", async ({ userCreator, page }) => {
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User enters the code into the ui
  // User is redirected to the app (default redirect url)
  await userCreator.create()
  await userCreator.addEmailOTPFactor();
  await loginWithPassword(page, userCreator.username, userCreator.password);
  await verifyEmailOTPFromMockServer(page, userCreator.username);
  await loginScreenExpect(page, userCreator.fullName);
});

test.skip("username, password and email otp login, click link in email", async ({ page }) => {
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User clicks link in the email
  // User is redirected to the app (default redirect url)
});

test("username, password and email otp login, resend code", async ({ userCreator, page }) => {
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User clicks resend code
  // User receives a new email with a verification code
  // User enters the new code in the ui
  // User is redirected to the app (default redirect url)
  await userCreator.create()
  await userCreator.addEmailOTPFactor();
  await loginWithPassword(page, userCreator.username, userCreator.password);
  // drain first code
  await eventualEmailOTP(userCreator.username, "oTP");
  await resendCode(page);
  await verifyEmailOTPFromMockServer(page, userCreator.username);
  await loginScreenExpect(page, userCreator.fullName);
});

test("username, password and email otp login, wrong code", async ({ userCreator, page }) => {
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User enters a wrong code
  // Error message - "Invalid code" is shown
  await userCreator.create()
  await userCreator.addEmailOTPFactor();
  await loginWithPassword(page, userCreator.username, userCreator.password);
  // await valid code exists
  await eventualEmailOTP(userCreator.username, "oTP");
  const c = "wrongcode";
  await verifyOTPCode(page, c);
  await codeScreenExpect(page, c);
});

test.skip("username, password and email otp login, multiple mfa options", async ({ page }) => {
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
