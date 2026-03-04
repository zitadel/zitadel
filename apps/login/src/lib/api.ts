import { newSystemToken } from "@zitadel/client/node";
import { readFile } from "fs/promises";
import { getLoginSystemUserId } from "./deployment";


// The key token is only loaded once from disk per process.
// If the file was loaded you need to restart the process to switch the key.
let keyToken: Promise<string> | undefined;

async function getTokenFromFile(): Promise<string> {
  keyToken ??= readFile(process.env.SYSTEM_USER_PRIVATE_KEY_FILE, "binary");

  try {
    return await keyToken;
  } catch (error) {
    // if the file doesn't exist, don't cache it
    keyToken = undefined;
    throw error;
  }
}

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
  const key = process.env.SYSTEM_USER_PRIVATE_KEY_FILE
    ? await getTokenFromFile()
    : Buffer.from(process.env.SYSTEM_USER_PRIVATE_KEY, "base64").toString("utf-8");

  return newSystemToken({
    audience: process.env.AUDIENCE,
    subject: process.env.SYSTEM_USER_ID,
    key,
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
    const key = await readFile(keyFile, "utf-8");

    return newSystemToken({
      audience: process.env.AUDIENCE || process.env.ZITADEL_API_URL,
      subject: getLoginSystemUserId()!,
      key: key,
    });
  } catch (err) {
    throw new Error(`Failed to read login service key file "${keyFile}": ${err instanceof Error ? err.message : err}`);
  }
}
