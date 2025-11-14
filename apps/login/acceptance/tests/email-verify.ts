import { Page } from "@playwright/test";
import { emailVerifyScreen } from "./email-verify-screen";

export async function startEmailVerify(page: Page, loginname: string) {
  await page.goto("./verify");
}

export async function emailVerify(page: Page, code: string) {
  await emailVerifyScreen(page, code);
  await page.getByTestId("submit-button").click();
}

export async function emailVerifyResend(page: Page) {
  await page.getByTestId("resend-button").click();
}
