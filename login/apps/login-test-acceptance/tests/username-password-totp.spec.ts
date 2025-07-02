import { faker } from "@faker-js/faker";
import { test as base } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
import { code } from "./code";
import { codeScreenExpect } from "./code-screen";
import { loginScreenExpect, loginWithPassword, loginWithPasswordAndTOTP } from "./login";
import { PasswordUserWithTOTP } from "./user";

// Read from ".env" file.
dotenv.config({ path: path.resolve(__dirname, "../../login/.env.test.local") });

const test = base.extend<{ user: PasswordUserWithTOTP; sink: any }>({
  user: async ({ page }, use) => {
    const user = new PasswordUserWithTOTP({
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

test("username, password and totp login", async ({ user, page }) => {
  // Given totp is enabled on the organization of the user
  // Given the user has only totp configured as second factor
  // User enters username
  // User enters password
  // Screen for entering the code is shown directly
  // User enters the code into the ui
  // User is redirected to the app (default redirect url)
  await loginWithPasswordAndTOTP(page, user.getUsername(), user.getPassword(), user.getSecret());
  await loginScreenExpect(page, user.getFullName());
});

test("username, password and totp otp login, wrong code", async ({ user, page }) => {
  // Given totp is enabled on the organization of the user
  // Given the user has only totp configured as second factor
  // User enters username
  // User enters password
  // Screen for entering the code is shown directly
  // User enters a wrond code
  // Error message - "Invalid code" is shown
  const c = "wrongcode";
  await loginWithPassword(page, user.getUsername(), user.getPassword());
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
