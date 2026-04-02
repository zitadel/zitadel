"use server";

import { create, toJson } from "@zitadel/client";
import {
  AddOrganizationRequestSchema,
  AddOrganizationResponseSchema,
} from "@zitadel/proto/zitadel/org/v2/org_service_pb";
import { OrganizationService } from "@zitadel/proto/zitadel/org/v2/org_service_pb";
import { createClient } from "@connectrpc/connect";
import { getTransport } from "./transport";

/**
 * Create a new organization via the v2 AddOrganization RPC.
 * Returns the created organization's ID.
 */
export async function createOrganization(name: string): Promise<{ organizationId: string }> {
  const transport = getTransport();
  const client = createClient(OrganizationService, transport);

  const request = create(AddOrganizationRequestSchema, {
    name,
  });

  const response = await client.addOrganization(request);
  const json = toJson(AddOrganizationResponseSchema, response) as any;
  return { organizationId: json.organizationId ?? "" };
}
