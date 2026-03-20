// @vitest-environment node
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";

describe("loginClientKeyToken", () => {
  const originalEnv = process.env;
  let mockReadFile: ReturnType<typeof vi.fn>;
  let mockNewSystemToken: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    process.env = { ...originalEnv };
    vi.resetModules();

    mockReadFile = vi.fn();
    mockNewSystemToken = vi.fn();

    vi.doMock("fs/promises", () => ({ readFile: mockReadFile }));
    vi.doMock("@zitadel/client/node", () => ({ newSystemToken: mockNewSystemToken }));
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  test("should read key file and create token with hardcoded subject", async () => {
    process.env.ZITADEL_LOGINCLIENT_KEYFILE = "/path/to/key.pem";
    process.env.AUDIENCE = "https://api.zitadel.cloud";

    mockReadFile.mockResolvedValue("-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----");
    mockNewSystemToken.mockResolvedValue("signed-jwt-token");

    const { loginClientKeyToken } = await import("./api");
    const token = await loginClientKeyToken();

    expect(token).toBe("signed-jwt-token");
    expect(mockReadFile).toHaveBeenCalledWith("/path/to/key.pem", "utf-8");
    expect(mockNewSystemToken).toHaveBeenCalledWith({
      audience: "https://api.zitadel.cloud",
      subject: "login-client",
      key: "-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----",
    });
  });

  test("should cache key and not re-read file on subsequent calls", async () => {
    process.env.ZITADEL_LOGINCLIENT_KEYFILE = "/path/to/key.pem";
    process.env.AUDIENCE = "https://api.zitadel.cloud";

    mockReadFile.mockResolvedValue("cached-key");
    mockNewSystemToken.mockResolvedValue("token");

    const { loginClientKeyToken } = await import("./api");
    await loginClientKeyToken();
    await loginClientKeyToken();

    expect(mockReadFile).toHaveBeenCalledTimes(1);
  });

  test("should fall back to ZITADEL_API_URL when AUDIENCE is not set", async () => {
    process.env.ZITADEL_LOGINCLIENT_KEYFILE = "/path/to/key.pem";
    process.env.AUDIENCE = undefined as any;
    process.env.ZITADEL_API_URL = "https://zitadel.example.com";

    mockReadFile.mockResolvedValue("key-content");
    mockNewSystemToken.mockResolvedValue("token");

    const { loginClientKeyToken } = await import("./api");
    await loginClientKeyToken();

    expect(mockNewSystemToken).toHaveBeenCalledWith(
      expect.objectContaining({ audience: "https://zitadel.example.com" }),
    );
  });

  test("should throw a clear error when key file cannot be read", async () => {
    process.env.ZITADEL_LOGINCLIENT_KEYFILE = "/nonexistent/key.pem";

    mockReadFile.mockRejectedValue(new Error("ENOENT: no such file or directory"));

    const { loginClientKeyToken } = await import("./api");
    await expect(loginClientKeyToken()).rejects.toThrow(
      'Failed to read login client key file "/nonexistent/key.pem"',
    );
  });
});
