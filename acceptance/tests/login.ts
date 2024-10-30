import {Page} from "@playwright/test";

export async function loginWithPassword(page: Page, username: string, password: string) {
    await page.goto("/loginname");
    await loginnameScreen(page, username)
    await page.getByTestId("submit-button").click()
    await passwordScreen(page, password)
    await page.getByTestId("submit-button").click()
}

export async function loginnameScreen(page: Page, username: string) {
    await page.getByTestId("username-text-input").pressSequentially(username);
}

export async function passwordScreen(page: Page, password: string) {
    await page.getByTestId("password-text-input").pressSequentially(password);
}

export async function loginWithPasskey(page: Page, username: string) {
    await page.goto("/loginname");
    await loginnameScreen(page, username)
    await page.getByTestId("submit-button").click()
    await page.getByTestId("submit-button").click()
}