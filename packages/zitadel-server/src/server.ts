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

export function getServer(name?: string): ZitadelServer | undefined {
  return name
    ? apps.find((a) => a.name === name)
    : apps.length === 1
    ? apps[0]
    : undefined;
}
