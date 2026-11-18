import { RegisterU2f } from "@/components/register-u2f";
import { Translated } from "@/components/translated";
import { isValidElement, ReactElement, ReactNode } from "react";
import { beforeEach, describe, expect, test, vi } from "vitest";
import Page from "./page";

vi.mock("@/lib/server/zitadel", () => ({
  getBrandingSettings: vi.fn(),
}));

vi.mock("@/lib/server/enrollment-guard", async (importOriginal) => ({
  ...(await importOriginal<typeof import("@/lib/server/enrollment-guard")>()),
  getSetupAuthorizationError: vi.fn(),
}));

vi.mock("@/lib/server/session", () => ({
  loadMostRecentSession: vi.fn(),
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

const session = {
  id: "session-1",
  factors: { user: { id: "user-1", loginName: "user@example.com", displayName: "User", organizationId: "org-1" } },
};

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

describe("u2f set page", () => {
  let sessionLib: any;
  let guard: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const zitadel = await import("@/lib/server/zitadel");
    sessionLib = await import("@/lib/server/session");
    guard = await import("@/lib/server/enrollment-guard");

    vi.mocked(zitadel.getBrandingSettings).mockResolvedValue({} as any);
    vi.mocked(sessionLib.loadMostRecentSession).mockResolvedValue(session as any);
    vi.mocked(guard.getSetupAuthorizationError).mockResolvedValue(null);
  });

  test("offers the registration when the session may set up U2F", async () => {
    const result = await Page({ searchParams: Promise.resolve({ loginName: "user@example.com" }) });

    expect(guard.getSetupAuthorizationError).toHaveBeenCalledWith(
      expect.objectContaining({ session, userId: "user-1", method: "u2f" }),
    );
    expect(findElements(result, RegisterU2f)).toHaveLength(1);
  });

  test.each([
    ["You have to verify your second factor first", "secondFactorRequired"],
    ["This second factor is not allowed", "authenticatorNotAllowed"],
    ["You have to authenticate or have a valid User Verification Check", "sessionExpired"],
  ])("hides the registration and shows a translated error for %s", async (error, i18nKey) => {
    vi.mocked(guard.getSetupAuthorizationError).mockResolvedValue(error);

    const result = await Page({ searchParams: Promise.resolve({ loginName: "user@example.com" }) });

    expect(findElements(result, RegisterU2f)).toHaveLength(0);
    const messages = findElements(result, Translated).map((element) => element.props);
    expect(messages).toContainEqual(expect.objectContaining({ i18nKey, namespace: "error" }));
  });

  test("keeps the session expired message when the session can't be resolved", async () => {
    vi.mocked(sessionLib.loadMostRecentSession).mockResolvedValue(undefined);

    const result = await Page({ searchParams: Promise.resolve({ loginName: "user@example.com" }) });

    expect(guard.getSetupAuthorizationError).not.toHaveBeenCalled();
    const messages = findElements(result, Translated).map((element) => element.props);
    expect(messages).toContainEqual(expect.objectContaining({ i18nKey: "sessionExpired", namespace: "error" }));
  });
});
