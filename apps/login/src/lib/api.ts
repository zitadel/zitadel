import { newSystemToken } from "@zitadel/client/node";
import { getInstanceDomainByHost } from "./zitadel";

export async function getInstanceUrl(host: string): Promise<string> {
  const [hostname, port] = host.split(":");

  if (hostname === "localhost") {
    console.log("fallback to ZITADEL_API_URL");
    return process.env.ZITADEL_API_URL || "";
  }

  const instanceDomain = await getInstanceDomainByHost(host).catch((error) => {
    console.error(`Could not get instance by host ${host}`, error);
    return null;
  });

  if (!instanceDomain) {
    throw new Error("No instance found");
  }

  console.log(`host: ${host}, api: ${instanceDomain}`);

  return instanceDomain;
}

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
