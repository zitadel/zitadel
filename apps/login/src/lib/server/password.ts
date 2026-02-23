"use server";

import { createSessionAndUpdateCookie, setSessionAndUpdateCookie } from "@/lib/server/cookie";
import {
  getLockoutSettings,
  getLoginSettings,
  getPasswordExpirySettings,
  getSession,
  getUserByID,
  listAuthenticationMethodTypes,
  passwordReset,
  searchUsers,
  setPassword,
  setUserPassword,
} from "@/lib/zitadel";
import { create, Duration } from "@zitadel/client";
import { Checks, ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { User, UserState } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { SetPasswordRequestSchema } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import { completeFlowOrGetUrl } from "../client";
import { getSessionCookieById, getSessionCookieByLoginName } from "../cookies";
import { getServiceConfig } from "../service-url";
import {
  checkEmailVerification,
  checkMFAFactors,
  checkPasswordChangeRequired,
  checkUserVerification,
} from "../verify-helper";
import { getPublicHostWithProtocol } from "./host";

type ResetPasswordCommand = {
  loginName: string;
  organization?: string;
  defaultOrganization?: string;
  requestId?: string;
};

export async function resetPassword(command: ResetPasswordCommand) {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const t = await getTranslations("password");

  // Get the original host that the user sees with protocol
  const hostWithProtocol = await getPublicHostWithProtocol(_headers);

  const loginSettings = await getLoginSettings({
    serviceConfig,
    organization: command.organization ?? command.defaultOrganization,
  });

  if (!loginSettings) {
    return { error: t("errors.couldNotSendResetLink") };
  }

  const searchResult = await searchUsers({
    serviceConfig,
    searchValue: command.loginName,
    organizationId: command.organization,
    loginSettings,
  });

  if (
    !searchResult ||
    !("result" in searchResult) ||
    !searchResult.result ||
    searchResult.result.length !== 1 ||
    !searchResult.result[0].userId
  ) {
    if (loginSettings?.ignoreUnknownUsernames) {
      await new Promise((resolve) => setTimeout(resolve, 2000));
      return {};
    }
    return { error: t("errors.couldNotSendResetLink") };
  }
  const user = searchResult.result[0];
  const humanUser = user.type.case === "human" ? user.type.value : undefined;

  const userLoginSettings = await getLoginSettings({ serviceConfig, organization: user.details?.resourceOwner });

  if (userLoginSettings?.disableLoginWithEmail && userLoginSettings?.disableLoginWithPhone) {
    if (user.preferredLoginName !== command.loginName) {
      if (userLoginSettings?.ignoreUnknownUsernames) {
        await new Promise((resolve) => setTimeout(resolve, 2000));
        return {};
      }
      return { error: t("errors.couldNotSendResetLink") };
    }
  } else if (userLoginSettings?.disableLoginWithEmail) {
    if (user.preferredLoginName !== command.loginName && humanUser?.phone?.phone !== command.loginName) {
      if (userLoginSettings?.ignoreUnknownUsernames) {
        await new Promise((resolve) => setTimeout(resolve, 2000));
        return {};
      }
      return { error: t("errors.couldNotSendResetLink") };
    }
  } else if (userLoginSettings?.disableLoginWithPhone) {
    if (user.preferredLoginName !== command.loginName && humanUser?.email?.email !== command.loginName) {
      if (userLoginSettings?.ignoreUnknownUsernames) {
        await new Promise((resolve) => setTimeout(resolve, 2000));
        return {};
      }
      return { error: t("errors.couldNotSendResetLink") };
    }
  }

  const userId = user.userId;
  const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

  return passwordReset({
    serviceConfig,
    userId,
    urlTemplate:
      `${hostWithProtocol}${basePath}/password/set?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}` +
      (command.requestId ? `&requestId=${command.requestId}` : ""),
  });
}

export type UpdateSessionCommand = {
  loginName: string;
  organization?: string;
  defaultOrganization?: string;
  checks: Checks;
  requestId?: string;
};

export async function sendPassword(
  command: UpdateSessionCommand,
): Promise<{ error: string } | { redirect: string } | { samlData: { url: string; fields: Record<string, string> } }> {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);
  const t = await getTranslations("password");

  let sessionCookie = await getSessionCookieByLoginName({
    loginName: command.loginName,
    organization: command.organization,
  });

  let session;
  let user: User | undefined;
  let loginSettingsByContext: LoginSettings | undefined;
  let loginSettingsByUser: LoginSettings | undefined;

  if (sessionCookie) {
    try {
      loginSettingsByUser = await getLoginSettings({ serviceConfig, organization: sessionCookie.organization });

      if (loginSettingsByUser) {
        let lifetime = loginSettingsByUser.passwordCheckLifetime;

        if (!lifetime || !lifetime.seconds) {
          console.warn("No password lifetime provided, defaulting to 24 hours");
          lifetime = {
            seconds: BigInt(60 * 60 * 24), // default to 24 hours
            nanos: 0,
          } as Duration;
        }

        session = await setSessionAndUpdateCookie({
          recentCookie: sessionCookie,
          checks: command.checks,
          requestId: command.requestId,
          lifetime,
        });
      } else {
        // Force fallback if settings can't be loaded
        throw new Error("Could not load login settings");
      }
    } catch {
      console.warn("[Password] Could not update session");
      // If the session was terminated or any other error occurred during update,
      // we fall back to creating a new session.
      sessionCookie = undefined;
      session = undefined;
    }
  }

  if (!sessionCookie) {
    if (!loginSettingsByContext) {
      loginSettingsByContext = await getLoginSettings({
        serviceConfig,
        organization: command.organization ?? command.defaultOrganization,
      });
    }

    // Force fallback if settings can't be loaded
    if (!loginSettingsByContext) {
      // this is a fake error message to hide that the user does not even exist
      return { error: t("errors.couldNotVerifyPassword") };
    }

    const searchResult = await searchUsers({
      serviceConfig,
      searchValue: command.loginName,
      organizationId: command.organization,
      loginSettings: loginSettingsByContext,
    });

    if (
      searchResult &&
      "result" in searchResult &&
      searchResult.result &&
      searchResult.result.length === 1 &&
      searchResult.result[0].userId
    ) {
      user = searchResult.result[0];
      const humanUser = user.type.case === "human" ? user.type.value : undefined;

      const userLoginSettings = await getLoginSettings({ serviceConfig, organization: user.details?.resourceOwner });

      // recheck login settings after user discovery, as the search might have been done without org scope
      if (userLoginSettings?.disableLoginWithEmail && userLoginSettings?.disableLoginWithPhone) {
        if (user.preferredLoginName !== command.loginName) {
          // emulate user not found to prevent enumeration (use context settings not user settings)
          if (loginSettingsByContext?.ignoreUnknownUsernames) {
            return { error: t("errors.failedToAuthenticateNoLimit") };
          }
          return { error: t("errors.couldNotVerifyPassword") };
        }
      } else if (userLoginSettings?.disableLoginWithEmail) {
        if (user.preferredLoginName !== command.loginName && humanUser?.phone?.phone !== command.loginName) {
          if (loginSettingsByContext?.ignoreUnknownUsernames) {
            return { error: t("errors.failedToAuthenticateNoLimit") };
          }
          return { error: t("errors.couldNotVerifyPassword") };
        }
      } else if (userLoginSettings?.disableLoginWithPhone) {
        if (user.preferredLoginName !== command.loginName && humanUser?.email?.email !== command.loginName) {
          if (loginSettingsByContext?.ignoreUnknownUsernames) {
            return { error: t("errors.failedToAuthenticateNoLimit") };
          }
          return { error: t("errors.couldNotVerifyPassword") };
        }
      }

      const checks = create(ChecksSchema, {
        user: { search: { case: "userId", value: user.userId } },
        password: { password: command.checks.password?.password },
      });

      try {
        const result = await createSessionAndUpdateCookie({
          checks,
          requestId: command.requestId,
          lifetime: loginSettingsByContext?.passwordCheckLifetime,
        });
        session = result.session;
        sessionCookie = result.sessionCookie;
      } catch (error: any) {
        if ("failedAttempts" in error && error.failedAttempts) {
          if (loginSettingsByContext?.ignoreUnknownUsernames) {
            return { error: t("errors.failedToAuthenticateNoLimit") };
          }
          const lockoutSettings = await getLockoutSettings({ serviceConfig, orgId: command.organization });

          const hasLimit =
            lockoutSettings?.maxPasswordAttempts !== undefined && lockoutSettings?.maxPasswordAttempts > BigInt(0);
          const locked = hasLimit && error.failedAttempts >= lockoutSettings?.maxPasswordAttempts;
          const messageKey = hasLimit ? "errors.failedToAuthenticate" : "errors.failedToAuthenticateNoLimit";

          return {
            error: t(messageKey, {
              failedAttempts: error.failedAttempts,
              maxPasswordAttempts: hasLimit ? (lockoutSettings?.maxPasswordAttempts).toString() : "?",
              lockoutMessage: locked ? t("errors.accountLockedContactAdmin") : "",
            }),
          };
        }
        if (loginSettingsByContext?.ignoreUnknownUsernames) {
          return { error: t("errors.failedToAuthenticateNoLimit") };
        }
        return { error: t("errors.couldNotCreateSessionForUser") };
      }
    } else {
      // this is a fake error message to hide that the user does not even exist
      if (loginSettingsByContext?.ignoreUnknownUsernames) {
        return { error: t("errors.failedToAuthenticateNoLimit") };
      }
      return { error: t("errors.couldNotVerifyPassword") };
    }
  }

  if (!session?.factors?.user?.id) {
    if (loginSettingsByContext?.ignoreUnknownUsernames) {
      return { error: t("errors.failedToAuthenticateNoLimit") };
    }
    return { error: t("errors.couldNotCreateSessionForUser") };
  }

  if (!user) {
    const userResponse = await getUserByID({ serviceConfig, userId: session?.factors?.user?.id });
    if (!userResponse.user) {
      return { error: t("errors.userNotFound") };
    }
    user = userResponse.user;
  }

  if (!session?.factors?.user?.id || !sessionCookie) {
    if (loginSettingsByContext?.ignoreUnknownUsernames) {
      return { error: t("errors.failedToAuthenticateNoLimit") };
    }
    return { error: t("errors.couldNotCreateSessionForUser") };
  }

  if (!loginSettingsByUser) {
    loginSettingsByUser = await getLoginSettings({
      serviceConfig,
      organization: command.organization ?? session.factors?.user?.organizationId ?? command.defaultOrganization,
    });
  }

  const humanUser = user.type.case === "human" ? user.type.value : undefined;

  const expirySettings = await getPasswordExpirySettings({
    serviceConfig,
    orgId: command.organization ?? session.factors?.user?.organizationId,
  });

  // check if the user has to change password first
  const passwordChangedCheck = checkPasswordChangeRequired(
    expirySettings,
    session,
    humanUser,
    command.organization,
    command.requestId,
  );

  if (passwordChangedCheck?.redirect) {
    return passwordChangedCheck;
  }

  // throw error if user is in initial state here and do not continue
  if (user.state === UserState.INITIAL) {
    return { error: t("errors.initialUserNotSupported") };
  }

  // check to see if user was verified
  const emailVerificationCheck = checkEmailVerification(session, humanUser, command.organization, command.requestId);

  if (emailVerificationCheck?.redirect) {
    return emailVerificationCheck;
  }

  // if password, check if user has MFA methods
  let authMethods;
  if (command.checks && command.checks.password && session.factors?.user?.id) {
    const response = await listAuthenticationMethodTypes({ serviceConfig, userId: session.factors.user.id });
    if (response.authMethodTypes && response.authMethodTypes.length) {
      authMethods = response.authMethodTypes;
    }
  }

  if (!authMethods) {
    return { error: t("errors.couldNotVerifyPassword") };
  }

  const mfaFactorCheck = await checkMFAFactors(
    serviceConfig,
    session,
    loginSettingsByUser,
    authMethods,
    command.organization,
    command.requestId,
  );

  if (mfaFactorCheck?.redirect) {
    return mfaFactorCheck;
  }

  let result: Awaited<ReturnType<typeof completeFlowOrGetUrl>>;
  
  if (command.requestId && session.id) {
    // OIDC/SAML flow
    console.log("Password auth: OIDC/SAML flow with requestId:", command.requestId, "sessionId:", session.id);
    result = await completeFlowOrGetUrl(
      {
        sessionId: session.id,
        requestId: command.requestId,
        organization: command.organization ?? session.factors?.user?.organizationId,
      },
      loginSettingsByUser?.defaultRedirectUri,
    );
  } else {
    // Regular flow (no requestId)
    console.log("Password auth: Regular flow with loginName:", session.factors.user.loginName);
    result = await completeFlowOrGetUrl(
      {
        loginName: session.factors.user.loginName,
        organization: session.factors?.user?.organizationId,
      },
      loginSettingsByUser?.defaultRedirectUri,
    );
  }

  if (result && typeof result === "object") {
    return result;
  }

  return { error: "Authentication completed but navigation failed" };
}

