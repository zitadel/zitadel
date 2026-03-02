import type { JWTPayload } from "jose";

import { createAuthorizationBearerInterceptor } from "../interceptors.js";
import { verifyJwt, type VerifyJwtOptions } from "../token.js";
import { createGrpcTransport, type NodeTransportOptions } from "../transport.js";

const BEARER_TOKEN_HEADER_REGEX = /^Bearer\s+(\S+)\s*$/i;

export type BearerTokenHeaderSource =
  | { get(name: string): string | null }
  | Record<string, string | string[] | undefined>;

export interface ValidateBearerTokenOptions extends VerifyJwtOptions {
  keysEndpoint: string;
}

export interface BearerTokenTransportOptions
  extends Omit<NodeTransportOptions, "interceptors"> {
  token: string;
  interceptors?: NonNullable<NodeTransportOptions["interceptors"]>;
}

export function extractBearerTokenFromAuthorizationHeader(
  authorizationHeader: string | null | undefined,
): string | null {
  if (!authorizationHeader) {
    return null;
  }

  const match = BEARER_TOKEN_HEADER_REGEX.exec(authorizationHeader.trim());
  return match?.[1] ?? null;
}

function getAuthorizationHeader(headers: BearerTokenHeaderSource): string | null {
  if ("get" in headers && typeof headers.get === "function") {
    return headers.get("authorization");
  }

  const headerRecord = headers as Record<string, string | string[] | undefined>;
  const direct = headerRecord.authorization ?? headerRecord.Authorization;
  if (Array.isArray(direct)) {
    return direct[0] ?? null;
  }
  if (typeof direct === "string") {
    return direct;
  }

  for (const [name, value] of Object.entries(headerRecord)) {
    if (name.toLowerCase() !== "authorization") {
      continue;
    }
    if (Array.isArray(value)) {
      return value[0] ?? null;
    }
    return value ?? null;
  }

  return null;
}

export function extractBearerTokenFromHeaders(
  headers: BearerTokenHeaderSource,
): string | null {
  return extractBearerTokenFromAuthorizationHeader(getAuthorizationHeader(headers));
}

export async function validateBearerToken<T = JWTPayload>(
  token: string,
  options: ValidateBearerTokenOptions,
): Promise<T & JWTPayload> {
  const { keysEndpoint, ...verifyOptions } = options;
  return verifyJwt<T>(token, keysEndpoint, verifyOptions);
}

export function createBearerTokenInterceptor(token: string) {
  return createAuthorizationBearerInterceptor(token);
}

export function createBearerTokenTransport(options: BearerTokenTransportOptions) {
  const { token, interceptors, ...transportOptions } = options;
  return createGrpcTransport({
    ...transportOptions,
    interceptors: [
      createAuthorizationBearerInterceptor(token),
      ...(interceptors ?? []),
    ],
  } as NodeTransportOptions);
}
