"use server";

import { createSessionAndUpdateCookie, setSessionAndUpdateCookie } from "@/lib/server/cookie";
import {
  getLockoutSettings,
  getLoginSettings,
  getPasswordExpirySettings,
  getSession,
  getUserByID,
  listAuthenticationMethodTypes,
  listUsers,
  passwordReset,
  setPassword,
  setUserPassword,
} from "@/lib/zitadel";
import { ConnectError, create, Duration } from "@zitadel/client";
import { createUserServiceClient } from "@zitadel/client/v2";
import { Checks, ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { User, UserState } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { AuthenticationMethodType, SetPasswordRequestSchema } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { headers } from "next/headers";
import { completeFlowOrGetUrl } from "../client";
import { getSessionCookieById, getSessionCookieByLoginName } from "../cookies";
import { getServiceUrlFromHeaders } from "../service-url";
import { getOriginalHostWithProtocol } from "./host";
import {
  checkEmailVerification,
  checkMFAFactors,
  checkPasswordChangeRequired,
  checkUserVerification,
} from "../verify-helper";
import { createServerTransport } from "../zitadel";

type ResetPasswordCommand = {
  loginName: string;
  organization?: string;
  requestId?: string;
};

export async function resetPassword(command: ResetPasswordCommand) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  // Get the original host that the user sees with protocol
  const hostWithProtocol = await getOriginalHostWithProtocol();

  const users = await listUsers({
    serviceUrl,
    loginName: command.loginName,
    organizationId: command.organization,
  });

  if (!users.details || users.details.totalResult !== BigInt(1) || !users.result[0].userId) {
    return { error: "Could not send Password Reset Link" };
  }
  const userId = users.result[0].userId;

  const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

  return passwordReset({
    serviceUrl,
    userId,
    urlTemplate:
      `${hostWithProtocol}${basePath}/password/set?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}` +
      (command.requestId ? `&requestId=${command.requestId}` : ""),
  });
}

export type UpdateSessionCommand = {
  loginName: string;
  organization?: string;
  checks: Checks;
  requestId?: string;
};

