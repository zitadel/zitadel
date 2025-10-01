import { Page } from "@playwright/test";
import { codeScreen } from "./code-screen.js";
import { eventualOtp } from "./mock.js";

export async function otpFromSink(page: Page, key: string) {
  const c = await eventualOtp(key);
  await code(page, c);
}

export async function code(page: Page, code: string) {
  await codeScreen(page, code);
  await page.getByTestId("submit-button").click();
}

export async function codeResend(page: Page) {
  await page.getByTestId("resend-button").click();
}
