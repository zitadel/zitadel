import { timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { PasswordExpirySettings } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";
import { HumanUser } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import crypto from "crypto";
import moment from "moment";
import { cookies } from "next/headers";
import { getFingerprintIdCookie } from "./fingerprint";
import { getUserByID } from "./zitadel";

export function checkPasswordChangeRequired(
  expirySettings: PasswordExpirySettings | undefined,
  session: Session,
  humanUser: HumanUser | undefined,
  organization?: string,
  requestId?: string,
) {
  let isOutdated = false;
  if (expirySettings?.maxAgeDays && humanUser?.passwordChanged) {
    const maxAgeDays = Number(expirySettings.maxAgeDays); // Convert bigint to number
    const passwordChangedDate = moment(
      timestampDate(humanUser.passwordChanged),
    );
    const outdatedPassword = passwordChangedDate.add(maxAgeDays, "days");
    isOutdated = moment().isAfter(outdatedPassword);
  }

  if (humanUser?.passwordChangeRequired || isOutdated) {
    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
    });

    if (organization || session.factors?.user?.organizationId) {
      params.append(
        "organization",
        session.factors?.user?.organizationId as string,
      );
    }

    if (requestId) {
      params.append("requestId", requestId);
    }

    return { redirect: "/password/change?" + params };
  }
}

export function checkEmailVerified(
  session: Session,
  humanUser?: HumanUser,
  organization?: string,
  requestId?: string,
) {
  if (!humanUser?.email?.isVerified) {
    const paramsVerify = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
      userId: session.factors?.user?.id as string, // verify needs user id
      send: "true", // we request a new email code once the page is loaded
    });

    if (organization || session.factors?.user?.organizationId) {
      paramsVerify.append(
        "organization",
        organization ?? (session.factors?.user?.organizationId as string),
      );
    }

    if (requestId) {
      paramsVerify.append("requestId", requestId);
    }

    return { redirect: "/verify?" + paramsVerify };
  }
}

export function checkEmailVerification(
  session: Session,
  humanUser?: HumanUser,
  organization?: string,
  requestId?: string,
) {
  if (
    !humanUser?.email?.isVerified &&
    process.env.EMAIL_VERIFICATION === "true"
  ) {
    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
      send: "true", // set this to true as we dont expect old email codes to be valid anymore
    });

    if (requestId) {
      params.append("requestId", requestId);
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

export async function checkMFAFactors(
  serviceUrl: string,
  session: Session,
  loginSettings: LoginSettings | undefined,
  authMethods: AuthenticationMethodType[],
  organization?: string,
  requestId?: string,
) {
  const availableMultiFactors = authMethods?.filter(
    (m: AuthenticationMethodType) =>
      m !== AuthenticationMethodType.PASSWORD &&
      m !== AuthenticationMethodType.PASSKEY,
  );

  const hasAuthenticatedWithPasskey =
    session.factors?.webAuthN?.verifiedAt &&
    session.factors?.webAuthN?.userVerified;

  // escape further checks if user has authenticated with passkey
  if (hasAuthenticatedWithPasskey) {
    return;
  }

  // if user has not authenticated with passkey and has only one additional mfa factor, redirect to that
  if (availableMultiFactors?.length == 1) {
    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
    });

    if (requestId) {
      params.append("requestId", requestId);
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
  } else if (availableMultiFactors?.length > 1) {
    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
    });

    if (requestId) {
      params.append("requestId", requestId);
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

    if (requestId) {
      params.append("requestId", requestId);
    }

    if (organization || session.factors?.user?.organizationId) {
      params.append(
        "organization",
        organization ?? (session.factors?.user?.organizationId as string),
      );
    }

    // TODO: provide a way to setup passkeys on mfa page?
    return { redirect: `/mfa/set?` + params };
  } else if (
    loginSettings?.mfaInitSkipLifetime &&
    (loginSettings.mfaInitSkipLifetime.nanos > 0 ||
      loginSettings.mfaInitSkipLifetime.seconds > 0) &&
    !availableMultiFactors.length &&
    session?.factors?.user?.id
  ) {
    const userResponse = await getUserByID({
      serviceUrl,
      userId: session.factors?.user?.id,
    });

    const humanUser =
      userResponse?.user?.type.case === "human"
        ? userResponse?.user.type.value
        : undefined;

    if (humanUser?.mfaInitSkipped) {
      const mfaInitSkippedTimestamp = timestampDate(humanUser.mfaInitSkipped);

      const mfaInitSkipLifetimeMillis =
        Number(loginSettings.mfaInitSkipLifetime.seconds) * 1000 +
        loginSettings.mfaInitSkipLifetime.nanos / 1000000;
      const currentTime = Date.now();
      const mfaInitSkippedTime = mfaInitSkippedTimestamp.getTime();
      const timeDifference = currentTime - mfaInitSkippedTime;

      if (!(timeDifference > mfaInitSkipLifetimeMillis)) {
        // if the time difference is smaller than the lifetime, skip the mfa setup
        return;
      }
    }

    // the user has never skipped the mfa init but we have a setting so we redirect

    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
      force: "false", // this defines if the mfa is not forced in the settings and can be skipped
      checkAfter: "true", // this defines if the check is directly made after the setup
    });

    if (requestId) {
      params.append("requestId", requestId);
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
}

export async function checkUserVerification(userId: string): Promise<boolean> {
  // check if a verification was done earlier
  const cookiesList = await cookies();

  // only read cookie to prevent issues on page.tsx
  const fingerPrintCookie = await getFingerprintIdCookie();

  if (!fingerPrintCookie || !fingerPrintCookie.value) {
    return false;
  }

  const verificationCheck = crypto
    .createHash("sha256")
    .update(`${userId}:${fingerPrintCookie.value}`)
    .digest("hex");

  const cookieValue = await cookiesList.get("verificationCheck")?.value;

  if (!cookieValue) {
    console.warn(
      "User verification check cookie not found. User verification check failed.",
    );
    return false;
  }

  if (cookieValue !== verificationCheck) {
    console.warn(
      `User verification check failed. Expected ${verificationCheck} but got ${cookieValue}`,
    );
    return false;
  }

  return true;
}
