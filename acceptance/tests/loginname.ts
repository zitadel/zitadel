import {Page} from "@playwright/test";

export async function loginnameScreen(page: Page, username: string) {
    await page.getByTestId("username-text-input").pressSequentially(username);
}
