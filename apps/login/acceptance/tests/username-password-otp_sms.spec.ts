import { test } from "./fixtures.js";
import { code } from "./code.js";
import { codeScreenExpect } from "./code-screen.js";
import { loginScreenExpect, loginWithPassword, loginWithPasswordAndPhoneOTP } from "./login.js";

test.skip("DOESN'T WORK: username, password and sms otp login, enter code manually", async ({ registeredUser, userService, page }) => {
  // Given sms otp is enabled on the organization of the user
  // Given the user has only sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives a sms with a verification code
  // User enters the code into the ui
  // User is redirected to the app (default redirect url)
  await registeredUser.create()
  await userService.native.addOTPSMS({
    userId: registeredUser.res?.id,
  })
  await loginWithPasswordAndPhoneOTP(page, registeredUser.username, registeredUser.password, registeredUser.phone);
  await loginScreenExpect(page, registeredUser.fullName);
});

test.skip("DOESN'T WORK: username, password and sms otp login, resend code", async ({ registeredUser, userService, page }) => {
  // Given sms otp is enabled on the organization of the user
  // Given the user has only sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives a sms with a verification code
  // User clicks resend code
  // User receives a new sms with a verification code
  // User is redirected to the app (default redirect url)
  await registeredUser.create()
  await userService.native.addOTPSMS({
    userId: registeredUser.res?.id,
  })
  await loginWithPasswordAndPhoneOTP(page, registeredUser.username, registeredUser.password, registeredUser.phone);
  await loginScreenExpect(page, registeredUser.fullName);
});

test("username, password and sms otp login, wrong code", async ({ registeredUser, userService, page }) => {
  // Given sms otp is enabled on the organization of the user
  // Given the user has only sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives a sms with a verification code
  // User enters a wrong code
  // Error message - "Invalid code" is shown
  await registeredUser.create()
  await userService.native.addOTPSMS({
    userId: registeredUser.res?.id,
  })
  const c = "wrongcode";
  await loginWithPassword(page, registeredUser.username, registeredUser.password);
  await code(page, c);
  await codeScreenExpect(page, c);
});
