import { emailVerify, emailVerifyResend } from "./email-verify.js";
import { emailVerifyScreenExpect } from "./email-verify-screen.js";
import { loginScreenExpect, loginWithPassword } from "./login.js";
import { eventualEmailOTP } from "./mock.js";
import { test } from "./user.js";

test.only("user email not verified, verify", async ({ user, page }) => {
  await user.create();
  await loginWithPassword(page, user.username, user.password);
  // Why does loginWithPassword send a code again?
/*  const code = await eventualEmailOTP(user.getUsername());
  await emailVerify(page, code);
  // wait for resend of the code
  await page.waitForTimeout(2000);*/
  await loginScreenExpect(page, user.password);
});

test("user email not verified, resend, verify", async ({ user, page }) => {
  await user.create();
  await loginWithPassword(page, user.username, user.password);
  // auto-redirect on /verify
  await emailVerifyResend(page);
  const code = await eventualEmailOTP(user.username);
  // wait for resend of the code
  await page.waitForTimeout(2000);
  await emailVerify(page, code);
  await loginScreenExpect(page, user.fullName);
});

test("user email not verified, resend, old code", async ({ user, page }) => {
  await user.create();
  await loginWithPassword(page, user.username, user.password);
  const c = await eventualEmailOTP(user.username);
  await emailVerifyResend(page);
  // wait for resend of the code
  await page.waitForTimeout(2000);
  await emailVerify(page, c);
  await emailVerifyScreenExpect(page, c);
});

test("user email not verified, wrong code", async ({ user, page }) => {
  await user.create();    
  await loginWithPassword(page, user.username, user.password);
  // auto-redirect on /verify
  const code = "wrong";
  await emailVerify(page, code);
  await emailVerifyScreenExpect(page, code);
});
