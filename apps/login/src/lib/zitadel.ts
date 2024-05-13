import { VerifyU2FRegistrationRequest } from "@zitadel/server";
import {
  GetUserByIDResponse,
  RegisterTOTPResponse,
  VerifyTOTPRegistrationResponse,
} from "@zitadel/server";
import {
  LegalAndSupportSettings,
  PasswordComplexitySettings,
  ZitadelServer,
  VerifyMyAuthFactorOTPResponse,
  ZitadelServerOptions,
  user,
  oidc,
  settings,
  getServers,
  auth,
  initializeServer,
  session,
  GetGeneralSettingsResponse,
  CreateSessionResponse,
  GetBrandingSettingsResponse,
  GetPasswordComplexitySettingsResponse,
  RegisterU2FResponse,
  GetLegalAndSupportSettingsResponse,
  AddHumanUserResponse,
  BrandingSettings,
  ListSessionsResponse,
  GetSessionResponse,
  VerifyEmailResponse,
  Checks,
  SetSessionResponse,
  SetSessionRequest,
  ListUsersResponse,
  management,
  DeleteSessionResponse,
  VerifyPasskeyRegistrationResponse,
  LoginSettings,
  GetOrgByDomainGlobalResponse,
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
  TextQueryMethod,
  ListHumanAuthFactorsResponse,
  AddHumanUserRequest,
  AddOTPEmailResponse,
  AddOTPSMSResponse,
} from "@zitadel/server";

const SESSION_LIFETIME_S = 3000;

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
  server: ZitadelServer,
  organization?: string,
): Promise<BrandingSettings | undefined> {
  const settingsService = settings.getSettings(server);
  return settingsService
    .getBrandingSettings(
      { ctx: organization ? { orgId: organization } : { instance: true } },
      {},
    )
    .then((resp: GetBrandingSettingsResponse) => resp.settings);
}

export async function getLoginSettings(
  server: ZitadelServer,
  orgId?: string,
): Promise<LoginSettings | undefined> {
  const settingsService = settings.getSettings(server);
  return settingsService
    .getLoginSettings({ ctx: orgId ? { orgId } : { instance: true } }, {})
    .then((resp: GetLoginSettingsResponse) => resp.settings);
}

export async function verifyMyAuthFactorOTP(
  code: string,
): Promise<VerifyMyAuthFactorOTPResponse> {
  const authService = auth.getAuth(server);
  return authService.verifyMyAuthFactorOTP({ code }, {});
}

export async function addOTPEmail(
  userId: string,
): Promise<AddOTPEmailResponse | undefined> {
  const userService = user.getUser(server);
  return userService.addOTPEmail(
    {
      userId,
    },
    {},
  );
}

export async function addOTPSMS(
  userId: string,
  token?: string,
): Promise<AddOTPSMSResponse | undefined> {
  let userService;
  if (token) {
    const authConfig: ZitadelServerOptions = {
      name: "zitadel login",
      apiUrl: process.env.ZITADEL_API_URL ?? "",
      token: token,
    };

    const sessionUser = initializeServer(authConfig);
    userService = user.getUser(sessionUser);
  } else {
    userService = user.getUser(server);
  }
  return userService.addOTPSMS({ userId }, {});
}

export async function registerTOTP(
  userId: string,
  token?: string,
): Promise<RegisterTOTPResponse | undefined> {
  let userService;
  if (token) {
    const authConfig: ZitadelServerOptions = {
      name: "zitadel login",
      apiUrl: process.env.ZITADEL_API_URL ?? "",
      token: token,
    };

    const sessionUser = initializeServer(authConfig);
    userService = user.getUser(sessionUser);
  } else {
    userService = user.getUser(server);
  }
  return userService.registerTOTP({ userId }, {});
}

export async function getGeneralSettings(
  server: ZitadelServer,
): Promise<string[] | undefined> {
  const settingsService = settings.getSettings(server);
  return settingsService
    .getGeneralSettings({}, {})
    .then((resp: GetGeneralSettingsResponse) => resp.supportedLanguages);
}