export async function sendPassword(command: UpdateSessionCommand): Promise<{ error: string } | { redirect: string }> {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  let sessionCookie = await getSessionCookieByLoginName({
    loginName: command.loginName,
    organization: command.organization,
  }).catch((error) => {
    console.warn("Ignored error:", error);
  });

  let session;
  let user: User;
  let loginSettings: LoginSettings | undefined;

  if (!sessionCookie) {
    const users = await listUsers({
      serviceUrl,
      loginName: command.loginName,
      organizationId: command.organization,
    });

    if (users.details?.totalResult == BigInt(1) && users.result[0].userId) {
      user = users.result[0];

      const checks = create(ChecksSchema, {
        user: { search: { case: "userId", value: users.result[0].userId } },
        password: { password: command.checks.password?.password },
      });

      loginSettings = await getLoginSettings({
        serviceUrl,
        organization: command.organization,
      });

      try {
        session = await createSessionAndUpdateCookie({
          checks,
          requestId: command.requestId,
          lifetime: loginSettings?.passwordCheckLifetime,
        });
      } catch (error: any) {
        if ("failedAttempts" in error && error.failedAttempts) {
          const lockoutSettings = await getLockoutSettings({
            serviceUrl,
            orgId: command.organization,
          });

          return {
            error:
              `Failed to authenticate. You had ${error.failedAttempts} of ${lockoutSettings?.maxPasswordAttempts} password attempts.` +
              (lockoutSettings?.maxPasswordAttempts && error.failedAttempts >= lockoutSettings?.maxPasswordAttempts
                ? "Contact your administrator to unlock your account"
                : ""),
          };
        }
        return { error: "Could not create session for user" };
      }
    } else {
      // this is a fake error message to hide that the user does not even exist
      return { error: "Could not verify password" };
    }
  } else {
    loginSettings = await getLoginSettings({
      serviceUrl,
      organization: sessionCookie.organization,
    });

    if (!loginSettings) {
      return { error: "Could not load login settings" };
    }

    let lifetime = loginSettings.passwordCheckLifetime;

    if (!lifetime) {
      console.warn("No password lifetime provided, defaulting to 24 hours");
      lifetime = {
        seconds: BigInt(60 * 60 * 24), // default to 24 hours
        nanos: 0,
      } as Duration;
    }

    try {
      session = await setSessionAndUpdateCookie({
        recentCookie: sessionCookie,
        checks: command.checks,
        requestId: command.requestId,
        lifetime,
      });
    } catch (error: any) {
      if ("failedAttempts" in error && error.failedAttempts) {
        const lockoutSettings = await getLockoutSettings({
          serviceUrl,
          orgId: command.organization,
        });

        return {
          error:
            `Failed to authenticate. You had ${error.failedAttempts} of ${lockoutSettings?.maxPasswordAttempts} password attempts.` +
            (lockoutSettings?.maxPasswordAttempts && error.failedAttempts >= lockoutSettings?.maxPasswordAttempts
              ? " Contact your administrator to unlock your account"
              : ""),
        };
      }
      throw error;
    }

    if (!session?.factors?.user?.id) {
      return { error: "Could not create session for user" };
    }

    const userResponse = await getUserByID({
      serviceUrl,
      userId: session?.factors?.user?.id,
    });

    if (!userResponse.user) {
      return { error: "User not found in the system" };
    }

    user = userResponse.user;
  }

  if (!loginSettings) {
    loginSettings = await getLoginSettings({
      serviceUrl,
      organization: command.organization ?? session.factors?.user?.organizationId,
    });
  }

  if (!session?.factors?.user?.id) {
    return { error: "Could not create session for user" };
  }

  const humanUser = user.type.case === "human" ? user.type.value : undefined;

  const expirySettings = await getPasswordExpirySettings({
    serviceUrl,
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
    return { error: "Initial User not supported" };
  }

  // check to see if user was verified
  const emailVerificationCheck = checkEmailVerification(session, humanUser, command.organization, command.requestId);

  if (emailVerificationCheck?.redirect) {
    return emailVerificationCheck;
  }

  // if password, check if user has MFA methods
  let authMethods;
  if (command.checks && command.checks.password && session.factors?.user?.id) {
    const response = await listAuthenticationMethodTypes({
      serviceUrl,
      userId: session.factors.user.id,
    });
    if (response.authMethodTypes && response.authMethodTypes.length) {
      authMethods = response.authMethodTypes;
    }
  }

  if (!authMethods) {
    return { error: "Could not verify password!" };
  }

  const mfaFactorCheck = await checkMFAFactors(
    serviceUrl,
    session,
    loginSettings,
    authMethods,
    command.organization,
    command.requestId,
  );

  if (mfaFactorCheck?.redirect) {
    return mfaFactorCheck;
  }

  if (command.requestId && session.id) {
    // OIDC/SAML flow - use completeFlowOrGetUrl for proper handling
    console.log("Password auth: OIDC/SAML flow with requestId:", command.requestId, "sessionId:", session.id);
    const result = await completeFlowOrGetUrl(
      {
        sessionId: session.id,
        requestId: command.requestId,
        organization: command.organization ?? session.factors?.user?.organizationId,
      },
      loginSettings?.defaultRedirectUri,
    );
    console.log("Password auth: OIDC/SAML flow result:", result);

    // Safety net - ensure we always return a valid object
    if (!result || typeof result !== "object" || (!("redirect" in result) && !("error" in result))) {
      console.error("Password auth: Invalid result from completeFlowOrGetUrl (OIDC/SAML):", result);
      return { error: "Authentication completed but navigation failed" };
    }

    return result;
  }

  // Regular flow (no requestId) - return URL for client-side navigation
  console.log("Password auth: Regular flow with loginName:", session.factors.user.loginName);
  const result = await completeFlowOrGetUrl(
    {
      loginName: session.factors.user.loginName,
      organization: session.factors?.user?.organizationId,
    },
    loginSettings?.defaultRedirectUri,
  );
  console.log("Password auth: Regular flow result:", result);

  // Safety net - ensure we always return a valid object
  if (!result || typeof result !== "object" || (!("redirect" in result) && !("error" in result))) {
    console.error("Password auth: Invalid result from completeFlowOrGetUrl:", result);
    return { error: "Authentication completed but navigation failed" };
  }

  return result;
}

// this function lets users with code set a password or users with valid User Verification Check
export async function changePassword(command: { code?: string; userId: string; password: string }) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  // check for init state
  const { user } = await getUserByID({
    serviceUrl,
    userId: command.userId,
  });

  if (!user || user.userId !== command.userId) {
    return { error: "Could not send Password Reset Link" };
  }
  const userId = user.userId;

  if (user.state === UserState.INITIAL) {
    return { error: "User Initial State is not supported" };
  }

  // check if the user has no password set in order to set a password
  if (!command.code) {
    const authmethods = await listAuthenticationMethodTypes({
      serviceUrl,
      userId,
    });

    // if the user has no authmethods set, we need to check if the user was verified
    if (authmethods.authMethodTypes.length !== 0) {
      return {
        error: "You have to provide a code or have a valid User Verification Check",
      };
    }

    // check if a verification was done earlier
    const hasValidUserVerificationCheck = await checkUserVerification(user.userId);

    if (!hasValidUserVerificationCheck) {
      return { error: "User Verification Check has to be done" };
    }
  }

  return setUserPassword({
    serviceUrl,
    userId,
    password: command.password,
    code: command.code,
  });
}

