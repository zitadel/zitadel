import { Page } from "@playwright/test";
import { codeScreen } from "./code-screen";
import { getOtpFromSink } from "./sink";

export async function otpFromSink(page: Page, key: string) {
  // wait for send of the code
  await page.waitForTimeout(10000);
  const c = await getOtpFromSink(key);
  await code(page, c);
}

export async function code(page: Page, code: string) {
  await codeScreen(page, code);
  await page.getByTestId("submit-button").click();
}

export async function codeResend(page: Page) {
  await page.getByTestId("resend-button").click();
}
