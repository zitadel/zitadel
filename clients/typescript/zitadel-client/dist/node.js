import {
  NewAuthorizationBearerInterceptor
} from "./chunk-H5CIUTXZ.js";

// src/node.ts
import { createGrpcTransport } from "@connectrpc/connect-node";
import { importPKCS8, SignJWT } from "jose";
function createServerTransport(token, opts) {
  return createGrpcTransport({
    ...opts,
    interceptors: [...opts.interceptors || [], NewAuthorizationBearerInterceptor(token)]
  });
}
async function newSystemToken({
  audience,
  subject,
  key,
  expirationTime
}) {
  return await new SignJWT({}).setProtectedHeader({ alg: "RS256" }).setIssuedAt().setExpirationTime(expirationTime ?? "1h").setIssuer(subject).setSubject(subject).setAudience(audience).sign(await importPKCS8(key, "RS256"));
}
export {
  createServerTransport,
  newSystemToken
};
//# sourceMappingURL=node.js.map