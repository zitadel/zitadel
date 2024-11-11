import { test } from "@playwright/test";

test("login is accessible", async ({ page }) => {
  await page.goto("http://localhost:3000/");
  await page.getByRole("heading", { name: "Welcome back!" }).isVisible();
});
