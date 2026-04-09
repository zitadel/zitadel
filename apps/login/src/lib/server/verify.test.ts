import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { initialSendVerification, sendVerification } from "./verify";

import {
  createInviteCode,
  getSession,
  getUserByID,
  listAuthenticationMethodTypes,
  verifyEmail,
  sendEmailCode as zitadelSendEmailCode,
} from "@/lib/zitadel";
import { cookies } from "next/headers";
import { getSessionCookieByLoginName } from "../cookies";
import { createSessionAndUpdateCookie } from "./cookie";

// Mock dependencies
vi.mock("@/lib/zitadel", () => ({
  verifyEmail: vi.fn(),
  verifyInviteCode: vi.fn(),
  getUserByID: vi.fn(),
  getSession: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
  getLoginSettings: vi.fn(),
  sendEmailCode: vi.fn(),
  createInviteCode: vi.fn(),
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(),
  cookies: vi.fn(),
}));

vi.mock("../cookies", () => ({
  getSessionCookieByLoginName: vi.fn(),
}));

vi.mock("./cookie", () => ({
  createSessionAndUpdateCookie: vi.fn(),
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: {} })),
}));

vi.mock("../fingerprint", () => ({
  getOrSetFingerprintId: vi.fn(),
}));

vi.mock("./host", () => ({
  getPublicHostWithProtocol: vi.fn(() => "https://example.com"),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

describe("sendVerification", () => {
  let mockVerifyEmail: any;
  let mockGetUserByID: any;
  let mockGetSession: any;
  let mockListAuthenticationMethodTypes: any;
  let mockGetSessionCookieByLoginName: any;
  let mockCreateSessionAndUpdateCookie: any;
  let mockCookies: any;

  beforeEach(() => {
    mockVerifyEmail = verifyEmail;
    mockGetUserByID = getUserByID;
    mockGetSession = getSession;
    mockListAuthenticationMethodTypes = listAuthenticationMethodTypes;
    mockGetSessionCookieByLoginName = getSessionCookieByLoginName;
    mockCreateSessionAndUpdateCookie = createSessionAndUpdateCookie;
    mockCookies = cookies;

    // Default valid checking setup
    mockVerifyEmail.mockResolvedValue({});
    mockGetUserByID.mockResolvedValue({
      user: { userId: "user-1", preferredLoginName: "test@example.com" },
    });
    mockCookies.mockResolvedValue({
      set: vi.fn(),
    });
  });

  test("should handle failing getSession when no auth methods exist (stale cookie scenario)", async () => {
    // 1. Cookie exists
    mockGetSessionCookieByLoginName.mockResolvedValue({
      id: "stale-session-id",
      token: "stale-token",
    });

    // 2. getSession fails (throws or returns null/invalid)
    // Simulating call returning nothing useful or throwing
    mockGetSession.mockRejectedValue(new Error("Session not found"));
    // OR: mockGetSession.mockResolvedValue({}); // depending on how generic client behaves

    // 3. No auth methods (user needs setup)
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [],
    });

    // 4. Mock session creation success
    mockCreateSessionAndUpdateCookie.mockResolvedValue({
      session: {
        id: "new-session-id",
        factors: { user: { id: "user-1", loginName: "test@example.com", organizationId: "org-1" } },
      },
    });

    const result = await sendVerification({
      userId: "user-1",
      code: "123456",
      isInvite: false,
    });

    // Expect redirect to authenticator setup
    expect(result).toEqual({
      redirect: "/authenticator/set?sessionId=new-session-id&loginName=test%40example.com",
    });

    // Verify getSession was called (and failed, but we recovered)
    expect(mockGetSession).toHaveBeenCalled();
    // Verify creation was attempted
    expect(mockCreateSessionAndUpdateCookie).toHaveBeenCalled();
  });

  test("should include requestId in /authenticator/set redirect when provided", async () => {
    // 1. No existing session cookie
    mockGetSessionCookieByLoginName.mockResolvedValue(undefined);

    // 2. No auth methods (user needs setup)
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [],
    });

    // 3. Mock session creation success
    mockCreateSessionAndUpdateCookie.mockResolvedValue({
      session: {
        id: "new-session-id",
        factors: { user: { id: "user-1", loginName: "test@example.com", organizationId: "org-1" } },
      },
    });

    const result = await sendVerification({
      userId: "user-1",
      code: "123456",
      isInvite: false,
      requestId: "oidc_auth-req-123",
    });

    // Expect redirect to include requestId
    expect(result).toEqual({
      redirect: expect.stringContaining("/authenticator/set?"),
    });
    const redirectUrl = (result as { redirect: string }).redirect;
    const params = new URLSearchParams(redirectUrl.split("?")[1]);
    expect(params.get("requestId")).toBe("oidc_auth-req-123");
    expect(params.get("sessionId")).toBe("new-session-id");
    expect(params.get("loginName")).toBe("test@example.com");
  });

  test("should NOT include requestId in /authenticator/set redirect when not provided", async () => {
    // 1. No existing session cookie
    mockGetSessionCookieByLoginName.mockResolvedValue(undefined);

    // 2. No auth methods (user needs setup)
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [],
    });

    // 3. Mock session creation success
    mockCreateSessionAndUpdateCookie.mockResolvedValue({
      session: {
        id: "new-session-id",
        factors: { user: { id: "user-1", loginName: "test@example.com", organizationId: "org-1" } },
      },
    });

    const result = await sendVerification({
      userId: "user-1",
      code: "123456",
      isInvite: false,
      // no requestId
    });

    const redirectUrl = (result as { redirect: string }).redirect;
    const params = new URLSearchParams(redirectUrl.split("?")[1]);
    expect(params.get("requestId")).toBeNull();
  });

  test("should fall back to user.preferredLoginName for session cookie lookup when loginName is undefined", async () => {
    // Simulate the email-link scenario: loginName is undefined but userId is provided
    mockGetSessionCookieByLoginName.mockResolvedValue(undefined);
    mockGetUserByID.mockResolvedValue({
      user: { userId: "user-1", preferredLoginName: "test@example.com" },
    });
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [],
    });
    mockCreateSessionAndUpdateCookie.mockResolvedValue({
      session: {
        id: "new-session-id",
        factors: { user: { id: "user-1", loginName: "test@example.com", organizationId: "org-1" } },
      },
    });

    await sendVerification({
      userId: "user-1",
      code: "123456",
      isInvite: false,
      loginName: undefined, // explicitly undefined, like from email link
    });

    // The key assertion: session cookie lookup should use preferredLoginName as fallback
    expect(mockGetSessionCookieByLoginName).toHaveBeenCalledWith({
      loginName: "test@example.com",
      organization: undefined,
    });
  });
});

