import { Page } from "@playwright/test";
import { loginnameScreen } from "./loginname-screen";

export async function loginname(page: Page, username: string) {
  await loginnameScreen(page, username);
  await page.getByTestId("submit-button").click();
}
