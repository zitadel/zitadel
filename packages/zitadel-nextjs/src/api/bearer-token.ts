import {
  extractBearerTokenFromHeaders,
  validateBearerToken,
  type BearerTokenHeaderSource,
  type ValidateBearerTokenOptions,
} from "@zitadel/zitadel-js/api/bearer-token";

type ValidatedBearerTokenPayload<T> = Awaited<
  ReturnType<typeof validateBearerToken<T>>
>;

export interface BearerTokenRequest {
  headers: BearerTokenHeaderSource;
}

/**
 * Extracts a bearer token from a Next.js API request Authorization header.
 */
export function extractBearerTokenFromRequest(
  request: BearerTokenRequest,
): string | null {
  return extractBearerTokenFromHeaders(request.headers);
}

/**
 * Extracts and validates a bearer token from a Next.js API request.
 *
 * Returns `null` when no bearer token is present.
 */
export async function validateBearerTokenFromRequest<
  T = Record<string, unknown>,
>(
  request: BearerTokenRequest,
  options: ValidateBearerTokenOptions,
): Promise<ValidatedBearerTokenPayload<T> | null> {
  const token = extractBearerTokenFromRequest(request);
  if (!token) {
    return null;
  }
  return (await validateBearerToken<T>(
    token,
    options,
  )) as ValidatedBearerTokenPayload<T>;
}

export type { BearerTokenHeaderSource, ValidateBearerTokenOptions };
