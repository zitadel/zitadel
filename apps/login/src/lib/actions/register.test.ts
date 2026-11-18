import { registerUserAndLinkToIDP } from "@/lib/actions/register";
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
