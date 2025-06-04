import { faker } from "@faker-js/faker";
import { test as base } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
import { emailVerify, emailVerifyResend } from "./email-verify";
import { emailVerifyScreenExpect } from "./email-verify-screen";
import { loginScreenExpect, loginWithPassword } from "./login";
import { getCodeFromSink } from "./sink";
import { PasswordUser } from "./user";

// Read from ".env" file.
dotenv.config({ path: path.resolve(__dirname, ".env.local") });

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
  // auto-redirect on /verify
  // wait for send of the code
  await page.waitForTimeout(3000);
  const c = await getCodeFromSink(user.getUsername());
  await emailVerify(page, c);
  await loginScreenExpect(page, user.getFullName());
});

test("user email not verified, resend, verify", async ({ user, page }) => {
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  // auto-redirect on /verify
  await emailVerifyResend(page);
  // wait for send of the code
  await page.waitForTimeout(3000);
  const c = await getCodeFromSink(user.getUsername());
  await emailVerify(page, c);
  await loginScreenExpect(page, user.getFullName());
});

test("user email not verified, resend, old code", async ({ user, page }) => {
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  // auto-redirect on /verify
  // wait for send of the code
  await page.waitForTimeout(3000);
  const c = await getCodeFromSink(user.getUsername());
  await emailVerifyResend(page);
  // wait for resend of the code
  await page.waitForTimeout(1000);
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
