"use server";

import { create, toJson } from "@zitadel/client";
import {
  type ListOrganizationsResponse,
  type AddOrganizationResponse,
  ListOrganizationsRequestSchema,
  ListOrganizationsResponseSchema,
  AddOrganizationRequestSchema,
  DeleteOrganizationRequestSchema,
  type AddOrganizationRequest,
} from "@zitadel/proto/zitadel/org/v2/org_service_pb";
import {
  SearchQuerySchema as OrgSearchQuerySchema,
  OrganizationIDQuerySchema,
  type SearchQuery as OrgSearchQuery,
} from "@zitadel/proto/zitadel/org/v2/query_pb";
import { TextQueryMethod } from "@zitadel/proto/zitadel/object/v2/object_pb";
import { getOrganizationService } from "./services";

/**
 * List organizations with optional search queries and pagination.
 */
export async function listOrganizations(opts?: {
  queries?: OrgSearchQuery[];
  pageSize?: number;
  offset?: number;
  sortingColumn?: number;
  asc?: boolean;
}): Promise<ListOrganizationsResponse> {
  const orgService = getOrganizationService();
  const request = create(ListOrganizationsRequestSchema, {
    query: {
      limit: opts?.pageSize ?? 10,
      offset: BigInt(opts?.offset ?? 0),
      asc: opts?.asc ?? true,
    },
    queries: opts?.queries ?? [],
  });
  return orgService.listOrganizations(request);
}

/**
 * Search organizations by name.
 */
export async function searchOrganizationsByName(
  name: string
): Promise<ListOrganizationsResponse> {
  const query = create(OrgSearchQuerySchema, {
    query: {
      case: "nameQuery",
      value: {
        name,
        method: TextQueryMethod.CONTAINS,
      },
    },
  });
  return listOrganizations({ queries: [query] });
}

/**
 * Search organizations by domain.
 */
export async function searchOrganizationsByDomain(
  domain: string
): Promise<ListOrganizationsResponse> {
  const query = create(OrgSearchQuerySchema, {
    query: {
      case: "domainQuery",
      value: {
        domain,
        method: TextQueryMethod.CONTAINS,
      },
    },
  });
  return listOrganizations({ queries: [query] });
}

/**
 * Get the default organization.
 */
export async function getDefaultOrganization(): Promise<ListOrganizationsResponse> {
  const query = create(OrgSearchQuerySchema, {
    query: {
      case: "defaultQuery",
      value: {},
    },
  });
  return listOrganizations({ queries: [query] });
}

/**
 * Create a new organization.
 */
export async function addOrganization(
  request: AddOrganizationRequest
): Promise<AddOrganizationResponse> {
  const orgService = getOrganizationService();
  const req = create(AddOrganizationRequestSchema, request);
  return orgService.addOrganization(req);
}

/**
 * Delete an organization by ID.
 */
export async function deleteOrganization(organizationId: string) {
  const orgService = getOrganizationService();
  const req = create(DeleteOrganizationRequestSchema, { organizationId });
  return orgService.deleteOrganization(req);
}

/**
 * Fetch a single organization by ID as JSON-safe data.
 */
export async function fetchOrganization(orgId: string) {
  const query = create(OrgSearchQuerySchema, {
    query: {
      case: "idQuery",
      value: create(OrganizationIDQuerySchema, { id: orgId }),
    },
  });
  const response = await listOrganizations({ queries: [query], pageSize: 1 });
  const json = toJson(ListOrganizationsResponseSchema, response) as any;
  const orgs = json.result ?? [];
  return orgs.length > 0 ? orgs[0] : null;
}
