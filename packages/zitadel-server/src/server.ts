let apps: ZitadelServer[] = [];

export interface ZitadelServerProps {
  apiUrl: string; // process.env.ZITADEL_API_URL
  token: string;
}

export interface ZitadelServerOptions extends ZitadelServerProps {
  name?: string;
}

export interface ZitadelServer {
  name: string | undefined;
  config: ZitadelServerProps;
}

export async function initializeServer(
  config: ZitadelServerProps,
  name?: string
): Promise<ZitadelServer> {
  const app = { config, name };
  return app;
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
