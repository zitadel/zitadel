import { faker } from "@faker-js/faker";
import { test as base } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
import { loginScreenExpect, loginWithPassword } from "./login";
import { recoveryCode, recoveryCodeScreenExpect, useRecoveryCode } from "./recovery-code-screen";
import { PasswordUserWithRecoveryCodes } from "./user";

// Read from ".env" file.
dotenv.config({ path: path.resolve(__dirname, "../../login/.env.test.local") });

const test = base.extend<{ user: PasswordUserWithRecoveryCodes }>({
  user: async ({ page }, use) => {
    const user = new PasswordUserWithRecoveryCodes({
      email: faker.internet.email(),
      isEmailVerified: true,
      firstName: faker.person.firstName(),
      lastName: faker.person.lastName(),
      organization: "",
      phone: faker.phone.number({ style: "international" }),
      isPhoneVerified: true,
      password: "Password1!",
      passwordChangeRequired: false,
    });

    await user.ensure(page);
    await use(user);
    await user.cleanup();
  },
});

test("username, password and recovery code login", async ({ user, page }) => {
  // Given totp and recovery codes are configured for the user
  // User enters username
  // User enters password
  // Screen for entering the totp code is shown, with "Use recovery code" as the alternative
  // User clicks the alternative and enters a recovery code
  // User is redirected to the app (default redirect url)
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  await useRecoveryCode(page);
  await recoveryCode(page, user.getRecoveryCode(0));
  await loginScreenExpect(page, user.getFullName());
});

test("username, password and recovery code login, wrong code", async ({ user, page }) => {
  // Given totp and recovery codes are configured for the user
  // User enters username
  // User enters password
  // User switches to the recovery code page and enters a wrong code
  // Error message - "The recovery code is invalid or has already been used." is shown
  // NOTE: only one failed attempt so the fixture user does not get locked out
  const c = "wrong-code";
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  await useRecoveryCode(page);
  await recoveryCode(page, c);
  await recoveryCodeScreenExpect(page, c);
});

test("username, password and recovery code login, code is single use", async ({ user, page }) => {
  // Given totp and recovery codes are configured for the user
  // User logs in with a recovery code
  // User logs in a second time and reuses the same recovery code
  // Error message - "The recovery code is invalid or has already been used." is shown
  const c = user.getRecoveryCode(1);

  await loginWithPassword(page, user.getUsername(), user.getPassword());
  await useRecoveryCode(page);
  await recoveryCode(page, c);
  await loginScreenExpect(page, user.getFullName());

  await page.context().clearCookies();

  await loginWithPassword(page, user.getUsername(), user.getPassword());
  await useRecoveryCode(page);
  await recoveryCode(page, c);
  await recoveryCodeScreenExpect(page, c);
});
