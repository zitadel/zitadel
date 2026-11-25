import { ClassifiedConnectError } from "@/lib/server/grpc/interceptors/error-classification";
import { Code, ConnectError } from "@connectrpc/connect";
import { NextRequest } from "next/server";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { GET } from "./route";

vi.mock("@/lib/server/auth-utils", () => ({
  isRSCRequest: vi.fn(() => false),
  validateAuthRequest: vi.fn((searchParams: URLSearchParams) => searchParams.get("requestId")),
}));

vi.mock("@/lib/server/cookies", () => ({
  getAllSessions: vi.fn(async () => []),
}));

vi.mock("@/lib/server/flow-initiation", () => ({
  handleOIDCFlowInitiation: vi.fn(),
  handleSAMLFlowInitiation: vi.fn(),
}));

vi.mock("@/lib/server/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://api.example.com" } })),
}));

vi.mock("@/lib/server/zitadel", () => ({
  listSessions: vi.fn(),
}));

const apiError = (code: Code, message: string) => new ClassifiedConnectError(new ConnectError(message, code));

describe("GET /login", () => {
  let handleOIDCFlowInitiation: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const flowInitiation = await import("@/lib/server/flow-initiation");
    handleOIDCFlowInitiation = vi.mocked(flowInitiation.handleOIDCFlowInitiation);
  });

  const get = () => GET(new NextRequest("https://example.com/login?requestId=oidc_abc123"));

  test.each([
    ["the auth request was already handled", Code.FailedPrecondition, 400],
    ["the auth request is gone", Code.NotFound, 404],
  ])("answers a rejection with its status when %s", async (_, code, status) => {
    handleOIDCFlowInitiation.mockRejectedValue(apiError(code, "Auth request rejected (OIDC-x)"));

    const res = await get();

    expect(res.status).toBe(status);
    expect(await res.json()).toEqual({ error: "Auth request rejected (OIDC-x)" });
  });

  test.each([
    ["a server fault", apiError(Code.Internal, "boom")],
    ["the API rejecting the login's own credentials", apiError(Code.Unauthenticated, "Token invalid")],
    ["a missing permission of the login", apiError(Code.PermissionDenied, "No permission")],
    ["an error that did not come from the API", new Error("boom")],
  ])("answers %s with a 500", async (_, error) => {
    handleOIDCFlowInitiation.mockRejectedValue(error);

    const res = await get();

    expect(res.status).toBe(500);
    expect(await res.json()).toEqual({ error: "Internal server error" });
  });
});
