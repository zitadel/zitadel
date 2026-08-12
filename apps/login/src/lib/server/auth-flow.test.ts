import { Code, ConnectError } from "@connectrpc/connect";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { ClassifiedConnectError } from "../grpc/interceptors/error-classification";
import { completeAuthFlow } from "./auth-flow";

vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("@/lib/cookies", () => ({
  getAllSessions: vi.fn(),
}));

vi.mock("@/lib/service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("@/lib/zitadel", () => ({
  listSessions: vi.fn(),
}));

vi.mock("@/lib/oidc", () => ({
  loginWithOIDCAndSession: vi.fn(),
}));

vi.mock("@/lib/saml", () => ({
  loginWithSAMLAndSession: vi.fn(),
}));

describe("completeAuthFlow", () => {
  let mockGetAllSessions: any;
  let mockListSessions: any;
  let mockLoginWithOIDCAndSession: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("@/lib/service-url");
    const { getAllSessions } = await import("@/lib/cookies");
    const { listSessions } = await import("@/lib/zitadel");
    const { loginWithOIDCAndSession } = await import("@/lib/oidc");

    vi.mocked(headers).mockResolvedValue({} as any);
    vi.mocked(getServiceConfig).mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } } as any);
    mockGetAllSessions = vi.mocked(getAllSessions);
    mockListSessions = vi.mocked(listSessions);
    mockLoginWithOIDCAndSession = vi.mocked(loginWithOIDCAndSession);

    mockGetAllSessions.mockResolvedValue([{ id: "session-1", token: "token-1", loginName: "user@example.com" }]);
  });

  test("rethrows when loading sessions fails with a server error", async () => {
    // listSessions resolves with an empty list for unknown ids — a rejection
    // is a genuine service failure and must keep the action failing
    mockListSessions.mockRejectedValue(new ClassifiedConnectError(new ConnectError("database down", Code.Internal)));

    await expect(completeAuthFlow({ sessionId: "session-1", requestId: "oidc_abc123" })).rejects.toThrow("database down");
    expect(mockLoginWithOIDCAndSession).not.toHaveBeenCalled();
  });

  test("degrades to no sessions when loading fails with a user error", async () => {
    mockListSessions.mockRejectedValue(new ClassifiedConnectError(new ConnectError("bad ids", Code.InvalidArgument)));
    mockLoginWithOIDCAndSession.mockResolvedValue({ error: "Session not found or invalid" });

    const result = await completeAuthFlow({ sessionId: "session-1", requestId: "oidc_abc123" });

    expect(result).toEqual({ error: "Session not found or invalid" });
    expect(mockLoginWithOIDCAndSession).toHaveBeenCalledWith(expect.objectContaining({ sessions: [] }));
  });

  test("passes loaded sessions through to the flow handler", async () => {
    const session = { id: "session-1", factors: { user: { loginName: "user@example.com" } } };
    mockListSessions.mockResolvedValue({ sessions: [session] });
    mockLoginWithOIDCAndSession.mockResolvedValue({ redirect: "https://app.example.com/callback" });

    const result = await completeAuthFlow({ sessionId: "session-1", requestId: "oidc_abc123" });

    expect(result).toEqual({ redirect: "https://app.example.com/callback" });
    expect(mockLoginWithOIDCAndSession).toHaveBeenCalledWith(expect.objectContaining({ sessions: [session] }));
  });
});
