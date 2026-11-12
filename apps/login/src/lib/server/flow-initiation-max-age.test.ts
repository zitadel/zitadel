import { create, timestampFromDate } from "@zitadel/client";
import { AuthRequestSchema, ErrorReason, Prompt } from "@zitadel/proto/zitadel/oidc/v2/authorization_pb";
import { CreateCallbackResponseSchema, GetAuthRequestResponseSchema } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { SessionSchema } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { ListAuthenticationMethodTypesResponseSchema } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { NextRequest } from "next/server";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import {
  createCallback,
  getAuthRequest,
  getLoginSettings,
  getSecuritySettings,
  listAuthenticationMethodTypes,
} from "../zitadel";
import { handleOIDCFlowInitiation } from "./flow-initiation";
import { sendLoginname } from "./loginname";

vi.mock("../cookies", () => ({ getLanguageCookie: vi.fn(), setLanguageCookie: vi.fn() }));
vi.mock("../verify-helper", () => ({ shouldEnforceMFA: () => false }));
vi.mock("./loginname", () => ({ sendLoginname: vi.fn() }));
vi.mock("../zitadel", () => ({
  createCallback: vi.fn(),
  getAuthRequest: vi.fn(),
  getLoginSettings: vi.fn(),
  getSecuritySettings: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
}));

describe("OIDC max_age flow initiation", () => {
  const requestCreatedAt = new Date("2026-01-01T12:00:00Z");
  const authRequest = create(AuthRequestSchema, {
    id: "request-id",
    creationDate: timestampFromDate(requestCreatedAt),
    maxAge: { seconds: BigInt(10) },
    prompt: [Prompt.NONE],
  });
  const session = create(SessionSchema, {
    id: "session-id",
    factors: {
      user: { id: "user-id", loginName: "user@example.com", organizationId: "org-id" },
      password: { verifiedAt: timestampFromDate(new Date("2026-01-01T11:00:00Z")) },
    },
  });

  beforeEach(() => {
    vi.clearAllMocks();
    vi.useFakeTimers();
    vi.setSystemTime(requestCreatedAt);
    vi.stubEnv("NEXT_PUBLIC_BASE_PATH", "/ui/v2/login");
    vi.stubEnv("EMAIL_VERIFICATION", "false");
    vi.mocked(getLoginSettings).mockResolvedValue(undefined);
    vi.mocked(listAuthenticationMethodTypes).mockResolvedValue(create(ListAuthenticationMethodTypesResponseSchema));
    vi.mocked(getAuthRequest).mockResolvedValue(create(GetAuthRequestResponseSchema, { authRequest }));
    vi.mocked(getSecuritySettings).mockResolvedValue(undefined);
    vi.mocked(createCallback).mockResolvedValue(
      create(CreateCallbackResponseSchema, {
        callbackUrl: "https://rp.example.com/callback?error=login_required&state=state",
      }),
    );
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.unstubAllEnvs();
  });

  test.each([true, false])("prompt=none returns login_required (has stale session: %s)", async (hasSession) => {
    const response = await handleOIDCFlowInitiation({
      serviceConfig: { baseUrl: "https://api.example.com" },
      requestId: "oidc_request-id",
      sessions: hasSession ? [session] : [],
      sessionCookies: [
        { id: session.id, token: "token", loginName: "user@example.com", creationTs: "", changeTs: "", expirationTs: "" },
      ],
      request: new NextRequest("https://login.example.com/ui/v2/login/login?requestId=oidc_request-id", {
        headers: { host: "login.example.com" },
      }),
    });

    expect(response.status).toBe(307);
    expect(response.headers.get("location")).toBe("https://rp.example.com/callback?error=login_required&state=state");
    expect(createCallback).toHaveBeenCalledWith(
      expect.objectContaining({
        req: expect.objectContaining({
          authRequestId: "request-id",
          callbackKind: { case: "error", value: expect.objectContaining({ error: ErrorReason.LOGIN_REQUIRED }) },
        }),
      }),
    );
    expect(sendLoginname).not.toHaveBeenCalled();
  });

  test("interactive requests send stale sessions to account selection under the login base path", async () => {
    vi.mocked(getAuthRequest).mockResolvedValue(
      create(GetAuthRequestResponseSchema, {
        authRequest: create(AuthRequestSchema, { ...authRequest, prompt: [] }),
      }),
    );
    const response = await handleOIDCFlowInitiation({
      serviceConfig: { baseUrl: "https://api.example.com" },
      requestId: "oidc_request-id",
      sessions: [session],
      sessionCookies: [],
      request: new NextRequest("https://login.example.com/ui/v2/login/login?requestId=oidc_request-id", {
        headers: { host: "login.example.com" },
      }),
    });

    const location = new URL(response.headers.get("location")!);
    expect(location.pathname).toBe("/ui/v2/login/accounts");
    expect(location.searchParams.get("requestId")).toBe("oidc_request-id");
    expect(createCallback).not.toHaveBeenCalled();
  });
});
