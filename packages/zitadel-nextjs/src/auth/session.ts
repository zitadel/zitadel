import { createBearerTokenTransport } from "@zitadel/zitadel-js/api/bearer-token";
import {
  createOIDCServiceClient,
  createSessionServiceClient,
  type Checks,
  type RequestChallenges,
  type UserAgent,
} from "@zitadel/zitadel-js/v2";
import { getSession as getOIDCSession } from "../session.js";

export interface SessionLifetime {
  seconds: bigint;
  nanos?: number;
}

/**
 * Shared authentication options for Session API and OIDC callback helpers.
 */
export interface SessionAuthOptions {
  /**
   * ZITADEL API base URL.
   * Falls back to `ZITADEL_API_URL`.
   */
  apiUrl?: string;
  /**
   * Access token used to call ZITADEL APIs.
   * Falls back to `ZITADEL_SERVICE_USER_TOKEN`, then to the current OIDC session token.
   */
  accessToken?: string;
  /**
   * Cookie secret used when reading the OIDC session fallback token.
   * Falls back to `ZITADEL_COOKIE_SECRET`.
   */
  cookieSecret?: string;
}

export interface CreateSessionOptions extends SessionAuthOptions {
  checks?: Checks;
  challenges?: RequestChallenges;
  metadata?: Record<string, Uint8Array>;
  userAgent?: UserAgent;
  lifetime?: SessionLifetime;
}

export interface SetSessionOptions extends SessionAuthOptions {
  sessionId: string;
  sessionToken: string;
  checks?: Checks;
  challenges?: RequestChallenges;
  metadata?: Record<string, Uint8Array>;
  lifetime?: SessionLifetime;
}

export interface GetSessionOptions extends SessionAuthOptions {
  sessionId: string;
  sessionToken?: string;
}

export interface DeleteSessionOptions extends SessionAuthOptions {
  sessionId: string;
  sessionToken?: string;
}

export interface CreateCallbackOptions extends SessionAuthOptions {
  /**
   * OIDC auth request ID.
   * Accepts both plain IDs and values prefixed with `oidc_`.
   */
  authRequestId: string;
  sessionId: string;
  sessionToken: string;
}

function resolveApiUrl(apiUrl?: string): string {
  const resolved = apiUrl ?? process.env.ZITADEL_API_URL;
  if (!resolved) {
    throw new Error(
      "apiUrl option or ZITADEL_API_URL environment variable is required",
    );
  }
  return resolved;
}

async function resolveAccessToken(options?: SessionAuthOptions): Promise<string> {
  if (options?.accessToken) {
    return options.accessToken;
  }

  if (process.env.ZITADEL_SERVICE_USER_TOKEN) {
    return process.env.ZITADEL_SERVICE_USER_TOKEN;
  }

  const session = await getOIDCSession(options?.cookieSecret);
  if (session?.accessToken) {
    return session.accessToken;
  }

  throw new Error(
    "accessToken option, ZITADEL_SERVICE_USER_TOKEN, or an active OIDC session is required",
  );
}

function normalizeAuthRequestId(authRequestId: string): string {
  return authRequestId.startsWith("oidc_")
    ? authRequestId.slice("oidc_".length)
    : authRequestId;
}

async function getClients(options?: SessionAuthOptions) {
  const apiUrl = resolveApiUrl(options?.apiUrl);
  const token = await resolveAccessToken(options);
  const transport = createBearerTokenTransport({
    baseUrl: apiUrl,
    httpVersion: "2",
    token,
  });

  return {
    sessionService: createSessionServiceClient(transport),
    oidcService: createOIDCServiceClient(transport),
  };
}

/**
 * Creates a new ZITADEL session.
 */
export async function createSession(options: CreateSessionOptions) {
  const { sessionService } = await getClients(options);

  return sessionService.createSession({
    checks: options.checks,
    challenges: options.challenges,
    metadata: options.metadata ?? {},
    userAgent: options.userAgent,
    lifetime: options.lifetime,
  });
}

/**
 * Updates an existing ZITADEL session with additional checks/challenges.
 */
export async function setSession(options: SetSessionOptions) {
  const { sessionService } = await getClients(options);

  return sessionService.setSession({
    sessionId: options.sessionId,
    sessionToken: options.sessionToken,
    checks: options.checks,
    challenges: options.challenges,
    metadata: options.metadata ?? {},
    lifetime: options.lifetime,
  });
}

/**
 * Retrieves a ZITADEL session by ID (and optional session token).
 */
export async function getSession(options: GetSessionOptions) {
  const { sessionService } = await getClients(options);

  return sessionService.getSession({
    sessionId: options.sessionId,
    sessionToken: options.sessionToken,
  });
}

/**
 * Deletes a ZITADEL session.
 */
export async function deleteSession(options: DeleteSessionOptions) {
  const { sessionService } = await getClients(options);

  return sessionService.deleteSession({
    sessionId: options.sessionId,
    sessionToken: options.sessionToken,
  });
}

/**
 * Finalizes an OIDC auth request with a verified session and returns the callback URL.
 */
export async function createCallback(options: CreateCallbackOptions) {
  const { oidcService } = await getClients(options);

  const response = await oidcService.createCallback({
    authRequestId: normalizeAuthRequestId(options.authRequestId),
    callbackKind: {
      case: "session",
      value: {
        sessionId: options.sessionId,
        sessionToken: options.sessionToken,
      },
    },
  });

  return response.callbackUrl;
}
