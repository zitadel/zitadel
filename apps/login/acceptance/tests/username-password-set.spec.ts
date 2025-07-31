import { faker } from "@faker-js/faker";
import { test as base } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
import { loginScreenExpect, loginWithPassword, startLogin } from "./login";
import { loginname } from "./loginname";
import { resetPassword, startResetPassword } from "./password";
import { resetPasswordScreen, resetPasswordScreenExpect } from "./password-screen";
import { PasswordUser } from "./user";

// Read from ".env" file.
dotenv.config({ path: path.resolve(__dirname, "../../login/.env.test.local") });

const test = base.extend<{ user: PasswordUser }>({
  user: async ({ page }, use) => {
    const user = new PasswordUser({
      email: faker.internet.email(),
      isEmailVerified: true,
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

test("username and password set login", async ({ user, page }) => {
  const changedPw = "ChangedPw1!";
  await startLogin(page);
  await loginname(page, user.getUsername());
  await resetPassword(page, user.getUsername(), changedPw);
  await loginScreenExpect(page, user.getFullName());

  await loginWithPassword(page, user.getUsername(), changedPw);
  await loginScreenExpect(page, user.getFullName());
});

test("password set not with desired complexity", async ({ user, page }) => {
  const changedPw1 = "change";
  const changedPw2 = "chang";
  await startLogin(page);
  await loginname(page, user.getUsername());
  await startResetPassword(page);
  await resetPasswordScreen(page, user.getUsername(), changedPw1, changedPw2);
  await resetPasswordScreenExpect(page, changedPw1, changedPw2, false, false, false, false, true, false);
});
