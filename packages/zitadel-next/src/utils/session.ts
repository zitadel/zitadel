import { createSessionServiceClient } from "@zitadel/client/v2";
import { createServerTransport } from "@zitadel/node";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { getMostRecentCookieWithLoginname } from "./cookies";

const SESSION_LIFETIME_S = 3000;

const transport = createServerTransport(process.env.ZITADEL_SERVICE_USER_TOKEN!, {
  baseUrl: process.env.ZITADEL_API_URL!,
  httpVersion: "2",
});

const sessionService = createSessionServiceClient(transport);

export async function loadMostRecentSession(loginName?: string, organization?: string): Promise<Session | undefined> {
  const recent = await getMostRecentCookieWithLoginname(loginName, organization);
  return getMostRecentSession(recent.id, recent.token);
}

async function getMostRecentSession(sessionId: string, sessionToken: string) {
  return sessionService.getSession({ sessionId, sessionToken }, {}).then((resp) => resp.session);
}
