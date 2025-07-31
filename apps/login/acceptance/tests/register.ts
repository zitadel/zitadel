import { Page } from "@playwright/test";
import { emailVerify } from "./email-verify";
import { passkeyRegister } from "./passkey";
import { registerPasswordScreen, registerUserScreenPasskey, registerUserScreenPassword } from "./register-screen";
import { getCodeFromSink } from "./sink";

export async function registerWithPassword(
  page: Page,
  firstname: string,
  lastname: string,
  email: string,
  password1: string,
  password2: string,
) {
  await page.goto("./register");
  await registerUserScreenPassword(page, firstname, lastname, email);
  await page.getByTestId("submit-button").click();
  await registerPasswordScreen(page, password1, password2);
  await page.getByTestId("submit-button").click();
  await verifyEmail(page, email);
}

export async function registerWithPasskey(page: Page, firstname: string, lastname: string, email: string): Promise<string> {
  await page.goto("./register");
  await registerUserScreenPasskey(page, firstname, lastname, email);
  await page.getByTestId("submit-button").click();

  // wait for projection of user
  await page.waitForTimeout(10000);
  const authId = await passkeyRegister(page);

  await verifyEmail(page, email);
  return authId;
}

async function verifyEmail(page: Page, email: string) {
  const c = await getCodeFromSink(email);
  await emailVerify(page, c);
}
