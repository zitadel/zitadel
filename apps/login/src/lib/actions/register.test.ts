import { registerUser, registerUserAndLinkToIDP } from "@/lib/actions/register";
import { ClassifiedConnectError } from "@/lib/server/grpc/interceptors/error-classification";
import { Code, ConnectError } from "@connectrpc/connect";
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

const apiError = (code: Code, message: string) => new ClassifiedConnectError(new ConnectError(message, code));

const serverFaults = [
  ["a server fault", apiError(Code.Internal, "boom")],
  ["the API rejecting the login's own credentials", apiError(Code.Unauthenticated, "Token invalid")],
] as const;

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

describe("registerUserAndLinkToIDP rejections", () => {
  let zitadel: any;
  let cookieHelpers: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    zitadel = vi.mocked(await import("@/lib/server/zitadel"));
    cookieHelpers = vi.mocked(await import("@/lib/server/session-cookie"));
    const { completeFlowOrGetUrl } = await import("@/lib/actions/flow");

    zitadel.getLoginSettings.mockResolvedValue({ allowRegister: true });
    zitadel.retrieveIDPIntent.mockResolvedValue({
      idpInformation: { idpId: "idp-1", userId: "external-user", userName: "intent-name" },
    });
    zitadel.addHumanUser.mockResolvedValue({ userId: "new-user" });
    zitadel.addIDPLink.mockResolvedValue({});
    zitadel.listAuthenticationMethodTypes.mockResolvedValue({ authMethodTypes: [] });
    cookieHelpers.createSessionForIdpAndUpdateCookie.mockResolvedValue({
      id: "session-1",
      factors: { user: { id: "new-user", loginName: "new@example.com", organizationId: "org-1" } },
    });
    vi.mocked(completeFlowOrGetUrl).mockResolvedValue({ redirect: "/signedin" } as any);
  });

  test("says the user exists when the email address is taken", async () => {
    zitadel.addHumanUser.mockRejectedValue(apiError(Code.AlreadyExists, "User already exists (V3-DKcYh)"));

    await expect(registerUserAndLinkToIDP(command)).resolves.toEqual({ error: "errors.userAlreadyExists" });
    expect(zitadel.addIDPLink).not.toHaveBeenCalled();
  });

  test("answers other rejected user data with a message", async () => {
    zitadel.addHumanUser.mockRejectedValue(apiError(Code.InvalidArgument, "Email is invalid (EMAIL-599BI)"));

    await expect(registerUserAndLinkToIDP(command)).resolves.toEqual({ error: "errors.couldNotCreateUser" });
  });

  test("answers an external identity linked to another user with a message", async () => {
    zitadel.addIDPLink.mockRejectedValue(apiError(Code.AlreadyExists, "External IDP already exists (COMMAND-Sgfd2)"));

    await expect(registerUserAndLinkToIDP(command)).resolves.toEqual({ error: "errors.couldNotLinkIDP" });
  });

  test("answers an expired or used intent with a message", async () => {
    cookieHelpers.createSessionForIdpAndUpdateCookie.mockRejectedValue(
      apiError(Code.FailedPrecondition, "Intent has expired (COMMAND-SAf42)"),
    );

    await expect(registerUserAndLinkToIDP(command)).resolves.toEqual({ error: "errors.couldNotCreateSession" });
  });

  test("answers a rejected lookup of the new user's methods with a message", async () => {
    zitadel.listAuthenticationMethodTypes.mockRejectedValue(apiError(Code.NotFound, "User not found (QUERY-Dfbg2)"));

    await expect(registerUserAndLinkToIDP(command)).resolves.toEqual({ error: "errors.couldNotRegisterUser" });
  });

  test.each(serverFaults)("rethrows %s when creating the user", async (_, error) => {
    zitadel.addHumanUser.mockRejectedValue(error);

    await expect(registerUserAndLinkToIDP(command)).rejects.toBe(error);
  });

  test.each(serverFaults)("rethrows %s when linking the identity", async (_, error) => {
    zitadel.addIDPLink.mockRejectedValue(error);

    await expect(registerUserAndLinkToIDP(command)).rejects.toBe(error);
  });

  test.each(serverFaults)("rethrows %s when creating the session", async (_, error) => {
    cookieHelpers.createSessionForIdpAndUpdateCookie.mockRejectedValue(error);

    await expect(registerUserAndLinkToIDP(command)).rejects.toBe(error);
  });
});

describe("registerUser", () => {
  const registerCommand = {
    email: "new@example.com",
    firstName: "New",
    lastName: "User",
    password: "Password1!",
    organization: "org-1",
  };

  let zitadel: any;
  let cookieHelpers: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    zitadel = vi.mocked(await import("@/lib/server/zitadel"));
    cookieHelpers = vi.mocked(await import("@/lib/server/session-cookie"));

    zitadel.getLoginSettings.mockResolvedValue({ allowRegister: true, allowLocalAuthentication: true });
    zitadel.addHumanUser.mockResolvedValue({ userId: "new-user" });
  });

  test("says the user exists when the email address is taken", async () => {
    zitadel.addHumanUser.mockRejectedValue(apiError(Code.AlreadyExists, "User already exists (V3-DKcYh)"));

    await expect(registerUser(registerCommand)).resolves.toEqual({ error: "errors.userAlreadyExists" });
    expect(cookieHelpers.createSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("answers a rejected password with a message", async () => {
    zitadel.addHumanUser.mockRejectedValue(apiError(Code.InvalidArgument, "Password is too short (COMMA-HuJf6)"));

    await expect(registerUser(registerCommand)).resolves.toEqual({ error: "errors.couldNotCreateUser" });
  });

  test("answers a rejected session with a message", async () => {
    cookieHelpers.createSessionAndUpdateCookie.mockRejectedValue(
      apiError(Code.FailedPrecondition, "User is locked (COMMAND-SF3fg)"),
    );

    await expect(registerUser(registerCommand)).resolves.toEqual({ error: "errors.couldNotCreateSession" });
  });

  test.each(serverFaults)("rethrows %s when creating the user", async (_, error) => {
    zitadel.addHumanUser.mockRejectedValue(error);

    await expect(registerUser(registerCommand)).rejects.toBe(error);
  });

  test.each(serverFaults)("rethrows %s when creating the session", async (_, error) => {
    cookieHelpers.createSessionAndUpdateCookie.mockRejectedValue(error);

    await expect(registerUser(registerCommand)).rejects.toBe(error);
  });
});
