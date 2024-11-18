"use server";

import {
  createSessionForIdpAndUpdateCookie,
  setSessionAndUpdateCookie,
} from "@/lib/server/cookie";
import { deleteSession, listAuthenticationMethodTypes } from "@/lib/zitadel";
import { RequestChallenges } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { Checks } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import {
  getMostRecentSessionCookie,
  getSessionCookieById,
  getSessionCookieByLoginName,
  removeSessionFromCookie,
} from "../cookies";

type CreateNewSessionCommand = {
  userId: string;
  idpIntent: {
    idpIntentId: string;
    idpIntentToken: string;
  };
  loginName?: string;
  password?: string;
  authRequestId?: string;
};

export async function createNewSessionForIdp(options: CreateNewSessionCommand) {
  const { userId, idpIntent, authRequestId } = options;

  if (!userId || !idpIntent) {
    throw new Error("No userId or loginName provided");
  }
  return createSessionForIdpAndUpdateCookie(userId, idpIntent, authRequestId);
}

export type UpdateSessionCommand = {
  loginName?: string;
  sessionId?: string;
  organization?: string;
  checks?: Checks;
  authRequestId?: string;
  challenges?: RequestChallenges;
};

export async function updateSession(options: UpdateSessionCommand) {
  let {
    loginName,
    sessionId,
    organization,
    checks,
    authRequestId,
    challenges,
  } = options;
  const sessionPromise = sessionId
    ? getSessionCookieById({ sessionId }).catch((error) => {
        return Promise.reject(error);
      })
    : loginName
      ? getSessionCookieByLoginName({ loginName, organization }).catch(
          (error) => {
            return Promise.reject(error);
          },
        )
      : getMostRecentSessionCookie().catch((error) => {
          return Promise.reject(error);
        });

  // TODO remove ports from host header for URL with port
  const host = "localhost";

  if (
    host &&
    challenges &&
    challenges.webAuthN &&
    !challenges.webAuthN.domain
  ) {
    const [hostname, port] = host.split(":");
    challenges.webAuthN.domain = hostname;
  }

  const recent = await sessionPromise;

  const session = await setSessionAndUpdateCookie(
    recent,
    checks,
    challenges,
    authRequestId,
  );

  // if password, check if user has MFA methods
  let authMethods;
  if (checks && checks.password && session.factors?.user?.id) {
    const response = await listAuthenticationMethodTypes(
      session.factors.user.id,
    );
    if (response.authMethodTypes && response.authMethodTypes.length) {
      authMethods = response.authMethodTypes;
    }
  }

  return {
    sessionId: session.id,
    factors: session.factors,
    challenges: session.challenges,
    authMethods,
  };
}

type ClearSessionOptions = {
  sessionId: string;
};

export async function clearSession(options: ClearSessionOptions) {
  const { sessionId } = options;

  const session = await getSessionCookieById({ sessionId });

  const deletedSession = await deleteSession(session.id, session.token);

  if (deletedSession) {
    return removeSessionFromCookie(session);
  }
}

type CleanupSessionCommand = {
  sessionId: string;
};

export async function cleanupSession({ sessionId }: CleanupSessionCommand) {
  const sessionCookie = await getSessionCookieById({ sessionId });

  const deleteResponse = await deleteSession(
    sessionCookie.id,
    sessionCookie.token,
  );

  if (!deleteResponse) {
    throw new Error("Could not delete session");
  }

  return removeSessionFromCookie(sessionCookie);
}
