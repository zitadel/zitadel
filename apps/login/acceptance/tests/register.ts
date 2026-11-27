import { Page } from "@playwright/test";
import { emailVerify } from "./email-verify";
import { passkeyRegister } from "./passkey";
import { registerUserScreen } from "./register-screen";
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
  await registerUserScreen(page, firstname, lastname, email, password1, password2);
  await page.getByTestId("submit-button").click();
  await verifyEmail(page, email);
}

// The registration always sets a password, the passkey is added with the signed in session afterwards.
export async function registerWithPasswordAndPasskey(
  page: Page,
  firstname: string,
  lastname: string,
  email: string,
  password: string,
): Promise<string> {
  await registerWithPassword(page, firstname, lastname, email, password, password);

  // wait for projection of user
  await page.waitForTimeout(10000);
  await page.goto("./passkey/set?" + new URLSearchParams({ loginName: email }));
  return passkeyRegister(page);
}

async function verifyEmail(page: Page, email: string) {
  const c = await getCodeFromSink(email);
  await emailVerify(page, c);
}
