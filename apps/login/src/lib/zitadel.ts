import { Client, create, Duration } from "@zitadel/client";
import { createServerTransport as libCreateServerTransport } from "@zitadel/client/node";
import { makeReqCtx } from "@zitadel/client/v2";
import { IdentityProviderService } from "@zitadel/proto/zitadel/idp/v2/idp_service_pb";
import { OrganizationSchema, TextQueryMethod } from "@zitadel/proto/zitadel/object/v2/object_pb";
import { CreateCallbackRequest, OIDCService } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { OrganizationService } from "@zitadel/proto/zitadel/org/v2/org_service_pb";
import { CreateResponseRequest, SAMLService } from "@zitadel/proto/zitadel/saml/v2/saml_service_pb";
import { RequestChallenges } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { Checks, SessionService } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb";
import { SendEmailVerificationCodeSchema } from "@zitadel/proto/zitadel/user/v2/email_pb";
import type { FormData, RedirectURLsJson } from "@zitadel/proto/zitadel/user/v2/idp_pb";
import { NotificationType, SendPasswordResetLinkSchema } from "@zitadel/proto/zitadel/user/v2/password_pb";
import { SearchQuery, SearchQuerySchema } from "@zitadel/proto/zitadel/user/v2/query_pb";
import { SendInviteCodeSchema } from "@zitadel/proto/zitadel/user/v2/user_pb";
import {
  AddHumanUserRequest,
  AddHumanUserRequestSchema,
  ResendEmailCodeRequest,
  ResendEmailCodeRequestSchema,
  SendEmailCodeRequestSchema,
  SetPasswordRequest,
  SetPasswordRequestSchema,
  UpdateHumanUserRequest,
  UserService,
  VerifyPasskeyRegistrationRequest,
  VerifyU2FRegistrationRequest,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { unstable_cacheLife as cacheLife } from "next/cache";
import { getTranslations } from "next-intl/server";
import { getUserAgent } from "./fingerprint";
import { setSAMLFormCookie } from "./saml";
import { createServiceForHost } from "./service";

const useCache = process.env.DEBUG !== "true";

async function cacheWrapper<T>(callback: Promise<T>) {
  "use cache";
  cacheLife("hours");

  return callback;
}

export async function getHostedLoginTranslation({
  serviceConfig,
  organization,
  locale,
}: WithServiceConfig<{
  organization?: string;
  locale?: string;
}>) {
  const settingsService: Client<typeof SettingsService> = await createServiceForHost(SettingsService, serviceConfig);

  const callback = settingsService
    .getHostedLoginTranslation(
      {
        level: organization
          ? {
              case: "organizationId",
              value: organization,
            }
          : {
              case: "instance",
              value: true,
            },
        locale: locale,
      },
      {},
    )
    .then((resp) => {
      return resp.translations ? resp.translations : undefined;
    });

  return useCache ? cacheWrapper(callback) : callback;
}

export async function getBrandingSettings({
  serviceConfig,
  organization,
}: WithServiceConfig<{
  organization?: string;
}>) {
  const settingsService: Client<typeof SettingsService> = await createServiceForHost(SettingsService, serviceConfig);

  const callback = settingsService
    .getBrandingSettings({ ctx: makeReqCtx(organization) }, {})
    .then((resp) => (resp.settings ? resp.settings : undefined));

  return useCache ? cacheWrapper(callback) : callback;
}

export async function getLoginSettings({
  serviceConfig,
  organization,
}: WithServiceConfig<{
  organization?: string;
}>) {
  const settingsService: Client<typeof SettingsService> = await createServiceForHost(SettingsService, serviceConfig);

  const callback = settingsService
    .getLoginSettings({ ctx: makeReqCtx(organization) }, {})
    .then((resp) => (resp.settings ? resp.settings : undefined));

  return useCache ? cacheWrapper(callback) : callback;
}

export async function getSecuritySettings({ serviceConfig }: WithServiceConfig) {
  const settingsService: Client<typeof SettingsService> = await createServiceForHost(SettingsService, serviceConfig);

  const callback = settingsService.getSecuritySettings({}).then((resp) => (resp.settings ? resp.settings : undefined));

  return useCache ? cacheWrapper(callback) : callback;
}

export async function getLockoutSettings({ serviceConfig, orgId }: WithServiceConfig<{ orgId?: string }>) {
  const settingsService: Client<typeof SettingsService> = await createServiceForHost(SettingsService, serviceConfig);

  const callback = settingsService
    .getLockoutSettings({ ctx: makeReqCtx(orgId) }, {})
    .then((resp) => (resp.settings ? resp.settings : undefined));

  return useCache ? cacheWrapper(callback) : callback;
}

export async function getPasswordExpirySettings({ serviceConfig, orgId }: WithServiceConfig<{ orgId?: string }>) {
  const settingsService: Client<typeof SettingsService> = await createServiceForHost(SettingsService, serviceConfig);

  const callback = settingsService
    .getPasswordExpirySettings({ ctx: makeReqCtx(orgId) }, {})
    .then((resp) => (resp.settings ? resp.settings : undefined));

  return useCache ? cacheWrapper(callback) : callback;
}

export async function listIDPLinks({ serviceConfig, userId }: WithServiceConfig<{ userId: string }>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.listIDPLinks({ userId }, {});
}

export async function addOTPEmail({ serviceConfig, userId }: WithServiceConfig<{ userId: string }>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.addOTPEmail({ userId }, {});
}

export async function addOTPSMS({ serviceConfig, userId }: WithServiceConfig<{ userId: string }>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.addOTPSMS({ userId }, {});
}

export async function registerTOTP({ serviceConfig, userId }: WithServiceConfig<{ userId: string }>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.registerTOTP({ userId }, {});
}

export async function getGeneralSettings({ serviceConfig }: WithServiceConfig) {
  const settingsService: Client<typeof SettingsService> = await createServiceForHost(SettingsService, serviceConfig);

  const callback = settingsService.getGeneralSettings({}, {}).then((resp) => resp.supportedLanguages);

  return useCache ? cacheWrapper(callback) : callback;
}

export async function getLegalAndSupportSettings({
  serviceConfig,
  organization,
}: WithServiceConfig<{
  organization?: string;
}>) {
  const settingsService: Client<typeof SettingsService> = await createServiceForHost(SettingsService, serviceConfig);

  const callback = settingsService
    .getLegalAndSupportSettings({ ctx: makeReqCtx(organization) }, {})
    .then((resp) => (resp.settings ? resp.settings : undefined));

  return useCache ? cacheWrapper(callback) : callback;
}

export async function getPasswordComplexitySettings({
  serviceConfig,
  organization,
}: WithServiceConfig<{
  organization?: string;
}>) {
  const settingsService: Client<typeof SettingsService> = await createServiceForHost(SettingsService, serviceConfig);

  const callback = settingsService
    .getPasswordComplexitySettings({ ctx: makeReqCtx(organization) })
    .then((resp) => (resp.settings ? resp.settings : undefined));

  return useCache ? cacheWrapper(callback) : callback;
}

export async function createSessionFromChecks({
  serviceConfig,
  checks,
  lifetime,
}: WithServiceConfig<{
  checks: Checks;
  lifetime: Duration;
}>) {
  const sessionService: Client<typeof SessionService> = await createServiceForHost(SessionService, serviceConfig);

  const userAgent = await getUserAgent();

  return sessionService.createSession({ checks, lifetime, userAgent }, {});
}

export async function createSessionForUserIdAndIdpIntent({
  serviceConfig,
  userId,
  idpIntent,
  lifetime,
}: WithServiceConfig<{
  userId: string;
  idpIntent: {
    idpIntentId?: string | undefined;
    idpIntentToken?: string | undefined;
  };
  lifetime: Duration;
}>) {
  console.log("Creating session for userId and IDP intent", { userId, idpIntent, lifetime });
  const sessionService: Client<typeof SessionService> = await createServiceForHost(SessionService, serviceConfig);

  const userAgent = await getUserAgent();

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
    userAgent,
  });
}

