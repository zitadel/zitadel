import {
  createOIDCServiceClient,
  createSessionServiceClient,
  createSettingsServiceClient,
  createUserServiceClient,
  createIdpServiceClient,
  makeReqCtx,
  createOrganizationServiceClient,
} from "@zitadel/client/v2";
import { createManagementServiceClient } from "@zitadel/client/v1";
import { createServerTransport } from "@zitadel/node";
import { Checks } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { RequestChallenges } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import {
  RetrieveIdentityProviderIntentRequest,
  VerifyPasskeyRegistrationRequest,
  VerifyU2FRegistrationRequest,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { IDPInformation } from "@zitadel/proto/zitadel/user/v2/idp_pb";

import { CreateCallbackRequest } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import type { RedirectURLsJson } from "@zitadel/proto/zitadel/user/v2/idp_pb";
import {
  SearchQuery,
  SearchQuerySchema,
} from "@zitadel/proto/zitadel/user/v2/query_pb";
import { PROVIDER_MAPPING } from "./idp";
import { create } from "@zitadel/client";
import { TextQueryMethod } from "@zitadel/proto/zitadel/object/v2/object_pb";

const SESSION_LIFETIME_S = 3000;

const transport = createServerTransport(
  process.env.ZITADEL_SERVICE_USER_TOKEN!,
  {
    baseUrl: process.env.ZITADEL_API_URL!,
    httpVersion: "2",
  },
);

export const sessionService = createSessionServiceClient(transport);
export const userService = createUserServiceClient(transport);
export const oidcService = createOIDCServiceClient(transport);
export const idpService = createIdpServiceClient(transport);
export const orgService = createOrganizationServiceClient(transport);

export const settingsService = createSettingsServiceClient(transport);

export async function getBrandingSettings(organization?: string) {
  return settingsService
    .getBrandingSettings({ ctx: makeReqCtx(organization) }, {})
    .then((resp) => resp.settings);
}

export async function getLoginSettings(orgId?: string) {
  return settingsService
    .getLoginSettings({ ctx: makeReqCtx(orgId) }, {})
    .then((resp) => resp.settings);
}

export async function addOTPEmail(userId: string) {
  return userService.addOTPEmail(
    {
      userId,
    },
    {},
  );
}

export async function addOTPSMS(userId: string, token?: string) {
  // TODO: Follow up here, I do not understand the branching
  // let userService;
  // if (token) {
  //   const authConfig: ZitadelServerOptions = {
  //     name: "zitadel login",
  //     apiUrl: process.env.ZITADEL_API_URL ?? "",
  //     token: token,
  //   };
  //   const sessionUser = initializeServer(authConfig);
  //   userService = user.getUser(sessionUser);
  // } else {
  //   userService = user.getUser(server);
  // }

  return userService.addOTPSMS({ userId }, {});
}

export async function registerTOTP(userId: string, token?: string) {
  // TODO: Follow up here, I do not understand the branching
  // let userService;
  // if (token) {
  //   const authConfig: ZitadelServerOptions = {
  //     name: "zitadel login",
  //     apiUrl: process.env.ZITADEL_API_URL ?? "",
  //     token: token,
  //   };
  //
  //   const sessionUser = initializeServer(authConfig);
  //   userService = user.getUser(sessionUser);
  // } else {
  //   userService = user.getUser(server);
  // }
  return userService.registerTOTP({ userId }, {});
}

export async function getGeneralSettings() {
  return settingsService
    .getGeneralSettings({}, {})
    .then((resp) => resp.supportedLanguages);
}

export async function getLegalAndSupportSettings(organization?: string) {
  return settingsService
    .getLegalAndSupportSettings({ ctx: makeReqCtx(organization) }, {})
    .then((resp) => {
      return resp.settings;
    });
}

export async function getPasswordComplexitySettings(organization?: string) {
  return settingsService
    .getPasswordComplexitySettings({ ctx: makeReqCtx(organization) })
    .then((resp) => resp.settings);
}

export async function createSessionFromChecks(
  checks: Checks,
  challenges: RequestChallenges | undefined,
) {
  return sessionService.createSession(
    {
      checks: checks,
      challenges,
      lifetime: {
        seconds: BigInt(SESSION_LIFETIME_S),
        nanos: 0,
      },
    },
    {},
  );
}

export async function createSessionForUserIdAndIdpIntent(
  userId: string,
  idpIntent: {
    idpIntentId?: string | undefined;
    idpIntentToken?: string | undefined;
  },
) {
  return sessionService.createSession({
    checks: {
      user: {
        search: {
          case: "userId",
          value: userId,
        },
      },
      idpIntent,
    },
    // lifetime: {
    //   seconds: 300,
    //   nanos: 0,
    // },
  });
}

export async function setSession(
  sessionId: string,
  sessionToken: string,
  challenges: RequestChallenges | undefined,
  checks?: Checks,
) {
  return sessionService.setSession(
    {
      sessionId,
      sessionToken,
      challenges,
      checks: checks ? checks : {},
      metadata: {},
    },
    {},
  );
}

export async function getSession(sessionId: string, sessionToken: string) {
  return sessionService.getSession({ sessionId, sessionToken }, {});
}

export async function deleteSession(sessionId: string, sessionToken: string) {
  return sessionService.deleteSession({ sessionId, sessionToken }, {});
}

export async function listSessions(ids: string[]) {
  return sessionService.listSessions(
    {
      queries: [
        {
          query: {
            case: "idsQuery",
            value: { ids: ids },
          },
        },
      ],
    },
    {},
  );
}

export type AddHumanUserData = {
  firstName: string;
  lastName: string;
  email: string;
  password: string | undefined;
  organization: string | undefined;
};

export async function addHumanUser({
  email,
  firstName,
  lastName,
  password,
  organization,
}: AddHumanUserData) {
  return userService.addHumanUser({
    email: { email },
    username: email,
    profile: { givenName: firstName, familyName: lastName },
    organization: organization
      ? { org: { case: "orgId", value: organization } }
      : undefined,
    passwordType: password
      ? { case: "password", value: { password: password } }
      : undefined,
  });
}

export async function verifyTOTPRegistration(code: string, userId: string) {
  return userService.verifyTOTPRegistration({ code, userId }, {});
}

export async function getUserByID(userId: string) {
  return userService.getUserByID({ userId }, {});
}

export async function listUsers({
  userName,
  email,
  organizationId,
}: {
  userName?: string;
  email?: string;
  organizationId?: string;
}) {
  const queries: SearchQuery[] = [];

  if (userName) {
    queries.push(
      create(SearchQuerySchema, {
        query: {
          case: "userNameQuery",
          value: {
            userName,
            method: TextQueryMethod.EQUALS,
          },
        },
      }),
    );
  }

  if (organizationId) {
    queries.push(
      create(SearchQuerySchema, {
        query: {
          case: "organizationIdQuery",
          value: {
            organizationId,
          },
        },
      }),
    );
  }

  if (email) {
    queries.push(
      create(SearchQuerySchema, {
        query: {
          case: "emailQuery",
          value: {
            emailAddress: email,
          },
        },
      }),
    );
  }

  return userService.listUsers({ queries: queries });
}

export async function getOrgsByDomain(domain: string) {
  return orgService.listOrganizations(
    {
      queries: [
        {
          query: {
            case: "domainQuery",
            value: { domain, method: TextQueryMethod.EQUALS },
          },
        },
      ],
    },
    {},
  );
}

export async function startIdentityProviderFlow({
  idpId,
  urls,
}: {
  idpId: string;
  urls: RedirectURLsJson;
}) {
  return userService.startIdentityProviderIntent({
    idpId,
    content: {
      case: "urls",
      value: urls,
    },
  });
}

export async function retrieveIdentityProviderInformation({
  idpIntentId,
  idpIntentToken,
}: RetrieveIdentityProviderIntentRequest) {
  return userService.retrieveIdentityProviderIntent({
    idpIntentId,
    idpIntentToken,
  });
}

export async function getAuthRequest({
  authRequestId,
}: {
  authRequestId: string;
}) {
  return oidcService.getAuthRequest({
    authRequestId,
  });
}

export async function createCallback(req: CreateCallbackRequest) {
  return oidcService.createCallback(req);
}

export async function verifyEmail(userId: string, verificationCode: string) {
  return userService.verifyEmail(
    {
      userId,
      verificationCode,
    },
    {},
  );
}

/**
 *
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */
export async function resendEmailCode(userId: string) {
  return userService.resendEmailCode(
    {
      userId,
    },
    {},
  );
}

export function retrieveIDPIntent(id: string, token: string) {
  return userService.retrieveIdentityProviderIntent(
    { idpIntentId: id, idpIntentToken: token },
    {},
  );
}

export function getIDPByID(id: string) {
  return idpService.getIDPByID({ id }, {}).then((resp) => resp.idp);
}

export function addIDPLink(
  idp: {
    id: string;
    userId: string;
    userName: string;
  },
  userId: string,
) {
  return userService.addIDPLink(
    {
      idpLink: {
        userId: idp.userId,
        idpId: idp.id,
        userName: idp.userName,
      },
      userId,
    },
    {},
  );
}

export function createUser(provider: string, info: IDPInformation) {
  const userData = PROVIDER_MAPPING[provider](info);
  console.log("ud", userData);
  return userService.addHumanUser(userData, {});
}

/**
 *
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */
export async function passwordReset(userId: string): Promise<any> {
  return userService.passwordReset(
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
) {
  // let userService;
  // if (token) {
  //   const authConfig: ZitadelServerOptions = {
  //     name: "zitadel login",
  //     apiUrl: process.env.ZITADEL_API_URL ?? "",
  //     token: token,
  //   };
  //
  //   const sessionUser = initializeServer(authConfig);
  //   userService = user.getUser(sessionUser);
  // } else {
  //   userService = user.getUser(server);
  // }

  return userService.createPasskeyRegistrationLink({
    userId,
    medium: {
      case: "returnCode",
      value: {},
    },
  });
}

/**
 *
 * @param userId the id of the user where the email should be set
 * @param domain the domain on which the factor is registered
 * @returns the newly set email
 */
export async function registerU2F(userId: string, domain: string) {
  return userService.registerU2F({
    userId,
    domain,
  });
}

/**
 *
 * @param userId the id of the user where the email should be set
 * @param domain the domain on which the factor is registered
 * @returns the newly set email
 */
export async function verifyU2FRegistration(
  request: VerifyU2FRegistrationRequest,
) {
  return userService.verifyU2FRegistration(request, {});
}

export async function getActiveIdentityProviders(orgId?: string) {
  return settingsService.getActiveIdentityProviders(
    { ctx: makeReqCtx(orgId) },
    {},
  );
}

/**
 *
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */
export async function verifyPasskeyRegistration(
  request: VerifyPasskeyRegistrationRequest,
) {
  return userService.verifyPasskeyRegistration(request, {});
}

/**
 *
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */
export async function registerPasskey(
  userId: string,
  code: { id: string; code: string },
  domain: string,
) {
  return userService.registerPasskey({
    userId,
    code,
    domain,
    // authenticator:
  });
}

/**
 *
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */
export async function listAuthenticationMethodTypes(userId: string) {
  return userService.listAuthenticationMethodTypes({
    userId,
  });
}
