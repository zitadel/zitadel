import { expect, Page } from "@playwright/test";
import { codeScreen } from "./code-screen.js";
import { eventualEmailOTP , eventualSMSOTP } from "./mock.js";

export async function verifyEmailCodeFromMockServer(page: Page, key: string) {
  const c = await eventualEmailOTP(key);
  await verifyCode(page, c);
}

export async function smsCodeFromMockServer(page: Page, key: string) {
  const c = await eventualSMSOTP(key);
  await verifyCode(page, c);
}

export async function verifyTOTPCode(page: Page, code: string) {
  await expect(page.getByRole("heading")).toContainText('Verify 2-Factor');
  await verifyCode(page, code); 
}

async function verifyCode(page: Page, code: string) {
  await codeScreen(page, code);
  await page.getByTestId("submit-button").click();
}

export async function resendCode(page: Page) {
  await page.getByTestId("resend-button").click();
}
