/**
 * Server-side helpers for calling the ZITADEL v2 API from Next.js.
 *
 * Uses the user's access token from the session cookie to authenticate.
 *
 * @example
 * ```ts
 * import { createZitadelApiClient } from "@zitadel/nextjs";
 *
 * const api = await createZitadelApiClient({
 *   apiUrl: "https://my.zitadel.cloud",
 * });
 *
 * const user = await api.userService.getUser({ userId: "123" });
 * ```
 */
import {
  createGrpcTransport,
  createAuthorizationBearerInterceptor,
} from "@zitadel/zitadel-js";
import {
  createUserServiceClient,
  createSettingsServiceClient,
  createSessionServiceClient,
  createOIDCServiceClient,
  createSAMLServiceClient,
  createOrganizationServiceClient,
  createFeatureServiceClient,
  createIdpServiceClient,
  createActionServiceClient,
} from "@zitadel/zitadel-js/v2";
import { getSession } from "./session.js";

export interface ApiClientOptions {
  /**
   * The ZITADEL API base URL, e.g. `https://my.zitadel.cloud`.
   * Falls back to the `ZITADEL_API_URL` environment variable.
   */
  apiUrl?: string;
  /**
   * An explicit access token to use instead of reading from the session.
   * Useful for service-to-service calls with a system token.
   */
  accessToken?: string;
  /**
   * Cookie secret for reading the session (when accessToken is not provided).
   * Falls back to `ZITADEL_COOKIE_SECRET`.
   */
  cookieSecret?: string;
}

export interface ZitadelApiClient {
  /** ZITADEL User Service v2 client. */
  userService: ReturnType<typeof createUserServiceClient>;
  /** ZITADEL Settings Service v2 client. */
  settingsService: ReturnType<typeof createSettingsServiceClient>;
  /** ZITADEL Session Service v2 client. */
  sessionService: ReturnType<typeof createSessionServiceClient>;
  /** ZITADEL OIDC Service v2 client. */
  oidcService: ReturnType<typeof createOIDCServiceClient>;
  /** ZITADEL SAML Service v2 client. */
  samlService: ReturnType<typeof createSAMLServiceClient>;
  /** ZITADEL Organization Service v2 client. */
  organizationService: ReturnType<typeof createOrganizationServiceClient>;
  /** ZITADEL Feature Service v2 client. */
  featureService: ReturnType<typeof createFeatureServiceClient>;
  /** ZITADEL Identity Provider Service v2 client. */
  idpService: ReturnType<typeof createIdpServiceClient>;
  /** ZITADEL Action Service v2 client. */
  actionService: ReturnType<typeof createActionServiceClient>;
}

/**
 * Creates a pre-authenticated ZITADEL v2 API client.
 *
 * If `accessToken` is not provided, reads the session cookie to get
 * the user's access token. Throws if no valid session exists.
 */
export async function createZitadelApiClient(
  options?: ApiClientOptions,
): Promise<ZitadelApiClient> {
  const apiUrl =
    options?.apiUrl ?? process.env.ZITADEL_API_URL;
  if (!apiUrl) {
    throw new Error(
      "apiUrl option or ZITADEL_API_URL environment variable is required",
    );
  }

  // Resolve the access token
  let token = options?.accessToken;
  if (!token) {
    const session = await getSession(options?.cookieSecret);
    if (!session) {
      throw new Error(
        "No valid session found. Call signIn() first or provide an accessToken.",
      );
    }
    token = session.accessToken;
  }

  // Create a gRPC transport with the bearer token interceptor
  const transport = createGrpcTransport({
    baseUrl: apiUrl,
    httpVersion: "2",
    interceptors: [createAuthorizationBearerInterceptor(token)],
  });

  return {
    userService: createUserServiceClient(transport),
    settingsService: createSettingsServiceClient(transport),
    sessionService: createSessionServiceClient(transport),
    oidcService: createOIDCServiceClient(transport),
    samlService: createSAMLServiceClient(transport),
    organizationService: createOrganizationServiceClient(transport),
    featureService: createFeatureServiceClient(transport),
    idpService: createIdpServiceClient(transport),
    actionService: createActionServiceClient(transport),
  };
}

/**
 * Higher-order function that provides an authenticated API client to a handler.
 *
 * Reads the session, creates the API client, and passes it to the handler.
 * Throws if the user is not authenticated.
 *
 * @example
 * ```ts
 * import { withApiClient } from "@zitadel/nextjs";
 *
 * const getUser = withApiClient(async (api, userId: string) => {
 *   return api.userService.getUser({ userId });
 * });
 *
 * // In a server component or server action:
 * const user = await getUser("123");
 * ```
 */
export function withApiClient<TArgs extends unknown[], TResult>(
  handler: (api: ZitadelApiClient, ...args: TArgs) => Promise<TResult>,
  options?: ApiClientOptions,
): (...args: TArgs) => Promise<TResult> {
  return async (...args: TArgs) => {
    const api = await createZitadelApiClient(options);
    return handler(api, ...args);
  };
}
