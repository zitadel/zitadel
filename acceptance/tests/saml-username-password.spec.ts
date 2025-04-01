import { faker } from "@faker-js/faker";
import { test as base } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
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

test("saml username and password login", async ({ user, page }) => {
  //TODO commented out because of prefetching while calling ACS of SAML SP
  /*
  await startSAML(page)
  await selectNewAccount(page)
  await loginname(page, user.getUsername());
  await password(page, user.getPassword());
  await expect(page.locator('html')).toContainText(user.getUsername());
   */
});
