import { Code, ConnectError } from "@connectrpc/connect";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { ClassifiedConnectError } from "../grpc/interceptors/error-classification";
import { registerUserAndLinkToIDP } from "./register";

vi.mock("next/headers", () => ({
  headers: vi.fn(),
  cookies: vi.fn(),
}));

// this returns the key itself so tests can assert on translation keys
vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("../zitadel", () => ({
  addHumanUser: vi.fn(),
  addIDPLink: vi.fn(),
  getLoginSettings: vi.fn(),
  getUserByID: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
}));

vi.mock("@/lib/server/cookie", () => ({
  createSessionAndUpdateCookie: vi.fn(),
  createSessionForIdpAndUpdateCookie: vi.fn(),
}));

vi.mock("../client", () => ({
  completeFlowOrGetUrl: vi.fn(),
}));

vi.mock("../fingerprint", () => ({
  getOrSetFingerprintId: vi.fn(),
}));

vi.mock("../verify-helper", () => ({
  checkEmailVerification: vi.fn(),
  checkMFAFactors: vi.fn(),
}));

function classifiedError(code: Code, message: string): ClassifiedConnectError {
  return new ClassifiedConnectError(new ConnectError(message, code));
}

describe("registerUserAndLinkToIDP", () => {
  let mockHeaders: any;
  let mockGetServiceConfig: any;
  let mockGetLoginSettings: any;
  let mockAddHumanUser: any;
  let mockAddIDPLink: any;
  let mockCreateSessionForIdpAndUpdateCookie: any;

  const command = {
    email: "user@example.com",
    firstName: "Jane",
    lastName: "Doe",
    organization: "org-1",
    idpIntent: { idpIntentId: "intent-1", idpIntentToken: "intent-token" },
    idpUserId: "idp-user-1",
    idpId: "idp-1",
    idpUserName: "user@example.com",
  };

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { addHumanUser, addIDPLink, getLoginSettings } = await import("../zitadel");
    const { createSessionForIdpAndUpdateCookie } = await import("@/lib/server/cookie");

    mockHeaders = vi.mocked(headers);
    mockGetServiceConfig = vi.mocked(getServiceConfig);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockAddHumanUser = vi.mocked(addHumanUser);
    mockAddIDPLink = vi.mocked(addIDPLink);
    mockCreateSessionForIdpAndUpdateCookie = vi.mocked(createSessionForIdpAndUpdateCookie);

    mockHeaders.mockResolvedValue({} as any);
    mockGetServiceConfig.mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } });
    mockGetLoginSettings.mockResolvedValue({ allowRegister: true });
  });

  test("returns userAlreadyExists error when the user already exists instead of throwing", async () => {
    mockAddHumanUser.mockRejectedValue(classifiedError(Code.AlreadyExists, "User already exists (V3-DKcYh)"));

    const response = await registerUserAndLinkToIDP(command);

    expect(response).toEqual({ error: "errors.userAlreadyExists" });
    expect(mockAddIDPLink).not.toHaveBeenCalled();
    expect(mockCreateSessionForIdpAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("returns couldNotCreateUser for other user errors during creation", async () => {
    mockAddHumanUser.mockRejectedValue(classifiedError(Code.InvalidArgument, "Errors.User.EmailInvalid"));

    const response = await registerUserAndLinkToIDP(command);

    expect(response).toEqual({ error: "errors.couldNotCreateUser" });
  });

  test("rethrows genuine server errors so they still surface as 500", async () => {
    mockAddHumanUser.mockRejectedValue(classifiedError(Code.Internal, "database down"));

    await expect(registerUserAndLinkToIDP(command)).rejects.toThrow("database down");
  });

  test("returns couldNotLinkIDP when linking the IDP fails with a user error", async () => {
    mockAddHumanUser.mockResolvedValue({ userId: "user-1" });
    mockAddIDPLink.mockRejectedValue(classifiedError(Code.AlreadyExists, "Errors.User.ExternalIDP.AlreadyExists"));

    const response = await registerUserAndLinkToIDP(command);

    expect(response).toEqual({ error: "errors.couldNotLinkIDP" });
    expect(mockCreateSessionForIdpAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("returns couldNotCreateSession when session creation fails with a user error", async () => {
    mockAddHumanUser.mockResolvedValue({ userId: "user-1" });
    mockAddIDPLink.mockResolvedValue({ details: {} });
    mockCreateSessionForIdpAndUpdateCookie.mockRejectedValue(
      classifiedError(Code.FailedPrecondition, "Errors.Intent.Expired"),
    );

    const response = await registerUserAndLinkToIDP(command);

    expect(response).toEqual({ error: "errors.couldNotCreateSession" });
  });

  test("returns registerNotAllowed when registration is disabled", async () => {
    mockGetLoginSettings.mockResolvedValue({ allowRegister: false });

    const response = await registerUserAndLinkToIDP(command);

    expect(response).toEqual({ error: "errors.registerNotAllowed" });
    expect(mockAddHumanUser).not.toHaveBeenCalled();
  });
});
