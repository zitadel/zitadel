import {expect, Page} from "@playwright/test";
import {loginnameScreen} from "./loginname";
import {passwordScreen} from "./password";
import {passkeyScreen} from "./passkey";

export async function loginWithPassword(page: Page, username: string, password: string) {
    await page.goto("/loginname");
    await loginnameScreen(page, username)
    await page.getByTestId("submit-button").click()
    await passwordScreen(page, password)
    await page.getByTestId("submit-button").click()
}

export async function loginWithPasskey(page: Page, username: string) {
    await page.goto("/loginname");
    await loginnameScreen(page, username)
    await page.getByTestId("submit-button").click()
    await passkeyScreen(page)
}

export async function checkLogin(page: Page, fullName: string) {
    await expect(page.getByRole('heading')).toContainText(fullName);
}