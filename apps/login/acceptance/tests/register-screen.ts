import { Page } from "@playwright/test";

export async function registerUserScreen(
  page: Page,
  firstname: string,
  lastname: string,
  email: string,
  password1: string,
  password2: string,
) {
  await page.getByTestId("firstname-text-input").pressSequentially(firstname);
  await page.getByTestId("lastname-text-input").pressSequentially(lastname);
  await page.getByTestId("email-text-input").pressSequentially(email);
  await page.getByTestId("password-text-input").pressSequentially(password1);
  await page.getByTestId("password-confirm-text-input").pressSequentially(password2);
  await page.getByTestId("privacy-policy-checkbox").check();
  await page.getByTestId("tos-checkbox").check();
}
