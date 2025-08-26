import { faker } from "@faker-js/faker";
import { test as base } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
import { code, codeResend, otpFromSink } from "./code";
import { codeScreenExpect } from "./code-screen";
import { loginScreenExpect, loginWithPassword, loginWithPasswordAndEmailOTP } from "./login";
import { OtpType, PasswordUserWithOTP } from "./user";

// Read from ".env" file.
dotenv.config({ path: path.resolve(__dirname, "../../login/.env.test.local") });

const test = base.extend<{ user: PasswordUserWithOTP; sink: any }>({
  user: async ({ page }, use) => {
    const user = new PasswordUserWithOTP({
      email: faker.internet.email(),
      isEmailVerified: true,
      firstName: faker.person.firstName(),
      lastName: faker.person.lastName(),
      organization: "",
      phone: faker.phone.number(),
      isPhoneVerified: false,
      password: "Password1!",
      passwordChangeRequired: false,
      type: OtpType.email,
    });

    await user.ensure(page);
    await use(user);
    await user.cleanup();
  },
});

test.skip("DOESN'T WORK: username, password and email otp login, enter code manually", async ({ user, page }) => {
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User enters the code into the ui
  // User is redirected to the app (default redirect url)
  await loginWithPasswordAndEmailOTP(page, user.getUsername(), user.getPassword(), user.getUsername());
  await loginScreenExpect(page, user.getFullName());
});

test("username, password and email otp login, click link in email", async ({ page }) => {
  base.skip();
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User clicks link in the email
  // User is redirected to the app (default redirect url)
});

test.skip("DOESN'T WORK: username, password and email otp login, resend code", async ({ user, page }) => {
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User clicks resend code
  // User receives a new email with a verification code
  // User enters the new code in the ui
  // User is redirected to the app (default redirect url)
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  await codeResend(page);
  await otpFromSink(page, user.getUsername());
  await loginScreenExpect(page, user.getFullName());
});

test("username, password and email otp login, wrong code", async ({ user, page }) => {
  // Given email otp is enabled on the organization of the user
  // Given the user has only email otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User enters a wrong code
  // Error message - "Invalid code" is shown
  const c = "wrongcode";
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  await code(page, c);
  await codeScreenExpect(page, c);
});

test("username, password and email otp login, multiple mfa options", async ({ page }) => {
  base.skip();
  // Given email otp and sms otp is enabled on the organization of the user
  // Given the user has email and sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives an email with a verification code
  // User clicks button to use sms otp as second factor
  // User receives a sms with a verification code
  // User enters code in ui
  // User is redirected to the app (default redirect url)
});
