import { createSession, getSession, server, setSession } from "#/lib/zitadel";
import { NextResponse } from "next/server";
import {
  SessionCookie,
  addSessionToCookie,
  updateSessionCookie,
} from "./cookies";
import { ChallengeKind, Session } from "@zitadel/server";

export async function createSessionAndUpdateCookie(
  loginName: string,
  password: string | undefined,
  domain: string,
  challenges: ChallengeKind[] | undefined
): Promise<Session> {
  const createdSession = await createSession(
    server,
    loginName,
    domain,
    password,
    challenges
  );

  if (createdSession) {
    return getSession(
      server,
      createdSession.sessionId,
      createdSession.sessionToken
    ).then((response) => {
      if (response?.session && response.session?.factors?.user?.loginName) {
        const sessionCookie: SessionCookie = {
          id: createdSession.sessionId,
          token: createdSession.sessionToken,
          changeDate: response.session.changeDate?.toString() ?? "",
          loginName: response.session?.factors?.user?.loginName ?? "",
        };

        return addSessionToCookie(sessionCookie).then(() => {
          return response.session as Session;
          //     {
          //     sessionId: createdSession.sessionId,
          //     factors: response?.session?.factors,
          //   });
        });
      } else {
        throw "could not get session or session does not have loginName";
      }
    });
  } else {
    throw "Could not create session";
  }
}

export async function setSessionAndUpdateCookie(
  sessionId: string,
  sessionToken: string,
  loginName: string,
  password: string | undefined,
  domain: string | undefined,
  challenges: ChallengeKind[] | undefined
): Promise<Session> {
  return setSession(
    server,
    sessionId,
    sessionToken,
    domain,
    password,
    challenges
  ).then((session) => {
    if (session) {
      const sessionCookie: SessionCookie = {
        id: sessionId,
        token: session.sessionToken,
        changeDate: session.details?.changeDate?.toString() ?? "",
        loginName: loginName,
      };

      return getSession(server, sessionCookie.id, sessionCookie.token).then(
        (response) => {
          if (response?.session && response.session.factors?.user?.loginName) {
            const { session } = response;
            const newCookie: SessionCookie = {
              id: sessionCookie.id,
              token: sessionCookie.token,
              changeDate: session.changeDate?.toString() ?? "",
              loginName: session.factors?.user?.loginName ?? "",
            };

            return updateSessionCookie(sessionCookie.id, newCookie).then(() => {
              return session;
            });
          } else {
            throw "could not get session or session does not have loginName";
          }
        }
      );
    } else {
      throw "Session not be set";
    }
  });
}
