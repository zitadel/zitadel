import { create, timestampFromDate } from "@zitadel/client";
import { AuthRequestSchema } from "@zitadel/proto/zitadel/oidc/v2/authorization_pb";
import { CreateCallbackResponseSchema, GetAuthRequestResponseSchema } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { SessionSchema } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { ListAuthenticationMethodTypesResponseSchema } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { loginWithOIDCAndSession } from "./oidc";
import { sendLoginname } from "./server/loginname";
import { findValidSession, isSessionValid } from "./session";
import { createCallback, getAuthRequest, getLoginSettings, listAuthenticationMethodTypes } from "./zitadel";

vi.mock("./cookies", () => ({}));
vi.mock("./verify-helper", () => ({ shouldEnforceMFA: () => false }));
vi.mock("./server/loginname", () => ({ sendLoginname: vi.fn() }));
vi.mock("./zitadel", () => ({
  createCallback: vi.fn(),
  getAuthRequest: vi.fn(),
  getLoginSettings: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
}));

describe("OIDC session authentication age", () => {
  const serviceConfig = { baseUrl: "https://api.example.com" };
  const requestCreatedAt = new Date("2026-01-01T12:00:00Z");
  const authRequest = create(AuthRequestSchema, {
    creationDate: timestampFromDate(requestCreatedAt),
    maxAge: { seconds: BigInt(10) },
  });
  const session = create(SessionSchema, {
    id: "session-id",
    factors: {
      user: { id: "user-id", organizationId: "org-id", loginName: "user@example.com" },
      password: { verifiedAt: timestampFromDate(new Date("2026-01-01T11:00:00Z")) },
    },
  });

  beforeEach(() => {
    vi.clearAllMocks();
    vi.useFakeTimers();
    vi.setSystemTime(requestCreatedAt);
    vi.mocked(getLoginSettings).mockResolvedValue(undefined);
    vi.mocked(listAuthenticationMethodTypes).mockResolvedValue(create(ListAuthenticationMethodTypesResponseSchema));
    vi.mocked(getAuthRequest).mockResolvedValue(create(GetAuthRequestResponseSchema, { authRequest }));
    vi.mocked(createCallback).mockResolvedValue(
      create(CreateCallbackResponseSchema, { callbackUrl: "https://rp.example.com/callback" }),
    );
    vi.mocked(sendLoginname).mockResolvedValue({ redirect: "/password?requestId=oidc_request-id" });
    vi.stubEnv("EMAIL_VERIFICATION", "false");
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.unstubAllEnvs();
  });

  test("does not select a stale session for an auth request with max_age", async () => {
    expect(await findValidSession({ serviceConfig, sessions: [session], authRequest })).toBeUndefined();
  });

  test.each([
    { name: "absent max_age", seconds: undefined, offset: -3600, valid: true },
    { name: "recent authentication", seconds: BigInt(10), offset: -5, valid: true },
    { name: "exact age limit", seconds: BigInt(10), offset: -10, valid: true },
    { name: "stale authentication", seconds: BigInt(10), offset: -11, valid: false },
    { name: "zero before reauthentication", seconds: BigInt(0), offset: -1, valid: false },
    { name: "zero after reauthentication", seconds: BigInt(0), offset: 1, valid: true },
  ])("$name", async ({ seconds, offset, valid }) => {
    const request = create(AuthRequestSchema, {
      creationDate: authRequest.creationDate,
      maxAge: seconds === undefined ? undefined : { seconds },
    });
    const candidate = create(SessionSchema, {
      factors: {
        user: session.factors?.user,
        password: { verifiedAt: timestampFromDate(new Date(requestCreatedAt.getTime() + offset * 1000)) },
      },
    });

    expect(await isSessionValid({ serviceConfig, session: candidate, authRequest: request })).toBe(valid);
  });

  test("uses the latest authentication factor, matching the backend auth_time", async () => {
    const candidate = create(SessionSchema, {
      factors: {
        user: session.factors?.user,
        password: session.factors?.password,
        intent: { verifiedAt: timestampFromDate(requestCreatedAt) },
      },
    });

    expect(await isSessionValid({ serviceConfig, session: candidate, authRequest })).toBe(true);
  });

  test("does not treat a recent session change or user lookup as authentication", async () => {
    const candidate = create(SessionSchema, {
      id: session.id,
      changeDate: timestampFromDate(requestCreatedAt),
      factors: {
        password: session.factors?.password,
        user: { id: "user-id", organizationId: "org-id", verifiedAt: timestampFromDate(requestCreatedAt) },
      },
    });

    expect(await isSessionValid({ serviceConfig, session: candidate, authRequest })).toBe(false);
  });

  test("fails closed when the auth request creation date is missing", async () => {
    const request = create(AuthRequestSchema, { maxAge: { seconds: BigInt(10) } });

    expect(await isSessionValid({ serviceConfig, session, authRequest: request })).toBe(false);
  });

  test.each([
    { name: "password", response: { redirect: "/password?requestId=oidc_request-id" } },
    { name: "SAML IdP", response: { samlData: { url: "https://idp.example.com/sso", fields: { SAMLRequest: "request" } } } },
    { name: "reauthentication failure", response: { error: "Could not reauthenticate" } },
  ])("reauthenticates a stale account selected from the account picker: $name", async ({ response }) => {
    vi.mocked(sendLoginname).mockResolvedValue(response);
    const result = await loginWithOIDCAndSession({
      serviceConfig,
      authRequest: "request-id",
      sessionId: session.id,
      sessions: [session],
      sessionCookies: [
        {
          id: session.id,
          token: "session-token",
          loginName: "user@example.com",
          creationTs: "",
          changeTs: "",
          expirationTs: "",
        },
      ],
    });

    expect(result).toEqual(response);
    expect(sendLoginname).toHaveBeenCalledWith({
      loginName: "user@example.com",
      organization: "org-id",
      requestId: "oidc_request-id",
    });
    expect(createCallback).not.toHaveBeenCalled();
  });
});
