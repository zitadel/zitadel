import { addU2F, verifyU2F } from "@/lib/actions/u2f";
import { ClassifiedConnectError } from "@/lib/server/grpc/interceptors/error-classification";
import { Code, ConnectError } from "@connectrpc/connect";
import { beforeEach, describe, expect, test, vi } from "vitest";

vi.mock("@zitadel/client", () => ({
  // `create(schema, data)` is only used to build the verify request — return the data.
  create: vi.fn((_schema: unknown, data: unknown) => data),
}));

vi.mock("@/lib/server/zitadel", () => ({
  getSession: vi.fn(),
  registerU2F: vi.fn(),
  verifyU2FRegistration: vi.fn(),
}));

vi.mock("@/lib/server/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://example.com" } })),
}));

vi.mock("@/lib/server/host", () => ({
  getPublicHost: vi.fn(() => "test.com"),
}));

vi.mock("@/lib/server/cookies", () => ({
  getSessionCookieById: vi.fn(),
}));

vi.mock("@/lib/server/enrollment-guard", () => ({
  getSetupAuthorizationError: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(async () => (key: string) => key),
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(() => new Headers()),
}));

vi.mock("next/server", () => ({
  userAgent: vi.fn(() => ({ browser: {}, device: {}, os: {} })),
}));

const userOnlySession = {
  session: { id: "session-1", factors: { user: { id: "victim-1", loginName: "victim@example.com" } } },
};

describe("addU2F", () => {
  let getSessionCookieById: any;
  let getSession: any;
  let registerU2F: any;
  let getSetupAuthorizationError: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    getSessionCookieById = vi.mocked((await import("@/lib/server/cookies")).getSessionCookieById);
    getSession = vi.mocked((await import("@/lib/server/zitadel")).getSession);
    registerU2F = vi.mocked((await import("@/lib/server/zitadel")).registerU2F);
    getSetupAuthorizationError = vi.mocked((await import("@/lib/server/enrollment-guard")).getSetupAuthorizationError);

    getSessionCookieById.mockResolvedValue({ id: "session-1", token: "token-1" });
    getSession.mockResolvedValue(userOnlySession);
  });

  test("rejects enrollment on an unauthenticated (identify-only) session without registering", async () => {
    getSetupAuthorizationError.mockResolvedValue("You have to authenticate or have a valid User Verification Check");

    const result = await addU2F({ sessionId: "session-1" });

    expect(result).toEqual({ error: "You have to authenticate or have a valid User Verification Check" });
    expect(getSetupAuthorizationError).toHaveBeenCalledWith({
      serviceConfig: { baseUrl: "https://example.com" },
      session: userOnlySession.session,
      userId: "victim-1",
      method: "u2f",
    });
    expect(registerU2F).not.toHaveBeenCalled();
  });

  test("registers the credential when the session is authorized", async () => {
    getSetupAuthorizationError.mockResolvedValue(null);
    registerU2F.mockResolvedValue({ u2fId: "u2f-1" });

    const result = await addU2F({ sessionId: "session-1" });

    expect(registerU2F).toHaveBeenCalledWith({
      serviceConfig: { baseUrl: "https://example.com" },
      userId: "victim-1",
      domain: "test.com",
    });
    expect(result).toEqual({ u2fId: "u2f-1" });
  });

  test("returns error and never calls the guard when no session cookie exists", async () => {
    getSessionCookieById.mockResolvedValue(undefined);

    const result = await addU2F({ sessionId: "missing" });

    expect(result).toEqual({ error: "errors.couldNotLoadSession" });
    expect(getSetupAuthorizationError).not.toHaveBeenCalled();
    expect(registerU2F).not.toHaveBeenCalled();
  });
});

describe("verifyU2F", () => {
  let getSessionCookieById: any;
  let getSession: any;
  let verifyU2FRegistration: any;
  let getSetupAuthorizationError: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    getSessionCookieById = vi.mocked((await import("@/lib/server/cookies")).getSessionCookieById);
    getSession = vi.mocked((await import("@/lib/server/zitadel")).getSession);
    verifyU2FRegistration = vi.mocked((await import("@/lib/server/zitadel")).verifyU2FRegistration);
    getSetupAuthorizationError = vi.mocked((await import("@/lib/server/enrollment-guard")).getSetupAuthorizationError);

    getSessionCookieById.mockResolvedValue({ id: "session-1", token: "token-1" });
    getSession.mockResolvedValue(userOnlySession);
  });

  test("rejects verification on an unauthenticated session without persisting the credential", async () => {
    getSetupAuthorizationError.mockResolvedValue("User Verification Check has to be done");

    const result = await verifyU2F({
      u2fId: "u2f-1",
      passkeyName: "my key",
      publicKeyCredential: {},
      sessionId: "session-1",
    });

    expect(result).toEqual({ error: "User Verification Check has to be done" });
    expect(verifyU2FRegistration).not.toHaveBeenCalled();
  });

  test("verifies the registration when the session is authorized", async () => {
    getSetupAuthorizationError.mockResolvedValue(null);
    verifyU2FRegistration.mockResolvedValue({ done: true });

    const result = await verifyU2F({
      u2fId: "u2f-1",
      passkeyName: "my key",
      publicKeyCredential: {},
      sessionId: "session-1",
    });

    expect(verifyU2FRegistration).toHaveBeenCalledTimes(1);
    expect(result).toEqual({ done: true });
  });
});

describe("U2F calls ZITADEL rejects", () => {
  const apiError = (code: Code) => new ClassifiedConnectError(new ConnectError("rejected", code));
  let zitadel: any;
  let getSetupAuthorizationError: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    vi.mocked((await import("@/lib/server/cookies")).getSessionCookieById).mockResolvedValue({
      id: "session-1",
      token: "token-1",
    } as any);
    zitadel = vi.mocked(await import("@/lib/server/zitadel"));
    getSetupAuthorizationError = vi.mocked((await import("@/lib/server/enrollment-guard")).getSetupAuthorizationError);
    zitadel.getSession.mockResolvedValue(userOnlySession);
    getSetupAuthorizationError.mockResolvedValue(undefined);
  });

  test.each([Code.NotFound, Code.PermissionDenied])(
    "answers a session the server no longer has (code %s) with a message",
    async (code) => {
      zitadel.getSession.mockRejectedValue(apiError(code));

      expect(await addU2F({ sessionId: "session-1" })).toEqual({ error: "errors.couldNotLoadSession" });
      expect(await verifyU2F({ u2fId: "u2f-1", publicKeyCredential: {}, sessionId: "session-1" })).toEqual({
        error: "errors.couldNotLoadSession",
      });
    },
  );

  test("answers a rejected registration with a message", async () => {
    zitadel.registerU2F.mockRejectedValue(apiError(Code.FailedPrecondition));

    expect(await addU2F({ sessionId: "session-1" })).toEqual({ error: "errors.couldNotRegister" });
  });

  test("answers a rejected attestation with a message", async () => {
    zitadel.verifyU2FRegistration.mockRejectedValue(apiError(Code.InvalidArgument));

    expect(await verifyU2F({ u2fId: "u2f-1", publicKeyCredential: {}, sessionId: "session-1" })).toEqual({
      error: "errors.couldNotVerify",
    });
  });

  test.each([Code.Internal, Code.Unauthenticated])("rethrows code %s", async (code) => {
    const error = apiError(code);
    zitadel.verifyU2FRegistration.mockRejectedValue(error);

    await expect(verifyU2F({ u2fId: "u2f-1", publicKeyCredential: {}, sessionId: "session-1" })).rejects.toBe(error);
  });
});
