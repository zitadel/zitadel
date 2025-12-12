import { describe, expect, test, vi, beforeEach } from "vitest";
import { sendPassword } from "./password";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { UserState } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { create } from "@zitadel/client";

// Mock dependencies
vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("@zitadel/client", () => ({
  create: vi.fn((schema, data) => data),
  ConnectError: class extends Error {
    code: number;
    constructor(msg: string, code: number) {
      super(msg);
      this.code = code;
    }
  },
  timestampDate: (ts: any) => new Date(ts.seconds * 1000),
}));

vi.mock("@zitadel/client/v2", () => ({
  createUserServiceClient: vi.fn(),
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("../zitadel", () => ({
  getLoginSettings: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
  getSession: vi.fn(),
  setPassword: vi.fn(),
  createServerTransport: vi.fn(),
  listUsers: vi.fn(),
  getLockoutSettings: vi.fn(),
  getPasswordExpirySettings: vi.fn(),
  getUserByID: vi.fn(),
}));

vi.mock("../cookies", () => ({
  getSessionCookieByLoginName: vi.fn(),
  getSessionCookieById: vi.fn(),
}));

vi.mock("../client", () => ({
  completeFlowOrGetUrl: vi.fn(),
}));

vi.mock("@/lib/server/cookie", () => ({
  createSessionAndUpdateCookie: vi.fn(),
  setSessionAndUpdateCookie: vi.fn(),
}));

vi.mock("../verify-helper", () => ({
  checkPasswordChangeRequired: vi.fn(),
  checkEmailVerification: vi.fn(),
  checkMFAFactors: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

describe("sendPassword", () => {
  let mockHeaders: any;
  let mockGetServiceConfig: any;
  let mockGetSessionCookieByLoginName: any;
  let mockListUsers: any;
  let mockGetLoginSettings: any;
  let mockCreateSessionAndUpdateCookie: any;
  let mockCompleteFlowOrGetUrl: any;
  let mockListAuthenticationMethodTypes: any;
  let mockCheckMFAFactors: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { getSessionCookieByLoginName } = await import("../cookies");
    const { listUsers, getLoginSettings, listAuthenticationMethodTypes } = await import("../zitadel");
    const { createSessionAndUpdateCookie } = await import("@/lib/server/cookie");
    const { completeFlowOrGetUrl } = await import("../client");
    const { checkMFAFactors } = await import("../verify-helper");

    mockHeaders = vi.mocked(headers);
    mockGetServiceConfig = vi.mocked(getServiceConfig);
    mockGetSessionCookieByLoginName = vi.mocked(getSessionCookieByLoginName);
    mockListUsers = vi.mocked(listUsers);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockCreateSessionAndUpdateCookie = vi.mocked(createSessionAndUpdateCookie);
    mockCompleteFlowOrGetUrl = vi.mocked(completeFlowOrGetUrl);
    mockListAuthenticationMethodTypes = vi.mocked(listAuthenticationMethodTypes);
    mockCheckMFAFactors = vi.mocked(checkMFAFactors);

    mockHeaders.mockResolvedValue({});
    mockGetServiceConfig.mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } });
  });

  test("should create session and verify password when no session exists", async () => {
    // 1. No existing session
    mockGetSessionCookieByLoginName.mockResolvedValue(null);

    // 2. User exists
    mockListUsers.mockResolvedValue({
      details: { totalResult: BigInt(1) },
      result: [{ userId: "user123", type: { case: "human", value: {} }, state: UserState.ACTIVE }],
    });

    // 3. Login settings
    mockGetLoginSettings.mockResolvedValue({ passwordCheckLifetime: { seconds: BigInt(86400) } });

    // 4. Create session success
    mockCreateSessionAndUpdateCookie.mockResolvedValue({
      session: {
        id: "new-session-id",
        factors: {
          user: { id: "user123", loginName: "testuser", organizationId: "org123" },
          password: { verifiedAt: { seconds: 100 } },
        },
      },
      sessionCookie: {
        id: "new-session-id",
        token: "token123",
        organization: "org123",
      },
    });

    // 5. Auth methods (no MFA)
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [AuthenticationMethodType.PASSWORD],
    });

    // 6. MFA check passes
    mockCheckMFAFactors.mockResolvedValue(null);

    // 7. Redirect
    mockCompleteFlowOrGetUrl.mockResolvedValue({ redirect: "/next-page" });

    const result = await sendPassword({
      loginName: "testuser",
      checks: create(ChecksSchema, { password: { password: "password123" } }),
    });

    // Verification
    expect(mockListUsers).toHaveBeenCalledWith(expect.objectContaining({ loginName: "testuser" }));
    expect(mockCreateSessionAndUpdateCookie).toHaveBeenCalled();
    expect(mockCompleteFlowOrGetUrl).toHaveBeenCalled();
    expect(result).toEqual({ redirect: "/next-page" });
  });
});
