import { expect, Page } from "@playwright/test";

const usernameTextInput = "username-text-input";

export async function loginnameScreen(page: Page, username: string) {
  await page.getByTestId(usernameTextInput).pressSequentially(username);
}

export async function loginnameScreenExpect(page: Page, username: string) {
  await expect(page.getByTestId(usernameTextInput)).toHaveValue(username);
  await expect(page.getByTestId("error").locator("div")).toContainText("User not found in the system");
}
