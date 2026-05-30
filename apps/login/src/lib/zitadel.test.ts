import { beforeEach, describe, expect, test, vi } from "vitest";

vi.mock("./cache", () => ({
  PromiseCache: class {
    getOrFetch(_key: string, fetcher: () => Promise<unknown>) {
      return fetcher();
    }
  },
}));

vi.mock("./service", () => ({
  createServiceForHost: vi.fn(),
}));

describe("zitadel settings helpers", () => {
  beforeEach(() => {
    vi.resetModules();
    vi.clearAllMocks();
  });

  test("getBrandingSettings returns undefined when the backend omits branding settings", async () => {
    const { createServiceForHost } = await import("./service");
    vi.mocked(createServiceForHost).mockResolvedValue({
      getBrandingSettings: vi.fn().mockResolvedValue({}),
    } as any);

    const { getBrandingSettings } = await import("./zitadel");
    const result = await getBrandingSettings({ serviceConfig: { baseUrl: "https://api.example.com" } as any });

    expect(result).toBeUndefined();
  });

  test("getHostedLoginTranslation returns undefined when the backend omits translations", async () => {
    const { createServiceForHost } = await import("./service");
    vi.mocked(createServiceForHost).mockResolvedValue({
      getHostedLoginTranslation: vi.fn().mockResolvedValue({}),
    } as any);

    const { getHostedLoginTranslation } = await import("./zitadel");
    const result = await getHostedLoginTranslation({
      serviceConfig: { baseUrl: "https://api.example.com" } as any,
      locale: "en",
    });

    expect(result).toBeUndefined();
  });

  test("getDefaultOrg returns null when the backend call fails", async () => {
    const { createServiceForHost } = await import("./service");
    vi.mocked(createServiceForHost).mockResolvedValue({
      listOrganizations: vi.fn().mockRejectedValue(new Error("backend unavailable")),
    } as any);

    const { getDefaultOrg } = await import("./zitadel");
    const result = await getDefaultOrg({ serviceConfig: { baseUrl: "https://api.example.com" } as any });

    expect(result).toBeNull();
  });
});