// this function lets users with code set a password or users with valid User Verification Check
export async function changePassword(command: { code?: string; userId: string; password: string; organization?: string }) {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);
  const t = await getTranslations("password");

  // check for init state
  const { user } = await getUserByID({ serviceConfig, userId: command.userId });

  if (!user || user.userId !== command.userId) {
    const loginSettings = await getLoginSettings({ serviceConfig, organization: command.organization });
    if (loginSettings?.ignoreUnknownUsernames) {
      return { error: t("set.errors.couldNotSetPassword") };
    }
    return { error: t("errors.couldNotSendResetLink") };
  }
  const userId = user.userId;

  if (user.state === UserState.INITIAL) {
    return { error: t("errors.userInitialStateNotSupported") };
  }

  // check if the user has no password set in order to set a password
  if (!command.code) {
    const authmethods = await listAuthenticationMethodTypes({ serviceConfig, userId });

    // if the user has no authmethods set, we need to check if the user was verified
    if (authmethods.authMethodTypes.length !== 0) {
      return {
        error: t("errors.codeOrVerificationRequired"),
      };
    }

    // check if a verification was done earlier
    const hasValidUserVerificationCheck = await checkUserVerification(user.userId);

    if (!hasValidUserVerificationCheck) {
      return { error: t("errors.verificationRequired") };
    }
  }

  return setUserPassword({ serviceConfig, userId, password: command.password, code: command.code });
}

