import { CompatServiceDefinition } from "nice-grpc/lib/service-definitions";

import {
  SettingsServiceClient,
  SettingsServiceDefinition,
} from "../../proto/server/zitadel/settings/v2beta/settings_service";

import { ZitadelServer, createClient, getServers } from "../../server";

export const getSettings = (server?: string | ZitadelServer) => {
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
    config.token,
  );
};
