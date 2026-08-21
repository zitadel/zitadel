import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { LoginSettings, PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { describe, expect, it } from "vitest";
import { checkSessionReuse } from "./session-reuse";

const verifiedAt = { seconds: BigInt(1), nanos: 0 } as any;

function makeSession(overrides: {
  organizationId?: string;
  password?: boolean;
  intent?: boolean;
  passkey?: boolean;
}): Session {
  return {
    factors: {
      user: {
        loginName: "test@example.com",
        organizationId: overrides.organizationId ?? "org-a",
      },
      ...(overrides.password ? { password: { verifiedAt } } : {}),
      ...(overrides.intent ? { intent: { verifiedAt } } : {}),
      ...(overrides.passkey ? { webAuthN: { verifiedAt } } : {}),
    },
  } as unknown as Session;
}

function makeSettings(overrides: Partial<LoginSettings>): LoginSettings {
  return {
    allowLocalAuthentication: true,
    allowExternalIdp: true,
    passkeysType: PasskeysType.ALLOWED,
    ...overrides,
  } as LoginSettings;
}

describe("checkSessionReuse", () => {
  describe("organization gate", () => {
    it("blocks a session whose user belongs to a different organization", () => {
      const result = checkSessionReuse({
        session: makeSession({ organizationId: "org-a", password: true }),
        targetOrganization: "org-b",
      });
      expect(result).toEqual({ reusable: false, reason: "orgMismatch" });
    });

    it("allows a session in the target organization", () => {
      const result = checkSessionReuse({
        session: makeSession({ organizationId: "org-b", password: true }),
        targetOrganization: "org-b",
      });
      expect(result).toEqual({ reusable: true });
    });

    it("does not enforce the org gate when no target organization is known", () => {
      const result = checkSessionReuse({
        session: makeSession({ organizationId: "org-a", password: true }),
      });
      expect(result).toEqual({ reusable: true });
    });
  });

  describe("auth-method gate", () => {
    it("treats a session as reusable when login settings are unknown", () => {
      const result = checkSessionReuse({
        session: makeSession({ password: true }),
      });
      expect(result).toEqual({ reusable: true });
    });

    it("blocks a password session when local authentication is not allowed", () => {
      const result = checkSessionReuse({
        session: makeSession({ password: true }),
        loginSettings: makeSettings({ allowLocalAuthentication: false }),
      });
      expect(result).toEqual({ reusable: false, reason: "localAuthNotAllowed" });
    });

    it("blocks an external-IdP session when external IdPs are not allowed", () => {
      const result = checkSessionReuse({
        session: makeSession({ intent: true }),
        loginSettings: makeSettings({ allowExternalIdp: false }),
      });
      expect(result).toEqual({ reusable: false, reason: "externalIdpNotAllowed" });
    });

    it("blocks a passkey session when passkeys are not allowed", () => {
      const result = checkSessionReuse({
        session: makeSession({ passkey: true }),
        loginSettings: makeSettings({ passkeysType: PasskeysType.NOT_ALLOWED }),
      });
      expect(result).toEqual({ reusable: false, reason: "passkeysNotAllowed" });
    });

    it("treats an IdP session as an IdP login even when a password factor is also present", () => {
      const result = checkSessionReuse({
        session: makeSession({ intent: true, password: true }),
        loginSettings: makeSettings({ allowLocalAuthentication: false, allowExternalIdp: true }),
      });
      // primary factor is IdP, which is allowed -> reusable despite local auth being disabled
      expect(result).toEqual({ reusable: true });
    });

    it("allows a password session when local authentication is allowed", () => {
      const result = checkSessionReuse({
        session: makeSession({ password: true }),
        loginSettings: makeSettings({ allowLocalAuthentication: true }),
      });
      expect(result).toEqual({ reusable: true });
    });
  });

  it("applies the org gate before the auth-method gate", () => {
    const result = checkSessionReuse({
      session: makeSession({ organizationId: "org-a", intent: true }),
      targetOrganization: "org-b",
      loginSettings: makeSettings({ allowExternalIdp: true }),
    });
    expect(result).toEqual({ reusable: false, reason: "orgMismatch" });
  });
});
