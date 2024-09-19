import { timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { GetSessionResponse } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { getMostRecentCookieWithLoginname } from "./cookies";
import { sessionService } from "./zitadel";

export async function loadMostRecentSession(sessionParams: {
  loginName?: string;
  organization?: string;
}): Promise<Session | undefined> {
  const recent = await getMostRecentCookieWithLoginname({
    loginName: sessionParams.loginName,
    organization: sessionParams.organization,
  });
  return sessionService
    .getSession({ sessionId: recent.id, sessionToken: recent.token }, {})
    .then((resp: GetSessionResponse) => resp.session);
}

export function isSessionValid(session: Partial<Session>) {
  const validPassword = session?.factors?.password?.verifiedAt;
  const validPasskey = session?.factors?.webAuthN?.verifiedAt;
  const stillValid = session.expirationDate
    ? timestampDate(session.expirationDate) > new Date()
    : true;

  const verifiedAt = validPassword || validPasskey;
  const valid = (validPassword || validPasskey) && stillValid;

  return { valid, verifiedAt };
}
