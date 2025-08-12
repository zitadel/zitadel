import { faker } from "@faker-js/faker";
import { test as base } from "@playwright/test";
import { emailVerify, emailVerifyResend } from "./email-verify";
import { emailVerifyScreenExpect } from "./email-verify-screen";
import { loginScreenExpect, loginWithPassword } from "./login";
import { getCodeFromSink } from "./sink";
import { PasswordUser } from "./user";
import { Config, ConfigReader } from "./config";

const test = base.extend<{ user: PasswordUser; cfg: Config }>({
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
  cfg: async ({}, use) => {
    await use(new ConfigReader().config);
  }
});

test.skip("FAILS: user email not verified, verify", async ({ user, page, cfg }) => {
  const since = new Date();
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  const c = await getCodeFromSink(cfg, user.getUsername(), since);
  await page.waitForTimeout(10_000);
  await emailVerify(page, c);
  // wait for resend of the code
  await loginScreenExpect(page, user.getFullName());
});

test.skip("FAILS: user email not verified, resend, verify", async ({ user, page, cfg }) => {
  const sinceFirst = new Date();
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  // await for the first code
  const first = await getCodeFromSink(cfg, user.getUsername(), sinceFirst);
  // auto-redirect on /verify
  const sinceSecond = new Date();
  await emailVerifyResend(page);
  const second = await getCodeFromSink(cfg, user.getUsername(), sinceSecond);
  if (first === second) {
    throw new Error("Resent code is the same as the first one, expected a different code.");
  }
  await page.waitForTimeout(10_000);
  await emailVerify(page, second);
  await loginScreenExpect(page, user.getFullName());
});

test("user email not verified, resend, old code", async ({ user, page, cfg }) => {
  const sinceFirst = new Date();
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  const first = await getCodeFromSink(cfg, user.getUsername(), sinceFirst);
  const sinceSecond = new Date();
  await emailVerifyResend(page);
  const second = await getCodeFromSink(cfg, user.getUsername(), sinceSecond);
  await page.waitForTimeout(10_000);
  await emailVerify(page, first);
  await emailVerifyScreenExpect(page, first);
});

test("user email not verified, wrong code", async ({ user, page }) => {
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  // auto-redirect on /verify
  const code = "wrong";
  await emailVerify(page, code);
  await emailVerifyScreenExpect(page, code);
});
