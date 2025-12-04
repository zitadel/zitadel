"use server";

import { timestampDate, timestampFromMs } from "@zitadel/client";
import { cookies } from "next/headers";
import { LANGUAGE_COOKIE_NAME } from "./i18n";

// TODO: improve this to handle overflow
const MAX_COOKIE_SIZE = 2048;

export type Cookie = {
  id: string;
  token: string;
  loginName: string;
  organization?: string;
  creationTs: string;
  expirationTs: string;
  changeTs: string;
  requestId?: string; // if its linked to an OIDC flow
};

type SessionCookie<T> = Cookie & T;

async function setSessionHttpOnlyCookie<T>(sessions: SessionCookie<T>[], iFrameEnabled: boolean = false) {
  const cookiesList = await cookies();

  // "none" is required for iframe embedding (with secure flag)
  let resolvedSameSite: "lax" | "strict" | "none";

  if (iFrameEnabled) {
    // When embedded in iframe, must use "none" with secure flag
    resolvedSameSite = "none";
  } else {
    // This allows cookies during top-level navigation while blocking cross-origin requests
    resolvedSameSite = "lax";
  }

  return cookiesList.set({
    name: "sessions",
    value: JSON.stringify(sessions),
    httpOnly: true,
    path: "/",
    sameSite: resolvedSameSite,
    secure: process.env.NODE_ENV === "production",
  });
}

export async function setLanguageCookie(language: string) {
  const cookiesList = await cookies();

  await cookiesList.set({
    name: LANGUAGE_COOKIE_NAME,
    value: language,
    httpOnly: true,
    path: "/",
  });
}

export async function addSessionToCookie<T>({
  session,
  cleanup,
  iFrameEnabled,
}: {
  session: SessionCookie<T>;
  cleanup?: boolean;
  iFrameEnabled?: boolean;
}): Promise<any> {
  const cookiesList = await cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  let currentSessions: SessionCookie<T>[] = stringifiedCookie?.value ? JSON.parse(stringifiedCookie?.value) : [];

  const index = currentSessions.findIndex((s) => s.loginName === session.loginName);

  if (index > -1) {
    currentSessions[index] = session;
  } else {
    const temp = [...currentSessions, session];

    if (JSON.stringify(temp).length >= MAX_COOKIE_SIZE) {
      console.log("WARNING COOKIE OVERFLOW");
      // TODO: improve cookie handling
      // this replaces the first session (oldest) with the new one
      currentSessions = [session].concat(currentSessions.slice(1));
    } else {
      currentSessions = [session].concat(currentSessions);
    }
  }

  if (cleanup) {
    const now = new Date();
    const filteredSessions = currentSessions.filter((session) =>
      session.expirationTs ? timestampDate(timestampFromMs(Number(session.expirationTs))) > now : true,
    );
    return setSessionHttpOnlyCookie(filteredSessions, iFrameEnabled);
  } else {
    return setSessionHttpOnlyCookie(currentSessions, iFrameEnabled);
  }
}

export async function updateSessionCookie<T>({
  id,
  session,
  cleanup,
  iFrameEnabled,
}: {
  id: string;
  session: SessionCookie<T>;
  cleanup?: boolean;
  iFrameEnabled?: boolean;
}): Promise<any> {
  const cookiesList = await cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  const sessions: SessionCookie<T>[] = stringifiedCookie?.value ? JSON.parse(stringifiedCookie?.value) : [session];

  const foundIndex = sessions.findIndex((session) => session.id === id);

  if (foundIndex > -1) {
    sessions[foundIndex] = session;
    if (cleanup) {
      const now = new Date();
      const filteredSessions = sessions.filter((session) =>
        session.expirationTs ? timestampDate(timestampFromMs(Number(session.expirationTs))) > now : true,
      );
      return setSessionHttpOnlyCookie(filteredSessions, iFrameEnabled);
    } else {
      return setSessionHttpOnlyCookie(sessions, iFrameEnabled);
    }
  } else {
    throw "updateSessionCookie<T>: session id not found";
  }
}

