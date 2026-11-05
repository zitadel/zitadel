import { Code, ConnectError } from "@connectrpc/connect";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { ClassifiedConnectError } from "../grpc/interceptors/error-classification";
import { sendRecoveryCode } from "./recovery-code";

// Mock dependencies
vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("@zitadel/client", () => ({
  create: (_schema: any, init: any) => init,
  Code: { InvalidArgument: 3, FailedPrecondition: 9 },
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("../zitadel", () => ({
  getLoginSettings: vi.fn(),
}));

vi.mock("./cookie", () => ({
  setSessionAndUpdateCookie: vi.fn(),
}));

vi.mock("../cookies", () => ({
  getSessionCookieById: vi.fn(),
  getSessionCookieByLoginName: vi.fn(),
}));

vi.mock("../client", () => ({
  completeFlowOrGetUrl: vi.fn(),
}));

vi.mock("../metrics", () => ({
  recordAuthAttempt: vi.fn(),
  recordAuthSuccess: vi.fn(),
  recordAuthFailure: vi.fn(),
}));

vi.mock("../logger", () => ({
  createLogger: () => ({
    debug: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  }),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

/**
 * Builds the kind of error the classified transport surfaces for a failed recovery code check.
 * The backend translates the message before building the status, so the text is
 * `"<translated text> (<error id>)"`.
 */
function classifiedError(message: string, code: Code) {
  return new ClassifiedConnectError(new ConnectError(message, code));
}

describe("sendRecoveryCode", () => {
  let mockHeaders: any;
  let mockGetServiceConfig: any;
  let mockGetSessionCookieById: any;
  let mockGetSessionCookieByLoginName: any;
  let mockGetLoginSettings: any;
  let mockSetSessionAndUpdateCookie: any;
  let mockCompleteFlowOrGetUrl: any;
  let mockRecordAuthAttempt: any;
  let mockRecordAuthSuccess: any;
  let mockRecordAuthFailure: any;

  const sessionCookie = {
    id: "session123",
    token: "token123",
    loginName: "test@example.com",
    organization: "org123",
  };

  const updatedSession = {
    id: "session123",
    factors: {
      user: {
        id: "user123",
        loginName: "test@example.com",
        organizationId: "org123",
      },
    },
  };

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { getSessionCookieById, getSessionCookieByLoginName } = await import("../cookies");
    const { getLoginSettings } = await import("../zitadel");
    const { setSessionAndUpdateCookie } = await import("./cookie");
    const { completeFlowOrGetUrl } = await import("../client");
    const { recordAuthAttempt, recordAuthSuccess, recordAuthFailure } = await import("../metrics");

    mockHeaders = vi.mocked(headers);
    mockGetServiceConfig = vi.mocked(getServiceConfig);
    mockGetSessionCookieById = vi.mocked(getSessionCookieById);
    mockGetSessionCookieByLoginName = vi.mocked(getSessionCookieByLoginName);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockSetSessionAndUpdateCookie = vi.mocked(setSessionAndUpdateCookie);
    mockCompleteFlowOrGetUrl = vi.mocked(completeFlowOrGetUrl);
    mockRecordAuthAttempt = vi.mocked(recordAuthAttempt);
    mockRecordAuthSuccess = vi.mocked(recordAuthSuccess);
    mockRecordAuthFailure = vi.mocked(recordAuthFailure);

    mockHeaders.mockResolvedValue({});
    mockGetServiceConfig.mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } });

    mockGetSessionCookieById.mockResolvedValue(sessionCookie);
    mockGetSessionCookieByLoginName.mockResolvedValue(sessionCookie);

    mockGetLoginSettings.mockResolvedValue({
      secondFactorCheckLifetime: { seconds: BigInt(600), nanos: 0 },
    });

    mockSetSessionAndUpdateCookie.mockResolvedValue(updatedSession);
    mockCompleteFlowOrGetUrl.mockResolvedValue({ redirect: "/signedin" });
  });

  describe("session lookup", () => {
    test("returns couldNotFindSession without checking the code when no session cookie exists", async () => {
      mockGetSessionCookieByLoginName.mockResolvedValue(undefined);

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(result).toEqual({ error: "verify.errors.couldNotFindSession" });
      // a failed check counts towards the lockout policy, so it must never be retried on a new session
      expect(mockSetSessionAndUpdateCookie).not.toHaveBeenCalled();
      expect(mockRecordAuthAttempt).toHaveBeenCalledWith("recovery_code", undefined);
      expect(mockRecordAuthFailure).toHaveBeenCalledWith("recovery_code", "session_not_found", undefined);
    });

    test("looks the session up by id when a sessionId is given", async () => {
      await sendRecoveryCode({ sessionId: "session123", organization: "org123", code: "abcde-fghij" });

      expect(mockGetSessionCookieById).toHaveBeenCalledWith({ sessionId: "session123", organization: "org123" });
      expect(mockGetSessionCookieByLoginName).not.toHaveBeenCalled();
    });

    test("looks the session up by login name when no sessionId is given", async () => {
      await sendRecoveryCode({ loginName: "test@example.com", organization: "org123", code: "abcde-fghij" });

      expect(mockGetSessionCookieByLoginName).toHaveBeenCalledWith({
        loginName: "test@example.com",
        organization: "org123",
      });
      expect(mockGetSessionCookieById).not.toHaveBeenCalled();
    });
  });

  describe("successful verification", () => {
    test("completes the flow with sessionId and requestId when a requestId is present", async () => {
      const result = await sendRecoveryCode({
        loginName: "test@example.com",
        organization: "org123",
        requestId: "oidc_123",
        code: "abcde-fghij",
      });

      expect(mockCompleteFlowOrGetUrl).toHaveBeenCalledWith(
        { sessionId: "session123", requestId: "oidc_123", organization: "org123" },
        undefined,
      );
      expect(result).toEqual({ redirect: "/signedin" });
      expect(mockRecordAuthSuccess).toHaveBeenCalledWith("recovery_code", "org123");
    });

    test("completes the flow with the login name when no requestId is present", async () => {
      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(mockCompleteFlowOrGetUrl).toHaveBeenCalledWith(
        { loginName: "test@example.com", organization: "org123" },
        undefined,
      );
      expect(result).toEqual({ redirect: "/signedin" });
      expect(mockRecordAuthSuccess).toHaveBeenCalledWith("recovery_code", undefined);
    });

    test("passes the default redirect uri from the login settings to the flow", async () => {
      mockGetLoginSettings.mockResolvedValue({
        secondFactorCheckLifetime: { seconds: BigInt(600), nanos: 0 },
        defaultRedirectUri: "https://example.com/app",
      });

      await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(mockCompleteFlowOrGetUrl).toHaveBeenCalledWith(expect.anything(), "https://example.com/app");
    });

    test("returns a samlData response unchanged", async () => {
      const samlData = { url: "https://sp.example.com/acs", fields: { SAMLResponse: "response" } };
      mockCompleteFlowOrGetUrl.mockResolvedValue({ samlData });

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(result).toEqual({ samlData });
      expect(mockRecordAuthSuccess).toHaveBeenCalledWith("recovery_code", undefined);
    });
  });

  describe("check payload", () => {
    test("trims the submitted code", async () => {
      await sendRecoveryCode({ loginName: "test@example.com", code: "  abcde-fghij  " });

      expect(mockSetSessionAndUpdateCookie).toHaveBeenCalledWith(
        expect.objectContaining({ checks: { recoveryCode: { code: "abcde-fghij" } } }),
      );
    });

    test("does not change the casing or the dashes of the code", async () => {
      await sendRecoveryCode({ loginName: "test@example.com", code: "AbCde-FGhij" });

      expect(mockSetSessionAndUpdateCookie).toHaveBeenCalledWith(
        expect.objectContaining({ checks: { recoveryCode: { code: "AbCde-FGhij" } } }),
      );
    });

    test("uses the configured second factor check lifetime", async () => {
      await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(mockSetSessionAndUpdateCookie).toHaveBeenCalledWith(
        expect.objectContaining({ lifetime: { seconds: BigInt(600), nanos: 0 } }),
      );
    });

    test("falls back to a 24 hour lifetime when the settings have no second factor lifetime", async () => {
      mockGetLoginSettings.mockResolvedValue({});

      await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(mockSetSessionAndUpdateCookie).toHaveBeenCalledWith(
        expect.objectContaining({ lifetime: { seconds: BigInt(60 * 60 * 24), nanos: 0 } }),
      );
    });

    test("forwards the requestId to the session update", async () => {
      await sendRecoveryCode({ loginName: "test@example.com", requestId: "oidc_123", code: "abcde-fghij" });

      expect(mockSetSessionAndUpdateCookie).toHaveBeenCalledWith(
        expect.objectContaining({ recentCookie: sessionCookie, requestId: "oidc_123" }),
      );
    });
  });

  describe("error mapping", () => {
    test("maps a translated InvalidArgument to invalidCode", async () => {
      mockSetSessionAndUpdateCookie.mockRejectedValue(
        classifiedError("Invalid recovery code provided (DOMAIN-6uvh0)", Code.InvalidArgument),
      );

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "wrong" });

      expect(result).toEqual({ error: "verify.errors.invalidCode" });
      expect(mockRecordAuthFailure).toHaveBeenCalledWith("recovery_code", "invalid_code", undefined);
    });

    test("maps an untranslated InvalidArgument to invalidCode", async () => {
      mockSetSessionAndUpdateCookie.mockRejectedValue(
        classifiedError("Errors.User.Code.Invalid (DOMAIN-9xrr0)", Code.InvalidArgument),
      );

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "wrong" });

      expect(result).toEqual({ error: "verify.errors.invalidCode" });
    });

    test("maps a locked user to userLocked", async () => {
      mockSetSessionAndUpdateCookie.mockRejectedValue(
        classifiedError("User is locked (COMMAND-2w6oa)", Code.FailedPrecondition),
      );

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(result).toEqual({ error: "verify.errors.userLocked" });
      expect(mockRecordAuthFailure).toHaveBeenCalledWith("recovery_code", "user_locked", undefined);
    });

    test("maps a locked user to userLocked even when the message is translated", async () => {
      // the backend translates the message before building the status, so only the id is stable
      mockSetSessionAndUpdateCookie.mockRejectedValue(
        classifiedError("Der Benutzer ist gesperrt (COMMAND-ASV12)", Code.FailedPrecondition),
      );

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(result).toEqual({ error: "verify.errors.userLocked" });
    });

    test("maps a locked user to userLocked from the untranslated key", async () => {
      mockSetSessionAndUpdateCookie.mockRejectedValue(classifiedError("Errors.User.Locked", Code.FailedPrecondition));

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(result).toEqual({ error: "verify.errors.userLocked" });
    });

    test("maps missing recovery codes to notReady", async () => {
      mockSetSessionAndUpdateCookie.mockRejectedValue(
        classifiedError("Der Benutzer hat keine Wiederherstellungscodes (COMMAND-84rgg)", Code.FailedPrecondition),
      );

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(result).toEqual({ error: "verify.errors.notReady" });
      expect(mockRecordAuthFailure).toHaveBeenCalledWith("recovery_code", "not_ready", undefined);
    });

    test("maps missing recovery codes to notReady from the untranslated key", async () => {
      mockSetSessionAndUpdateCookie.mockRejectedValue(
        classifiedError("Errors.User.MFA.RecoveryCodes.NotReady", Code.FailedPrecondition),
      );

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(result).toEqual({ error: "verify.errors.notReady" });
    });

    test("maps an unrelated FailedPrecondition to couldNotVerifyCode", async () => {
      mockSetSessionAndUpdateCookie.mockRejectedValue(
        classifiedError("Errors.Session.Terminated (COMMAND-SAF3t)", Code.FailedPrecondition),
      );

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(result).toEqual({ error: "verify.errors.couldNotVerifyCode" });
      expect(mockRecordAuthFailure).toHaveBeenCalledWith("recovery_code", "session_update_failed", undefined);
    });

    test("maps a credentials check error shape to invalidCode", async () => {
      // defensive: passwordAttemptsHandler rethrows this shape when details carry failedAttempts
      mockSetSessionAndUpdateCookie.mockRejectedValue({ error: "Failed to authenticate", failedAttempts: 1 });

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(result).toEqual({ error: "verify.errors.invalidCode" });
      expect(mockRecordAuthFailure).toHaveBeenCalledWith("recovery_code", "invalid_code", undefined);
    });

    test("maps an unknown error to couldNotVerifyCode", async () => {
      mockSetSessionAndUpdateCookie.mockRejectedValue(new Error("boom"));

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(result).toEqual({ error: "verify.errors.couldNotVerifyCode" });
      expect(mockRecordAuthFailure).toHaveBeenCalledWith("recovery_code", "session_update_failed", undefined);
    });

    test("never completes the flow after a failed check", async () => {
      mockSetSessionAndUpdateCookie.mockRejectedValue(
        classifiedError("Invalid recovery code provided (DOMAIN-6uvh0)", Code.InvalidArgument),
      );

      await sendRecoveryCode({ loginName: "test@example.com", code: "wrong" });

      expect(mockCompleteFlowOrGetUrl).not.toHaveBeenCalled();
    });
  });

  describe("post-check validation", () => {
    test("returns couldNotVerifyCode when the updated session has no login name", async () => {
      mockSetSessionAndUpdateCookie.mockResolvedValue({ id: "session123", factors: {} });

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(result).toEqual({ error: "verify.errors.couldNotVerifyCode" });
      expect(mockRecordAuthFailure).toHaveBeenCalledWith("recovery_code", "session_invalid", undefined);
    });

    test("returns the flow error when the flow could not be completed", async () => {
      mockCompleteFlowOrGetUrl.mockResolvedValue({ error: "Could not complete flow" });

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(result).toEqual({ error: "Could not complete flow" });
      expect(mockRecordAuthFailure).toHaveBeenCalledWith("recovery_code", "flow_error", undefined);
      expect(mockRecordAuthSuccess).not.toHaveBeenCalled();
    });

    test("returns couldNotDetermineRedirect when the flow returns nothing", async () => {
      mockCompleteFlowOrGetUrl.mockResolvedValue(undefined);

      const result = await sendRecoveryCode({ loginName: "test@example.com", code: "abcde-fghij" });

      expect(result).toEqual({ error: "verify.errors.couldNotDetermineRedirect" });
      expect(mockRecordAuthFailure).toHaveBeenCalledWith("recovery_code", "navigation_failed", undefined);
    });
  });
});
