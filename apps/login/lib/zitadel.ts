import {
  management,
  ZitadelServer,
  ZitadelServerOptions,
  getManagement,
  orgMetadata,
  getServer,
  getServers,
  LabelPolicy,
  initializeServer,
  PrivacyPolicy,
  PasswordComplexityPolicy,
} from "@zitadel/server";
// import { getAuth } from "@zitadel/server/auth";

export const zitadelConfig: ZitadelServerOptions = {
  name: "zitadel login",
  apiUrl: process.env.ZITADEL_API_URL ?? "",
  token: process.env.ZITADEL_SERVICE_USER_TOKEN ?? "",
};

let server: ZitadelServer;

if (!getServers().length) {
  console.log("initialize server");
  server = initializeServer(zitadelConfig);
}

export function getBranding(
  server: ZitadelServer
): Promise<LabelPolicy | undefined> {
  const mgmt = getManagement(server);
  return mgmt
    .getLabelPolicy(
      {},
      { metadata: orgMetadata(process.env.ZITADEL_ORG_ID ?? "") }
    )
    .then((resp) => resp.policy);
}

export function getPrivacyPolicy(
  server: ZitadelServer
): Promise<PrivacyPolicy | undefined> {
  const mgmt = getManagement(server);
  return mgmt
    .getPrivacyPolicy(
      {},
      { metadata: orgMetadata(process.env.ZITADEL_ORG_ID ?? "") }
    )
    .then((resp) => resp.policy);
}

export function getPasswordComplexityPolicy(
  server: ZitadelServer
): Promise<PasswordComplexityPolicy | undefined> {
  const mgmt = getManagement(server);
  return mgmt
    .getPasswordComplexityPolicy(
      {},
      { metadata: orgMetadata(process.env.ZITADEL_ORG_ID ?? "") }
    )
    .then((resp) => resp.policy);
}

export { server };
