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
import { IDPInformation } from "@zitadel/proto/zitadel/user/v2/idp_pb";
import {
  RetrieveIdentityProviderIntentRequest,
  VerifyPasskeyRegistrationRequest,
  VerifyU2FRegistrationRequest,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";

import { create, fromJson, toJson } from "@zitadel/client";
import { TextQueryMethod } from "@zitadel/proto/zitadel/object/v2/object_pb";
import { CreateCallbackRequest } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { BrandingSettingsSchema } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { LegalAndSupportSettingsSchema } from "@zitadel/proto/zitadel/settings/v2/legal_settings_pb";
import {
  IdentityProviderType,
  LoginSettingsSchema,
} from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { PasswordComplexitySettingsSchema } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";
import type { RedirectURLsJson } from "@zitadel/proto/zitadel/user/v2/idp_pb";
import {
  SearchQuery,
  SearchQuerySchema,
} from "@zitadel/proto/zitadel/user/v2/query_pb";
import { unstable_cache } from "next/cache";
import { PROVIDER_MAPPING } from "./idp";

const SESSION_LIFETIME_S = 3600; // TODO load from oidc settings
const CACHE_REVALIDATION_INTERVAL_IN_SECONDS = process.env
  .CACHE_REVALIDATION_INTERVAL_IN_SECONDS
  ? Number(process.env.CACHE_REVALIDATION_INTERVAL_IN_SECONDS)
  : 3600;

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
  return unstable_cache(
    async () => {
      return await settingsService
        .getBrandingSettings({ ctx: makeReqCtx(organization) }, {})
        .then((resp) =>
          resp.settings
            ? toJson(BrandingSettingsSchema, resp.settings)
            : undefined,
        );
    },
    ["brandingSettings", organization ?? "default"],
    {
      revalidate: CACHE_REVALIDATION_INTERVAL_IN_SECONDS,
      tags: ["brandingSettings"],
    },
  )().then((resp) =>
    resp ? fromJson(BrandingSettingsSchema, resp) : undefined,
  );
}

export async function getLoginSettings(orgId?: string) {
  return unstable_cache(
    async () => {
      return await settingsService
        .getLoginSettings({ ctx: makeReqCtx(orgId) }, {})
        .then((resp) =>
          resp.settings
            ? toJson(LoginSettingsSchema, resp.settings)
            : undefined,
        );
    },
    ["loginSettings", orgId ?? "default"],
    {
      revalidate: CACHE_REVALIDATION_INTERVAL_IN_SECONDS,
      tags: ["loginSettings"],
    },
  )().then((resp) => (resp ? fromJson(LoginSettingsSchema, resp) : undefined));
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
  return settingsService
    .getGeneralSettings({}, {})
    .then((resp) => resp.supportedLanguages);
}

export async function getLegalAndSupportSettings(organization?: string) {
  return unstable_cache(
    async () => {
      return await settingsService
        .getLegalAndSupportSettings({ ctx: makeReqCtx(organization) }, {})
        .then((resp) =>
          resp.settings
            ? toJson(LegalAndSupportSettingsSchema, resp.settings)
            : undefined,
        );
    },
    ["legalAndSupportSettings", organization ?? "default"],
    {
      revalidate: CACHE_REVALIDATION_INTERVAL_IN_SECONDS,
      tags: ["legalAndSupportSettings"],
    },
  )().then((resp) =>
    resp ? fromJson(LegalAndSupportSettingsSchema, resp) : undefined,
  );
}

export async function getPasswordComplexitySettings(organization?: string) {
  return unstable_cache(
    async () => {
      return await settingsService
        .getPasswordComplexitySettings({ ctx: makeReqCtx(organization) })
        .then((resp) =>
          resp.settings
            ? toJson(PasswordComplexitySettingsSchema, resp.settings)
            : undefined,
        );
    },
    ["complexitySettings", organization ?? "default"],
    {
      revalidate: CACHE_REVALIDATION_INTERVAL_IN_SECONDS,
      tags: ["complexitySettings"],
    },
  )().then((resp) =>
    resp ? fromJson(PasswordComplexitySettingsSchema, resp) : undefined,
  );
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

export function createUser(
  provider: IdentityProviderType,
  info: IDPInformation,
) {
  const userData = PROVIDER_MAPPING[provider](info);
  return userService.addHumanUser(userData, {});
}

/**
 *
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */
export async function passwordReset(userId: string) {
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

// TODO check for token requirements!
export async function createPasskeyRegistrationLink(
  userId: string,
  // token: string,
) {
  // const transport = createServerTransport(token, {
  //   baseUrl: process.env.ZITADEL_API_URL!,
  //   httpVersion: "2",
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

// TODO check for token requirements!
export async function registerU2F(
  userId: string,
  domain: string,
  // token: string,
) {
  // const transport = createServerTransport(token, {
  //   baseUrl: process.env.ZITADEL_API_URL!,
  //   httpVersion: "2",
  // });

  // const service = createUserServiceClient(transport);
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
