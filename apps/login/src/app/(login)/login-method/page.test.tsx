import { getLoginMethods } from "@/lib/login-methods";
import { loadMostRecentSession } from "@/lib/session";
import { getLoginSettings, listAuthenticationMethodTypes } from "@/lib/zitadel";
import { cleanup, render, screen } from "@testing-library/react";
import { beforeEach, expect, test, vi } from "vitest";
import Page from "./page";
vi.mock("@/lib/session", () => ({ loadMostRecentSession: vi.fn() }));
vi.mock("@/lib/service-url", () => ({ getServiceConfig: () => ({ serviceConfig: {} }) }));
vi.mock("@/lib/zitadel", () => ({
  getLoginSettings: vi.fn(),
  getBrandingSettings: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
}));
vi.mock("@/lib/login-methods", () => ({ getLoginMethods: vi.fn() }));
vi.mock("next/headers", () => ({ headers: vi.fn() }));
vi.mock("next/navigation", () => ({
  redirect: (url: string) => {
    throw new Error(url);
  },
}));
vi.mock("next-intl/server", () => ({ getTranslations: vi.fn() }));
vi.mock("@/components/dynamic-theme", () => ({ DynamicTheme: ({ children }: any) => <>{children}</> }));
vi.mock("@/components/translated", () => ({ Translated: ({ i18nKey }: any) => <span>{i18nKey}</span> }));
vi.mock("@/components/user-avatar", () => ({ UserAvatar: () => null }));
vi.mock("@/components/sign-in-with-idp", () => ({
  SignInWithIdp: ({ sessionId, loginHint, organization, requestId, identityProviders }: any) => (
    <button data-session={sessionId} data-hint={loginHint} data-org={organization} data-request={requestId}>
      {identityProviders[0].name}
    </button>
  ),
}));
const page = (params = {}) =>
  Page({ searchParams: Promise.resolve({ loginName: "local@example.org", requestId: "req", ...params }) });
beforeEach(() => {
  cleanup();
  vi.clearAllMocks();
  vi.mocked(loadMostRecentSession).mockResolvedValue({
    factors: { user: { id: "user", loginName: "local@example.org", organizationId: "org" } },
  } as any);
  vi.mocked(getLoginSettings).mockResolvedValue({ allowLocalAuthentication: true } as any);
  vi.mocked(listAuthenticationMethodTypes).mockResolvedValue({ authMethodTypes: [] } as any);
  vi.mocked(getLoginMethods).mockResolvedValue({
    password: true,
    passkey: true,
    count: 3,
    identityProviders: [{ id: "google", name: "Google", loginHint: "external@example.org" }],
  } as any);
});
test("renders choices without mounting automatic passkey login or initiating a link", async () => {
  render(await page());
  expect(screen.getByRole("link", { name: /methods.password/ })).toHaveAttribute(
    "href",
    "/password?loginName=local%40example.org&requestId=req&organization=org",
  );
  expect(screen.getByRole("link", { name: /methods.passkey/ })).toHaveAttribute(
    "href",
    "/passkey?loginName=local%40example.org&requestId=req&organization=org",
  );
  const google = screen.getByRole("button", { name: "Google" });
  expect(google).not.toHaveAttribute("data-session");
  expect(google).toHaveAttribute("data-hint", "external@example.org");
  expect(google).toHaveAttribute("data-org", "org");
  expect(google).toHaveAttribute("data-request", "req");
});
test("requires a browser-owned session even for a direct URL", async () => {
  vi.mocked(loadMostRecentSession).mockResolvedValue(undefined);
  await expect(page()).rejects.toThrow("/loginname?");
  expect(getLoginMethods).not.toHaveBeenCalled();
});
test("offers policy-level choices without looking up a protected username", async () => {
  vi.mocked(getLoginSettings).mockResolvedValue({ ignoreUnknownUsernames: true } as any);
  render(await page());
  expect(listAuthenticationMethodTypes).not.toHaveBeenCalled();
  expect(loadMostRecentSession).not.toHaveBeenCalled();
  expect(getLoginMethods).toHaveBeenCalledWith(expect.objectContaining({ userId: undefined, authMethods: undefined }));
});
test("rejects a session for a different login name", async () => {
  await expect(page({ loginName: "other@example.org" })).rejects.toThrow("/loginname?");
  expect(getLoginMethods).not.toHaveBeenCalled();
});
test("derives provider scope from the session user", async () => {
  render(await page({ organization: "org" }));
  expect(getLoginMethods).toHaveBeenCalledWith(expect.objectContaining({ organization: "org", userId: "user" }));
});

test("rejects a session outside the requested organization", async () => {
  await expect(page({ organization: "other-org" })).rejects.toThrow("/loginname?");
  expect(getLoginMethods).not.toHaveBeenCalled();
});
