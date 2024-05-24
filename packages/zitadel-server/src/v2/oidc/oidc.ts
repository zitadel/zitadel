import { CompatServiceDefinition } from "nice-grpc/lib/service-definitions";

import { ZitadelServer, createClient, getServers } from "../../server";
import { OIDCServiceClient, OIDCServiceDefinition } from ".";

export const getOidc = (server?: string | ZitadelServer) => {
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

  return createClient<OIDCServiceClient>(
    OIDCServiceDefinition as CompatServiceDefinition,
    config.apiUrl,
    config.token,
  );
};
