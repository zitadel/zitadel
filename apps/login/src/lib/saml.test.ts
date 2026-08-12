import { beforeEach, describe, expect, it, vi } from "vitest";
import { loginWithSAMLAndSession } from "./saml";
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
  create: vi.fn((_schema: unknown, data: unknown) => data),
}));

describe("loginWithSAMLAndSession", () => {
  const mockSamlRequest = "saml-123";
  const mockSessionId = "session-123";

  let mockSessions: any[];
  let mockCookies: any[];

  beforeEach(() => {
    vi.clearAllMocks();
    vi.spyOn(console, "log").mockImplementation(() => {});
    vi.spyOn(console, "warn").mockImplementation(() => {});
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

  it("should return the redirect binding when the session is valid", async () => {
    vi.mocked(sessionModule.isSessionValid).mockResolvedValue(true);
    vi.mocked(zitadelModule.createResponse).mockResolvedValue({
      url: "https://sp.example.com/acs",
      binding: { case: "redirect", value: {} },
    } as any);

    const result = await loginWithSAMLAndSession({
      serviceConfig: {} as any,
      samlRequest: mockSamlRequest,
      sessionId: mockSessionId,
      sessions: mockSessions,
      sessionCookies: mockCookies,
    });

    expect(result).toEqual({ redirect: "https://sp.example.com/acs" });
  });

  it("should rethrow when session validation fails with a server error", async () => {
    const { ConnectError } = await import("@zitadel/client");
    const serverError = new ConnectError("internal failure", 13) as any;
    serverError.isUserError = false;
    vi.mocked(sessionModule.isSessionValid).mockRejectedValue(serverError);

    await expect(
      loginWithSAMLAndSession({
        serviceConfig: {} as any,
        samlRequest: mockSamlRequest,
        sessionId: mockSessionId,
        sessions: mockSessions,
        sessionCookies: mockCookies,
      }),
    ).rejects.toThrow("internal failure");
    expect(zitadelModule.createResponse).not.toHaveBeenCalled();
  });

  it("should degrade to re-authentication when session validation fails with a user error", async () => {
    const { ConnectError } = await import("@zitadel/client");
    const userError = new ConnectError("user was removed", 5) as any;
    userError.isUserError = true;
    vi.mocked(sessionModule.isSessionValid).mockRejectedValue(userError);
    vi.mocked(loginnameModule.sendLoginname).mockResolvedValue({
      redirect: "/password",
    });

    const result = await loginWithSAMLAndSession({
      serviceConfig: {} as any,
      samlRequest: mockSamlRequest,
      sessionId: mockSessionId,
      sessions: mockSessions,
      sessionCookies: mockCookies,
    });

    expect(result).toEqual({ redirect: "/password" });
    expect(loginnameModule.sendLoginname).toHaveBeenCalledWith({
      loginName: "test@example.com",
      organization: "org-123",
      requestId: `saml_${mockSamlRequest}`,
    });
  });
});
