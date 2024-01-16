import {
  ZitadelServer,
  ZitadelServerOptions,
  user,
  oidc,
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
  SetSessionRequest,
  DeleteSessionResponse,
  VerifyPasskeyRegistrationResponse,
  LoginSettings,
  GetLoginSettingsResponse,
  ListAuthenticationMethodTypesResponse,
  StartIdentityProviderIntentRequest,
  StartIdentityProviderIntentResponse,
  RetrieveIdentityProviderIntentRequest,
  RetrieveIdentityProviderIntentResponse,
  GetAuthRequestResponse,
  GetAuthRequestRequest,
  CreateCallbackRequest,
  CreateCallbackResponse,
  RequestChallenges,
  AddHumanUserRequest,
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

export async function getBrandingSettings(
  server: ZitadelServer
): Promise<BrandingSettings | undefined> {
  const settingsService = settings.getSettings(server);
  return settingsService
    .getBrandingSettings({}, {})
    .then((resp: GetBrandingSettingsResponse) => resp.settings);
}

export async function getLoginSettings(
  server: ZitadelServer
): Promise<LoginSettings | undefined> {
  const settingsService = settings.getSettings(server);
  return settingsService
    .getLoginSettings({}, {})
    .then((resp: GetLoginSettingsResponse) => resp.settings);
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
  loginName: string,
  password: string | undefined,
  challenges: RequestChallenges | undefined
): Promise<CreateSessionResponse | undefined> {
  const sessionService = session.getSession(server);
  return password
    ? sessionService.createSession(
        {
          checks: { user: { loginName }, password: { password } },
          challenges,
        },
        {}
      )
    : sessionService.createSession(
        { checks: { user: { loginName } }, challenges },

        {}
      );
}

export async function setSession(
  server: ZitadelServer,
  sessionId: string,
  sessionToken: string,
  password: string | undefined,
  webAuthN: { credentialAssertionData: any } | undefined,
  challenges: RequestChallenges | undefined
): Promise<SetSessionResponse | undefined> {
  const sessionService = session.getSession(server);

  const payload: SetSessionRequest = {
    sessionId,
    sessionToken,
    challenges,
    checks: {},
    metadata: {},
  };

  if (password && payload.checks) {
    payload.checks.password = { password };
  }

  if (webAuthN && payload.checks) {
    payload.checks.webAuthN = webAuthN;
  }

  return sessionService.setSession(payload, {});
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
  password: string | undefined;
};

export async function addHumanUser(
  server: ZitadelServer,
  { email, firstName, lastName, password }: AddHumanUserData
): Promise<string> {
  const userService = user.getUser(server);

  const payload: Partial<AddHumanUserRequest> = {
    email: { email },
    username: email,
    profile: { givenName: firstName, familyName: lastName },
  };
  return userService
    .addHumanUser(
      password
        ? {
            ...payload,
            password: { password },
          }
        : payload,
      {}
    )
    .then((resp: AddHumanUserResponse) => {
      return resp.userId;
    });
}

export async function startIdentityProviderFlow(
  server: ZitadelServer,
  { idpId, urls }: StartIdentityProviderIntentRequest
): Promise<StartIdentityProviderIntentResponse> {
  const userService = user.getUser(server);

  return userService.startIdentityProviderIntent({
    idpId,
    urls,
  });
}

export async function retrieveIdentityProviderInformation(
  server: ZitadelServer,
  { idpIntentId, idpIntentToken }: RetrieveIdentityProviderIntentRequest
): Promise<RetrieveIdentityProviderIntentResponse> {
  const userService = user.getUser(server);

  return userService.retrieveIdentityProviderIntent({
    idpIntentId,
    idpIntentToken,
  });
}

export async function getAuthRequest(
  server: ZitadelServer,
  { authRequestId }: GetAuthRequestRequest
): Promise<GetAuthRequestResponse> {
  const oidcService = oidc.getOidc(server);

  return oidcService.getAuthRequest({
    authRequestId,
  });
}

export async function createCallback(
  server: ZitadelServer,
  req: CreateCallbackRequest
): Promise<CreateCallbackResponse> {
  const oidcService = oidc.getOidc(server);

  return oidcService.createCallback(req);
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

/**
 *
 * @param server
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */
export async function createPasskeyRegistrationLink(
  userId: string
): Promise<any> {
  const userservice = user.getUser(server);

  return userservice.createPasskeyRegistrationLink({
    userId,
    returnCode: {},
  });
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
  publicKeyCredential:
    | {
        [key: string]: any;
      }
    | undefined,
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
  code: { id: string; code: string },
  domain: string
): Promise<any> {
  const userservice = user.getUser(server);
  return userservice.registerPasskey({
    userId,
    code,
    domain,
    // authenticator:
  });
}

/**
 *
 * @param server
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */
export async function listAuthenticationMethodTypes(
  userId: string
): Promise<ListAuthenticationMethodTypesResponse> {
  const userservice = user.getUser(server);
  return userservice.listAuthenticationMethodTypes({
    userId,
  });
}

export { server };
