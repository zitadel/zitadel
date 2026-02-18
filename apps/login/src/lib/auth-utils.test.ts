import { describe, it, expect } from "vitest";
import { validateAuthRequest, isRSCRequest, isValidLanguage, getValidLocaleFromUILocales } from "./auth-utils";

describe("auth-utils", () => {
  describe("isValidLanguage", () => {
    it("should return true for valid language codes", () => {
      expect(isValidLanguage("en")).toBe(true);
      expect(isValidLanguage("de")).toBe(true);
      expect(isValidLanguage("fr")).toBe(true);
      expect(isValidLanguage("zh")).toBe(true);
    });

    it("should return false for invalid language codes", () => {
      expect(isValidLanguage("invalid")).toBe(false);
      expect(isValidLanguage("xx")).toBe(false);
      expect(isValidLanguage("")).toBe(false);
    });

    it("should be case-sensitive", () => {
      expect(isValidLanguage("EN")).toBe(false);
      expect(isValidLanguage("De")).toBe(false);
    });
  });

  describe("getValidLocaleFromUILocales", () => {
    it("should return null for undefined or empty array", () => {
      expect(getValidLocaleFromUILocales(undefined)).toBeNull();
      expect(getValidLocaleFromUILocales([])).toBeNull();
    });

    it("should extract valid language from simple locale codes", () => {
      expect(getValidLocaleFromUILocales(["de"])).toBe("de");
      expect(getValidLocaleFromUILocales(["en", "de"])).toBe("en");
    });

    it("should extract language code from language tags", () => {
      expect(getValidLocaleFromUILocales(["de-CH"])).toBe("de");
      expect(getValidLocaleFromUILocales(["en-US"])).toBe("en");
      expect(getValidLocaleFromUILocales(["zh-CN"])).toBe("zh");
    });

    it("should return first valid language when multiple provided", () => {
      expect(getValidLocaleFromUILocales(["xx-XX", "de", "en"])).toBe("de");
    });

    it("should return null when no valid languages found", () => {
      expect(getValidLocaleFromUILocales(["xx", "yy", "zz"])).toBeNull();
    });

    it("should handle case-insensitive matching", () => {
      expect(getValidLocaleFromUILocales(["DE"])).toBe("de");
      expect(getValidLocaleFromUILocales(["En-US"])).toBe("en");
    });
  });

  describe("validateAuthRequest", () => {
    it("should return null when no auth parameters are present", () => {
      const searchParams = new URLSearchParams();
      const result = validateAuthRequest(searchParams);
      expect(result).toBeNull();
    });

    it("should return requestId when explicitly provided", () => {
      const searchParams = new URLSearchParams({
        requestId: "explicit-request-123",
      });
      const result = validateAuthRequest(searchParams);
      expect(result).toBe("explicit-request-123");
    });

    it("should generate OIDC requestId from authRequest parameter", () => {
      const searchParams = new URLSearchParams({
        authRequest: "oidc-auth-456",
      });
      const result = validateAuthRequest(searchParams);
      expect(result).toBe("oidc_oidc-auth-456");
    });

    it("should generate SAML requestId from samlRequest parameter", () => {
      const searchParams = new URLSearchParams({
        samlRequest: "saml-req-789",
      });
      const result = validateAuthRequest(searchParams);
      expect(result).toBe("saml_saml-req-789");
    });

    it("should prioritize explicit requestId over authRequest", () => {
      const searchParams = new URLSearchParams({
        requestId: "explicit-123",
        authRequest: "oidc-456",
      });
      const result = validateAuthRequest(searchParams);
      expect(result).toBe("explicit-123");
    });

    it("should prioritize explicit requestId over samlRequest", () => {
      const searchParams = new URLSearchParams({
        requestId: "explicit-123",
        samlRequest: "saml-456",
      });
      const result = validateAuthRequest(searchParams);
      expect(result).toBe("explicit-123");
    });

    it("should handle both authRequest and samlRequest, preferring OIDC", () => {
      const searchParams = new URLSearchParams({
        authRequest: "oidc-123",
        samlRequest: "saml-456",
      });
      const result = validateAuthRequest(searchParams);
      expect(result).toBe("oidc_oidc-123");
    });

    it("should handle empty string requestId", () => {
      const searchParams = new URLSearchParams({
        requestId: "",
      });
      const result = validateAuthRequest(searchParams);
      expect(result).toBeNull();
    });

    it("should handle whitespace in requestId", () => {
      const searchParams = new URLSearchParams({
        requestId: "  request-with-spaces  ",
      });
      const result = validateAuthRequest(searchParams);
      expect(result).toBe("  request-with-spaces  ");
    });

    it("should handle special characters in requestId", () => {
      const searchParams = new URLSearchParams({
        requestId: "request-!@#$%^&*()",
      });
      const result = validateAuthRequest(searchParams);
      expect(result).toBe("request-!@#$%^&*()");
    });

    it("should handle URL-encoded values", () => {
      const searchParams = new URLSearchParams({
        authRequest: "oidc encoded",
      });
      const result = validateAuthRequest(searchParams);
      expect(result).toBe("oidc_oidc encoded");
    });

    it("should handle very long requestId values", () => {
      const longId = "a".repeat(1000);
      const searchParams = new URLSearchParams({
        requestId: longId,
      });
      const result = validateAuthRequest(searchParams);
      expect(result).toBe(longId);
    });
  });

  describe("isRSCRequest", () => {
    it("should return true when _rsc parameter is present", () => {
      const searchParams = new URLSearchParams({
        _rsc: "1",
      });
      const result = isRSCRequest(searchParams);
      expect(result).toBe(true);
    });

    it("should return true when _rsc parameter has empty value", () => {
      const searchParams = new URLSearchParams({
        _rsc: "",
      });
      const result = isRSCRequest(searchParams);
      expect(result).toBe(true);
    });

    it("should return false when _rsc parameter is not present", () => {
      const searchParams = new URLSearchParams({
        other: "param",
      });
      const result = isRSCRequest(searchParams);
      expect(result).toBe(false);
    });

    it("should return false for empty search params", () => {
      const searchParams = new URLSearchParams();
      const result = isRSCRequest(searchParams);
      expect(result).toBe(false);
    });

    it("should return true regardless of _rsc value", () => {
      const searchParams = new URLSearchParams({
        _rsc: "any-value-123",
      });
      const result = isRSCRequest(searchParams);
      expect(result).toBe(true);
    });

    it("should work with multiple parameters", () => {
      const searchParams = new URLSearchParams({
        authRequest: "123",
        _rsc: "1",
        other: "param",
      });
      const result = isRSCRequest(searchParams);
      expect(result).toBe(true);
    });

    it("should be case-sensitive for parameter name", () => {
      const searchParams = new URLSearchParams({
        _RSC: "1",
      });
      const result = isRSCRequest(searchParams);
      expect(result).toBe(false);
    });
  });

  describe("integration: validateAuthRequest and isRSCRequest", () => {
    it("should handle typical OIDC auth flow with RSC", () => {
      const searchParams = new URLSearchParams({
        authRequest: "oidc-123",
        _rsc: "1",
      });

      const requestId = validateAuthRequest(searchParams);
      const isRSC = isRSCRequest(searchParams);

      expect(requestId).toBe("oidc_oidc-123");
      expect(isRSC).toBe(true);
    });

    it("should handle SAML flow without RSC", () => {
      const searchParams = new URLSearchParams({
        samlRequest: "saml-456",
      });

      const requestId = validateAuthRequest(searchParams);
      const isRSC = isRSCRequest(searchParams);

      expect(requestId).toBe("saml_saml-456");
      expect(isRSC).toBe(false);
    });
  });
});
