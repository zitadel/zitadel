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
import type { AddHumanUserResponse } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import type { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import type { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import type { GetUserByIDResponse } from "@zitadel/proto/zitadel/user/v2/user_service_pb";

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
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://zitadel-test.zitadel.cloud" } })),
}));

describe("registerUser server action", () => {
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
      vi.mocked(zitadelModule.addHumanUser).mockResolvedValue({
        userId: mockUserId,
      } as AddHumanUserResponse);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        passkeysType: PasskeysType.ALLOWED,
        allowRegister: true,
      } as LoginSettings);

      vi.mocked(cookieModule.createSessionAndUpdateCookie).mockResolvedValue({
        id: "session-id",
        factors: {
          user: {
            id: mockUserId,
            loginName: mockEmail,
            organizationId: mockOrganization,
          },
        },
      } as Session);

      const result = await registerUser({
        email: mockEmail,
        firstName: "John",
        lastName: "Doe",
        organization: mockOrganization,
        requestId: mockRequestId,
      });

      expect(result).toHaveProperty("redirect");
      if ("redirect" in result) {
        expect(result.redirect).toContain("/passkey/set");
        expect(result.redirect).toContain(`loginName=${encodeURIComponent(mockEmail)}`);
        expect(result.redirect).toContain(`organization=${mockOrganization}`);
        expect(result.redirect).toContain(`requestId=${mockRequestId}`);
      }
    });
  });

  describe("password registration", () => {
    it("should create user with password and check email verification", async () => {
      const mockPassword = "SecurePassword123!";

      vi.mocked(zitadelModule.addHumanUser).mockResolvedValue({
        userId: mockUserId,
      } as AddHumanUserResponse);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        allowRegister: true,
        allowUsernamePassword: true,
        passwordCheckLifetime: { seconds: BigInt(300) },
      } as LoginSettings);

      vi.mocked(cookieModule.createSessionAndUpdateCookie).mockResolvedValue({
        id: "session-id",
        factors: {
          user: {
            id: mockUserId,
            loginName: mockEmail,
            organizationId: mockOrganization,
          },
        },
      } as Session);

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
      } as GetUserByIDResponse);

      const { completeFlowOrGetUrl } = await import("../client");
      vi.mocked(completeFlowOrGetUrl).mockResolvedValue({
        redirect: "/dashboard",
      } as { redirect: string });

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
      if ("error" in result) {
        expect(result.error).toBe("errors.couldNotCreateUser");
      }
    });

    it("should return error when session creation fails", async () => {
      vi.mocked(zitadelModule.addHumanUser).mockResolvedValue({
        userId: mockUserId,
      } as AddHumanUserResponse);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        allowRegister: true,
      } as LoginSettings);

      vi.mocked(cookieModule.createSessionAndUpdateCookie).mockResolvedValue(undefined as unknown as Session);

      const result = await registerUser({
        email: mockEmail,
        firstName: "John",
        lastName: "Doe",
        organization: mockOrganization,
      });

      expect(result).toHaveProperty("error");
      if ("error" in result) {
        expect(result.error).toBe("errors.couldNotCreateSession");
      }
    });

    it("should return error when user cannot be found after creation", async () => {
      const mockPassword = "SecurePassword123!";

      vi.mocked(zitadelModule.addHumanUser).mockResolvedValue({
        userId: mockUserId,
      } as AddHumanUserResponse);

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        allowRegister: true,
      } as LoginSettings);

      vi.mocked(cookieModule.createSessionAndUpdateCookie).mockResolvedValue({
        id: "session-id",
        factors: {
          user: {
            id: mockUserId,
            loginName: mockEmail,
            organizationId: mockOrganization,
          },
        },
      } as Session);

      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: undefined,
      } as GetUserByIDResponse);

      const result = await registerUser({
        email: mockEmail,
        firstName: "John",
        lastName: "Doe",
        password: mockPassword,
        organization: mockOrganization,
      });

      expect(result).toHaveProperty("error");
      if ("error" in result) {
        expect(result.error).toBe("errors.userNotFound");
      }
    });
  });
});