type CheckSessionAndSetPasswordCommand = {
  sessionId: string;
  password: string;
};

export async function checkSessionAndSetPassword({ sessionId, password }: CheckSessionAndSetPasswordCommand) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const sessionCookie = await getSessionCookieById({ sessionId });

  const { session } = await getSession({
    serviceUrl,
    sessionId: sessionCookie.id,
    sessionToken: sessionCookie.token,
  });

  if (!session || !session.factors?.user?.id) {
    return { error: "Could not load session" };
  }

  const payload = create(SetPasswordRequestSchema, {
    userId: session.factors.user.id,
    newPassword: {
      password,
    },
  });

  // check if the user has no password set in order to set a password
  const authmethods = await listAuthenticationMethodTypes({
    serviceUrl,
    userId: session.factors.user.id,
  });

  if (!authmethods) {
    return { error: "Could not load auth methods" };
  }

  const requiredAuthMethodsForForceMFA = [
    AuthenticationMethodType.OTP_EMAIL,
    AuthenticationMethodType.OTP_SMS,
    AuthenticationMethodType.TOTP,
    AuthenticationMethodType.U2F,
  ];

  const hasNoMFAMethods = requiredAuthMethodsForForceMFA.every((method) => !authmethods.authMethodTypes.includes(method));

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization: session.factors.user.organizationId,
  });

  const forceMfa = !!(loginSettings?.forceMfa || loginSettings?.forceMfaLocalOnly);

  // if the user has no MFA but MFA is enforced, we can set a password otherwise we use the token of the user
  if (forceMfa && hasNoMFAMethods) {
    return setPassword({ serviceUrl, payload }).catch((error) => {
      // throw error if failed precondition (ex. User is not yet initialized)
      if (error.code === 9 && error.message) {
        return { error: "Failed precondition" };
      } else {
        throw error;
      }
    });
  } else {
    const transport = async (serviceUrl: string, token: string) => {
      return createServerTransport(token, serviceUrl);
    };

    const myUserService = async (serviceUrl: string, sessionToken: string) => {
      const transportPromise = await transport(serviceUrl, sessionToken);
      return createUserServiceClient(transportPromise);
    };

    const selfService = await myUserService(serviceUrl, `${sessionCookie.token}`);

    return selfService
      .setPassword(
        {
          userId: session.factors.user.id,
          newPassword: { password, changeRequired: false },
        },
        {},
      )
      .catch((error: ConnectError) => {
        console.log(error);
        if (error.code === 7) {
          return { error: "Session is not valid." };
        }
        throw error;
      });
  }
}
