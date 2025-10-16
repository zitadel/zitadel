"use server";

import { create } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import { idpTypeToIdentityProviderType, idpTypeToSlug } from "../idp";

import { PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { UserState } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { getServiceUrlFromHeaders } from "../service-url";
import {
  getActiveIdentityProviders,
  getIDPByID,
  getLoginSettings,
  getOrgsByDomain,
  listAuthenticationMethodTypes,
  listIDPLinks,
  searchUsers,
  SearchUsersCommand,
  startIdentityProviderFlow,
} from "../zitadel";
import { createSessionAndUpdateCookie } from "./cookie";
import { getOriginalHost } from "./host";
import { IDPLink } from "@zitadel/proto/zitadel/user/v2/idp_pb";

export type SendLoginnameCommand = {
  loginName: string;
  requestId?: string;
  organization?: string;
  suffix?: string;
};

const ORG_SUFFIX_REGEX = /(?<=@)(.+)/;

/**
 * Validates if the user's login name matches the login settings constraints.
 * Returns true if validation passes, false otherwise.
 */
function validateUserLoginName(params: {
  user: {
    preferredLoginName: string;
    type: { case?: string; value?: any };
  };
  loginName: string;
  concatLoginname: string;
  userLoginSettings?: {
    disableLoginWithEmail?: boolean;
    disableLoginWithPhone?: boolean;
  } | null;
}): boolean {
  const { user, loginName, concatLoginname, userLoginSettings } = params;

  const humanUser = user.type.case === "human" ? user.type.value : undefined;

  // recheck login settings after user discovery, as the search might have been done without org scope
  if (userLoginSettings?.disableLoginWithEmail && userLoginSettings?.disableLoginWithPhone) {
    return user.preferredLoginName === concatLoginname;
  } else if (userLoginSettings?.disableLoginWithEmail) {
    return user.preferredLoginName === concatLoginname || humanUser?.phone?.phone === loginName;
  } else if (userLoginSettings?.disableLoginWithPhone) {
    return user.preferredLoginName === concatLoginname || humanUser?.email?.email === loginName;
  }

  return true;
}

/**
 * Routes user to the appropriate authentication method based on available auth methods and login settings.
 */
async function handleAuthenticationMethodRouting(params: {
  methods: { authMethodTypes: AuthenticationMethodType[] };
  userLoginSettings?: {
    allowUsernamePassword?: boolean;
    passkeysType?: PasskeysType;
  } | null;
  loginName: string;
  userId: string;
  organization?: string;
  requestId?: string;
  serviceUrl: string;
  t: Awaited<ReturnType<typeof getTranslations<"loginname">>>;
}): Promise<{ redirect: string } | { error: string }> {
  const { methods, userLoginSettings, loginName, userId, organization, requestId, serviceUrl, t } = params;

  if (methods.authMethodTypes.length == 1) {
    const method = methods.authMethodTypes[0];
    switch (method) {
      case AuthenticationMethodType.PASSWORD: // user has only password as auth method
        if (!userLoginSettings?.allowUsernamePassword) {
          // Check if user has IDPs available as alternative, that could eventually be used to register/link.
          const idpResp = await redirectUserToIDP({
            serviceUrl,
            userId,
            organization,
            requestId,
            t,
          });
          if (idpResp && "redirect" in idpResp) {
            return idpResp;
          }

          return {
            error: t("errors.usernamePasswordNotAllowed"),
          };
        }

        const paramsPassword = new URLSearchParams({
          loginName,
        });

        // TODO: does this have to be checked in loginSettings.allowDomainDiscovery

        if (organization) {
          paramsPassword.append("organization", organization);
        }

        if (requestId) {
          paramsPassword.append("requestId", requestId);
        }

        return {
          redirect: "/password?" + paramsPassword,
        };

      case AuthenticationMethodType.PASSKEY: // AuthenticationMethodType.AUTHENTICATION_METHOD_TYPE_PASSKEY
        if (userLoginSettings?.passkeysType === PasskeysType.NOT_ALLOWED) {
          return {
            error: t("errors.passkeysNotAllowed"),
          };
        }

        const paramsPasskey = new URLSearchParams({
          loginName,
        });
        if (requestId) {
          paramsPasskey.append("requestId", requestId);
        }

        if (organization) {
          paramsPasskey.append("organization", organization);
        }

        return { redirect: "/passkey?" + paramsPasskey };

      case AuthenticationMethodType.IDP:
        const resp = await redirectUserToIDP({
          serviceUrl,
          userId,
          organization,
          requestId,
          t,
        });

        if (resp && "error" in resp) {
          return { error: resp.error };
        }

        if (resp && "redirect" in resp) {
          return resp;
        }

        // IDP is the user's only auth method but no suitable IDP found
        return { error: t("errors.noSuitableIDPFound") };
    }
  } else {
    // prefer passkey in favor of other methods
    if (methods.authMethodTypes.includes(AuthenticationMethodType.PASSKEY)) {
      const passkeyParams = new URLSearchParams({
        loginName,
        altPassword: `${methods.authMethodTypes.includes(AuthenticationMethodType.PASSWORD) && userLoginSettings?.allowUsernamePassword}`, // show alternative password option only if allowed
      });

      if (requestId) {
        passkeyParams.append("requestId", requestId);
      }

      if (organization) {
        passkeyParams.append("organization", organization);
      }

      return { redirect: "/passkey?" + passkeyParams };
    } else if (methods.authMethodTypes.includes(AuthenticationMethodType.IDP)) {
      const idpResp = await redirectUserToIDP({
        serviceUrl,
        userId,
        organization,
        requestId,
        t,
      });

      if (idpResp) {
        return idpResp;
      }

      // IDP is one of the user's auth methods but no suitable IDP found
      return { error: t("errors.noSuitableIDPFound") };
    } else if (methods.authMethodTypes.includes(AuthenticationMethodType.PASSWORD)) {
      // Check if password authentication is allowed
      if (!userLoginSettings?.allowUsernamePassword) {
        return {
          error: "Username Password not allowed! Contact your administrator for more information.",
        };
      }

      // user has no passkey setup and login settings allow passwords
      const paramsPasswordDefault = new URLSearchParams({
        loginName,
      });

      if (requestId) {
        paramsPasswordDefault.append("requestId", requestId);
      }

      if (organization) {
        paramsPasswordDefault.append("organization", organization);
      }

      return {
        redirect: "/password?" + paramsPasswordDefault,
      };
    }
  }

  // No matching authentication method found (should not happen in normal cases)
  return { error: t("errors.noAuthenticationMethodAvailable") };
}

/**
 * Helper function to redirect user to their identity provider.
 * Checks user-specific IDP links first, then falls back to organization-level active IDPs.
 * Returns:
 * - { redirect: string } if exactly one IDP is found and flow started successfully
 * - { error: string } if IDP was found but flow failed to start
 * - undefined if no suitable IDP is found (0 or 2+ IDPs) - allows caller to handle fallback
 */
async function redirectUserToIDP(params: {
  serviceUrl: string;
  userId?: string;
  organization?: string;
  requestId?: string;
  t: Awaited<ReturnType<typeof getTranslations<"loginname">>>;
}): Promise<{ redirect: string } | { error: string } | undefined> {
  const { serviceUrl, userId, organization, requestId, t } = params;

  // If userId is provided, check for user-specific IDP links first
  let identityProviders: IDPLink[] = [];
  if (userId) {
    identityProviders = await listIDPLinks({
      serviceUrl,
      userId,
    }).then((resp) => {
      return resp.result;
    });
  }

  // If no IDP links exist for the user (or no userId provided), try to get active IDPs from the organization
  if (identityProviders.length === 0) {
    const activeIdps = await getActiveIdentityProviders({
      serviceUrl,
      orgId: organization,
    }).then((resp) => {
      return resp.identityProviders;
    });

    // If exactly one active IDP exists in the organization, redirect to it
    if (activeIdps.length === 1) {
      const host = await getOriginalHost();

      const identityProviderType = activeIdps[0].type;
      const provider = idpTypeToSlug(identityProviderType);

      const urlParams = new URLSearchParams();

      if (userId) {
        urlParams.set("userId", userId);
      }

      if (requestId) {
        urlParams.set("requestId", requestId);
      }

      if (organization) {
        urlParams.set("organization", organization);
      }

      const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

      const url = await startIdentityProviderFlow({
        serviceUrl,
        idpId: activeIdps[0].id,
        urls: {
          successUrl:
            `${host.includes("localhost") ? "http://" : "https://"}${host}${basePath}/idp/${provider}/success?` +
            new URLSearchParams(urlParams),
          failureUrl:
            `${host.includes("localhost") ? "http://" : "https://"}${host}${basePath}/idp/${provider}/failure?` +
            new URLSearchParams(urlParams),
        },
      });

      if (!url) {
        return { error: t("errors.couldNotStartIDPFlow") };
      }

      return { redirect: url };
    }
  }

  if (identityProviders.length === 1) {
    const host = await getOriginalHost();

    const identityProviderId = identityProviders[0].idpId;

    const idp = await getIDPByID({
      serviceUrl,
      id: identityProviderId,
    });

    const idpType = idp?.type;

    if (!idp || !idpType) {
      throw new Error(t("errors.couldNotFindIdentityProvider"));
    }

    const identityProviderType = idpTypeToIdentityProviderType(idpType);
    const provider = idpTypeToSlug(identityProviderType);

    const urlParams = new URLSearchParams();

    if (userId) {
      urlParams.set("userId", userId);
    }

    if (requestId) {
      urlParams.set("requestId", requestId);
    }

    if (organization) {
      urlParams.set("organization", organization);
    }

    const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

    const url = await startIdentityProviderFlow({
      serviceUrl,
      idpId: idp.id,
      urls: {
        successUrl:
          `${host.includes("localhost") ? "http://" : "https://"}${host}${basePath}/idp/${provider}/success?` +
          new URLSearchParams(urlParams),
        failureUrl:
          `${host.includes("localhost") ? "http://" : "https://"}${host}${basePath}/idp/${provider}/failure?` +
          new URLSearchParams(urlParams),
      },
    });

    if (!url) {
      return { error: t("errors.couldNotStartIDPFlow") };
    }

    return { redirect: url };
  }

  // No suitable IDP found (0 or multiple IDPs) - return undefined to allow caller to handle fallback
  return undefined;
}

export async function sendLoginname(command: SendLoginnameCommand): Promise<{ redirect: string } | { error: string }> {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const t = await getTranslations("loginname");

  const loginSettingsByContext = await getLoginSettings({
    serviceUrl,
    organization: command.organization,
  });

  if (!loginSettingsByContext) {
    return { error: t("errors.couldNotGetLoginSettings") };
  }

  let searchUsersRequest: SearchUsersCommand = {
    serviceUrl,
    searchValue: command.loginName,
    organizationId: command.organization,
    loginSettings: loginSettingsByContext,
    suffix: command.suffix,
  };

  const searchResult = await searchUsers(searchUsersRequest);

  if ("error" in searchResult && searchResult.error) {
    return searchResult;
  }

  if (!("result" in searchResult)) {
    return { error: t("errors.couldNotSearchUsers") };
  }

  const { result: potentialUsers } = searchResult;

  if (potentialUsers.length > 1) {
    return { error: t("errors.moreThanOneUserFound") };
  } else if (potentialUsers.length == 1 && potentialUsers[0].userId) {
    const user = potentialUsers[0];
    const userId = potentialUsers[0].userId;

    const userLoginSettings = await getLoginSettings({
      serviceUrl,
      organization: user.details?.resourceOwner,
    });

    // compare with the concatenated suffix when set
    const concatLoginname = command.suffix ? `${command.loginName}@${command.suffix}` : command.loginName;

    // Validate user login name against login settings
    if (
      !validateUserLoginName({
        user,
        loginName: command.loginName,
        concatLoginname,
        userLoginSettings,
      })
    ) {
      return { error: t("errors.userNotFound") };
    }

    const checks = create(ChecksSchema, {
      user: { search: { case: "userId", value: userId } },
    });

    const sessionOrError = await createSessionAndUpdateCookie({
      checks,
      requestId: command.requestId,
    }).catch((error) => {
      if (error?.rawMessage === "Errors.User.NotActive (SESSION-Gj4ko)") {
        return { error: t("errors.userNotActive") };
      }
      throw error;
    });

    if ("error" in sessionOrError) {
      return sessionOrError;
    }

    const session = sessionOrError;

    if (!session.factors?.user?.id) {
      return { error: t("errors.couldNotCreateSession") };
    }

    // We return an error since initial users are not supported with Login V2
    if (user.state === UserState.INITIAL) {
      return { error: t("errors.initialUserNotSupported") };
    }

    // Resolve organization from command or session
    const organization = command.organization ?? session.factors?.user?.organizationId;

    const methods = await listAuthenticationMethodTypes({
      serviceUrl,
      userId: session.factors?.user?.id,
    });

    // always resend invite if user has no auth method set
    if (!methods.authMethodTypes || !methods.authMethodTypes.length) {
      const params = new URLSearchParams({
        loginName: session.factors?.user?.loginName as string,
        send: "true", // set this to true to request a new code immediately
        invite: "true",
      });

      if (command.requestId) {
        params.append("requestId", command.requestId);
      }

      if (organization) {
        params.append("organization", organization);
      }

      return { redirect: `/verify?` + params };
    }

    // Route to appropriate authentication method
    return handleAuthenticationMethodRouting({
      methods,
      userLoginSettings,
      loginName: session.factors?.user?.loginName as string,
      userId,
      organization,
      requestId: command.requestId,
      serviceUrl,
      t,
    });
  }

  // user not found, check if register is enabled on instance / organization context
  if (loginSettingsByContext?.allowRegister && !loginSettingsByContext?.allowUsernamePassword) {
    const resp = await redirectUserToIDP({
      serviceUrl,
      userId: undefined,
      organization: command.organization,
      requestId: command.requestId,
      t,
    });
    if (resp) {
      return resp;
    }
    return { error: t("errors.userNotFound") };
  } else if (loginSettingsByContext?.allowRegister && loginSettingsByContext?.allowUsernamePassword) {
    let orgToRegisterOn: string | undefined = command.organization;

    if (
      !loginSettingsByContext?.ignoreUnknownUsernames &&
      !orgToRegisterOn &&
      command.loginName &&
      ORG_SUFFIX_REGEX.test(command.loginName)
    ) {
      const matched = ORG_SUFFIX_REGEX.exec(command.loginName);
      const suffix = matched?.[1] ?? "";

      // this just returns orgs where the suffix is set as primary domain
      const orgs = await getOrgsByDomain({
        serviceUrl,
        domain: suffix,
      });
      const orgToCheckForDiscovery = orgs.result && orgs.result.length === 1 ? orgs.result[0].id : undefined;

      const orgLoginSettings = await getLoginSettings({
        serviceUrl,
        organization: orgToCheckForDiscovery,
      });
      if (orgLoginSettings?.allowDomainDiscovery) {
        orgToRegisterOn = orgToCheckForDiscovery;
      }
    }

    // do not register user if ignoreUnknownUsernames is set
    if (orgToRegisterOn && !loginSettingsByContext?.ignoreUnknownUsernames) {
      const params = new URLSearchParams({ organization: orgToRegisterOn });

      if (command.requestId) {
        params.set("requestId", command.requestId);
      }

      if (command.loginName) {
        params.set("email", command.loginName);
      }

      return { redirect: "/register?" + params };
    }
  }

  if (loginSettingsByContext?.ignoreUnknownUsernames) {
    const paramsPasswordDefault = new URLSearchParams({
      loginName: command.loginName,
    });

    if (command.requestId) {
      paramsPasswordDefault.append("requestId", command.requestId);
    }

    if (command.organization) {
      paramsPasswordDefault.append("organization", command.organization);
    }

    return { redirect: "/password?" + paramsPasswordDefault };
  }

  // fallbackToPassword

  return { error: t("errors.userNotFound") };
}