export async function getLegalAndSupportSettings(
  server: ZitadelServer,
  organization?: string,
): Promise<LegalAndSupportSettings | undefined> {
  const settingsService = settings.getSettings(server);
  return settingsService
    .getLegalAndSupportSettings(
      { ctx: organization ? { orgId: organization } : { instance: true } },
      {},
    )
    .then((resp: GetLegalAndSupportSettingsResponse) => {
      return resp.settings;
    });
}

export async function getPasswordComplexitySettings(
  server: ZitadelServer,
  organization?: string,
): Promise<PasswordComplexitySettings | undefined> {
  const settingsService = settings.getSettings(server);

  return settingsService
    .getPasswordComplexitySettings(
      organization
        ? { ctx: { orgId: organization } }
        : { ctx: { instance: true } },
      {},
    )
    .then((resp: GetPasswordComplexitySettingsResponse) => resp.settings);
}

export async function createSessionFromChecks(
  server: ZitadelServer,
  checks: Checks,
  challenges: RequestChallenges | undefined,
): Promise<CreateSessionResponse | undefined> {
  const sessionService = session.getSession(server);
  return sessionService.createSession(
    {
      checks: checks,
      challenges,
      lifetime: {
        seconds: SESSION_LIFETIME_S,
        nanos: 0,
      },
    },
    {},
  );
}

export async function createSessionForUserIdAndIdpIntent(
  server: ZitadelServer,
  userId: string,
  idpIntent: {
    idpIntentId?: string | undefined;
    idpIntentToken?: string | undefined;
  },
): Promise<CreateSessionResponse | undefined> {
  const sessionService = session.getSession(server);

  return sessionService.createSession(
    {
      checks: { user: { userId }, idpIntent },
      // lifetime: {
      //   seconds: 300,
      //   nanos: 0,
      // },
    },
    {},
  );
}

export async function setSession(
  server: ZitadelServer,
  sessionId: string,
  sessionToken: string,
  challenges: RequestChallenges | undefined,
  checks: Checks,
): Promise<SetSessionResponse | undefined> {
  const sessionService = session.getSession(server);

  const payload: SetSessionRequest = {
    sessionId,
    sessionToken,
    challenges,
    checks: {},
    metadata: {},
  };

  if (checks && payload.checks) {
    payload.checks = checks;
  }

  return sessionService.setSession(payload, {});
}

export async function getSession(
  server: ZitadelServer,
  sessionId: string,
  sessionToken: string,
): Promise<GetSessionResponse | undefined> {
  const sessionService = session.getSession(server);
  return sessionService.getSession({ sessionId, sessionToken }, {});
}

export async function deleteSession(
  server: ZitadelServer,
  sessionId: string,
  sessionToken: string,
): Promise<DeleteSessionResponse | undefined> {
  const sessionService = session.getSession(server);
  return sessionService.deleteSession({ sessionId, sessionToken }, {});
}

export async function listSessions(
  server: ZitadelServer,
  ids: string[],
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
  organization: string | undefined;
};

export async function addHumanUser(
  server: ZitadelServer,
  { email, firstName, lastName, password, organization }: AddHumanUserData,
): Promise<AddHumanUserResponse> {
  const userService = user.getUser(server);

  const payload: Partial<AddHumanUserRequest> = {
    email: { email },
    username: email,
    profile: { givenName: firstName, familyName: lastName },
  };

  if (organization) {
    payload.organization = { orgId: organization };
  }

  return userService.addHumanUser(
    password
      ? {
          ...payload,
          password: { password },
        }
      : payload,
    {},
  );
}

export async function verifyTOTPRegistration(
  code: string,
  userId: string,
  token?: string,
): Promise<VerifyTOTPRegistrationResponse> {
  let userService;
  if (token) {
    const authConfig: ZitadelServerOptions = {
      name: "zitadel login",
      apiUrl: process.env.ZITADEL_API_URL ?? "",
      token: token,
    };

    const sessionUser = initializeServer(authConfig);
    userService = user.getUser(sessionUser);
  } else {
    userService = user.getUser(server);
  }
  return userService.verifyTOTPRegistration({ code, userId }, {});
}

export async function getUserByID(
  userId: string,
): Promise<GetUserByIDResponse> {
  const userService = user.getUser(server);

  return userService.getUserByID({ userId }, {});
}

