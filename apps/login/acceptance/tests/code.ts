import { expect, Page } from "@playwright/test";
import { codeScreen } from "./code-screen";
import { getOtpFromSink } from "./sink";

export async function otpFromSink(page: Page, key: string) {
  const c = await getOtpFromSink(key);
  await code(page, c);
}

export async function code(page: Page, code: string) {
  await codeScreen(page, code);
  await page.getByTestId("submit-button").click();
}

export async function codeResend(page: Page) {
  // The button is disabled for the resend cooldown (30s, see RESEND_COOLDOWN_SECONDS)
  // after a code was just sent, so wait for it to enable.
  await expect(page.getByTestId("resend-button")).toBeEnabled({ timeout: 40_000 });
  await page.getByTestId("resend-button").click();
}
