import { faker } from "@faker-js/faker";
import { test as base } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
import { code } from "./code";
import { codeScreenExpect } from "./code-screen";
import { loginScreenExpect, loginWithPassword, loginWithPasswordAndPhoneOTP } from "./login";
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
      phone: faker.phone.number({ style: "international" }),
      isPhoneVerified: true,
      password: "Password1!",
      passwordChangeRequired: false,
      type: OtpType.sms,
    });

    await user.ensure(page);
    await use(user);
    await user.cleanup();
  },
});

test.skip("DOESN'T WORK: username, password and sms otp login, enter code manually", async ({ user, page }) => {
  // Given sms otp is enabled on the organization of the user
  // Given the user has only sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives a sms with a verification code
  // User enters the code into the ui
  // User is redirected to the app (default redirect url)
  await loginWithPasswordAndPhoneOTP(page, user.getUsername(), user.getPassword(), user.getPhone());
  await loginScreenExpect(page, user.getFullName());
});

test.skip("DOESN'T WORK: username, password and sms otp login, resend code", async ({ user, page }) => {
  // Given sms otp is enabled on the organization of the user
  // Given the user has only sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives a sms with a verification code
  // User clicks resend code
  // User receives a new sms with a verification code
  // User is redirected to the app (default redirect url)
  await loginWithPasswordAndPhoneOTP(page, user.getUsername(), user.getPassword(), user.getPhone());
  await loginScreenExpect(page, user.getFullName());
});

test("username, password and sms otp login, wrong code", async ({ user, page }) => {
  // Given sms otp is enabled on the organization of the user
  // Given the user has only sms otp configured as second factor
  // User enters username
  // User enters password
  // User receives a sms with a verification code
  // User enters a wrong code
  // Error message - "Invalid code" is shown
  const c = "wrongcode";
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  await code(page, c);
  await codeScreenExpect(page, c);
});
