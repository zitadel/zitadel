"use server";

import { loadMostRecentSession } from "./session";
import { verifyTOTPRegistration } from "./zitadel";

export async function verifyTOTP(
  code: string,
  loginName?: string,
  organization?: string,
) {
  return loadMostRecentSession({
    loginName,
    organization,
  }).then((session) => {
    if (session?.factors?.user?.id) {
      return verifyTOTPRegistration(code, session.factors.user.id);
    } else {
      throw Error("No user id found in session.");
    }
  });
}
