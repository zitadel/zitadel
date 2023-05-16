import {
  ZitadelServer,
  ZitadelServerOptions,
  management,
  settings,
  getServers,
  initializeServer,
  session,
} from "@zitadel/server";

export const zitadelConfig: ZitadelServerOptions = {
  name: "zitadel login",
  apiUrl: process.env.ZITADEL_API_URL ?? "",
  token: process.env.ZITADEL_SERVICE_USER_TOKEN ?? "",
};

let server: ZitadelServer;

if (!getServers().length) {
  console.log("initialize server");
  server = initializeServer(zitadelConfig);
}

export function getBrandingSettings(
  server: ZitadelServer
): Promise<any | undefined> {
  // settings.branding_settings.BrandingSettings
  const settingsService = settings.getSettings(server);
  return settingsService
    .getBrandingSettings(
      {},
      {
        // metadata: orgMetadata(process.env.ZITADEL_ORG_ID ?? "")
      }
    )
    .then((resp) => resp.settings);
}

export function getGeneralSettings(
  server: ZitadelServer
): Promise<any | undefined> {
  // settings.branding_settings.BrandingSettings
  const settingsService = settings.getSettings(server);
  return settingsService
    .getGeneralSettings(
      {},
      {
        // metadata: orgMetadata(process.env.ZITADEL_ORG_ID ?? "")
      }
    )
    .then((resp) => resp.supportedLanguages);
}

export function getLegalAndSupportSettings(
  server: ZitadelServer
): Promise<any | undefined> {
  const settingsService = settings.getSettings(server);
  return settingsService
    .getLegalAndSupportSettings(
      {},
      {
        //  metadata: orgMetadata(process.env.ZITADEL_ORG_ID ?? "")
      }
    )
    .then((resp) => resp.settings);
}

export function getPasswordComplexitySettings(
  server: ZitadelServer
): Promise<any | undefined> {
  const settingsService = settings.getSettings(server);

  return settingsService
    .getPasswordComplexitySettings(
      {},
      {
        // metadata: orgMetadata(process.env.ZITADEL_ORG_ID ?? "")
      }
    )
    .then((resp) => resp.settings);
}

export function createSession(
  server: ZitadelServer,
  loginName: string
): Promise<any | undefined> {
  const sessionService = session.getSession(server);
  return sessionService.createSession({ checks: { user: { loginName } } }, {});
}

export function setSession(
  server: ZitadelServer,
  loginName: string
): Promise<any | undefined> {
  const sessionService = session.getSession(server);
  return sessionService.setSession({ checks: { user: { loginName } } }, {});
}

export type AddHumanUserData = {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
};
export function addHumanUser(
  server: ZitadelServer,
  { email, firstName, lastName, password }: AddHumanUserData
): Promise<string> {
  const mgmt = management.getManagement(server);
  return mgmt
    .addHumanUser(
      {
        email: { email, isEmailVerified: false },
        userName: email,
        profile: { firstName, lastName },
        initialPassword: password,
      },
      {
        // metadata: orgMetadata(process.env.ZITADEL_ORG_ID ?? "")
      }
    )
    .then((resp) => {
      console.log("added user", resp.userId);
      return resp.userId;
    });
}

export { server };
