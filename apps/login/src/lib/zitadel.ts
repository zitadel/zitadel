import {
  createOIDCServiceClient,
  createSessionServiceClient,
  createSettingsServiceClient,
  createUserServiceClient,
  makeReqCtx,
} from "@zitadel/client/v2";
import { createManagementServiceClient } from "@zitadel/client/v1";
import { createServerTransport } from "@zitadel/node";
import { GetActiveIdentityProvidersRequest } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb";
import { Checks } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { RequestChallenges } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import {
  RetrieveIdentityProviderIntentRequest,
  VerifyU2FRegistrationRequest,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";

import { CreateCallbackRequest } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { TextQueryMethod } from "@zitadel/proto/zitadel/object/v2/object_pb";
import type { RedirectURLs } from "@zitadel/proto/zitadel/user/v2/idp_pb";
import { ProviderSlug } from "./demos";
import { PlainMessage } from "@zitadel/client";

const SESSION_LIFETIME_S = 3000;

const transport = createServerTransport(
  process.env.ZITADEL_SERVICE_USER_TOKEN!,
  {
    baseUrl: process.env.ZITADEL_API_URL!,
    httpVersion: "2",
  },
);

export const sessionService = createSessionServiceClient(transport);
export const managementService = createManagementServiceClient(transport);
export const userService = createUserServiceClient(transport);
export const oidcService = createOIDCServiceClient(transport);
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
  checks: PlainMessage<Checks>,
  challenges: PlainMessage<RequestChallenges> | undefined,
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
  checks?: PlainMessage<Checks>,
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

export async function verifyTOTPRegistration(
  code: string,
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
  return userService.verifyTOTPRegistration({ code, userId }, {});
}

export async function getUserByID(userId: string) {
  return userService.getUserByID({ userId }, {});
}

export async function listUsers(userName: string, organizationId: string) {
  return userService.listUsers(
    {
      queries: organizationId
        ? [
            {
              query: {
                case: "userNameQuery",
                value: {
                  userName,
                  method: TextQueryMethod.EQUALS,
                },
              },
            },
            {
              query: {
                case: "organizationIdQuery",
                value: {
                  organizationId,
                },
              },
            },
          ]
        : [
            {
              query: {
                case: "userNameQuery",
                value: {
                  userName,
                  method: TextQueryMethod.EQUALS,
                },
              },
            },
          ],
    },
    {},
  );
}

export async function getOrgByDomain(domain: string) {
  return managementService.getOrgByDomainGlobal({ domain }, {});
}

export const PROVIDER_NAME_MAPPING: {
  [provider: string]: string;
} = {
  [ProviderSlug.GOOGLE]: "Google",
  [ProviderSlug.GITHUB]: "GitHub",
};

export async function startIdentityProviderFlow({
  idpId,
  urls,
}: {
  idpId: string;
  urls: PlainMessage<RedirectURLs>;
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

export async function createCallback(req: PlainMessage<CreateCallbackRequest>) {
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
  request: PlainMessage<VerifyU2FRegistrationRequest>,
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
  passkeyId: string,
  passkeyName: string,
  publicKeyCredential:
    | {
        [key: string]: any;
      }
    | undefined,
  userId: string,
) {
  return userService.verifyPasskeyRegistration(
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
