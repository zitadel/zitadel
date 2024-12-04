import { Page } from "@playwright/test";
import { passkeyRegister } from "./passkey";
import { registerPasswordScreen, registerUserScreenPasskey, registerUserScreenPassword } from "./register-screen";

export async function registerWithPassword(
  page: Page,
  firstname: string,
  lastname: string,
  email: string,
  password1: string,
  password2: string,
) {
  await page.goto("/register");
  await registerUserScreenPassword(page, firstname, lastname, email);
  await page.getByTestId("submit-button").click();
  await registerPasswordScreen(page, password1, password2);
  await page.getByTestId("submit-button").click();
}

export async function registerWithPasskey(page: Page, firstname: string, lastname: string, email: string): Promise<string> {
  await page.goto("/register");
  await registerUserScreenPasskey(page, firstname, lastname, email);
  await page.getByTestId("submit-button").click();

  // wait for projection of user
  await page.waitForTimeout(2000);

  return await passkeyRegister(page);
}
