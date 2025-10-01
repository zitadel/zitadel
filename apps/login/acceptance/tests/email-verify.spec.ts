import { faker } from "@faker-js/faker";
import { test as base } from "@playwright/test";
import { emailVerify, emailVerifyResend } from "./email-verify.js";
import { emailVerifyScreenExpect } from "./email-verify-screen.js";
import { loginScreenExpect, loginWithPassword } from "./login.js";
import { eventualCode } from "./mock.js";
import { PasswordUser } from "./user.js";

const test = base.extend<{ user: PasswordUser }>({
  user: async ({ page }, use) => {
    const user = new PasswordUser({
      email: faker.internet.email(),
      isEmailVerified: false,
      firstName: faker.person.firstName(),
      lastName: faker.person.lastName(),
      organization: "",
      phone: faker.phone.number(),
      isPhoneVerified: false,
      password: "Password1!",
      passwordChangeRequired: false,
    });
    await user.ensure(page);
    await use(user);
    await user.cleanup();
  },
});

test("user email not verified, verify", async ({ user, page }) => {
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  const c = await eventualCode(user.getUsername());
  await emailVerify(page, c);
  // wait for resend of the code
  await page.waitForTimeout(2000);
  await loginScreenExpect(page, user.getFullName());
});

test("user email not verified, resend, verify", async ({ user, page }) => {
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  // auto-redirect on /verify
  await emailVerifyResend(page);
  const c = await eventualCode(user.getUsername());
  // wait for resend of the code
  await page.waitForTimeout(2000);
  await emailVerify(page, c);
  await loginScreenExpect(page, user.getFullName());
});

test("user email not verified, resend, old code", async ({ user, page }) => {
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  const c = await eventualCode(user.getUsername());
  await emailVerifyResend(page);
  // wait for resend of the code
  await page.waitForTimeout(2000);
  await emailVerify(page, c);
  await emailVerifyScreenExpect(page, c);
});

test("user email not verified, wrong code", async ({ user, page }) => {
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  // auto-redirect on /verify
  const code = "wrong";
  await emailVerify(page, code);
  await emailVerifyScreenExpect(page, code);
});
