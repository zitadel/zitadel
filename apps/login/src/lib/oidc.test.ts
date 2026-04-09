import { beforeEach, describe, expect, it, vi } from "vitest";
import { loginWithOIDCAndSession } from "./oidc";
import * as loginnameModule from "./server/loginname";
import * as sessionModule from "./session";
import * as zitadelModule from "./zitadel";

vi.mock("./session");
vi.mock("./zitadel");
vi.mock("./server/loginname");
vi.mock("@/lib/grpc/interceptors/error-classification", () => ({
  isClassifiedError: (error: unknown): boolean =>
    typeof error === "object" && error !== null && "code" in error && typeof (error as any).code === "number",
}));

vi.mock("@zitadel/client", () => ({
  Code: { FailedPrecondition: 9 },
  ConnectError: class MockConnectError extends Error {
    code: number;
    constructor(msg: string, code: number) {
      super(msg);
      this.code = code;
    }
  },
  create: vi.fn(),
}));

describe("loginWithOIDCAndSession", () => {
  const mockServiceUrl = "https://zitadel.example.com";
  const mockAuthRequest = "auth-123";
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
          password: {
            verifiedAt: { seconds: BigInt(Math.floor(Date.now() / 1000)) },
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

  it("should redirect to callback URL when session is valid", async () => {
    vi.mocked(sessionModule.isSessionValid).mockResolvedValue(true);
    vi.mocked(zitadelModule.createCallback).mockResolvedValue({
      callbackUrl: "https://app.example.com/callback",
    } as any);

    const result = await loginWithOIDCAndSession({
      serviceUrl: mockServiceUrl,
      authRequest: mockAuthRequest,
      sessionId: mockSessionId,
      sessions: mockSessions,
      sessionCookies: mockCookies,
    });

    expect(result).toEqual({ redirect: "https://app.example.com/callback" });
  });

  it("should redirect to re-authenticate when session is invalid", async () => {
    vi.mocked(sessionModule.isSessionValid).mockResolvedValue(false);
    vi.mocked(loginnameModule.sendLoginname).mockResolvedValue({
      redirect: "/password",
    });

    const result = await loginWithOIDCAndSession({
      serviceUrl: mockServiceUrl,
      authRequest: mockAuthRequest,
      sessionId: mockSessionId,
      sessions: mockSessions,
      sessionCookies: mockCookies,
    });

    expect(result).toEqual({ redirect: "/password" });
    expect(loginnameModule.sendLoginname).toHaveBeenCalledWith({
      loginName: "test@example.com",
      organization: "org-123",
      requestId: `oidc_${mockAuthRequest}`,
    });
  });

  it("should return error when session not found", async () => {
    const result = await loginWithOIDCAndSession({
      serviceUrl: mockServiceUrl,
      authRequest: mockAuthRequest,
      sessionId: "nonexistent",
      sessions: mockSessions,
      sessionCookies: mockCookies,
    });

    expect(result).toEqual({ error: "Session not found or invalid" });
  });

  it("should return error when cookie not found", async () => {
    vi.mocked(sessionModule.isSessionValid).mockResolvedValue(true);

    const result = await loginWithOIDCAndSession({
      serviceUrl: mockServiceUrl,
      authRequest: mockAuthRequest,
      sessionId: mockSessionId,
      sessions: mockSessions,
      sessionCookies: [],
    });

    expect(result).toEqual({ error: "Session not found or invalid" });
  });

  it("should handle error code 9 with default redirect", async () => {
    const { ConnectError } = await import("@zitadel/client");
    vi.mocked(sessionModule.isSessionValid).mockResolvedValue(true);
    vi.mocked(zitadelModule.createCallback).mockRejectedValue(new ConnectError("already handled", 9));
    vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
      defaultRedirectUri: "https://default.example.com",
    } as any);

    const result = await loginWithOIDCAndSession({
      serviceUrl: mockServiceUrl,
      authRequest: mockAuthRequest,
      sessionId: mockSessionId,
      sessions: mockSessions,
      sessionCookies: mockCookies,
    });

    expect(result).toEqual({ redirect: "https://default.example.com" });
  });

  it("should redirect to /signedin when error code 9 and no default URI", async () => {
    const { ConnectError } = await import("@zitadel/client");
    vi.mocked(sessionModule.isSessionValid).mockResolvedValue(true);
    vi.mocked(zitadelModule.createCallback).mockRejectedValue(new ConnectError("already handled", 9));
    vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({} as any);

    const result = await loginWithOIDCAndSession({
      serviceUrl: mockServiceUrl,
      authRequest: mockAuthRequest,
      sessionId: mockSessionId,
      sessions: mockSessions,
      sessionCookies: mockCookies,
    });

    expect(result).toHaveProperty("redirect");
    if ("redirect" in result) {
      expect(result.redirect).toContain("/signedin");
      expect(result.redirect).toContain("loginName=test%40example.com");
      expect(result.redirect).toContain("organization=org-123");
    }
  });

  it("should return unknown error for non-code-9 errors", async () => {
    vi.mocked(sessionModule.isSessionValid).mockResolvedValue(true);
    vi.mocked(zitadelModule.createCallback).mockRejectedValue({
      code: 13,
      message: "Internal error",
    });

    const result = await loginWithOIDCAndSession({
      serviceUrl: mockServiceUrl,
      authRequest: mockAuthRequest,
      sessionId: mockSessionId,
      sessions: mockSessions,
      sessionCookies: mockCookies,
    });

    expect(result).toEqual({ error: "Unknown error occurred" });
  });
});
