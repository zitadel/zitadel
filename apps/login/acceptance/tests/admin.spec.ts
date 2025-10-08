import { test } from "@playwright/test";
import { loginScreenExpect, loginWithPassword } from "./login";

test("admin login", async ({ page }) => {
  await loginWithPassword(page, process.env["ZITADEL_ADMIN_USER"], "Password1!");
  await loginScreenExpect(page, "ZITADEL Admin");
});
