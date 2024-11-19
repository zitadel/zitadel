import { test as base } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
import { OtpType, PasswordUserWithOTP } from "./user";

// Read from ".env" file.
dotenv.config({ path: path.resolve(__dirname, ".env.local") });

const test = base.extend<{ user: PasswordUserWithOTP }>({
  user: async ({ page }, use) => {
    const user = new PasswordUserWithOTP({
      email: "otp_sms@example.com",
      firstName: "first",
      lastName: "last",
      password: "Password1!",
      organization: "",
      type: OtpType.sms,
    });

    await user.ensure(page);
    await use(user);
  },
});

test("username, password and totp login", async ({ user, page }) => {
  // Given totp is enabled on the organizaiton of the user
  // Given the user has only totp configured as second factor
  // User enters username
  // User enters password
  // Screen for entering the code is shown directly
  // User enters the code into the ui
  // User is redirected to the app (default redirect url)
});

test("username, password and totp otp login, wrong code", async ({ user, page }) => {
  // Given totp is enabled on the organizaiton of the user
  // Given the user has only totp configured as second factor
  // User enters username
  // User enters password
  // Screen for entering the code is shown directly
  // User enters a wrond code
  // Error message - "Invalid code" is shown
});

test("username, password and totp login, multiple mfa options", async ({ user, page }) => {
  // Given totp and email otp is enabled on the organizaiton of the user
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
