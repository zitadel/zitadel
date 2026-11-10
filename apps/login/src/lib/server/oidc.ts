"use server";

import { getAllSessions } from "@/lib/cookies";
import { createLogger } from "@/lib/logger";
import {
  listSessions,
  ServiceConfig,
  getDeviceAuthorizationRequest as zitadelGetDeviceAuthorizationRequest,
} from "@/lib/zitadel";
import { headers } from "next/headers";
import { setDeviceCodeCheckCookie } from "../device-code-check";
import { getServiceConfig } from "../service-url";

const logger = createLogger("device-code");

// Same rule as flow-initiation.ts for OIDC/SAML: offer the account picker only if
// this browser has at least one live session, otherwise start at /loginname.
async function hasLiveSession(serviceConfig: ServiceConfig): Promise<boolean> {
  const ids = (await getAllSessions()).map((s) => s.id).filter((id) => !!id);

  if (!ids.length) {
    return false;
  }

  try {
    const response = await listSessions({ serviceConfig, ids });
    return (response?.sessions?.length ?? 0) > 0;
  } catch (error) {
    logger.warn("Failed to load sessions for device flow, falling back to loginname", { error });
    return false;
  }
}

export async function getDeviceAuthorizationRequest(userCode: string) {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const { deviceAuthorizationRequest } = await zitadelGetDeviceAuthorizationRequest({ serviceConfig, userCode });

  if (!deviceAuthorizationRequest) {
    return { deviceAuthorizationRequest: undefined, nextUrl: undefined };
  }

  await setDeviceCodeCheckCookie(userCode);

  // authenticate first (picker or loginname), consent is asked for afterwards on /device/consent
  const path = (await hasLiveSession(serviceConfig)) ? "/accounts" : "/loginname";
  const nextUrl = `${path}?` + new URLSearchParams({ requestId: `device_${userCode}` });

  return { deviceAuthorizationRequest, nextUrl };
}
