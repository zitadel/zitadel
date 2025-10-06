import { test } from "./fixtures.js";
import { loginScreenExpect, loginWithPassword } from "./login.js";
import { changePassword } from "./password.js";
import { CreateUserRequestBuilder } from "./user-registered.js";

test("username and password login, change required", async ({ registeredUser, page }) => {
  const changedPw = "ChangedPw1!";
  const createUserReq = new CreateUserRequestBuilder().withPasswordChangeRequired()
  await registeredUser.create(createUserReq);
  await loginWithPassword(page, registeredUser.username, registeredUser.password);
  await page.waitForTimeout(10000);
  await changePassword(page, changedPw);
  await loginScreenExpect(page, registeredUser.fullName);
  await loginWithPassword(page, registeredUser.username, changedPw);
  await loginScreenExpect(page, registeredUser.fullName);
});
