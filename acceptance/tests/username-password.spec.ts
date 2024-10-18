import { test } from "@playwright/test";

test("username and password", async ({ page }) => {
  await page.goto("/");
  const loginname = page.getByLabel("Loginname");
  await loginname.pressSequentially("zitadel-admin@zitadel.localhost");
  await loginname.press("Enter");
  const password = page.getByLabel("Password");
  await password.pressSequentially("Password1!");
  await password.press("Enter");
  await page.getByRole("heading", {name: "Welcome ZITADEL Admin!"}).click();
});
