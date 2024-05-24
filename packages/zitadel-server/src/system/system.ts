import { CompatServiceDefinition } from "nice-grpc/lib/service-definitions";
import { importPKCS8, SignJWT } from "jose";

import { createChannel, createClientFactory } from "nice-grpc";
import {
  SystemServiceClient,
  SystemServiceDefinition,
} from "../proto/server/zitadel/system";
import { authMiddleware } from "../middleware";

const createSystemClient = <Client>(
  definition: CompatServiceDefinition,
  accessToken: string,
) => {
  const channel = createChannel(process.env.ZITADEL_SYSTEM_API_URL ?? "");
  return createClientFactory()
    .use(authMiddleware(accessToken))
    .create(definition, channel) as Client;
};

export const getSystem = async () => {
  const token = await new SignJWT({})
    .setProtectedHeader({ alg: "RS256" })
    .setIssuedAt()
    .setExpirationTime("1h")
    .setIssuer(process.env.ZITADEL_SYSTEM_API_USERID ?? "")
    .setSubject(process.env.ZITADEL_SYSTEM_API_USERID ?? "")
    .setAudience(process.env.ZITADEL_ISSUER ?? "")
    .sign(await importPKCS8(process.env.ZITADEL_SYSTEM_API_KEY ?? "", "RS256"));

  return createSystemClient<SystemServiceClient>(
    SystemServiceDefinition as CompatServiceDefinition,
    token,
  );
};
