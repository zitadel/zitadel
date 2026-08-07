"use server";

import { createLogger } from "@/lib/logger";
import {
  createPasskeyRegistrationLink,
  getLoginSettings,
  getSession,
  getUserByID,
  registerPasskey,
  verifyPasskeyRegistration as zitadelVerifyPasskeyRegistration,
} from "@/lib/zitadel";
import { create, Duration } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { Checks, ChecksSchema, GetSessionResponse } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import {
  RegisterPasskeyResponse,
  VerifyPasskeyRegistrationRequestSchema,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import { userAgent } from "next/server";
import { completeFlowOrGetUrl } from "../client";
import { getSessionCookieById } from "../cookies";
import { getServiceConfig } from "../service-url";
import { checkEmailVerification } from "../verify-helper";
import { createSessionAndUpdateCookie } from "./cookie";
import { getEnrollmentAuthorizationError } from "./enrollment-guard";
import { catchUserError } from "./error-utils";
import { getPublicHost } from "./host";
import { updateOrCreateSession } from "./session";

const logger = createLogger("passkeys");

type VerifyPasskeyCommand = {
  passkeyId: string;
  passkeyName?: string;
  publicKeyCredential: any;
  sessionId?: string;
  userId?: string;
};

type RegisterPasskeyCommand = {
  sessionId?: string;
  userId?: string;
  code?: string;
  codeId?: string;
};

export async function registerPasskeyLink(
  command: RegisterPasskeyCommand,
): Promise<RegisterPasskeyResponse | { error: string }> {
  const t = await getTranslations("passkey");

  if (!command.sessionId && !command.userId) {
    return { error: t("set.errors.missingContext") };
  }

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);
  const host = getPublicHost(_headers);

  let session: GetSessionResponse | undefined;
  let createdSession: Session | undefined;
  let currentUserId: string | undefined = undefined;
  let registerCode: { id: string; code: string } | undefined = undefined;

  if (command.sessionId) {
    // Session-based flow (existing logic)
    const sessionCookie = await getSessionCookieById({ sessionId: command.sessionId });

    if (!sessionCookie) {
      return { error: t("set.errors.couldNotLoadSession") };
    }

    try {
      session = await getSession({ serviceConfig, sessionId: sessionCookie.id, sessionToken: sessionCookie.token });
    } catch (error) {
      // the cookie can outlive the server-side session — an expected user state
      return catchUserError(error, t("set.errors.couldNotLoadSession"));
    }

    if (!session?.session?.factors?.user?.id) {
      return { error: t("set.errors.couldNotLoadSession") };
    }

    currentUserId = session.session.factors.user.id;

    const enrollmentError = await getEnrollmentAuthorizationError({
      serviceConfig,
      session: session.session,
      userId: currentUserId,
    });

    if (enrollmentError) {
      return { error: enrollmentError };
    }

    // Generate registration code if not provided
    if (command.code && command.codeId) {
      registerCode = {
        id: command.codeId,
        code: command.code,
      };
    } else {
      let codeResponse;
      try {
        codeResponse = await createPasskeyRegistrationLink({ serviceConfig, userId: currentUserId });
      } catch (error) {
        return catchUserError(error, t("set.errors.couldNotRegisterPasskey"));
      }

      if (!codeResponse?.code?.code) {
        return { error: t("set.errors.couldNotRegisterPasskey") };
      }

      registerCode = codeResponse.code;
    }
  } else if (command.userId && command.code && command.codeId) {
    currentUserId = command.userId;
    registerCode = {
      id: command.codeId,
      code: command.code,
    };

    // Check if user exists
    let userResponse;
    try {
      userResponse = await getUserByID({ serviceConfig, userId: currentUserId });
    } catch (error) {
      return catchUserError(error, t("set.errors.userNotFound"));
    }

    if (!userResponse || !userResponse.user) {
      return { error: t("set.errors.userNotFound") };
    }

    // Create a session for the user to continue the flow after passkey registration
    const checks = create(ChecksSchema, {
      user: {
        search: {
          case: "loginName",
          value: userResponse.user.preferredLoginName,
        },
      },
    });

    let result;
    try {
      result = await createSessionAndUpdateCookie({
        checks,
        requestId: undefined, // No requestId in passkey registration context, TODO: consider if needed
      });
    } catch (error) {
      return catchUserError(error, t("set.errors.couldNotCreateSession"));
    }
    createdSession = result.session;

    if (!createdSession) {
      return { error: t("set.errors.couldNotCreateSession") };
    }
  }

  if (!registerCode) {
    return { error: t("set.errors.missingContext") };
  }

  const [hostname] = host.split(":");

  if (!hostname) {
    throw new Error("Could not get hostname");
  }

  if (!currentUserId) {
    return { error: t("set.errors.missingContext") };
  }

  return registerPasskey({ serviceConfig, userId: currentUserId, code: registerCode, domain: hostname }).catch((error) =>
    catchUserError(error, t("set.errors.couldNotRegisterPasskey")),
  );
}

