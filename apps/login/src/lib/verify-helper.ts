import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { HumanUser } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";

export function checkPasswordChangeRequired(
  session: Session,
  humanUser: HumanUser | undefined,
  organization?: string,
  authRequestId?: string,
) {
  if (humanUser?.passwordChangeRequired) {
    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
    });

    if (organization || session.factors?.user?.organizationId) {
      params.append(
        "organization",
        session.factors?.user?.organizationId as string,
      );
    }

    if (authRequestId) {
      params.append("authRequestId", authRequestId);
    }

    return { redirect: "/password/change?" + params };
  }
}

export function checkEmailVerification(
  session: Session,
  humanUser?: HumanUser,
  organization?: string,
  authRequestId?: string,
) {
  console.log(
    humanUser?.email,
    process.env.EMAIL_VERIFICATION,
    process.env.EMAIL_VERIFICATION === "true",
  );
  if (
    !humanUser?.email?.isVerified &&
    process.env.EMAIL_VERIFICATION === "true"
  ) {
    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
    });

    if (authRequestId) {
      params.append("authRequestId", authRequestId);
    }

    if (organization || session.factors?.user?.organizationId) {
      params.append(
        "organization",
        organization ?? (session.factors?.user?.organizationId as string),
      );
    }

    return { redirect: `/verify?` + params };
  }
}

export function checkMFAFactors(
  session: Session,
  loginSettings: LoginSettings | undefined,
  authMethods: AuthenticationMethodType[],
  organization?: string,
  authRequestId?: string,
) {
  const availableMultiFactors = authMethods?.filter(
    (m: AuthenticationMethodType) =>
      m !== AuthenticationMethodType.PASSWORD &&
      m !== AuthenticationMethodType.PASSKEY,
  );

  if (availableMultiFactors?.length == 1) {
    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
    });

    if (authRequestId) {
      params.append("authRequestId", authRequestId);
    }

    if (organization || session.factors?.user?.organizationId) {
      params.append(
        "organization",
        organization ?? (session.factors?.user?.organizationId as string),
      );
    }

    const factor = availableMultiFactors[0];
    // if passwordless is other method, but user selected password as alternative, perform a login
    if (factor === AuthenticationMethodType.TOTP) {
      return { redirect: `/otp/time-based?` + params };
    } else if (factor === AuthenticationMethodType.OTP_SMS) {
      return { redirect: `/otp/sms?` + params };
    } else if (factor === AuthenticationMethodType.OTP_EMAIL) {
      return { redirect: `/otp/email?` + params };
    } else if (factor === AuthenticationMethodType.U2F) {
      return { redirect: `/u2f?` + params };
    }
  } else if (availableMultiFactors?.length >= 1) {
    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
    });

    if (authRequestId) {
      params.append("authRequestId", authRequestId);
    }

    if (organization || session.factors?.user?.organizationId) {
      params.append(
        "organization",
        organization ?? (session.factors?.user?.organizationId as string),
      );
    }

    return { redirect: `/mfa?` + params };
  } else if (
    (loginSettings?.forceMfa || loginSettings?.forceMfaLocalOnly) &&
    !availableMultiFactors.length
  ) {
    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
      force: "true", // this defines if the mfa is forced in the settings
      checkAfter: "true", // this defines if the check is directly made after the setup
    });

    if (authRequestId) {
      params.append("authRequestId", authRequestId);
    }

    if (organization || session.factors?.user?.organizationId) {
      params.append(
        "organization",
        organization ?? (session.factors?.user?.organizationId as string),
      );
    }

    // TODO: provide a way to setup passkeys on mfa page?
    return { redirect: `/mfa/set?` + params };
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
}
