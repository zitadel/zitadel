import { Code, ConnectError } from "@connectrpc/connect";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { ClassifiedConnectError } from "../grpc/interceptors/error-classification";
import { completeDeviceAuthorization } from "./device";

vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

// this returns the key itself so tests can assert on translation keys
vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("@/lib/zitadel", () => ({
  authorizeOrDenyDeviceAuthorization: vi.fn(),
}));

describe("completeDeviceAuthorization", () => {
  let mockAuthorize: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { authorizeOrDenyDeviceAuthorization } = await import("@/lib/zitadel");

    vi.mocked(headers).mockResolvedValue({} as any);
    vi.mocked(getServiceConfig).mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } } as any);
    mockAuthorize = vi.mocked(authorizeOrDenyDeviceAuthorization);
  });

  test("resolves without error on success", async () => {
    mockAuthorize.mockResolvedValue({});

    await expect(
      completeDeviceAuthorization("device-1", { sessionId: "session-1", sessionToken: "token-1" }),
    ).resolves.toBeUndefined();
  });

  test("returns an error for an expired or consumed device code instead of throwing", async () => {
    mockAuthorize.mockRejectedValue(
      new ClassifiedConnectError(new ConnectError("Errors.DeviceAuth.NotFound", Code.FailedPrecondition)),
    );

    await expect(completeDeviceAuthorization("device-1")).resolves.toEqual({ error: "errors.couldNotComplete" });
  });

  test("rethrows genuine server errors so they still surface as 500", async () => {
    mockAuthorize.mockRejectedValue(new ClassifiedConnectError(new ConnectError("database down", Code.Internal)));

    await expect(completeDeviceAuthorization("device-1")).rejects.toThrow("database down");
  });
});
