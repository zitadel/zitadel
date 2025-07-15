import { Page } from "@playwright/test";

const passwordField = "password-text-input";
const passwordConfirmField = "password-confirm-text-input";

export async function registerUserScreenPassword(page: Page, firstname: string, lastname: string, email: string) {
  await registerUserScreen(page, firstname, lastname, email);
  await page.getByTestId("password-radio").click();
}

export async function registerUserScreenPasskey(page: Page, firstname: string, lastname: string, email: string) {
  await registerUserScreen(page, firstname, lastname, email);
  await page.getByTestId("passkey-radio").click();
}

export async function registerPasswordScreen(page: Page, password1: string, password2: string) {
  await page.getByTestId(passwordField).pressSequentially(password1);
  await page.getByTestId(passwordConfirmField).pressSequentially(password2);
}

export async function registerUserScreen(page: Page, firstname: string, lastname: string, email: string) {
  await page.getByTestId("firstname-text-input").pressSequentially(firstname);
  await page.getByTestId("lastname-text-input").pressSequentially(lastname);
  await page.getByTestId("email-text-input").pressSequentially(email);
  await page.getByTestId("privacy-policy-checkbox").check();
  await page.getByTestId("tos-checkbox").check();
}
