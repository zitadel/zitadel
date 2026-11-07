import { timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { PasswordExpirySettings } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";
import { HumanUser } from "@zitadel/proto/zitadel/user/v2/user_pb";
import moment from "moment";
import { trySendVerification } from "./server/verify";
import { hasValidUserVerificationCheckCookie } from "./verification-check";

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
    const passwordChangedDate = moment(timestampDate(humanUser.passwordChanged));
    const outdatedPassword = passwordChangedDate.add(maxAgeDays, "days");
    isOutdated = moment().isAfter(outdatedPassword);
  }

  if (humanUser?.passwordChangeRequired || isOutdated) {
    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
    });

    if (organization || session.factors?.user?.organizationId) {
      params.append("organization", session.factors?.user?.organizationId as string);
    }

    if (requestId) {
      params.append("requestId", requestId);
    }

    return { redirect: "/password/change?" + params };
  }
}

export async function checkEmailVerified(
  session: Session,
  humanUser?: HumanUser,
  organization?: string,
  requestId?: string,
) {
  if (!humanUser?.email?.isVerified) {
    const codeSent = await trySendVerification({
      userId: session.factors?.user?.id as string,
      isInvite: false,
      requestId,
    });

    const paramsVerify = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
      userId: session.factors?.user?.id as string, // verify needs user id
    });

    if (codeSent) {
      paramsVerify.append("codeSent", "true");
    }

    if (organization || session.factors?.user?.organizationId) {
      paramsVerify.append("organization", organization ?? (session.factors?.user?.organizationId as string));
    }

    if (requestId) {
      paramsVerify.append("requestId", requestId);
    }

    return { redirect: "/verify?" + paramsVerify };
  }
}

export async function checkEmailVerification(
  session: Session,
  humanUser?: HumanUser,
  organization?: string,
  requestId?: string,
) {
  if (!humanUser?.email?.isVerified && process.env.EMAIL_VERIFICATION === "true") {
    const codeSent = await trySendVerification({
      userId: session.factors?.user?.id as string,
      isInvite: false,
      requestId,
    });

    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName as string,
    });

    if (codeSent) {
      params.append("codeSent", "true");
    }

    if (requestId) {
      params.append("requestId", requestId);
    }

    if (organization || session.factors?.user?.organizationId) {
      params.append("organization", organization ?? (session.factors?.user?.organizationId as string));
    }

    return { redirect: `/verify?` + params };
  }
}

// Re-export MFA helpers for backward compatibility
export { checkMFAFactors, shouldEnforceMFA } from "./mfa-helper";

export async function checkUserVerification(userId: string): Promise<boolean> {
  // check if a verification was done earlier in this browser. The cookie value is an HMAC keyed with a
  // server-only secret (see verification-check.ts), so it can't be computed by the client.
  const valid = await hasValidUserVerificationCheckCookie(userId);

  if (!valid) {
    console.warn("User verification check failed.");
  }

  return valid;
}
