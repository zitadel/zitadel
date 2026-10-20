import { cleanup, render } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import Page from "./page";

const { redirect } = vi.hoisted(() => ({
  redirect: vi.fn((url: string) => {
    throw new Error(`NEXT_REDIRECT ${url}`);
  }),
}));

vi.mock("next/navigation", () => ({ redirect }));

vi.mock("@/lib/zitadel", () => ({
  getBrandingSettings: vi.fn(),
  getDefaultOrg: vi.fn(),
  getDeviceAuthorizationRequest: vi.fn(),
}));

vi.mock("@/lib/session", () => ({
  loadMostRecentSession: vi.fn(),
  loadSessionById: vi.fn(),
}));

vi.mock("@/lib/device-code-check", () => ({
  checkDeviceCodeCheckCookie: vi.fn(),
}));

vi.mock("@/lib/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://example.com" } })),
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(() => new Headers()),
}));

vi.mock("@/components/consent", () => ({ ConsentScreen: vi.fn(() => null) }));
vi.mock("@/components/user-avatar", () => ({ UserAvatar: vi.fn(() => null) }));
vi.mock("@/components/dynamic-theme", () => ({ DynamicTheme: ({ children }: { children: React.ReactNode }) => children }));
vi.mock("@/components/translated", () => ({ Translated: vi.fn(() => null) }));

const session = {
  id: "session-1",
  factors: { user: { id: "user-1", loginName: "user@example.com", displayName: "User", organizationId: "org-1" } },
};

describe("device consent page", () => {
  let checkDeviceCodeCheckCookie: any;
  let loadSessionById: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const zitadel = await import("@/lib/zitadel");
    const sessionLib = await import("@/lib/session");
    const deviceCodeCheck = await import("@/lib/device-code-check");

    checkDeviceCodeCheckCookie = vi.mocked(deviceCodeCheck.checkDeviceCodeCheckCookie);
    loadSessionById = vi.mocked(sessionLib.loadSessionById);

    checkDeviceCodeCheckCookie.mockResolvedValue(true);
    loadSessionById.mockResolvedValue(session);
    vi.mocked(sessionLib.loadMostRecentSession).mockResolvedValue(session as any);
    vi.mocked(zitadel.getDeviceAuthorizationRequest).mockResolvedValue({
      deviceAuthorizationRequest: { id: "device-request-1", appName: "My Device", scope: ["openid"] },
    } as any);
    vi.mocked(zitadel.getBrandingSettings).mockResolvedValue({} as any);
    vi.mocked(zitadel.getDefaultOrg).mockResolvedValue(null as any);
  });

  afterEach(cleanup);

  test("redirects to the code entry page when the device code was not entered in this browser", async () => {
    checkDeviceCodeCheckCookie.mockResolvedValue(false);

    await expect(
      Page({ searchParams: Promise.resolve({ requestId: "device_ABCD-1234", sessionId: "session-1" }) }),
    ).rejects.toThrow("NEXT_REDIRECT");

    expect(checkDeviceCodeCheckCookie).toHaveBeenCalledWith("ABCD-1234");
    expect(redirect).toHaveBeenCalledWith("/device?user_code=ABCD-1234");
    expect(loadSessionById).not.toHaveBeenCalled();
  });

  test("redirects to the login funnel when the session is stale", async () => {
    loadSessionById.mockResolvedValue(undefined);

    await expect(
      Page({ searchParams: Promise.resolve({ requestId: "device_ABCD-1234", sessionId: "session-1" }) }),
    ).rejects.toThrow("NEXT_REDIRECT");

    expect(redirect).toHaveBeenCalledWith("/loginname?requestId=device_ABCD-1234");
  });

  test("passes the displayed session to the consent screen", async () => {
    const { ConsentScreen } = await import("@/components/consent");
    const { UserAvatar } = await import("@/components/user-avatar");

    render(
      await Page({
        searchParams: Promise.resolve({
          requestId: "device_ABCD-1234",
          sessionId: "session-1",
          loginName: "other@example.com",
          organization: "org-1",
        }),
      }),
    );

    // React 19 calls function components with (props, undefined)
    expect(ConsentScreen).toHaveBeenCalledWith(
      expect.objectContaining({ userCode: "ABCD-1234", sessionId: "session-1", organization: "org-1" }),
      undefined,
    );
    expect(ConsentScreen).toHaveBeenCalledWith(expect.not.objectContaining({ loginName: expect.anything() }), undefined);
    expect(UserAvatar).toHaveBeenCalledWith(expect.objectContaining({ loginName: "user@example.com" }), undefined);
  });
});
