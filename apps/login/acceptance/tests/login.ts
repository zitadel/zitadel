import { expect, Page } from "@playwright/test";
import { code, otpFromSink } from "./code";
import { loginname } from "./loginname";
import { password } from "./password";
import { totp } from "./zitadel";

export async function startLogin(page: Page) {
  await page.goto(`./loginname`);
}

export async function loginWithPassword(page: Page, username: string, pw: string) {
  await startLogin(page);
  await loginname(page, username);
  await password(page, pw);
}

export async function loginWithPasskey(page: Page, authenticatorId: string, username: string) {
  await startLogin(page);
  await loginname(page, username);
  // await passkey(page, authenticatorId);
}

export async function loginScreenExpect(page: Page, fullName: string) {
  await expect(page).toHaveURL(/.*signedin.*/);
  await expect(page.getByRole("heading")).toContainText(fullName);
}

export async function loginWithPasswordAndEmailOTP(page: Page, username: string, password: string, email: string) {
  await loginWithPassword(page, username, password);
  await otpFromSink(page, email);
}

export async function loginWithPasswordAndPhoneOTP(page: Page, username: string, password: string, phone: string) {
  await loginWithPassword(page, username, password);
  await otpFromSink(page, phone);
}

export async function loginWithPasswordAndTOTP(page: Page, username: string, password: string, secret: string) {
  await loginWithPassword(page, username, password);
  await code(page, totp(secret));
}
