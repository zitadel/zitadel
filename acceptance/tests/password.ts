import { Page } from "@playwright/test";
import { changePasswordScreen, passwordScreen } from "./password-screen";

const passwordSubmitButton = "submit-button";

export async function startChangePassword(page: Page, loginname: string) {
  await page.goto("/password/change?" + new URLSearchParams({ loginName: loginname }));
}

export async function changePassword(page: Page, loginname: string, password: string) {
  await startChangePassword(page, loginname);
  await changePasswordScreen(page, password, password);
  await page.getByTestId(passwordSubmitButton).click();
}

export async function password(page: Page, password: string) {
  await passwordScreen(page, password);
  await page.getByTestId(passwordSubmitButton).click();
}
