import { create, timestampMs } from "@zitadel/client";
import { RequestChallengesSchema } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { beforeEach, describe, expect, test, vi } from "vitest";

vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("@/lib/zitadel", () => ({
  createSessionForUserIdAndIdpIntent: vi.fn(),
  createSessionFromChecksAndChallenges: vi.fn(),
  getSecuritySettings: vi.fn(),
  getSession: vi.fn(),
  setSession: vi.fn(),
}));

vi.mock("@/lib/cookies", () => ({
  addSessionToCookie: vi.fn(),
  updateSessionCookie: vi.fn(),
}));

import { addSessionToCookie, updateSessionCookie } from "@/lib/cookies";
import { createSessionFromChecksAndChallenges, getSecuritySettings, getSession, setSession } from "@/lib/zitadel";
import { headers } from "next/headers";
import { getServiceConfig } from "../service-url";
import { createSessionAndUpdateCookie, setSessionAndUpdateCookie } from "./cookie";

// A challenge request as an attacker-controlled browser could construct it: both OTP methods
// asking for the code to come back in the response instead of being sent to the user.
function returnCodeChallenges() {
  return create(RequestChallengesSchema, {
    otpSms: { returnCode: true },
    otpEmail: { deliveryType: { case: "returnCode", value: {} } },
  });
}

function expectSanitized(challenges: any) {
  expect(challenges.otpSms.returnCode).toBe(false);
  expect(challenges.otpEmail.deliveryType.case).toBe("sendCode");
}

describe("session challenge sanitization", () => {
  beforeEach(() => {
    vi.clearAllMocks();

    vi.mocked(headers).mockResolvedValue({} as any);
    vi.mocked(getServiceConfig).mockReturnValue({ serviceConfig: {} } as any);
    vi.mocked(getSecuritySettings).mockResolvedValue({ embeddedIframe: { enabled: false } } as any);
    vi.mocked(getSession).mockResolvedValue({
      session: {
        id: "sessionId",
        factors: { user: { loginName: "victim@example.com", organizationId: "orgId" } },
      },
    } as any);
    vi.mocked(addSessionToCookie).mockResolvedValue(undefined as any);
    vi.mocked(updateSessionCookie).mockResolvedValue(undefined as any);
  });

  test("createSessionAndUpdateCookie strips returnCode before calling the API", async () => {
    vi.mocked(createSessionFromChecksAndChallenges).mockResolvedValue({
      sessionId: "sessionId",
      sessionToken: "sessionToken",
    } as any);

    await createSessionAndUpdateCookie({
      checks: {} as any,
      requestId: undefined,
      challenges: returnCodeChallenges(),
    });

    expect(createSessionFromChecksAndChallenges).toHaveBeenCalledTimes(1);
    expectSanitized(vi.mocked(createSessionFromChecksAndChallenges).mock.calls[0][0].challenges);
  });

  test("setSessionAndUpdateCookie strips returnCode before calling the API", async () => {
    vi.mocked(setSession).mockResolvedValue({
      sessionToken: "sessionToken",
      details: {},
    } as any);

    await setSessionAndUpdateCookie({
      recentCookie: {
        id: "sessionId",
        token: "sessionToken",
        loginName: "victim@example.com",
        creationTs: "",
        expirationTs: "",
        changeTs: "",
      },
      challenges: returnCodeChallenges(),
      lifetime: {} as any,
    });

    expect(setSession).toHaveBeenCalledTimes(1);
    expectSanitized(vi.mocked(setSession).mock.calls[0][0].challenges);
  });
});

describe("session cookie renewal", () => {
  beforeEach(() => {
    vi.clearAllMocks();

    vi.mocked(headers).mockResolvedValue({} as any);
    vi.mocked(getServiceConfig).mockReturnValue({ serviceConfig: {} } as any);
    vi.mocked(getSecuritySettings).mockResolvedValue({ embeddedIframe: { enabled: false } } as any);
    vi.mocked(addSessionToCookie).mockResolvedValue(undefined as any);
    vi.mocked(updateSessionCookie).mockResolvedValue(undefined as any);
    vi.mocked(setSession).mockResolvedValue({
      sessionToken: "sessionToken",
      details: {},
    } as any);
  });

  test("refreshes expirationTs from the newly fetched session instead of the stale cookie", async () => {
    const freshExpirationDate = { seconds: BigInt(2000), nanos: 0 };

    vi.mocked(getSession).mockResolvedValue({
      session: {
        id: "sessionId",
        expirationDate: freshExpirationDate,
        factors: { user: { loginName: "victim@example.com", organizationId: "orgId" } },
      },
    } as any);

    await setSessionAndUpdateCookie({
      recentCookie: {
        id: "sessionId",
        token: "sessionToken",
        loginName: "victim@example.com",
        creationTs: "",
        expirationTs: "1000000", // stale value carried over from before this renewal
        changeTs: "",
      },
      lifetime: {} as any,
    });

    const newCookie = vi.mocked(updateSessionCookie).mock.calls[0][0].session;
    expect(newCookie.expirationTs).toBe(`${timestampMs(freshExpirationDate as any)}`);
    expect(newCookie.expirationTs).not.toBe("1000000");
  });

  test("falls back to the previous expirationTs when the fetched session has none", async () => {
    vi.mocked(getSession).mockResolvedValue({
      session: {
        id: "sessionId",
        factors: { user: { loginName: "victim@example.com", organizationId: "orgId" } },
      },
    } as any);

    await setSessionAndUpdateCookie({
      recentCookie: {
        id: "sessionId",
        token: "sessionToken",
        loginName: "victim@example.com",
        creationTs: "",
        expirationTs: "1000000",
        changeTs: "",
      },
      lifetime: {} as any,
    });

    const newCookie = vi.mocked(updateSessionCookie).mock.calls[0][0].session;
    expect(newCookie.expirationTs).toBe("1000000");
  });
});
