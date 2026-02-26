import { describe, expect, test } from "vitest";

describe("@zitadel/react", () => {
  test("exports are defined", async () => {
    const mod = await import("./index.js");
    expect(mod.ZitadelProvider).toBeDefined();
    expect(mod.useZitadel).toBeDefined();
    expect(mod.useSession).toBeDefined();
    expect(mod.useUser).toBeDefined();
    expect(mod.useToken).toBeDefined();
    expect(mod.SignedIn).toBeDefined();
    expect(mod.SignedOut).toBeDefined();
  });
});
