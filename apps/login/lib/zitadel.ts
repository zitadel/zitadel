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
  DeleteSessionResponse,
  VerifyPasskeyRegistrationResponse,
} from "@zitadel/server";
import { Metadata } from "nice-grpc";

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

export async function getBrandingSettings(
  server: ZitadelServer
): Promise<BrandingSettings | undefined> {
  const settingsService = settings.getSettings(server);
  return settingsService
    .getBrandingSettings({}, {})
    .then((resp: GetBrandingSettingsResponse) => resp.settings);
}

export async function getGeneralSettings(
  server: ZitadelServer
): Promise<string[] | undefined> {
  const settingsService = settings.getSettings(server);
  return settingsService
    .getGeneralSettings({}, {})
    .then((resp: GetGeneralSettingsResponse) => resp.supportedLanguages);
}

export async function getLegalAndSupportSettings(
  server: ZitadelServer
): Promise<LegalAndSupportSettings | undefined> {
  const settingsService = settings.getSettings(server);
  return settingsService
    .getLegalAndSupportSettings({}, {})
    .then((resp: GetLegalAndSupportSettingsResponse) => {
      return resp.settings;
    });
}

export async function getPasswordComplexitySettings(
  server: ZitadelServer
): Promise<PasswordComplexitySettings | undefined> {
  const settingsService = settings.getSettings(server);

  return settingsService
    .getPasswordComplexitySettings({}, {})
    .then((resp: GetPasswordComplexitySettingsResponse) => resp.settings);
}

export async function createSession(
  server: ZitadelServer,
  loginName: string
): Promise<CreateSessionResponse | undefined> {
  const sessionService = session.getSession(server);
  return sessionService.createSession({ checks: { user: { loginName } } }, {});
}

export async function setSession(
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

export async function getSession(
  server: ZitadelServer,
  sessionId: string,
  sessionToken: string
): Promise<GetSessionResponse | undefined> {
  const sessionService = session.getSession(server);
  return sessionService.getSession({ sessionId, sessionToken }, {});
}

export async function deleteSession(
  server: ZitadelServer,
  sessionId: string,
  sessionToken: string
): Promise<DeleteSessionResponse | undefined> {
  const sessionService = session.getSession(server);
  return sessionService.deleteSession({ sessionId, sessionToken }, {});
}

export async function listSessions(
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

export async function addHumanUser(
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

export async function verifyEmail(
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
export async function setEmail(
  server: ZitadelServer,
  userId: string
): Promise<any> {
  const userservice = user.getUser(server);
  return userservice.setEmail(
    {
      userId,
    },
    {}
  );
}

const bearerTokenMetadata = (token: string) =>
  new Metadata({ authorization: `Bearer ${token}` });

/**
 *
 * @param server
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */
export async function createPasskeyRegistrationLink(
  userId: string,
  sessionToken: string
): Promise<any> {
  //   this actions will be made from the currently seleected user
  //   const zitadelConfig: ZitadelServerOptions = {
  //     name: "zitadel login",
  //     apiUrl: process.env.ZITADEL_API_URL ?? "",
  //     token: "",
  //   };

  //   const authserver: ZitadelServer = initializeServer(zitadelConfig);
  //   console.log("server", authserver);
  const userservice = user.getUser(server);
  return userservice.createPasskeyRegistrationLink(
    {
      userId,
      returnCode: {},
    }
    // { metadata: bearerTokenMetadata(sessionToken) }
  );
}

/**
 *
 * @param server
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */
export async function verifyPasskeyRegistration(
  server: ZitadelServer,
  passkeyId: string,
  passkeyName: string,
  publicKeyCredential: any,
  userId: string
): Promise<VerifyPasskeyRegistrationResponse> {
  const userservice = user.getUser(server);
  return userservice.verifyPasskeyRegistration(
    {
      passkeyId,
      passkeyName,
      publicKeyCredential,
      userId,
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
export async function registerPasskey(
  userId: string,
  code: { id: string; code: string }
): Promise<any> {
  //   this actions will be made from the currently seleected user
  const zitadelConfig: ZitadelServerOptions = {
    name: "zitadel login",
    apiUrl: process.env.ZITADEL_API_URL ?? "",
    token: "",
  };

  const authserver: ZitadelServer = initializeServer(zitadelConfig);
  console.log("server", authserver);
  const userservice = user.getUser(server);
  return userservice.registerPasskey({
    userId,
    code,
    //   returnCode: new ReturnPasskeyRegistrationCode(),
  });
}

export { server };
