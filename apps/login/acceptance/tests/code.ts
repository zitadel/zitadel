import { Page } from "@playwright/test";
import { codeScreen } from "./code-screen.js";
import { eventualEmailOTP , eventualSMSOTP } from "./mock.js";

export async function emailOtpFromMockServer(page: Page, key: string) {
  const c = await eventualEmailOTP(key);
  await code(page, c);
}

export async function smsOtpFromMockServer(page: Page, key: string) {
  const c = await eventualSMSOTP(key);
  await code(page, c);
}

export async function code(page: Page, code: string) {
  await codeScreen(page, code);
  await page.getByTestId("submit-button").click();
}

export async function codeResend(page: Page) {
  await page.getByTestId("resend-button").click();
}
