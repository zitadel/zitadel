let apps: ZitadelApp[] = [];

export interface ZitadelClientProps {
  clientId: string;
  apiUrl: string; // process.env.ZITADEL_API_URL
  token: string;
  adminToken?: string;
  managementToken?: string;
}

export interface ZitadelOptions extends ZitadelClientProps {
  name?: string;
}

export interface ZitadelApp {
  name: string | undefined;
  config: ZitadelClientProps;
}

export async function initializeApp(
  config: ZitadelClientProps,
  name?: string
): Promise<ZitadelApp> {
  const app = { config, name };
  return app;
}

export function getApps(): ZitadelApp[] {
  return apps;
}

export function getApp(name?: string): ZitadelApp | undefined {
  return name
    ? apps.find((a) => a.name === name)
    : apps.length === 1
    ? apps[0]
    : undefined;
}
