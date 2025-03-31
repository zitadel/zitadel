import { Page } from "@playwright/test";

export async function startOIDC(page: Page) {
  await page.goto("http://localhost:8000/login");
}
