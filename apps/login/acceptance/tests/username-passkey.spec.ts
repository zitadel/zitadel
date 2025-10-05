import { test } from "./registered.js";
import { loginScreenExpect, loginWithPasskey } from "./login.js";

test("username and passkey login", async ({ registeredUser, page }) => {
  await registeredUser.create({
    ...registeredUser.minimal
  })
  await loginWithPasskey(page, registeredUser.getAuthenticatorId(), registeredUser.username);
  await loginScreenExpect(page, registeredUser.fullName);
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
