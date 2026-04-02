"use server";

import { create, toJson } from "@zitadel/client";
import {
  ListSessionsRequestSchema,
  ListSessionsResponseSchema,
  DeleteSessionRequestSchema,
  SessionService,
} from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import {
  SearchQuerySchema as SessionSearchQuerySchema,
  UserIDQuerySchema,
  CreatorQuerySchema,
} from "@zitadel/proto/zitadel/session/v2/session_pb";
import { createClient } from "@connectrpc/connect";
import { getTransport } from "./transport";

export interface PaginatedSessions {
  sessions: any[];
  totalResult: number;
}

/**
 * List sessions for a user with pagination.
 * Queries both by UserIDQuery (checked user) and CreatorQuery (session creator)
 * to find all sessions associated with the user, then deduplicates.
 */
export async function listUserSessions(
  userId: string,
  limit: number = 10,
  offset: number = 0,
): Promise<PaginatedSessions> {
  const transport = getTransport();
  const client = createClient(SessionService, transport);

  // Query 1: sessions where user is the checked user
  const userIdQuery = create(SessionSearchQuerySchema, {
    query: {
      case: "userIdQuery",
      value: create(UserIDQuerySchema, { id: userId }),
    },
  });

  // Query 2: sessions created by this user
  const creatorQuery = create(SessionSearchQuerySchema, {
    query: {
      case: "creatorQuery",
      value: create(CreatorQuerySchema, { id: userId }),
    },
  });

  const [byUser, byCreator] = await Promise.all([
    client.listSessions(create(ListSessionsRequestSchema, {
      query: { limit, offset: BigInt(offset) },
      queries: [userIdQuery],
    })),
    client.listSessions(create(ListSessionsRequestSchema, {
      query: { limit, offset: BigInt(offset) },
      queries: [creatorQuery],
    })),
  ]);

  const byUserJson = toJson(ListSessionsResponseSchema, byUser) as any;
  const byCreatorJson = toJson(ListSessionsResponseSchema, byCreator) as any;

  // Merge and deduplicate
  const allSessions = [...(byUserJson.sessions ?? []), ...(byCreatorJson.sessions ?? [])];
  const seen = new Set<string>();
  const deduplicated = allSessions.filter((s: any) => {
    if (seen.has(s.id)) return false;
    seen.add(s.id);
    return true;
  });

  const totalByUser = parseInt(byUserJson.details?.totalResult ?? "0", 10);
  const totalByCreator = parseInt(byCreatorJson.details?.totalResult ?? "0", 10);

  return {
    sessions: deduplicated,
    totalResult: Math.max(totalByUser, totalByCreator, deduplicated.length),
  };
}

/**
 * Revoke (delete) a session by ID.
 */
export async function deleteSession(sessionId: string) {
  const transport = getTransport();
  const client = createClient(SessionService, transport);
  const req = create(DeleteSessionRequestSchema, { sessionId });
  await client.deleteSession(req);
}

