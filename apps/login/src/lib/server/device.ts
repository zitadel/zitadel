"use server";

import { authorizeOrDenyDeviceAuthorization } from "@/lib/zitadel";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import { getServiceConfig } from "../service-url";
import { catchUserError } from "./error-utils";

export async function completeDeviceAuthorization(
  deviceAuthorizationId: string,
  session?: { sessionId: string; sessionToken: string },
): Promise<{ error: string } | void> {
  const t = await getTranslations("device");
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  // without the session, device auth request is denied
  try {
    await authorizeOrDenyDeviceAuthorization({ serviceConfig, deviceAuthorizationId, session });
  } catch (error) {
    // an expired or already-used device code is a routine user state
    return catchUserError(error, t("errors.couldNotComplete"));
  }
}
