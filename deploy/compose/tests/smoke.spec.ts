/**
 * Compose stack wiring smoke tests
 *
 * PURPOSE: verify that all compose services are correctly wired together.
 * These are INFRASTRUCTURE tests, not login-feature tests.
 * They answer: "does the compose stack work end-to-end?"
 *
 * Test path:  browser → Traefik (port 8080) → zitadel-login → zitadel-api → postgres
 * A single successful login proves every link in that chain is functioning.
 *
 * For feature-level login tests (MFA, OIDC flows, IdP, etc.) see:
 *   apps/login/acceptance/  – run via @zitadel/compose:test-login-acceptance
 *
 * Credentials are read from .env.test so nothing is hardcoded here.
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
    // Step 1 – enter login name
    await page.goto("/ui/v2/login/loginname");
    await page.getByTestId("username-text-input").pressSequentially(username);
    await page.getByTestId("submit-button").click();

    // Step 2 – enter password
    await page.getByTestId("password-text-input").pressSequentially(password);
    await page.getByTestId("submit-button").click();

    // Step 3 – assert signed-in page (ZITADEL redirects here after a successful
    // direct login that has no original auth-request to return to).
    await expect(page).toHaveURL(/signedin/, { timeout: 15_000 });
  });
});
