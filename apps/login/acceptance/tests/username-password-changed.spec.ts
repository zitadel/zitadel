import { test } from "./fixtures.js";
import { loginScreenExpect, loginWithPassword } from "./login.js";
import { changePassword, startChangePassword } from "./password.js";
import { changePasswordScreen, changePasswordScreenExpect } from "./password-screen.js";

test("username and password changed login", async ({ userCreator, page }) => {
  const changedPw = "ChangedPw1!";
  await userCreator.create();
  await loginWithPassword(page, userCreator.username, userCreator.password);

  // wait for projection of token
  await page.waitForTimeout(10000);

  await startChangePassword(page, userCreator.username);
  await changePassword(page, changedPw);
  await loginScreenExpect(page, userCreator.fullName);

  await loginWithPassword(page, userCreator.username, changedPw);
  await loginScreenExpect(page, userCreator.fullName);
});

test("password change not with desired complexity", async ({ userCreator, page }) => {
  const changedPw1 = "change";
  const changedPw2 = "chang";
  await userCreator.create();
  await loginWithPassword(page, userCreator.username, userCreator.password);
  await startChangePassword(page, userCreator.username);
  await changePasswordScreen(page, changedPw1, changedPw2);
  await changePasswordScreenExpect(page, changedPw1, changedPw2, false, false, false, false, true, false);
});
