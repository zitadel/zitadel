import {
  ZitadelServerOptions,
  getServer,
  getServers,
  initializeServer,
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

// export async function getMyUser(): Promise<GetMyUserResponse> {
//   const auth = await getAuth();
//   const response = await auth.getMyUser({});
//   return response;
// }
