// Not a server action module: these helpers read and sign the httpOnly `sessions` cookie, so they
// must only be called by server code. Browser-facing actions live in `lib/server/*`.
import "server-only";

import { LANGUAGE_COOKIE_NAME } from "@/lib/i18n";
import { createLogger } from "@/lib/server/logger";
import { parseAndVerifySessions, signSession } from "@/lib/server/session-cookie-signature";
import { timestampDate, timestampFromMs } from "@zitadel/client";
import { cookies, headers } from "next/headers";

const logger = createLogger("cookies");

const SESSIONS_COOKIE_NAME = "sessions";

// Browsers cap a single cookie at ~4096 bytes and silently drop a Set-Cookie above it, which leaves
// the previous cookie (with an outdated session token) in place. The budget applies to
// `name=value` as it is sent, i.e. after Next.js URL-encodes the value, and keeps a margin for the
// attributes (path, sameSite, secure, httpOnly) that some browsers count towards the limit.
const MAX_COOKIE_SIZE = 3800;

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

function isExpired(session: Cookie, now: Date): boolean {
  return session.expirationTs ? timestampDate(timestampFromMs(Number(session.expirationTs))) <= now : false;
}

function encodedCookieSize(value: string): number {
  // encodeURIComponent only emits ASCII, so the string length equals the byte length
  return `${SESSIONS_COOKIE_NAME}=${encodeURIComponent(value)}`.length;
}

/**
 * Signs the sessions and evicts entries until the cookie fits into MAX_COOKIE_SIZE.
 * Entries the cookie considers expired go first, then the least recently changed ones.
 * The entry with `keepId` (the session of the current flow) is never evicted.
 */
function serializeWithinBudget<T>(sessions: SessionCookie<T>[], keepId?: string): string {
  const signed = sessions.map((session) => signSession(session));
  let kept = signed;
  let value = JSON.stringify(kept);

  if (encodedCookieSize(value) <= MAX_COOKIE_SIZE) {
    return value;
  }

  const sizeBefore = encodedCookieSize(value);
  const now = new Date();
  const evictionOrder = signed
    .filter((session) => session.id !== keepId)
    .sort((a, b) => {
      const expiredA = isExpired(a, now);
      if (expiredA !== isExpired(b, now)) {
        return expiredA ? -1 : 1;
      }
      return Number(a.changeTs || 0) - Number(b.changeTs || 0);
    });

  for (const evicted of evictionOrder) {
    kept = kept.filter((session) => session !== evicted);
    value = JSON.stringify(kept);
    if (encodedCookieSize(value) <= MAX_COOKIE_SIZE) {
      break;
    }
  }

  const sizeAfter = encodedCookieSize(value);
  const meta = { sizeBefore, sizeAfter, evicted: sessions.length - kept.length, remaining: kept.length };
  if (sizeAfter > MAX_COOKIE_SIZE) {
    logger.error("sessions cookie exceeds the size budget even after evicting all other entries", meta);
  } else {
    logger.warn("sessions cookie exceeded the size budget, evicted the oldest entries", meta);
  }

  return value;
}

async function setSessionHttpOnlyCookie<T>(sessions: SessionCookie<T>[], iFrameEnabled: boolean = false, keepId?: string) {
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
    name: SESSIONS_COOKIE_NAME,
    value: serializeWithinBudget(sessions, keepId),
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
    // the size budget is enforced when the cookie is written
    currentSessions = [session].concat(currentSessions);
  }

  if (cleanup) {
    const now = new Date();
    const filteredSessions = currentSessions.filter((s) => !isExpired(s, now));
    return setSessionHttpOnlyCookie(filteredSessions, iFrameEnabled, session.id);
  } else {
    return setSessionHttpOnlyCookie(currentSessions, iFrameEnabled, session.id);
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
      const filteredSessions = sessions.filter((s) => !isExpired(s, now));
      return setSessionHttpOnlyCookie(filteredSessions, iFrameEnabled, session.id);
    } else {
      return setSessionHttpOnlyCookie(sessions, iFrameEnabled, session.id);
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
    const filteredSessions = reducedSessions.filter((s) => !isExpired(s, now));
    return setSessionHttpOnlyCookie(filteredSessions, iFrameEnabled);
  } else {
    return setSessionHttpOnlyCookie(reducedSessions, iFrameEnabled);
  }
}

/**
 * Which entry of the `sessions` cookie is meant. Each kind stands for one of the getters below:
 * `id` for `getSessionCookieById`, `loginName` for `getSessionCookieByLoginName` and `mostRecent`
 * for `getMostRecentCookieWithLoginname` and `getMostRecentSessionCookie`.
 */
export type SessionCookieLookup =
  | { by: "id"; sessionId: string; organization?: string }
  | { by: "loginName"; loginName?: string; organization?: string }
  | { by: "mostRecent"; loginName?: string; organization?: string };

/**
 * The one rule for which entry a lookup means. The getters apply it to the verified entries and
 * `explainSessionCookieMiss` to the entries as the browser sent them, so the two cannot drift apart.
 */
