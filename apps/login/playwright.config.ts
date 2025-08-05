import { defineConfig, devices } from "@playwright/test";

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  testDir: './acceptance/tests',
  testMatch: '*.spec.ts',
  /* Run tests in files in parallel */
  fullyParallel: true,
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,
  /* Retry on CI only */
  retries: process.env.CI ? 2 : 0,
//  expect: {
//    timeout: 10_000, // 10 seconds
//  },
//  timeout: 300 * 1000, // 5 minutes
//  globalTimeout: 30 * 60_000, // 30 minutes
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  reporter: [
    ["line"],
    ["html", { open: process.env.CI ? "never" : "on-failure", host: "0.0.0.0", outputFolder: "./playwright-report/html" }],
  ],
  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    baseURL: process.env.LOGIN_BASE_URL || "http://localhost:3000/ui/v2/login/",
    trace: "retain-on-failure",
    headless: true,
    screenshot: "only-on-failure",
    video: "retain-on-failure",
    ignoreHTTPSErrors: true,
  },
  outputDir: "test-results/results",

  /* Configure projects for major browsers */
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
    /*
            {
              name: "firefox",
              use: { ...devices["Desktop Firefox"] },
            },
            TODO: webkit fails. Is this a bug?
            {
              name: 'webkit',
              use: { ...devices['Desktop Safari'] },
            },
        */

    /* Test against mobile viewports. */
    // {
    //   name: 'Mobile Chrome',
    //   use: { ...devices['Pixel 5'] },
    // },
    // {
    //   name: 'Mobile Safari',
    //   use: { ...devices['iPhone 12'] },
    // },

    /* Test against branded browsers. */
    // {
    //   name: 'Microsoft Edge',
    //   use: { ...devices['Desktop Edge'], channel: 'msedge' },
    // },
    // {
    //   name: 'Google Chrome',
    //   use: { ...devices['Desktop Chrome'], channel: 'chrome' },
    // },
  ],
});
