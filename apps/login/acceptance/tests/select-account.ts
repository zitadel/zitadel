import { Page } from "@playwright/test";

export async function selectNewAccount(page: Page) {
  await page.getByRole("link", { name: "Add another account" }).click();
}
