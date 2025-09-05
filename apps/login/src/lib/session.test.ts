/**
 * Unit tests for the isSessionValid function.
 * 
 * This test suite covers the comprehensive session validation logic including:
 * - Session expiration checks
 * - User presence validation
 * - Authentication factor verification (password, passkey, IDP)
 * - MFA validation with configured authentication methods (TOTP, OTP Email/SMS, U2F)
 * - MFA validation with login settings (forceMfa, forceMfaLocalOnly)
 * - Email verification when EMAIL_VERIFICATION environment variable is enabled
 * - Edge cases like sessions without expiration date
 */

import { timestampDate } from "@zitadel/client";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { isSessionValid } from "./session";
import * as zitadelModule from "./zitadel";

// Mock the zitadel client timestampDate function
vi.mock("@zitadel/client", () => ({
  timestampDate: vi.fn(),
}));

// Mock the zitadel module
vi.mock("./zitadel", () => ({
  listAuthenticationMethodTypes: vi.fn(),
  getLoginSettings: vi.fn(),
  getUserByID: vi.fn(),
}));

// Mock environment variables
const originalEnv = process.env;

describe("isSessionValid", () => {
  const mockServiceUrl = "https://zitadel-abc123.zitadel.cloud";
  const mockUserId = "test-user-id";
  const mockOrganizationId = "test-org-id";

  beforeEach(() => {
    vi.clearAllMocks();
    process.env = { ...originalEnv };
    // @ts-ignore - delete is OK for test environment variables
    delete process.env.EMAIL_VERIFICATION;

    // Setup timestampDate mock to return valid dates
    vi.mocked(timestampDate).mockImplementation((timestamp: any) => {
      if (!timestamp || !timestamp.seconds) {
        return new Date(); // Return current date for invalid timestamps
      }
      return new Date(Number(timestamp.seconds) * 1000);
    });
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  const createMockTimestamp = (offsetMs = 3600000): any => ({
    seconds: BigInt(Math.floor((Date.now() + offsetMs) / 1000)),
  });

  const createMockSession = (overrides: any = {}): any => {
    const futureTimestamp = createMockTimestamp();

    const defaultSession = {
      id: "session-id",
      expirationDate: futureTimestamp,
      factors: {
        user: {
          id: mockUserId,
          organizationId: mockOrganizationId,
          loginName: "test@example.com",
          displayName: "Test User",
          verifiedAt: futureTimestamp,
        },
        password: {
          verifiedAt: futureTimestamp,
        },
      },
      ...overrides,
    };

    return defaultSession;
  };

  describe("when session has no user", () => {
    test("should return false and log warning", async () => {
      const consoleSpy = vi.spyOn(console, "warn").mockImplementation(() => {});
      const session = createMockSession({
        factors: undefined,
      });

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(false);
      expect(consoleSpy).toHaveBeenCalledWith("Session has no user");
      consoleSpy.mockRestore();
    });
  });

  describe("when session is expired", () => {
    test("should return false and log warning", async () => {
      const consoleSpy = vi.spyOn(console, "warn").mockImplementation(() => {});
      const session = createMockSession({
        expirationDate: createMockTimestamp(-3600000), // 1 hour ago
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [],
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(false);
      expect(consoleSpy).toHaveBeenCalledWith(
        "Session is expired",
        expect.any(String)
      );
      consoleSpy.mockRestore();
    });
  });

  describe("when session has no valid authentication factors", () => {
    test("should return false when no password, passkey, or IDP verification", async () => {
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: createMockTimestamp(),
          },
          // No password, webAuthN, or intent factors
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [],
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(false);
    });
  });

  describe("MFA validation with configured authentication methods", () => {
    test("should return true when TOTP is configured and verified", async () => {
      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          password: {
            verifiedAt: verifiedTimestamp,
          },
          totp: {
            verifiedAt: verifiedTimestamp,
          },
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.TOTP],
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(true);
    });

    test("should return false when TOTP is configured but not verified", async () => {
      const consoleSpy = vi.spyOn(console, "warn").mockImplementation(() => {});
      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          password: {
            verifiedAt: verifiedTimestamp,
          },
          // No TOTP verification
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.TOTP],
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(false);
      expect(consoleSpy).toHaveBeenCalledWith(
        "Session has no valid MFA factor. Configured methods:",
        [AuthenticationMethodType.TOTP],
        "Session factors:",
        expect.objectContaining({
          totp: undefined,
          otpEmail: undefined,
          otpSms: undefined,
          webAuthN: undefined,
        })
      );
      consoleSpy.mockRestore();
    });

    test("should return true when OTP Email is configured and verified", async () => {
      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          password: {
            verifiedAt: verifiedTimestamp,
          },
          otpEmail: {
            verifiedAt: verifiedTimestamp,
          },
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.OTP_EMAIL],
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(true);
    });

    test("should return true when U2F is configured and verified", async () => {
      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          password: {
            verifiedAt: verifiedTimestamp,
          },
          webAuthN: {
            verifiedAt: verifiedTimestamp,
          },
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.U2F],
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(true);
    });

    test("should return true when multiple auth methods are configured and one is verified", async () => {
      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          password: {
            verifiedAt: verifiedTimestamp,
          },
          otpEmail: {
            verifiedAt: verifiedTimestamp,
          },
          // TOTP not verified
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.TOTP, AuthenticationMethodType.OTP_EMAIL],
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(true);
    });
  });

  describe("MFA validation with login settings (no configured auth methods)", () => {
    test("should return true when MFA is not forced and no auth methods configured", async () => {
      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          password: {
            verifiedAt: verifiedTimestamp,
          },
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [],
      } as any);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        forceMfa: false,
        forceMfaLocalOnly: false,
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(true);
    });

    test("should return false when MFA is forced but no factors are verified", async () => {
      const consoleSpy = vi.spyOn(console, "warn").mockImplementation(() => {});
      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          password: {
            verifiedAt: verifiedTimestamp,
          },
          // No MFA factors verified
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [],
      } as any);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        forceMfa: true,
        forceMfaLocalOnly: false,
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(false);
      expect(consoleSpy).toHaveBeenCalledWith(
        "Session has no valid multifactor",
        expect.any(Object)
      );
      consoleSpy.mockRestore();
    });

    test("should return true when MFA is forced and TOTP is verified", async () => {
      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          password: {
            verifiedAt: verifiedTimestamp,
          },
          totp: {
            verifiedAt: verifiedTimestamp,
          },
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [],
      } as any);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        forceMfa: true,
        forceMfaLocalOnly: false,
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(true);
    });

    test("should return true when forceMfaLocalOnly is enabled and WebAuthn is verified", async () => {
      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          password: {
            verifiedAt: verifiedTimestamp,
          },
          webAuthN: {
            verifiedAt: verifiedTimestamp,
          },
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [],
      } as any);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        forceMfa: false,
        forceMfaLocalOnly: true,
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(true);
    });
  });

  describe("email verification", () => {
    test("should return false when EMAIL_VERIFICATION is enabled and email is not verified", async () => {
      process.env.EMAIL_VERIFICATION = "true";
      const consoleSpy = vi.spyOn(console, "warn").mockImplementation(() => {});

      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          password: {
            verifiedAt: verifiedTimestamp,
          },
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [],
      } as any);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        forceMfa: false,
        forceMfaLocalOnly: false,
      } as any);

      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: {
          type: {
            case: "human",
            value: {
              email: {
                email: "test@example.com",
                isVerified: false,
              },
            },
          },
        },
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(false);
      expect(consoleSpy).toHaveBeenCalledWith(
        "Session invalid: Email not verified and EMAIL_VERIFICATION is enabled",
        mockUserId
      );
      consoleSpy.mockRestore();
    });

    test("should return true when EMAIL_VERIFICATION is enabled and email is verified", async () => {
      process.env.EMAIL_VERIFICATION = "true";

      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          password: {
            verifiedAt: verifiedTimestamp,
          },
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [],
      } as any);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        forceMfa: false,
        forceMfaLocalOnly: false,
      } as any);

      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: {
          type: {
            case: "human",
            value: {
              email: {
                email: "test@example.com",
                isVerified: true,
              },
            },
          },
        },
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(true);
    });

    test("should return true when EMAIL_VERIFICATION is disabled", async () => {
      // EMAIL_VERIFICATION is not set, so it's disabled by default

      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          password: {
            verifiedAt: verifiedTimestamp,
          },
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [],
      } as any);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        forceMfa: false,
        forceMfaLocalOnly: false,
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(true);
      // getUserByID should not be called when EMAIL_VERIFICATION is disabled
      expect(zitadelModule.getUserByID).not.toHaveBeenCalled();
    });
  });

  describe("passkey authentication", () => {
    test("should return true when authenticated with passkey (WebAuthn)", async () => {
      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          webAuthN: {
            verifiedAt: verifiedTimestamp,
          },
          // No password factor
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [],
      } as any);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        forceMfa: false,
        forceMfaLocalOnly: false,
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(true);
    });
  });

  describe("IDP authentication", () => {
    test("should return true when authenticated with IDP intent", async () => {
      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          intent: {
            verifiedAt: verifiedTimestamp,
          },
          // No password factor
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [],
      } as any);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        forceMfa: false,
        forceMfaLocalOnly: false,
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(true);
    });
  });

  describe("edge cases", () => {
    test("should handle session without expiration date", async () => {
      const verifiedTimestamp = createMockTimestamp();
      const session = createMockSession({
        expirationDate: undefined, // No expiration date
        factors: {
          user: {
            id: mockUserId,
            organizationId: mockOrganizationId,
            loginName: "test@example.com",
            displayName: "Test User",
            verifiedAt: verifiedTimestamp,
          },
          password: {
            verifiedAt: verifiedTimestamp,
          },
        },
      });

      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [],
      } as any);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        forceMfa: false,
        forceMfaLocalOnly: false,
      } as any);

      const result = await isSessionValid({ serviceUrl: mockServiceUrl, session });

      expect(result).toBe(true);
    });
  });
});