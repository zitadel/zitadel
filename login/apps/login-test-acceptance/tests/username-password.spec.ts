import { faker } from "@faker-js/faker";
import { test as base } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
import { loginScreenExpect, loginWithPassword, startLogin } from "./login";
import { loginname } from "./loginname";
import { loginnameScreenExpect } from "./loginname-screen";
import { password } from "./password";
import { passwordScreenExpect } from "./password-screen";
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

test("username and password login", async ({ user, page }) => {
  await loginWithPassword(page, user.getUsername(), user.getPassword());
  await loginScreenExpect(page, user.getFullName());
});

test("username and password login, unknown username", async ({ page }) => {
  const username = "unknown";
  await startLogin(page);
  await loginname(page, username);
  await loginnameScreenExpect(page, username);
});

test("username and password login, wrong password", async ({ user, page }) => {
  await startLogin(page);
  await loginname(page, user.getUsername());
  await password(page, "wrong");
  await passwordScreenExpect(page, "wrong");
});

test("username and password login, wrong username, ignore unknown usernames", async ({ user, page }) => {
  test.skip();
  // Given user doesn't exist but ignore unknown usernames setting is set to true
  // Given username password login is enabled on the users organization
  // enter login name
  // enter password
  // redirect to loginname page --> error message username or password wrong
});

test("username and password login, initial password change", async ({ user, page }) => {
  test.skip();
  // Given user is created and has changePassword set to true
  // Given username password login is enabled on the users organization
  // enter login name
  // enter password
  // create new password
});

test("username and password login, reset password hidden", async ({ user, page }) => {
  test.skip();
  // Given the organization has enabled "Password reset hidden" in the login policy
  // Given username password login is enabled on the users organization
  // enter login name
  // password reset link should not be shown on password screen
});

test("username and password login, reset password - enter code manually", async ({ user, page }) => {
  test.skip();
  // Given user has forgotten password and clicks the forgot password button
  // Given username password login is enabled on the users organization
  // enter login name
  // click password forgotten
  // enter code from email
  // user is redirected to app (default redirect url)
});

test("username and password login, reset password - click link", async ({ user, page }) => {
  test.skip();
  // Given user has forgotten password and clicks the forgot password button, and then the link in the email
  // Given username password login is enabled on the users organization
  // enter login name
  // click password forgotten
  // click link in email
  // set new password
  // redirect to app (default redirect url)
});

test("username and password login, reset password, resend code", async ({ user, page }) => {
  test.skip();
  // Given user has forgotten password and clicks the forgot password button and then resend code
  // Given username password login is enabled on the users organization
  // enter login name
  // click password forgotten
  // click resend code
  // enter code from second email
  // user is redirected to app (default redirect url)
});

test("email login enabled", async ({ user, page }) => {
  test.skip();
  // Given user with the username "testuser", email test@zitadel.com and phone number 0711111111 exists
  // Given no other user with the same email address exists
  // enter email address "test@zitadel.com " in login screen
  // user will get to password screen
});

test("email login disabled", async ({ user, page }) => {
  test.skip();
  // Given user with the username "testuser", email test@zitadel.com and phone number 0711111111 exists
  // Given no other user with the same email address exists
  // enter email address "test@zitadel.com" in login screen
  // user will see error message "user not found"
});

test("email login enabled - multiple users", async ({ user, page }) => {
  test.skip();
  // Given user with the username "testuser", email test@zitadel.com and phone number 0711111111 exists
  // Given a second user with the username "testuser2", email test@zitadel.com and phone number 0711111111 exists
  // enter email address "test@zitadel.com" in login screen
  // user will see error message "user not found"
});

test("phone login enabled", async ({ user, page }) => {
  test.skip();
  // Given user with the username "testuser", email test@zitadel.com and phone number 0711111111 exists
  // Given no other user with the same phon number exists
  // enter phone number "0711111111" in login screen
  // user will get to password screen
});

test("phone login disabled", async ({ user, page }) => {
  test.skip();
  // Given user with the username "testuser", email test@zitadel.com and phone number 0711111111 exists
  // Given no other user with the same phone number exists
  // enter phone number "0711111111" in login screen
  // user will see error message "user not found"
});

test("phone login enabled - multiple users", async ({ user, page }) => {
  test.skip();
  // Given user with the username "testuser", email test@zitadel.com and phone number 0711111111 exists
  // Given a second user with the username "testuser2", email test@zitadel.com and phone number 0711111111 exists
  // enter phone number "0711111111" in login screen
  // user will see error message "user not found"
});
