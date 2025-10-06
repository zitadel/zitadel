import { faker } from "@faker-js/faker";
import { code } from "./code.js";
import { codeScreenExpect } from "./code-screen.js";
import { loginScreenExpect, loginWithPassword } from "./login.js";
import { test } from "./fixtures.js";

test("username, password and totp login", async ({ registeredUser, userService, page }) => {
  // Given totp is enabled on the organization of the user
  // Given the user has only totp configured as second factor
  // User enters username
  // User enters password
  // Screen for entering the code is shown directly
  // User enters the code into the ui
  // User is redirected to the app (default redirect url)
  const user = await registeredUser.create()
  const secret = await userService.addTOTP(user.id);
  await loginWithPassword(page, registeredUser.username, registeredUser.password);
  await code(page, userService.totp(secret));
  await loginScreenExpect(page, registeredUser.fullName);
});

test("username, password and totp otp login, wrong code", async ({ registeredUser, page }) => {
  // Given totp is enabled on the organization of the user
  // Given the user has only totp configured as second factor
  // User enters username
  // User enters password
  // Screen for entering the code is shown directly
  // User enters a wrong code
  // Error message - "Invalid code" is shown
  await registeredUser.create()
  const c = "wrongcode";
  await loginWithPassword(page, registeredUser.username, registeredUser.password);
  await code(page, c);
  await codeScreenExpect(page, c);
});

test("username, password and totp login, multiple mfa options", async ({ page }) => {
  test.skip();
  // Given totp and email otp is enabled on the organization of the user
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
