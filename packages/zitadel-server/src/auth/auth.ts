import { CompatServiceDefinition } from "nice-grpc/lib/service-definitions";
import { createChannel, createClientFactory } from "nice-grpc";
import {
  AuthServiceClient,
  AuthServiceDefinition,
  GetMyUserResponse,
} from "../proto/server/zitadel/auth";
import { ZitadelServer, getServers } from "../server";
import { authMiddleware } from "../middleware";

const createClient = <Client>(
  definition: CompatServiceDefinition,
  apiUrl: string,
  token: string,
) => {
  if (!apiUrl) {
    throw Error("ZITADEL_API_URL not set");
  }

  const channel = createChannel(process.env.ZITADEL_API_URL ?? "");
  return createClientFactory()
    .use(authMiddleware(token))
    .create(definition, channel) as Client;
};

export const getAuth = (app?: string | ZitadelServer) => {
  let config;
  if (app && typeof app === "string") {
    const apps = getServers();
    config = apps.find((a) => a.name === app)?.config;
  } else if (app && typeof app === "object") {
    config = app.config;
  }

  if (!config) {
    throw Error("No ZITADEL app found");
  }

  return createClient<AuthServiceClient>(
    AuthServiceDefinition as CompatServiceDefinition,
    config.apiUrl,
    config.token,
  );
};

export async function getMyUser(): Promise<GetMyUserResponse> {
  const auth = await getAuth();
  const response = await auth.getMyUser({});
  return response;
}
