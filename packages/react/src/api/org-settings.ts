"use server"

import { create } from "@zitadel/client"
import {
  ListOrganizationDomainsRequestSchema,
  AddOrganizationDomainRequestSchema,
  DeleteOrganizationDomainRequestSchema,
  ListOrganizationMetadataRequestSchema,
  SetOrganizationMetadataRequestSchema,
  MetadataSchema,
  DeleteOrganizationMetadataRequestSchema,
  OrganizationService,
} from "@zitadel/proto/zitadel/org/v2/org_service_pb"
import { createClient } from "@connectrpc/connect"
import { getTransport } from "./transport"

// ─── Domains ─────────────────────────────────────────────────────────────────

export interface OrgDomain {
  domain: string
  isVerified: boolean
  isPrimary: boolean
  validationType: string
}

/**
 * List all domains for an organization.
 */
export async function listOrgDomains(organizationId: string): Promise<OrgDomain[]> {
  const transport = getTransport()
  const client = createClient(OrganizationService, transport)
  const request = create(ListOrganizationDomainsRequestSchema, {
    organizationId,
  })
  const response = await client.listOrganizationDomains(request)
  return (response.domains ?? []).map((d: any) => ({
    domain: d.domain ?? "",
    isVerified: d.isVerified ?? false,
    isPrimary: d.isPrimary ?? false,
    validationType: d.validationType ?? "",
  }))
}

/**
 * Add a domain to an organization.
 */
export async function addOrgDomain(organizationId: string, domain: string): Promise<void> {
  const transport = getTransport()
  const client = createClient(OrganizationService, transport)
  const request = create(AddOrganizationDomainRequestSchema, {
    organizationId,
    domain,
  })
  await client.addOrganizationDomain(request)
}

/**
 * Delete a domain from an organization.
 */
export async function deleteOrgDomain(organizationId: string, domain: string): Promise<void> {
  const transport = getTransport()
  const client = createClient(OrganizationService, transport)
  const request = create(DeleteOrganizationDomainRequestSchema, {
    organizationId,
    domain,
  })
  await client.deleteOrganizationDomain(request)
}

// ─── Metadata ────────────────────────────────────────────────────────────────

export interface OrgMetadataEntry {
  key: string
  value: string // base64-decoded to string
}

/**
 * List all metadata for an organization.
 */
export async function listOrgMetadata(organizationId: string): Promise<OrgMetadataEntry[]> {
  const transport = getTransport()
  const client = createClient(OrganizationService, transport)
  const request = create(ListOrganizationMetadataRequestSchema, {
    organizationId,
  })
  const response = await client.listOrganizationMetadata(request)
  return (response.metadata ?? []).map((m: any) => ({
    key: m.key ?? "",
    value: m.value
      ? typeof m.value === "string"
        ? m.value
        : new TextDecoder().decode(m.value)
      : "",
  }))
}

/**
 * Set a metadata key-value pair on an organization.
 * An empty value will delete the key.
 */
export async function setOrgMetadata(
  organizationId: string,
  key: string,
  value: string
): Promise<void> {
  const transport = getTransport()
  const client = createClient(OrganizationService, transport)
  const entry = create(MetadataSchema, {
    key,
    value: new TextEncoder().encode(value),
  })
  const request = create(SetOrganizationMetadataRequestSchema, {
    organizationId,
    metadata: [entry],
  })
  await client.setOrganizationMetadata(request)
}

/**
 * Delete metadata keys from an organization.
 */
export async function deleteOrgMetadata(
  organizationId: string,
  keys: string[]
): Promise<void> {
  const transport = getTransport()
  const client = createClient(OrganizationService, transport)
  const request = create(DeleteOrganizationMetadataRequestSchema, {
    organizationId,
    keys,
  })
  await client.deleteOrganizationMetadata(request)
}
