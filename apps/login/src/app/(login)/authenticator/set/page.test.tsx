import { Translated } from "@/components/translated";
import { isValidElement, ReactElement, ReactNode } from "react";
import { beforeEach, describe, expect, test, vi } from "vitest";
import Page from "./page";

vi.mock("@/lib/server/zitadel", () => ({
  getActiveIdentityProviders: vi.fn(),
  getBrandingSettings: vi.fn(),
  getLoginSettings: vi.fn(),
  getSession: vi.fn(),
  getUserByID: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
}));

vi.mock("@/lib/server/session", () => ({
  loadMostRecentSession: vi.fn(),
}));

vi.mock("@/lib/server/cookies", () => ({
  getSessionCookieById: vi.fn(),
}));

vi.mock("@/lib/server/verify-helper", () => ({
  checkUserVerification: vi.fn(),
}));

vi.mock("@/lib/server/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://example.com" } })),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(() => new Headers()),
}));

// collects all elements of the given type from the (unrendered) element tree returned by the page
function findElements(node: ReactNode, type: unknown): ReactElement<any>[] {
  if (Array.isArray(node)) {
    return node.flatMap((child) => findElements(child, type));
  }
  if (!isValidElement(node)) {
    return [];
  }
  const element = node as ReactElement<any>;
  return [...(element.type === type ? [element] : []), ...findElements(element.props?.children, type)];
}

describe("authenticator set page", () => {
  let zitadel: any;
  let sessionLib: any;
  let cookies: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    zitadel = await import("@/lib/server/zitadel");
    sessionLib = await import("@/lib/server/session");
    cookies = await import("@/lib/server/cookies");
  });

  test("shows the session expired alert when the most recent session is gone", async () => {
    vi.mocked(sessionLib.loadMostRecentSession).mockResolvedValue(undefined);

    const result = await Page({ searchParams: Promise.resolve({ loginName: "user@example.com" }) });

    expect(findElements(result, Translated).map((t) => t.props.i18nKey)).toEqual(["sessionExpired"]);
    expect(zitadel.listAuthenticationMethodTypes).not.toHaveBeenCalled();
  });

  test("shows the session expired alert when the session of the cookie is gone", async () => {
    vi.mocked(cookies.getSessionCookieById).mockResolvedValue({ id: "session-1", token: "token-1" });
    vi.mocked(zitadel.getSession).mockResolvedValue({ session: undefined });

    const result = await Page({ searchParams: Promise.resolve({ sessionId: "session-1" }) });

    expect(findElements(result, Translated).map((t) => t.props.i18nKey)).toEqual(["sessionExpired"]);
    expect(zitadel.listAuthenticationMethodTypes).not.toHaveBeenCalled();
  });
});
