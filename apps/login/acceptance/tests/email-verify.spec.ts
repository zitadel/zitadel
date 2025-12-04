import { emailVerify, emailVerifyResend } from "./email-verify.js";
import { emailVerifyFailedScreenExpect } from "./email-verify-screen.js";
import { loginScreenExpect, loginWithPassword } from "./login.js";
import { eventualEmailOTP } from "./mock.js";
import { test } from "./fixtures.js";

test("user email not verified, verify", async ({ userCreator, page }) => {
  // Create user with password but unverified email
  await userCreator.withEmailUnverified().create();
  console.log("Created user", userCreator.username);
  await loginWithPassword(page, userCreator.username, userCreator.password);
  const code = await eventualEmailOTP(userCreator.username);
  await emailVerify(page, code);
  await loginScreenExpect(page, userCreator.fullName);
});

test("user email not verified, resend, verify", async ({ userCreator, page }) => {
  await userCreator.withEmailUnverified().create();
  await loginWithPassword(page, userCreator.username, userCreator.password);
  // Drain first code
  const _ = await eventualEmailOTP(userCreator.username);
  await emailVerifyResend(page);
  const secondCode = await eventualEmailOTP(userCreator.username);
  await emailVerify(page, secondCode);
  await loginScreenExpect(page, userCreator.fullName);
});

test("user email not verified, resend, old code", async ({ userCreator, page }) => {
  await userCreator.withEmailUnverified().create();
  await loginWithPassword(page, userCreator.username, userCreator.password);
  const firstCode = await eventualEmailOTP(userCreator.username);
  await emailVerifyResend(page);
  // Await second code
  const _ = await eventualEmailOTP(userCreator.username);
  await emailVerify(page, firstCode);
  await emailVerifyFailedScreenExpect(page, firstCode);
});

test("user email not verified, wrong code", async ({ userCreator, page }) => {
  await userCreator.withEmailUnverified().create();
  await loginWithPassword(page, userCreator.username, userCreator.password);
  // auto-redirect on /verify
  const code = "wrong";
  await emailVerify(page, code);
  await emailVerifyFailedScreenExpect(page, code);
});
