import { expect, Page } from "@playwright/test";

const codeTextInput = "code-text-input";

export async function emailVerifyScreen(page: Page, code: string) {
  await page.getByTestId(codeTextInput).pressSequentially(code);
}

export async function emailVerifyScreenExpect(page: Page, code: string) {
  await expect(page.getByTestId(codeTextInput)).toHaveValue(code);
  await expect(page.getByTestId("error").locator("div")).toContainText("Could not verify email");
}
