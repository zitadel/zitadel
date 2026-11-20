import { registerUser, registerUserAndLinkToIDP } from "@/lib/actions/register";
import { beforeEach, describe, expect, test, vi } from "vitest";

vi.mock("@/lib/server/zitadel", () => ({
  addHumanUser: vi.fn(),
  addIDPLink: vi.fn(),
  getLoginSettings: vi.fn(),
  getUserByID: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
  retrieveIDPIntent: vi.fn(),
}));

vi.mock("@/lib/server/session-cookie", () => ({
  createSessionAndUpdateCookie: vi.fn(),
  createSessionForIdpAndUpdateCookie: vi.fn(),
}));

vi.mock("@/lib/actions/flow", () => ({
  completeFlowOrGetUrl: vi.fn(),
}));

vi.mock("@/lib/server/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://example.com" } })),
}));

vi.mock("@/lib/server/verification-check", () => ({
  setUserVerificationCheckCookie: vi.fn(),
}));

vi.mock("@/lib/server/verify-helper", () => ({
  checkEmailVerification: vi.fn(),
  checkMFAFactors: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(async () => (key: string) => key),
}));

// fphg: the registration must come from the browser that started the IdP flow
vi.mock("@/lib/server/fingerprint", () => ({
  getBindingIdCookie: vi.fn(async () => ({ value: "binding-id" })),
}));

vi.mock("@/lib/server/idp-flow-proof", () => ({
  isIdpFlowProofValid: vi.fn(() => true),
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(() => new Headers()),
}));

const command = {
  email: "new@example.com",
  firstName: "New",
  lastName: "User",
  organization: "org-1",
  idpIntent: { idpIntentId: "intent-1", idpIntentToken: "intent-token-1" },
  idpUserName: "typed-name",
};

describe("registerUserAndLinkToIDP", () => {
  let zitadel: any;
  let cookieHelpers: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    zitadel = vi.mocked(await import("@/lib/server/zitadel"));
    cookieHelpers = vi.mocked(await import("@/lib/server/session-cookie"));

    zitadel.getLoginSettings.mockResolvedValue({ allowRegister: true });
    zitadel.addHumanUser.mockResolvedValue({ userId: "new-user" });
    zitadel.addIDPLink.mockResolvedValue({});
    cookieHelpers.createSessionForIdpAndUpdateCookie.mockResolvedValue(undefined);
  });

  test("links the external identity from the intent, not identity fields from the request", async () => {
    zitadel.retrieveIDPIntent.mockResolvedValue({
      idpInformation: { idpId: "idp-from-intent", userId: "external-user-from-intent", userName: "intent-name" },
    });

    // a direct POST can still send the fields the action used to accept
    await registerUserAndLinkToIDP({ ...command, idpId: "victim-idp", idpUserId: "victim-external-user" } as any);

    expect(zitadel.retrieveIDPIntent).toHaveBeenCalledWith(
      expect.objectContaining({ id: "intent-1", token: "intent-token-1" }),
    );
    expect(zitadel.addIDPLink).toHaveBeenCalledWith(
      expect.objectContaining({
        idp: { id: "idp-from-intent", userId: "external-user-from-intent", userName: "intent-name" },
        userId: "new-user",
      }),
    );
  });

  test("uses the typed user name only when the IDP provided none", async () => {
    zitadel.retrieveIDPIntent.mockResolvedValue({
      idpInformation: { idpId: "idp-1", userId: "external-user", userName: "" },
    });

    await registerUserAndLinkToIDP(command);

    expect(zitadel.addIDPLink).toHaveBeenCalledWith(
      expect.objectContaining({ idp: { id: "idp-1", userId: "external-user", userName: "typed-name" } }),
    );
  });

  test("creates no user when the intent can't be retrieved", async () => {
    zitadel.retrieveIDPIntent.mockRejectedValue(new Error("invalid intent token"));

    const res = await registerUserAndLinkToIDP(command);

    expect(res).toEqual({ error: "errors.couldNotLinkIDP" });
    expect(zitadel.addHumanUser).not.toHaveBeenCalled();
    expect(zitadel.addIDPLink).not.toHaveBeenCalled();
  });
});

