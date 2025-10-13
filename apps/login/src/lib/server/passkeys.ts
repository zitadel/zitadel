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
import { getMostRecentSessionCookie, getSessionCookieById, getSessionCookieByLoginName } from "../cookies";
import { getServiceUrlFromHeaders } from "../service-url";
import { checkEmailVerification, checkUserVerification } from "../verify-helper";
import { createSessionAndUpdateCookie, setSessionAndUpdateCookie } from "./cookie";
import { getOriginalHost } from "./host";
import { completeFlowOrGetUrl } from "../client";

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
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);
  const host = await getOriginalHost();

  let session: GetSessionResponse | undefined;
  let createdSession: Session | undefined;
  let currentUserId: string | undefined = undefined;
  let registerCode: { id: string; code: string } | undefined = undefined;

  if (command.sessionId) {
    // Session-based flow (existing logic)
    const sessionCookie = await getSessionCookieById({ sessionId: command.sessionId });
    session = await getSession({
      serviceUrl,
      sessionId: sessionCookie.id,
      sessionToken: sessionCookie.token,
    });

    if (!session?.session?.factors?.user?.id) {
      return { error: "Could not determine user from session" };
    }

    currentUserId = session.session.factors.user.id;

    const sessionValid = isSessionValid(session.session);

    if (!sessionValid.valid) {
      const authmethods = await listAuthenticationMethodTypes({
        serviceUrl,
        userId: currentUserId,
      });

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

      if (!command.code) {
        // request a new code if no code is provided
        const codeResponse = await createPasskeyRegistrationLink({
          serviceUrl,
          userId: currentUserId,
        });

        if (!codeResponse?.code?.code) {
          return { error: "Could not create registration link" };
        }

        registerCode = codeResponse.code;
      }
    }
  } else if (command.userId && command.code && command.codeId) {
    currentUserId = command.userId;
    registerCode = {
      id: command.codeId,
      code: command.code,
    };

    // Check if user exists
    const userResponse = await getUserByID({
      serviceUrl,
      userId: currentUserId,
    });

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

    createdSession = await createSessionAndUpdateCookie({
      checks,
      requestId: undefined, // No requestId in passkey registration context, TODO: consider if needed
    });

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

  return registerPasskey({
    serviceUrl,
    userId: currentUserId,
    code: registerCode,
    domain: hostname,
  });
}

export async function verifyPasskeyRegistration(command: VerifyPasskeyCommand) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

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

  let currentUserId: string;

  if (command.sessionId) {
    // Session-based flow
    const sessionCookie = await getSessionCookieById({
      sessionId: command.sessionId,
    });
    const session = await getSession({
      serviceUrl,
      sessionId: sessionCookie.id,
      sessionToken: sessionCookie.token,
    });
    const userId = session?.session?.factors?.user?.id;

    if (!userId) {
      throw new Error("Could not get session");
    }

    currentUserId = userId;
  } else {
    // UserId-based flow
    currentUserId = command.userId!;

    // Verify user exists
    const userResponse = await getUserByID({
      serviceUrl,
      userId: currentUserId,
    });

    if (!userResponse || !userResponse.user) {
      throw new Error("User not found");
    }
  }

  return zitadelVerifyPasskeyRegistration({
    serviceUrl,
    request: create(VerifyPasskeyRegistrationRequestSchema, {
      passkeyId: command.passkeyId,
      publicKeyCredential: command.publicKeyCredential,
      passkeyName,
      userId: currentUserId,
    }),
  });
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
  const recentSession = sessionId
    ? await getSessionCookieById({ sessionId })
    : loginName
      ? await getSessionCookieByLoginName({ loginName, organization })
      : await getMostRecentSessionCookie();

  if (!recentSession) {
    return {
      error: "Could not find session",
    };
  }

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization,
  });

  let lifetime = checks?.webAuthN
    ? loginSettings?.multiFactorCheckLifetime // TODO different lifetime for webauthn u2f/passkey
    : checks?.otpEmail || checks?.otpSms
      ? loginSettings?.secondFactorCheckLifetime
      : undefined;

  if (!lifetime) {
    console.warn("No passkey lifetime provided, defaulting to 24 hours");

    lifetime = {
      seconds: BigInt(60 * 60 * 24), // default to 24 hours
      nanos: 0,
    } as Duration;
  }

  const session = await setSessionAndUpdateCookie({
    recentCookie: recentSession,
    checks,
    requestId,
    lifetime,
  });

  if (!session || !session?.factors?.user?.id) {
    return { error: "Could not update session" };
  }

  const userResponse = await getUserByID({
    serviceUrl,
    userId: session?.factors?.user?.id,
  });

  if (!userResponse.user) {
    return { error: "User not found in the system" };
  }

  const humanUser = userResponse.user.type.case === "human" ? userResponse.user.type.value : undefined;

  const emailVerificationCheck = checkEmailVerification(session, humanUser, organization, requestId);

  if (emailVerificationCheck?.redirect) {
    return emailVerificationCheck;
  }

  if (requestId && session.id) {
    return completeFlowOrGetUrl(
      {
        sessionId: session.id,
        requestId: requestId,
        organization: organization,
      },
      loginSettings?.defaultRedirectUri,
    );
  } else if (session?.factors?.user?.loginName) {
    return completeFlowOrGetUrl(
      {
        loginName: session.factors.user.loginName,
        organization: organization,
      },
      loginSettings?.defaultRedirectUri,
    );
  }
}
