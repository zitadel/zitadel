import { test } from "./fixtures.js";
import { codeScreenExpect } from "./code-screen.js";
import { loginScreenExpect, loginWithPassword, loginWithPasswordAndPhoneOTP } from "./login.js";
import { resendCode, verifySMSOTPFromMockServer, verifyOTPCode } from "./code.js";
import { eventualSMSOTP } from "./mock.js";

test("username, password and sms otp login, enter code manually", async ({ userCreator, userService, page }) => {
  // Given sms otp is enabled on the organization of the user
  // Given the user has only sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives a sms with a verification code
  // User enters the code into the ui
  // User is redirected to the app (default redirect url)
  await userCreator.create()
  await userCreator.addSMSOTPFactor();
  await loginWithPassword(page, userCreator.username, userCreator.password);
  await verifySMSOTPFromMockServer(page, userCreator.phone);
  await loginScreenExpect(page, userCreator.fullName);
});

test("username, password and sms otp login, resend code", async ({ userCreator, userService, page }) => {
  // Given sms otp is enabled on the organization of the user
  // Given the user has only sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives a sms with a verification code
  // User clicks resend code
  // User receives a new sms with a verification code
  // User is redirected to the app (default redirect url)
  await userCreator.create()
  await userCreator.addSMSOTPFactor();
  await loginWithPassword(page, userCreator.username, userCreator.password);
  // drain first code
  await eventualSMSOTP(userCreator.phone);
  await resendCode(page);
  await verifySMSOTPFromMockServer(page, userCreator.phone);
  await loginScreenExpect(page, userCreator.fullName);
});

test("username, password and sms otp login, wrong code", async ({ userCreator, userService, page }) => {
  // Given sms otp is enabled on the organization of the user
  // Given the user has only sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives a sms with a verification code
  // User enters a wrong code
  // Error message - "Invalid code" is shown
  await userCreator.create()
  await userCreator.addSMSOTPFactor();
  await loginWithPassword(page, userCreator.username, userCreator.password);
  // await valid code exists
  await eventualSMSOTP(userCreator.phone);
  const c = "wrongcode";
  await verifyOTPCode(page, c);
  await codeScreenExpect(page, c);
});
