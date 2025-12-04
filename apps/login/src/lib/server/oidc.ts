"use server";

import { getDeviceAuthorizationRequest as zitadelGetDeviceAuthorizationRequest } from "@/lib/zitadel";
import { headers } from "next/headers";
import { getServiceUrlFromHeaders } from "../service-url";

export async function getDeviceAuthorizationRequest(userCode: string) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  return zitadelGetDeviceAuthorizationRequest({
    serviceUrl,
    userCode,
  });
}
