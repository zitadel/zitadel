import { ZitadelApp } from "./app";
import { authMiddleware } from "./middleware";

// const createClient = <Client>(
//   definition: CompatServiceDefinition,
//   accessToken: string
// ) => {
//   const channel = createChannel(process.env.ZITADEL_API_URL ?? "");
//   return createClientFactory()
//     .use(authMiddleware(accessToken))
//     .create(definition, channel) as Client;
// };

export async function getAuth(app?: ZitadelApp) {
  //   return createClient<AuthServiceClient>(
  //     AuthServiceDefinition as CompatServiceDefinition,
  //     ""
  //   );
}
