import { cleanup, render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import Page from "./page";

vi.mock("@/lib/zitadel", () => ({
  getBrandingSettings: vi.fn(),
  getLoginSettings: vi.fn(),
}));

vi.mock("@/lib/session", async () => {
  const actual = await vi.importActual<typeof import("@/lib/session")>("@/lib/session");
  return {
    resolveMostRecentSession: vi.fn(),
    resolveSessionById: vi.fn(),
    SESSION_RESOLUTION_MESSAGE_KEYS: actual.SESSION_RESOLUTION_MESSAGE_KEYS,
  };
});

vi.mock("@/lib/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://example.com" } })),
}));

vi.mock("@/lib/server/host", () => ({
  getPublicHost: vi.fn(() => "example.com"),
}));

vi.mock("@/components/dynamic-theme", () => ({
  DynamicTheme: ({ children }: any) => <div>{children}</div>,
}));

vi.mock("@/components/alert", () => ({
  Alert: ({ children }: any) => <div role="alert">{children}</div>,
}));

vi.mock("@/components/translated", () => ({
  Translated: ({ i18nKey }: any) => <span>{i18nKey}</span>,
}));

vi.mock("@/components/user-avatar", () => ({
  UserAvatar: () => <div data-testid="user-avatar" />,
}));

vi.mock("@/components/login-otp", () => ({
  LoginOTP: () => <div data-testid="login-otp" />,
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(async () => (key: string) => key),
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(() => new Headers()),
}));

describe("otp page", () => {
  let resolveSessionById: any;
  let resolveMostRecentSession: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const zitadel = await import("@/lib/zitadel");
    const sessionLib = await import("@/lib/session");

    resolveSessionById = vi.mocked(sessionLib.resolveSessionById);
    resolveMostRecentSession = vi.mocked(sessionLib.resolveMostRecentSession);
    vi.mocked(zitadel.getBrandingSettings).mockResolvedValue({} as any);
    vi.mocked(zitadel.getLoginSettings).mockResolvedValue({} as any);
  });

  afterEach(cleanup);

  async function renderPage(searchParams: Record<string, string>) {
    const result = await Page({
      searchParams: Promise.resolve(searchParams),
      params: Promise.resolve({ method: "time-based" }),
    });
    render(result);
  }

  // A session that cannot be resolved must not crash the render: the page shows the alert for the
  // reason on every branch.
  test.each([
    {
      name: "sessionId and loginName",
      searchParams: { sessionId: "session-1", loginName: "user@example.com", organization: "org-1" },
      expectedAlert: "sessionResolution.sessionNotFound",
    },
    {
      name: "loginName only",
      searchParams: { loginName: "user@example.com", organization: "org-1" },
      expectedAlert: "sessionResolution.sessionNotFound",
    },
    {
      name: "neither sessionId nor loginName",
      searchParams: {},
      expectedAlert: "unknownContext",
    },
  ])(
    "renders the $expectedAlert alert and no OTP form when the session cannot be loaded ($name)",
    async ({ searchParams, expectedAlert }) => {
      resolveSessionById.mockResolvedValue({ reason: "server_session_not_found" });
      resolveMostRecentSession.mockResolvedValue({ reason: "server_session_not_found" });

      await renderPage(searchParams as Record<string, string>);

      if (searchParams.sessionId) {
        expect(resolveSessionById).toHaveBeenCalledWith(
          { baseUrl: "https://example.com" },
          "session-1",
          "org-1",
          "otp_page",
        );
        expect(resolveMostRecentSession).not.toHaveBeenCalled();
      } else {
        expect(resolveMostRecentSession).toHaveBeenCalledWith(
          {
            serviceConfig: { baseUrl: "https://example.com" },
            sessionParams: { loginName: searchParams.loginName, organization: searchParams.organization },
          },
          "otp_page",
        );
        expect(resolveSessionById).not.toHaveBeenCalled();
      }
      expect(screen.getByRole("alert").textContent).toBe(expectedAlert);
      expect(screen.queryByTestId("login-otp")).toBeNull();
    },
  );

  test.each([
    ["cookie_missing", "sessionResolution.cookieMissing"],
    ["cookie_invalid", "sessionResolution.cookieInvalid"],
    ["session_entry_not_found", "sessionResolution.entryNotFound"],
    ["server_session_not_found", "sessionResolution.sessionNotFound"],
    ["session_token_rejected", "sessionResolution.cookieInvalid"],
    ["session_lookup_failed", "sessionResolution.lookupFailed"],
  ])("shows the message for %s", async (reason, expectedAlert) => {
    resolveMostRecentSession.mockResolvedValue({ reason });

    await renderPage({ loginName: "user@example.com" });

    expect(screen.getByRole("alert").textContent).toBe(expectedAlert);
  });

  test("offers to restart the sign-in with the auth request when the session cannot be resolved", async () => {
    resolveMostRecentSession.mockResolvedValue({ reason: "cookie_missing" });

    await renderPage({ loginName: "user@example.com", organization: "org-1", requestId: "oidc_123" });

    const href = screen.getByTestId("restart-sign-in").getAttribute("href");
    expect(href).toBe("/loginname?loginName=user%40example.com&organization=org-1&requestId=oidc_123");
  });

  test("renders the OTP form and no alert when the session resolves", async () => {
    resolveMostRecentSession.mockResolvedValue({ session: { id: "session-1", factors: { user: { loginName: "user" } } } });

    await renderPage({ loginName: "user@example.com" });

    expect(screen.getByTestId("login-otp")).toBeTruthy();
    expect(screen.queryByRole("alert")).toBeNull();
    expect(screen.queryByTestId("restart-sign-in")).toBeNull();
  });
});
