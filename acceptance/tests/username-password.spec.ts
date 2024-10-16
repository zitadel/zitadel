import { test, expect } from '@playwright/test';

test('get started link', async ({ page }) => {
  await page.goto('http://localhost:8080/');

  await page.getByRole('heading', { name: 'Welcome back!' }).isVisible();
});
