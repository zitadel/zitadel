"use server"

import { create, toJson } from "@zitadel/client"
import {
  UpdateOrganizationRequestSchema,
  UpdateOrganizationResponseSchema,
  DeleteOrganizationRequestSchema,
  DeleteOrganizationResponseSchema,
  DeactivateOrganizationRequestSchema,
  DeactivateOrganizationResponseSchema,
  ActivateOrganizationRequestSchema,
  ActivateOrganizationResponseSchema,
} from "@zitadel/proto/zitadel/org/v2/org_service_pb"
import { OrganizationService } from "@zitadel/proto/zitadel/org/v2/org_service_pb"
import { createClient } from "@connectrpc/connect"
import { getTransport } from "./transport"

/**
 * Update an organization's name via the v2 UpdateOrganization RPC.
 */
export async function updateOrganization(
  organizationId: string,
  name: string
): Promise<void> {
  const transport = getTransport()
  const client = createClient(OrganizationService, transport)
  const request = create(UpdateOrganizationRequestSchema, {
    organizationId,
    name,
  })
  await client.updateOrganization(request)
}

/**
 * Delete an organization via the v2 DeleteOrganization RPC.
 * This permanently removes the organization and all its resources.
 */
export async function deleteOrganization(
  organizationId: string
): Promise<void> {
  const transport = getTransport()
  const client = createClient(OrganizationService, transport)
  const request = create(DeleteOrganizationRequestSchema, {
    organizationId,
  })
  await client.deleteOrganization(request)
}

/**
 * Deactivate an organization via the v2 DeactivateOrganization RPC.
 * Users of this organization will not be able to log in.
 */
export async function deactivateOrganization(
  organizationId: string
): Promise<void> {
  const transport = getTransport()
  const client = createClient(OrganizationService, transport)
  const request = create(DeactivateOrganizationRequestSchema, {
    organizationId,
  })
  await client.deactivateOrganization(request)
}

/**
 * Activate an organization via the v2 ActivateOrganization RPC.
 * Only works if the organization is currently deactivated.
 */
export async function activateOrganization(
  organizationId: string
): Promise<void> {
  const transport = getTransport()
  const client = createClient(OrganizationService, transport)
  const request = create(ActivateOrganizationRequestSchema, {
    organizationId,
  })
  await client.activateOrganization(request)
}
