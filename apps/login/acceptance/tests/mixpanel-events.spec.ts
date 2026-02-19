import { faker } from "@faker-js/faker";
import { test as base, expect, Browser, BrowserContext } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
import { fileURLToPath } from "url";
import { PasswordUser } from "./user";

// Read from ".env" file.
const __dirname = path.dirname(fileURLToPath(import.meta.url));
dotenv.config({ path: path.resolve(__dirname, "../../.env.test.local") });

/**
 * Helper to create a browser context with a consent cookie and capture mixpanel requests
 */
async function createMixpanelTrackingContext(browser: Browser) {
  const baseURL = process.env.LOGIN_BASE_URL || "http://127.0.0.1:3000";
  const url = new URL(baseURL);

  const context = await browser.newContext({
    baseURL,
  });

  // Set the cc_cookie to simulate mixpanel consent
  // Set on both 127.0.0.1 and localhost to handle either baseURL
  const cookieValue = encodeURIComponent(
    JSON.stringify({
      services: {
        analytics: ["mixpanel"],
      },
    }),
  );
  await context.addCookies([
    {
      name: "cc_cookie",
      value: cookieValue,
      domain: url.hostname,
      path: "/",
    },
    {
      name: "cc_cookie",
      value: cookieValue,
      domain: "localhost",
      path: "/",
    },
  ]);

  const mixpanelRequests: string[] = [];

  context.on("request", (request) => {
    const requestUrl = request.url();
    const postData = request.postData();
    const decodedPostData = postData ? decodeURIComponent(postData) : null;
    if (requestUrl.indexOf("api-eu.mixpanel") !== -1) {
      mixpanelRequests.push(`${request.method()} ${requestUrl}${decodedPostData ? " " + decodedPostData : ""}`);
    }
  });

  return { context, mixpanelRequests };
}

/**
 * Helper to find a specific mixpanel event in the captured requests
 */
function findMixpanelEvent(requests: string[], eventName: string): string | undefined {
  return requests.find(
    (req) =>
      req.includes("api-eu.mixpanel.com/track") && req.startsWith("POST") && req.indexOf(eventName) > -1,
  );
}

/**
 * Helper to ensure a user exists for testing, returns cleanup function
 */
async function ensurePasswordUser(context: BrowserContext) {
  const page = await context.newPage();
  const user = new PasswordUser({
    email: faker.internet.email(),
    isEmailVerified: true,
    firstName: faker.person.firstName(),
    lastName: faker.person.lastName(),
    organization: "",
    phone: faker.phone.number(),
    isPhoneVerified: false,
    password: "Password1!",
    passwordChangeRequired: false,
  });
  await user.ensure(page);
  await page.close();
  return user;
}

const test = base.extend<{ user: PasswordUser }>({
  user: async ({ page }, use) => {
    const user = new PasswordUser({
      email: faker.internet.email(),
      isEmailVerified: true,
      firstName: faker.person.firstName(),
      lastName: faker.person.lastName(),
      organization: "",
      phone: faker.phone.number(),
      isPhoneVerified: false,
      password: "Password1!",
      passwordChangeRequired: false,
    });
    await user.ensure(page);
    await use(user);
    await user.cleanup();
  },
});

test.describe("Mixpanel Events - Username Submission", () => {
  test("should send username_submitted event when username is submitted", async ({ browser }) => {
    const { context, mixpanelRequests } = await createMixpanelTrackingContext(browser);
    const user = await ensurePasswordUser(context);

    try {
      const page = await context.newPage();

      await page.goto("http://localhost:3000/ui/v2/login/loginname");
      await page.waitForTimeout(2000);

      await page.getByTestId("username-text-input").pressSequentially(user.getUsername());
      await page.getByTestId("submit-button").click();

      await page.waitForTimeout(4000);
      // We need to set themixpanel cookie
      const event = findMixpanelEvent(mixpanelRequests, "username_submitted");
      expect(event).toBeTruthy();
      expect(event).toContain('"source":"login"');
    } finally {
      await user.cleanup();
      await context.close();
    }
  });
});

test.describe("Mixpanel Events - Password Submission", () => {
  test("should send password_submitted and login_success events on successful login", async ({ browser }) => {
    const { context, mixpanelRequests } = await createMixpanelTrackingContext(browser);
    const user = await ensurePasswordUser(context);

    try {
      const page = await context.newPage();
      await page.goto("http://localhost:3000/ui/v2/login/loginname");
      // await page.goto("./loginname");
      await page.waitForTimeout(2000);

      // Submit username
      await page.getByTestId("username-text-input").pressSequentially(user.getUsername());
      await page.getByTestId("submit-button").click();
      await page.waitForTimeout(2000);

      // Submit password
      await page.getByTestId("password-text-input").pressSequentially(user.getPassword());
      await page.getByTestId("submit-button").click();

      await page.waitForTimeout(8000);

      const passwordEvent = findMixpanelEvent(mixpanelRequests, "password_submitted");
      expect(passwordEvent).toBeTruthy();
      expect(passwordEvent).toContain('"source":"login"');

      const successEvent = findMixpanelEvent(mixpanelRequests, "login_success");
      expect(successEvent).toBeTruthy();
      expect(successEvent).toContain('"source":"login"');
    } finally {
      await user.cleanup();
      await context.close();
    }
  });

  test("should send password_submitted and login_failure events on wrong password", async ({ browser }) => {
    const { context, mixpanelRequests } = await createMixpanelTrackingContext(browser);
    const user = await ensurePasswordUser(context);

    try {
      const page = await context.newPage();
      await page.goto("http://localhost:3000/ui/v2/login/loginname");
      // await page.goto("./loginname");
      await page.waitForTimeout(2000);

      // Submit username
      await page.getByTestId("username-text-input").pressSequentially(user.getUsername());
      await page.getByTestId("submit-button").click();
      await page.waitForTimeout(2000);

      // Submit wrong password
      await page.getByTestId("password-text-input").pressSequentially("WrongPassword1!");
      await page.getByTestId("submit-button").click();

      await page.waitForTimeout(8000);

      const passwordEvent = findMixpanelEvent(mixpanelRequests, "password_submitted");
      expect(passwordEvent).toBeTruthy();

      const failureEvent = findMixpanelEvent(mixpanelRequests, "login_failure");
      expect(failureEvent).toBeTruthy();
      expect(failureEvent).toContain('"source":"login"');
    } finally {
      await user.cleanup();
      await context.close();
    }
  });
});

