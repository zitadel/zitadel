import { describe, it, expect, beforeEach } from "vitest";
import { shouldEnforceMFA } from "./verify-helper";

// Mock function to create timestamps - following the same pattern as session.test.ts
function createMockTimestamp(offsetMs = 3600000): any {
  return {
    seconds: BigInt(Math.floor((Date.now() + offsetMs) / 1000)),
    nanos: 0,
  };
}

// Mock function to create a basic session - following the same pattern as session.test.ts
function createMockSession(overrides: any = {}): any {
  const futureTimestamp = createMockTimestamp();

  const defaultSession = {
    id: "test-session-id",
    factors: {
      user: {
        id: "test-user-id",
        loginName: "test@example.com",
        displayName: "Test User",
        organizationId: "test-org-id",
        verifiedAt: futureTimestamp,
      },
    },
    ...overrides,
  };

  return defaultSession;
}

// Mock function to create login settings
function createMockLoginSettings(overrides: any = {}): any {
  return {
    forceMfa: false,
    forceMfaLocalOnly: false,
    ...overrides,
  };
}

describe("shouldEnforceMFA", () => {
  let mockSession: any;
  let mockLoginSettings: any;

  beforeEach(() => {
    mockSession = createMockSession();
    mockLoginSettings = createMockLoginSettings();
  });

  describe("when loginSettings is undefined", () => {
    it("should return false", () => {
      const result = shouldEnforceMFA(mockSession, undefined);
      expect(result).toBe(false);
    });
  });

  describe("passkey authentication", () => {
    beforeEach(() => {
      mockSession = createMockSession({
        factors: {
          ...mockSession.factors,
          webAuthN: {
            verifiedAt: createMockTimestamp(),
            userVerified: true,
          },
        },
      });
    });

    it("should return false when user authenticated with passkey, even with forceMfa enabled", () => {
      mockLoginSettings = createMockLoginSettings({ forceMfa: true });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(false);
    });

    it("should return false when user authenticated with passkey, even with forceMfaLocalOnly enabled", () => {
      mockLoginSettings = createMockLoginSettings({ forceMfaLocalOnly: true });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(false);
    });

    it("should return false when user authenticated with passkey and both force settings enabled", () => {
      mockLoginSettings = createMockLoginSettings({
        forceMfa: true,
        forceMfaLocalOnly: true,
      });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(false);
    });

    it("should return true when passkey is not user verified", () => {
      mockSession = createMockSession({
        factors: {
          ...mockSession.factors,
          webAuthN: {
            verifiedAt: createMockTimestamp(),
            userVerified: false, // Not user verified
          },
        },
      });
      mockLoginSettings = createMockLoginSettings({ forceMfa: true });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      // Should return true because passkey is not user verified, so it doesn't count as passkey auth
      expect(result).toBe(true);
    });
  });

  describe("forceMfa setting", () => {
    beforeEach(() => {
      mockLoginSettings = createMockLoginSettings({ forceMfa: true });
    });

    it("should return true when forceMfa is enabled and user authenticated with password", () => {
      mockSession = createMockSession({
        factors: {
          ...mockSession.factors,
          password: {
            verifiedAt: createMockTimestamp(),
          },
        },
      });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(true);
    });

    it("should return true when forceMfa is enabled and user authenticated with IDP", () => {
      mockSession = createMockSession({
        factors: {
          ...mockSession.factors,
          intent: {
            verifiedAt: createMockTimestamp(),
          },
        },
      });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(true);
    });

    it("should return true when forceMfa is enabled with no specific authentication method", () => {
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(true);
    });
  });

  describe("forceMfaLocalOnly setting", () => {
    beforeEach(() => {
      mockLoginSettings = createMockLoginSettings({ forceMfaLocalOnly: true });
    });

    it("should return true when forceMfaLocalOnly is enabled and user authenticated with password", () => {
      mockSession = createMockSession({
        factors: {
          ...mockSession.factors,
          password: {
            verifiedAt: createMockTimestamp(),
          },
        },
      });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(true);
    });

    it("should return false when forceMfaLocalOnly is enabled and user authenticated with IDP", () => {
      mockSession = createMockSession({
        factors: {
          ...mockSession.factors,
          intent: {
            verifiedAt: createMockTimestamp(),
          },
        },
      });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(false);
    });

    it("should return false when forceMfaLocalOnly is enabled with no specific authentication method", () => {
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(false);
    });
  });

  describe("mixed authentication scenarios", () => {
    it("should prioritize passkey over password when both are present", () => {
      mockSession = createMockSession({
        factors: {
          ...mockSession.factors,
          password: {
            verifiedAt: createMockTimestamp(),
          },
          webAuthN: {
            verifiedAt: createMockTimestamp(),
            userVerified: true,
          },
        },
      });
      mockLoginSettings = createMockLoginSettings({ forceMfa: true });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(false); // Passkey should override password
    });

    it("should prioritize passkey over IDP when both are present", () => {
      mockSession = createMockSession({
        factors: {
          ...mockSession.factors,
          intent: {
            verifiedAt: createMockTimestamp(),
          },
          webAuthN: {
            verifiedAt: createMockTimestamp(),
            userVerified: true,
          },
        },
      });
      mockLoginSettings = createMockLoginSettings({ forceMfaLocalOnly: true });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(false); // Passkey should override IDP
    });

    it("should handle password + IDP scenario with forceMfaLocalOnly", () => {
      mockSession = createMockSession({
        factors: {
          ...mockSession.factors,
          password: {
            verifiedAt: createMockTimestamp(),
          },
          intent: {
            verifiedAt: createMockTimestamp(),
          },
        },
      });
      mockLoginSettings = createMockLoginSettings({ forceMfaLocalOnly: true });
      // With both password and IDP, the current logic should return false for IDP
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(false);
    });
  });

  describe("no MFA enforcement", () => {
    it("should return false when neither forceMfa nor forceMfaLocalOnly is enabled", () => {
      mockLoginSettings = createMockLoginSettings({
        forceMfa: false,
        forceMfaLocalOnly: false,
      });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(false);
    });
  });

  describe("edge cases", () => {
    it("should handle session with no factors", () => {
      mockSession = createMockSession({
        factors: undefined,
      });
      mockLoginSettings = createMockLoginSettings({ forceMfa: true });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(true);
    });

    it("should handle session with empty factors", () => {
      mockSession = createMockSession({
        factors: {
          user: {
            id: "test-user-id",
            loginName: "test@example.com",
            displayName: "Test User",
            organizationId: "test-org-id",
            verifiedAt: createMockTimestamp(),
          },
        },
      });
      mockLoginSettings = createMockLoginSettings({ forceMfa: true });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(true);
    });

    it("should handle webAuthN factor without userVerified", () => {
      mockSession = createMockSession({
        factors: {
          ...mockSession.factors,
          webAuthN: {
            verifiedAt: createMockTimestamp(),
            userVerified: false,
          },
        },
      });
      mockLoginSettings = createMockLoginSettings({ forceMfa: true });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(true); // Should require MFA since it's not a proper passkey
    });

    it("should handle webAuthN factor without verifiedAt", () => {
      mockSession = createMockSession({
        factors: {
          ...mockSession.factors,
          webAuthN: {
            userVerified: true,
            // verifiedAt is undefined
          },
        },
      });
      mockLoginSettings = createMockLoginSettings({ forceMfa: true });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(true); // Should require MFA since webAuthN wasn't actually verified
    });

    it("should handle webAuthN factor with verifiedAt but no userVerified property", () => {
      mockSession = createMockSession({
        factors: {
          ...mockSession.factors,
          webAuthN: {
            verifiedAt: createMockTimestamp(),
            // userVerified is undefined (should be falsy)
          },
        },
      });
      mockLoginSettings = createMockLoginSettings({ forceMfa: true });
      const result = shouldEnforceMFA(mockSession, mockLoginSettings);
      expect(result).toBe(true); // Should require MFA since userVerified is falsy
    });
  });
});
