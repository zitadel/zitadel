import {
  ZitadelServer,
  ZitadelServerOptions,
  user,
  settings,
  getServers,
  initializeServer,
  session,
  GetGeneralSettingsResponse,
  GetBrandingSettingsResponse,
  GetPasswordComplexitySettingsResponse,
  GetLegalAndSupportSettingsResponse,
  AddHumanUserResponse,
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
    .then((resp: GetBrandingSettingsResponse) => resp.settings);
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
    .then((resp: GetGeneralSettingsResponse) => resp.supportedLanguages);
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
    .then((resp: GetLegalAndSupportSettingsResponse) => resp.settings);
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
    .then((resp: GetPasswordComplexitySettingsResponse) => resp.settings);
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
  sessionId: string,
  sessionToken: string,
  password: string
): Promise<any | undefined> {
  const sessionService = session.getSession(server);
  return sessionService.setSession(
    { sessionId, sessionToken, checks: { password: { password } } },
    {}
  );
}

export function getSession(
  server: ZitadelServer,
  sessionId: string,
  sessionToken: string
): Promise<any | undefined> {
  const sessionService = session.getSession(server);
  return sessionService.getSession({ sessionId, sessionToken }, {});
}

export function listSessions(
  server: ZitadelServer,
  ids: string[]
): Promise<any | undefined> {
  const sessionService = session.getSession(server);
  const query = { offset: 0, limit: 100, asc: true };
  console.log(ids);
  const queries = [{ idsQuery: { ids } }];
  return sessionService.listSessions({ queries: queries }, {});
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
  const mgmt = user.getUser(server);
  return mgmt
    .addHumanUser(
      {
        email: { email },
        username: email,
        profile: { firstName, lastName },
        password: { password },
      },
      {
        // metadata: orgMetadata(process.env.ZITADEL_ORG_ID ?? "")
      }
    )
    .then((resp: AddHumanUserResponse) => {
      console.log("added user", resp.userId);
      return resp.userId;
    });
}

export function verifyEmail(
  server: ZitadelServer,
  userId: string,
  verificationCode: string
): Promise<any> {
  const userservice = user.getUser(server);
  return userservice.verifyEmail(
    {
      userId,
      verificationCode,
    },
    {}
  );
}

/**
 *
 * @param server
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */
export function setEmail(server: ZitadelServer, userId: string): Promise<any> {
  const userservice = user.getUser(server);
  return userservice.setEmail(
    {
      userId,
    },
    {}
  );
}

export { server };
