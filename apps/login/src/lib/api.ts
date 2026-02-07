import { newSystemToken } from "@zitadel/client/node";
import { readFile } from 'fs/promises'

// The key token is only loaded once from disk per process.
// If the file was loaded you need to restart the process to switch the key.
let keyToken: Promise<string> | undefined

async function getTokenFromFile(): Promise<string> {
  keyToken ??= readFile(process.env.SYSTEM_USER_PRIVATE_KEY_FILE, "binary");

  try {
    return await keyToken
  } catch (error) {
    // if the file doesn't exist, don't cache it
    keyToken = undefined
    throw error
  }
}

export async function systemAPIToken() {
  const token = {
    audience: process.env.AUDIENCE,
    userID: process.env.SYSTEM_USER_ID,
    token: process.env.SYSTEM_USER_PRIVATE_KEY_FILE ? await getTokenFromFile() : Buffer.from(process.env.SYSTEM_USER_PRIVATE_KEY, "base64").toString(
      "utf-8",
    ),
  };

  return newSystemToken({
    audience: token.audience,
    subject: token.userID,
    key: token.token,
  });
}
