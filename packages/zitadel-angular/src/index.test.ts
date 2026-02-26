import { describe, expect, test } from "vitest";

describe("@zitadel/angular", () => {
  test("exports are defined", async () => {
    const mod = await import("./index.js");
    expect(mod.ZitadelAuthService).toBeDefined();
    expect(mod.zitadelAuthGuard).toBeDefined();
    expect(mod.zitadelAuthInterceptor).toBeDefined();
    expect(mod.provideZitadel).toBeDefined();
    expect(mod.createWebhookHandler).toBeDefined();
  });
});
