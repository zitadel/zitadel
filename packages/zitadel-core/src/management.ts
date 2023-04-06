import { CompatServiceDefinition } from "nice-grpc/lib/service-definitions";

import { createChannel, createClientFactory } from "nice-grpc";
import {
  ManagementServiceClient,
  ManagementServiceDefinition,
} from "./proto/server/zitadel/management";

import { authMiddleware } from "./middleware";
import { ZitadelApp } from "./core";

const createClient = <Client>(
  definition: CompatServiceDefinition,
  accessToken: string
) => {
  const apiUrl = process.env.ZITADEL_API_URL;

  if (!apiUrl) {
    throw Error("ZITADEL_API_URL not set");
  }

  const channel = createChannel(process.env.ZITADEL_API_URL);
  return createClientFactory()
    .use(authMiddleware(accessToken))
    .create(definition, channel) as Client;
};

export const getManagement = (app?: ZitadelApp) =>
  createClient<ManagementServiceClient>(
    ManagementServiceDefinition,
    process.env.ZITADEL_ADMIN_TOKEN ?? ""
  );
