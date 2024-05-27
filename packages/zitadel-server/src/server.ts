import { createChannel, createClientFactory } from "nice-grpc";
import { authMiddleware } from "./middleware";
import { CompatServiceDefinition } from "nice-grpc/lib/service-definitions";

let apps: ZitadelServer[] = [];

export interface ZitadelServerProps {
  apiUrl: string; // process.env.ZITADEL_API_URL
  token: string;
}

export interface ZitadelServerOptions extends ZitadelServerProps {
  name?: string;
}

export function initializeServer(
  config: ZitadelServerProps,
  name?: string,
): ZitadelServer {
  const server = new ZitadelServer(config, name);
  return server;
}

export class ZitadelServer {
  name: string | undefined;
  config: ZitadelServerProps;

  constructor(config: ZitadelServerProps, name?: string) {
    if (name) {
      this.name = name;
    }
    this.config = config;
  }
}

export function getServers(): ZitadelServer[] {
  return apps;
}

export function getServer(name?: string): ZitadelServer {
  if (name) {
    const found = apps.find((a) => a.name === name);
    if (found) {
      return found;
    } else {
      throw new Error("No server found");
    }
  } else {
    if (apps.length) {
      return apps[0];
    } else {
      throw new Error("No server found");
    }
  }
}

export const createClient = <Client>(
  definition: CompatServiceDefinition,
  apiUrl: string,
  token: string,
) => {
  if (!apiUrl) {
    throw Error("ZITADEL_API_URL not set");
  }

  const channel = createChannel(process.env.ZITADEL_API_URL ?? "");
  return createClientFactory()
    .use(authMiddleware(token))
    .create(definition, channel) as Client;
};
