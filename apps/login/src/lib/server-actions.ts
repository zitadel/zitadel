"use server";

import { getMostRecentCookieWithLoginname } from "@/utils/cookies";
import { getSession, verifyTOTPRegistration } from "./zitadel";

export async function verifyTOTP(
  code: string,
  loginName?: string,
  organization?: string,
) {
  return getMostRecentCookieWithLoginname(loginName, organization)
    .then((recent) => {
      return getSession(recent.id, recent.token).then((response) => {
        return { session: response?.session, token: recent.token };
      });
    })
    .then(({ session, token }) => {
      if (session?.factors?.user?.id) {
        return verifyTOTPRegistration(code, session.factors.user.id, token);
      } else {
        throw Error("No user id found in session.");
      }
    });
}
