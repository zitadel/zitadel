import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { shouldEnforceMFA } from "./verify-helper";
import { cookies } from "next/headers";
import { getFingerprintIdCookie } from "./fingerprint";
import crypto from "crypto";

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

import {
  checkPasswordChangeRequired,
  checkEmailVerified,
  checkEmailVerification,
  checkMFAFactors,
  checkUserVerification,
} from "./verify-helper";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";

// Mock dependencies
vi.mock("@zitadel/client", () => ({
  timestampDate: vi.fn((ts: any) => new Date(Number(ts.seconds) * 1000)),
  timestampFromMs: vi.fn((ms: number) => ({ seconds: BigInt(Math.floor(ms / 1000)) }) as any),
}));

vi.mock("next/headers", () => ({
  cookies: vi.fn(),
}));

vi.mock("./fingerprint", () => ({
  getFingerprintIdCookie: vi.fn(),
}));

vi.mock("./zitadel", () => ({
  getUserByID: vi.fn(),
}));

vi.mock("crypto", () => ({
  default: {
    createHash: vi.fn(() => ({
      update: vi.fn().mockReturnThis(),
      digest: vi.fn(() => "mockedhash123"),
    })),
  },
}));

describe("checkPasswordChangeRequired", () => {
  const mockSession: any = {
    id: "session-123",
    factors: {
      user: {
        id: "user-123",
        loginName: "user@example.com",
        organizationId: "org-123",
      },
    },
  };

  const createTimestamp = (daysAgo: number): any => {
    const date = new Date();
    date.setDate(date.getDate() - daysAgo);
    return {
      seconds: BigInt(Math.floor(date.getTime() / 1000)),
      nanos: 0,
    };
  };

  it("should redirect if password change is required on user", () => {
    const humanUser: any = {
      passwordChangeRequired: true,
      passwordChanged: createTimestamp(0),
    };

    const result = checkPasswordChangeRequired(undefined, mockSession, humanUser);

    expect(result).toEqual({
      redirect: expect.stringContaining("/password/change"),
    });
  });

  it("should redirect if password is expired based on maxAgeDays", () => {
    const expirySettings: any = {
      maxAgeDays: BigInt(30),
    };

    const humanUser: any = {
      passwordChangeRequired: false,
      passwordChanged: createTimestamp(35), // 35 days ago
    };

    const result = checkPasswordChangeRequired(expirySettings, mockSession, humanUser);

    expect(result).toEqual({
      redirect: expect.stringContaining("/password/change"),
    });
  });

  it("should not redirect if password is not expired", () => {
    const expirySettings: any = {
      maxAgeDays: BigInt(30),
    };

    const humanUser: any = {
      passwordChangeRequired: false,
      passwordChanged: createTimestamp(20), // 20 days ago
    };

    const result = checkPasswordChangeRequired(expirySettings, mockSession, humanUser);

    expect(result).toBeUndefined();
  });

  it("should include organization in redirect params", () => {
    const humanUser: any = {
      passwordChangeRequired: true,
    };

    const result = checkPasswordChangeRequired(undefined, mockSession, humanUser, "custom-org");

    expect(result?.redirect).toContain("organization=");
  });

  it("should include requestId in redirect params", () => {
    const humanUser: any = {
      passwordChangeRequired: true,
    };

    const result = checkPasswordChangeRequired(undefined, mockSession, humanUser, undefined, "request-123");

    expect(result?.redirect).toContain("requestId=request-123");
  });

  it("should handle no expiry settings", () => {
    const humanUser: any = {
      passwordChangeRequired: false,
      passwordChanged: createTimestamp(100),
    };

    const result = checkPasswordChangeRequired(undefined, mockSession, humanUser);

    expect(result).toBeUndefined();
  });

  it("should handle no password changed timestamp", () => {
    const expirySettings: any = {
      maxAgeDays: BigInt(30),
    };

    const humanUser: any = {
      passwordChangeRequired: false,
    };

    const result = checkPasswordChangeRequired(expirySettings, mockSession, humanUser);

    expect(result).toBeUndefined();
  });
});

