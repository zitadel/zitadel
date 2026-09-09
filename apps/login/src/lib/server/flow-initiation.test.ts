import type { Cookie } from "@/lib/cookies";
import { NextRequest } from "next/server";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { FlowInitiationParams, handleOIDCFlowInitiation } from "./flow-initiation";

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
  getLoginSettings: vi.fn(),
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

function makeRequest(url = "https://example.com/login?requestId=oidc_abc123"): NextRequest {
  return new NextRequest(url);
}

function makeCookie(id: string, overrides?: Partial<Cookie>): Cookie {
  return { id, token: "tok", loginName: "", creationTs: "", expirationTs: "", changeTs: "", ...overrides };
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

describe("handleOIDCFlowInitiation — org-scoped session filtering", () => {
  let mockGetAuthRequest: ReturnType<typeof vi.fn>;
  let mockConstructUrl: ReturnType<typeof vi.fn>;
  let mockFindValidSession: ReturnType<typeof vi.fn>;
  let mockSendLoginname: ReturnType<typeof vi.fn>;
  let mockGetLoginSettings: ReturnType<typeof vi.fn>;

  const orgSession = {
    id: "session-org",
    factors: { user: { id: "user1", organizationId: "111111", loginName: "user@org.com" } },
  };
  const otherOrgSession = {
    id: "session-other",
    factors: { user: { id: "user2", organizationId: "999999", loginName: "user@other.com" } },
  };

  beforeEach(async () => {
    vi.clearAllMocks();
    vi.unstubAllEnvs();

    const zitadel = await import("@/lib/zitadel");
    const serviceUrl = await import("@/lib/service-url");
    const session = await import("@/lib/session");
    const authUtils = await import("@/lib/auth-utils");
    const loginname = await import("@/lib/server/loginname");

    mockGetAuthRequest = vi.mocked(zitadel.getAuthRequest);
    mockConstructUrl = vi.mocked(serviceUrl.constructUrl);
    mockFindValidSession = vi.mocked(session.findValidSession);
    mockSendLoginname = vi.mocked(loginname.sendLoginname);
    mockGetLoginSettings = vi.mocked(zitadel.getLoginSettings);
    mockGetLoginSettings.mockResolvedValue(undefined);
    vi.mocked(authUtils.getValidLocaleFromUILocales).mockReturnValue(null);

    mockConstructUrl.mockImplementation((_req: any, path: string) => {
      return new URL(`https://example.com${path}`);
    });

    mockFindValidSession.mockResolvedValue(null);
    // Default: hint cannot be resolved to a redirect, exercising the prefilled
    // /loginname fallback. Individual tests override this to assert the happy path.
    mockSendLoginname.mockResolvedValue(undefined);
  });

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  test("should redirect to /loginname when org scope filters out all sessions (default prompt)", async () => {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: ["urn:zitadel:iam:org:id:111111"],
        prompt: [],
        loginHint: undefined,
      },
    });

    const res = await handleOIDCFlowInitiation(
      makeBaseParams({
        sessions: [otherOrgSession] as any,
        sessionCookies: [makeCookie("session-other")],
      }),
    );

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/loginname");
    expect(location).not.toContain("/accounts");
  });

  test("should resolve loginHint server-side (straight to next step) when org scope filters out all sessions", async () => {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: ["urn:zitadel:iam:org:id:111111"],
        prompt: [],
        loginHint: "user@example.com",
      },
    });

    mockSendLoginname.mockResolvedValue({
      redirect: "/password?loginName=user%40example.com",
    });

    const res = await handleOIDCFlowInitiation(
      makeBaseParams({
        sessions: [otherOrgSession] as any,
        sessionCookies: [makeCookie("session-other")],
      }),
    );

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/password");
    expect(location).not.toContain("/loginname");
    expect(location).not.toContain("/accounts");
  });

  test("should prefill /loginname WITHOUT submit=true when loginHint cannot be resolved and org scope filters out all sessions", async () => {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: ["urn:zitadel:iam:org:id:111111"],
        prompt: [],
        loginHint: "user@example.com",
      },
    });

    // mockSendLoginname resolves to undefined (default) → fallback path.

    const res = await handleOIDCFlowInitiation(
      makeBaseParams({
        sessions: [otherOrgSession] as any,
        sessionCookies: [makeCookie("session-other")],
      }),
    );

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/loginname");
    expect(location).toContain("loginName=user%40example.com");
    expect(location).not.toContain("submit=true");
    expect(location).not.toContain("/accounts");
  });

  test("should redirect to /accounts when org-eligible sessions exist but none are valid (default prompt)", async () => {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: ["urn:zitadel:iam:org:id:111111"],
        prompt: [],
        loginHint: undefined,
      },
    });

    const res = await handleOIDCFlowInitiation(
      makeBaseParams({
        sessions: [orgSession] as any,
        sessionCookies: [makeCookie("session-org")],
      }),
    );

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/accounts");
  });

  test("should resolve loginHint server-side when it matches no session but an unrelated session exists (default prompt)", async () => {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: [],
        prompt: [],
        loginHint: "user@example.com",
      },
    });

    // An unrelated session is present (so eligibleSessions is non-empty), but
    // findValidSession filtered by the hint returns nothing.
    mockFindValidSession.mockResolvedValue(null);

    mockSendLoginname.mockResolvedValue({
      redirect: "/password?loginName=user%40example.com",
    });

    const res = await handleOIDCFlowInitiation(
      makeBaseParams({
        sessions: [otherOrgSession] as any,
        sessionCookies: [makeCookie("session-other")],
      }),
    );

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/password");
    expect(location).not.toContain("/loginname");
    expect(location).not.toContain("/accounts");
  });

  test("should fall back to prefilled /loginname (no submit) when loginHint cannot be resolved and an unrelated session exists (default prompt)", async () => {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: [],
        prompt: [],
        loginHint: "user@example.com",
      },
    });

    mockFindValidSession.mockResolvedValue(null);
    // mockSendLoginname resolves to undefined (default) → fallback path.

    const res = await handleOIDCFlowInitiation(
      makeBaseParams({
        sessions: [otherOrgSession] as any,
        sessionCookies: [makeCookie("session-other")],
      }),
    );

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/loginname");
    expect(location).toContain("loginName=user%40example.com");
    expect(location).not.toContain("submit=true");
    expect(location).not.toContain("/accounts");
  });

  test("should redirect to /loginname when org scope filters out all sessions (SELECT_ACCOUNT prompt)", async () => {
    const { Prompt } = await import("@zitadel/proto/zitadel/oidc/v2/authorization_pb");

    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: ["urn:zitadel:iam:org:id:111111"],
        prompt: [Prompt.SELECT_ACCOUNT],
        loginHint: undefined,
      },
    });

    const res = await handleOIDCFlowInitiation(
      makeBaseParams({
        sessions: [otherOrgSession] as any,
        sessionCookies: [makeCookie("session-other")],
      }),
    );

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/loginname");
    expect(location).not.toContain("/accounts");
  });

  test("should redirect to /accounts when org-eligible sessions exist (SELECT_ACCOUNT prompt)", async () => {
    const { Prompt } = await import("@zitadel/proto/zitadel/oidc/v2/authorization_pb");

    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: ["urn:zitadel:iam:org:id:111111"],
        prompt: [Prompt.SELECT_ACCOUNT],
        loginHint: undefined,
      },
    });

    const res = await handleOIDCFlowInitiation(
      makeBaseParams({
        sessions: [orgSession] as any,
        sessionCookies: [makeCookie("session-org")],
      }),
    );

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/accounts");
  });

  test("should redirect to /loginname when no org scope and no sessions", async () => {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: [],
        prompt: [],
        loginHint: undefined,
      },
    });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [] }));

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/loginname");
  });

  test("should resolve loginHint server-side when there are no sessions (default prompt)", async () => {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: [],
        prompt: [],
        loginHint: "user@example.com",
      },
    });

    mockSendLoginname.mockResolvedValue({
      redirect: "/password?loginName=user%40example.com",
    });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [] }));

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/password");
    expect(location).not.toContain("/loginname");
  });

  test("should prefill /loginname WITHOUT submit=true when there are no sessions and loginHint cannot be resolved", async () => {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: [],
        prompt: [],
        loginHint: "user@example.com",
      },
    });

    // mockSendLoginname resolves to undefined (default) → fallback path.

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [] }));

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/loginname");
    expect(location).toContain("loginName=user%40example.com");
    expect(location).not.toContain("submit=true");
  });

  test("should resolve loginHint without forwarding enumeration settings (derived server-side by sendLoginname)", async () => {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: [],
        prompt: [],
        loginHint: "unknown@example.com",
      },
    });

    // Enumeration protection is derived inside sendLoginname from the request-context
    // login settings; resolveLoginHint only passes the hint and the org context, and
    // an unknown hint still redirects to the (fake) /password step.
    mockSendLoginname.mockResolvedValue({
      redirect: "/password?loginName=unknown%40example.com",
    });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [] }));

    expect(mockSendLoginname).toHaveBeenCalledWith(
      expect.objectContaining({
        loginName: "unknown@example.com",
      }),
    );
    expect(mockSendLoginname).not.toHaveBeenCalledWith(
      expect.objectContaining({
        ignoreUnknownUsernames: expect.anything(),
      }),
    );

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/password");
    expect(location).not.toContain("/loginname");
  });

  test("should redirect to an absolute IdP URL as-is without prepending the base path (domain discovery auto-redirect)", async () => {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: [],
        prompt: [],
        loginHint: "user@discovered-org.com",
      },
    });

    // Simulate a base path being configured, as in ZITADEL Cloud (/ui/v2/login).
    mockConstructUrl.mockImplementation((_req: any, path: string) => {
      return new URL(`https://example.com/ui/v2/login${path}`);
    });

    // sendLoginname resolved the hint via domain discovery to an org with a
    // single external IdP and returns the absolute authorize URL of that IdP.
    const idpUrl = "https://login.microsoftonline.com/tenant-id/oauth2/v2.0/authorize?client_id=xyz&state=abc";
    mockSendLoginname.mockResolvedValue({ redirect: idpUrl });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [] }));

    const location = res.headers.get("location") ?? "";
    expect(location).toBe(idpUrl);
    // The base path must never be glued onto an absolute URL
    // (regression: https://<host>/ui/v2/loginhttps://login.microsoftonline.com/...).
    expect(location).not.toContain("/ui/v2/login");
    expect(mockConstructUrl).not.toHaveBeenCalledWith(expect.anything(), idpUrl);
  });

  test("should render an auto-submit form when loginHint resolves to a SAML POST-binding IdP", async () => {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: [],
        prompt: [],
        loginHint: "user@discovered-org.com",
      },
    });

    // sendLoginname resolved the hint via domain discovery to an org whose
    // single IdP is SAML with POST binding: the AuthnRequest is delivered as
    // form fields, not a redirect URL.
    mockSendLoginname.mockResolvedValue({
      samlData: {
        url: "https://adfs.example.com/adfs/ls",
        fields: { SAMLRequest: "PHNhbWxwOkF1dGhuUmVxdWVzdD4=", RelayState: "relay-123" },
      },
    });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [] }));

    // No redirect: the response is an HTML page auto-posting the form to the IdP.
    expect(res.headers.get("location")).toBeNull();
    expect(res.headers.get("content-type")).toContain("text/html");

    const html = await res.text();
    expect(html).toContain('action="https://adfs.example.com/adfs/ls"');
    expect(html).toContain('name="SAMLRequest"');
    expect(html).toContain('value="PHNhbWxwOkF1dGhuUmVxdWVzdD4="');
    expect(html).toContain('name="RelayState"');
    expect(html).toContain("document.forms[0].submit()");
  });

  test("should block unsafe SAML post URLs from loginHint resolution and fall back to /loginname", async () => {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: [],
        prompt: [],
        loginHint: "user@example.com",
      },
    });

    mockSendLoginname.mockResolvedValue({
      samlData: {
        url: "javascript:alert(1)",
        fields: { SAMLRequest: "abc" },
      },
    });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [] }));

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/loginname");
    expect(location).not.toContain("javascript:");
  });

  test("should block unsafe absolute redirect URLs from loginHint resolution and fall back to /loginname", async () => {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: [],
        prompt: [],
        loginHint: "user@example.com",
      },
    });

    mockSendLoginname.mockResolvedValue({ redirect: "javascript:alert(1)" });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [] }));

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/loginname");
    expect(location).not.toContain("javascript:");
  });
});

