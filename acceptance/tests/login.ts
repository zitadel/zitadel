import {Page} from "@playwright/test";

export async function loginWithPassword(page: Page, username: string, password: string) {
    await page.goto("/loginname");
    const loginname = page.getByLabel("Loginname");
    await loginname.pressSequentially(username);
    await loginname.press("Enter");
    const pw = page.getByLabel("Password");
    await pw.pressSequentially(password);
    await pw.press("Enter");
}


export async function loginWithPasskey(page: Page, username: string) {
    await page.goto("/loginname");
    const loginname = page.getByLabel("Loginname");
    await loginname.pressSequentially(username);
    await loginname.press("Enter");
}