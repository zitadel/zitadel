import {
  ZitadelServer,
  ZitadelServerOptions,
  user,
  settings,
  getServers,
  initializeServer,
  session,
  GetGeneralSettingsResponse,
  CreateSessionResponse,
  GetBrandingSettingsResponse,
  GetPasswordComplexitySettingsResponse,
  GetLegalAndSupportSettingsResponse,
  AddHumanUserResponse,
  BrandingSettings,
  ListSessionsResponse,
  LegalAndSupportSettings,
  PasswordComplexitySettings,
  GetSessionResponse,
  VerifyEmailResponse,
  SetSessionResponse,
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
): Promise<BrandingSettings | undefined> {
  const settingsService = settings.getSettings(server);
  return settingsService
    .getBrandingSettings({}, {})
    .then((resp: GetBrandingSettingsResponse) => resp.settings);
}

export function getGeneralSettings(
  server: ZitadelServer
): Promise<string[] | undefined> {
  const settingsService = settings.getSettings(server);
  return settingsService
    .getGeneralSettings({}, {})
    .then((resp: GetGeneralSettingsResponse) => resp.supportedLanguages);
}

export function getLegalAndSupportSettings(
  server: ZitadelServer
): Promise<LegalAndSupportSettings | undefined> {
  const settingsService = settings.getSettings(server);
  return settingsService
    .getLegalAndSupportSettings({}, {})
    .then((resp: GetLegalAndSupportSettingsResponse) => {
      return resp.settings;
    });
}

export function getPasswordComplexitySettings(
  server: ZitadelServer
): Promise<PasswordComplexitySettings | undefined> {
  const settingsService = settings.getSettings(server);

  return settingsService
    .getPasswordComplexitySettings({}, {})
    .then((resp: GetPasswordComplexitySettingsResponse) => resp.settings);
}

export function createSession(
  server: ZitadelServer,
  loginName: string
): Promise<CreateSessionResponse | undefined> {
  const sessionService = session.getSession(server);
  return sessionService.createSession({ checks: { user: { loginName } } }, {});
}

export function setSession(
  server: ZitadelServer,
  sessionId: string,
  sessionToken: string,
  password: string
): Promise<SetSessionResponse | undefined> {
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
): Promise<GetSessionResponse | undefined> {
  const sessionService = session.getSession(server);
  return sessionService.getSession({ sessionId, sessionToken }, {});
}

export function listSessions(
  server: ZitadelServer,
  ids: string[]
): Promise<ListSessionsResponse | undefined> {
  const sessionService = session.getSession(server);
  const query = { offset: 0, limit: 100, asc: true };
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
      {}
    )
    .then((resp: AddHumanUserResponse) => {
      return resp.userId;
    });
}

export function verifyEmail(
  server: ZitadelServer,
  userId: string,
  verificationCode: string
): Promise<VerifyEmailResponse> {
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