describe("handleOIDCFlowInitiation — stale session cookie fallback (#12252)", () => {
  let mockGetAuthRequest: ReturnType<typeof vi.fn>;
  let mockConstructUrl: ReturnType<typeof vi.fn>;
  let mockFindValidSession: ReturnType<typeof vi.fn>;
  let mockSendLoginname: ReturnType<typeof vi.fn>;

  // Cookie entries left behind after an RP-initiated logout terminated the
  // server-side sessions: no live Session objects, but loginName is still known.
  const staleOrgCookie = makeCookie("session-stale-org", { loginName: "user@org.com", organization: "111111" });
  const staleOtherOrgCookie = makeCookie("session-stale-other", { loginName: "user@other.com", organization: "999999" });
  const otherOrgSession = {
    id: "session-other",
    factors: { user: { id: "user2", organizationId: "999999", loginName: "user@other.com" } },
  };

  function authRequestWith(overrides: Record<string, unknown>) {
    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: [],
        prompt: [],
        loginHint: undefined,
        ...overrides,
      },
    });
  }

  beforeEach(async () => {
    vi.clearAllMocks();
    vi.unstubAllEnvs();

    const zitadel = await import("@/lib/zitadel");
    const serviceUrl = await import("@/lib/service-url");
    const session = await import("@/lib/session");
    const authUtils = await import("@/lib/auth-utils");
    const loginname = await import("@/lib/server/loginname");

    mockGetAuthRequest = vi.mocked(zitadel.getAuthRequest);
    mockConstructUrl = vi.mocked(serviceUrl.constructUrl);
    mockFindValidSession = vi.mocked(session.findValidSession);
    mockSendLoginname = vi.mocked(loginname.sendLoginname);
    vi.mocked(zitadel.getLoginSettings).mockResolvedValue(undefined);
    vi.mocked(zitadel.getSecuritySettings).mockResolvedValue(undefined);
    vi.mocked(authUtils.getValidLocaleFromUILocales).mockReturnValue(null);

    mockConstructUrl.mockImplementation((_req: any, path: string) => {
      return new URL(`https://example.com${path}`);
    });

    mockFindValidSession.mockResolvedValue(null);
    mockSendLoginname.mockResolvedValue(undefined);
  });

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  test("should redirect to /accounts when there are no live sessions but the cookie references an account", async () => {
    authRequestWith({});

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [], sessionCookies: [staleOrgCookie] }));

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/accounts");
    expect(location).not.toContain("/loginname");
  });

  test("should ignore cookie entries without a loginName", async () => {
    authRequestWith({});

    const res = await handleOIDCFlowInitiation(
      makeBaseParams({ sessions: [], sessionCookies: [{ ...staleOrgCookie, loginName: "" }] }),
    );

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/loginname");
  });

  test("should prefer login_hint over the cookie fallback when there are no live sessions", async () => {
    authRequestWith({ loginHint: "hinted@org.com" });
    mockSendLoginname.mockResolvedValue({ redirect: "/password?requestId=oidc_abc123" });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [], sessionCookies: [staleOrgCookie] }));

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/password");
    expect(location).not.toContain("/accounts");
  });

  test("should not show the account picker for prompt=none even if the cookie references an account", async () => {
    const { Prompt } = await import("@zitadel/proto/zitadel/oidc/v2/authorization_pb");
    authRequestWith({ prompt: [Prompt.NONE] });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [], sessionCookies: [staleOrgCookie] }));

    const location = res.headers.get("location") ?? "";
    expect(location).not.toContain("/accounts");
  });

  test("should redirect to /loginname when the cookie account belongs to a different org than the org scope", async () => {
    authRequestWith({ scope: ["urn:zitadel:iam:org:id:111111"] });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [], sessionCookies: [staleOtherOrgCookie] }));

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/loginname");
    expect(location).not.toContain("/accounts");
  });

  test("should redirect to /accounts when the cookie account matches the org scope", async () => {
    authRequestWith({ scope: ["urn:zitadel:iam:org:id:111111"] });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [], sessionCookies: [staleOrgCookie] }));

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/accounts");
    expect(location).toContain("organization=111111");
  });

  test("should keep /accounts for SELECT_ACCOUNT when live sessions are filtered out by org but the cookie has a matching account", async () => {
    const { Prompt } = await import("@zitadel/proto/zitadel/oidc/v2/authorization_pb");
    authRequestWith({ scope: ["urn:zitadel:iam:org:id:111111"], prompt: [Prompt.SELECT_ACCOUNT] });

    const res = await handleOIDCFlowInitiation(
      makeBaseParams({
        sessions: [otherOrgSession] as any,
        sessionCookies: [{ ...staleOtherOrgCookie, id: "session-other" }, staleOrgCookie],
      }),
    );

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/accounts");
  });

  test("should not show the account picker for prompt=login even if the cookie references an account", async () => {
    const { Prompt } = await import("@zitadel/proto/zitadel/oidc/v2/authorization_pb");
    authRequestWith({ prompt: [Prompt.LOGIN] });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [], sessionCookies: [staleOrgCookie] }));

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/loginname");
    expect(location).not.toContain("/accounts");
  });

  test("should let login_hint win over the cookie fallback for SELECT_ACCOUNT when no live session is eligible", async () => {
    const { Prompt } = await import("@zitadel/proto/zitadel/oidc/v2/authorization_pb");
    authRequestWith({
      scope: ["urn:zitadel:iam:org:id:111111"],
      prompt: [Prompt.SELECT_ACCOUNT],
      loginHint: "hinted@org.com",
    });

    const res = await handleOIDCFlowInitiation(
      makeBaseParams({
        sessions: [otherOrgSession] as any,
        sessionCookies: [
          makeCookie("session-other", { loginName: "user@other.com", organization: "999999" }),
          staleOrgCookie,
        ],
      }),
    );

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/loginname");
    expect(location).toContain("loginName=hinted%40org.com");
    expect(location).not.toContain("/accounts");
  });

  test("should keep /accounts (default prompt) when live sessions are filtered out by org but the cookie has a matching account", async () => {
    authRequestWith({ scope: ["urn:zitadel:iam:org:id:111111"] });

    const res = await handleOIDCFlowInitiation(
      makeBaseParams({
        sessions: [otherOrgSession] as any,
        sessionCookies: [{ ...staleOtherOrgCookie, id: "session-other" }, staleOrgCookie],
      }),
    );

    const location = res.headers.get("location") ?? "";
    expect(location).toContain("/accounts");
  });
});

