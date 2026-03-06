import type { Interceptor } from "@connectrpc/connect";
import { beforeEach, describe, expect, test, vi } from "vitest";
import {
  createBearerTokenInterceptor,
  createBearerTokenTransport,
  extractBearerTokenFromAuthorizationHeader,
  extractBearerTokenFromHeaders,
  validateBearerToken,
} from "./bearer-token.js";

const {
  createAuthorizationBearerInterceptorMock,
  createGrpcTransportMock,
  verifyJwtMock,
} = vi.hoisted(() => ({
  createAuthorizationBearerInterceptorMock: vi.fn(),
  createGrpcTransportMock: vi.fn(),
  verifyJwtMock: vi.fn(),
}));

vi.mock("../interceptors.js", () => ({
  createAuthorizationBearerInterceptor: createAuthorizationBearerInterceptorMock,
}));

vi.mock("../transport.js", () => ({
  createGrpcTransport: createGrpcTransportMock,
}));

vi.mock("../token.js", () => ({
  verifyJwt: verifyJwtMock,
}));

describe("api/bearer-token helpers", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  test("extracts a token from an Authorization header", () => {
    expect(
      extractBearerTokenFromAuthorizationHeader("Bearer access-token"),
    ).toBe("access-token");
    expect(
      extractBearerTokenFromAuthorizationHeader("bearer another-token"),
    ).toBe("another-token");
  });

  test("returns null for malformed Authorization headers", () => {
    expect(extractBearerTokenFromAuthorizationHeader(undefined)).toBeNull();
    expect(extractBearerTokenFromAuthorizationHeader("Basic abc")).toBeNull();
    expect(extractBearerTokenFromAuthorizationHeader("Bearer ")).toBeNull();
  });

  test("extracts bearer token from header containers", () => {
    expect(
      extractBearerTokenFromHeaders({ authorization: "Bearer from-object" }),
    ).toBe("from-object");
    expect(
      extractBearerTokenFromHeaders({ Authorization: ["Bearer from-array"] }),
    ).toBe("from-array");
    expect(
      extractBearerTokenFromHeaders({ AUTHORIZATION: "Bearer from-uppercase" }),
    ).toBe("from-uppercase");
    expect(
      extractBearerTokenFromHeaders(
        new Headers({ Authorization: "Bearer from-headers" }),
      ),
    ).toBe("from-headers");
  });

  test("validates bearer token via verifyJwt helper", async () => {
    verifyJwtMock.mockResolvedValueOnce({ sub: "user-1" });

    const payload = await validateBearerToken("jwt-token", {
      keysEndpoint: "https://issuer.example.com/oauth/v2/keys",
      issuer: "https://issuer.example.com",
      audience: "client-id",
    });

    expect(verifyJwtMock).toHaveBeenCalledWith(
      "jwt-token",
      "https://issuer.example.com/oauth/v2/keys",
      {
        issuer: "https://issuer.example.com",
        audience: "client-id",
      },
    );
    expect(payload).toEqual({ sub: "user-1" });
  });

  test("creates bearer token interceptor via existing interceptor helper", () => {
    const interceptor = ((next) => (req) => next(req)) as Interceptor;
    createAuthorizationBearerInterceptorMock.mockReturnValueOnce(interceptor);

    const result = createBearerTokenInterceptor("api-token");

    expect(createAuthorizationBearerInterceptorMock).toHaveBeenCalledWith(
      "api-token",
    );
    expect(result).toBe(interceptor);
  });

  test("creates grpc transport with bearer interceptor prepended", () => {
    const bearerInterceptor = ((next) => (req) => next(req)) as Interceptor;
    const customInterceptor = ((next) => (req) => next(req)) as Interceptor;
    const transport = { mocked: true };
    createAuthorizationBearerInterceptorMock.mockReturnValueOnce(
      bearerInterceptor,
    );
    createGrpcTransportMock.mockReturnValueOnce(transport);

    const result = createBearerTokenTransport({
      baseUrl: "https://api.example.com",
      httpVersion: "2",
      token: "api-token",
      interceptors: [customInterceptor],
    });

    expect(createGrpcTransportMock).toHaveBeenCalledWith({
      baseUrl: "https://api.example.com",
      httpVersion: "2",
      interceptors: [bearerInterceptor, customInterceptor],
    });
    expect(result).toBe(transport);
  });
});
