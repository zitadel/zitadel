import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { GetSessionResponse } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { getMostRecentCookieWithLoginname } from "./cookies";

export async function loadMostRecentSession(
  sessionService: any, // TODO: SessionServiceClient,
  loginName?: string,
  organization?: string,
): Promise<Session | undefined> {
  const recent = await getMostRecentCookieWithLoginname(loginName, organization);
  return sessionService
    .getSession({ sessionId: recent.id, sessionToken: recent.token }, {})
    .then((resp: GetSessionResponse) => resp.session);
}
