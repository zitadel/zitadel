import { expect, Page } from "@playwright/test";
import { emailVerifyScreen } from "./email-verify-screen";

export async function startEmailVerify(page: Page, loginname: string) {
  await page.goto("./verify");
}

export async function emailVerify(page: Page, code: string) {
  await emailVerifyScreen(page, code);
  await page.getByTestId("submit-button").click();
}

export async function emailVerifyResend(page: Page) {
  // The button is disabled for the resend cooldown (30s, see RESEND_COOLDOWN_SECONDS)
  // after a code was just sent, e.g. right after login, so wait for it to enable.
  await expect(page.getByTestId("resend-button")).toBeEnabled({ timeout: 40_000 });
  await page.getByTestId("resend-button").click();
}