export async function setSession({
  serviceConfig,
  sessionId,
  sessionToken,
  challenges,
  checks,
  lifetime,
}: WithServiceConfig<{
  sessionId: string;
  sessionToken: string;
  challenges: RequestChallenges | undefined;
  checks?: Checks;
  lifetime: Duration;
}>) {
  const sessionService: Client<typeof SessionService> = await createServiceForHost(SessionService, serviceConfig);

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
  serviceConfig,
  sessionId,
  sessionToken,
}: WithServiceConfig<{
  sessionId: string;
  sessionToken: string;
}>) {
  const sessionService: Client<typeof SessionService> = await createServiceForHost(SessionService, serviceConfig);

  return sessionService.getSession({ sessionId, sessionToken }, {});
}

export async function deleteSession({
  serviceConfig,
  sessionId,
  sessionToken,
}: WithServiceConfig<{
  sessionId: string;
  sessionToken: string;
}>) {
  const sessionService: Client<typeof SessionService> = await createServiceForHost(SessionService, serviceConfig);

  return sessionService.deleteSession({ sessionId, sessionToken }, {});
}

type ListSessionsCommand = WithServiceConfig<{
  ids: string[];
}>;

export async function listSessions({ serviceConfig, ids }: ListSessionsCommand) {
  const sessionService: Client<typeof SessionService> = await createServiceForHost(SessionService, serviceConfig);

  return sessionService.listSessions(
    {
      queries: [
        {
          query: {
            case: "idsQuery",
            value: { ids },
          },
        },
      ],
    },
    {},
  );
}

