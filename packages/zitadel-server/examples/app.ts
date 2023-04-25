import {
  ZitadelServerOptions,
  getServer,
  getServers,
  initializeServer,
} from "#";
import { GetMyUserResponse, getAuth } from "#/auth";

async function getMyUser(): Promise<GetMyUserResponse> {
  const auth = await getAuth();
  const response = await auth.getMyUser({});
  return response;
}

async function main() {
  const zitadelConfig: ZitadelServerOptions = {
    apiUrl: "https://dev-mfhquc.zitadel.cloud/",
    token: "123",
  };

  if (!getServers().length) {
    initializeServer(zitadelConfig);
  }

  const app = getServer();
}

main();
