import { cleanup, render } from "@testing-library/react";
import { IdentityProviderType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { ReactNode } from "react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import Page from "./page";

const mocks = vi.hoisted(() => ({
  redirect: vi.fn((url: string) => {
    // like next/navigation, redirect() never returns
    throw new Error(`NEXT_REDIRECT:${url}`);
  }),
  getDefaultOrg: vi.fn(),
  getLoginSettings: vi.fn(),
  getActiveIdentityProviders: vi.fn(),
  getBrandingSettings: vi.fn(),
  startIdentityProviderFlow: vi.fn(),
}));

vi.mock("next/navigation", () => ({ redirect: mocks.redirect }));
vi.mock("next/headers", () => ({ headers: vi.fn().mockResolvedValue(new Headers()) }));
vi.mock("next-intl/server", () => ({ getTranslations: vi.fn() }));

vi.mock("@/lib/service-url", () => ({
  getServiceConfig: () => ({ serviceConfig: { baseUrl: "https://api.example.com" } }),
}));

vi.mock("@/lib/server/host", () => ({
  getPublicHost: () => "login.example.com",
}));

vi.mock("@/lib/zitadel", () => ({
  getDefaultOrg: mocks.getDefaultOrg,
  getLoginSettings: mocks.getLoginSettings,
  getActiveIdentityProviders: mocks.getActiveIdentityProviders,
  getBrandingSettings: mocks.getBrandingSettings,
  startIdentityProviderFlow: mocks.startIdentityProviderFlow,
}));

vi.mock("@/components/dynamic-theme", () => ({
  DynamicTheme: ({ children }: { children: ReactNode }) => <div>{children}</div>,
}));
vi.mock("@/components/translated", () => ({
  Translated: ({ i18nKey }: { i18nKey: string }) => <span>{i18nKey}</span>,
}));
vi.mock("@/components/username-form", () => ({
  UsernameForm: () => <div data-testid="username-form" />,
}));
vi.mock("@/components/sign-in-with-idp", () => ({
  SignInWithIdp: () => <div data-testid="sign-in-with-idp" />,
}));
vi.mock("@/components/auto-submit-form", () => ({
  AutoSubmitForm: ({ url, fields }: { url: string; fields: Record<string, string> }) => (
    <div data-testid="auto-submit-form" data-url={url} data-fields={JSON.stringify(fields)} />
  ),
}));

const oidcIdp = { id: "idp-oidc", name: "OIDC", type: IdentityProviderType.OIDC };
const samlIdp = { id: "idp-saml", name: "SAML", type: IdentityProviderType.SAML };
const ldapIdp = { id: "idp-ldap", name: "LDAP", type: IdentityProviderType.LDAP };

function setup({
  allowLocalAuthentication = false,
  allowExternalIdp = true,
  identityProviders = [oidcIdp] as unknown[],
} = {}) {
  mocks.getLoginSettings.mockResolvedValue({ allowLocalAuthentication, allowExternalIdp });
  mocks.getActiveIdentityProviders.mockResolvedValue({ identityProviders });
}

function renderPage(searchParams: Record<string, string> = {}) {
  return Page({ searchParams: Promise.resolve(searchParams) });
}

describe("loginname page", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    delete process.env.NEXT_PUBLIC_BASE_PATH;
    mocks.getDefaultOrg.mockResolvedValue({ id: "default-org" });
    mocks.getBrandingSettings.mockResolvedValue({});
    mocks.startIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/authorize?x=1" });
  });

  afterEach(cleanup);

  describe("single IdP and local authentication disabled", () => {
    test("should redirect to the IdP without rendering the page", async () => {
      setup();

      await expect(
        renderPage({ requestId: "oidc_123", organization: "org-1", loginName: "jane", orgDomain: "acme.com" }),
      ).rejects.toThrow("NEXT_REDIRECT:https://idp.example.com/authorize?x=1");

      const params = "requestId=oidc_123&organization=org-1&postErrorRedirectUrl=%2Floginname";
      expect(mocks.startIdentityProviderFlow).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        idpId: "idp-oidc",
        urls: {
          successUrl: `https://login.example.com/idp/oidc/process?${params}`,
          failureUrl: `https://login.example.com/idp/oidc/failure?${params}`,
          loginHint: "jane@acme.com",
        },
      });
    });

    test("should include the base path in the IdP callback URLs", async () => {
      process.env.NEXT_PUBLIC_BASE_PATH = "/ui/v2/login";
      setup();

      await expect(renderPage()).rejects.toThrow("NEXT_REDIRECT:");

      const { urls } = mocks.startIdentityProviderFlow.mock.calls[0][0];
      expect(urls.successUrl).toBe(
        "https://login.example.com/ui/v2/login/idp/oidc/process?postErrorRedirectUrl=%2Floginname",
      );
      expect(urls.failureUrl).toBe(
        "https://login.example.com/ui/v2/login/idp/oidc/failure?postErrorRedirectUrl=%2Floginname",
      );
    });

    test("should render an auto-submitting form for SAML POST binding", async () => {
      setup({ identityProviders: [samlIdp] });
      const fields = { SAMLRequest: "req", RelayState: "state" };
      mocks.startIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/saml", fields });

      const { getByTestId, queryByTestId } = render(await renderPage());

      expect(mocks.redirect).not.toHaveBeenCalled();
      const form = getByTestId("auto-submit-form");
      expect(form).toHaveAttribute("data-url", "https://idp.example.com/saml");
      expect(form).toHaveAttribute("data-fields", JSON.stringify(fields));
      expect(queryByTestId("sign-in-with-idp")).not.toBeInTheDocument();
    });

    test("should redirect to the LDAP page without starting an IdP flow", async () => {
      setup({ identityProviders: [ldapIdp] });

      await expect(renderPage({ requestId: "oidc_123" })).rejects.toThrow(
        "NEXT_REDIRECT:/idp/ldap?requestId=oidc_123&postErrorRedirectUrl=%2Floginname&idpId=idp-ldap",
      );
      expect(mocks.startIdentityProviderFlow).not.toHaveBeenCalled();
    });

    test("should render the IdP button when the IdP flow cannot be started", async () => {
      setup();
      mocks.startIdentityProviderFlow.mockRejectedValue(new Error("unavailable"));

      const { getByTestId } = render(await renderPage());

      expect(mocks.redirect).not.toHaveBeenCalled();
      expect(getByTestId("sign-in-with-idp")).toBeInTheDocument();
    });

    test("should render the IdP button when the IdP flow returns no URL", async () => {
      setup();
      mocks.startIdentityProviderFlow.mockResolvedValue(null);

      const { getByTestId } = render(await renderPage());

      expect(mocks.redirect).not.toHaveBeenCalled();
      expect(getByTestId("sign-in-with-idp")).toBeInTheDocument();
    });

  describe("no automatic redirect", () => {
    test.each([
      ["local authentication is enabled", { allowLocalAuthentication: true }],
      ["more than one IdP is active", { identityProviders: [oidcIdp, samlIdp] }],
      ["external IdPs are not allowed", { allowExternalIdp: false }],
      ["no IdP is active", { identityProviders: [] }],
    ])("should render the page when %s", async (_, settings) => {
      setup(settings);

      render(await renderPage());

      expect(mocks.startIdentityProviderFlow).not.toHaveBeenCalled();
      expect(mocks.redirect).not.toHaveBeenCalled();
    });
  });
});