export type AddHumanUserData = WithServiceConfig<{
  firstName: string;
  lastName: string;
  email: string;
  password?: string;
  organization: string;
}>;

export async function addHumanUser({ serviceConfig, email, firstName, lastName, password, organization }: AddHumanUserData) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  let addHumanUserRequest: AddHumanUserRequest = create(AddHumanUserRequestSchema, {
    email: {
      email,
      verification: {
        case: "isVerified",
        value: false,
      },
    },
    username: email,
    profile: { givenName: firstName, familyName: lastName },
    passwordType: password ? { case: "password", value: { password } } : undefined,
  });

  if (organization) {
    const organizationSchema = create(OrganizationSchema, {
      org: { case: "orgId", value: organization },
    });

    addHumanUserRequest = {
      ...addHumanUserRequest,
      organization: organizationSchema,
    };
  }

  return userService.addHumanUser(addHumanUserRequest);
}

export async function addHuman({ serviceConfig, request }: WithServiceConfig<{ request: AddHumanUserRequest }>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.addHumanUser(request);
}

export async function updateHuman({
  serviceConfig,
  request,
}: WithServiceConfig<{
  request: UpdateHumanUserRequest;
}>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.updateHumanUser(request);
}

export async function verifyTOTPRegistration({
  serviceConfig,
  code,
  userId,
}: WithServiceConfig<{
  code: string;
  userId: string;
}>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.verifyTOTPRegistration({ code, userId }, {});
}

export async function getUserByID({ serviceConfig, userId }: WithServiceConfig<{ userId: string }>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.getUserByID({ userId }, {});
}

export async function humanMFAInitSkipped({ serviceConfig, userId }: WithServiceConfig<{ userId: string }>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.humanMFAInitSkipped({ userId }, {});
}

export async function verifyInviteCode({
  serviceConfig,
  userId,
  verificationCode,
}: WithServiceConfig<{
  userId: string;
  verificationCode: string;
}>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.verifyInviteCode({ userId, verificationCode }, {});
}

export async function sendEmailCode({
  serviceConfig,
  userId,
  urlTemplate,
}: WithServiceConfig<{
  userId: string;
  urlTemplate: string;
}>) {
  let medium = create(SendEmailCodeRequestSchema, { userId });

  medium = create(SendEmailCodeRequestSchema, {
    ...medium,
    verification: {
      case: "sendCode",
      value: create(SendEmailVerificationCodeSchema, {
        urlTemplate,
      }),
    },
  });

  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.sendEmailCode(medium, {});
}

