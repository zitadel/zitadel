import { beforeEach, describe, expect, test, vi } from "vitest";
import {
  extractBearerTokenFromRequest,
  validateBearerTokenFromRequest,
} from "./bearer-token.js";

const { extractBearerTokenFromHeadersMock, validateBearerTokenMock } = vi.hoisted(
  () => ({
    extractBearerTokenFromHeadersMock: vi.fn(),
    validateBearerTokenMock: vi.fn(),
  }),
);

vi.mock("@zitadel/zitadel-js/api/bearer-token", () => ({
  extractBearerTokenFromHeaders: extractBearerTokenFromHeadersMock,
  validateBearerToken: validateBearerTokenMock,
}));

describe("nextjs api/bearer-token helpers", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  test("extracts bearer token from request headers", () => {
    const headers = new Headers();
    extractBearerTokenFromHeadersMock.mockReturnValueOnce("access-token");

    const token = extractBearerTokenFromRequest({ headers });

    expect(extractBearerTokenFromHeadersMock).toHaveBeenCalledWith(headers);
    expect(token).toBe("access-token");
  });

  test("returns null when no bearer token is present", async () => {
    extractBearerTokenFromHeadersMock.mockReturnValueOnce(null);

    const payload = await validateBearerTokenFromRequest(
      { headers: new Headers() },
      { keysEndpoint: "https://issuer.example.com/oauth/v2/keys" },
    );

    expect(validateBearerTokenMock).not.toHaveBeenCalled();
    expect(payload).toBeNull();
  });

  test("validates bearer token extracted from request", async () => {
    const options = {
      keysEndpoint: "https://issuer.example.com/oauth/v2/keys",
      issuer: "https://issuer.example.com",
      audience: "client-id",
    };
    extractBearerTokenFromHeadersMock.mockReturnValueOnce("jwt-token");
    validateBearerTokenMock.mockResolvedValueOnce({ sub: "user-1" });

    const payload = await validateBearerTokenFromRequest(
      { headers: new Headers() },
      options,
    );

    expect(validateBearerTokenMock).toHaveBeenCalledWith("jwt-token", options);
    expect(payload).toEqual({ sub: "user-1" });
  });
});