describe("checkEmailVerified", () => {
  const mockSession: any = {
    factors: {
      user: {
        id: "user-123",
        loginName: "user@example.com",
        organizationId: "org-123",
      },
    },
  };

  it("should redirect if email is not verified", () => {
    const humanUser: any = {
      email: {
        isVerified: false,
      },
    };

    const result = checkEmailVerified(mockSession, humanUser);

    expect(result).toEqual({
      redirect: expect.stringContaining("/verify"),
    });
  });

  it("should not redirect if email is verified", () => {
    const humanUser: any = {
      email: {
        isVerified: true,
      },
    };

    const result = checkEmailVerified(mockSession, humanUser);

    expect(result).toBeUndefined();
  });

  it("should include userId in verify redirect", () => {
    const humanUser: any = {
      email: {
        isVerified: false,
      },
    };

    const result = checkEmailVerified(mockSession, humanUser);

    expect(result?.redirect).toContain("userId=user-123");
  });

  it("should include send=true parameter", () => {
    const humanUser: any = {
      email: {
        isVerified: false,
      },
    };

    const result = checkEmailVerified(mockSession, humanUser);

    expect(result?.redirect).toContain("send=true");
  });

  it("should include organization in redirect", () => {
    const humanUser: any = {
      email: {
        isVerified: false,
      },
    };

    const result = checkEmailVerified(mockSession, humanUser, "custom-org");

    expect(result?.redirect).toContain("organization=custom-org");
  });

  it("should include requestId in redirect", () => {
    const humanUser: any = {
      email: {
        isVerified: false,
      },
    };

    const result = checkEmailVerified(mockSession, humanUser, undefined, "request-123");

    expect(result?.redirect).toContain("requestId=request-123");
  });

  it("should handle no email on user", () => {
    const humanUser: any = {};

    const result = checkEmailVerified(mockSession, humanUser);

    expect(result).toEqual({
      redirect: expect.stringContaining("/verify"),
    });
  });
});

describe("checkEmailVerification", () => {
  const mockSession: any = {
    factors: {
      user: {
        loginName: "user@example.com",
        organizationId: "org-123",
      },
    },
  };

  const originalEnv = process.env.EMAIL_VERIFICATION;

  afterEach(() => {
    process.env.EMAIL_VERIFICATION = originalEnv;
  });

  it("should redirect if email not verified and EMAIL_VERIFICATION is true", () => {
    process.env.EMAIL_VERIFICATION = "true";

    const humanUser: any = {
      email: {
        isVerified: false,
      },
    };

    const result = checkEmailVerification(mockSession, humanUser);

    expect(result).toEqual({
      redirect: expect.stringContaining("/verify"),
    });
  });

  it("should not redirect if EMAIL_VERIFICATION is not true", () => {
    process.env.EMAIL_VERIFICATION = "false";

    const humanUser: any = {
      email: {
        isVerified: false,
      },
    };

    const result = checkEmailVerification(mockSession, humanUser);

    expect(result).toBeUndefined();
  });

  it("should not redirect if email is verified", () => {
    process.env.EMAIL_VERIFICATION = "true";

    const humanUser: any = {
      email: {
        isVerified: true,
      },
    };

    const result = checkEmailVerification(mockSession, humanUser);

    expect(result).toBeUndefined();
  });

  it("should include send=true parameter", () => {
    process.env.EMAIL_VERIFICATION = "true";

    const humanUser: any = {
      email: {
        isVerified: false,
      },
    };

    const result = checkEmailVerification(mockSession, humanUser);

    expect(result?.redirect).toContain("send=true");
  });

  it("should include organization in redirect", () => {
    process.env.EMAIL_VERIFICATION = "true";

    const humanUser: any = {
      email: {
        isVerified: false,
      },
    };

    const result = checkEmailVerification(mockSession, humanUser, "custom-org");

    expect(result?.redirect).toContain("organization=custom-org");
  });
});

describe("checkUserVerification", () => {
  let mockCookies: any;

  beforeEach(() => {
    vi.clearAllMocks();
    mockCookies = {
      get: vi.fn(),
    };
    vi.mocked(cookies).mockResolvedValue(mockCookies);
  });

  it("should return true if verification hash matches", async () => {
    vi.mocked(getFingerprintIdCookie).mockResolvedValue({
      name: "fingerprintId",
      value: "fingerprint-123",
    } as any);

    mockCookies.get.mockReturnValue({
      value: "mockedhash123",
    });

    const result = await checkUserVerification("user-123");

    expect(result).toBe(true);
  });

  it("should return false if fingerprint cookie not found", async () => {
    vi.mocked(getFingerprintIdCookie).mockResolvedValue(undefined);

    const result = await checkUserVerification("user-123");

    expect(result).toBe(false);
  });

  it("should return false if fingerprint cookie has no value", async () => {
    vi.mocked(getFingerprintIdCookie).mockResolvedValue({
      name: "fingerprintId",
      value: "",
    } as any);

    const result = await checkUserVerification("user-123");

    expect(result).toBe(false);
  });

  it("should return false if verification cookie not found", async () => {
    vi.mocked(getFingerprintIdCookie).mockResolvedValue({
      name: "fingerprintId",
      value: "fingerprint-123",
    } as any);

    mockCookies.get.mockReturnValue(undefined);

    const result = await checkUserVerification("user-123");

    expect(result).toBe(false);
  });

  it("should return false if verification hash does not match", async () => {
    vi.mocked(getFingerprintIdCookie).mockResolvedValue({
      name: "fingerprintId",
      value: "fingerprint-123",
    } as any);

    mockCookies.get.mockReturnValue({
      value: "wronghash",
    });

    const result = await checkUserVerification("user-123");

    expect(result).toBe(false);
  });

  it("should create hash from userId and fingerprint", async () => {
    vi.mocked(getFingerprintIdCookie).mockResolvedValue({
      name: "fingerprintId",
      value: "fingerprint-456",
    } as any);

    mockCookies.get.mockReturnValue({
      value: "mockedhash123",
    });

    await checkUserVerification("user-456");

    expect(crypto.createHash).toHaveBeenCalledWith("sha256");
  });
});

