import { faker } from "@faker-js/faker";
import { test as base } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
import { loginScreenExpect, loginWithPasskey } from "./login";
import { PasskeyUser } from "./user";

// Read from ".env" file.
dotenv.config({ path: path.resolve(__dirname, "../../login/.env.test.local") });

const test = base.extend<{ user: PasskeyUser }>({
  user: async ({ page }, use) => {
    const user = new PasskeyUser({
      email: faker.internet.email(),
      isEmailVerified: true,
      firstName: faker.person.firstName(),
      lastName: faker.person.lastName(),
      organization: "",
      phone: faker.phone.number(),
      isPhoneVerified: false,
    });
    await user.ensure(page);
    await use(user);
    await user.cleanup();
  },
});

test("username and passkey login", async ({ user, page }) => {
  await loginWithPasskey(page, user.getAuthenticatorId(), user.getUsername());
  await loginScreenExpect(page, user.getFullName());
});

test("username and passkey login, multiple auth methods", async ({ page }) => {
  test.skip();
  // Given passkey and password is enabled on the organization of the user
  // Given the user has password and passkey registered
  // enter username
  // passkey popup is directly shown
  // user aborts passkey authentication
  // user switches to password authentication
  // user enters password
  // user is redirected to app
});
