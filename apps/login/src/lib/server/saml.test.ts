import * as loginnameModule from "@/lib/actions/loginname";
import { loginWithSAMLAndSession } from "@/lib/server/saml";
import * as sessionModule from "@/lib/server/session";
import * as zitadelModule from "@/lib/server/zitadel";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("@/lib/server/session");
vi.mock("@/lib/server/zitadel");
vi.mock("@/lib/actions/loginname");
vi.mock("@/lib/server/grpc/interceptors/error-classification", () => ({
  isClassifiedError: (error: unknown): boolean =>
    typeof error === "object" && error !== null && "code" in error && typeof (error as any).code === "number",
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

vi.mock("@zitadel/client", () => ({
  Code: { FailedPrecondition: 9, PermissionDenied: 7 },
  create: vi.fn(),
}));

describe("loginWithSAMLAndSession", () => {
  const mockServiceConfig = { baseUrl: "https://zitadel.example.com" };
  const mockSamlRequest = "saml-123";
  const mockSessionId = "session-123";

  let mockSessions: any[];
  let mockCookies: any[];

  beforeEach(() => {
    vi.clearAllMocks();
    vi.spyOn(console, "log").mockImplementation(() => {});
    vi.spyOn(console, "error").mockImplementation(() => {});

    mockSessions = [
      {
        id: mockSessionId,
        factors: {
          user: {
            id: "user-123",
            loginName: "test@example.com",
            organizationId: "org-123",
          },
        },
      },
    ];

    mockCookies = [
      {
        id: mockSessionId,
        token: "token-123",
        loginName: "test@example.com",
        creationTs: new Date().toISOString(),
        expirationTs: new Date(Date.now() + 3600000).toISOString(),
        changeTs: new Date().toISOString(),
      },
    ];
  });

  const login = (sessions = mockSessions) =>
    loginWithSAMLAndSession({
      serviceConfig: mockServiceConfig as any,
      samlRequest: mockSamlRequest,
      sessionId: mockSessionId,
      sessions,
      sessionCookies: mockCookies,
    });

  it("should redirect to the SAML response URL when session is valid", async () => {
    vi.mocked(sessionModule.isSessionValid).mockResolvedValue(true);
    vi.mocked(zitadelModule.createResponse).mockResolvedValue({
      url: "https://sp.example.com/acs",
      binding: { case: "redirect", value: {} },
    } as any);

    const result = await login();

    expect(result).toEqual({ redirect: "https://sp.example.com/acs" });
  });

  it("should redirect to re-authenticate when session is invalid", async () => {
    vi.mocked(sessionModule.isSessionValid).mockResolvedValue(false);
    vi.mocked(loginnameModule.sendLoginname).mockResolvedValue({ redirect: "/password" });

    const result = await login();

    expect(result).toEqual({ redirect: "/password" });
    expect(loginnameModule.sendLoginname).toHaveBeenCalledWith({
      loginName: "test@example.com",
      organization: "org-123",
      requestId: `saml_${mockSamlRequest}`,
    });
    expect(zitadelModule.createResponse).not.toHaveBeenCalled();
  });

  it.each([
    ["undefined", undefined],
    ["an error", { error: "errors.localAuthenticationNotAllowed" }],
  ])("should not create a SAML response for an invalid session when re-authentication returns %s", async (_, response) => {
    vi.mocked(sessionModule.isSessionValid).mockResolvedValue(false);
    vi.mocked(loginnameModule.sendLoginname).mockResolvedValue(response as any);

    const result = await login();

    expect(result).toEqual({ error: "Session not found or invalid" });
    expect(zitadelModule.createResponse).not.toHaveBeenCalled();
  });

  it("should not create a SAML response for an invalid session without a user factor", async () => {
    vi.mocked(sessionModule.isSessionValid).mockResolvedValue(false);

    const result = await login([{ id: mockSessionId, factors: {} }]);

    expect(result).toEqual({ error: "Session not found or invalid" });
    expect(loginnameModule.sendLoginname).not.toHaveBeenCalled();
    expect(zitadelModule.createResponse).not.toHaveBeenCalled();
  });

  it.each([
    ["Errors.User.GrantRequired (SAML-foSyH49RvL)", "grantRequired"],
    ["Errors.User.ProjectRequired (SAML-foSyH49RvL)", "projectRequired"],
  ])("should return the translated authorization error for %s", async (rawMessage, expected) => {
    vi.mocked(sessionModule.isSessionValid).mockResolvedValue(true);
    vi.mocked(zitadelModule.createResponse).mockRejectedValue({ code: 7, rawMessage });

    const result = await login();

    expect(result).toEqual({ error: expected });
  });

  it("should return unknown error for other permission denied errors", async () => {
    vi.mocked(sessionModule.isSessionValid).mockResolvedValue(true);
    vi.mocked(zitadelModule.createResponse).mockRejectedValue({ code: 7, rawMessage: "Errors.AuthRequest.NoPermission" });

    const result = await login();

    expect(result).toEqual({ error: "Unknown error occurred" });
  });
});
