import { expect, Page } from "@playwright/test";
import { getCodeFromSink } from "./sink";

const codeField = "code-text-input";
const passwordField = "password-text-input";
const passwordChangeField = "password-change-text-input";
const passwordChangeConfirmField = "password-change-confirm-text-input";
const passwordSetField = "password-set-text-input";
const passwordSetConfirmField = "password-set-confirm-text-input";
const lengthCheck = "length-check";
const symbolCheck = "symbol-check";
const numberCheck = "number-check";
const uppercaseCheck = "uppercase-check";
const lowercaseCheck = "lowercase-check";
const equalCheck = "equal-check";

const matchText = "Matches";
const noMatchText = "Doesn't match";

export async function changePasswordScreen(page: Page, password1: string, password2: string) {
  await page.getByTestId(passwordChangeField).pressSequentially(password1);
  await page.getByTestId(passwordChangeConfirmField).pressSequentially(password2);
}

export async function passwordScreen(page: Page, password: string) {
  await page.getByTestId(passwordField).pressSequentially(password);
}

export async function passwordScreenExpect(page: Page, password: string) {
  await expect(page.getByTestId(passwordField)).toHaveValue(password);
  await expect(page.getByTestId("error").locator("div")).toContainText("Failed to authenticate.");
}

export async function changePasswordScreenExpect(
  page: Page,
  password1: string,
  password2: string,
  length: boolean,
  symbol: boolean,
  number: boolean,
  uppercase: boolean,
  lowercase: boolean,
  equals: boolean,
) {
  await expect(page.getByTestId(passwordChangeField)).toHaveValue(password1);
  await expect(page.getByTestId(passwordChangeConfirmField)).toHaveValue(password2);

  await checkComplexity(page, length, symbol, number, uppercase, lowercase, equals);
}

async function checkComplexity(
  page: Page,
  length: boolean,
  symbol: boolean,
  number: boolean,
  uppercase: boolean,
  lowercase: boolean,
  equals: boolean,
) {
  await checkContent(page, lengthCheck, length);
  await checkContent(page, symbolCheck, symbol);
  await checkContent(page, numberCheck, number);
  await checkContent(page, uppercaseCheck, uppercase);
  await checkContent(page, lowercaseCheck, lowercase);
  await checkContent(page, equalCheck, equals);
}

async function checkContent(page: Page, testid: string, match: boolean) {
  if (match) {
    await expect(page.getByTestId(testid)).toContainText(matchText);
  } else {
    await expect(page.getByTestId(testid)).toContainText(noMatchText);
  }
}

export async function resetPasswordScreen(page: Page, username: string, password1: string, password2: string) {
  const c = await getCodeFromSink(username);
  await page.getByTestId(codeField).pressSequentially(c);
  await page.getByTestId(passwordSetField).pressSequentially(password1);
  await page.getByTestId(passwordSetConfirmField).pressSequentially(password2);
}

export async function resetPasswordScreenExpect(
  page: Page,
  password1: string,
  password2: string,
  length: boolean,
  symbol: boolean,
  number: boolean,
  uppercase: boolean,
  lowercase: boolean,
  equals: boolean,
) {
  await expect(page.getByTestId(passwordSetField)).toHaveValue(password1);
  await expect(page.getByTestId(passwordSetConfirmField)).toHaveValue(password2);

  await checkComplexity(page, length, symbol, number, uppercase, lowercase, equals);
}