export async function createInviteCode({
  serviceConfig,
  urlTemplate,
  userId,
}: WithServiceConfig<{
  urlTemplate: string;
  userId: string;
}>) {
  let medium = create(SendInviteCodeSchema, {
    applicationName: process.env.NEXT_PUBLIC_APPLICATION_NAME || "Zitadel Login",
  });

  medium = {
    ...medium,
    urlTemplate,
  };

  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

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

export type ListUsersCommand = WithServiceConfig<{
  loginName?: string;
  userName?: string;
  email?: string;
  phone?: string;
  organizationId?: string;
}>;

export async function listUsers({ serviceConfig, loginName, userName, phone, email, organizationId }: ListUsersCommand) {
  const queries: SearchQuery[] = [];

  // either use loginName or userName, email, phone
  if (loginName) {
    queries.push(
      create(SearchQuerySchema, {
        query: {
          case: "loginNameQuery",
          value: {
            loginName,
            method: TextQueryMethod.EQUALS,
          },
        },
      }),
    );
  } else if (userName || email || phone) {
    const orQueries: SearchQuery[] = [];

    if (userName) {
      const userNameQuery = create(SearchQuerySchema, {
        query: {
          case: "userNameQuery",
          value: {
            userName,
            method: TextQueryMethod.EQUALS,
          },
        },
      });
      orQueries.push(userNameQuery);
    }

    if (email) {
      const emailQuery = create(SearchQuerySchema, {
        query: {
          case: "emailQuery",
          value: {
            emailAddress: email,
            method: TextQueryMethod.EQUALS,
          },
        },
      });
      orQueries.push(emailQuery);
    }

    if (phone) {
      const phoneQuery = create(SearchQuerySchema, {
        query: {
          case: "phoneQuery",
          value: {
            number: phone,
            method: TextQueryMethod.EQUALS,
          },
        },
      });
      orQueries.push(phoneQuery);
    }

    queries.push(
      create(SearchQuerySchema, {
        query: {
          case: "orQuery",
          value: {
            queries: orQueries,
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

  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.listUsers({ queries });
}

export type SearchUsersCommand = WithServiceConfig<{
  searchValue: string;
  loginSettings: LoginSettings;
  organizationId?: string;
  suffix?: string;
}>;

const PhoneQuery = (searchValue: string) =>
  create(SearchQuerySchema, {
    query: {
      case: "phoneQuery",
      value: {
        number: searchValue,
        method: TextQueryMethod.EQUALS,
      },
    },
  });

const LoginNameQuery = (searchValue: string) =>
  create(SearchQuerySchema, {
    query: {
      case: "loginNameQuery",
      value: {
        loginName: searchValue,
        method: TextQueryMethod.EQUALS_IGNORE_CASE,
      },
    },
  });

const EmailQuery = (searchValue: string) =>
  create(SearchQuerySchema, {
    query: {
      case: "emailQuery",
      value: {
        emailAddress: searchValue,
        method: TextQueryMethod.EQUALS_IGNORE_CASE,
      },
    },
  });

/**
 * this is a dedicated search function to search for users from the loginname page
 * it searches users based on the loginName or userName and org suffix combination, and falls back to email and phone if no users are found
 *  */
export async function searchUsers({
  serviceConfig,
  searchValue,
  loginSettings,
  organizationId,
  suffix,
}: SearchUsersCommand) {
  const queries: SearchQuery[] = [];

  const t = await getTranslations("zitadel");

  // if a suffix is provided, we search for the userName concatenated with the suffix
  if (suffix) {
    const searchValueWithSuffix = `${searchValue}@${suffix}`;
    const loginNameQuery = LoginNameQuery(searchValueWithSuffix);
    queries.push(loginNameQuery);
  } else {
    const loginNameQuery = LoginNameQuery(searchValue);
    queries.push(loginNameQuery);
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

  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  const loginNameResult = await userService.listUsers({ queries });

  if (!loginNameResult || !loginNameResult.details) {
    return { error: t("errors.errorOccured") };
  }

  if (loginNameResult.result.length > 1) {
    return { error: t("errors.multipleUsersFound") };
  }

  if (loginNameResult.result.length == 1) {
    return loginNameResult;
  }

  const emailAndPhoneQueries: SearchQuery[] = [];
  if (loginSettings.disableLoginWithEmail && loginSettings.disableLoginWithPhone) {
    // Both email and phone login are disabled, return empty result
    return { result: [] };
  } else if (loginSettings.disableLoginWithEmail && searchValue.length <= 20) {
    const phoneQuery = PhoneQuery(searchValue);
    emailAndPhoneQueries.push(phoneQuery);
  } else if (loginSettings.disableLoginWithPhone) {
    const emailQuery = EmailQuery(searchValue);
    emailAndPhoneQueries.push(emailQuery);
  } else {
    const orQuery: SearchQuery[] = [];

    const emailQuery = EmailQuery(searchValue);
    orQuery.push(emailQuery);

    let phoneQuery;
    if (searchValue.length <= 20) {
      phoneQuery = PhoneQuery(searchValue);
      orQuery.push(phoneQuery);
    }

    emailAndPhoneQueries.push(
      create(SearchQuerySchema, {
        query: {
          case: "orQuery",
          value: {
            queries: orQuery,
          },
        },
      }),
    );
  }

  if (organizationId) {
    emailAndPhoneQueries.push(
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

  const emailOrPhoneResult = await userService.listUsers({
    queries: emailAndPhoneQueries,
  });

  if (!emailOrPhoneResult || !emailOrPhoneResult.details) {
    return { error: t("errors.errorOccured") };
  }

  if (emailOrPhoneResult.result.length > 1) {
    return { error: t("errors.multipleUsersFound") };
  }

  if (emailOrPhoneResult.result.length == 1) {
    return emailOrPhoneResult;
  }

  // No users found - return empty result, not an error
  return { result: [] };
}

export async function getDefaultOrg({ serviceConfig }: WithServiceConfig): Promise<Organization | null> {
  const orgService: Client<typeof OrganizationService> = await createServiceForHost(OrganizationService, serviceConfig);

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

export async function getOrgsByDomain({ serviceConfig, domain }: WithServiceConfig<{ domain: string }>) {
  const orgService: Client<typeof OrganizationService> = await createServiceForHost(OrganizationService, serviceConfig);

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
  serviceConfig,
  idpId,
  urls,
}: WithServiceConfig<{
  idpId: string;
  urls: RedirectURLsJson;
}>): Promise<string | null> {
  // Use empty publicHost to avoid issues with redirect URIs pointing to the login UI instead of the zitadel API
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, {
    ...serviceConfig,
    publicHost: "",
  });

  return userService
    .startIdentityProviderIntent({
      idpId,
      content: {
        case: "urls",
        value: urls,
      },
    })
    .then(async (resp) => {
      if (resp.nextStep.case === "authUrl" && resp.nextStep.value) {
        return resp.nextStep.value;
      } else if (resp.nextStep.case === "formData" && resp.nextStep.value) {
        const formData: FormData = resp.nextStep.value;
        const redirectUrl = "/saml-post";

        try {
          // Log the attempt with structure inspection
          console.log("Attempting to stringify formData.fields:", {
            fields: formData.fields,
            fieldsType: typeof formData.fields,
            fieldsKeys: Object.keys(formData.fields || {}),
            fieldsEntries: Object.entries(formData.fields || {}),
          });

          const stringifiedFields = JSON.stringify(formData.fields);
          console.log("Successfully stringified formData.fields, length:", stringifiedFields.length);

          // Check cookie size limits (typical limit is 4KB)
          if (stringifiedFields.length > 4000) {
            console.warn(
              `SAML form cookie value is large (${stringifiedFields.length} characters), may exceed browser limits`,
            );
          }

          const dataId = await setSAMLFormCookie(stringifiedFields);
          const params = new URLSearchParams({ url: formData.url, id: dataId });

          return `${redirectUrl}?${params.toString()}`;
        } catch (stringifyError) {
          console.error("JSON serialization failed:", stringifyError);
          throw new Error(
            `Failed to serialize SAML form data: ${stringifyError instanceof Error ? stringifyError.message : String(stringifyError)}`,
          );
        }
      } else {
        return null;
      }
    });
}

export async function startLDAPIdentityProviderFlow({
  serviceConfig,
  idpId,
  username,
  password,
}: WithServiceConfig<{
  idpId: string;
  username: string;
  password: string;
}>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.startIdentityProviderIntent({
    idpId,
    content: {
      case: "ldap",
      value: {
        username,
        password,
      },
    },
  });
}

export async function getAuthRequest({
  serviceConfig,
  authRequestId,
}: WithServiceConfig<{
  authRequestId: string;
}>) {
  const oidcService = await createServiceForHost(OIDCService, serviceConfig);

  return oidcService.getAuthRequest({
    authRequestId,
  });
}

export async function getDeviceAuthorizationRequest({
  serviceConfig,
  userCode,
}: WithServiceConfig<{
  userCode: string;
}>) {
  const oidcService = await createServiceForHost(OIDCService, serviceConfig);

  return oidcService.getDeviceAuthorizationRequest({
    userCode,
  });
}

export async function authorizeOrDenyDeviceAuthorization({
  serviceConfig,
  deviceAuthorizationId,
  session,
}: WithServiceConfig<{
  deviceAuthorizationId: string;
  session?: { sessionId: string; sessionToken: string };
}>) {
  const oidcService = await createServiceForHost(OIDCService, serviceConfig);

  return oidcService.authorizeOrDenyDeviceAuthorization({
    deviceAuthorizationId,
    decision: session
      ? {
          case: "session",
          value: session,
        }
      : {
          case: "deny",
          value: {},
        },
  });
}

export async function createCallback({ serviceConfig, req }: WithServiceConfig<{ req: CreateCallbackRequest }>) {
  const oidcService = await createServiceForHost(OIDCService, serviceConfig);

  return oidcService.createCallback(req);
}

export async function getSAMLRequest({
  serviceConfig,
  samlRequestId,
}: WithServiceConfig<{
  samlRequestId: string;
}>) {
  const samlService = await createServiceForHost(SAMLService, serviceConfig);

  return samlService.getSAMLRequest({
    samlRequestId,
  });
}

export async function createResponse({ serviceConfig, req }: WithServiceConfig<{ req: CreateResponseRequest }>) {
  const samlService = await createServiceForHost(SAMLService, serviceConfig);

  return samlService.createResponse(req);
}

export async function verifyEmail({
  serviceConfig,
  userId,
  verificationCode,
}: WithServiceConfig<{
  userId: string;
  verificationCode: string;
}>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.verifyEmail(
    {
      userId,
      verificationCode,
    },
    {},
  );
}

export async function resendEmailCode({
  serviceConfig,
  userId,
  urlTemplate,
}: WithServiceConfig<{
  userId: string;
  urlTemplate: string;
}>) {
  let request: ResendEmailCodeRequest = create(ResendEmailCodeRequestSchema, {
    userId,
  });

  const medium = create(SendEmailVerificationCodeSchema, {
    urlTemplate,
  });

  request = { ...request, verification: { case: "sendCode", value: medium } };

  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.resendEmailCode(request, {});
}

export async function retrieveIDPIntent({
  serviceConfig,
  id,
  token,
}: WithServiceConfig<{
  id: string;
  token: string;
}>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.retrieveIdentityProviderIntent({ idpIntentId: id, idpIntentToken: token }, {});
}

export async function getIDPByID({ serviceConfig, id }: WithServiceConfig<{ id: string }>) {
  const idpService: Client<typeof IdentityProviderService> = await createServiceForHost(
    IdentityProviderService,
    serviceConfig,
  );

  return idpService.getIDPByID({ id }, {}).then((resp) => resp.idp);
}

export async function addIDPLink({
  serviceConfig,
  idp,
  userId,
}: WithServiceConfig<{
  idp: { id: string; userId: string; userName: string };
  userId: string;
}>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

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

export async function passwordReset({
  serviceConfig,
  userId,
  urlTemplate,
}: WithServiceConfig<{
  userId: string;
  urlTemplate?: string;
}>) {
  let medium = create(SendPasswordResetLinkSchema, {
    notificationType: NotificationType.Email,
  });

  medium = {
    ...medium,
    urlTemplate,
  };

  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

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

export async function setUserPassword({
  serviceConfig,
  userId,
  password,
  code,
}: WithServiceConfig<{
  userId: string;
  password: string;
  code?: string;
}>) {
  let payload = create(SetPasswordRequestSchema, {
    userId,
    newPassword: {
      password,
    },
  });

  if (code) {
    payload = {
      ...payload,
      verification: {
        case: "verificationCode",
        value: code,
      },
    };
  }

  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.setPassword(payload, {}).catch((error) => {
    // throw error if failed precondition (ex. User is not yet initialized)
    if (error.code === 9 && error.message) {
      return { error: error.message };
    } else {
      throw error;
    }
  });
}

export async function setPassword({
  serviceConfig,
  payload,
}: WithServiceConfig<{
  payload: SetPasswordRequest;
}>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.setPassword(payload, {});
}

/**
 *
 * @param host
 * @param userId the id of the user where the email should be set
 * @returns the newly set email
 */
export async function createPasskeyRegistrationLink({
  serviceConfig,
  userId,
}: WithServiceConfig<{
  userId: string;
}>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

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
 * @param host
 * @param userId the id of the user where the email should be set
 * @param domain the domain on which the factor is registered
 * @returns the newly set email
 */
export async function registerU2F({
  serviceConfig,
  userId,
  domain,
}: WithServiceConfig<{
  userId: string;
  domain: string;
}>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.registerU2F({
    userId,
    domain,
  });
}

/**
 *
 * @param host
 * @param request the request object for verifying U2F registration
 * @returns the result of the verification
 */
export async function verifyU2FRegistration({
  serviceConfig,
  request,
}: WithServiceConfig<{
  request: VerifyU2FRegistrationRequest;
}>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.verifyU2FRegistration(request, {});
}

/**
 *
 * @param host
 * @param orgId the organization ID
 * @param linking_allowed whether linking is allowed
 * @returns the active identity providers
 */
export async function getActiveIdentityProviders({
  serviceConfig,
  orgId,
  linking_allowed,
}: WithServiceConfig<{
  orgId?: string;
  linking_allowed?: boolean;
}>) {
  const props: any = { ctx: makeReqCtx(orgId) };
  if (linking_allowed) {
    props.linkingAllowed = linking_allowed;
  }
  const settingsService: Client<typeof SettingsService> = await createServiceForHost(SettingsService, serviceConfig);

  return settingsService.getActiveIdentityProviders(props, {});
}

/**
 *
 * @param host
 * @param request the request object for verifying passkey registration
 * @returns the result of the verification
 */
export async function verifyPasskeyRegistration({
  serviceConfig,
  request,
}: WithServiceConfig<{
  request: VerifyPasskeyRegistrationRequest;
}>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.verifyPasskeyRegistration(request, {});
}

/**
 *
 * @param host
 * @param userId the id of the user where the email should be set
 * @param code the code for registering the passkey
 * @param domain the domain on which the factor is registered
 * @returns the newly set email
 */
export async function registerPasskey({
  serviceConfig,
  userId,
  code,
  domain,
}: WithServiceConfig<{
  userId: string;
  code: { id: string; code: string };
  domain: string;
}>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.registerPasskey({
    userId,
    code,
    domain,
  });
}

/**
 *
 * @param host
 * @param userId the id of the user where the email should be set
 * @returns the list of authentication method types
 */
export async function listAuthenticationMethodTypes({
  serviceConfig,
  userId,
}: WithServiceConfig<{
  userId: string;
}>) {
  const userService: Client<typeof UserService> = await createServiceForHost(UserService, serviceConfig);

  return userService.listAuthenticationMethodTypes({
    userId,
  });
}

export interface ServiceConfig {
  baseUrl: string;
  instanceHost?: string; // only for multi-tenant
  publicHost?: string; // only for multi-tenant
}

/**
 * Base type that all function parameters must extend to ensure serviceConfig is always required
 */
export type WithServiceConfig<T = {}> = T & {
  serviceConfig: ServiceConfig;
};

export function createServerTransport(token: string, serviceConfig: ServiceConfig) {
  return libCreateServerTransport(token, {
    baseUrl: serviceConfig.baseUrl,
    interceptors:
      !process.env.CUSTOM_REQUEST_HEADERS && !serviceConfig.instanceHost && !serviceConfig.publicHost
        ? undefined
        : [
            (next) => {
              return (req) => {
                // Apply headers from serviceConfig
                if (serviceConfig.instanceHost) {
                  req.header.set("x-zitadel-instance-host", serviceConfig.instanceHost);
                }
                if (serviceConfig.publicHost) {
                  req.header.set("x-zitadel-public-host", serviceConfig.publicHost);
                }

                // Apply headers from CUSTOM_REQUEST_HEADERS environment variable
                if (process.env.CUSTOM_REQUEST_HEADERS) {
                  process.env.CUSTOM_REQUEST_HEADERS.split(",").forEach((header) => {
                    const kv = header.indexOf(":");
                    if (kv > 0) {
                      req.header.set(header.slice(0, kv).trim(), header.slice(kv + 1).trim());
                    } else {
                      console.warn(`Skipping malformed header: ${header}`);
                    }
                  });
                }

                return next(req);
              };
            },
          ],
  });
}
