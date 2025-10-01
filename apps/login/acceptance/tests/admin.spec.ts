import { test } from "@playwright/test";
import { loginScreenExpect, loginWithPassword } from "./login.js";

test("admin login", async ({ page }) => {
  await loginWithPassword(page, process.env["ZITADEL_ADMIN_USER"]!, "Password1!");
  await loginScreenExpect(page, "Welcome ZITADEL Admin!");
});
