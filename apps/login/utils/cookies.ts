"use server";

import { cookies } from "next/headers";

export type SessionCookie = {
  id: string;
  token: string;
  loginName: string;
  changeDate: string;
  authRequestId?: string; // if its linked to an OIDC flow
};

function setSessionHttpOnlyCookie(sessions: SessionCookie[]) {
  const cookiesList = cookies();
  // @ts-ignore
  return cookiesList.set({
    name: "sessions",
    value: JSON.stringify(sessions),
    httpOnly: true,
    path: "/",
  });
}

export async function addSessionToCookie(session: SessionCookie): Promise<any> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  let currentSessions: SessionCookie[] = stringifiedCookie?.value
    ? JSON.parse(stringifiedCookie?.value)
    : [];

  const index = currentSessions.findIndex(
    (s) => s.loginName === session.loginName
  );

  if (index > -1) {
    currentSessions[index] = session;
  } else {
    currentSessions = [...currentSessions, session];
  }

  return setSessionHttpOnlyCookie(currentSessions);
}

export async function updateSessionCookie(
  id: string,
  session: SessionCookie
): Promise<any> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  const sessions: SessionCookie[] = stringifiedCookie?.value
    ? JSON.parse(stringifiedCookie?.value)
    : [session];

  const foundIndex = sessions.findIndex((session) => session.id === id);
  if (foundIndex > -1) {
    sessions[foundIndex] = session;
    return setSessionHttpOnlyCookie(sessions);
  } else {
    throw "updateSessionCookie: session id now found";
  }
}

export async function removeSessionFromCookie(
  session: SessionCookie
): Promise<any> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  const sessions: SessionCookie[] = stringifiedCookie?.value
    ? JSON.parse(stringifiedCookie?.value)
    : [session];

  const filteredSessions = sessions.filter((s) => s.id !== session.id);

  return setSessionHttpOnlyCookie(filteredSessions);
}

export async function getMostRecentSessionCookie(): Promise<any> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie[] = JSON.parse(stringifiedCookie?.value);

    const latest = sessions.reduce((prev, current) => {
      return new Date(prev.changeDate).getTime() >
        new Date(current.changeDate).getTime()
        ? prev
        : current;
    });

    return latest;
  } else {
    return Promise.reject("no session cookie found");
  }
}

export async function getSessionCookieById(id: string): Promise<SessionCookie> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie[] = JSON.parse(stringifiedCookie?.value);

    const found = sessions.find((s) => s.id === id);
    if (found) {
      return found;
    } else {
      return Promise.reject();
    }
  } else {
    return Promise.reject();
  }
}

export async function getSessionCookieByLoginName(
  loginName: string
): Promise<SessionCookie> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie[] = JSON.parse(stringifiedCookie?.value);

    const found = sessions.find((s) => s.loginName === loginName);
    if (found) {
      return found;
    } else {
      return Promise.reject("no cookie found with loginName: " + loginName);
    }
  } else {
    return Promise.reject("no session cookie found");
  }
}

export async function getAllSessionCookieIds(): Promise<any> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie[] = JSON.parse(stringifiedCookie?.value);
    return sessions.map((session) => session.id);
  } else {
    return [];
  }
}

export async function getAllSessions(): Promise<SessionCookie[]> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie[] = JSON.parse(stringifiedCookie?.value);
    return sessions;
  } else {
    return [];
  }
}

/**
 * Returns most recent session filtered by optinal loginName
 * @param loginName
 * @returns most recent session
 */
export async function getMostRecentCookieWithLoginname(
  loginName?: string
): Promise<any> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie[] = JSON.parse(stringifiedCookie?.value);
    const filtered = sessions.filter((cookie) => {
      return !!loginName ? cookie.loginName === loginName : true;
    });

    const latest =
      filtered && filtered.length
        ? filtered.reduce((prev, current) => {
            return new Date(prev.changeDate).getTime() >
              new Date(current.changeDate).getTime()
              ? prev
              : current;
          })
        : undefined;

    if (latest) {
      return latest;
    } else {
      return Promise.reject("Could not get the context or retrieve a session");
    }
  } else {
    return Promise.reject("Could not read session cookie");
  }
}

export async function clearSessions() {}
