import { ZitadelOptions } from "@zitadel/server";
import { getAuth } from "@zitadel/server/auth";

import { getApp, getApps, initializeApp } from "@zitadel/server/app";

export const zitadelConfig: ZitadelOptions = {
  apiUrl: process.env.ZITADEL_API_URL ?? "",
  projectId: process.env.ZITADEL_PROJECT_ID ?? "",
  appId: process.env.ZITADEL_APP_ID ?? "",
  token: "this should be a pat",
};

if (!getApps().length) {
  initializeApp(zitadelConfig);
}

const app = getApp();

export async function getMyUser(): Promise<GetMyUserResponse> {
  const auth = await getAuth();
  const response = await auth.getMyUser({});
  return response;
}
