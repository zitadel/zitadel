import { CompatServiceDefinition } from "nice-grpc/lib/service-definitions";
import { createChannel, createClientFactory } from "nice-grpc";
import {
  AuthServiceClient,
  AuthServiceDefinition,
  GetMyUserResponse,
} from "../proto/server/zitadel/auth";
import { ZitadelServer } from "../server";
import { authMiddleware } from "../middleware";

const createClient = <Client>(
  definition: CompatServiceDefinition,
  accessToken: string
) => {
  const channel = createChannel(process.env.ZITADEL_API_URL ?? "");
  return createClientFactory()
    .use(authMiddleware(accessToken))
    .create(definition, channel) as Client;
};

export async function getAuth(app?: ZitadelServer): Promise<AuthServiceClient> {
  return createClient<AuthServiceClient>(
    AuthServiceDefinition as CompatServiceDefinition,
    ""
  );
}

export async function getMyUser(): Promise<GetMyUserResponse> {
  const auth = await getAuth();
  const response = await auth.getMyUser({});
  return response;
}
