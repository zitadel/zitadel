"use server";

import { getSessionCookieById } from "@/lib/server/cookies";
import { getEnrollmentAuthorizationError } from "@/lib/server/enrollment-guard";
import { isClassifiedError } from "@/lib/server/grpc/interceptors/error-classification";
import { createLogger } from "@/lib/server/logger";
import { getServiceConfig } from "@/lib/server/service-url";
import { addOTPEmail, addOTPSMS, getSession } from "@/lib/server/zitadel";
import { Code } from "@connectrpc/connect";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

const logger = createLogger("otp");

type AddOTPCommand = {
  method: "email" | "sms";
  sessionId: string;
};

/**
 * Adds OTP via email or SMS as second factor of the user of the session.
 * This must only be called from an explicit user action (never while rendering a page),
 * so that a link alone can't change the authentication methods of a signed-in user.
 */
export async function addOTP(command: AddOTPCommand): Promise<{ error: string } | { success: true }> {
  const { serviceConfig } = getServiceConfig(await headers());
  const t = await getTranslations("otp");

  const sessionCookie = await getSessionCookieById({
    sessionId: command.sessionId,
  });

  if (!sessionCookie) {
    return { error: t("errors.couldNotGetSession") };
  }

  const session = await getSession({ serviceConfig, sessionId: sessionCookie.id, sessionToken: sessionCookie.token });
  const userId = session?.session?.factors?.user?.id;

  if (!session?.session || !userId) {
    return { error: t("errors.couldNotGetSession") };
  }

  // Enrollment must be authorized: a bare identify-only session (only the user factor set)
  // must not be able to attach a new authenticator to the account (GHSA-45f2-5q3r-xgg6).
  const enrollmentError = await getEnrollmentAuthorizationError({ serviceConfig, session: session.session, userId });
  if (enrollmentError) {
    return { error: enrollmentError };
  }

  try {
    if (command.method === "email") {
      await addOTPEmail({ serviceConfig, userId });
    } else if (command.method === "sms") {
      await addOTPSMS({ serviceConfig, userId });
    } else {
      return { error: t("errors.invalidMethod") };
    }
  } catch (error) {
    // the user asked for the factor and already has it (e.g. when coming back to this page)
    if (isClassifiedError(error) && error.code === Code.AlreadyExists) {
      return { success: true };
    }
    // the email address or phone number has to be verified first
    if (isClassifiedError(error) && error.code === Code.FailedPrecondition) {
      return { error: t(command.method === "email" ? "errors.emailNotVerified" : "errors.phoneNotVerified") };
    }
    logger.warn("Could not add OTP", { method: command.method, error });
    return { error: t(command.method === "email" ? "errors.couldNotAddEmail" : "errors.couldNotAddSms") };
  }

  return { success: true };
}
