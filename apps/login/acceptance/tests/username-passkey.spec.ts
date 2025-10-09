import { AuthFactorState } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { test } from "./fixtures.js";
import { loginScreenExpect, loginWithPasskey } from "./login.js";
import { AuthFactors } from "@zitadel/proto/zitadel/user/v2/user_service_pb";

test("username and passkey login", async ({ userRegistrator, page }) => {
  const authId = await userRegistrator.registerWithPasskey();
  await loginWithPasskey(page, authId, userRegistrator.username!);
  await loginScreenExpect(page, userRegistrator.fullName!);
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
