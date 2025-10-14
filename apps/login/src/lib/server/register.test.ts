/**
 * Unit tests for the registerUser server action.
 *
 * These tests replace the integration tests from register.cy.ts which tested:
 * - User registration with passwordless (passkey) setup
 * - User registration flow and session creation
 * - Redirect to passkey setup after successful registration
 */

import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { registerUser } from "./register";
import * as zitadelModule from "../zitadel";
import * as cookieModule from "./cookie";
import { PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";

// Mock all dependencies
vi.mock("../zitadel");
vi.mock("./cookie");
vi.mock("../client");
vi.mock("../fingerprint");
vi.mock("next/headers", () => ({
  headers: vi.fn(() => Promise.resolve(new Map())),
  cookies: vi.fn(() =>
    Promise.resolve({
      set: vi.fn(),
    }),
  ),
}));
vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => Promise.resolve((key: string) => key)),
}));
vi.mock("../service-url", () => ({
  getServiceUrlFromHeaders: vi.fn(() => ({ serviceUrl: "https://zitadel-test.zitadel.cloud" })),
}));

describe("registerUser server action", () => {
  const mockServiceUrl = "https://zitadel-test.zitadel.cloud";
  const mockUserId = "221394658884845598";
  const mockEmail = "john@example.com";
  const mockOrganization = "256088834543534543";
  const mockRequestId = "req123";

  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("passwordless registration (passkey setup)", () => {
    it("should redirect to /passkey/set when user registers without password", async () => {
      // Mock add user response
      vi.mocked(zitadelModule.addHumanUser).mockResolvedValue({
        userId: mockUserId,
      } as any);

      // Mock login settings
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        passkeysType: PasskeysType.ALLOWED,
        allowRegister: true,
      } as any);

      // Mock session creation
      vi.mocked(cookieModule.createSessionAndUpdateCookie).mockResolvedValue({
        id: "session-id",
        factors: {
          user: {
            id: mockUserId,
            loginName: mockEmail,
            organizationId: mockOrganization,
          },
        },
      } as any);

      const result = await registerUser({
        email: mockEmail,
        firstName: "John",
        lastName: "Doe",
        organization: mockOrganization,
        requestId: mockRequestId,
      });

      expect(result).toHaveProperty("redirect");
      expect(result.redirect).toContain("/passkey/set");
      expect(result.redirect).toContain(`loginName=${encodeURIComponent(mockEmail)}`);
      expect(result.redirect).toContain(`organization=${mockOrganization}`);
      expect(result.redirect).toContain(`requestId=${mockRequestId}`);
    });
  });

  describe("password registration", () => {
    it("should create user with password and check email verification", async () => {
      const mockPassword = "SecurePassword123!";

      // Mock add user response
      vi.mocked(zitadelModule.addHumanUser).mockResolvedValue({
        userId: mockUserId,
      } as any);

      // Mock login settings
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        allowRegister: true,
        allowUsernamePassword: true,
        passwordCheckLifetime: { seconds: BigInt(300) },
      } as any);

      // Mock session creation
      vi.mocked(cookieModule.createSessionAndUpdateCookie).mockResolvedValue({
        id: "session-id",
        factors: {
          user: {
            id: mockUserId,
            loginName: mockEmail,
            organizationId: mockOrganization,
          },
        },
      } as any);

      // Mock user retrieval
      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: {
          userId: mockUserId,
          type: {
            case: "human",
            value: {
              email: {
                email: mockEmail,
                isVerified: true,
              },
            },
          },
        },
      } as any);

      const { completeFlowOrGetUrl } = await import("../client");
      vi.mocked(completeFlowOrGetUrl).mockResolvedValue({
        redirect: "/dashboard",
      } as any);

      const result = await registerUser({
        email: mockEmail,
        firstName: "John",
        lastName: "Doe",
        password: mockPassword,
        organization: mockOrganization,
        requestId: mockRequestId,
      });

      expect(vi.mocked(zitadelModule.addHumanUser)).toHaveBeenCalledWith(
        expect.objectContaining({
          email: mockEmail,
          password: mockPassword,
          organization: mockOrganization,
        }),
      );

      expect(result).toHaveProperty("redirect");
    });
  });

  describe("error handling", () => {
    it("should return error when user creation fails", async () => {
      vi.mocked(zitadelModule.addHumanUser).mockResolvedValue(undefined as any);

      const result = await registerUser({
        email: mockEmail,
        firstName: "John",
        lastName: "Doe",
        organization: mockOrganization,
      });

      expect(result).toHaveProperty("error");
      expect(result.error).toBe("errors.couldNotCreateUser");
    });

    it("should return error when session creation fails", async () => {
      vi.mocked(zitadelModule.addHumanUser).mockResolvedValue({
        userId: mockUserId,
      } as any);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        allowRegister: true,
      } as any);

      vi.mocked(cookieModule.createSessionAndUpdateCookie).mockResolvedValue(undefined as any);

      const result = await registerUser({
        email: mockEmail,
        firstName: "John",
        lastName: "Doe",
        organization: mockOrganization,
      });

      expect(result).toHaveProperty("error");
      expect(result.error).toBe("errors.couldNotCreateSession");
    });

    it("should return error when user cannot be found after creation", async () => {
      const mockPassword = "SecurePassword123!";

      vi.mocked(zitadelModule.addHumanUser).mockResolvedValue({
        userId: mockUserId,
      } as any);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        allowRegister: true,
      } as any);

      vi.mocked(cookieModule.createSessionAndUpdateCookie).mockResolvedValue({
        id: "session-id",
        factors: {
          user: {
            id: mockUserId,
            loginName: mockEmail,
            organizationId: mockOrganization,
          },
        },
      } as any);

      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: undefined,
      } as any);

      const result = await registerUser({
        email: mockEmail,
        firstName: "John",
        lastName: "Doe",
        password: mockPassword,
        organization: mockOrganization,
      });

      expect(result).toHaveProperty("error");
      expect(result.error).toBe("errors.userNotFound");
    });
  });
});
