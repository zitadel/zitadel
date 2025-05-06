"use server";

import {
  authorizeOrDenyDeviceAuthorization,
  getDeviceAuthorizationRequest as zitadelGetDeviceAuthorizationRequest,
} from "@/lib/zitadel";
import { headers } from "next/headers";
import { getServiceUrlFromHeaders } from "../service";

export async function getDeviceAuthorizationRequest(userCode: string) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  return zitadelGetDeviceAuthorizationRequest({
    serviceUrl,
    userCode,
  });
}

export async function denyDeviceAuthorization(deviceAuthorizationId: string) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  // without the session, device auth request is denied
  return authorizeOrDenyDeviceAuthorization({
    serviceUrl,
    deviceAuthorizationId,
  });
}
