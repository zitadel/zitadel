import { describe, expect, test, beforeEach, afterEach } from "vitest";
import { hasSystemUserCredentials, hasServiceUserToken, hasLoginServiceKey, getLoginSystemUserId } from "./deployment";

describe("Deployment utilities", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    // Reset environment before each test
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  describe("hasSystemUserCredentials", () => {
    test("should return true when all system user credentials are present", () => {
      process.env.AUDIENCE = "https://api.zitadel.cloud";
      process.env.SYSTEM_USER_ID = "12345";
      process.env.SYSTEM_USER_PRIVATE_KEY = "-----BEGIN PRIVATE KEY-----\n...";

      expect(hasSystemUserCredentials()).toBe(true);
    });

    test("should return false when AUDIENCE is missing", () => {
      process.env.AUDIENCE = undefined as any;
      process.env.SYSTEM_USER_ID = "12345";
      process.env.SYSTEM_USER_PRIVATE_KEY = "-----BEGIN PRIVATE KEY-----\n...";

      expect(hasSystemUserCredentials()).toBe(false);
    });

    test("should return false when SYSTEM_USER_ID is missing", () => {
      process.env.AUDIENCE = "https://api.zitadel.cloud";
      process.env.SYSTEM_USER_ID = undefined as any;
      process.env.SYSTEM_USER_PRIVATE_KEY = "-----BEGIN PRIVATE KEY-----\n...";

      expect(hasSystemUserCredentials()).toBe(false);
    });

    test("should return false when SYSTEM_USER_PRIVATE_KEY is missing", () => {
      process.env.AUDIENCE = "https://api.zitadel.cloud";
      process.env.SYSTEM_USER_ID = "12345";
      process.env.SYSTEM_USER_PRIVATE_KEY = undefined as any;

      expect(hasSystemUserCredentials()).toBe(false);
    });

    test("should return false when all credentials are missing", () => {
      process.env.AUDIENCE = undefined as any;
      process.env.SYSTEM_USER_ID = undefined as any;
      process.env.SYSTEM_USER_PRIVATE_KEY = undefined as any;

      expect(hasSystemUserCredentials()).toBe(false);
    });
  });

  describe("hasLoginServiceKey", () => {
    test("should return true with ZITADEL_LOGIN_SERVICE_KEY_FILE and ZITADEL_LOGIN_SYSTEM_USER_ID", () => {
      process.env.ZITADEL_LOGIN_SERVICE_KEY_FILE = "/path/to/key.pem";
      process.env.ZITADEL_LOGIN_SYSTEM_USER_ID = "login-user";

      expect(hasLoginServiceKey()).toBe(true);
    });

    test("should return true with ZITADEL_LOGIN_SERVICE_KEY_FILE and SYSTEM_USER_ID fallback", () => {
      process.env.ZITADEL_LOGIN_SERVICE_KEY_FILE = "/path/to/key.pem";
      process.env.ZITADEL_LOGIN_SYSTEM_USER_ID = undefined as any;
      process.env.SYSTEM_USER_ID = "system-user";

      expect(hasLoginServiceKey()).toBe(true);
    });

    test("should return false when ZITADEL_LOGIN_SERVICE_KEY_FILE is missing", () => {
      process.env.ZITADEL_LOGIN_SERVICE_KEY_FILE = undefined as any;
      process.env.ZITADEL_LOGIN_SYSTEM_USER_ID = "login-user";

      expect(hasLoginServiceKey()).toBe(false);
    });

    test("should return false when no user ID is set", () => {
      process.env.ZITADEL_LOGIN_SERVICE_KEY_FILE = "/path/to/key.pem";
      process.env.ZITADEL_LOGIN_SYSTEM_USER_ID = undefined as any;
      process.env.SYSTEM_USER_ID = undefined as any;

      expect(hasLoginServiceKey()).toBe(false);
    });
  });

  describe("getLoginSystemUserId", () => {
    test("should return ZITADEL_LOGIN_SYSTEM_USER_ID when set", () => {
      process.env.ZITADEL_LOGIN_SYSTEM_USER_ID = "login-user";
      process.env.SYSTEM_USER_ID = "system-user";

      expect(getLoginSystemUserId()).toBe("login-user");
    });

    test("should return SYSTEM_USER_ID as fallback", () => {
      process.env.ZITADEL_LOGIN_SYSTEM_USER_ID = undefined as any;
      process.env.SYSTEM_USER_ID = "system-user";

      expect(getLoginSystemUserId()).toBe("system-user");
    });

    test("should return undefined when neither is set", () => {
      process.env.ZITADEL_LOGIN_SYSTEM_USER_ID = undefined as any;
      process.env.SYSTEM_USER_ID = undefined as any;

      expect(getLoginSystemUserId()).toBeUndefined();
    });
  });

  describe("hasServiceUserToken", () => {
    test("should return true when ZITADEL_SERVICE_USER_TOKEN is present", () => {
      process.env.ZITADEL_SERVICE_USER_TOKEN = "token123";

      expect(hasServiceUserToken()).toBe(true);
    });

    test("should return false when ZITADEL_SERVICE_USER_TOKEN is not set", () => {
      process.env.ZITADEL_SERVICE_USER_TOKEN = undefined as any;

      expect(hasServiceUserToken()).toBe(false);
    });

    test("should return false when ZITADEL_SERVICE_USER_TOKEN is empty string", () => {
      process.env.ZITADEL_SERVICE_USER_TOKEN = "";

      expect(hasServiceUserToken()).toBe(false);
    });
  });
});
