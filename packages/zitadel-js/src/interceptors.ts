import type { Interceptor } from "@connectrpc/connect";

/**
 * Creates a Connect interceptor that injects an `Authorization: Bearer <token>` header.
 *
 * If the request already has an Authorization header, it is left untouched.
 *
 * @example
 * ```ts
 * import { createGrpcTransport } from "@zitadel/zitadel-js";
 * import { createAuthorizationBearerInterceptor } from "@zitadel/zitadel-js";
 *
 * const transport = createGrpcTransport({
 *   baseUrl: "https://my.zitadel.cloud",
 *   interceptors: [createAuthorizationBearerInterceptor(token)],
 * });
 * ```
 */
export function createAuthorizationBearerInterceptor(
  token: string,
): Interceptor {
  return (next) => (req) => {
    if (!req.header.get("Authorization")) {
      req.header.set("Authorization", `Bearer ${token}`);
    }
    return next(req);
  };
}