export async function removeSessionFromCookie<T>({
  session,
  cleanup,
  iFrameEnabled,
}: {
  session: SessionCookie<T>;
  cleanup?: boolean;
  iFrameEnabled?: boolean;
}) {
  const cookiesList = await cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  const sessions: SessionCookie<T>[] = stringifiedCookie?.value ? JSON.parse(stringifiedCookie?.value) : [session];

  const reducedSessions = sessions.filter((s) => s.id !== session.id);
  if (cleanup) {
    const now = new Date();
    const filteredSessions = reducedSessions.filter((session) =>
      session.expirationTs ? timestampDate(timestampFromMs(Number(session.expirationTs))) > now : true,
    );
    return setSessionHttpOnlyCookie(filteredSessions, iFrameEnabled);
  } else {
    return setSessionHttpOnlyCookie(reducedSessions, iFrameEnabled);
  }
}

export async function getMostRecentSessionCookie<T>(): Promise<Cookie> {
  const cookiesList = await cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie<T>[] = JSON.parse(stringifiedCookie?.value);

    const latest = sessions.reduce((prev, current) => {
      return prev.changeTs > current.changeTs ? prev : current;
    });

    return latest;
  } else {
    return Promise.reject("no session cookie found");
  }
}

export async function getSessionCookieById<T>({
  sessionId,
  organization,
}: {
  sessionId: string;
  organization?: string;
}): Promise<SessionCookie<T>> {
  const cookiesList = await cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie<T>[] = JSON.parse(stringifiedCookie?.value);

    const found = sessions.find((s) =>
      organization ? s.organization === organization && s.id === sessionId : s.id === sessionId,
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

export async function getSessionCookieByLoginName<T>({
  loginName,
  organization,
}: {
  loginName?: string;
  organization?: string;
}): Promise<SessionCookie<T>> {
  const cookiesList = await cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie<T>[] = JSON.parse(stringifiedCookie?.value);
    const found = sessions.find((s) =>
      organization ? s.organization === organization && s.loginName === loginName : s.loginName === loginName,
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
export async function getAllSessionCookieIds<T>(cleanup: boolean = false): Promise<string[]> {
  const cookiesList = await cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie<T>[] = JSON.parse(stringifiedCookie?.value);

    if (cleanup) {
      const now = new Date();
      return sessions
        .filter((session) =>
          session.expirationTs ? timestampDate(timestampFromMs(Number(session.expirationTs))) > now : true,
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
export async function getAllSessions<T>(cleanup: boolean = false): Promise<SessionCookie<T>[]> {
  const cookiesList = await cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie<T>[] = JSON.parse(stringifiedCookie?.value);

    if (cleanup) {
      const now = new Date();
      return sessions.filter((session) =>
        session.expirationTs ? timestampDate(timestampFromMs(Number(session.expirationTs))) > now : true,
      );
    } else {
      return sessions;
    }
  } else {
    console.log("getAllSessions: No session cookie found, returning empty array");
    return [];
  }
}

/**
 * Returns most recent session filtered by optinal loginName
 * @param loginName optional loginName to filter cookies, if non provided, returns most recent session
 * @param organization optional organization to filter cookies
 * @returns most recent session
 */
export async function getMostRecentCookieWithLoginname<T>({
  loginName,
  organization,
}: {
  loginName?: string;
  organization?: string;
}): Promise<any> {
  const cookiesList = await cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  if (stringifiedCookie?.value) {
    const sessions: SessionCookie<T>[] = JSON.parse(stringifiedCookie?.value);
    let filtered = sessions.filter((cookie) => {
      return loginName ? cookie.loginName === loginName : true;
    });

    if (organization) {
      filtered = filtered.filter((cookie) => {
        return cookie.organization === organization;
      });
    }

    const latest =
      filtered && filtered.length
        ? filtered.reduce((prev, current) => {
            return prev.changeTs > current.changeTs ? prev : current;
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
