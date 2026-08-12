import { Code, ConnectError } from "@connectrpc/connect";
import { create } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { ClassifiedConnectError } from "../grpc/interceptors/error-classification";
import { updateOrCreateSession } from "./session";

vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

// this returns the key itself so tests can assert on translation keys
vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("../zitadel", () => ({
  deleteSession: vi.fn(),
  getLoginSettings: vi.fn(),
  getSecuritySettings: vi.fn(),
  humanMFAInitSkipped: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
  listUsers: vi.fn(),
}));

vi.mock("./cookie", () => ({
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

vi.mock("../session", () => ({
  isSessionValid: vi.fn(),
}));

vi.mock("./host", () => ({
  getPublicHost: vi.fn(),
}));

vi.mock("./loginname", () => ({
  sendLoginname: vi.fn(),
}));

function classifiedError(code: Code, message: string): ClassifiedConnectError {
  return new ClassifiedConnectError(new ConnectError(message, code));
}

describe("updateOrCreateSession", () => {
  let mockHeaders: any;
  let mockGetServiceConfig: any;
  let mockGetLoginSettings: any;
  let mockListUsers: any;
  let mockSetSessionAndUpdateCookie: any;
  let mockCreateSessionAndUpdateCookie: any;
  let mockGetSessionCookieByLoginName: any;
  let mockGetPublicHost: any;

  const sessionCookie = {
    id: "session-1",
    token: "token-1",
    loginName: "user@example.com",
    organization: "org-1",
    creationTs: "",
    expirationTs: "",
    changeTs: "",
  };

  const otpChecks = () => create(ChecksSchema, { otpSms: { code: "123456" } });

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { getLoginSettings, listUsers } = await import("../zitadel");
    const { createSessionAndUpdateCookie, setSessionAndUpdateCookie } = await import("./cookie");
    const { getSessionCookieByLoginName } = await import("../cookies");
    const { getPublicHost } = await import("./host");

    mockHeaders = vi.mocked(headers);
    mockGetServiceConfig = vi.mocked(getServiceConfig);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockListUsers = vi.mocked(listUsers);
    mockSetSessionAndUpdateCookie = vi.mocked(setSessionAndUpdateCookie);
    mockCreateSessionAndUpdateCookie = vi.mocked(createSessionAndUpdateCookie);
    mockGetSessionCookieByLoginName = vi.mocked(getSessionCookieByLoginName);
    mockGetPublicHost = vi.mocked(getPublicHost);

    mockHeaders.mockResolvedValue({} as any);
    mockGetServiceConfig.mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } });
    mockGetPublicHost.mockReturnValue("login.example.com");
    mockGetLoginSettings.mockResolvedValue({ secondFactorCheckLifetime: { seconds: BigInt(3600), nanos: 0 } });
    mockGetSessionCookieByLoginName.mockResolvedValue(sessionCookie);
  });

  test("returns invalidCode error for a wrong OTP code instead of throwing", async () => {
    mockSetSessionAndUpdateCookie.mockRejectedValue(classifiedError(Code.InvalidArgument, "Errors.User.Code.Invalid"));

    const response = await updateOrCreateSession({
      loginName: "user@example.com",
      organization: "org-1",
      checks: otpChecks(),
    });

    expect(response).toEqual({ error: "invalidCode" });
    // a failed check must never trigger the session-recreation fallback
    expect(mockListUsers).not.toHaveBeenCalled();
    expect(mockCreateSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("recreates the session when it is gone server-side (NotFound)", async () => {
    mockSetSessionAndUpdateCookie.mockRejectedValue(classifiedError(Code.NotFound, "Errors.Session.NotExisting"));
    mockListUsers.mockResolvedValue({
      details: { totalResult: BigInt(1) },
      result: [{ userId: "user-1" }],
    });
    mockCreateSessionAndUpdateCookie.mockResolvedValue({
      session: {
        id: "session-2",
        factors: { user: { id: "user-1", loginName: "user@example.com", organizationId: "org-1" } },
      },
      sessionCookie,
      challenges: undefined,
    });

    const response = await updateOrCreateSession({
      loginName: "user@example.com",
      organization: "org-1",
      checks: otpChecks(),
    });

    expect(response).toMatchObject({ sessionId: "session-2" });
    expect(mockCreateSessionAndUpdateCookie).toHaveBeenCalledWith(
      expect.objectContaining({
        checks: expect.objectContaining({
          user: expect.objectContaining({ search: { case: "userId", value: "user-1" } }),
        }),
      }),
    );
  });

  test("returns couldNotFindSession when the session is gone and the user cannot be resolved", async () => {
    mockSetSessionAndUpdateCookie.mockRejectedValue(classifiedError(Code.NotFound, "Errors.Session.NotExisting"));
    mockListUsers.mockResolvedValue({ details: { totalResult: BigInt(0) }, result: [] });

    const response = await updateOrCreateSession({
      loginName: "user@example.com",
      checks: otpChecks(),
    });

    expect(response).toEqual({ error: "couldNotFindSession" });
    expect(mockCreateSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("returns an error when the recreated session rejects the checks", async () => {
    mockSetSessionAndUpdateCookie.mockRejectedValue(classifiedError(Code.NotFound, "Errors.Session.NotExisting"));
    mockListUsers.mockResolvedValue({
      details: { totalResult: BigInt(1) },
      result: [{ userId: "user-1" }],
    });
    mockCreateSessionAndUpdateCookie.mockRejectedValue(classifiedError(Code.InvalidArgument, "Errors.User.Code.Invalid"));

    const response = await updateOrCreateSession({
      loginName: "user@example.com",
      checks: otpChecks(),
    });

    expect(response).toEqual({ error: "invalidCode" });
  });

  test("rethrows genuine server errors so they still surface as 500", async () => {
    mockSetSessionAndUpdateCookie.mockRejectedValue(classifiedError(Code.Internal, "database down"));

    await expect(
      updateOrCreateSession({
        loginName: "user@example.com",
        checks: otpChecks(),
      }),
    ).rejects.toThrow("database down");
  });

  test("does not claim invalidCode for non-InvalidArgument failures of a code check", async () => {
    mockSetSessionAndUpdateCookie.mockRejectedValue(classifiedError(Code.ResourceExhausted, "Errors.Limits.ExceededQuota"));

    const response = await updateOrCreateSession({
      loginName: "user@example.com",
      checks: otpChecks(),
    });

    expect(response).toEqual({ error: "couldNotUpdateSession" });
  });

  test("passes through preformatted password-attempt errors", async () => {
    mockSetSessionAndUpdateCookie.mockRejectedValue({
      error: "Failed to authenticate: You had 2 password attempts.",
      failedAttempts: 2,
    });

    const response = await updateOrCreateSession({
      loginName: "user@example.com",
      checks: create(ChecksSchema, { password: { password: "wrong" } }),
    });

    expect(response).toEqual({ error: "Failed to authenticate: You had 2 password attempts." });
    expect(mockListUsers).not.toHaveBeenCalled();
  });

  test("returns couldNotUpdateSession for non-code check failures", async () => {
    mockSetSessionAndUpdateCookie.mockRejectedValue(classifiedError(Code.FailedPrecondition, "Errors.User.NotActive"));

    const response = await updateOrCreateSession({
      loginName: "user@example.com",
      checks: create(ChecksSchema, { webAuthN: { credentialAssertionData: {} } }),
    });

    expect(response).toEqual({ error: "couldNotUpdateSession" });
  });
});
