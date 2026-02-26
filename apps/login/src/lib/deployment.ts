/**
 * Checks if system user credentials are available for JWT authentication.
 *
 * System user authentication requires:
 * - AUDIENCE: The API audience for JWT authentication
 * - SYSTEM_USER_ID: The system user's ID
 * - SYSTEM_USER_PRIVATE_KEY: The private key for JWT signing
 *
 * Both multi-tenant and self-hosted deployments can use system user authentication.
 *
 * @returns true if system user credentials are present, false otherwise
 */
export function hasSystemUserCredentials(): boolean {
  return !!process.env.AUDIENCE && !!process.env.SYSTEM_USER_ID && !!process.env.SYSTEM_USER_PRIVATE_KEY;
}

/**
 * Checks if login service key file is available for JWT authentication.
 *
 * @returns true if ZITADEL_LOGIN_SERVICE_KEY_FILE and a user ID are present, false otherwise
 */
export function hasLoginServiceKey(): boolean {
  return (
    !!process.env.ZITADEL_LOGIN_SERVICE_KEY_FILE &&
    !!(process.env.ZITADEL_LOGIN_SYSTEM_USER_ID || process.env.SYSTEM_USER_ID)
  );
}

/**
 * Gets the system user ID for login service key authentication.
 *
 * @returns ZITADEL_LOGIN_SYSTEM_USER_ID if set, otherwise SYSTEM_USER_ID
 */
export function getLoginSystemUserId(): string | undefined {
  return process.env.ZITADEL_LOGIN_SYSTEM_USER_ID || process.env.SYSTEM_USER_ID;
}

/**
 * Checks if service account token is available for authentication.
 *
 * @returns true if ZITADEL_SERVICE_USER_TOKEN is present, false otherwise
 */
export function hasServiceUserToken(): boolean {
  return !!process.env.ZITADEL_SERVICE_USER_TOKEN;
}
