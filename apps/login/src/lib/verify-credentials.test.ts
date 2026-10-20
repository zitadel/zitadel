import { createServiceForHost } from "@/lib/service";
import { Code, ConnectError } from "@connectrpc/connect";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { verifyApiCredentials } from "./verify-credentials";

vi.mock("@/lib/service", () => ({
  createServiceForHost: vi.fn(),
}));

vi.mock("@/lib/logger", () => ({
  createLogger: () => ({ info: vi.fn(), warn: vi.fn(), error: vi.fn() }),
}));

const ENV_KEYS = [
  "ZITADEL_API_URL",
  "ZITADEL_SERVICE_USER_TOKEN",
  "ZITADEL_LOGINCLIENT_KEYFILE",
  "AUDIENCE",
  "SYSTEM_USER_ID",
  "SYSTEM_USER_PRIVATE_KEY",
  "SYSTEM_USER_PRIVATE_KEY_FILE",
] as const;

describe("verifyApiCredentials", () => {
  const savedEnv: Partial<Record<(typeof ENV_KEYS)[number], string | undefined>> = {};

  beforeEach(() => {
    for (const key of ENV_KEYS) {
      savedEnv[key] = process.env[key];
      delete process.env[key];
    }
    process.env.ZITADEL_API_URL = "http://localhost:8080";
    process.env.ZITADEL_SERVICE_USER_TOKEN = "token";
    vi.mocked(createServiceForHost).mockReset();
  });

  afterEach(() => {
    for (const key of ENV_KEYS) {
      if (savedEnv[key] === undefined) {
        delete process.env[key];
      } else {
        process.env[key] = savedEnv[key];
      }
    }
    vi.restoreAllMocks();
  });

  test("returns ok when the API accepts the credentials", async () => {
    const getGeneralSettings = vi.fn().mockResolvedValue({ allowedLanguages: ["en"] });
    vi.mocked(createServiceForHost).mockResolvedValue({ getGeneralSettings } as any);

    await expect(verifyApiCredentials()).resolves.toBe("ok");

    expect(createServiceForHost).toHaveBeenCalledWith(expect.anything(), { baseUrl: "http://localhost:8080" });
    expect(getGeneralSettings).toHaveBeenCalledTimes(1);
  });

  test("bounds the startup RPC with a timeout", async () => {
    const getGeneralSettings = vi.fn().mockResolvedValue({});
    vi.mocked(createServiceForHost).mockResolvedValue({ getGeneralSettings } as any);

    await verifyApiCredentials();

    expect(getGeneralSettings).toHaveBeenCalledWith({}, { timeoutMs: expect.any(Number) });
    expect(getGeneralSettings.mock.calls[0][1].timeoutMs).toBeGreaterThan(0);
  });

  test("returns skipped when ZITADEL_API_URL is not set", async () => {
    delete process.env.ZITADEL_API_URL;

    await expect(verifyApiCredentials()).resolves.toBe("skipped");

    expect(createServiceForHost).not.toHaveBeenCalled();
  });

  test("returns missing when no credentials are configured", async () => {
    delete process.env.ZITADEL_SERVICE_USER_TOKEN;

    await expect(verifyApiCredentials()).resolves.toBe("missing");

    expect(createServiceForHost).not.toHaveBeenCalled();
  });

  test("returns invalid when the credentials cannot be loaded or signed", async () => {
    vi.mocked(createServiceForHost).mockRejectedValue(
      new Error('Failed to read login client key file "/keys/login.json": ENOENT: no such file or directory'),
    );

    await expect(verifyApiCredentials()).resolves.toBe("invalid");
  });

  test.each([
    ["Unauthenticated", Code.Unauthenticated],
    ["PermissionDenied", Code.PermissionDenied],
  ])("returns rejected when the API responds with %s", async (_, code) => {
    vi.mocked(createServiceForHost).mockResolvedValue({
      getGeneralSettings: vi.fn().mockRejectedValue(new ConnectError("invalid token", code)),
    } as any);

    await expect(verifyApiCredentials()).resolves.toBe("rejected");
  });

  test("returns skipped when ZITADEL_API_URL does not resolve to an instance", async () => {
    vi.mocked(createServiceForHost).mockResolvedValue({
      getGeneralSettings: vi.fn().mockRejectedValue(new ConnectError("no instanceHost specified", Code.NotFound)),
    } as any);

    await expect(verifyApiCredentials()).resolves.toBe("skipped");
  });

  test.each([
    ["Unavailable", Code.Unavailable],
    ["DeadlineExceeded", Code.DeadlineExceeded],
    ["Internal", Code.Internal],
  ])("returns unreachable when the API responds with %s", async (_, code) => {
    vi.mocked(createServiceForHost).mockResolvedValue({
      getGeneralSettings: vi.fn().mockRejectedValue(new ConnectError("failed", code)),
    } as any);

    await expect(verifyApiCredentials()).resolves.toBe("unreachable");
  });

  test("returns unreachable on non-Connect errors", async () => {
    vi.mocked(createServiceForHost).mockResolvedValue({
      getGeneralSettings: vi.fn().mockRejectedValue(new Error("ECONNRESET")),
    } as any);

    await expect(verifyApiCredentials()).resolves.toBe("unreachable");
  });
});
