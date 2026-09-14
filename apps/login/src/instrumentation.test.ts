import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";

const verifyApiCredentials = vi.fn();
const registerNode = vi.fn();

vi.mock("./lib/verify-credentials", () => ({ verifyApiCredentials }));
vi.mock("./instrumentation.node", () => ({ registerNode }));

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
  });

  afterEach(() => {
    for (const [key, value] of Object.entries(savedEnv)) {
      if (value === undefined) {
        delete process.env[key];
      } else {
        process.env[key] = value;
      }
    }
    vi.restoreAllMocks();
  });

  test("verifies credentials on startup and continues when they are accepted", async () => {
    const register = await loadRegister();

    await register();

    expect(verifyApiCredentials).toHaveBeenCalledTimes(1);
    expect(exitSpy).not.toHaveBeenCalled();
  });

  test.each(["rejected", "missing"])("exits the process when the credential check returns %s", async (result) => {
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
});