describe("registerUser", () => {
  let zitadel: any;
  let cookieHelpers: any;
  let verifyHelper: any;
  let flow: any;

  const registration = {
    email: "jane@example.com",
    firstName: "Jane",
    lastName: "Doe",
    password: "Password1!",
    organization: "org-1",
  };

  beforeEach(async () => {
    vi.clearAllMocks();
    zitadel = vi.mocked(await import("@/lib/server/zitadel"));
    cookieHelpers = vi.mocked(await import("@/lib/server/session-cookie"));
    verifyHelper = vi.mocked(await import("@/lib/server/verify-helper"));
    flow = vi.mocked(await import("@/lib/actions/flow"));

    zitadel.getLoginSettings.mockResolvedValue({
      allowRegister: true,
      allowLocalAuthentication: true,
      passwordCheckLifetime: { seconds: BigInt(864000), nanos: 0 },
    });
    zitadel.addHumanUser.mockResolvedValue({ userId: "new-user" });
    zitadel.getUserByID.mockResolvedValue({ user: { type: { case: "human", value: {} } } });
    cookieHelpers.createSessionAndUpdateCookie.mockResolvedValue({
      session: {
        id: "session-1",
        factors: { user: { id: "new-user", loginName: "jane@example.com", organizationId: "org-1" } },
      },
    });
    verifyHelper.checkEmailVerification.mockResolvedValue(undefined);
    flow.completeFlowOrGetUrl.mockResolvedValue({ redirect: "/signedin?loginName=jane%40example.com" });
  });

  test("creates the user with the password and signs in with a password check", async () => {
    const res = await registerUser({ ...registration, requestId: "oidc_1" });

    expect(zitadel.addHumanUser).toHaveBeenCalledWith(
      expect.objectContaining({
        email: "jane@example.com",
        firstName: "Jane",
        lastName: "Doe",
        password: "Password1!",
        organization: "org-1",
      }),
    );
    const sessionCommand = cookieHelpers.createSessionAndUpdateCookie.mock.calls[0][0];
    expect(sessionCommand.checks.user.search).toEqual({ case: "userId", value: "new-user" });
    expect(sessionCommand.checks.password.password).toBe("Password1!");
    expect(sessionCommand.lifetime).toEqual({ seconds: BigInt(864000), nanos: 0 });
    expect(sessionCommand.requestId).toBe("oidc_1");
    expect(res).toEqual({ redirect: "/signedin?loginName=jane%40example.com" });
  });

  test("creates no user without a password", async () => {
    const res = await registerUser({ ...registration, password: "" });

    expect(res).toEqual({ error: "password.required.password" });
    expect(zitadel.addHumanUser).not.toHaveBeenCalled();
  });

  test("creates no user when the registration is disabled", async () => {
    zitadel.getLoginSettings.mockResolvedValue({ allowRegister: false, allowLocalAuthentication: true });

    const res = await registerUser(registration);

    expect(res).toEqual({ error: "disabled.description" });
    expect(zitadel.addHumanUser).not.toHaveBeenCalled();
  });

  test("creates no user when passwords are not allowed", async () => {
    zitadel.getLoginSettings.mockResolvedValue({ allowRegister: true, allowLocalAuthentication: false });

    const res = await registerUser(registration);

    expect(res).toEqual({ error: "errors.localAuthenticationNotAllowed" });
    expect(zitadel.addHumanUser).not.toHaveBeenCalled();
  });

  test("continues with the email verification when it is required", async () => {
    verifyHelper.checkEmailVerification.mockResolvedValue({ redirect: "/verify?loginName=jane%40example.com" });

    const res = await registerUser(registration);

    expect(res).toEqual({ redirect: "/verify?loginName=jane%40example.com" });
    expect(flow.completeFlowOrGetUrl).not.toHaveBeenCalled();
  });
});
