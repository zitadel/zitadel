import { loadMostRecentSession } from "@/lib/session";
import { checkUserVerification } from "@/lib/verify-helper";
import { getActiveIdentityProviders, listAuthenticationMethodTypes } from "@/lib/zitadel";
import { cleanup, render, screen } from "@testing-library/react";
import { redirect } from "next/navigation";
import { beforeEach, describe, expect, test, vi } from "vitest";
import Page from "./page";

vi.mock("@zitadel/client", () => ({ timestampDate: (ts: { seconds: bigint }) => new Date(Number(ts.seconds) * 1000) }));
vi.mock("@/lib/session", async (original) => ({ ...(await original<object>()), loadMostRecentSession: vi.fn() }));
vi.mock("@/lib/cookies", () => ({ getSessionCookieById: vi.fn() }));
vi.mock("@/lib/verify-helper", () => ({ checkUserVerification: vi.fn() }));
vi.mock("@/lib/service-url", () => ({ getServiceConfig: () => ({ serviceConfig: {} }) }));
vi.mock("@/lib/zitadel", () => ({
  getActiveIdentityProviders: vi.fn(),
  getBrandingSettings: vi.fn(),
  getLoginSettings: vi.fn().mockResolvedValue({ allowExternalIdp: true, allowRegister: false }),
  getSession: vi.fn(),
  getUserByID: vi.fn().mockResolvedValue({ user: { type: { case: "human", value: {} } } }),
  listAuthenticationMethodTypes: vi.fn(),
}));
vi.mock("next/headers", () => ({ headers: vi.fn() }));
vi.mock("next/navigation", () => ({
  redirect: vi.fn((url: string) => {
    throw new Error(url);
  }),
}));
vi.mock("next-intl/server", () => ({ getTranslations: vi.fn() }));
vi.mock("@/components/dynamic-theme", () => ({ DynamicTheme: ({ children }: any) => <>{children}</> }));
vi.mock("@/components/alert", () => ({ Alert: ({ children }: any) => <div>{children}</div> }));
vi.mock("@/components/translated", () => ({ Translated: ({ i18nKey }: any) => <span>{i18nKey}</span> }));
vi.mock("@/components/user-avatar", () => ({ UserAvatar: () => null }));
vi.mock("@/components/back-button", () => ({ BackButton: () => null }));
vi.mock("@/components/choose-authenticator-to-setup", () => ({ ChooseAuthenticatorToSetup: () => null }));
vi.mock("@/components/sign-in-with-idp", () => ({
  SignInWithIdp: ({ sessionId }: any) => <button data-session={sessionId}>Google</button>,
}));

const future = { seconds: BigInt(Math.floor(Date.now() / 1000) + 3600), nanos: 0 };
const past = { seconds: BigInt(Math.floor(Date.now() / 1000) - 3600), nanos: 0 };
const user = { id: "local-user", loginName: "local@example.org", organizationId: "org" };
const session = { id: "local-session", factors: { user, password: { verifiedAt: future } }, expirationDate: future };

beforeEach(() => {
  cleanup();
  vi.clearAllMocks();
  vi.mocked(loadMostRecentSession).mockResolvedValue(session as any);
  vi.mocked(listAuthenticationMethodTypes).mockResolvedValue({ authMethodTypes: [1] } as any);
  vi.mocked(checkUserVerification).mockResolvedValue(false);
  vi.mocked(getActiveIdentityProviders).mockResolvedValue({ identityProviders: [{ id: "google" }] } as any);
});

async function page() {
  return Page({ searchParams: Promise.resolve({ loginName: user.loginName }) });
}

describe("existing-account IDP enrollment", () => {
  test("authenticated password session can link Google with registration closed and no invitation cookie", async () => {
    render(await page());
    expect(screen.getByRole("button", { name: "Google" })).toHaveAttribute("data-session", session.id);
    expect(checkUserVerification).not.toHaveBeenCalled();
    expect(redirect).not.toHaveBeenCalled();
    expect(getActiveIdentityProviders).toHaveBeenCalledWith(
      expect.objectContaining({ orgId: "org", linking_allowed: true }),
    );
  });

  test("identify-only existing account cannot enroll or request an invitation", async () => {
    vi.mocked(loadMostRecentSession).mockResolvedValue({ ...session, factors: { user } } as any);
    render(await page());
    expect(screen.getByText("sessionExpired")).toBeInTheDocument();
    expect(getActiveIdentityProviders).not.toHaveBeenCalled();
    expect(redirect).not.toHaveBeenCalled();
  });

  test("expired authenticated account cannot enroll", async () => {
    vi.mocked(loadMostRecentSession).mockResolvedValue({ ...session, expirationDate: past } as any);
    render(await page());
    expect(screen.getByText("sessionExpired")).toBeInTheDocument();
    expect(getActiveIdentityProviders).not.toHaveBeenCalled();
  });

  test("new account without prior verification still follows invitation onboarding", async () => {
    vi.mocked(loadMostRecentSession).mockResolvedValue({ ...session, factors: { user } } as any);
    vi.mocked(listAuthenticationMethodTypes).mockResolvedValue({ authMethodTypes: [] } as any);
    await expect(page()).rejects.toThrow("/verify?");
    expect(redirect).toHaveBeenCalledWith(expect.stringContaining("invite=true"));
    expect(getActiveIdentityProviders).not.toHaveBeenCalled();
  });

  test("new account with prior verification can still enroll", async () => {
    vi.mocked(loadMostRecentSession).mockResolvedValue({ ...session, factors: { user } } as any);
    vi.mocked(listAuthenticationMethodTypes).mockResolvedValue({ authMethodTypes: [] } as any);
    vi.mocked(checkUserVerification).mockResolvedValue(true);
    render(await page());
    expect(screen.getByRole("button", { name: "Google" })).toHaveAttribute("data-session", session.id);
  });
});
