import { expect, Page } from '@playwright/test';
import { API_CALLS_DOMAIN, CONSOLE_URL } from '../playwright.config'
import { User } from '../models/users';

export async function login(page: Page, user: User) {

  // Open login page
  await page.goto(CONSOLE_URL);

  // Fill username field
  await page.fill('[placeholder="username@domain"]', user.username);

  // Click Next button
  await Promise.all([
    page.waitForNavigation(),
    page.click('text=next'),
  ])

  await expect(page).toHaveURL(`https://accounts.${API_CALLS_DOMAIN}/loginname`);

  // Fill password field
  await page.fill('input[name="password"]', user.password);

  // Login
  await Promise.all([
    page.waitForNavigation(),
    page.click('text=next')
  ]);

}