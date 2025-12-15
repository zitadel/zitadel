"use server";

import { authorizeOrDenyDeviceAuthorization } from "@/lib/zitadel";
import { headers } from "next/headers";
import { getServiceConfig } from "../service-url";

export async function completeDeviceAuthorization(
  deviceAuthorizationId: string,
  session?: { sessionId: string; sessionToken: string },
) {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  // without the session, device auth request is denied
  return authorizeOrDenyDeviceAuthorization({ serviceConfig, deviceAuthorizationId,
    session,
  });
}
