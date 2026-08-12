"use server";

import { getSession, registerU2F, verifyU2FRegistration } from "@/lib/zitadel";
import { create } from "@zitadel/client";
import { VerifyU2FRegistrationRequestSchema } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import { userAgent } from "next/server";
import { getSessionCookieById } from "../cookies";
import { getServiceConfig } from "../service-url";
import { getEnrollmentAuthorizationError } from "./enrollment-guard";
import { catchUserError } from "./error-utils";
import { getPublicHost } from "./host";

type RegisterU2FCommand = {
  sessionId: string;
};

type VerifyU2FCommand = {
  u2fId: string;
  passkeyName?: string;
  publicKeyCredential: any;
  sessionId: string;
};

export async function addU2F(command: RegisterU2FCommand) {
  const t = await getTranslations("u2f");
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);
  const host = getPublicHost(_headers);

  const sessionCookie = await getSessionCookieById({
    sessionId: command.sessionId,
  });

  if (!sessionCookie) {
    return { error: t("errors.couldNotLoadSession") };
  }

  let session;
  try {
    session = await getSession({ serviceConfig, sessionId: sessionCookie.id, sessionToken: sessionCookie.token });
  } catch (error) {
    // the cookie can outlive the server-side session — an expected user state
    return catchUserError(error, t("errors.couldNotLoadSession"));
  }

  const [hostname] = host.split(":");

  if (!hostname) {
    throw new Error("Could not get hostname");
  }

  const userId = session?.session?.factors?.user?.id;

  if (!session || !userId) {
    return { error: t("errors.couldNotLoadSession") };
  }

  // Enrollment must be authorized: a bare identify-only session (only the user factor set)
  // must not be able to attach a new authenticator to the account (GHSA-45f2-5q3r-xgg6).
  const enrollmentError = await getEnrollmentAuthorizationError({ serviceConfig, session: session.session!, userId });
  if (enrollmentError) {
    return { error: enrollmentError };
  }

  return registerU2F({ serviceConfig, userId, domain: hostname }).catch((error) =>
    catchUserError(error, t("errors.couldNotRegister")),
  );
}

export async function verifyU2F(command: VerifyU2FCommand) {
  const t = await getTranslations("u2f");
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);
  let passkeyName = command.passkeyName;
  if (!passkeyName) {
    const headersList = await headers();
    const userAgentStructure = { headers: headersList };
    const { browser, device, os } = userAgent(userAgentStructure);

    passkeyName = `${device.vendor ?? ""} ${device.model ?? ""}${
      device.vendor || device.model ? ", " : ""
    }${os.name}${os.name ? ", " : ""}${browser.name}`;
  }
  const sessionCookie = await getSessionCookieById({
    sessionId: command.sessionId,
  });

  if (!sessionCookie) {
    return { error: t("errors.couldNotLoadSession") };
  }

  let session;
  try {
    session = await getSession({ serviceConfig, sessionId: sessionCookie.id, sessionToken: sessionCookie.token });
  } catch (error) {
    return catchUserError(error, t("errors.couldNotLoadSession"));
  }

  const userId = session?.session?.factors?.user?.id;

  if (!userId) {
    return { error: t("errors.couldNotLoadSession") };
  }

  // Enrollment must be authorized: only an authenticated session (or a valid onboarding
  // verification) may finish registering a new authenticator (GHSA-45f2-5q3r-xgg6).
  const enrollmentError = await getEnrollmentAuthorizationError({ serviceConfig, session: session.session!, userId });
  if (enrollmentError) {
    return { error: enrollmentError };
  }

  const request = create(VerifyU2FRegistrationRequestSchema, {
    u2fId: command.u2fId,
    publicKeyCredential: command.publicKeyCredential,
    tokenName: passkeyName,
    userId,
  });

  // a failed WebAuthn attestation is a user/browser-side failure, not a server fault
  return verifyU2FRegistration({ serviceConfig, request }).catch((error) =>
    catchUserError(error, t("errors.couldNotVerify")),
  );
}
