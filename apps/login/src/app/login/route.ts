import { isRSCRequest, validateAuthRequest } from "@/lib/auth-utils";
import { getAllSessions } from "@/lib/cookies";
import { createLogger } from "@/lib/logger";
import { FlowInitiationParams, handleOIDCFlowInitiation, handleSAMLFlowInitiation } from "@/lib/server/flow-initiation";
import { getServiceConfig } from "@/lib/service-url";
import { listSessions, ServiceConfig } from "@/lib/zitadel";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";

import { NextRequest, NextResponse } from "next/server";

export const dynamic = "force-dynamic";
export const revalidate = false;
export const fetchCache = "default-no-store";

const logger = createLogger("login-route");

async function loadSessions({ serviceConfig, ids }: { serviceConfig: ServiceConfig; ids: string[] }): Promise<Session[]> {
  const response = await listSessions({ serviceConfig, ids: ids.filter((id: string | undefined) => !!id) });

  return response?.sessions ?? [];
}

export async function GET(request: NextRequest) {
  const { serviceConfig } = getServiceConfig(request.headers);

  const searchParams = request.nextUrl.searchParams;

  // Defensive check: block RSC requests early
  if (isRSCRequest(searchParams)) {
    return NextResponse.json({ error: "RSC requests not supported" }, { status: 400 });
  }

  // Early validation: if no valid request parameters, return error immediately
  const requestId = validateAuthRequest(searchParams);
  if (!requestId) {
    return NextResponse.json({ error: "No valid authentication request found" }, { status: 400 });
  }

  const sessionCookies = await getAllSessions();
  const ids = sessionCookies.map((s) => s.id);
  let sessions: Session[] = [];
  if (ids && ids.length) {
    try {
      sessions = await loadSessions({ serviceConfig, ids });
    } catch (error) {
      logger.warn("Failed to load sessions", { error });
      // listSessions can fail for various reasons (stale/expired session IDs
      // still in cookies, API errors, etc.).  Treat any failure as "no valid
      // sessions" so the user is redirected to loginname instead of a 500.
      sessions = [];
    }
  }

  // Flow initiation - delegate to appropriate handler
  const flowParams: FlowInitiationParams = { serviceConfig, requestId, sessions, sessionCookies, request };

  try {
    if (requestId.startsWith("oidc_")) {
      return await handleOIDCFlowInitiation(flowParams);
    } else if (requestId.startsWith("saml_")) {
      return await handleSAMLFlowInitiation(flowParams);
    } else if (requestId.startsWith("device_")) {
      // Device Authorization does not need to start here as it is handled on the /device endpoint
      return NextResponse.json({ error: "Device authorization should use /device endpoint" }, { status: 400 });
    } else {
      return NextResponse.json({ error: "Invalid request ID format" }, { status: 400 });
    }
  } catch (error: unknown) {
    logger.error("Flow initiation failed", { requestId, error });
    return NextResponse.json({ error: "Internal server error" }, { status: 500 });
  }
}
