import { newSystemToken } from "@zitadel/client/node";

export async function systemAPIToken() {
  const audience = process.env.AUDIENCE;
  const userID = process.env.SYSTEM_USER_ID;
  const key = process.env.SYSTEM_USER_PRIVATE_KEY;

  const decodedToken = Buffer.from(key, "base64").toString("utf-8");

  const token = newSystemToken({
    audience: audience,
    subject: userID,
    key: decodedToken,
  });

  return token;
}
