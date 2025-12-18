"use server";

import { getDeviceAuthorizationRequest as zitadelGetDeviceAuthorizationRequest } from "@/lib/zitadel";
import { headers } from "next/headers";
import { getServiceConfig } from "../service-url";

export async function getDeviceAuthorizationRequest(userCode: string) {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  return zitadelGetDeviceAuthorizationRequest({ serviceConfig, userCode,
  });
}
