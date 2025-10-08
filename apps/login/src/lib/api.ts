import { newSystemToken } from "@zitadel/client/node";

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
