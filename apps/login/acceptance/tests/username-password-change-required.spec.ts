import { test } from "./fixtures.js";
import { loginScreenExpect, loginWithPassword } from "./login.js";
import { changePassword } from "./password.js";

test("username and password login, change required", async ({ userCreator, page }) => {
  const changedPw = "ChangedPw1!";
  await userCreator.withPasswordChangeRequired().create();
  await loginWithPassword(page, userCreator.username, userCreator.password);
  await page.waitForTimeout(10000);
  await changePassword(page, changedPw);
  await loginScreenExpect(page, userCreator.fullName);
  await loginWithPassword(page, userCreator.username, changedPw);
  await loginScreenExpect(page, userCreator.fullName);
});
