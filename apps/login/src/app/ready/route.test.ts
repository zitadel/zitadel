import { createServiceForHost } from "@/lib/service";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { GET } from "./route";

vi.mock("@/lib/service", () => ({
  createServiceForHost: vi.fn(),
}));

vi.mock("@/lib/logger", () => ({
  createLogger: () => ({ error: vi.fn() }),
}));

describe("GET /ready", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
    process.env.ZITADEL_API_URL = "http://localhost:8080";
  });

  afterEach(() => {
    process.env = originalEnv;
    vi.restoreAllMocks();
  });

  test("should return 200 when gRPC call succeeds", async () => {
    vi.mocked(createServiceForHost).mockResolvedValue({
      getGeneralSettings: vi.fn().mockResolvedValue({ allowedLanguages: ["en"] }),
    } as any);

    const response = await GET();

    expect(response.status).toBe(200);
    expect(response.headers.get("Content-Type")).toBe("text/plain");
    expect(await response.text()).toBe("OK");
  });

  test("should return 503 when gRPC call fails", async () => {
    vi.mocked(createServiceForHost).mockResolvedValue({
      getGeneralSettings: vi.fn().mockRejectedValue(new Error("connection refused")),
    } as any);

    const response = await GET();

    expect(response.status).toBe(503);
    expect(response.headers.get("Content-Type")).toBe("text/plain");
    expect(await response.text()).toBe("Service unavailable");
  });

  test("should return 503 when authentication fails", async () => {
    vi.mocked(createServiceForHost).mockRejectedValue(new Error("No authentication credentials found"));

    const response = await GET();

    expect(response.status).toBe(503);
    expect(response.headers.get("Content-Type")).toBe("text/plain");
    expect(await response.text()).toBe("Service unavailable");
  });

  test("should return 503 when ZITADEL_API_URL is not set", async () => {
    process.env.ZITADEL_API_URL = undefined as any;

    const response = await GET();

    expect(response.status).toBe(503);
    expect(response.headers.get("Content-Type")).toBe("text/plain");
    expect(await response.text()).toBe("Service unavailable");
  });
});
