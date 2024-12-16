import { expect, Page } from "@playwright/test";
import { getCodeFromSink } from "./sink";

const codeField = "code-text-input";
const passwordField = "password-text-input";
const passwordConfirmField = "password-confirm-text-input";
const lengthCheck = "length-check";
const symbolCheck = "symbol-check";
const numberCheck = "number-check";
const uppercaseCheck = "uppercase-check";
const lowercaseCheck = "lowercase-check";
const equalCheck = "equal-check";

const matchText = "Matches";
const noMatchText = "Doesn't match";

export async function changePasswordScreen(page: Page, password1: string, password2: string) {
  await page.getByTestId(passwordField).pressSequentially(password1);
  await page.getByTestId(passwordConfirmField).pressSequentially(password2);
}

export async function passwordScreen(page: Page, password: string) {
  await page.getByTestId(passwordField).pressSequentially(password);
}

export async function passwordScreenExpect(page: Page, password: string) {
  await expect(page.getByTestId(passwordField)).toHaveValue(password);
  await expect(page.getByTestId("error").locator("div")).toContainText("Could not verify password");
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
  await expect(page.getByTestId(passwordField)).toHaveValue(password1);
  await expect(page.getByTestId(passwordConfirmField)).toHaveValue(password2);

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
  // wait for send of the code
  await page.waitForTimeout(3000);
  const c = await getCodeFromSink(username);
  await page.getByTestId(codeField).pressSequentially(c);
  await page.getByTestId(passwordField).pressSequentially(password1);
  await page.getByTestId(passwordConfirmField).pressSequentially(password2);
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
  await changePasswordScreenExpect(page, password1, password2, length, symbol, number, uppercase, lowercase, equals);
}