function selectSessionCookie<E extends Partial<Cookie>>(entries: E[], lookup: SessionCookieLookup): E | undefined {
  const candidates = lookup.organization ? entries.filter((entry) => entry.organization === lookup.organization) : entries;

  switch (lookup.by) {
    case "id":
      return candidates.find((entry) => entry.id === lookup.sessionId);
    case "loginName":
      return candidates.find((entry) => entry.loginName === lookup.loginName);
    case "mostRecent": {
      const matching = lookup.loginName ? candidates.filter((entry) => entry.loginName === lookup.loginName) : candidates;
      return matching.length
        ? matching.reduce((prev, current) => ((prev.changeTs ?? "") > (current.changeTs ?? "") ? prev : current))
        : undefined;
    }
  }
}

async function findSessionCookie<T>(lookup: SessionCookieLookup): Promise<SessionCookie<T> | undefined> {
  const cookiesList = await cookies();
  const stringifiedCookie = cookiesList.get(SESSIONS_COOKIE_NAME);

  return selectSessionCookie(readSessionCookies<T>(stringifiedCookie?.value), lookup);
}

export type SessionCookieMiss = {
  /**
   * - `cookie_missing`: the browser sent no `sessions` cookie
   * - `cookie_invalid`: the cookie is unreadable, or the wanted entry has a missing or invalid signature
   * - `session_entry_not_found`: the cookie is fine but holds no entry for the lookup
   */
  reason: "cookie_missing" | "cookie_invalid" | "session_entry_not_found";
  context: Record<string, unknown>;
};

// The entries of the cookie before any signature check, undefined when the value is not a JSON array.
function readUnverifiedEntries(value: string | undefined): Partial<Cookie>[] | undefined {
  if (!value) {
    return [];
  }
  try {
    const parsed = JSON.parse(value);
    return Array.isArray(parsed) ? parsed.filter((entry) => !!entry && typeof entry === "object") : undefined;
  } catch {
    return undefined;
  }
}

/**
 * Explains why a lookup found no entry in the `sessions` cookie. Call it with the lookup of the
 * getter that returned nothing. It reads the cookie again, so it is only meant for that failure
 * path. The context holds counts and sizes only: no token, signature or login name.
 */
export async function explainSessionCookieMiss(lookup: SessionCookieLookup): Promise<SessionCookieMiss> {
  const value = (await cookies()).get(SESSIONS_COOKIE_NAME)?.value;
  const cookieHeader = (await headers()).get("cookie") ?? "";

  const unverified = readUnverifiedEntries(value);
  const entries = unverified ?? [];
  const verified = parseAndVerifySessions<Cookie>(value);

  const context: Record<string, unknown> = {
    lookupBy: lookup.by,
    hasLoginName: lookup.by !== "id" && !!lookup.loginName,
    hasOrganization: !!lookup.organization,
    hasSessionId: lookup.by === "id",
    sessionId: lookup.by === "id" ? lookup.sessionId : undefined,
    cookieEntries: entries.length,
    cookieEntriesDropped: entries.length - verified.length,
    // as the browser sent it (URL-encoded), browsers cap a single cookie at ~4096 bytes
    cookieSize: value ? encodedCookieSize(value) : 0,
    cookieHeaderSize: cookieHeader.length,
    // more than one means cookies with different domain or path attributes shadow each other
    sessionsCookieCount: cookieHeader.match(/(?:^|;\s*)sessions=/g)?.length ?? 0,
  };

  if (lookup.by !== "id" && lookup.loginName) {
    const loginName = lookup.loginName;
    // an entry exists for the login name but not for the organization, or only with different casing
    context.entriesMatchingLoginName = verified.filter((entry) => entry.loginName === loginName).length;
    context.entriesMatchingLoginNameIgnoringCase = verified.filter(
      (entry) => entry.loginName?.toLowerCase() === loginName.toLowerCase(),
    ).length;
  }

  if (!value) {
    return { reason: "cookie_missing", context };
  }

  // the wanted entry is in the cookie, but it was dropped for its missing or invalid signature
  const wantedEntryDropped = !!selectSessionCookie(entries, lookup) && !selectSessionCookie(verified, lookup);
  if (!unverified || wantedEntryDropped) {
    return { reason: "cookie_invalid", context };
  }

  return { reason: "session_entry_not_found", context };
}

export async function getMostRecentSessionCookie<T>(): Promise<Cookie | undefined> {
  return findSessionCookie<T>({ by: "mostRecent" });
}

export async function getSessionCookieById<T>({
  sessionId,
  organization,
}: {
  sessionId: string;
  organization?: string;
}): Promise<SessionCookie<T> | undefined> {
  return findSessionCookie<T>({ by: "id", sessionId, organization });
}

export async function getSessionCookieByLoginName<T>({
  loginName,
  organization,
}: {
  loginName?: string;
  organization?: string;
}): Promise<SessionCookie<T> | undefined> {
  return findSessionCookie<T>({ by: "loginName", loginName, organization });
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
    return sessions.filter((session) => !isExpired(session, now));
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
  return findSessionCookie<T>({ by: "mostRecent", loginName, organization });
}
