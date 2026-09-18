import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { GET } from "./route";

vi.mock("@/lib/logger", () => ({
  createLogger: () => ({ error: vi.fn(), warn: vi.fn() }),
}));

function mockApiReady(status = 200) {
  const fetchMock = vi.fn().mockResolvedValue(
    new Response(status === 200 ? '"ok"' : "{}", {
      status,
      headers: { "Content-Type": "application/json" },
    }),
  );
  vi.stubGlobal("fetch", fetchMock);
  return fetchMock;
}

describe("GET /ready", () => {
  let savedApiUrl: string | undefined;
  let savedCustomHeaders: string | undefined;

  beforeEach(() => {
    savedApiUrl = process.env.ZITADEL_API_URL;
    savedCustomHeaders = process.env.CUSTOM_REQUEST_HEADERS;
    process.env.ZITADEL_API_URL = "http://localhost:8080";
    delete process.env.CUSTOM_REQUEST_HEADERS;
  });

  afterEach(() => {
    if (savedApiUrl === undefined) {
      delete process.env.ZITADEL_API_URL;
    } else {
      process.env.ZITADEL_API_URL = savedApiUrl;
    }
    if (savedCustomHeaders === undefined) {
      delete process.env.CUSTOM_REQUEST_HEADERS;
    } else {
      process.env.CUSTOM_REQUEST_HEADERS = savedCustomHeaders;
    }
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });

  test("should return 200 when the API is ready", async () => {
    const fetchMock = mockApiReady();

    const response = await GET();

    expect(response.status).toBe(200);
    expect(response.headers.get("Content-Type")).toBe("text/plain");
    expect(await response.text()).toBe("OK");
    expect(fetchMock).toHaveBeenCalledTimes(1);
    expect(fetchMock.mock.calls[0][0].toString()).toBe("http://localhost:8080/debug/ready");
  });

  test("should return 503 when the API is not ready", async () => {
    mockApiReady(412);

    const response = await GET();

    expect(response.status).toBe(503);
    expect(response.headers.get("Content-Type")).toBe("text/plain");
    expect(await response.text()).toBe("Service unavailable");
  });

  test("should return 503 when the API is unreachable", async () => {
    vi.stubGlobal("fetch", vi.fn().mockRejectedValue(new Error("connection refused")));

    const response = await GET();

    expect(response.status).toBe(503);
    expect(await response.text()).toBe("Service unavailable");
  });

  test("should return 503 when ZITADEL_API_URL is not set", async () => {
    const fetchMock = mockApiReady();
    delete process.env.ZITADEL_API_URL;

    const response = await GET();

    expect(response.status).toBe(503);
    expect(response.headers.get("Content-Type")).toBe("text/plain");
    expect(await response.text()).toBe("Service unavailable");
    expect(fetchMock).not.toHaveBeenCalled();
  });

  test("should forward CUSTOM_REQUEST_HEADERS to the readiness request", async () => {
    process.env.CUSTOM_REQUEST_HEADERS = "Host: zitadel.internal, X-Forwarded-Host: login.example.com";
    const fetchMock = mockApiReady();

    await GET();

    const headers = fetchMock.mock.calls[0][1].headers as Headers;
    expect(headers.get("Host")).toBe("zitadel.internal");
    expect(headers.get("X-Forwarded-Host")).toBe("login.example.com");
  });
});