describe("handleOIDCFlowInitiation — Prompt.LOGIN + loginHint requestId prefix", () => {
  let mockGetAuthRequest: ReturnType<typeof vi.fn>;
  let mockConstructUrl: ReturnType<typeof vi.fn>;
  let mockSendLoginname: ReturnType<typeof vi.fn>;

  const existingSession = {
    id: "session-1",
    factors: { user: { id: "user1", organizationId: "org1", loginName: "user@example.com" } },
  };

  beforeEach(async () => {
    vi.clearAllMocks();
    vi.unstubAllEnvs();

    const zitadel = await import("@/lib/zitadel");
    const serviceUrl = await import("@/lib/service-url");
    const authUtils = await import("@/lib/auth-utils");
    const loginname = await import("@/lib/server/loginname");

    mockGetAuthRequest = vi.mocked(zitadel.getAuthRequest);
    mockConstructUrl = vi.mocked(serviceUrl.constructUrl);
    mockSendLoginname = vi.mocked(loginname.sendLoginname);
    vi.mocked(authUtils.getValidLocaleFromUILocales).mockReturnValue(null);

    mockConstructUrl.mockImplementation((_req: any, path: string) => {
      return new URL(`https://example.com${path}`);
    });
  });

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  test("should pass requestId with oidc_ prefix to sendLoginname (not authRequest.id)", async () => {
    const { Prompt } = await import("@zitadel/proto/zitadel/oidc/v2/authorization_pb");

    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: [],
        prompt: [Prompt.LOGIN],
        loginHint: "user@example.com",
      },
    });

    mockSendLoginname.mockResolvedValue({
      redirect: "/password?loginName=user%40example.com",
    });

    await handleOIDCFlowInitiation(
      makeBaseParams({
        requestId: "oidc_abc123",
        sessions: [existingSession] as any,
        sessionCookies: [makeCookie("session-1")],
      }),
    );

    expect(mockSendLoginname).toHaveBeenCalledWith(
      expect.objectContaining({
        loginName: "user@example.com",
        requestId: "oidc_abc123",
      }),
    );
  });

  test("should NOT pass raw authRequest.id without oidc_ prefix to sendLoginname", async () => {
    const { Prompt } = await import("@zitadel/proto/zitadel/oidc/v2/authorization_pb");

    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: [],
        prompt: [Prompt.LOGIN],
        loginHint: "user@example.com",
      },
    });

    mockSendLoginname.mockResolvedValue({
      redirect: "/password?loginName=user%40example.com",
    });

    await handleOIDCFlowInitiation(
      makeBaseParams({
        requestId: "oidc_abc123",
        sessions: [existingSession] as any,
        sessionCookies: [makeCookie("session-1")],
      }),
    );

    // Ensure the raw authRequest.id is NOT what's passed
    expect(mockSendLoginname).not.toHaveBeenCalledWith(
      expect.objectContaining({
        requestId: "abc123",
      }),
    );
  });
});

