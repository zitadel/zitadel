"use server";

import { cookies } from "next/headers";

export type SessionCookie = {
  id: string;
  token: string;
  loginName: string;
  organization?: string;
  creationDate: string;
  expirationDate: string;
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

export async function addSessionToCookie(
  session: SessionCookie,
  cleanup: boolean = false,
): Promise<any> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  let currentSessions: SessionCookie[] = stringifiedCookie?.value
    ? JSON.parse(stringifiedCookie?.value)
    : [];

  const index = currentSessions.findIndex(
    (s) => s.loginName === session.loginName,
  );

  if (index > -1) {
    currentSessions[index] = session;
  } else {
    currentSessions = [...currentSessions, session];
  }

  if (cleanup) {
    const now = new Date();
    const filteredSessions = currentSessions.filter((session) =>
      session.expirationDate ? new Date(session.expirationDate) > now : true,
    );
    return setSessionHttpOnlyCookie(filteredSessions);
  } else {
    return setSessionHttpOnlyCookie(currentSessions);
  }
}

export async function updateSessionCookie(
  id: string,
  session: SessionCookie,
  cleanup: boolean = false,
): Promise<any> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  const sessions: SessionCookie[] = stringifiedCookie?.value
    ? JSON.parse(stringifiedCookie?.value)
    : [session];

  const foundIndex = sessions.findIndex((session) => session.id === id);

  if (foundIndex > -1) {
    sessions[foundIndex] = session;
    if (cleanup) {
      const now = new Date();
      const filteredSessions = sessions.filter((session) =>
        session.expirationDate ? new Date(session.expirationDate) > now : true,
      );
      return setSessionHttpOnlyCookie(filteredSessions);
    } else {
      return setSessionHttpOnlyCookie(sessions);
    }
  } else {
    throw "updateSessionCookie: session id now found";
  }
}

export async function removeSessionFromCookie(
  session: SessionCookie,
  cleanup: boolean = false,
): Promise<any> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  const sessions: SessionCookie[] = stringifiedCookie?.value
    ? JSON.parse(stringifiedCookie?.value)
    : [session];

  const reducedSessions = sessions.filter((s) => s.id !== session.id);
  if (cleanup) {
    const now = new Date();
    const filteredSessions = reducedSessions.filter((session) =>
      session.expirationDate ? new Date(session.expirationDate) > now : true,
    );
    return setSessionHttpOnlyCookie(filteredSessions);
  } else {
    return setSessionHttpOnlyCookie(reducedSessions);
  }
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

export async function getSessionCookieById(
  id: string,
  organization?: string,
): Promise<SessionCookie> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie[] = JSON.parse(stringifiedCookie?.value);

    const found = sessions.find((s) =>
      organization
        ? s.organization === organization && s.id === id
        : s.id === id,
    );
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
  loginName: string,
  organization?: string,
): Promise<SessionCookie> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie[] = JSON.parse(stringifiedCookie?.value);
    const found = sessions.find((s) =>
      organization
        ? s.organization === organization && s.loginName === loginName
        : s.loginName === loginName,
    );
    if (found) {
      return found;
    } else {
      return Promise.reject("no cookie found with loginName: " + loginName);
    }
  } else {
    return Promise.reject("no session cookie found");
  }
}

/**
 *
 * @param cleanup when true, removes all expired sessions, default true
 * @returns Session Cookies
 */
export async function getAllSessionCookieIds(
  cleanup: boolean = false,
): Promise<any> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie[] = JSON.parse(stringifiedCookie?.value);

    if (cleanup) {
      const now = new Date();
      return sessions
        .filter((session) =>
          session.expirationDate
            ? new Date(session.expirationDate) > now
            : true,
        )
        .map((session) => session.id);
    } else {
      return sessions.map((session) => session.id);
    }
  } else {
    return [];
  }
}

/**
 *
 * @param cleanup when true, removes all expired sessions, default true
 * @returns Session Cookies
 */
export async function getAllSessions(
  cleanup: boolean = false,
): Promise<SessionCookie[]> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie[] = JSON.parse(stringifiedCookie?.value);

    if (cleanup) {
      const now = new Date();
      return sessions.filter((session) =>
        session.expirationDate ? new Date(session.expirationDate) > now : true,
      );
    } else {
      return sessions;
    }
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
  loginName?: string,
  organization?: string,
): Promise<any> {
  const cookiesList = cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie[] = JSON.parse(stringifiedCookie?.value);
    let filtered = sessions.filter((cookie) => {
      return !!loginName ? cookie.loginName === loginName : true;
    });

    if (organization) {
      filtered = filtered.filter((cookie) => {
        return cookie.organization === organization;
      });
    }

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
      console.error("sessions", sessions, loginName, organization);
      return Promise.reject("Could not get the context or retrieve a session");
    }
  } else {
    return Promise.reject("Could not read session cookie");
  }
}

export async function clearSessions() {}
