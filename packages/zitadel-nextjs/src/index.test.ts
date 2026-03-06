import { describe, expect, test } from "vitest";

describe("@zitadel/nextjs", () => {
  test("exports are defined", async () => {
    const mod = await import("./index.js");
    // Auth — OIDC
    expect(mod.signIn).toBeDefined();
    expect(mod.signOut).toBeDefined();
    expect(mod.handleCallback).toBeDefined();
    expect(mod.getOIDCSession).toBeDefined();
    // Middleware
    expect(mod.createZitadelMiddleware).toBeDefined();
    // Webhook
    expect(mod.createZitadelWebhookHandler).toBeDefined();
    // Server actions
    expect(mod.protectedAction).toBeDefined();
    // Session
    expect(mod.getSession).toBeDefined();
    // API client
    expect(mod.createZitadelApiClient).toBeDefined();
    expect(mod.withApiClient).toBeDefined();
  });
});
