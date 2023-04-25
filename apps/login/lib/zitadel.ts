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

export { server };
// export async function getMyUser(): Promise<GetMyUserResponse> {
//   const auth = await getAuth();
//   const response = await auth.getMyUser({});
//   return response;
// }
