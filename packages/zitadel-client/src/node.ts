import { createGrpcTransport, GrpcTransportOptions } from "@connectrpc/connect-node";
import { importPKCS8, SignJWT } from "jose";
import { NewAuthorizationBearerInterceptor } from "./interceptors";

/**
 * Create a server transport with the given token and configuration options.
 * @param token
 * @param opts
 */
export function createServerTransport(token: string, opts: GrpcTransportOptions) {
  return createGrpcTransport({
    ...opts,
    interceptors: [...(opts.interceptors || []), NewAuthorizationBearerInterceptor(token)],
  });
}

export async function newSystemToken() {
  return await new SignJWT({})
    .setProtectedHeader({ alg: "RS256" })
    .setIssuedAt()
    .setExpirationTime("1h")
    .setIssuer(process.env.ZITADEL_SYSTEM_API_USERID ?? "")
    .setSubject(process.env.ZITADEL_SYSTEM_API_USERID ?? "")
    .setAudience(process.env.ZITADEL_ISSUER ?? "")
    .sign(await importPKCS8(process.env.ZITADEL_SYSTEM_API_KEY ?? "", "RS256"));
}
