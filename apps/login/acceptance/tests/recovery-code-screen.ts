import { expect, Page } from "@playwright/test";

const codeTextInput = "code-text-input";
const recoveryCodeButton = "recovery-code-button";

export async function useRecoveryCode(page: Page) {
  await page.getByTestId(recoveryCodeButton).click();
}

export async function recoveryCodeScreen(page: Page, code: string) {
  await page.getByTestId(codeTextInput).pressSequentially(code);
}

export async function recoveryCode(page: Page, code: string) {
  await recoveryCodeScreen(page, code);
  await page.getByTestId("submit-button").click();
}

export async function recoveryCodeScreenExpect(page: Page, code: string) {
  await expect(page.getByTestId(codeTextInput)).toHaveValue(code);
  await expect(page.getByTestId("error").locator("div")).toContainText(
    "The recovery code is invalid or has already been used.",
  );
}
