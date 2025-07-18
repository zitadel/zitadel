import { createGrpcTransport, GrpcTransportOptions } from "@connectrpc/connect-node";
import { importPKCS8, SignJWT } from "jose";
import { NewAuthorizationBearerInterceptor } from "./interceptors.js";

/**
 * Create a server transport using grpc with the given token and configuration options.
 * @param token
 * @param opts
 */
export function createServerTransport(token: string, opts: GrpcTransportOptions) {
  return createGrpcTransport({
    ...opts,
    interceptors: [...(opts.interceptors || []), NewAuthorizationBearerInterceptor(token)],
  });
}

export async function newSystemToken({
  audience,
  subject,
  key,
  expirationTime,
}: {
  audience: string;
  subject: string;
  key: string;
  expirationTime?: number | string | Date;
}) {
  return await new SignJWT({})
    .setProtectedHeader({ alg: "RS256" })
    .setIssuedAt()
    .setExpirationTime(expirationTime ?? "1h")
    .setIssuer(subject)
    .setSubject(subject)
    .setAudience(audience)
    .sign(await importPKCS8(key, "RS256"));
}
