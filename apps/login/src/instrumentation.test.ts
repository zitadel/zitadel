import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";

const verifyApiCredentials = vi.fn();
const registerNode = vi.fn();

// Keep the real FATAL_CREDENTIAL_CHECK_RESULTS, but replace the check itself
// and cut the import chain into the API client.
vi.mock("./lib/verify-credentials", async (importOriginal) => ({
  ...(await importOriginal<typeof import("./lib/verify-credentials")>()),
  verifyApiCredentials,
}));
vi.mock("@/lib/service", () => ({ createServiceForHost: vi.fn() }));
vi.mock("./instrumentation.node", () => ({ registerNode }));
const startupLogger = { warn: vi.fn(), error: vi.fn() };
vi.mock("./lib/logger", () => ({ createLogger: () => startupLogger }));

async function loadRegister() {
  vi.resetModules();
  const { register } = await import("./instrumentation");
  return register;
}

describe("register", () => {
  const savedEnv = {
    NEXT_RUNTIME: process.env.NEXT_RUNTIME,
    NODE_ENV: process.env.NODE_ENV,
    OTEL_SDK_DISABLED: process.env.OTEL_SDK_DISABLED,
  };
  let exitSpy: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    process.env.NEXT_RUNTIME = "nodejs";
    process.env.OTEL_SDK_DISABLED = "true";
    verifyApiCredentials.mockReset().mockResolvedValue("ok");
    registerNode.mockReset().mockResolvedValue(null);
    exitSpy = vi.spyOn(process, "exit").mockImplementation((() => undefined) as never);
    startupLogger.warn.mockReset();
    startupLogger.error.mockReset();
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "test-session-cookie-secret-at-least-32-chars");
  });

  afterEach(() => {
    for (const [key, value] of Object.entries(savedEnv)) {
      if (value === undefined) {
        delete process.env[key];
      } else {
        process.env[key] = value;
      }
    }
    vi.unstubAllEnvs();
    vi.restoreAllMocks();
  });

  test("verifies credentials on startup and continues when they are accepted", async () => {
    const register = await loadRegister();

    await register();

    expect(verifyApiCredentials).toHaveBeenCalledTimes(1);
    expect(exitSpy).not.toHaveBeenCalled();
  });

  test.each(["rejected", "invalid", "missing"])("exits the process when the credential check returns %s", async (result) => {
    verifyApiCredentials.mockResolvedValue(result);
    const register = await loadRegister();

    await register();

    expect(exitSpy).toHaveBeenCalledWith(1);
  });

  test.each(["unreachable", "skipped"])("continues when the credential check returns %s", async (result) => {
    verifyApiCredentials.mockResolvedValue(result);
    const register = await loadRegister();

    await register();

    expect(exitSpy).not.toHaveBeenCalled();
  });

  test("does nothing outside the nodejs runtime", async () => {
    process.env.NEXT_RUNTIME = "edge";
    const register = await loadRegister();

    await register();

    expect(verifyApiCredentials).not.toHaveBeenCalled();
    expect(registerNode).not.toHaveBeenCalled();
  });

  test("registers OpenTelemetry before the credential check when enabled", async () => {
    delete process.env.OTEL_SDK_DISABLED;
    const calls: string[] = [];
    registerNode.mockImplementation(async () => {
      calls.push("otel");
      return null;
    });
    verifyApiCredentials.mockImplementation(async () => {
      calls.push("credentials");
      return "ok";
    });
    const register = await loadRegister();

    await register();

    expect(calls).toEqual(["otel", "credentials"]);
  });

  test("still verifies credentials when OpenTelemetry is disabled", async () => {
    const register = await loadRegister();

    await register();

    expect(registerNode).not.toHaveBeenCalled();
    expect(verifyApiCredentials).toHaveBeenCalledTimes(1);
  });

  describe("session cookie secret notice", () => {
    test("logs nothing when a dedicated secret is configured", async () => {
      const register = await loadRegister();

      await register();

      expect(startupLogger.warn).not.toHaveBeenCalled();
      expect(startupLogger.error).not.toHaveBeenCalled();
    });

    test("warns about the deprecated credential fallback, also in development", async () => {
      vi.stubEnv("NODE_ENV", "development");
      vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", undefined);
      vi.stubEnv("ZITADEL_SERVICE_USER_TOKEN", "service-user-token");
      const register = await loadRegister();

      await register();

      expect(startupLogger.warn).toHaveBeenCalledWith(expect.stringMatching(/fallback is deprecated/));
    });

    test("logs an error for a too short secret before the credential check", async () => {
      vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "too-short");
      const calls: string[] = [];
      startupLogger.error.mockImplementation(() => calls.push("notice"));
      verifyApiCredentials.mockImplementation(async () => {
        calls.push("credentials");
        return "ok";
      });
      const register = await loadRegister();

      await register();

      expect(calls).toEqual(["notice", "credentials"]);
      expect(startupLogger.error).toHaveBeenCalledWith(expect.stringMatching(/at least 32 characters/));
    });
  });
});
