import { getLocale } from "next-intl/server";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { createServiceForHost } from "./service";
import { addHumanUser, startIdentityProviderFlow } from "./zitadel";

vi.mock("./service", () => ({
  createServiceForHost: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getLocale: vi.fn(),
  getTranslations: vi.fn(),
}));

vi.mock("./fingerprint", () => ({
  getUserAgent: vi.fn(),
}));

describe("addHumanUser", () => {
  const add = vi.fn().mockResolvedValue({ userId: "user123" });

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(createServiceForHost).mockResolvedValue({ addHumanUser: add } as any);
  });

  test.each(["de", "fr"])("saves the selected %s locale as the user's preferred language", async (locale) => {
    vi.mocked(getLocale).mockResolvedValue(locale);

    await addHumanUser({
      serviceConfig: { baseUrl: "https://api.example.com" },
      email: "user@example.com",
      firstName: "Test",
      lastName: "User",
      organization: "test-org",
    });

    expect(add).toHaveBeenCalledWith(
      expect.objectContaining({
        profile: expect.objectContaining({ givenName: "Test", familyName: "User", preferredLanguage: locale }),
      }),
    );
  });
});

describe("startIdentityProviderFlow", () => {
  const serviceConfig = { baseUrl: "https://api.example.com" } as any;
  let startIdentityProviderIntent: ReturnType<typeof vi.fn>;

  beforeEach(async () => {
    vi.clearAllMocks();
    const { createServiceForHost } = await import("./service");
    startIdentityProviderIntent = vi.fn().mockResolvedValue({
      nextStep: { case: "authUrl", value: "https://idp.example.com/auth" },
    });
    vi.mocked(createServiceForHost).mockResolvedValue({ startIdentityProviderIntent } as any);
  });

  test("forwards the login hint to the API", async () => {
    const result = await startIdentityProviderFlow({
      serviceConfig,
      idpId: "idp123",
      urls: { successUrl: "https://login/success", failureUrl: "https://login/failure", loginHint: "user@example.com" },
    });

    expect(result).toEqual({ url: "https://idp.example.com/auth" });
    expect(startIdentityProviderIntent).toHaveBeenCalledWith({
      idpId: "idp123",
      content: {
        case: "urls",
        value: { successUrl: "https://login/success", failureUrl: "https://login/failure", loginHint: "user@example.com" },
      },
    });
  });

  test("drops a login hint the API would reject instead of failing the IdP flow", async () => {
    const loginHint = "a".repeat(201);

    const result = await startIdentityProviderFlow({
      serviceConfig,
      idpId: "idp123",
      urls: { successUrl: "https://login/success", failureUrl: "https://login/failure", loginHint },
    });

    expect(result).toEqual({ url: "https://idp.example.com/auth" });
    expect(startIdentityProviderIntent).toHaveBeenCalledWith({
      idpId: "idp123",
      content: {
        case: "urls",
        value: { successUrl: "https://login/success", failureUrl: "https://login/failure" },
      },
    });
  });

  test("does not set an empty login hint", async () => {
    await startIdentityProviderFlow({
      serviceConfig,
      idpId: "idp123",
      urls: { successUrl: "https://login/success", failureUrl: "https://login/failure", loginHint: "" },
    });

    const value = startIdentityProviderIntent.mock.calls[0][0].content.value;
    expect(value).not.toHaveProperty("loginHint");
  });
});
