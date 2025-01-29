"use server";

import {
  createSessionAndUpdateCookie,
  setSessionAndUpdateCookie,
} from "@/lib/server/cookie";
import {
  getLoginSettings,
  getSession,
  getUserByID,
  listAuthenticationMethodTypes,
  listUsers,
  passwordReset,
  setPassword,
  setUserPassword,
} from "@/lib/zitadel";
import { create } from "@zitadel/client";
import { createServerTransport } from "@zitadel/client/node";
import { createUserServiceClient } from "@zitadel/client/v2";
import {
  Checks,
  ChecksSchema,
} from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { User, UserState } from "@zitadel/proto/zitadel/user/v2/user_pb";
import {
  AuthenticationMethodType,
  SetPasswordRequestSchema,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { headers } from "next/headers";
import { getNextUrl } from "../client";
import { getSessionCookieById, getSessionCookieByLoginName } from "../cookies";
import { getApiUrlOfHeaders } from "../service";
import {
  checkEmailVerification,
  checkMFAFactors,
  checkPasswordChangeRequired,
} from "../verify-helper";

type ResetPasswordCommand = {
  loginName: string;
  organization?: string;
  authRequestId?: string;
};

export async function resetPassword(command: ResetPasswordCommand) {
  const _headers = await headers();
  const serviceUrl = getApiUrlOfHeaders(_headers);
  const host = _headers.get("host");

  if (!host || typeof host !== "string") {
    throw new Error("No host found");
  }

  const users = await listUsers({
    serviceUrl,
    loginName: command.loginName,
    organizationId: command.organization,
  });

  if (
    !users.details ||
    users.details.totalResult !== BigInt(1) ||
    !users.result[0].userId
  ) {
    return { error: "Could not send Password Reset Link" };
  }
  const userId = users.result[0].userId;

  return passwordReset({
    serviceUrl,
    userId,
    urlTemplate:
      `${host.includes("localhost") ? "http://" : "https://"}${host}/password/set?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}` +
      (command.authRequestId ? `&authRequestId=${command.authRequestId}` : ""),
  });
}

export type UpdateSessionCommand = {
  loginName: string;
  organization?: string;
  checks: Checks;
  authRequestId?: string;
};

export async function sendPassword(command: UpdateSessionCommand) {
  const _headers = await headers();
  const serviceUrl = getApiUrlOfHeaders(_headers);

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

      session = await createSessionAndUpdateCookie(
        checks,
        undefined,
        command.authRequestId,
        loginSettings?.passwordCheckLifetime,
      );
    }

    // this is a fake error message to hide that the user does not even exist
    return { error: "Could not verify password" };
  } else {
    session = await setSessionAndUpdateCookie(
      sessionCookie,
      command.checks,
      undefined,
      command.authRequestId,
      loginSettings?.passwordCheckLifetime,
    );

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
      organization:
        command.organization ?? session.factors?.user?.organizationId,
    });
  }

  if (!session?.factors?.user?.id || !sessionCookie) {
    return { error: "Could not create session for user" };
  }

  const humanUser = user.type.case === "human" ? user.type.value : undefined;

  // check if the user has to change password first
  const passwordChangedCheck = checkPasswordChangeRequired(
    session,
    humanUser,
    command.organization,
    command.authRequestId,
  );

  if (passwordChangedCheck?.redirect) {
    return passwordChangedCheck;
  }

  // throw error if user is in initial state here and do not continue
  if (user.state === UserState.INITIAL) {
    return { error: "Initial User not supported" };
  }

  // check to see if user was verified
  const emailVerificationCheck = checkEmailVerification(
    session,
    humanUser,
    command.organization,
    command.authRequestId,
  );

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

  const mfaFactorCheck = checkMFAFactors(
    session,
    loginSettings,
    authMethods,
    command.organization,
    command.authRequestId,
  );

  if (mfaFactorCheck?.redirect) {
    return mfaFactorCheck;
  }

  if (command.authRequestId && session.id) {
    const nextUrl = await getNextUrl(
      {
        sessionId: session.id,
        authRequestId: command.authRequestId,
        organization:
          command.organization ?? session.factors?.user?.organizationId,
      },
      loginSettings?.defaultRedirectUri,
    );

    return { redirect: nextUrl };
  }

  const url = await getNextUrl(
    {
      loginName: session.factors.user.loginName,
      organization: session.factors?.user?.organizationId,
    },
    loginSettings?.defaultRedirectUri,
  );

  return { redirect: url };
}

export async function changePassword(command: {
  code?: string;
  userId: string;
  password: string;
}) {
  const _headers = await headers();
  const serviceUrl = getApiUrlOfHeaders(_headers);

  // check for init state
  const { user } = await getUserByID({ serviceUrl, userId: command.userId });

  if (!user || user.userId !== command.userId) {
    return { error: "Could not send Password Reset Link" };
  }
  const userId = user.userId;

  return setUserPassword({
    serviceUrl,
    userId,
    password: command.password,
    user,
    code: command.code,
  });
}

type CheckSessionAndSetPasswordCommand = {
  sessionId: string;
  password: string;
};

export async function checkSessionAndSetPassword({
  sessionId,
  password,
}: CheckSessionAndSetPasswordCommand) {
  const _headers = await headers();
  const serviceUrl = getApiUrlOfHeaders(_headers);

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

  const hasNoMFAMethods = requiredAuthMethodsForForceMFA.every(
    (method) => !authmethods.authMethodTypes.includes(method),
  );

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization: session.factors.user.organizationId,
  });

  const forceMfa = !!(
    loginSettings?.forceMfa || loginSettings?.forceMfaLocalOnly
  );

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
    const transport = async (host: string, token: string) => {
      return createServerTransport(token, {
        baseUrl: serviceUrl,
      });
    };

    const myUserService = async (host: string, sessionToken: string) => {
      const transportPromise = await transport(serviceUrl, sessionToken);
      return createUserServiceClient(transportPromise);
    };

    const selfService = await myUserService(
      serviceUrl,
      `${sessionCookie.token}`,
    );

    return selfService
      .setPassword(
        {
          userId: session.factors.user.id,
          newPassword: { password, changeRequired: false },
        },
        {},
      )
      .catch((error) => {
        console.log(error);
        if (error.code === 7) {
          return { error: "Session is not valid." };
        }
        throw error;
      });
  }
}
