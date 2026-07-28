import { Code, ConnectError } from "@connectrpc/connect";
import { create, toJson } from "@zitadel/client";
import { Checks, ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { ClassifiedConnectError } from "../grpc/interceptors/error-classification";
import { updateOrCreateSession } from "./session";

vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("@/lib/server/cookie", () => ({
  createSessionAndUpdateCookie: vi.fn(),
  setSessionAndUpdateCookie: vi.fn(),
}));

vi.mock("@/lib/zitadel", () => ({
  deleteSession: vi.fn(),
  getLoginSettings: vi.fn(),
  getSecuritySettings: vi.fn(),
  humanMFAInitSkipped: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
  listUsers: vi.fn(),
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

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("./host", () => ({
  getPublicHost: vi.fn(),
}));

// Mock translations - returns the key itself for testing
vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

const USER_ID = "1416cf78-43e0-423d-9691-49fd4f4df776";
const LOGIN_NAME = "user@example.com";

describe("updateOrCreateSession", () => {
  let mockCreateSessionAndUpdateCookie: any;
  let mockSetSessionAndUpdateCookie: any;
  let mockListUsers: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { getPublicHost } = await import("./host");
    const { getSessionCookieByLoginName } = await import("../cookies");
    const { getLoginSettings, listUsers } = await import("@/lib/zitadel");
    const { createSessionAndUpdateCookie, setSessionAndUpdateCookie } = await import("@/lib/server/cookie");

    vi.mocked(headers).mockResolvedValue(new Headers() as any);
    vi.mocked(getServiceConfig).mockReturnValue({ serviceConfig: {} } as any);
    vi.mocked(getPublicHost).mockReturnValue("localhost:3000");
    vi.mocked(getSessionCookieByLoginName).mockResolvedValue({
      id: "session-id",
      token: "session-token",
      loginName: LOGIN_NAME,
      organization: "org-id",
    } as any);
    vi.mocked(getLoginSettings).mockResolvedValue({} as any);
    vi.mocked(listUsers).mockResolvedValue({
      details: { totalResult: BigInt(1) },
      result: [{ userId: USER_ID }],
    } as any);
    vi.mocked(createSessionAndUpdateCookie).mockResolvedValue({
      session: { id: "new-session-id", factors: { user: { id: USER_ID, loginName: LOGIN_NAME } } },
    } as any);

    mockCreateSessionAndUpdateCookie = vi.mocked(createSessionAndUpdateCookie);
    mockSetSessionAndUpdateCookie = vi.mocked(setSessionAndUpdateCookie);
    mockListUsers = vi.mocked(listUsers);
  });

  test("returns an invalid code error instead of re-creating the session when the code check is rejected", async () => {
    // ZITADEL answers a wrong TOTP/OTP code with invalid_argument (Errors.User.MFA.OTP.InvalidCode).
    mockSetSessionAndUpdateCookie.mockRejectedValue(
      new ClassifiedConnectError(new ConnectError("Code is invalid (EVENT-8isk2)", Code.InvalidArgument)),
    );

    const result = await updateOrCreateSession({
      loginName: LOGIN_NAME,
      organization: "org-id",
      checks: create(ChecksSchema, { totp: { code: "000000" } }),
    });

    expect(result).toEqual({ error: "invalidCode" });
    // The session is fine - re-creating it would only replay the same rejected code.
    expect(mockListUsers).not.toHaveBeenCalled();
    expect(mockCreateSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("re-creates the session with a serializable Checks message when the session is gone", async () => {
    // A terminated/expired session answers with failed_precondition, which is what the
    // fallback below is meant to recover from.
    mockSetSessionAndUpdateCookie.mockRejectedValue(
      new ClassifiedConnectError(new ConnectError("Errors.Session.Terminated", Code.FailedPrecondition)),
    );

    const result = await updateOrCreateSession({
      loginName: LOGIN_NAME,
      organization: "org-id",
      checks: create(ChecksSchema, { totp: { code: "123456" } }),
    });

    expect(mockCreateSessionAndUpdateCookie).toHaveBeenCalledOnce();
    const newChecks: Checks = mockCreateSessionAndUpdateCookie.mock.calls[0][0].checks;

    // Regression guard: a plain object literal for `user` survives `create()` unconverted
    // (protobuf-es returns an init value that already is a message of the target type as-is),
    // which made serialization fail with "cannot use field
    // zitadel.session.v2.CheckUser.user_id with message undefined" and turned a wrong code
    // into an HTTP 500.
    expect(() => toJson(ChecksSchema, newChecks)).not.toThrow();
    expect(toJson(ChecksSchema, newChecks)).toEqual({
      user: { userId: USER_ID },
      totp: { code: "123456" },
    });

    expect(result).toEqual(expect.objectContaining({ sessionId: "new-session-id" }));
  });
});
