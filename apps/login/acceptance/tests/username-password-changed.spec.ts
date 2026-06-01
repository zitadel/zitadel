import { faker } from "@faker-js/faker";
import { test as base } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
import { loginScreenExpect, loginWithPassword } from "./login";
import { changePassword, startChangePassword } from "./password";
import { changePasswordScreen, changePasswordScreenExpect } from "./password-screen";
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

test("username and password changed login", async ({ user, page }) => {
  const changedPw = "ChangedPw1!";
  await loginWithPassword(page, user.getUsername(), user.getPassword());

  // wait for projection of token
  await page.waitForTimeout(10000);

  await startChangePassword(page, user.getUsername());
  await changePassword(page, changedPw);
  await loginScreenExpect(page, user.getFullName());

  await loginWithPassword(page, user.getUsername(), changedPw);
  await loginScreenExpect(page, user.getFullName());
});

test("password change not with desired complexity", async ({ user, page }) => {
  const changedPw1 = "change";
  const changedPw2 = "chang";
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  await startChangePassword(page, user.getUsername());
  await changePasswordScreen(page, changedPw1, changedPw2);
  await changePasswordScreenExpect(page, changedPw1, changedPw2, false, false, false, false, true, false);
});
