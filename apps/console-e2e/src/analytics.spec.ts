import { test, expect } from "@playwright/test";

test.describe("Analytics tests", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/active-user");
  });

  test("daily range picker works", async ({ page }) => {
    // instead of dynamically get the current month to make the test more robust
    // we would usually use the clock api https://playwright.dev/docs/clock
    // but this messes with token expiration in the console so the test
    // has to be dynamic over the current date
    const now = new Date();
    now.setMonth(now.getMonth() - 1);

    const currentMonth = now.toLocaleString("en-us", {
      month: "long",
    });

    await page
      .locator("cnsl-active-user-card-daily")
      .getByRole("radio", { name: "Custom" })
      .click();

    await page.getByRole("button", { name: `${currentMonth} 4,` }).click();
    await page.getByRole("button", { name: `${currentMonth} 4,` }).click();
    await page.getByRole("button", { name: `${currentMonth} 21,` }).click();

    await page.keyboard.press("Escape");

    const start = new Date(now);
    start.setDate(4);

    const end = new Date(now);
    end.setDate(21);

    await expect(
      page.getByText(
        `active users between ${start.toLocaleDateString(
          "en-US"
        )} - ${end.toLocaleDateString("en-US")}`
      )
    ).toBeVisible();
  });

  test("monthly range picker works", async ({ page }) => {
    await expect(
      page.getByText("active users in the last 3 months ")
    ).toBeVisible();

    await page
      .locator("cnsl-active-user-card-monthly")
      .getByRole("radio", { name: "6 months" })
      .click();
    await expect(
      page.getByText("active users in the last 6 months ")
    ).toBeVisible();

    await page
      .locator("cnsl-active-user-card-monthly")
      .getByRole("radio", { name: "12 months" })
      .click();
    await expect(
      page.getByText("active users in the last 12 months ")
    ).toBeVisible();

    await page
      .locator("cnsl-active-user-card-monthly")
      .getByRole("radio", { name: "custom" })
      .click();

    await page
      .getByRole("button", {
        name: (new Date().getUTCFullYear() - 1).toString(),
      })
      .click();
    await page.getByRole("button", { name: "2022" }).click();
    await page.getByRole("button", { name: "Feb" }).click();

    await page
      .getByRole("button", { name: new Date().getUTCFullYear().toString() })
      .click();
    await page.getByRole("button", { name: "2024" }).click();
    await page.getByRole("button", { name: "Aug" }).click();

    await expect(
      page.getByText("active users between 2/1/2022 - 8/31/2024")
    ).toBeVisible();
  });
});
