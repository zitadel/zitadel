import { CompatServiceDefinition } from "nice-grpc/lib/service-definitions";

import { createChannel, createClientFactory } from "nice-grpc";
import {
  AuthServiceClient,
  AuthServiceDefinition,
} from "./proto/server/zitadel/auth";
import {
  AdminServiceClient,
  AdminServiceDefinition,
} from "./proto/server/zitadel/admin";
import { authMiddleware } from "./middleware";

const createClient = <Client>(
  definition: CompatServiceDefinition,
  accessToken: string
) => {
  const channel = createChannel(process.env.ZITADEL_API_URL);
  return createClientFactory()
    .use(authMiddleware(accessToken))
    .create(definition, channel) as Client;
};

export const getAuth = async () =>
  createClient<AuthServiceClient>(AuthServiceDefinition, "");

export const getAdmin = () =>
  createClient<AdminServiceClient>(
    AdminServiceDefinition,
    process.env.ZITADEL_ADMIN_TOKEN ?? ""
  );
