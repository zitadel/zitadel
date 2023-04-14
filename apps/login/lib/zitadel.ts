import { ZitadelOptions, getApps, initializeApp } from "@zitadel/core";

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
const auth = getAuth();
