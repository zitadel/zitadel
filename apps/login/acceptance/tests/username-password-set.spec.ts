import { test } from "./fixtures.js";
import { loginScreenExpect, loginWithPassword, startLogin } from "./login.js";
import { loginname } from "./loginname.js";
import { resetPassword, startResetPassword } from "./password.js";
import { resetPasswordScreen, resetPasswordScreenExpect } from "./password-screen.js";

test("username and password set login", async ({ registeredUser, page }) => {
  const changedPw = "ChangedPw1!";
  await registeredUser.create();
  await startLogin(page);
  await loginname(page, registeredUser.username);
  await resetPassword(page, registeredUser.username, changedPw);
  await loginScreenExpect(page, registeredUser.fullName);

  await loginWithPassword(page, registeredUser.username, changedPw);
  await loginScreenExpect(page, registeredUser.fullName);
});

test("password set not with desired complexity", async ({ registeredUser, page }) => {
  const changedPw1 = "change";
  const changedPw2 = "chang";
  await registeredUser.create();
  await startLogin(page);
  await loginname(page, registeredUser.username);
  await startResetPassword(page);
  await resetPasswordScreen(page, registeredUser.username, changedPw1, changedPw2);
  await resetPasswordScreenExpect(page, changedPw1, changedPw2, false, false, false, false, true, false);
});
