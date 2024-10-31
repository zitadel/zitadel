import {Page} from "@playwright/test";

export async function changePassword(page: Page, loginname: string, password: string) {
    await page.goto('password/change?' + new URLSearchParams({loginName: loginname}));
    await changePasswordScreen(page, password, password)
    await page.getByTestId("submit-button").click();
}

async function changePasswordScreen(page: Page, password1: string, password2: string) {
    await page.getByTestId('password-text-input').pressSequentially(password1);
    await page.getByTestId('password-confirm-text-input').pressSequentially(password2);
}

export async function passwordScreen(page: Page, password: string) {
    await page.getByTestId("password-text-input").pressSequentially(password);
}

export async function registerPasswordScreen(page: Page, password1: string, password2: string) {
    await page.getByTestId('password-text-input').pressSequentially(password1);
    await page.getByTestId('password-confirm-text-input').pressSequentially(password2);
}
