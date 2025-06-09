import { test } from "@playwright/test";
import { loginScreenExpect, loginWithPassword } from "./login";

test("admin login", async ({ page }) => {
  await loginWithPassword(page, "zitadel-admin@zitadel.localhost", "Password1!");
  await loginScreenExpect(page, "ZITADEL Admin");
});
