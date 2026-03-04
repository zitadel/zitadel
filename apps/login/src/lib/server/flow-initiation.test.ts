import { describe, expect, test, vi, beforeEach, afterEach } from "vitest";
import { NextRequest } from "next/server";

vi.mock("@/lib/cookies", () => ({
  getLanguageCookie: vi.fn(),
  setLanguageCookie: vi.fn(),
}));

vi.mock("@/lib/auth-utils", () => ({
  getValidLocaleFromUILocales: vi.fn(),
}));

vi.mock("@/lib/idp", () => ({
  idpTypeToSlug: vi.fn(),
}));

vi.mock("@/lib/server/loginname", () => ({
  sendLoginname: vi.fn(),
}));

vi.mock("@/lib/service-url", () => ({
  constructUrl: vi.fn(),
}));

vi.mock("@/lib/session", () => ({
  findValidSession: vi.fn(),
}));

vi.mock("@/lib/zitadel", () => ({
  createCallback: vi.fn(),
  createResponse: vi.fn(),
  getActiveIdentityProviders: vi.fn(),
  getAuthRequest: vi.fn(),
  getOrgsByDomain: vi.fn(),
  getSAMLRequest: vi.fn(),
  getSecuritySettings: vi.fn(),
  startIdentityProviderFlow: vi.fn(),
}));

vi.mock("@zitadel/client", () => ({
  create: vi.fn(),
}));

vi.mock("escape-html", () => ({
  default: (s: string) => s,
}));

import { handleOIDCFlowInitiation, FlowInitiationParams } from "./flow-initiation";

function makeRequest(url = "https://example.com/login?requestId=oidc_abc123"): NextRequest {
  return new NextRequest(url);
}

function makeBaseParams(overrides?: Partial<FlowInitiationParams>): FlowInitiationParams {
  return {
    serviceConfig: { baseUrl: "https://api.example.com" } as any,
    requestId: "oidc_abc123",
    sessions: [],
    sessionCookies: [],
    request: makeRequest(),
    ...overrides,
  };
}

describe("handleOIDCFlowInitiation — locale / cookie handling", () => {
  let mockGetLanguageCookie: ReturnType<typeof vi.fn>;
  let mockSetLanguageCookie: ReturnType<typeof vi.fn>;
  let mockGetValidLocaleFromUILocales: ReturnType<typeof vi.fn>;
  let mockGetAuthRequest: ReturnType<typeof vi.fn>;
  let mockConstructUrl: ReturnType<typeof vi.fn>;
  let mockFindValidSession: ReturnType<typeof vi.fn>;

  beforeEach(async () => {
    vi.clearAllMocks();
    vi.unstubAllEnvs();

    const cookies = await import("@/lib/cookies");
    const authUtils = await import("@/lib/auth-utils");
    const zitadel = await import("@/lib/zitadel");
    const serviceUrl = await import("@/lib/service-url");
    const session = await import("@/lib/session");

    mockGetLanguageCookie = vi.mocked(cookies.getLanguageCookie);
    mockSetLanguageCookie = vi.mocked(cookies.setLanguageCookie);
    mockGetValidLocaleFromUILocales = vi.mocked(authUtils.getValidLocaleFromUILocales);
    mockGetAuthRequest = vi.mocked(zitadel.getAuthRequest);
    mockConstructUrl = vi.mocked(serviceUrl.constructUrl);
    mockFindValidSession = vi.mocked(session.findValidSession);

    // Default auth request: no prompts, no special scopes
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: [],
        prompt: [],
        loginHint: undefined,
      },
    });

    // constructUrl returns a real URL object so redirect paths resolve correctly
    mockConstructUrl.mockImplementation((_req: any, path: string) => {
      return new URL(`https://example.com${path}`);
    });

    mockFindValidSession.mockResolvedValue(null);
  });

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  describe("when uiLocales yields no valid locale", () => {
    test("should not read or set the language cookie", async () => {
      mockGetValidLocaleFromUILocales.mockReturnValue(null);

      await handleOIDCFlowInitiation(makeBaseParams());

      expect(mockGetLanguageCookie).not.toHaveBeenCalled();
      expect(mockSetLanguageCookie).not.toHaveBeenCalled();
    });
  });

  describe("when uiLocales yields a valid locale and ZITADEL_UI_LOCALES_OVERRIDE_COOKIE is unset (default)", () => {
    beforeEach(() => {
      mockGetValidLocaleFromUILocales.mockReturnValue("de");
    });

    test("should set cookie when no existing language cookie is present", async () => {
      mockGetLanguageCookie.mockResolvedValue(undefined);

      await handleOIDCFlowInitiation(makeBaseParams());

      expect(mockGetLanguageCookie).toHaveBeenCalledOnce();
      expect(mockSetLanguageCookie).toHaveBeenCalledWith("de");
    });

    test("should NOT overwrite an existing language cookie", async () => {
      mockGetLanguageCookie.mockResolvedValue("fr");

      await handleOIDCFlowInitiation(makeBaseParams());

      expect(mockGetLanguageCookie).toHaveBeenCalledOnce();
      expect(mockSetLanguageCookie).not.toHaveBeenCalled();
    });
  });

  describe("when uiLocales yields a valid locale and ZITADEL_UI_LOCALES_OVERRIDE_COOKIE=true", () => {
    beforeEach(() => {
      vi.stubEnv("ZITADEL_UI_LOCALES_OVERRIDE_COOKIE", "true");
      mockGetValidLocaleFromUILocales.mockReturnValue("ja");
    });

    test("should set cookie even when an existing language cookie is already set", async () => {
      mockGetLanguageCookie.mockResolvedValue("fr");

      await handleOIDCFlowInitiation(makeBaseParams());

      expect(mockSetLanguageCookie).toHaveBeenCalledWith("ja");
    });

    test("should set cookie when no existing cookie is present", async () => {
      mockGetLanguageCookie.mockResolvedValue(undefined);

      await handleOIDCFlowInitiation(makeBaseParams());

      expect(mockSetLanguageCookie).toHaveBeenCalledWith("ja");
    });
  });

  describe("when ZITADEL_UI_LOCALES_OVERRIDE_COOKIE is explicitly 'false'", () => {
    test("should preserve an existing cookie", async () => {
      vi.stubEnv("ZITADEL_UI_LOCALES_OVERRIDE_COOKIE", "false");
      mockGetValidLocaleFromUILocales.mockReturnValue("es");
      mockGetLanguageCookie.mockResolvedValue("it");

      await handleOIDCFlowInitiation(makeBaseParams());

      expect(mockSetLanguageCookie).not.toHaveBeenCalled();
    });
  });
});