describe("handleOIDCFlowInitiation — idp scope (urn:zitadel:iam:org:idp:id)", () => {
  let mockGetAuthRequest: ReturnType<typeof vi.fn>;
  let mockConstructUrl: ReturnType<typeof vi.fn>;
  let mockGetActiveIdentityProviders: ReturnType<typeof vi.fn>;
  let mockStartIdentityProviderFlow: ReturnType<typeof vi.fn>;
  let mockIdpTypeToSlug: ReturnType<typeof vi.fn>;

  beforeEach(async () => {
    vi.clearAllMocks();
    vi.unstubAllEnvs();

    const zitadel = await import("@/lib/zitadel");
    const serviceUrl = await import("@/lib/service-url");
    const authUtils = await import("@/lib/auth-utils");
    const idpLib = await import("@/lib/idp");

    mockGetAuthRequest = vi.mocked(zitadel.getAuthRequest);
    mockConstructUrl = vi.mocked(serviceUrl.constructUrl);
    mockGetActiveIdentityProviders = vi.mocked(zitadel.getActiveIdentityProviders);
    mockStartIdentityProviderFlow = vi.mocked(zitadel.startIdentityProviderFlow);
    mockIdpTypeToSlug = vi.mocked(idpLib.idpTypeToSlug);
    vi.mocked(authUtils.getValidLocaleFromUILocales).mockReturnValue(null);

    mockConstructUrl.mockImplementation((_req: any, path: string) => {
      return new URL(`https://example.com${path}`);
    });
    mockIdpTypeToSlug.mockReturnValue("azure");
  });

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  test("should use the type of the scoped IdP, not the first active IdP (multi-IdP org)", async () => {
    const { IdentityProviderType } = await import("@zitadel/proto/zitadel/settings/v2/login_settings_pb");

    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: ["urn:zitadel:iam:org:idp:id:idp-2"],
        prompt: [],
        loginHint: undefined,
      },
    });

    // Two active IdPs: the scope selects the SECOND one. Regression: the slug
    // was previously derived from identityProviders[0].type.
    mockGetActiveIdentityProviders.mockResolvedValue({
      identityProviders: [
        { id: "idp-1", type: IdentityProviderType.GITHUB },
        { id: "idp-2", type: IdentityProviderType.AZURE_AD },
      ],
    });

    mockStartIdentityProviderFlow.mockResolvedValue({
      url: "https://login.microsoftonline.com/tenant/oauth2/v2.0/authorize?client_id=xyz",
    });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [] }));

    expect(mockIdpTypeToSlug).toHaveBeenCalledWith(IdentityProviderType.AZURE_AD);
    expect(mockIdpTypeToSlug).not.toHaveBeenCalledWith(IdentityProviderType.GITHUB);
    expect(res.headers.get("location")).toBe("https://login.microsoftonline.com/tenant/oauth2/v2.0/authorize?client_id=xyz");
  });

  test("should block unsafe IdP URLs with a 400 instead of redirecting or rendering a form", async () => {
    const { IdentityProviderType } = await import("@zitadel/proto/zitadel/settings/v2/login_settings_pb");

    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: ["urn:zitadel:iam:org:idp:id:idp-1"],
        prompt: [],
        loginHint: undefined,
      },
    });

    mockGetActiveIdentityProviders.mockResolvedValue({
      identityProviders: [{ id: "idp-1", type: IdentityProviderType.SAML }],
    });

    mockStartIdentityProviderFlow.mockResolvedValue({
      url: "javascript:alert(1)",
      fields: { SAMLRequest: "abc" },
    });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [] }));

    expect(res.status).toBe(400);
    const body = await res.json();
    expect(body.error).toContain("Unsafe redirect URI");
  });

  test("should render the auto-submit form for a scoped SAML POST-binding IdP", async () => {
    const { IdentityProviderType } = await import("@zitadel/proto/zitadel/settings/v2/login_settings_pb");

    mockGetAuthRequest.mockResolvedValue({
      authRequest: {
        id: "abc123",
        uiLocales: [],
        scope: ["urn:zitadel:iam:org:idp:id:idp-1"],
        prompt: [],
        loginHint: undefined,
      },
    });

    mockGetActiveIdentityProviders.mockResolvedValue({
      identityProviders: [{ id: "idp-1", type: IdentityProviderType.SAML }],
    });

    mockStartIdentityProviderFlow.mockResolvedValue({
      url: "https://adfs.example.com/adfs/ls",
      fields: { SAMLRequest: "PHNhbWxwOkF1dGhuUmVxdWVzdD4=", RelayState: "relay-123" },
    });

    const res = await handleOIDCFlowInitiation(makeBaseParams({ sessions: [] }));

    expect(res.headers.get("location")).toBeNull();
    expect(res.headers.get("content-type")).toContain("text/html");
    const html = await res.text();
    expect(html).toContain('action="https://adfs.example.com/adfs/ls"');
    expect(html).toContain('name="SAMLRequest"');
  });
});
