import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { getIframeOrigins } from "./security-settings";

// Suppress logger output during tests
vi.mock("@/lib/logger", () => ({
  createLogger: () => ({
    warn: vi.fn(),
    info: vi.fn(),
    error: vi.fn(),
    debug: vi.fn(),
  }),
}));

vi.mock("@/lib/deployment", () => ({
  hasSystemUserCredentials: () => false,
  hasLoginClientKey: () => false,
  hasServiceUserToken: () => true,
}));

const RESPONSE = {
  settings: {
    embeddedIframe: {
      enabled: true,
      allowedOrigins: ["https://allowed.example.com"],
    },
  },
};

describe("getIframeOrigins", () => {
  // Results are cached per instance host, so each test uses a fresh host
  let instanceHost: string;
  let seq = 0;
  const fetchMock = vi.fn();

  beforeEach(() => {
    instanceHost = `instance-${++seq}.example.com`;
    fetchMock.mockResolvedValue(new Response(JSON.stringify(RESPONSE), { status: 200 }));
    vi.stubGlobal("fetch", fetchMock);
    vi.stubEnv("ZITADEL_SERVICE_USER_TOKEN", "token");
  });

  afterEach(() => {
    fetchMock.mockReset();
    vi.unstubAllGlobals();
    vi.unstubAllEnvs();
  });

  // Headers combines duplicate names into a comma-separated value, which
  // exposes headers that were accidentally sent twice
  const sentHeaders = () => new Headers(fetchMock.mock.calls[0][1].headers);

  test("sends instance and public host headers", async () => {
    const origins = await getIframeOrigins("http://zitadel:8080", instanceHost, "public.example.com");

    expect(origins).toEqual(["https://allowed.example.com"]);
    expect(sentHeaders().get("x-zitadel-instance-host")).toBe(instanceHost);
    expect(sentHeaders().get("x-zitadel-public-host")).toBe("public.example.com");
  });

  test("CUSTOM_REQUEST_HEADERS replaces derived headers regardless of casing", async () => {
    vi.stubEnv("CUSTOM_REQUEST_HEADERS", "X-Zitadel-Public-Host:custom.example.com");

    await getIframeOrigins("http://zitadel:8080", instanceHost, "public.example.com");

    expect(sentHeaders().get("x-zitadel-public-host")).toBe("custom.example.com");
  });

  test("CUSTOM_REQUEST_HEADERS removes headers regardless of casing", async () => {
    vi.stubEnv("CUSTOM_REQUEST_HEADERS", "X-Zitadel-Public-Host:");

    await getIframeOrigins("http://zitadel:8080", instanceHost, "public.example.com");

    expect(sentHeaders().get("x-zitadel-public-host")).toBeNull();
  });
});
