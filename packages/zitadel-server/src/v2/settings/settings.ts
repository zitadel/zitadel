import { CompatServiceDefinition } from "nice-grpc/lib/service-definitions";

import { createChannel, createClientFactory } from "nice-grpc";
import {
  SettingsServiceClient,
  SettingsServiceDefinition,
} from "../../proto/server/zitadel/settings/v2alpha/settings_service";

import { authMiddleware } from "../../middleware";
import { ZitadelServer, getServers } from "../../server";

const createClient = <Client>(
  definition: CompatServiceDefinition,
  apiUrl: string,
  token: string
) => {
  if (!apiUrl) {
    throw Error("ZITADEL_API_URL not set");
  }

  const channel = createChannel(process.env.ZITADEL_API_URL ?? "");
  return createClientFactory()
    .use(authMiddleware(token))
    .create(definition, channel) as Client;
};

export const getSettings = (server?: string | ZitadelServer) => {
  console.log("init settings");
  let config;
  if (server && typeof server === "string") {
    const apps = getServers();
    config = apps.find((a) => a.name === server)?.config;
  } else if (server && typeof server === "object") {
    config = server.config;
  }

  if (!config) {
    throw Error("No ZITADEL server found");
  }

  return createClient<SettingsServiceClient>(
    SettingsServiceDefinition as CompatServiceDefinition,
    config.apiUrl,
    config.token
  );
};
