import {
  management,
  ZitadelServer,
  ZitadelServerOptions,
  getManagement,
  getServer,
  getServers,
  initializeServer,
  LabelPolicy,
} from "@zitadel/server";
// import { getAuth } from "@zitadel/server/auth";

export const zitadelConfig: ZitadelServerOptions = {
  apiUrl: process.env.ZITADEL_API_URL ?? "",
  token: process.env.ZITADEL_SERVICE_USER_TOKEN ?? "",
};

if (!getServers().length) {
  initializeServer(zitadelConfig);
}

const server = getServer();
console.log(server);

export function getBranding(
  server: ZitadelServer
): Promise<LabelPolicy | undefined> {
  const mgmt = getManagement(server);

  return mgmt.getLabelPolicy({}).then((resp) => resp.policy);
}

export { server };
// export async function getMyUser(): Promise<GetMyUserResponse> {
//   const auth = await getAuth();
//   const response = await auth.getMyUser({});
//   return response;
// }
