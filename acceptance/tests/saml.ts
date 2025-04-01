import { Page } from "@playwright/test";

export async function startSAML(page: Page) {
  await page.goto("http://localhost:8001/hello");
}
