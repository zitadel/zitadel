import {Page} from "@playwright/test";

export async function changePassword(page: Page, loginname: string, password: string) {
    await page.goto('password/change?' + new URLSearchParams({loginName: loginname}));
    await changePasswordScreen(page, loginname, password, password)
    await page.getByRole('button', {name: 'Continue'}).click();
}

async function changePasswordScreen(page: Page, loginname: string, password1: string, password2: string) {
    await page.getByLabel('New Password *').pressSequentially(password1);
    await page.getByLabel('Confirm Password *').pressSequentially(password2);
}