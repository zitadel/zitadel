import { emailVerify, emailVerifyResend } from "./email-verify.js";
import { emailVerifyScreenExpect } from "./email-verify-screen.js";
import { loginScreenExpect, loginWithPassword } from "./login.js";
import { eventualEmailOTP } from "./mock.js";
import { test } from "./fixtures.js";

test("user email not verified, verify", async ({ registeredUser, page }) => {
  // Create user with password but unverified email
  await registeredUser.create();
  console.log("Created user", registeredUser.username);
  await loginWithPassword(page, registeredUser.username, registeredUser.password);
  const code = await eventualEmailOTP(registeredUser.username);
  await emailVerify(page, code);
  // wait for resend of the code
  await page.waitForTimeout(2000);
  await loginScreenExpect(page, registeredUser.fullName);
});

test("user email not verified, resend, verify", async ({ registeredUser, page }) => {
  await registeredUser.create();
  await loginWithPassword(page, registeredUser.username, registeredUser.password);
  // auto-redirect on /verify
  await emailVerifyResend(page);
  const code = await eventualEmailOTP(registeredUser.username);
  // wait for resend of the code
  await page.waitForTimeout(2000);
  await emailVerify(page, code);
  await loginScreenExpect(page, registeredUser.fullName);
});

test("user email not verified, resend, old code", async ({ registeredUser, page }) => {
  await registeredUser.create();
  await loginWithPassword(page, registeredUser.username, registeredUser.password);
  const c = await eventualEmailOTP(registeredUser.username);
  await emailVerifyResend(page);
  // wait for resend of the code
  await page.waitForTimeout(2000);
  await emailVerify(page, c);
  await emailVerifyScreenExpect(page, c);
});

test("user email not verified, wrong code", async ({ registeredUser, page }) => {
  await registeredUser.create();
  await loginWithPassword(page, registeredUser.username, registeredUser.password);
  // auto-redirect on /verify
  const code = "wrong";
  await emailVerify(page, code);
  await emailVerifyScreenExpect(page, code);
});
