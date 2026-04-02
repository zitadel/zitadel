"use server";

import { create, toJson } from "@zitadel/client";
import {
  ListUserMetadataRequestSchema,
  ListUserMetadataResponseSchema,
  SetUserMetadataRequestSchema,
  SetUserMetadataResponseSchema,
  DeleteUserMetadataRequestSchema,
  DeleteUserMetadataResponseSchema,
  MetadataSchema,
  UserService,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { createClient } from "@connectrpc/connect";
import { getTransport } from "./transport";

export interface UserMetadataEntry {
  key: string;
  value: string; // decoded from base64
  creationDate?: string;
  changeDate?: string;
}

/**
 * List all metadata for a user.
 */
export async function listUserMetadata(userId: string): Promise<UserMetadataEntry[]> {
  const transport = getTransport();
  const client = createClient(UserService, transport);
  const request = create(ListUserMetadataRequestSchema, { userId });
  const response = await client.listUserMetadata(request);
  const json = toJson(ListUserMetadataResponseSchema, response) as any;
  const entries: UserMetadataEntry[] = (json.metadata ?? []).map((m: any) => ({
    key: m.key ?? "",
    value: m.value ? atob(m.value) : "",
    creationDate: m.creationDate,
    changeDate: m.changeDate,
  }));
  return entries;
}

/**
 * Set (create or update) a metadata entry for a user.
 */
export async function setUserMetadata(userId: string, key: string, value: string) {
  const transport = getTransport();
  const client = createClient(UserService, transport);
  const request = create(SetUserMetadataRequestSchema, {
    userId,
    metadata: [
      create(MetadataSchema, {
        key,
        value: new TextEncoder().encode(value),
      }),
    ],
  });
  await client.setUserMetadata(request);
}

/**
 * Delete metadata entries by key for a user.
 */
export async function deleteUserMetadata(userId: string, keys: string[]) {
  const transport = getTransport();
  const client = createClient(UserService, transport);
  const request = create(DeleteUserMetadataRequestSchema, {
    userId,
    keys,
  });
  await client.deleteUserMetadata(request);
}
