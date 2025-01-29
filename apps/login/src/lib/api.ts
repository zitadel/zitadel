import { newSystemToken } from "@zitadel/client/node";

export async function systemAPIToken({
  serviceRegion,
}: {
  serviceRegion: string;
}) {
  const QA = {
    audience: process.env.QA_AUDIENCE,
    userID: process.env.QA_SYSTEM_USER_ID,
    token: Buffer.from(
      process.env.QA_SYSTEM_USER_PRIVATE_KEY,
      "base64",
    ).toString("utf-8"),
  };

  const PROD = {
    audience: process.env.QA_AUDIENCE,
    userID: process.env.QA_SYSTEM_USER_ID,
    token: Buffer.from(
      process.env.PROD_SYSTEM_USER_PRIVATE_KEY,
      "base64",
    ).toString("utf-8"),
  };

  let token;

  switch (serviceRegion) {
    case "eu1":
      token = newSystemToken({
        audience: QA.audience,
        subject: QA.userID,
        key: QA.token,
      });
      break;
    case "us1":
      token = newSystemToken({
        audience: PROD.audience,
        subject: PROD.userID,
        key: PROD.token,
      });
      break;
    default:
      token = newSystemToken({
        audience: QA.audience,
        subject: QA.userID,
        key: QA.token,
      });
  }

  return token;
}
