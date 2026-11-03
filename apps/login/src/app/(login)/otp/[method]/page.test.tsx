import { cleanup, render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import Page from "./page";

vi.mock("@/lib/zitadel", () => ({
  getBrandingSettings: vi.fn(),
  getLoginSettings: vi.fn(),
}));

vi.mock("@/lib/session", () => ({
  loadMostRecentSession: vi.fn(),
  loadSessionById: vi.fn(),
}));

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
  let loadSessionById: any;
  let loadMostRecentSession: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const zitadel = await import("@/lib/zitadel");
    const sessionLib = await import("@/lib/session");

    loadSessionById = vi.mocked(sessionLib.loadSessionById);
    loadMostRecentSession = vi.mocked(sessionLib.loadMostRecentSession);
    vi.mocked(zitadel.getBrandingSettings).mockResolvedValue({} as any);
    vi.mocked(zitadel.getLoginSettings).mockResolvedValue({} as any);
  });

  afterEach(cleanup);

  // A rejected session cookie must not crash the render: the shared loaders resolve undefined and
  // the page falls back to its no-session alert on every branch.
  test.each([
    {
      name: "sessionId and loginName",
      searchParams: { sessionId: "session-1", loginName: "user@example.com", organization: "org-1" },
      expectedAlert: "sessionExpired",
    },
    {
      name: "loginName only",
      searchParams: { loginName: "user@example.com", organization: "org-1" },
      expectedAlert: "sessionExpired",
    },
    {
      name: "neither sessionId nor loginName",
      searchParams: {},
      expectedAlert: "unknownContext",
    },
  ])(
    "renders the $expectedAlert alert and no OTP form when the session cannot be loaded ($name)",
    async ({ searchParams, expectedAlert }) => {
      loadSessionById.mockResolvedValue(undefined);
      loadMostRecentSession.mockResolvedValue(undefined);

      const result = await Page({
        searchParams: Promise.resolve(searchParams),
        params: Promise.resolve({ method: "time-based" }),
      });
      render(result);

      if (searchParams.sessionId) {
        expect(loadSessionById).toHaveBeenCalledWith({ baseUrl: "https://example.com" }, "session-1", "org-1");
        expect(loadMostRecentSession).not.toHaveBeenCalled();
      } else {
        expect(loadMostRecentSession).toHaveBeenCalledWith({
          serviceConfig: { baseUrl: "https://example.com" },
          sessionParams: { loginName: searchParams.loginName, organization: searchParams.organization },
        });
        expect(loadSessionById).not.toHaveBeenCalled();
      }
      expect(screen.getByRole("alert").textContent).toBe(expectedAlert);
      expect(screen.queryByTestId("login-otp")).toBeNull();
    },
  );
});