test.describe("Mixpanel Events - Password Reset", () => {
  test("should send password_reset_requested event when reset button is clicked", async ({ browser }) => {
    const { context, mixpanelRequests } = await createMixpanelTrackingContext(browser);
    const user = await ensurePasswordUser(context);

    try {
      const page = await context.newPage();
      await page.goto("http://localhost:3000/ui/v2/login/loginname");
      // await page.goto("./loginname");
      // await page.goto("./loginname");
      await page.waitForTimeout(2000);

      // Submit username
      await page.getByTestId("username-text-input").pressSequentially(user.getUsername());
      await page.getByTestId("submit-button").click();
      await page.waitForTimeout(2000);

      // Click reset password
      await page.getByTestId("reset-button").click();

      await page.waitForTimeout(8000);

      const resetEvent = findMixpanelEvent(mixpanelRequests, "password_reset_requested");
      expect(resetEvent).toBeTruthy();
      expect(resetEvent).toContain('"source":"login"');
    } finally {
      await user.cleanup();
      await context.close();
    }
  });
});

test.describe("Mixpanel Events - Registration", () => {
  test("should send register_method_selected event when registering with password", async ({ browser }) => {
    const { context, mixpanelRequests } = await createMixpanelTrackingContext(browser);
    const page = await context.newPage();

    try {
      await page.goto("./register");
      await page.goto("http://localhost:3000/ui/v2/login/register");
      // await page.goto("./loginname");
      await page.waitForTimeout(2000);

      const firstname = faker.person.firstName();
      const lastname = faker.person.lastName();
      const email = faker.internet.email();

      await page.getByTestId("firstname-text-input").pressSequentially(firstname);
      await page.getByTestId("lastname-text-input").pressSequentially(lastname);
      await page.getByTestId("email-text-input").pressSequentially(email);
      await page.getByTestId("submit-button").click();

      await page.waitForTimeout(4000);

      const methodEvent = findMixpanelEvent(mixpanelRequests, "register_method_selected");
      expect(methodEvent).toBeTruthy();
      expect(methodEvent).toContain('"source":"login"');
    } finally {
      await context.close();
    }
  });
});

test.describe("Mixpanel Events - Page View", () => {
  test("should send page_view event on navigation", async ({ browser }) => {
    const { context, mixpanelRequests } = await createMixpanelTrackingContext(browser);
    const page = await context.newPage();

    try {
      await page.goto("http://localhost:3000/ui/v2/login/loginname");
      // await page.goto("./loginname");
      // We need a long timeout for mixpanel sometimes
      await page.waitForTimeout(8000);

      const pageViewEvent = findMixpanelEvent(mixpanelRequests, "page_view");
      expect(pageViewEvent).toBeTruthy();
      expect(pageViewEvent).toContain('"source":"login"');
      expect(pageViewEvent).toContain("loginname");
    } finally {
      await context.close();
    }
  });
});

test.describe("Mixpanel Events - Consent Required", () => {
  test("should not send events when no consent cookie is set", async ({ browser }) => {
    const baseURL = process.env.LOGIN_BASE_URL || "http://127.0.0.1:3000";

    // Create context WITHOUT consent cookie
    const context = await browser.newContext({ baseURL });
    const mixpanelRequests: string[] = [];

    context.on("request", (request) => {
      const requestUrl = request.url();
      if (requestUrl.indexOf("api-eu.mixpanel") !== -1) {
        mixpanelRequests.push(requestUrl);
      }
    });

    const page = await context.newPage();

    try {
      await page.goto("http://localhost:3000/ui/v2/login/loginname");

      await page.waitForTimeout(8000);

      // No mixpanel requests should have been made
      expect(mixpanelRequests.length).toBe(0);
    } finally {
      await context.close();
    }
  });
});

test.describe("Mixpanel Events - Language Selection", () => {
  test("should send language_selected event when language is changed", async ({ browser }) => {
    const { context, mixpanelRequests } = await createMixpanelTrackingContext(browser);
    const page = await context.newPage();

    try {
      await page.goto("http://localhost:3000/ui/v2/login/loginname");
      await page.waitForTimeout(2000);

      // Open the language switcher and select a different language
      const languageSwitcher = page.locator("[data-headlessui-state]").first();
      await languageSwitcher.click();
      await page.waitForTimeout(500);

      // Click on a non-English language option (e.g., Deutsch)
      const germanOption = page.getByText("Deutsch");
      if (await germanOption.isVisible({ timeout: 2000 }).catch(() => false)) {
        await germanOption.click();
        await page.waitForTimeout(4000);

        const languageEvent = findMixpanelEvent(mixpanelRequests, "language_selected");
        expect(languageEvent).toBeTruthy();
        expect(languageEvent).toContain('"source":"login"');
        expect(languageEvent).toContain('"language":"de"');
      }
    } finally {
      await context.close();
    }
  });
});
