import { loginScreenExpect, loginWithPassword, startLogin } from "./login.js";
import { loginname } from "./loginname.js";
import { loginnameScreenExpect } from "./loginname-screen.js";
import { password } from "./password.js";
import { passwordScreenExpect } from "./password-screen.js";
import { test } from "./fixtures.js";

test("username and password login", async ({ registeredUser, page }) => {
  await registeredUser.create()
  await registeredUser.create()
  await loginWithPassword(page, registeredUser.username, registeredUser.getPassword());
  await loginScreenExpect(page, registeredUser.fullName);
});

test("username and password login, unknown username", async ({ page }) => {
  const username = "unknown";
  await startLogin(page);
  await loginname(page, username);
  await loginnameScreenExpect(page, username);
});

test("username and password login, wrong password", async ({ registeredUser, page }) => {
  await registeredUser.create()
  await startLogin(page);
  await loginname(page, registeredUser.username);
  await password(page, "wrong");
  await passwordScreenExpect(page, "wrong");
});

test("username and password login, wrong username, ignore unknown usernames", async ({ }) => {
  test.skip();
  // Given user doesn't exist but ignore unknown usernames setting is set to true
  // Given username password login is enabled on the users organization
  // enter login name
  // enter password
  // redirect to loginname page --> error message username or password wrong
});

test("username and password login, initial password change", async ({ }) => {
  test.skip();
  // Given user is created and has changePassword set to true
  // Given username password login is enabled on the users organization
  // enter login name
  // enter password
  // create new password
});

test("username and password login, reset password hidden", async ({ }) => {
  test.skip();
  // Given the organization has enabled "Password reset hidden" in the login policy
  // Given username password login is enabled on the users organization
  // enter login name
  // password reset link should not be shown on password screen
});

test("username and password login, reset password - enter code manually", async ({ }) => {
  test.skip();
  // Given user has forgotten password and clicks the forgot password button
  // Given username password login is enabled on the users organization
  // enter login name
  // click password forgotten
  // enter code from email
  // user is redirected to app (default redirect url)
});

test("username and password login, reset password - click link", async ({ }) => {
  test.skip();
  // Given user has forgotten password and clicks the forgot password button, and then the link in the email
  // Given username password login is enabled on the users organization
  // enter login name
  // click password forgotten
  // click link in email
  // set new password
  // redirect to app (default redirect url)
});

test("username and password login, reset password, resend code", async ({ }) => {
  test.skip();
  // Given user has forgotten password and clicks the forgot password button and then resend code
  // Given username password login is enabled on the users organization
  // enter login name
  // click password forgotten
  // click resend code
  // enter code from second email
  // user is redirected to app (default redirect url)
});

test("email login enabled", async ({ }) => {
  test.skip();
  // Given user with the username "testuser", email test@zitadel.com and phone number 0711111111 exists
  // Given no other user with the same email address exists
  // enter email address "test@zitadel.com " in login screen
  // user will get to password screen
});

test("email login disabled", async ({ }) => {
  test.skip();
  // Given user with the username "testuser", email test@zitadel.com and phone number 0711111111 exists
  // Given no other user with the same email address exists
  // enter email address "test@zitadel.com" in login screen
  // user will see error message "user not found"
});

test("email login enabled - multiple users", async ({ }) => {
  test.skip();
  // Given user with the username "testuser", email test@zitadel.com and phone number 0711111111 exists
  // Given a second user with the username "testuser2", email test@zitadel.com and phone number 0711111111 exists
  // enter email address "test@zitadel.com" in login screen
  // user will see error message "user not found"
});

test("phone login enabled", async ({ }) => {
  test.skip();
  // Given user with the username "testuser", email test@zitadel.com and phone number 0711111111 exists
  // Given no other user with the same phon number exists
  // enter phone number "0711111111" in login screen
  // user will get to password screen
});

test("phone login disabled", async ({ }) => {
  test.skip();
  // Given user with the username "testuser", email test@zitadel.com and phone number 0711111111 exists
  // Given no other user with the same phone number exists
  // enter phone number "0711111111" in login screen
  // user will see error message "user not found"
});

test("phone login enabled - multiple users", async ({ }) => {
  test.skip();
  // Given user with the username "testuser", email test@zitadel.com and phone number 0711111111 exists
  // Given a second user with the username "testuser2", email test@zitadel.com and phone number 0711111111 exists
  // enter phone number "0711111111" in login screen
  // user will see error message "user not found"
});
