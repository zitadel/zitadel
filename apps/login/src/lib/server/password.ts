"use server";

import {
  createSessionAndUpdateCookie,
  setSessionAndUpdateCookie,
} from "@/lib/server/cookie";
import {
  getUserByID,
  listAuthenticationMethodTypes,
  listUsers,
  passwordReset,
  setPassword,
} from "@/lib/zitadel";
import { create } from "@zitadel/client";
import {
  Checks,
  ChecksSchema,
} from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { User, UserState } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { headers } from "next/headers";
import { redirect } from "next/navigation";
import { getSessionCookieByLoginName } from "../cookies";

type ResetPasswordCommand = {
  loginName: string;
  organization?: string;
};

export async function resetPassword(command: ResetPasswordCommand) {
  const host = headers().get("host");

  const users = await listUsers({
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

  return passwordReset(userId, host);
}

export type UpdateSessionCommand = {
  loginName: string;
  organization?: string;
  checks: Checks;
  authRequestId?: string;
  forceMfa?: boolean;
};

export async function sendPassword(command: UpdateSessionCommand) {
  let sessionCookie = await getSessionCookieByLoginName({
    loginName: command.loginName,
    organization: command.organization,
  }).catch((error) => {
    console.warn("Ignored error:", error);
  });

  let session;
  let user: User;
  if (!sessionCookie) {
    const users = await listUsers({
      loginName: command.loginName,
      organizationId: command.organization,
    });

    if (users.details?.totalResult == BigInt(1) && users.result[0].userId) {
      user = users.result[0];

      const checks = create(ChecksSchema, {
        user: { search: { case: "userId", value: users.result[0].userId } },
        password: { password: command.checks.password?.password },
      });

      session = await createSessionAndUpdateCookie(
        checks,
        undefined,
        command.authRequestId,
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
    );

    if (!session?.factors?.user?.id) {
      return { error: "Could not create session for user" };
    }

    const userResponse = await getUserByID(session?.factors?.user?.id);

    if (!userResponse.user) {
      return { error: "Could not find user" };
    }

    user = userResponse.user;
  }

  if (!session?.factors?.user?.id || !sessionCookie) {
    return { error: "Could not create session for user" };
  }

  // if password, check if user has MFA methods
  let authMethods;
  if (command.checks && command.checks.password && session.factors?.user?.id) {
    const response = await listAuthenticationMethodTypes(
      session.factors.user.id,
    );
    if (response.authMethodTypes && response.authMethodTypes.length) {
      authMethods = response.authMethodTypes;
    }
  }

  const submitted = {
    sessionId: session.id,
    factors: session.factors,
    challenges: session.challenges,
    authMethods,
    userState: user.state,
  };

  if (
    !submitted ||
    !submitted.authMethods ||
    !submitted.factors?.user?.loginName
  ) {
    return { error: "Could not verify password!" };
  }

  const availableSecondFactors = submitted?.authMethods?.filter(
    (m: AuthenticationMethodType) =>
      m !== AuthenticationMethodType.PASSWORD &&
      m !== AuthenticationMethodType.PASSKEY,
  );

  if (availableSecondFactors?.length == 1) {
    const params = new URLSearchParams({
      loginName: submitted.factors?.user.loginName,
    });

    if (command.authRequestId) {
      params.append("authRequestId", command.authRequestId);
    }

    if (command.organization) {
      params.append("organization", command.organization);
    }

    const factor = availableSecondFactors[0];
    // if passwordless is other method, but user selected password as alternative, perform a login
    if (factor === AuthenticationMethodType.TOTP) {
      return redirect(`/otp/time-based?` + params);
    } else if (factor === AuthenticationMethodType.OTP_SMS) {
      return redirect(`/otp/sms?` + params);
    } else if (factor === AuthenticationMethodType.OTP_EMAIL) {
      return redirect(`/otp/email?` + params);
    } else if (factor === AuthenticationMethodType.U2F) {
      return redirect(`/u2f?` + params);
    }
  } else if (availableSecondFactors?.length >= 1) {
    const params = new URLSearchParams({
      loginName: submitted.factors.user.loginName,
    });

    if (command.authRequestId) {
      params.append("authRequestId", command.authRequestId);
    }

    if (command.organization) {
      params.append("organization", command.organization);
    }

    return redirect(`/mfa?` + params);
  } else if (submitted.userState === UserState.INITIAL) {
    const params = new URLSearchParams({
      loginName: submitted.factors.user.loginName,
    });

    if (command.authRequestId) {
      params.append("authRequestId", command.authRequestId);
    }

    if (command.organization) {
      params.append("organization", command.organization);
    }

    return redirect(`/password/change?` + params);
  } else if (command.forceMfa && !availableSecondFactors.length) {
    const params = new URLSearchParams({
      loginName: submitted.factors.user.loginName,
      force: "true", // this defines if the mfa is forced in the settings
      checkAfter: "true", // this defines if the check is directly made after the setup
    });

    if (command.authRequestId) {
      params.append("authRequestId", command.authRequestId);
    }

    if (command.organization) {
      params.append("organization", command.organization);
    }

    // TODO: provide a way to setup passkeys on mfa page?
    return redirect(`/mfa/set?` + params);
  }
  // TODO: implement passkey setup

  //  else if (
  //   submitted.factors &&
  //   !submitted.factors.webAuthN && // if session was not verified with a passkey
  //   promptPasswordless && // if explicitly prompted due policy
  //   !isAlternative // escaped if password was used as an alternative method
  // ) {
  //   const params = new URLSearchParams({
  //     loginName: submitted.factors.user.loginName,
  //     prompt: "true",
  //   });

  //   if (authRequestId) {
  //     params.append("authRequestId", authRequestId);
  //   }

  //   if (organization) {
  //     params.append("organization", organization);
  //   }

  //   return router.push(`/passkey/set?` + params);
  // }
  else if (command.authRequestId && submitted.sessionId) {
    const params = new URLSearchParams({
      sessionId: submitted.sessionId,
      authRequest: command.authRequestId,
    });

    if (command.organization) {
      params.append("organization", command.organization);
    }

    // move this to browser
    return { nextStep: `/login?${params}` };
  }

  // without OIDC flow
  const params = new URLSearchParams(
    command.authRequestId
      ? {
          loginName: submitted.factors.user.loginName,
          authRequestId: command.authRequestId,
        }
      : {
          loginName: submitted.factors.user.loginName,
        },
  );

  if (command.organization) {
    params.append("organization", command.organization);
  }

  return redirect(`/signedin?` + params);
}

export async function changePassword(command: {
  code?: string;
  userId: string;
  password: string;
}) {
  // check for init state
  const { user } = await getUserByID(command.userId);

  if (!user || user.userId !== command.userId) {
    return { error: "Could not send Password Reset Link" };
  }
  const userId = user.userId;

  return setPassword(userId, command.password, user, command.code);
}
