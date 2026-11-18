import { RegisterPasskey } from "@/components/register-passkey";
import { Translated } from "@/components/translated";
import { UserAvatar } from "@/components/user-avatar";
import { isValidElement, ReactElement, ReactNode } from "react";
import { beforeEach, describe, expect, test, vi } from "vitest";
import Page from "./page";

vi.mock("@/lib/server/zitadel", () => ({
  getBrandingSettings: vi.fn(),
  getUserByID: vi.fn(),
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

describe("passkey set page", () => {
  let zitadel: any;
  let sessionLib: any;
  let guard: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    zitadel = await import("@/lib/server/zitadel");
    sessionLib = await import("@/lib/server/session");
    guard = await import("@/lib/server/enrollment-guard");

    vi.mocked(zitadel.getBrandingSettings).mockResolvedValue({} as any);
    vi.mocked(sessionLib.loadMostRecentSession).mockResolvedValue(session as any);
    vi.mocked(guard.getSetupAuthorizationError).mockResolvedValue(null);
  });

  test("offers the registration when the session may set up a passkey", async () => {
    const result = await Page({ searchParams: Promise.resolve({ loginName: "user@example.com" }) });

    expect(guard.getSetupAuthorizationError).toHaveBeenCalledWith(
      expect.objectContaining({ session, userId: "user-1", method: "passkey" }),
    );
    expect(findElements(result, RegisterPasskey)).toHaveLength(1);
  });

  test.each([
    ["You have to verify your second factor first", "secondFactorRequired"],
    ["Passkeys are not allowed", "authenticatorNotAllowed"],
    ["You have to authenticate or have a valid User Verification Check", "sessionExpired"],
  ])("hides the registration and shows a translated error for %s", async (error, i18nKey) => {
    vi.mocked(guard.getSetupAuthorizationError).mockResolvedValue(error);

    const result = await Page({ searchParams: Promise.resolve({ loginName: "user@example.com" }) });

    expect(findElements(result, RegisterPasskey)).toHaveLength(0);
    const messages = findElements(result, Translated).map((element) => element.props);
    expect(messages).toContainEqual(expect.objectContaining({ i18nKey, namespace: "error" }));
  });

  test("leaves the emailed code flow to the code even if a session is in the browser", async () => {
    vi.mocked(guard.getSetupAuthorizationError).mockResolvedValue("You have to verify your second factor first");
    vi.mocked(zitadel.getUserByID).mockResolvedValue({
      user: {
        preferredLoginName: "code-user@example.com",
        type: { case: "human", value: { profile: { displayName: "User" } } },
      },
    } as any);

    const result = await Page({
      searchParams: Promise.resolve({
        loginName: "user@example.com",
        userId: "user-1",
        code: "code",
        codeId: "code-id",
        organization: "org-1",
      }),
    });

    expect(guard.getSetupAuthorizationError).not.toHaveBeenCalled();
    const forms = findElements(result, RegisterPasskey);
    expect(forms).toHaveLength(1);
    // the page shows and continues with the code's user, not the user of the browser session
    expect(forms[0].props.loginName).toBe("code-user@example.com");
    expect(findElements(result, UserAvatar).map((a) => a.props.loginName)).toEqual(["code-user@example.com"]);
  });

  test("leaves the emailed code flow without a session to the code", async () => {
    vi.mocked(zitadel.getUserByID).mockResolvedValue({
      user: {
        preferredLoginName: "user@example.com",
        type: { case: "human", value: { profile: { displayName: "User" } } },
      },
    } as any);

    const result = await Page({
      searchParams: Promise.resolve({ userId: "user-1", code: "code", codeId: "code-id", organization: "org-1" }),
    });

    expect(guard.getSetupAuthorizationError).not.toHaveBeenCalled();
    const forms = findElements(result, RegisterPasskey);
    expect(forms).toHaveLength(1);
    expect(forms[0].props).toMatchObject({ userId: "user-1", code: "code", codeId: "code-id" });
  });
});
