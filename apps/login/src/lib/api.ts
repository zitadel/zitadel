import { newSystemToken } from "@zitadel/client/node";
import { readFile } from "fs/promises";

// Keys are only loaded once from disk per process.
// If a file changes you need to restart the process to pick up the new key.
let keyToken: string | undefined;
let loginClientKeyCache: string | undefined;

async function getTokenFromFile(): Promise<string> {
  if (keyToken) {
    return keyToken;
  }

  keyToken = await readFile(process.env.SYSTEM_USER_PRIVATE_KEY_FILE, "binary");
  return keyToken;
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
 * specified in ZITADEL_LOGINCLIENT_KEYFILE. Uses a hardcoded subject of
 * "login-client". The audience is resolved from AUDIENCE or ZITADEL_API_URL.
 *
 * @returns A signed JWT token string for authenticating API requests.
 * @throws If the key file cannot be read or the token signing fails.
 */
export async function loginClientKeyToken() {
  const keyFile = process.env.ZITADEL_LOGINCLIENT_KEYFILE!;

  if (!loginClientKeyCache) {
    try {
      loginClientKeyCache = await readFile(keyFile, "utf-8");
    } catch (err) {
      throw new Error(`Failed to read login client key file "${keyFile}": ${err instanceof Error ? err.message : err}`, {
        cause: err,
      });
    }
  }

  return newSystemToken({
    audience: process.env.AUDIENCE || process.env.ZITADEL_API_URL,
    subject: "login-client",
    key: loginClientKeyCache,
  });
}
