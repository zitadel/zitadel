import {Page} from "@playwright/test";

export async function registerWithPassword(page: Page, firstname: string, lastname: string, email: string, password1: string, password2: string) {
    await page.goto('/register');
    await registerUserScreenPassword(page, firstname, lastname, email)
    await page.getByTestId('submit-button').click();
    await registerPasswordScreen(page, password1, password2)
    await page.getByTestId('submit-button').click();
}

async function registerUserScreenPassword(page: Page, firstname: string, lastname: string, email: string) {
    await registerUserScreen(page, firstname, lastname, email)
    await page.getByLabel('Password').click();
}

async function registerPasswordScreen(page: Page, password1: string, password2: string) {
    await page.getByTestId('password-text-input').fill(password1);
    await page.getByTestId('password-confirm-text-input').fill(password2);
}

export async function registerWithPasskey(page: Page, firstname: string, lastname: string, email: string) {
    await page.goto('/register');
    await registerUserScreenPasskey(page, firstname, lastname, email)
    await page.getByTestId('submit-button').click();
    await page.getByTestId('submit-button').click();
}

async function registerUserScreenPasskey(page: Page, firstname: string, lastname: string, email: string) {
    await registerUserScreen(page, firstname, lastname, email)
    await page.getByLabel('Passkey').click();
}

async function registerUserScreen(page: Page, firstname: string, lastname: string, email: string) {
    await page.getByTestId('firstname-text-input').pressSequentially(firstname);
    await page.getByTestId('lastname-text-input').pressSequentially(lastname);
    await page.getByTestId('email-text-input').pressSequentially(email);
    await page.getByTestId('privacy-policy-checkbox').check();
    await page.getByTestId('tos-checkbox').check();
}