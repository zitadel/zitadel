"use server";

import {
  createPasskeyRegistrationLink,
  getLoginSettings,
  getSession,
  getUserByID,
  listAuthenticationMethodTypes,
  registerPasskey,
  verifyPasskeyRegistration as zitadelVerifyPasskeyRegistration,
} from "@/lib/zitadel";
import { create, Duration, Timestamp, timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { Checks, ChecksSchema, GetSessionResponse } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import {
  RegisterPasskeyResponse,
  VerifyPasskeyRegistrationRequestSchema,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { headers } from "next/headers";
import { userAgent } from "next/server";
import { getTranslations } from "next-intl/server";
import { getSessionCookieById } from "../cookies";
import { getServiceConfig } from "../service-url";
import { checkEmailVerification, checkUserVerification } from "../verify-helper";
import { getPublicHost } from "./host";
import { updateOrCreateSession } from "./session";
import { completeFlowOrGetUrl } from "../client";
import { createSessionAndUpdateCookie } from "./cookie";

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

function isSessionValid(session: Partial<Session>): {
  valid: boolean;
  verifiedAt?: Timestamp;
} {
  const validPassword = session?.factors?.password?.verifiedAt;
  const validPasskey = session?.factors?.webAuthN?.verifiedAt;
  const stillValid = session.expirationDate ? timestampDate(session.expirationDate) > new Date() : true;

  const verifiedAt = validPassword || validPasskey;
  const valid = !!((validPassword || validPasskey) && stillValid);

  return { valid, verifiedAt };
}

export async function registerPasskeyLink(
  command: RegisterPasskeyCommand,
): Promise<RegisterPasskeyResponse | { error: string }> {
  if (!command.sessionId && !command.userId) {
    return { error: "Either sessionId or userId must be provided" };
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
      return { error: "Could not get session cookie" };
    }

    session = await getSession({ serviceConfig, sessionId: sessionCookie.id, sessionToken: sessionCookie.token });

    if (!session?.session?.factors?.user?.id) {
      return { error: "Could not determine user from session" };
    }

    currentUserId = session.session.factors.user.id;

    const sessionValid = isSessionValid(session.session);

    if (!sessionValid.valid) {
      const authmethods = await listAuthenticationMethodTypes({ serviceConfig, userId: currentUserId });

      // if the user has no authmethods set, we need to check if the user was verified
      if (authmethods.authMethodTypes.length !== 0) {
        return {
          error: "You have to authenticate or have a valid User Verification Check",
        };
      }

      // check if a verification was done earlier
      const hasValidUserVerificationCheck = await checkUserVerification(currentUserId);

      console.log("hasValidUserVerificationCheck", hasValidUserVerificationCheck);
      if (!hasValidUserVerificationCheck) {
        return { error: "User Verification Check has to be done" };
      }
    }

    // Generate registration code if not provided
    if (command.code && command.codeId) {
      registerCode = {
        id: command.codeId,
        code: command.code,
      };
    } else {
      const codeResponse = await createPasskeyRegistrationLink({ serviceConfig, userId: currentUserId });

      if (!codeResponse?.code?.code) {
        return { error: "Could not create registration link" };
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
    const userResponse = await getUserByID({ serviceConfig, userId: currentUserId });

    if (!userResponse || !userResponse.user) {
      return { error: "User not found" };
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

    const result = await createSessionAndUpdateCookie({
      checks,
      requestId: undefined, // No requestId in passkey registration context, TODO: consider if needed
    });
    createdSession = result.session;

    if (!createdSession) {
      return { error: "Could not create session" };
    }
  }

  if (!registerCode) {
    throw new Error("Missing code in response");
  }

  const [hostname] = host.split(":");

  if (!hostname) {
    throw new Error("Could not get hostname");
  }

  if (!currentUserId) {
    throw new Error("Could not determine user");
  }

  return registerPasskey({ serviceConfig, userId: currentUserId, code: registerCode, domain: hostname });
}

export async function verifyPasskeyRegistration(command: VerifyPasskeyCommand) {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  if (!command.sessionId && !command.userId) {
    throw new Error("Either sessionId or userId must be provided");
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
      throw new Error("Could not get session cookie");
    }

    const session = await getSession({ serviceConfig, sessionId: sessionCookie.id, sessionToken: sessionCookie.token });
    const userId = session?.session?.factors?.user?.id;

    if (!userId) {
      throw new Error("Could not get session");
    }

    currentUserId = userId;
    loginName = session?.session?.factors?.user?.loginName;
  } else {
    // UserId-based flow
    currentUserId = command.userId!;

    // Verify user exists
    const userResponse = await getUserByID({ serviceConfig, userId: currentUserId });

    if (!userResponse || !userResponse.user) {
      throw new Error("User not found");
    }

    loginName = userResponse.user.preferredLoginName;
  }

  const response = await zitadelVerifyPasskeyRegistration({
    serviceConfig,
    request: create(VerifyPasskeyRegistrationRequestSchema, {
      passkeyId: command.passkeyId,
      publicKeyCredential: command.publicKeyCredential,
      passkeyName,
      userId: currentUserId,
    }),
  });

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

  if (result.error) {
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
  const loginSettings = await getLoginSettings({ serviceConfig, organization });

  const userId = session?.factors?.user?.id;
  if (!userId) {
    return { error: t("verify.errors.couldNotFindSession") };
  }

  let userResponse;
  try {
    userResponse = await getUserByID({ serviceConfig, userId });
  } catch (error) {
    console.error("Error fetching user by ID:", error);
    return { error: t("verify.errors.couldNotGetUser") };
  }

  if (!userResponse.user) {
    return { error: t("verify.errors.userNotFound") };
  }

  const humanUser = userResponse.user.type.case === "human" ? userResponse.user.type.value : undefined;

  const emailVerificationCheck = checkEmailVerification(session as any, humanUser, organization, requestId);

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