export async function listUsers(
  userName: string,
  organizationId: string,
): Promise<ListUsersResponse> {
  const userService = user.getUser(server);

  return userService.listUsers(
    {
      queries: organizationId
        ? [
            {
              userNameQuery: {
                userName,
                method: TextQueryMethod.TEXT_QUERY_METHOD_EQUALS,
              },
            },
            {
              organizationIdQuery: {
                organizationId,
              },
            },
          ]
        : [
            {
              userNameQuery: {
                userName,
                method: TextQueryMethod.TEXT_QUERY_METHOD_EQUALS,
              },
            },
          ],
    },
    {},
  );
}

export async function getOrgByDomain(
  domain: string,
): Promise<GetOrgByDomainGlobalResponse> {
  const mgmtService = management.getManagement(server);
  return mgmtService.getOrgByDomainGlobal({ domain }, {});
}

export async function startIdentityProviderFlow(
  server: ZitadelServer,
  { idpId, urls }: StartIdentityProviderIntentRequest,
): Promise<StartIdentityProviderIntentResponse> {
  const userService = user.getUser(server);

  return userService.startIdentityProviderIntent({
    idpId,
    urls,
  });
}

export async function retrieveIdentityProviderInformation(
  server: ZitadelServer,
  { idpIntentId, idpIntentToken }: RetrieveIdentityProviderIntentRequest,
): Promise<RetrieveIdentityProviderIntentResponse> {
  const userService = user.getUser(server);

  return userService.retrieveIdentityProviderIntent({
    idpIntentId,
    idpIntentToken,
  });
}

export async function getAuthRequest(
  server: ZitadelServer,
  { authRequestId }: GetAuthRequestRequest,
): Promise<GetAuthRequestResponse> {
  const oidcService = oidc.getOidc(server);

  return oidcService.getAuthRequest({
    authRequestId,
  });
}

export async function createCallback(
  server: ZitadelServer,
  req: CreateCallbackRequest,
): Promise<CreateCallbackResponse> {
  const oidcService = oidc.getOidc(server);

  return oidcService.createCallback(req);
}

export async function verifyEmail(
  server: ZitadelServer,
  userId: string,
  verificationCode: string,
): Promise<VerifyEmailResponse> {
  const userservice = user.getUser(server);
  return userservice.verifyEmail(
    {
      userId,
      verificationCode,
    },
    {},
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
  userId: string,
): Promise<any> {
  const userservice = user.getUser(server);
  return userservice.setEmail(
    {
      userId,
    },
    {},
  );
}

/**
 *
 * @param server
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */
export async function createPasskeyRegistrationLink(
  userId: string,
  token?: string,
): Promise<any> {
  let userService;
  if (token) {
    const authConfig: ZitadelServerOptions = {
      name: "zitadel login",
      apiUrl: process.env.ZITADEL_API_URL ?? "",
      token: token,
    };

    const sessionUser = initializeServer(authConfig);
    userService = user.getUser(sessionUser);
  } else {
    userService = user.getUser(server);
  }

  return userService.createPasskeyRegistrationLink({
    userId,
    returnCode: {},
  });
}

/**
 *
 * @param server
 * @param userId the id of the user where the email should be set
 * @param domain the domain on which the factor is registered
 * @returns the newly set email
 */
export async function registerU2F(
  userId: string,
  domain: string,
): Promise<RegisterU2FResponse> {
  const userservice = user.getUser(server);

  return userservice.registerU2F({
    userId,
    domain,
  });
}

/**
 *
 * @param server
 * @param userId the id of the user where the email should be set
 * @param domain the domain on which the factor is registered
 * @returns the newly set email
 */
export async function verifyU2FRegistration(
  request: VerifyU2FRegistrationRequest,
): Promise<any> {
  const userservice = user.getUser(server);

  return userservice.verifyU2FRegistration(request, {});
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
  userId: string,
): Promise<VerifyPasskeyRegistrationResponse> {
  const userservice = user.getUser(server);
  return userservice.verifyPasskeyRegistration(
    {
      passkeyId,
      passkeyName,
      publicKeyCredential,
      userId,
    },
    {},
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
  domain: string,
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
  userId: string,
): Promise<ListAuthenticationMethodTypesResponse> {
  const userservice = user.getUser(server);
  return userservice.listAuthenticationMethodTypes({
    userId,
  });
}

export { server };
