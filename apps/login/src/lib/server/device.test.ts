/**
 * Unit tests for device server actions.
 *
 * Tests core business logic for:
 * - Device authorization completion
 * - Device authorization request retrieval
 * - Session handling in device flows
 */

import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { completeDeviceAuthorization } from "./device";
import { getDeviceAuthorizationRequest } from "./oidc";
import * as zitadelModule from "@/lib/zitadel";

// Mock all dependencies
vi.mock("@/lib/zitadel");
vi.mock("next/headers", () => ({
  headers: vi.fn(() => Promise.resolve(new Map())),
}));
vi.mock("@/lib/service-url", () => ({
  getServiceUrlFromHeaders: vi.fn(() => ({ serviceUrl: "https://zitadel-test.zitadel.cloud" })),
}));

describe("Device server actions", () => {
  const mockServiceUrl = "https://zitadel-test.zitadel.cloud";
  const mockDeviceAuthorizationId = "device123";
  const mockUserCode = "USER-CODE-123";

  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("completeDeviceAuthorization", () => {
    it("should authorize device with session", async () => {
      const mockSession = {
        sessionId: "session123",
        sessionToken: "token123",
      };

      const mockResponse = {
        details: {
          id: mockDeviceAuthorizationId,
          approved: true,
        },
      };

      vi.mocked(zitadelModule.authorizeOrDenyDeviceAuthorization).mockResolvedValue(mockResponse as any);

      const result = await completeDeviceAuthorization(mockDeviceAuthorizationId, mockSession);

      expect(result).toEqual(mockResponse);
      expect(zitadelModule.authorizeOrDenyDeviceAuthorization).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        deviceAuthorizationId: mockDeviceAuthorizationId,
        session: mockSession,
      });
    });

    it("should deny device authorization without session", async () => {
      const mockResponse = {
        details: {
          id: mockDeviceAuthorizationId,
          approved: false,
        },
      };

      vi.mocked(zitadelModule.authorizeOrDenyDeviceAuthorization).mockResolvedValue(mockResponse as any);

      const result = await completeDeviceAuthorization(mockDeviceAuthorizationId);

      expect(result).toEqual(mockResponse);
      expect(zitadelModule.authorizeOrDenyDeviceAuthorization).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        deviceAuthorizationId: mockDeviceAuthorizationId,
        session: undefined,
      });
    });

    it("should handle authorization errors", async () => {
      vi.mocked(zitadelModule.authorizeOrDenyDeviceAuthorization).mockRejectedValue(new Error("Authorization failed"));

      await expect(completeDeviceAuthorization(mockDeviceAuthorizationId)).rejects.toThrow("Authorization failed");
    });

    it("should pass session parameters correctly", async () => {
      const mockSession = {
        sessionId: "session456",
        sessionToken: "token456",
      };

      vi.mocked(zitadelModule.authorizeOrDenyDeviceAuthorization).mockResolvedValue({} as any);

      await completeDeviceAuthorization(mockDeviceAuthorizationId, mockSession);

      const callArgs = vi.mocked(zitadelModule.authorizeOrDenyDeviceAuthorization).mock.calls[0][0];
      expect(callArgs.session).toEqual(mockSession);
      expect(callArgs.session?.sessionId).toBe("session456");
      expect(callArgs.session?.sessionToken).toBe("token456");
    });
  });

  describe("getDeviceAuthorizationRequest", () => {
    it("should retrieve device authorization request by user code", async () => {
      const mockRequest = {
        id: mockDeviceAuthorizationId,
        userCode: mockUserCode,
        clientId: "client123",
        scope: ["openid", "profile"],
      };

      vi.mocked(zitadelModule.getDeviceAuthorizationRequest).mockResolvedValue(mockRequest as any);

      const result = await getDeviceAuthorizationRequest(mockUserCode);

      expect(result).toEqual(mockRequest);
      expect(zitadelModule.getDeviceAuthorizationRequest).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        userCode: mockUserCode,
      });
    });

    it("should handle invalid user codes", async () => {
      vi.mocked(zitadelModule.getDeviceAuthorizationRequest).mockRejectedValue(new Error("User code not found"));

      await expect(getDeviceAuthorizationRequest("INVALID-CODE")).rejects.toThrow("User code not found");
    });

    it("should handle expired user codes", async () => {
      vi.mocked(zitadelModule.getDeviceAuthorizationRequest).mockRejectedValue(new Error("User code expired"));

      await expect(getDeviceAuthorizationRequest(mockUserCode)).rejects.toThrow("User code expired");
    });
  });
});
