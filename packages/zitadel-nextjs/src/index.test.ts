import { describe, expect, test } from "vitest";

describe("@zitadel/nextjs", () => {
  test("exports are defined", async () => {
    const mod = await import("./index.js");
    expect(mod.signIn).toBeDefined();
    expect(mod.signOut).toBeDefined();
    expect(mod.handleCallback).toBeDefined();
    expect(mod.createZitadelMiddleware).toBeDefined();
    expect(mod.createZitadelWebhookHandler).toBeDefined();
    expect(mod.protectedAction).toBeDefined();
    expect(mod.getSession).toBeDefined();
  });
});
