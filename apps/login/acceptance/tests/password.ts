import { Page } from "@playwright/test";
import { changePasswordScreen, passwordScreen, resetPasswordScreen } from "./password-screen";
import { Config } from "./config";

const passwordSubmitButton = "submit-button";
const passwordResetButton = "reset-button";

export async function startChangePassword(page: Page, loginname: string) {
  await page.goto("./password/change?" + new URLSearchParams({ loginName: loginname }));
}

export async function changePassword(page: Page, password: string) {
  await changePasswordScreen(page, password, password);
  await page.getByTestId(passwordSubmitButton).click();
}

export async function password(page: Page, password: string) {
  await passwordScreen(page, password);
  await page.getByTestId(passwordSubmitButton).click();
}

export async function startResetPassword(page: Page) {
  await page.getByTestId(passwordResetButton).click();
}

export async function resetPassword(cfg: Config, page: Page, username: string, password: string) {
  const codeSince = new Date();
  await startResetPassword(page);
  await resetPasswordScreen(cfg, page, codeSince, username, password, password);
  await page.getByTestId(passwordSubmitButton).click();
}
