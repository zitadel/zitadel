import { beforeEach, describe, expect, test, vi } from "vitest";
import { getDeviceAuthorizationRequest } from "./oidc";

vi.mock("@/lib/zitadel", () => ({
  getDeviceAuthorizationRequest: vi.fn(),
  listSessions: vi.fn(),
}));

vi.mock("@/lib/cookies", () => ({
  getAllSessions: vi.fn(),
}));

vi.mock("../device-code-check", () => ({
  setDeviceCodeCheckCookie: vi.fn(),
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://example.com" } })),
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(() => new Headers()),
}));

const deviceAuthorizationRequest = { id: "encrypted-device-code", appName: "TV", clientId: "client-1" };

describe("getDeviceAuthorizationRequest", () => {
  let zitadelGetDeviceAuthorizationRequest: any;
  let listSessions: any;
  let getAllSessions: any;
  let setDeviceCodeCheckCookie: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const zitadel = await import("@/lib/zitadel");
    const cookies = await import("@/lib/cookies");
    const deviceCodeCheck = await import("../device-code-check");

    zitadelGetDeviceAuthorizationRequest = vi.mocked(zitadel.getDeviceAuthorizationRequest);
    listSessions = vi.mocked(zitadel.listSessions);
    getAllSessions = vi.mocked(cookies.getAllSessions);
    setDeviceCodeCheckCookie = vi.mocked(deviceCodeCheck.setDeviceCodeCheckCookie);

    zitadelGetDeviceAuthorizationRequest.mockResolvedValue({ deviceAuthorizationRequest });
    getAllSessions.mockResolvedValue([]);
  });

  test("does not set the check cookie or look up sessions when the lookup returns no request", async () => {
    zitadelGetDeviceAuthorizationRequest.mockResolvedValue({ deviceAuthorizationRequest: undefined });

    const result = await getDeviceAuthorizationRequest("ABCD-1234");

    expect(result).toEqual({ deviceAuthorizationRequest: undefined, nextUrl: undefined });
    expect(setDeviceCodeCheckCookie).not.toHaveBeenCalled();
    expect(getAllSessions).not.toHaveBeenCalled();
    expect(listSessions).not.toHaveBeenCalled();
  });

  test("does not set the check cookie when the lookup fails", async () => {
    zitadelGetDeviceAuthorizationRequest.mockRejectedValue(new Error("not found"));

    await expect(getDeviceAuthorizationRequest("ABCD-1234")).rejects.toThrow("not found");
    expect(setDeviceCodeCheckCookie).not.toHaveBeenCalled();
    expect(listSessions).not.toHaveBeenCalled();
  });

  test("continues to /loginname without session cookies and does not call listSessions", async () => {
    const result = await getDeviceAuthorizationRequest("ABCD-1234");

    expect(setDeviceCodeCheckCookie).toHaveBeenCalledWith("ABCD-1234");
    expect(listSessions).not.toHaveBeenCalled();
    expect(result).toEqual({
      deviceAuthorizationRequest,
      nextUrl: "/loginname?requestId=device_ABCD-1234",
    });
  });

  test("continues to /accounts when the browser has a live session", async () => {
    getAllSessions.mockResolvedValue([{ id: "session-1" }, { id: "" }, { id: "session-2" }]);
    listSessions.mockResolvedValue({ sessions: [{ id: "session-1" }] });

    const result = await getDeviceAuthorizationRequest("ABCD-1234");

    expect(listSessions).toHaveBeenCalledWith({
      serviceConfig: { baseUrl: "https://example.com" },
      ids: ["session-1", "session-2"],
    });
    expect(setDeviceCodeCheckCookie).toHaveBeenCalledWith("ABCD-1234");
    expect(result.nextUrl).toBe("/accounts?requestId=device_ABCD-1234");
  });

  test("continues to /loginname when all session cookies point to ended sessions", async () => {
    getAllSessions.mockResolvedValue([{ id: "session-1" }]);
    listSessions.mockResolvedValue({ sessions: [] });

    const result = await getDeviceAuthorizationRequest("ABCD-1234");

    expect(result.nextUrl).toBe("/loginname?requestId=device_ABCD-1234");
  });

  test("continues to /loginname and does not throw when listSessions fails", async () => {
    getAllSessions.mockResolvedValue([{ id: "session-1" }]);
    listSessions.mockRejectedValue(new Error("unavailable"));

    const result = await getDeviceAuthorizationRequest("ABCD-1234");

    expect(setDeviceCodeCheckCookie).toHaveBeenCalledWith("ABCD-1234");
    expect(result.nextUrl).toBe("/loginname?requestId=device_ABCD-1234");
  });

  test("builds the next url only from the user code", async () => {
    const result = await getDeviceAuthorizationRequest("AB&next=https://evil.example");

    expect(result.nextUrl).toBe("/loginname?requestId=device_AB%26next%3Dhttps%3A%2F%2Fevil.example");
  });
});
