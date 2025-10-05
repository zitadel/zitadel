import { faker } from "@faker-js/faker";
import { test as base } from "@playwright/test";
import { loginScreenExpect, loginWithPassword } from "./login.js";
import { changePassword } from "./password.js";
import { PasswordUser } from "./registered.js";

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
      passwordChangeRequired: true,
    });
    await user.ensure(page);
    await use(user);
    await user.cleanup();
  },
});

test("username and password login, change required", async ({ user, page }) => {
  const changedPw = "ChangedPw1!";

  await loginWithPassword(page, user.getUsername(), user.getPassword());
  await page.waitForTimeout(10000);
  await changePassword(page, changedPw);
  await loginScreenExpect(page, user.getFullName());

  await loginWithPassword(page, user.getUsername(), changedPw);
  await loginScreenExpect(page, user.getFullName());
});
