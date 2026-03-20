import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { loginClientKeyToken } from "./api";

vi.mock("fs/promises", () => ({
  readFile: vi.fn(),
}));

vi.mock("@zitadel/client/node", () => ({
  newSystemToken: vi.fn(),
}));

import { readFile } from "fs/promises";
import { newSystemToken } from "@zitadel/client/node";

const mockedReadFile = vi.mocked(readFile);
const mockedNewSystemToken = vi.mocked(newSystemToken);

describe("loginClientKeyToken", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
    vi.clearAllMocks();
    // Reset the module-level cache by re-importing would be ideal,
    // but since we mock readFile we can verify call counts instead.
  });

  afterEach(() => {
    process.env = originalEnv;
    vi.restoreAllMocks();
  });

  test("should read key file and create token with hardcoded subject", async () => {
    process.env.ZITADEL_LOGINCLIENT_KEYFILE = "/path/to/key.pem";
    process.env.AUDIENCE = "https://api.zitadel.cloud";

    mockedReadFile.mockResolvedValue("-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----");
    mockedNewSystemToken.mockResolvedValue("signed-jwt-token");

    const token = await loginClientKeyToken();

    expect(token).toBe("signed-jwt-token");
    expect(mockedReadFile).toHaveBeenCalledWith("/path/to/key.pem", "utf-8");
    expect(mockedNewSystemToken).toHaveBeenCalledWith({
      audience: "https://api.zitadel.cloud",
      subject: "login-client",
      key: "-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----",
    });
  });

  test("should fall back to ZITADEL_API_URL when AUDIENCE is not set", async () => {
    process.env.ZITADEL_LOGINCLIENT_KEYFILE = "/path/to/key.pem";
    process.env.AUDIENCE = undefined as any;
    process.env.ZITADEL_API_URL = "https://zitadel.example.com";

    mockedReadFile.mockResolvedValue("key-content");
    mockedNewSystemToken.mockResolvedValue("token");

    await loginClientKeyToken();

    expect(mockedNewSystemToken).toHaveBeenCalledWith(
      expect.objectContaining({ audience: "https://zitadel.example.com" }),
    );
  });

  test("should throw a clear error when key file cannot be read", async () => {
    process.env.ZITADEL_LOGINCLIENT_KEYFILE = "/nonexistent/key.pem";

    mockedReadFile.mockRejectedValue(new Error("ENOENT: no such file or directory"));

    await expect(loginClientKeyToken()).rejects.toThrow(
      'Failed to read login client key file "/nonexistent/key.pem"',
    );
  });
});
