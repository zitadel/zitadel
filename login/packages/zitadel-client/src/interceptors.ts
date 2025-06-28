import type { Interceptor } from "@connectrpc/connect";

/**
 * Creates an interceptor that adds an Authorization header with a Bearer token.
 * @param token
 */
export function NewAuthorizationBearerInterceptor(token: string): Interceptor {
  return (next) => (req) => {
    // TODO: I am not what is the intent of checking for the Authorization header
    //  and setting it if it is not present.
    if (!req.header.get("Authorization")) {
      req.header.set("Authorization", `Bearer ${token}`);
    }
    return next(req);
  };
}
