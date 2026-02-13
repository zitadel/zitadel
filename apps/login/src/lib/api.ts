import { newSystemToken } from "@zitadel/client/node";
import { readFileSync } from "fs";
import { getLoginSystemUserId } from "./deployment";

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

export async function loginServiceKeyToken() {
  const key = readFileSync(process.env.LOGIN_SERVICE_KEY_FILE!, "utf-8");
  const audience = process.env.AUDIENCE || process.env.ZITADEL_API_URL!;

  return newSystemToken({
    audience: audience,
    subject: getLoginSystemUserId()!,
    key: key,
  });
}
