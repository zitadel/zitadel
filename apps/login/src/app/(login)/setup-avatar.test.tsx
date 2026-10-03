import { isValidElement, ReactNode } from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

const { getUserByID, avatarComponent } = vi.hoisted(() => ({
  getUserByID: vi.fn(),
  avatarComponent: vi.fn(),
}));

vi.mock("@/components/user-avatar", () => ({ UserAvatar: avatarComponent }));
vi.mock("@/components/alert", () => ({ Alert: "div" }));
vi.mock("@/components/back-button", () => ({ BackButton: "button" }));
vi.mock("@/components/choose-authenticator-to-setup", () => ({ ChooseAuthenticatorToSetup: "div" }));
vi.mock("@/components/choose-second-factor-to-setup", () => ({ ChooseSecondFactorToSetup: "div" }));
vi.mock("@/components/dynamic-theme", () => ({ DynamicTheme: "div" }));
vi.mock("@/components/sign-in-with-idp", () => ({ SignInWithIdp: "div" }));
vi.mock("@/components/translated", () => ({ Translated: "span" }));
vi.mock("@/lib/cookies", () => ({ getSessionCookieById: vi.fn() }));
vi.mock("@/lib/service-url", () => ({ getServiceConfig: () => ({ serviceConfig: {} }) }));
vi.mock("@/lib/session", () => ({
  loadMostRecentSession: async () => ({
    id: "session-1",
    factors: { user: { id: "user-1", loginName: "sam@example.com", organizationId: "org-1" } },
  }),
  hasVerifiedPrimaryFactor: () => ({ valid: true }),
}));
vi.mock("@/lib/verify-helper", () => ({ checkUserVerification: async () => true }));
vi.mock("@/lib/zitadel", () => ({
  getUserByID,
  listAuthenticationMethodTypes: async () => ({ authMethodTypes: [] }),
  getBrandingSettings: async () => ({}),
  getLoginSettings: async () => ({ secondFactors: [] }),
  getActiveIdentityProviders: async () => ({ identityProviders: [] }),
  getSession: vi.fn(),
}));
vi.mock("next/headers", () => ({ headers: async () => new Headers() }));
vi.mock("next-intl/server", () => ({ getTranslations: async () => (key: string) => key }));
vi.mock("next/navigation", () => ({ redirect: vi.fn() }));

import AuthenticatorPage from "./authenticator/set/page";
import MfaPage from "./mfa/set/page";

function findAvatar(node: ReactNode): { imageUrl?: string } | undefined {
  if (Array.isArray(node)) {
    return node.map(findAvatar).find(Boolean);
  }
  if (!isValidElement<{ children?: ReactNode; imageUrl?: string }>(node)) return;
  if (node.type === avatarComponent) return node.props;
  return findAvatar(node.props.children);
}

describe.each([
  ["authenticator setup", AuthenticatorPage],
  ["MFA setup", MfaPage],
])("%s avatar", (_name, Page) => {
  beforeEach(() => vi.clearAllMocks());

  it.each(["https://example.com/photo.png", undefined])("reuses the loaded profile photo (%s)", async (avatarUrl) => {
    getUserByID.mockResolvedValue({
      user: { type: { case: "human", value: { profile: { avatarUrl } } } },
    });
    const page = await Page({ searchParams: Promise.resolve({ loginName: "sam@example.com" }) });
    expect(findAvatar(page)).toEqual(expect.objectContaining({ imageUrl: avatarUrl }));
    expect(getUserByID).toHaveBeenCalledExactlyOnceWith({ serviceConfig: {}, userId: "user-1" });
  });
});