describe("initialSendVerification", () => {
  let mockSendEmailCode: any;
  let mockCreateInviteCode: any;
  let originalBasePath: string | undefined;

  beforeEach(() => {
    originalBasePath = process.env.NEXT_PUBLIC_BASE_PATH;
    process.env.NEXT_PUBLIC_BASE_PATH = "/ui/v2/login";

    vi.clearAllMocks();
    mockSendEmailCode = zitadelSendEmailCode;
    mockCreateInviteCode = createInviteCode;
    mockSendEmailCode.mockResolvedValue({});
    mockCreateInviteCode.mockResolvedValue({});
  });

  afterEach(() => {
    if (originalBasePath === undefined) {
      delete process.env.NEXT_PUBLIC_BASE_PATH;
    } else {
      process.env.NEXT_PUBLIC_BASE_PATH = originalBasePath;
    }
  });

  test("should call sendEmailCode with correct URL template for non-invite", async () => {
    await initialSendVerification({
      userId: "user-1",
      isInvite: false,
    });

    expect(mockSendEmailCode).toHaveBeenCalledWith({
      serviceConfig: {},
      userId: "user-1",
      urlTemplate: "https://example.com/ui/v2/login/verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}",
    });
    expect(mockCreateInviteCode).not.toHaveBeenCalled();
  });

  test("should call createInviteCode with correct URL template for invite", async () => {
    await initialSendVerification({
      userId: "user-1",
      isInvite: true,
    });

    expect(mockCreateInviteCode).toHaveBeenCalledWith({
      serviceConfig: {},
      userId: "user-1",
      urlTemplate:
        "https://example.com/ui/v2/login/verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}&invite=true",
    });
    expect(mockSendEmailCode).not.toHaveBeenCalled();
  });

  test("should include URL-encoded requestId in URL template", async () => {
    await initialSendVerification({
      userId: "user-1",
      isInvite: false,
      requestId: "req-123",
    });

    expect(mockSendEmailCode).toHaveBeenCalledWith({
      serviceConfig: {},
      userId: "user-1",
      urlTemplate:
        "https://example.com/ui/v2/login/verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}&requestId=req-123",
    });
  });

  test("should URL-encode special characters in requestId", async () => {
    await initialSendVerification({
      userId: "user-1",
      isInvite: false,
      requestId: "req&id=injected",
    });

    expect(mockSendEmailCode).toHaveBeenCalledWith({
      serviceConfig: {},
      userId: "user-1",
      urlTemplate:
        "https://example.com/ui/v2/login/verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}&requestId=req%26id%3Dinjected",
    });
  });

  test("should include invite=true and requestId for invite with requestId", async () => {
    await initialSendVerification({
      userId: "user-1",
      isInvite: true,
      requestId: "req-456",
    });

    expect(mockCreateInviteCode).toHaveBeenCalledWith({
      serviceConfig: {},
      userId: "user-1",
      urlTemplate:
        "https://example.com/ui/v2/login/verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}&invite=true&requestId=req-456",
    });
  });
});
