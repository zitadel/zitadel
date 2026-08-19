import { setSessionAndUpdateCookie } from "@/lib/server/cookie";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { updateOrCreateSession } from "./session";

vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("@zitadel/client", () => ({
  create: vi.fn((_schema: any, value: any) => value),
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("./host", () => ({
  getPublicHost: vi.fn(),
}));

vi.mock("@/lib/zitadel", () => ({
  deleteSession: vi.fn(),
  getLoginSettings: vi.fn(),
  getSecuritySettings: vi.fn(),
  humanMFAInitSkipped: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
  listUsers: vi.fn(),
}));

vi.mock("@/lib/server/cookie", () => ({
  createSessionAndUpdateCookie: vi.fn(),
  setSessionAndUpdateCookie: vi.fn(),
}));

vi.mock("../cookies", () => ({
  getMostRecentSessionCookie: vi.fn(),
  getSessionCookieById: vi.fn(),
  getSessionCookieByLoginName: vi.fn(),
  removeSessionFromCookie: vi.fn(),
}));

vi.mock("../client", () => ({
  completeFlowOrGetUrl: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

vi.mock("@/lib/logger", () => ({
  createLogger: () => ({ debug: vi.fn(), info: vi.fn(), warn: vi.fn(), error: vi.fn() }),
}));

vi.mock("@/lib/grpc/interceptors/error-classification", () => ({
  isClassifiedError: vi.fn(() => false),
}));

import { getLoginSettings } from "@/lib/zitadel";
import { getSessionCookieById } from "../cookies";
import { getServiceConfig } from "../service-url";
import { getPublicHost } from "./host";

const THIRTY_DAYS = { seconds: BigInt(720 * 60 * 60), nanos: 0 };
const TWELVE_HOURS = { seconds: BigInt(12 * 60 * 60), nanos: 0 };
const TWENTY_FOUR_HOURS = BigInt(24 * 60 * 60);

const SESSION_COOKIE = {
  id: "session-1",
  token: "token-1",
  loginName: "user@example.com",
  organization: "org-1",
  creationTs: "1700000000000",
  expirationTs: "1800000000000",
  changeTs: "1700000000000",
};

describe("updateOrCreateSession session lifetime", () => {
  beforeEach(() => {
    vi.clearAllMocks();

    vi.mocked(getPublicHost).mockReturnValue("localhost:3000");
    vi.mocked(getServiceConfig).mockReturnValue({ serviceConfig: {} } as any);
    vi.mocked(getSessionCookieById).mockResolvedValue(SESSION_COOKIE as any);
    vi.mocked(setSessionAndUpdateCookie).mockResolvedValue({
      id: "session-1",
      factors: { user: { id: "user-1", loginName: "user@example.com" } },
    } as any);
    vi.mocked(getLoginSettings).mockResolvedValue({
      passwordCheckLifetime: THIRTY_DAYS,
      secondFactorCheckLifetime: THIRTY_DAYS,
      multiFactorCheckLifetime: THIRTY_DAYS,
      externalLoginCheckLifetime: THIRTY_DAYS,
    } as any);
  });

  test("uses secondFactorCheckLifetime for a TOTP check", async () => {
    await updateOrCreateSession({
      sessionId: "session-1",
      checks: { totp: { code: "123456" } } as any,
    });

    expect(setSessionAndUpdateCookie).toHaveBeenCalledWith(expect.objectContaining({ lifetime: THIRTY_DAYS }));
  });

  test("uses secondFactorCheckLifetime for an email OTP check", async () => {
    await updateOrCreateSession({
      sessionId: "session-1",
      checks: { otpEmail: { code: "123456" } } as any,
    });

    expect(setSessionAndUpdateCookie).toHaveBeenCalledWith(expect.objectContaining({ lifetime: THIRTY_DAYS }));
  });

  test("uses secondFactorCheckLifetime for an SMS OTP check", async () => {
    await updateOrCreateSession({
      sessionId: "session-1",
      checks: { otpSms: { code: "123456" } } as any,
    });

    expect(setSessionAndUpdateCookie).toHaveBeenCalledWith(expect.objectContaining({ lifetime: THIRTY_DAYS }));
  });

  test("uses multiFactorCheckLifetime for a WebAuthN check", async () => {
    await updateOrCreateSession({
      sessionId: "session-1",
      checks: { webAuthN: { credentialAssertionData: {} } } as any,
    });

    expect(setSessionAndUpdateCookie).toHaveBeenCalledWith(expect.objectContaining({ lifetime: THIRTY_DAYS }));
  });

  test("prefers an explicitly passed lifetime over the login settings", async () => {
    await updateOrCreateSession({
      sessionId: "session-1",
      checks: { totp: { code: "123456" } } as any,
      lifetime: TWELVE_HOURS,
    });

    expect(setSessionAndUpdateCookie).toHaveBeenCalledWith(expect.objectContaining({ lifetime: TWELVE_HOURS }));
  });

  test("falls back to 24 hours when the login settings carry no lifetime", async () => {
    vi.mocked(getLoginSettings).mockResolvedValue({} as any);

    await updateOrCreateSession({
      sessionId: "session-1",
      checks: { totp: { code: "123456" } } as any,
    });

    expect(setSessionAndUpdateCookie).toHaveBeenCalledWith(
      expect.objectContaining({ lifetime: expect.objectContaining({ seconds: TWENTY_FOUR_HOURS }) }),
    );
  });
});
