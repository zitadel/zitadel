import { registerUserAndLinkToIDP } from "@/lib/actions/register";
import { createIdpFlowProof } from "@/lib/server/idp-flow-proof";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";

vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("@/lib/server/service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("@/lib/server/zitadel", () => ({
  addHumanUser: vi.fn(),
  addIDPLink: vi.fn(),
  getLoginSettings: vi.fn(),
  getUserByID: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
  // vg5j: the identity to link is read from the intent
  retrieveIDPIntent: vi.fn().mockResolvedValue({
    idpInformation: { idpId: "idp123", userId: "idp-user-123", userName: "testuser" },
  }),
}));

vi.mock("@/lib/server/session-cookie", () => ({
  createSessionAndUpdateCookie: vi.fn(),
  createSessionForIdpAndUpdateCookie: vi.fn(),
}));

vi.mock("@/lib/actions/flow", () => ({
  completeFlowOrGetUrl: vi.fn(),
}));

vi.mock("@/lib/server/verify-helper", () => ({
  checkEmailVerification: vi.fn(),
  checkMFAFactors: vi.fn(),
}));

vi.mock("@/lib/server/verification-check", () => ({
  setUserVerificationCheckCookie: vi.fn(),
}));

vi.mock("@/lib/server/fingerprint", () => ({
  getBindingIdCookie: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

describe("registerUserAndLinkToIDP", () => {
  const command = {
    email: "test@example.com",
    firstName: "Test",
    lastName: "User",
    organization: "org123",
    requestId: "req123",
    idpIntent: { idpIntentId: "intent123", idpIntentToken: "token123" },
    idpUserName: "testuser",
  };

  beforeEach(async () => {
    vi.clearAllMocks();
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "test-session-cookie-secret-at-least-32-chars");

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("@/lib/server/service-url");
    const { addHumanUser, addIDPLink, getLoginSettings, listAuthenticationMethodTypes } =
      await import("@/lib/server/zitadel");
    const { createSessionForIdpAndUpdateCookie } = await import("@/lib/server/session-cookie");
    const { completeFlowOrGetUrl } = await import("@/lib/actions/flow");
    const { getBindingIdCookie } = await import("@/lib/server/fingerprint");

    vi.mocked(headers).mockResolvedValue({} as any);
    vi.mocked(getServiceConfig).mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } } as any);
    vi.mocked(getLoginSettings).mockResolvedValue({ allowRegister: true } as any);
    vi.mocked(addHumanUser).mockResolvedValue({ userId: "user123" } as any);
    vi.mocked(addIDPLink).mockResolvedValue({} as any);
    vi.mocked(createSessionForIdpAndUpdateCookie).mockResolvedValue({
      id: "session123",
      factors: { user: { id: "user123", loginName: "testuser", organizationId: "org123" } },
    } as any);
    vi.mocked(listAuthenticationMethodTypes).mockResolvedValue({ authMethodTypes: [] } as any);
    vi.mocked(completeFlowOrGetUrl).mockResolvedValue({ redirect: "https://app.example.com/success" } as any);
    vi.mocked(getBindingIdCookie).mockResolvedValue({ name: "bindingId", value: "fp123" });
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllEnvs();
  });

  test.each([
    ["a missing proof", () => undefined],
    ["a proof for another browser", () => createIdpFlowProof("attacker-fp")],
    ["a proof of a linking flow", () => createIdpFlowProof("fp123", "session123")],
    ["a garbage value", () => "not-a-proof"],
  ])("should reject %s without creating a user", async (_name, makeFlowFingerprint) => {
    const { addHumanUser, addIDPLink } = await import("@/lib/server/zitadel");
    const { createSessionForIdpAndUpdateCookie } = await import("@/lib/server/session-cookie");

    const result = await registerUserAndLinkToIDP({ ...command, flowFingerprint: makeFlowFingerprint() });

    expect(result).toEqual({ error: "errors.sessionMismatch" });
    expect(addHumanUser).not.toHaveBeenCalled();
    expect(addIDPLink).not.toHaveBeenCalled();
    expect(createSessionForIdpAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("should reject a valid proof when the browser has no binding id cookie", async () => {
    const { addHumanUser } = await import("@/lib/server/zitadel");
    const { getBindingIdCookie } = await import("@/lib/server/fingerprint");
    vi.mocked(getBindingIdCookie).mockResolvedValue(undefined);

    const result = await registerUserAndLinkToIDP({ ...command, flowFingerprint: createIdpFlowProof("fp123") });

    expect(result).toEqual({ error: "errors.sessionMismatch" });
    expect(addHumanUser).not.toHaveBeenCalled();
  });

  test("should register and link the user when the proof matches this browser", async () => {
    const { addHumanUser } = await import("@/lib/server/zitadel");

    const result = await registerUserAndLinkToIDP({ ...command, flowFingerprint: createIdpFlowProof("fp123") });

    expect(addHumanUser).toHaveBeenCalled();
    expect(result).toEqual({ redirect: "https://app.example.com/success" });
  });
});
