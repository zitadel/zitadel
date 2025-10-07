import { emailVerify, emailVerifyResend } from "./email-verify.js";
import { emailVerifyScreenExpect } from "./email-verify-screen.js";
import { loginScreenExpect, loginWithPassword } from "./login.js";
import { eventualEmailOTP } from "./mock.js";
import { test } from "./fixtures.js";

test("user email not verified, verify", async ({ userCreator, page }) => {
  // Create user with password but unverified email
  await userCreator.create();
  console.log("Created user", userCreator.username);
  await loginWithPassword(page, userCreator.username, userCreator.password);
  const code = await eventualEmailOTP(userCreator.username);
  await emailVerify(page, code);
  // wait for resend of the code
  await page.waitForTimeout(2000);
  await loginScreenExpect(page, userCreator.fullName);
});

test("user email not verified, resend, verify", async ({ userCreator, page }) => {
  await userCreator.create();
  await loginWithPassword(page, userCreator.username, userCreator.password);
  // auto-redirect on /verify
  await emailVerifyResend(page);
  const code = await eventualEmailOTP(userCreator.username);
  // wait for resend of the code
  await page.waitForTimeout(2000);
  await emailVerify(page, code);
  await loginScreenExpect(page, userCreator.fullName);
});

test("user email not verified, resend, old code", async ({ userCreator, page }) => {
  await userCreator.create();
  await loginWithPassword(page, userCreator.username, userCreator.password);
  const c = await eventualEmailOTP(userCreator.username);
  await emailVerifyResend(page);
  // wait for resend of the code
  await page.waitForTimeout(2000);
  await emailVerify(page, c);
  await emailVerifyScreenExpect(page, c);
});

test("user email not verified, wrong code", async ({ userCreator, page }) => {
  await userCreator.create();
  await loginWithPassword(page, userCreator.username, userCreator.password);
  // auto-redirect on /verify
  const code = "wrong";
  await emailVerify(page, code);
  await emailVerifyScreenExpect(page, code);
});
