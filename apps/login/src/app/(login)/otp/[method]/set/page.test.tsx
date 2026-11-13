import { OTPAddForm } from "@/components/otp-add-form";
import { Translated } from "@/components/translated";
import { isValidElement, ReactElement, ReactNode } from "react";
import { beforeEach, describe, expect, test, vi } from "vitest";
import Page from "./page";

vi.mock("@/lib/zitadel", () => ({
  getBrandingSettings: vi.fn(),
  getLoginSettings: vi.fn(),
  registerTOTP: vi.fn(),
  addOTPEmail: vi.fn(),
  addOTPSMS: vi.fn(),
}));

vi.mock("@/lib/server/enrollment-guard", () => ({
  getEnrollmentAuthorizationError: vi.fn(),
}));

vi.mock("@/lib/server/otp", () => ({
  addOTP: vi.fn(),
}));

vi.mock("@/lib/session", () => ({
  loadMostRecentSession: vi.fn(),
  loadSessionById: vi.fn(),
}));

vi.mock("@/lib/service-url", () => ({
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

describe("otp set page", () => {
  let zitadel: any;
  let sessionLib: any;
  let guard: any;
  let addOTP: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    zitadel = await import("@/lib/zitadel");
    sessionLib = await import("@/lib/session");
    guard = await import("@/lib/server/enrollment-guard");
    addOTP = vi.mocked((await import("@/lib/server/otp")).addOTP);

    vi.mocked(zitadel.getBrandingSettings).mockResolvedValue({} as any);
    vi.mocked(zitadel.getLoginSettings).mockResolvedValue({} as any);
    vi.mocked(sessionLib.loadMostRecentSession).mockResolvedValue(session as any);
    vi.mocked(sessionLib.loadSessionById).mockResolvedValue(session as any);
    vi.mocked(guard.getEnrollmentAuthorizationError).mockResolvedValue(null);
  });

  // Rendering the page (e.g. by following a link) must never add a second factor to the
  // signed-in user; OTP via email or SMS is only added on explicit confirmation.
  test.each(["email", "sms"])("does not add OTP via %s while rendering but asks for confirmation", async (method) => {
    const result = await Page({
      params: Promise.resolve({ method }),
      searchParams: Promise.resolve({}),
    });

    expect(zitadel.addOTPEmail).not.toHaveBeenCalled();
    expect(zitadel.addOTPSMS).not.toHaveBeenCalled();
    expect(addOTP).not.toHaveBeenCalled();

    const forms = findElements(result, OTPAddForm);
    expect(forms).toHaveLength(1);
    expect(forms[0].props).toMatchObject({ method, sessionId: "session-1", urlToContinue: "/accounts", checkAfter: false });
  });

  test.each(["email", "sms"])("does not redirect for %s with checkAfter while rendering", async (method) => {
    const result = await Page({
      params: Promise.resolve({ method }),
      searchParams: Promise.resolve({ checkAfter: "true", loginName: "user@example.com" }),
    });

    expect(addOTP).not.toHaveBeenCalled();
    const forms = findElements(result, OTPAddForm);
    expect(forms).toHaveLength(1);
    expect(forms[0].props.checkAfter).toBe(true);
    expect(forms[0].props.urlToContinue).toBe(`/otp/${method}?loginName=user%40example.com`);
  });

  test("hides the confirmation and shows the error when enrollment is not authorized", async () => {
    vi.mocked(guard.getEnrollmentAuthorizationError).mockResolvedValue("You have to authenticate");

    const result = await Page({
      params: Promise.resolve({ method: "email" }),
      searchParams: Promise.resolve({}),
    });

    expect(findElements(result, OTPAddForm)).toHaveLength(0);
    expect(JSON.stringify(result)).toContain("You have to authenticate");
  });

  test("renders an alert instead of the confirmation without a session", async () => {
    vi.mocked(sessionLib.loadMostRecentSession).mockResolvedValue(undefined);

    const result = await Page({
      params: Promise.resolve({ method: "email" }),
      searchParams: Promise.resolve({}),
    });

    expect(findElements(result, OTPAddForm)).toHaveLength(0);
    expect(findElements(result, Translated).map((t) => t.props.i18nKey)).toContain("unknownContext");
  });

  test("acts on the session given by sessionId instead of the most recent one", async () => {
    vi.mocked(sessionLib.loadSessionById).mockResolvedValue({ ...session, id: "session-a" } as any);

    const result = await Page({
      params: Promise.resolve({ method: "email" }),
      searchParams: Promise.resolve({ sessionId: "session-a", organization: "org-1" }),
    });

    expect(sessionLib.loadSessionById).toHaveBeenCalledWith({ baseUrl: "https://example.com" }, "session-a", "org-1");
    expect(sessionLib.loadMostRecentSession).not.toHaveBeenCalled();
    expect(findElements(result, OTPAddForm)[0].props.sessionId).toBe("session-a");
  });

  test("renders the session expired alert when the session given by sessionId can't be loaded", async () => {
    vi.mocked(sessionLib.loadSessionById).mockResolvedValue(undefined);

    const result = await Page({
      params: Promise.resolve({ method: "email" }),
      searchParams: Promise.resolve({ sessionId: "session-a" }),
    });

    expect(sessionLib.loadMostRecentSession).not.toHaveBeenCalled();
    expect(findElements(result, OTPAddForm)).toHaveLength(0);
    expect(findElements(result, Translated).map((t) => t.props.i18nKey)).toContain("sessionExpired");
  });
});
