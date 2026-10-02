"use server";

import { createSessionAndUpdateCookie, createSessionForIdpAndUpdateCookie } from "@/lib/server/cookie";
import {
  addHumanUser,
  addIDPLink,
  getLoginSettings,
  getUserByID,
  listAuthenticationMethodTypes,
  listUsers,
  ServiceConfig,
  setUserPassword,
} from "@/lib/zitadel";
import { Code, ConnectError, create, Duration } from "@zitadel/client";
import { Factors } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { Checks, ChecksJson, ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import crypto from "crypto";
import { getTranslations } from "next-intl/server";
import { cookies, headers } from "next/headers";
import { completeFlowOrGetUrl } from "../client";
import { getOrSetFingerprintId } from "../fingerprint";
import { createLogger } from "../logger";
import {
  createRegistrationMarker,
  readRegistrationMarker,
  REGISTRATION_MARKER_COOKIE_NAME,
  REGISTRATION_MARKER_LIFETIME_SECONDS,
} from "../registration-marker";
import { getServiceConfig } from "../service-url";
import { checkEmailVerification, checkMFAFactors } from "../verify-helper";

const logger = createLogger("register");

const MAX_SESSION_RETRIES = 3;
const RETRY_DELAYS_MS = [500, 1000, 2000];

/**
 * After user creation, backend projections (users, login_names) may not be updated yet.
 * This helper retries createSessionAndUpdateCookie on NotFound errors with increasing delays.
 */
async function createSessionWithRetry(command: { checks: Checks; requestId: string | undefined; lifetime?: Duration }) {
  let lastError: unknown;
  for (let attempt = 0; attempt < MAX_SESSION_RETRIES; attempt++) {
    try {
      return await createSessionAndUpdateCookie(command);
    } catch (error) {
      lastError = error;
      const isNotFound = error instanceof ConnectError && error.code === Code.NotFound;
      const isLastAttempt = attempt + 1 >= MAX_SESSION_RETRIES;
      if (!isNotFound || isLastAttempt) {
        throw error;
      }
      const delay = RETRY_DELAYS_MS[attempt] ?? 2000;
      logger.warn(
        `Session creation failed with NotFound (attempt ${attempt + 1}/${MAX_SESSION_RETRIES}), retrying in ${delay}ms...`,
      );
      await new Promise((resolve) => setTimeout(resolve, delay));
    }
  }
  throw lastError;
}

/**
 * Sets the registration marker cookie for a user this request just created on the passkey path.
 * It is never issued again on a resume, so the 15-minute window cannot be extended.
 * Without a configured signing secret no marker is issued, so the registration cannot be resumed.
 */
async function issueRegistrationMarker(userId: string, organization: string) {
  const marker = createRegistrationMarker(userId, organization);
  if (!marker) {
    logger.warn("No signing secret configured, an unfinished registration cannot be resumed");
    return;
  }

  const cookiesList = await cookies();
  await cookiesList.set({
    name: REGISTRATION_MARKER_COOKIE_NAME,
    value: marker,
    httpOnly: true,
    path: "/",
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    maxAge: REGISTRATION_MARKER_LIFETIME_SECONDS,
  });
}

/**
 * The passkey path of `registerUser` creates the user before any passkey exists. If the person
 * leaves `/passkey/set` without registering one (or changes their mind and picks a password),
 * the next registration attempt for the same e-mail fails because the user already exists.
 *
 * `registerUser` continues with that user only when ALL of the following hold:
 *  1. `addHumanUser` failed with `AlreadyExists` (checked by the caller),
 *  2. exactly one human user with this e-mail as username exists in the organization,
 *  3. its e-mail address is not verified,
 *  4. it has no authentication method at all (no password, passkey or IdP link),
 *  5. the request carries a registration marker for exactly this user id and organization that has
 *     not expired and is signed with the login app's signing secret. Only `registerUser` issues it,
 *     once, when it created this account on the passkey path; it is valid for 15 minutes from then
 *     and a resume does not extend it.
 *
 * An account an administrator created and invited can be in the same state as 2–4. Condition 5
 * limits the resume to a request holding the marker issued when this account was registered.
 *
 * The marker is checked first, so without a valid marker no user lookup is made.
 *
 * Returns the id of the user (conditions 2–5), otherwise undefined.
 */
async function resumableRegistrationUserId(
  serviceConfig: ServiceConfig,
  command: { email: string; organization: string },
): Promise<string | undefined> {
  const cookiesList = await cookies();
  const markerUserId = readRegistrationMarker(cookiesList.get(REGISTRATION_MARKER_COOKIE_NAME)?.value, command.organization);
  if (!markerUserId) {
    return undefined;
  }

  try {
    const users = await listUsers({ serviceConfig, userName: command.email, organizationId: command.organization });
    if (users.result.length !== 1) {
      return undefined;
    }

    const user = users.result[0];
    if (user.userId !== markerUserId || user.type.case !== "human" || user.type.value.email?.isVerified) {
      return undefined;
    }

    const authMethods = await listAuthenticationMethodTypes({ serviceConfig, userId: user.userId });
    if (authMethods.authMethodTypes.length !== 0) {
      return undefined;
    }

    logger.info("Resuming an unfinished registration", { userId: user.userId });
    return user.userId;
  } catch (error) {
    logger.error("Failed to check for an unfinished registration", { error });
    return undefined;
  }
}

type RegisterUserCommand = {
  email: string;
  firstName: string;
  lastName: string;
  password?: string;
  organization: string;
  requestId?: string;
};

export type RegisterUserResponse = {
  userId: string;
  sessionId: string;
  factors: Factors | undefined;
};
export async function registerUser(
  command: RegisterUserCommand,
): Promise<{ error: string } | { redirect: string } | { samlData: { url: string; fields: Record<string, string> } }> {
  const t = await getTranslations("register");
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const loginSettings = await getLoginSettings({ serviceConfig, organization: command.organization });

  if (!loginSettings) {
    return { error: t("errors.couldNotGetLoginSettings") };
  }

  if (!loginSettings.allowRegister) {
    return { error: t("errors.registerNotAllowed") };
  }

  if (command.password && !loginSettings.allowLocalAuthentication) {
    return { error: t("errors.localAuthenticationNotAllowed") };
  }

  let userId: string | undefined;
  // true only if this call created the user (not when resuming an unfinished registration)
  let created = false;
  try {
    const addResponse = await addHumanUser({
      serviceConfig,
      email: command.email,
      firstName: command.firstName,
      lastName: command.lastName,
      password: command.password ? command.password : undefined,
      organization: command.organization,
    });
    userId = addResponse.userId;
    created = true;
  } catch (error) {
    logger.error("Failed to create user", { error });

    if (error instanceof ConnectError && error.code === Code.AlreadyExists) {
      userId = await resumableRegistrationUserId(serviceConfig, command);
    }

    if (!userId) {
      return { error: t("errors.couldNotCreateUser") };
    }

    if (command.password) {
      const passwordSet = await setUserPassword({ serviceConfig, userId, password: command.password })
        .then((response) => !("error" in response))
        .catch((error) => {
          logger.error("Failed to set the password on a resumed registration", { error });
          return false;
        });

      if (!passwordSet) {
        return { error: t("errors.couldNotCreateUser") };
      }

      // the account has a password now, the marker has served its purpose
      const cookiesList = await cookies();
      cookiesList.delete(REGISTRATION_MARKER_COOKIE_NAME);
    }
  }

  let checkPayload: any = {
    user: { search: { case: "userId", value: userId } },
  };

  if (command.password) {
    checkPayload = {
      ...checkPayload,
      password: { password: command.password },
    } as ChecksJson;
  }

  const checks = create(ChecksSchema, checkPayload);

  const result = await createSessionWithRetry({
    checks,
    requestId: command.requestId,
    lifetime: command.password ? loginSettings?.passwordCheckLifetime : undefined,
  }).catch((error) => {
    logger.error("Failed to create session after user creation", { error });
    return null;
  });

  const session = result?.session;

  if (!session || !session.factors?.user) {
    return { error: t("errors.couldNotCreateSession") };
  }

  if (!command.password) {
    const params = new URLSearchParams({
      loginName: session.factors.user.loginName,
      organization: session.factors.user.organizationId,
    });

    if (command.requestId) {
      params.append("requestId", command.requestId);
    }

    // Set verification cookie for users registering with passkey (no password)
    // This allows them to proceed with passkey registration without additional verification
    const cookiesList = await cookies();
    const userAgentId = await getOrSetFingerprintId();

    const verificationCheck = crypto.createHash("sha256").update(`${session.factors.user.id}:${userAgentId}`).digest("hex");

    await cookiesList.set({
      name: "verificationCheck",
      value: verificationCheck,
      httpOnly: true,
      path: "/",
      maxAge: 300, // 5 minutes
    });

    // Allows this browser to resume the registration if the passkey is not set up (see resumableRegistrationUserId).
    // Only issued when the user was created now; a resume keeps the original marker and its expiry.
    if (created) {
      await issueRegistrationMarker(session.factors.user.id, command.organization);
    }

    return { redirect: "/passkey/set?" + params };
  } else {
    const userResponse = await getUserByID({ serviceConfig, userId: session?.factors?.user?.id }).catch((error) => {
      logger.error("Failed to get user after session creation", { error });
      return null;
    });

    if (!userResponse?.user) {
      return { error: t("errors.userNotFound") };
    }

    const humanUser = userResponse.user.type.case === "human" ? userResponse.user.type.value : undefined;

    const emailVerificationCheck = await checkEmailVerification(
      session,
      humanUser,
      session.factors.user.organizationId,
      command.requestId,
    );

    if (emailVerificationCheck?.redirect) {
      return emailVerificationCheck;
    }

    return completeFlowOrGetUrl(
      command.requestId && session.id
        ? {
            sessionId: session.id,
            requestId: command.requestId,
            organization: session.factors.user.organizationId,
          }
        : {
            loginName: session.factors.user.loginName,
            organization: session.factors.user.organizationId,
          },
      loginSettings?.defaultRedirectUri,
    );
  }
}

type RegisterUserAndLinkToIDPommand = {
  email: string;
  firstName: string;
  lastName: string;
  organization: string;
  requestId?: string;
  idpIntent: {
    idpIntentId: string;
    idpIntentToken: string;
  };
  idpUserId: string;
  idpId: string;
  idpUserName: string;
};

export type registerUserAndLinkToIDPResponse = {
  userId: string;
  sessionId: string;
  factors: Factors | undefined;
};
export async function registerUserAndLinkToIDP(
  command: RegisterUserAndLinkToIDPommand,
): Promise<{ error: string } | { redirect: string } | { samlData: { url: string; fields: Record<string, string> } }> {
  const t = await getTranslations("register");

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const loginSettings = await getLoginSettings({ serviceConfig, organization: command.organization });

  if (!loginSettings) {
    return { error: t("errors.couldNotGetLoginSettings") };
  }

  if (!loginSettings.allowRegister) {
    return { error: t("errors.registerNotAllowed") };
  }

  const addUserResponse = await addHumanUser({
    serviceConfig,
    email: command.email,
    firstName: command.firstName,
    lastName: command.lastName,
    organization: command.organization,
  });

  const idpLink = await addIDPLink({
    serviceConfig,
    idp: {
      id: command.idpId,
      userId: command.idpUserId,
      userName: command.idpUserName,
    },
    userId: addUserResponse.userId,
  });

  if (!idpLink) {
    return { error: t("errors.couldNotLinkIDP") };
  }

  const session = await createSessionForIdpAndUpdateCookie({
    requestId: command.requestId,
    userId: addUserResponse.userId, // the user we just created
    idpIntent: command.idpIntent,
    lifetime: loginSettings?.externalLoginCheckLifetime,
  });

  if (!session || !session.factors?.user) {
    return { error: t("errors.couldNotCreateSession") };
  }

  // const userResponse = await getUserByID({
  //   serviceConfig.baseUrl,
  //   userId: session?.factors?.user?.id,
  // });

  // if (!userResponse.user) {
  //   return { error: "User not found in the system" };
  // }

  // const humanUser = userResponse.user.type.case === "human" ? userResponse.user.type.value : undefined;

  // check to see if user was verified
  // const emailVerificationCheck = checkEmailVerification(session, humanUser, command.organization, command.requestId);

  // if (emailVerificationCheck?.redirect) {
  //   return emailVerificationCheck;
  // }

  // check if user has MFA methods
  let authMethods;
  if (session.factors?.user?.id) {
    const response = await listAuthenticationMethodTypes({ serviceConfig, userId: session.factors.user.id });
    if (response.authMethodTypes && response.authMethodTypes.length) {
      authMethods = response.authMethodTypes;
    }
  }

  // Always check MFA factors, even if no auth methods are configured
  // This ensures that force MFA settings are respected
  const mfaFactorCheck = await checkMFAFactors(
    serviceConfig,
    session,
    loginSettings,
    authMethods || [], // Pass empty array if no auth methods
    command.organization,
    command.requestId,
  );

  if (mfaFactorCheck?.redirect) {
    return mfaFactorCheck;
  }

  return completeFlowOrGetUrl(
    command.requestId && session.id
      ? {
          sessionId: session.id,
          requestId: command.requestId,
          organization: session.factors.user.organizationId,
        }
      : {
          loginName: session.factors.user.loginName,
          organization: session.factors.user.organizationId,
        },
    loginSettings?.defaultRedirectUri,
  );
}