type CheckSessionAndSetPasswordCommand = {
  sessionId: string;
  currentPassword: string;
  password: string;
};

export async function checkSessionAndSetPassword({
  sessionId,
  currentPassword,
  password,
}: CheckSessionAndSetPasswordCommand) {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);
  const t = await getTranslations("password");

  const sessionCookie = await getSessionCookieById({ sessionId });

  if (!sessionCookie) {
    return { error: "Could not load session cookie" };
  }

  let session;
  try {
    const sessionResponse = await getSession({
      serviceConfig,
      sessionId: sessionCookie.id,
      sessionToken: sessionCookie.token,
    });
    session = sessionResponse.session;
  } catch (error) {
    console.error("Error getting session:", error);
    return { error: "Could not load session" };
  }

  if (!session || !session.factors?.user?.id) {
    return { error: t("errors.couldNotLoadSession") };
  }

  const loginSettings = await getLoginSettings({
    serviceConfig,
    organization: sessionCookie.organization,
  });

  let lifetime = loginSettings?.passwordCheckLifetime;
  if (!lifetime || !lifetime.seconds) {
    lifetime = {
      seconds: BigInt(60 * 60 * 24),
      nanos: 0,
    } as Duration;
  }

  const checks = create(ChecksSchema, {
    password: { password: currentPassword },
  });

  try {
    await setSessionAndUpdateCookie({
      recentCookie: sessionCookie,
      checks,
      lifetime,
      requestId: sessionCookie.requestId,
    });
  } catch (error: any) {
    if ("failedAttempts" in error && error.failedAttempts) {
      if (loginSettings?.ignoreUnknownUsernames) {
        return { error: t("errors.failedToAuthenticateNoLimit") };
      }
      const lockoutSettings = await getLockoutSettings({ serviceConfig, orgId: sessionCookie.organization });

      const hasLimit =
        lockoutSettings?.maxPasswordAttempts !== undefined && lockoutSettings?.maxPasswordAttempts > BigInt(0);
      const locked = hasLimit && error.failedAttempts >= lockoutSettings?.maxPasswordAttempts;
      const messageKey = hasLimit ? "errors.failedToAuthenticate" : "errors.failedToAuthenticateNoLimit";

      return {
        error: t(messageKey, {
          failedAttempts: error.failedAttempts,
          maxPasswordAttempts: hasLimit ? (lockoutSettings?.maxPasswordAttempts).toString() : "?",
          lockoutMessage: locked ? t("errors.accountLockedContactAdmin") : "",
        }),
      };
    }
    if (loginSettings?.ignoreUnknownUsernames) {
      return { error: t("change.errors.couldNotVerifyPassword") };
    }
    return { error: t("change.errors.currentPasswordInvalid") };
  }

  const payload = create(SetPasswordRequestSchema, {
    userId: session.factors.user.id,
    newPassword: {
      password,
    },
  });

  return setPassword({ serviceConfig, payload }).catch((error) => {
    // throw error if failed precondition (ex. User is not yet initialized)
    if (error.code === 9 && error.message) {
      return { error: t("errors.failedPrecondition") };
    }
    return { error: "Could not set password" };
  });
}
