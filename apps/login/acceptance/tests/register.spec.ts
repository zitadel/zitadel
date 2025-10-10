import { faker } from "@faker-js/faker";
import { test } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
import { loginScreenExpect } from "./login";
import { registerWithPasskey, registerWithPassword } from "./register";
import { removeUserByUsername } from "./zitadel";

// Read from ".env" file.
dotenv.config({ path: path.resolve(__dirname, "../../login/.env.test.local") });

test("register with password", async ({ page }) => {
  const username = faker.internet.email();
  const password = "Password1!";
  const firstname = faker.person.firstName();
  const lastname = faker.person.lastName();

  await registerWithPassword(page, firstname, lastname, username, password, password);
  await loginScreenExpect(page, firstname + " " + lastname);

  // wait for projection of user
  await page.waitForTimeout(10000);
  await removeUserByUsername(username);
});

test("register with passkey", async ({ page }) => {
  const username = faker.internet.email();
  const firstname = faker.person.firstName();
  const lastname = faker.person.lastName();

  await registerWithPasskey(page, firstname, lastname, username);
  await loginScreenExpect(page, firstname + " " + lastname);

  // wait for projection of user
  await page.waitForTimeout(10000);
  await removeUserByUsername(username);
});

test("register with username and password - only password enabled", async ({ page }) => {
  test.skip();
  // Given on the default organization "username and password is allowed" is enabled
  // Given on the default organization "username registeration allowed" is enabled
  // Given on the default organization no idp is configured and enabled
  // Given on the default organization passkey is not enabled
  // Given user doesn't exist
  // Click on button "register new user"
  // User is redirected to registration page
  // Only password is shown as an option - no passkey
  // User enters "firstname", "lastname", "username" and "password"
  // User is redirected to app (default redirect url)
});

test("register with username and password - wrong password not enough characters", async ({ page }) => {
  test.skip();
  // Given on the default organization "username and password is allowed" is enabled
  // Given on the default organization "username registeration allowed" is enabled
  // Given on the default organization no idp is configured and enabled
  // Given on the default organization passkey is not enabled
  // Given password policy is set to 8 characters and must include number, symbol, lower and upper letter
  // Given user doesn't exist
  // Click on button "register new user"
  // User is redirected to registration page
  // Only password is shown as an option - no passkey
  // User enters "firstname", "lastname", "username" and a password thats to short
  // Error is shown "Password doesn't match the policy - it must have at least 8 characters"
});

test("register with username and password - wrong password number missing", async ({ page }) => {
  test.skip();
  // Given on the default organization "username and password is allowed" is enabled
  // Given on the default organization "username registeration allowed" is enabled
  // Given on the default organization no idp is configured and enabled
  // Given on the default organization passkey is not enabled
  // Given password policy is set to 8 characters and must include number, symbol, lower and upper letter
  // Given user doesn't exist
  // Click on button "register new user"
  // User is redirected to registration page
  // Only password is shown as an option - no passkey
  // User enters "firstname", "lastname", "username" and a password without a number
  // Error is shown "Password doesn't match the policy - number missing"
});

test("register with username and password - wrong password upper case missing", async ({ page }) => {
  test.skip();
  // Given on the default organization "username and password is allowed" is enabled
  // Given on the default organization "username registeration allowed" is enabled
  // Given on the default organization no idp is configured and enabled
  // Given on the default organization passkey is not enabled
  // Given password policy is set to 8 characters and must include number, symbol, lower and upper letter
  // Given user doesn't exist
  // Click on button "register new user"
  // User is redirected to registration page
  // Only password is shown as an option - no passkey
  // User enters "firstname", "lastname", "username" and a password without an upper case
  // Error is shown "Password doesn't match the policy - uppercase letter missing"
});

test("register with username and password - wrong password lower case missing", async ({ page }) => {
  test.skip();
  // Given on the default organization "username and password is allowed" is enabled
  // Given on the default organization "username registeration allowed" is enabled
  // Given on the default organization no idp is configured and enabled
  // Given on the default organization passkey is not enabled
  // Given password policy is set to 8 characters and must include number, symbol, lower and upper letter
  // Given user doesn't exist
  // Click on button "register new user"
  // User is redirected to registration page
  // Only password is shown as an option - no passkey
  // User enters "firstname", "lastname", "username" and a password without an lower case
  // Error is shown "Password doesn't match the policy - lowercase letter missing"
});

test("register with username and password - wrong password symboo missing", async ({ page }) => {
  test.skip();
  // Given on the default organization "username and password is allowed" is enabled
  // Given on the default organization "username registeration allowed" is enabled
  // Given on the default organization no idp is configured and enabled
  // Given on the default organization passkey is not enabled
  // Given password policy is set to 8 characters and must include number, symbol, lower and upper letter
  // Given user doesn't exist
  // Click on button "register new user"
  // User is redirected to registration page
  // Only password is shown as an option - no passkey
  // User enters "firstname", "lastname", "username" and a password without an symbol
  // Error is shown "Password doesn't match the policy - symbol missing"
});

test("register with username and password - password and passkey enabled", async ({ page }) => {
  test.skip();
  // Given on the default organization "username and password is allowed" is enabled
  // Given on the default organization "username registeration allowed" is enabled
  // Given on the default organization no idp is configured and enabled
  // Given on the default organization passkey is enabled
  // Given user doesn't exist
  // Click on button "register new user"
  // User is redirected to registration page
  // User enters "firstname", "lastname", "username"
  // Password and passkey are shown as authentication option
  // User clicks password
  // User enters password
  // User is redirected to app (default redirect url)
});

test("register with username and passkey - password and passkey enabled", async ({ page }) => {
  test.skip();
  // Given on the default organization "username and password is allowed" is enabled
  // Given on the default organization "username registeration allowed" is enabled
  // Given on the default organization no idp is configured and enabled
  // Given on the default organization passkey is enabled
  // Given user doesn't exist
  // Click on button "register new user"
  // User is redirected to registration page
  // User enters "firstname", "lastname", "username"
  // Password and passkey are shown as authentication option
  // User clicks passkey
  // Passkey is opened automatically
  // User verifies passkey
  // User is redirected to app (default redirect url)
});

test("register with username and password - registration disabled", async ({ page }) => {
  test.skip();
  // Given on the default organization "username and password is allowed" is enabled
  // Given on the default organization "username registeration allowed" is enabled
  // Given on the default organization no idp is configured and enabled
  // Given user doesn't exist
  // Button "register new user" is not available
});

test("register with username and password - multiple registration options", async ({ page }) => {
  test.skip();
  // Given on the default organization "username and password is allowed" is enabled
  // Given on the default organization "username registeration allowed" is enabled
  // Given on the default organization one idp is configured and enabled
  // Given user doesn't exist
  // Click on button "register new user"
  // User is redirected to registration options
  // Local User and idp button are shown
  // User clicks idp button
  // User enters "firstname", "lastname", "username" and "password"
  // User clicks next
  // User is redirected to app (default redirect url)
});
