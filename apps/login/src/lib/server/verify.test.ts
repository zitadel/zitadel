import { vi, describe, expect, test, beforeEach } from "vitest";
import { sendVerification } from "./verify";

// Mock dependencies
vi.mock("@/lib/zitadel", () => ({
  verifyEmail: vi.fn(),
  verifyInviteCode: vi.fn(),
  getUserByID: vi.fn(),
  getSession: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
  getLoginSettings: vi.fn(),
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

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

import { verifyEmail, getUserByID, getSession, listAuthenticationMethodTypes } from "@/lib/zitadel";
import { getSessionCookieByLoginName } from "../cookies";
import { createSessionAndUpdateCookie } from "./cookie";
import { cookies } from "next/headers";

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
});