export async function verifyPasskeyRegistration(command: VerifyPasskeyCommand) {
  const t = await getTranslations("passkey");
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  if (!command.sessionId && !command.userId) {
    return { error: t("set.errors.missingContext") };
  }

  // if no name is provided, try to generate one from the user agent
  let passkeyName = command.passkeyName;
  if (!passkeyName) {
    const headersList = await headers();
    const userAgentStructure = { headers: headersList };
    const { browser, device, os } = userAgent(userAgentStructure);

    passkeyName = `${device.vendor ?? ""} ${device.model ?? ""}${
      device.vendor || device.model ? ", " : ""
    }${os.name}${os.name ? ", " : ""}${browser.name}`;
  }

  let loginName: string | undefined;
  let currentUserId: string;

  if (command.sessionId) {
    // Session-based flow
    const sessionCookie = await getSessionCookieById({
      sessionId: command.sessionId,
    });

    if (!sessionCookie) {
      return { error: t("set.errors.couldNotLoadSession") };
    }

    let session;
    try {
      session = await getSession({ serviceConfig, sessionId: sessionCookie.id, sessionToken: sessionCookie.token });
    } catch (error) {
      // the cookie can outlive the server-side session — an expected user state
      return catchUserError(error, t("set.errors.couldNotLoadSession"));
    }
    const userId = session?.session?.factors?.user?.id;

    if (!userId) {
      return { error: t("set.errors.couldNotLoadSession") };
    }

    currentUserId = userId;
    loginName = session?.session?.factors?.user?.loginName;
  } else {
    // UserId-based flow
    currentUserId = command.userId!;

    // Verify user exists
    let userResponse;
    try {
      userResponse = await getUserByID({ serviceConfig, userId: currentUserId });
    } catch (error) {
      return catchUserError(error, t("set.errors.userNotFound"));
    }

    if (!userResponse || !userResponse.user) {
      return { error: t("set.errors.userNotFound") };
    }

    loginName = userResponse.user.preferredLoginName;
  }

  let response;
  try {
    response = await zitadelVerifyPasskeyRegistration({
      serviceConfig,
      request: create(VerifyPasskeyRegistrationRequestSchema, {
        passkeyId: command.passkeyId,
        publicKeyCredential: command.publicKeyCredential,
        passkeyName,
        userId: currentUserId,
      }),
    });
  } catch (error) {
    // a failed WebAuthn attestation is a user/browser-side failure, not a server fault
    return catchUserError(error, t("set.errors.couldNotVerifyPasskey"));
  }

  return { ...response, loginName };
}

type SendPasskeyCommand = {
  loginName?: string;
  sessionId?: string;
  organization?: string;
  checks?: Checks;
  requestId?: string;
  lifetime?: Duration;
};

export async function sendPasskey(command: SendPasskeyCommand) {
  let { loginName, sessionId, organization, checks, requestId } = command;

  const t = await getTranslations("passkey");

  const result = await updateOrCreateSession({
    loginName,
    sessionId,
    organization,
    checks,
    requestId,
    lifetime: command.lifetime,
  });

  if ("error" in result) {
    // try to interpret validation errors as translation keys if possible, or fallback to generic
    // For now returning the error string directly as key or default
    return { error: result.error };
  }

  // transformation to partial session for compatibility
  const session = {
    id: result.sessionId,
    factors: result.factors,
    // @ts-ignore
    challenges: result.challenges,
  };

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);
  // the passkey check already succeeded — a failed settings lookup must not
  // fail the login, it only loses the org's default redirect
  const loginSettings = await getLoginSettings({ serviceConfig, organization }).catch((error) => {
    logger.warn("Could not load login settings after passkey check", { error });
    return undefined;
  });

  const userId = session?.factors?.user?.id;
  if (!userId) {
    return { error: t("verify.errors.couldNotFindSession") };
  }

  let userResponse;
  try {
    userResponse = await getUserByID({ serviceConfig, userId });
  } catch (error) {
    logger.error("Error fetching user by ID:", { error });
    return { error: t("verify.errors.couldNotGetUser") };
  }

  if (!userResponse.user) {
    return { error: t("verify.errors.userNotFound") };
  }

  const humanUser = userResponse.user.type.case === "human" ? userResponse.user.type.value : undefined;

  const emailVerificationCheck = await checkEmailVerification(session as any, humanUser, organization, requestId);

  if (emailVerificationCheck?.redirect) {
    return emailVerificationCheck;
  }

  let redirectResult;
  if (requestId && session.id) {
    redirectResult = await completeFlowOrGetUrl(
      {
        sessionId: session.id,
        requestId: requestId,
        organization: organization,
      },
      loginSettings?.defaultRedirectUri,
    );
  } else if (session?.factors?.user?.loginName) {
    redirectResult = await completeFlowOrGetUrl(
      {
        loginName: session.factors.user.loginName,
        organization: organization,
      },
      loginSettings?.defaultRedirectUri,
    );
  }

  // Check if we got a valid redirect result
  if (redirectResult && typeof redirectResult === "object") {
    return redirectResult;
  }

  // Fallback error if we couldn't determine where to redirect
  return { error: t("verify.errors.couldNotDetermineRedirect") };
}
