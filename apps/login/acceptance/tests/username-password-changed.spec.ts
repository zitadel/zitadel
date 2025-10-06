import { test } from "./fixtures.js";
import { loginScreenExpect, loginWithPassword } from "./login.js";
import { changePassword, startChangePassword } from "./password.js";
import { changePasswordScreen, changePasswordScreenExpect } from "./password-screen.js";

test("username and password changed login", async ({ registeredUser, page }) => {
  const changedPw = "ChangedPw1!";
  await registeredUser.create();
  await loginWithPassword(page, registeredUser.username, registeredUser.password);

  // wait for projection of token
  await page.waitForTimeout(10000);

  await startChangePassword(page, registeredUser.username);
  await changePassword(page, changedPw);
  await loginScreenExpect(page, registeredUser.fullName);

  await loginWithPassword(page, registeredUser.username, changedPw);
  await loginScreenExpect(page, registeredUser.fullName);
});

test("password change not with desired complexity", async ({ registeredUser, page }) => {
  const changedPw1 = "change";
  const changedPw2 = "chang";
  await registeredUser.create();
  await loginWithPassword(page, registeredUser.username, registeredUser.password);
  await startChangePassword(page, registeredUser.username);
  await changePasswordScreen(page, changedPw1, changedPw2);
  await changePasswordScreenExpect(page, changedPw1, changedPw2, false, false, false, false, true, false);
});
