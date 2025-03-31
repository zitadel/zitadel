import { faker } from "@faker-js/faker";
import { test as base, expect } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
import { loginname } from "./loginname";
import { startOIDC } from "./oidc";
import { password } from "./password";
import { PasswordUser } from "./user";

// Read from ".env" file.
dotenv.config({ path: path.resolve(__dirname, ".env.local") });

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

test("oidc username and password login", async ({ user, page }) => {
  await startOIDC(page);
  await loginname(page, user.getUsername());
  await password(page, user.getPassword());
  await expect(page.locator("pre")).toContainText(user.getUsername());
});
