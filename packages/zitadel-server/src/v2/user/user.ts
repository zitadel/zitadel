import { CompatServiceDefinition } from "nice-grpc/lib/service-definitions";

import {
  UserServiceClient,
  UserServiceDefinition,
} from "../../proto/server/zitadel/user/v2beta/user_service";

import { ZitadelServer, createClient, getServers } from "../../server";

export const getUser = (server?: string | ZitadelServer) => {
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

  return createClient<UserServiceClient>(
    UserServiceDefinition as CompatServiceDefinition,
    config.apiUrl,
    config.token,
  );
};
