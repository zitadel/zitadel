"use server";

import { timestampDate, timestampFromMs } from "@zitadel/client";
import { cookies } from "next/headers";
import { LANGUAGE_COOKIE_NAME } from "./i18n";
import { createLogger } from "./logger";
import { parseAndVerifySessions, signSession } from "./session-cookie-signature";

const logger = createLogger("cookies");

// TODO: improve this to handle overflow
// Browsers cap a single cookie at ~4096 bytes (name, value and attributes). Keep a margin for
// the attributes and for the signature that is added to every entry when the cookie is written.
const MAX_COOKIE_SIZE = 3500;

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
    value: JSON.stringify(sessions.map((session) => signSession(session))),
    httpOnly: true,
    path: "/",
    sameSite: resolvedSameSite,
    // "none" is only accepted by browsers together with "secure", otherwise the cookie is silently dropped
    secure: process.env.NODE_ENV === "production" || resolvedSameSite === "none",
  });
}

function readSessionCookies<T>(value: string | undefined): SessionCookie<T>[] {
  const sessions = parseAndVerifySessions<SessionCookie<T>>(value);

  if (value && sessions.length < countCookieEntries(value)) {
    logger.warn(
      `readSessionCookies: ignoring ${countCookieEntries(value) - sessions.length} session cookie entries with a missing or invalid signature (unsigned legacy cookie or signing secret changed)`,
    );
  }

  return sessions;
}

function countCookieEntries(value: string): number {
  try {
    const parsed = JSON.parse(value);
    return Array.isArray(parsed) ? parsed.length : 0;
  } catch {
    return 0;
  }
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

export async function getLanguageCookie(): Promise<string | undefined> {
  const cookiesList = await cookies();
  const languageCookie = cookiesList.get(LANGUAGE_COOKIE_NAME);
  return languageCookie?.value;
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

  let currentSessions: SessionCookie<T>[] = readSessionCookies<T>(stringifiedCookie?.value);

  const index = currentSessions.findIndex((s) => s.loginName === session.loginName);

  if (index > -1) {
    currentSessions[index] = session;
  } else {
    const temp = [...currentSessions, session];

    // measure the value as it will be written (including signatures)
    if (JSON.stringify(temp.map((s) => signSession(s))).length >= MAX_COOKIE_SIZE) {
      logger.warn("WARNING COOKIE OVERFLOW");
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

  const sessions: SessionCookie<T>[] = stringifiedCookie?.value ? readSessionCookies<T>(stringifiedCookie.value) : [session];

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

  const sessions: SessionCookie<T>[] = stringifiedCookie?.value ? readSessionCookies<T>(stringifiedCookie.value) : [session];

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

export async function getMostRecentSessionCookie<T>(): Promise<Cookie | undefined> {
  const cookiesList = await cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  const sessions = readSessionCookies<T>(stringifiedCookie?.value);
  if (!sessions.length) {
    return undefined;
  }

  return sessions.reduce((prev, current) => {
    return prev.changeTs > current.changeTs ? prev : current;
  });
}

export async function getSessionCookieById<T>({
  sessionId,
  organization,
}: {
  sessionId: string;
  organization?: string;
}): Promise<SessionCookie<T> | undefined> {
  const cookiesList = await cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  const sessions = readSessionCookies<T>(stringifiedCookie?.value);
  return sessions.find((s) => (organization ? s.organization === organization && s.id === sessionId : s.id === sessionId));
}

export async function getSessionCookieByLoginName<T>({
  loginName,
  organization,
}: {
  loginName?: string;
  organization?: string;
}): Promise<SessionCookie<T> | undefined> {
  const cookiesList = await cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  const sessions = readSessionCookies<T>(stringifiedCookie?.value);
  return sessions.find((s) =>
    organization ? s.organization === organization && s.loginName === loginName : s.loginName === loginName,
  );
}

/**
 *
 * @param cleanup when true, excludes expired sessions from the result, default false
 * @returns Session Cookies
 */
export async function getAllSessionCookieIds<T>(cleanup: boolean = false): Promise<string[]> {
  const sessions = await getAllSessions<T>(cleanup);
  return sessions.map(({ id }) => id);
}

/**
 *
 * @param cleanup when true, excludes expired sessions from the result, default false
 * @returns Session Cookies
 */
export async function getAllSessions<T>(cleanup: boolean = false): Promise<SessionCookie<T>[]> {
  const cookiesList = await cookies();
  const stringifiedCookie = cookiesList.get("sessions");

  const sessions = readSessionCookies<T>(stringifiedCookie?.value);

  if (!sessions.length) {
    logger.info("getAllSessions: No session cookie found, returning empty array");
    return [];
  }

  if (cleanup) {
    const now = new Date();
    return sessions.filter((session) =>
      session.expirationTs ? timestampDate(timestampFromMs(Number(session.expirationTs))) > now : true,
    );
  }

  return sessions;
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

  const sessions = readSessionCookies<T>(stringifiedCookie?.value);

  if (!sessions.length) {
    return undefined;
  }

  let filtered = sessions;

  if (loginName) {
    filtered = filtered.filter((cookie) => cookie.loginName === loginName);
  }

  if (organization) {
    filtered = filtered.filter((cookie) => cookie.organization === organization);
  }

  if (!filtered || !filtered.length) {
    return undefined;
  }

  return filtered.reduce((prev, current) => {
    return prev.changeTs > current.changeTs ? prev : current;
  });
}