describe("checkMFAFactors", () => {
  const mockSession: any = {
    id: "session-123",
    factors: {
      user: {
        id: "user-123",
        loginName: "user@example.com",
        organizationId: "org-123",
      },
    },
  };

  const mockLoginSettings: any = {
    forceMfa: false,
    forceMfaLocalOnly: false,
  };

  beforeEach(() => {
    vi.clearAllMocks();
    vi.spyOn(console, "log").mockImplementation(() => {});
  });

  it("should not require MFA if authenticated with passkey", async () => {
    const sessionWithPasskey = {
      ...mockSession,
      factors: {
        ...mockSession.factors,
        webAuthN: {
          verifiedAt: createMockTimestamp(),
          userVerified: true,
        },
      },
    };

    const result = await checkMFAFactors("https://example.com", sessionWithPasskey, mockLoginSettings, []);

    expect(result).toBeUndefined();
  });

  it("should redirect to TOTP if only TOTP is available", async () => {
    const authMethods = [AuthenticationMethodType.TOTP];

    const result = await checkMFAFactors("https://example.com", mockSession, mockLoginSettings, authMethods);

    expect(result).toEqual({
      redirect: expect.stringContaining("/otp/time-based"),
    });
  });

  it("should redirect to SMS OTP if only OTP_SMS is available", async () => {
    const authMethods = [AuthenticationMethodType.OTP_SMS];

    const result = await checkMFAFactors("https://example.com", mockSession, mockLoginSettings, authMethods);

    expect(result).toEqual({
      redirect: expect.stringContaining("/otp/sms"),
    });
  });

  it("should redirect to Email OTP if only OTP_EMAIL is available", async () => {
    const authMethods = [AuthenticationMethodType.OTP_EMAIL];

    const result = await checkMFAFactors("https://example.com", mockSession, mockLoginSettings, authMethods);

    expect(result).toEqual({
      redirect: expect.stringContaining("/otp/email"),
    });
  });

  it("should redirect to U2F if only U2F is available", async () => {
    const authMethods = [AuthenticationMethodType.U2F];

    const result = await checkMFAFactors("https://example.com", mockSession, mockLoginSettings, authMethods);

    expect(result).toEqual({
      redirect: expect.stringContaining("/u2f"),
    });
  });

  it("should redirect to MFA selection page if multiple factors available", async () => {
    const authMethods = [AuthenticationMethodType.TOTP, AuthenticationMethodType.OTP_SMS];

    const result = await checkMFAFactors("https://example.com", mockSession, mockLoginSettings, authMethods);

    expect(result).toEqual({
      redirect: expect.stringContaining("/mfa?"),
    });
  });

  it("should include organization in redirect params", async () => {
    const authMethods = [AuthenticationMethodType.TOTP];

    const result = await checkMFAFactors("https://example.com", mockSession, mockLoginSettings, authMethods, "custom-org");

    expect(result?.redirect).toContain("organization=custom-org");
  });

  it("should include requestId in redirect params", async () => {
    const authMethods = [AuthenticationMethodType.TOTP];

    const result = await checkMFAFactors(
      "https://example.com",
      mockSession,
      mockLoginSettings,
      authMethods,
      undefined,
      "request-123",
    );

    expect(result?.redirect).toContain("requestId=request-123");
  });

  it("should ignore non-MFA authentication methods", async () => {
    const authMethods = [AuthenticationMethodType.PASSWORD, AuthenticationMethodType.PASSKEY];

    const result = await checkMFAFactors("https://example.com", mockSession, mockLoginSettings, authMethods);

    expect(result).toBeUndefined();
  });
});
