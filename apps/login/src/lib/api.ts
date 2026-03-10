import { newSystemToken } from "@zitadel/client/node";
import { readFileSync } from "fs";
import { getLoginSystemUserId } from "./deployment";

/**
 * Creates a signed JWT token using system user credentials from environment
 * variables. The SYSTEM_USER_PRIVATE_KEY is expected to be base64-encoded and
 * is decoded before signing. Requires AUDIENCE, SYSTEM_USER_ID, and
 * SYSTEM_USER_PRIVATE_KEY to be set in the environment.
 *
 * @returns A signed JWT token string for authenticating API requests.
 * @throws If the underlying token signing fails.
 */
export async function systemAPIToken() {
  const token = {
    audience: process.env.AUDIENCE,
    userID: process.env.SYSTEM_USER_ID,
    token: Buffer.from(process.env.SYSTEM_USER_PRIVATE_KEY, "base64").toString(
      "utf-8",
    ),
  };

  return newSystemToken({
    audience: token.audience,
    subject: token.userID,
    key: token.token,
  });
}

/**
 * Creates a signed JWT token by reading a private key from the file path
 * specified in ZITADEL_LOGIN_SERVICE_KEY_FILE. The audience is resolved from
 * AUDIENCE or falls back to ZITADEL_API_URL, and the subject is resolved via
 * {@link getLoginSystemUserId}.
 *
 * @returns A signed JWT token string for authenticating API requests.
 * @throws If the key file cannot be read or the token signing fails.
 */
export async function loginServiceKeyToken() {
  const keyFile = process.env.ZITADEL_LOGIN_SERVICE_KEY_FILE!;

  try {
    const key = readFileSync(keyFile, "utf-8");

    return newSystemToken({
      audience: process.env.AUDIENCE || process.env.ZITADEL_API_URL,
      subject: getLoginSystemUserId()!,
      key: key,
    });
  } catch (err) {
    throw new Error(`Failed to read login service key file "${keyFile}": ${err instanceof Error ? err.message : err}`);
  }
}
