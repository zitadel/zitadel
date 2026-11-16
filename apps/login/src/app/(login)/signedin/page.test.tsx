import { beforeEach, describe, expect, test, vi } from "vitest";
import Page from "./page";

vi.mock("@/lib/server/zitadel", () => ({
  getBrandingSettings: vi.fn(),
  getLoginSettings: vi.fn(),
  authorizeOrDenyDeviceAuthorization: vi.fn(),
}));

vi.mock("@/lib/actions/device", () => ({
  approveDeviceAuthorization: vi.fn(),
  denyDeviceAuthorization: vi.fn(),
  completeDeviceAuthorization: vi.fn(),
}));

vi.mock("@/lib/server/session", () => ({
  loadMostRecentSession: vi.fn(),
  loadSessionById: vi.fn(),
}));

vi.mock("@/lib/server/next-url", () => ({
  resolveRedirectUri: vi.fn(),
}));

vi.mock("@/lib/server/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://example.com" } })),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(async () => (key: string) => key),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(() => new Headers()),
}));

const session = {
  id: "session-1",
  factors: { user: { id: "user-1", loginName: "user@example.com", displayName: "User", organizationId: "org-1" } },
};

describe("signedin page", () => {
  let authorizeOrDenyDeviceAuthorization: any;
  let completeDeviceAuthorization: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const zitadel = await import("@/lib/server/zitadel");
    const device = await import("@/lib/actions/device");
    const sessionLib = await import("@/lib/server/session");
    const client = await import("@/lib/server/next-url");

    authorizeOrDenyDeviceAuthorization = vi.mocked(zitadel.authorizeOrDenyDeviceAuthorization);
    completeDeviceAuthorization = vi.mocked((device as any).completeDeviceAuthorization);

    vi.mocked(zitadel.getBrandingSettings).mockResolvedValue({} as any);
    vi.mocked(zitadel.getLoginSettings).mockResolvedValue({} as any);
    vi.mocked(sessionLib.loadMostRecentSession).mockResolvedValue(session as any);
    vi.mocked(sessionLib.loadSessionById).mockResolvedValue(session as any);
    vi.mocked(client.resolveRedirectUri).mockResolvedValue("/signedin");
  });

  // Regression test for the device authorization consent bypass: rendering this page with an
  // ambient session cookie must never complete a device authorization grant.
  test("does not complete a device authorization when rendered with a device requestId", async () => {
    const result = await Page({ searchParams: Promise.resolve({ requestId: "device_ABCD-1234" }) });

    expect(result).toBeTruthy();
    expect(authorizeOrDenyDeviceAuthorization).not.toHaveBeenCalled();
    expect(completeDeviceAuthorization).not.toHaveBeenCalled();
  });

  test("does not complete a device authorization when a sessionId is given", async () => {
    const result = await Page({
      searchParams: Promise.resolve({ requestId: "device_ABCD-1234", sessionId: "session-1", organization: "org-1" }),
    });

    expect(result).toBeTruthy();
    expect(authorizeOrDenyDeviceAuthorization).not.toHaveBeenCalled();
    expect(completeDeviceAuthorization).not.toHaveBeenCalled();
  });
});
