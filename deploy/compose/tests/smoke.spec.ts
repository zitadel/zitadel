/**
 * Compose stack wiring smoke tests
 *
 * PURPOSE: verify that all compose services are correctly wired together.
 * These are INFRASTRUCTURE tests, not login-feature tests.
 * They answer: "does the compose stack work end-to-end?"
 *
 * Test path:  browser → Traefik (port 8888) → zitadel-login → zitadel-api → postgres
 * A single successful login proves every link in that chain is functioning.
 *
 * For feature-level login tests (MFA, OIDC flows, IdP, etc.) see:
 *   apps/login/acceptance/  – run via @zitadel/login:test-acceptance
 *
 * Credentials are read from SMOKE_TEST_ADMIN_USERNAME / SMOKE_TEST_ADMIN_PASSWORD
 * (set in .env.test). The fallback values below match the defaults in .env.test and
 * docker-compose.test.yml — they are intentional so the test works out-of-the-box
 * without any extra configuration. If you change the seeded admin password, update
 * .env.test accordingly.
 * The initial admin password is seeded via ZITADEL_FIRSTINSTANCE_ORG_HUMAN_PASSWORD
 * in docker-compose.test.yml.
 */
import { expect, test } from "@playwright/test";

const username =
  process.env.SMOKE_TEST_ADMIN_USERNAME ?? "zitadel-admin@zitadel.localhost";
const password = process.env.SMOKE_TEST_ADMIN_PASSWORD ?? "Password1!";

test.describe("compose stack wiring", () => {
  test("admin can complete username/password login (proves Traefik → Login → API → DB chain)", async ({
    page,
  }) => {
    // Step 1 – enter login name.
    // Use a long timeout for the first navigation: even though the compose stack
    // is healthy, ZITADEL's LoginName page makes server-side gRPC calls to the
    // API on first render and can be slow when the instance is still warming up.
    await page.goto("/ui/v2/login/loginname");
    // The login page renders two inputs with the same test ID (one with autofocus,
    // one without — likely an accessibility duplicate). Use .first() to avoid a
    // strict-mode violation and target the active, focused input.
    const usernameInput = page.getByTestId("username-text-input").first();
    await usernameInput.waitFor({ state: "visible", timeout: 120_000 });
    await usernameInput.fill(username);
    // Wait for the button to be enabled before clicking — the page may still be
    // hydrating when the input becomes visible.
    await page.getByTestId("submit-button").first().waitFor({ state: "visible" });
    await page.getByTestId("submit-button").first().click();

    // Step 2 – enter password.
    // Wait for the password page to fully render before typing; clicking submit
    // above triggers a navigation and the input must be present and interactive.
    const passwordInput = page.getByTestId("password-text-input").first();
    await passwordInput.waitFor({ state: "visible", timeout: 30_000 });
    await passwordInput.fill(password);
    await page.getByTestId("submit-button").first().waitFor({ state: "visible" });
    await page.getByTestId("submit-button").first().click();

    // Step 3 – assert signed-in page (ZITADEL redirects here after a successful
    // direct login that has no original auth-request to return to).
    await expect(page).toHaveURL(/signedin/, { timeout: 30_000 });
  });
});
