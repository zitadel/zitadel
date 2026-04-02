"use server";

import { create, toJson } from "@zitadel/client";
import {
  ListSessionsRequestSchema,
  ListSessionsResponseSchema,
  SessionService,
} from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { createClient } from "@connectrpc/connect";
import { getTransport } from "./transport";

/**
 * List all sessions with pagination (no user filter).
 */
export async function fetchAllSessions(limit: number = 50, offset: number = 0) {
  const transport = getTransport();
  const client = createClient(SessionService, transport);

  const request = create(ListSessionsRequestSchema, {
    query: {
      limit,
      offset: BigInt(offset),
    },
  });

  const response = await client.listSessions(request);
  const json = toJson(ListSessionsResponseSchema, response) as any;

  return {
    sessions: json.sessions ?? [],
    totalResult: parseInt(json.details?.totalResult ?? "0", 10),
  };
}
