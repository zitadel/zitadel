import { Page } from "@playwright/test";
import { codeScreen } from "./code-screen";
import { getOtpFromSink } from "./sink";
import { Config } from "./config";

export async function otpFromSink(page: Page, key: string, cfg: Config, since: Date) {
  const c = await getOtpFromSink(cfg, key, since);
  await code(page, c);
}

export async function code(page: Page, code: string) {
  await codeScreen(page, code);
  await page.getByTestId("submit-button").click();
}

export async function codeResend(page: Page) {
  await page.getByTestId("resend-button").click();
}
