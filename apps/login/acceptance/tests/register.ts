import { Page } from "@playwright/test";
import { emailVerify } from "./email-verify";
import { passkeyRegister } from "./passkey";
import { registerPasswordScreen, registerUserScreenPasskey, registerUserScreenPassword } from "./register-screen";
import { getCodeFromSink } from "./sink";
import { Config } from "./config";

export async function registerWithPassword(
  cfg: Config, 
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
  await verifyEmail(cfg, page, email);
}

export async function registerWithPasskey(cfg: Config, page: Page, firstname: string, lastname: string, email: string): Promise<string> {
  await page.goto("./register");
  await registerUserScreenPasskey(page, firstname, lastname, email);
  await page.getByTestId("submit-button").click();

  // wait for projection of user
  await page.waitForTimeout(10000);
  const authId = await passkeyRegister(page);

  await verifyEmail(cfg, page, email);
  return authId;
}

async function verifyEmail(cfg: Config, page: Page, email: string) {
  const c = await getCodeFromSink(cfg, email);
  await emailVerify(page, c);
}
