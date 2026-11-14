import { render } from "@testing-library/react";
import { beforeEach, describe, expect, test, vi } from "vitest";
import Page from "./page";

vi.mock("@/lib/server/zitadel", () => ({
  getBrandingSettings: vi.fn(),
  getLoginSettings: vi.fn(),
  getUserByID: vi.fn(),
  searchUsers: vi.fn(),
}));

vi.mock("@/lib/server/session", () => ({
  loadMostRecentSession: vi.fn(),
}));

vi.mock("@/lib/server/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://example.com" } })),
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(() => new Headers()),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(async () => (key: string) => key),
}));

vi.mock("@/components/user-avatar", () => ({ UserAvatar: vi.fn(() => null) }));
vi.mock("@/components/verify-form", () => ({ VerifyForm: vi.fn(() => null) }));
vi.mock("@/components/alert", () => ({ Alert: vi.fn(() => null) }));
vi.mock("@/components/dynamic-theme", () => ({ DynamicTheme: ({ children }: { children: React.ReactNode }) => children }));
vi.mock("@/components/translated", () => ({ Translated: vi.fn(() => null) }));

const knownUser = {
  userId: "user-1",
  preferredLoginName: "anna@example.com@studio.example.com",
  details: { resourceOwner: "org-of-user" },
  type: { case: "human", value: { profile: { displayName: "Anna" } } },
};

async function renderVerify(ignoreUnknownUsernames: boolean, found: boolean) {
  const zitadel = await import("@/lib/server/zitadel");
  vi.mocked(zitadel.getBrandingSettings).mockResolvedValue({ hideLoginNameSuffix: true } as never);
  vi.mocked(zitadel.getLoginSettings).mockResolvedValue({ ignoreUnknownUsernames } as never);
  vi.mocked(zitadel.searchUsers).mockResolvedValue({ result: found ? [knownUser] : [] } as never);

  render(
    await Page({
      searchParams: Promise.resolve({ loginName: "anna@example.com", organization: "org-1" }),
    }),
  );

  const { UserAvatar } = await import("@/components/user-avatar");
  return vi.mocked(UserAvatar).mock.calls[0][0];
}

describe("verify page avatar", () => {
  beforeEach(async () => {
    vi.clearAllMocks();
    const session = await import("@/lib/server/session");
    vi.mocked(session.loadMostRecentSession).mockResolvedValue(undefined);
  });

  test("formats a found user's login name with the requested organization under enumeration protection", async () => {
    const props = await renderVerify(true, true);

    expect(props).toMatchObject({ organizationId: "org-1", hideLoginNameSuffix: true });
  });

  test("formats an unknown login name with the same organization under enumeration protection", async () => {
    const props = await renderVerify(true, false);

    expect(props).toMatchObject({ organizationId: "org-1", hideLoginNameSuffix: true });
  });

  test("formats a found user's login name with the user's organization without enumeration protection", async () => {
    const props = await renderVerify(false, true);

    expect(props).toMatchObject({ organizationId: "org-of-user" });
  });
});
