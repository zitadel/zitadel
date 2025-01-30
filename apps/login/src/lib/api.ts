import { newSystemToken } from "@zitadel/client/node";

export async function systemAPIToken({
  serviceRegion,
}: {
  serviceRegion: string;
}) {
  const prefix = serviceRegion.toUpperCase();
  const token = {
    audience: process.env[prefix + "_AUDIENCE"],
    userID: process.env[prefix + "_SYSTEM_USER_ID"],
    token: Buffer.from(
      process.env[prefix.toUpperCase() + "_SYSTEM_USER_PRIVATE_KEY"] as string,
      "base64",
    ).toString("utf-8"),
  };

  if (!token.audience || !token.userID || !token.token) {
    const fallbackToken = {
      audience: process.env.AUDIENCE,
      userID: process.env.SYSTEM_USER_ID,
      token: Buffer.from(
        process.env.SYSTEM_USER_PRIVATE_KEY,
        "base64",
      ).toString("utf-8"),
    };

    return newSystemToken({
      audience: fallbackToken.audience,
      subject: fallbackToken.userID,
      key: fallbackToken.token,
    });
  }

  return newSystemToken({
    audience: token.audience,
    subject: token.userID,
    key: token.token,
  });
}
