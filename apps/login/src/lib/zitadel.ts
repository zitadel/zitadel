import {
  createIdpServiceClient,
  createOIDCServiceClient,
  createOrganizationServiceClient,
  createSessionServiceClient,
  createSettingsServiceClient,
  createUserServiceClient,
  makeReqCtx,
} from "@zitadel/client/v2";
import { createServerTransport } from "@zitadel/node";
import { RequestChallenges } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { Checks } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import {
  AddHumanUserRequest,
  RetrieveIdentityProviderIntentRequest,
  SetPasswordRequest,
  SetPasswordRequestSchema,
  VerifyPasskeyRegistrationRequest,
  VerifyU2FRegistrationRequest,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";

import { create, Duration } from "@zitadel/client";
import { TextQueryMethod } from "@zitadel/proto/zitadel/object/v2/object_pb";
import { CreateCallbackRequest } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import type { RedirectURLsJson } from "@zitadel/proto/zitadel/user/v2/idp_pb";
import {
  NotificationType,
  SendPasswordResetLinkSchema,
} from "@zitadel/proto/zitadel/user/v2/password_pb";
import {
  SearchQuery,
  SearchQuerySchema,
} from "@zitadel/proto/zitadel/user/v2/query_pb";
import {
  SendInviteCodeSchema,
  User,
  UserState,
} from "@zitadel/proto/zitadel/user/v2/user_pb";
import { unstable_cacheLife as cacheLife } from "next/cache";

const transport = createServerTransport(
  process.env.ZITADEL_SERVICE_USER_TOKEN!,
  { baseUrl: process.env.ZITADEL_API_URL! },
);

export const sessionService = createSessionServiceClient(transport);
export const userService = createUserServiceClient(transport);
export const oidcService = createOIDCServiceClient(transport);
export const idpService = createIdpServiceClient(transport);
export const orgService = createOrganizationServiceClient(transport);
export const settingsService = createSettingsServiceClient(transport);

const useCache = process.env.DEBUG !== "true";

async function cacheWrapper<T>(callback: Promise<T>) {
  "use cache";
  cacheLife("hours");

  return callback;
}

export async function getBrandingSettings(organization?: string) {
  const callback = settingsService
    .getBrandingSettings({ ctx: makeReqCtx(organization) }, {})
    .then((resp) => (resp.settings ? resp.settings : undefined));

  return useCache ? cacheWrapper(callback) : callback;
}

export async function getLoginSettings(orgId?: string) {
  const callback = settingsService
    .getLoginSettings({ ctx: makeReqCtx(orgId) }, {})
    .then((resp) => (resp.settings ? resp.settings : undefined));

  return useCache ? cacheWrapper(callback) : callback;
}

export async function listIDPLinks(userId: string) {
  return userService.listIDPLinks(
    {
      userId,
    },
    {},
  );
}

export async function addOTPEmail(userId: string) {
  return userService.addOTPEmail(
    {
      userId,
    },
    {},
  );
}

export async function addOTPSMS(userId: string) {
  return userService.addOTPSMS({ userId }, {});
}

export async function registerTOTP(userId: string) {
  return userService.registerTOTP({ userId }, {});
}

export async function getGeneralSettings() {
  const callback = settingsService
    .getGeneralSettings({}, {})
    .then((resp) => resp.supportedLanguages);

  return useCache ? cacheWrapper(callback) : callback;
}

export async function getLegalAndSupportSettings(organization?: string) {
  const callback = settingsService
    .getLegalAndSupportSettings({ ctx: makeReqCtx(organization) }, {})
    .then((resp) => (resp.settings ? resp.settings : undefined));

  return useCache ? cacheWrapper(callback) : callback;
}

export async function getPasswordComplexitySettings(organization?: string) {
  const callback = settingsService
    .getPasswordComplexitySettings({ ctx: makeReqCtx(organization) })
    .then((resp) => (resp.settings ? resp.settings : undefined));

  return useCache ? cacheWrapper(callback) : callback;
}

export async function createSessionFromChecks(
  checks: Checks,
  challenges: RequestChallenges | undefined,
  lifetime?: Duration,
) {
  return sessionService.createSession(
    {
      checks: checks,
      challenges,
      lifetime,
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
  lifetime?: Duration,
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
    lifetime,
  });
}

export async function setSession(
  sessionId: string,
  sessionToken: string,
  challenges: RequestChallenges | undefined,
  checks?: Checks,
  lifetime?: Duration,
) {
  return sessionService.setSession(
    {
      sessionId,
      sessionToken,
      challenges,
      checks: checks ? checks : {},
      metadata: {},
      lifetime,
    },
    {},
  );
}

export async function getSession({
  sessionId,
  sessionToken,
}: {
  sessionId: string;
  sessionToken: string;
}) {
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
    email: {
      email,
      verification: {
        case: "isVerified",
        value: false,
      },
    },
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

export async function addHuman(request: AddHumanUserRequest) {
  return userService.addHumanUser(request);
}

export async function verifyTOTPRegistration(code: string, userId: string) {
  return userService.verifyTOTPRegistration({ code, userId }, {});
}

export async function getUserByID(userId: string) {
  return userService.getUserByID({ userId }, {});
}

export async function verifyInviteCode(
  userId: string,
  verificationCode: string,
) {
  return userService.verifyInviteCode({ userId, verificationCode }, {});
}

export async function resendInviteCode(userId: string) {
  return userService.resendInviteCode({ userId }, {});
}

export async function createInviteCode(userId: string, host: string | null) {
  let medium = create(SendInviteCodeSchema, {
    applicationName: "Typescript Login",
  });

  if (host) {
    medium = {
      ...medium,
      urlTemplate: `${host.includes("localhost") ? "http://" : "https://"}${host}/verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}&invite=true`,
    };
  }

  return userService.createInviteCode(
    {
      userId,
      verification: {
        case: "sendCode",
        value: medium,
      },
    },
    {},
  );
}

export async function listUsers({
  loginName,
  userName,
  email,
  organizationId,
}: {
  loginName?: string;
  userName?: string;
  email?: string;
  organizationId?: string;
}) {
  const queries: SearchQuery[] = [];

  if (loginName) {
    queries.push(
      create(SearchQuerySchema, {
        query: {
          case: "loginNameQuery",
          value: {
            loginName: loginName,
            method: TextQueryMethod.EQUALS,
          },
        },
      }),
    );
  }

  if (userName) {
    queries.push(
      create(SearchQuerySchema, {
        query: {
          case: "userNameQuery",
          value: {
            userName: userName,
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

export async function getDefaultOrg(): Promise<Organization | null> {
  return orgService
    .listOrganizations(
      {
        queries: [
          {
            query: {
              case: "defaultQuery",
              value: {},
            },
          },
        ],
      },
      {},
    )
    .then((resp) => (resp?.result && resp.result[0] ? resp.result[0] : null));
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

/**
 *
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */
export async function passwordReset(
  userId: string,
  host: string | null,
  authRequestId?: string,
) {
  let medium = create(SendPasswordResetLinkSchema, {
    notificationType: NotificationType.Email,
  });

  if (host) {
    medium = {
      ...medium,
      urlTemplate:
        `${host.includes("localhost") ? "http://" : "https://"}${host}/password/set?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}` +
        (authRequestId ? `&authRequestId=${authRequestId}` : ""),
    };
  }

  return userService.passwordReset(
    {
      userId,
      medium: {
        case: "sendLink",
        value: medium,
      },
    },
    {},
  );
}

/**
 *
 * @param userId userId of the user to set the password for
 * @param password the new password
 * @param code optional if the password should be set with a code (reset), no code for initial setup of password
 * @returns
 */
export async function setUserPassword(
  userId: string,
  password: string,
  user: User,
  code?: string,
) {
  let payload = create(SetPasswordRequestSchema, {
    userId,
    newPassword: {
      password,
    },
  });

  // check if the user has no password set in order to set a password
  if (!code) {
    const authmethods = await listAuthenticationMethodTypes(userId);

    // if the user has no authmethods set, we can set a password otherwise we need a code
    if (
      !(authmethods.authMethodTypes.length === 0) &&
      user.state !== UserState.INITIAL
    ) {
      return { error: "Provide a code to set a password" };
    }
  }

  if (code) {
    payload = {
      ...payload,
      verification: {
        case: "verificationCode",
        value: code,
      },
    };
  }

  return userService.setPassword(payload, {}).catch((error) => {
    // throw error if failed precondition (ex. User is not yet initialized)
    if (error.code === 9 && error.message) {
      return { error: error.message };
    } else {
      throw error;
    }
  });
}

export async function setPassword(payload: SetPasswordRequest) {
  return userService.setPassword(payload, {});
}

/**
 *
 * @param server
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */

// TODO check for token requirements!
export async function createPasskeyRegistrationLink(
  userId: string,
  // token: string,
) {
  // const transport = createServerTransport(token, {
  //   baseUrl: process.env.ZITADEL_API_URL!,
  // });

  // const service = createUserServiceClient(transport);
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
