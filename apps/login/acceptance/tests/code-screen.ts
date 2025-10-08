import { expect, Page } from "@playwright/test";

const codeTextInput = "code-text-input";

export async function codeScreen(page: Page, code: string) {
  await page.getByTestId(codeTextInput).pressSequentially(code);
}

export async function codeScreenExpect(page: Page, code: string) {
  await expect(page.getByTestId(codeTextInput)).toHaveValue(code);
  await expect(page.getByTestId("error").locator("div")).toContainText("Could not verify OTP code");
}
