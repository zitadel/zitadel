import { expect, Page } from "@playwright/test";

const usernameUserInput = "username-text-input";

export async function loginnameScreen(page: Page, username: string) {
  await page.getByTestId(usernameUserInput).pressSequentially(username);
}

export async function loginnameScreenExpect(page: Page, username: string) {
  await expect(page.getByTestId(usernameUserInput)).toHaveValue(username);
  await expect(page.getByTestId("error").locator("div")).toContainText("Could not find user");
}
