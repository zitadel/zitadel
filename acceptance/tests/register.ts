import {Page} from "@playwright/test";

export async function registerWithPassword(page: Page, firstname: string, lastname: string, email: string, password1: string, password2: string) {
    await page.goto('/register');
    await registerUserScreen(page, firstname, lastname, email)
    await page.getByLabel('Password').click();
    await page.getByRole('button', {name: 'Continue'}).click();
    await registerPasswordScreen(page, password1, password2)
    await page.getByRole('button', {name: 'Continue'}).click();
}

export async function registerWithPasskey(page: Page, firstname: string, lastname: string, email: string) {
    await page.goto('/register');
    await registerUserScreen(page, firstname, lastname, email)
    await page.getByLabel('Passkey').click();
    await page.getByRole('button', {name: 'Continue'}).click();
    await page.getByRole('button', {name: 'Continue'}).click();
}

async function registerUserScreen(page: Page, firstname: string, lastname: string, email: string) {
    await page.getByLabel('First name *').pressSequentially(firstname);
    await page.getByLabel('Last name *').pressSequentially(lastname);
    await page.getByLabel('E-mail *').pressSequentially(email);
    await page.getByRole('checkbox').first().check();
    await page.getByRole('checkbox').nth(1).check();
}

async function registerPasswordScreen(page: Page, password1: string, password2: string) {
    await page.getByLabel('Password *', {exact: true}).fill(password1);
    await page.getByLabel('Confirm Password *').fill(password2);
}