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
});
