import "server-only";

import { getActiveIdentityProviders, getLoginSettings, ServiceConfig } from "@/lib/server/zitadel";

/**
 * Validates if IDP linking is allowed for a user's organization.
 * Checks:
 * 1. Organization allows external IDP login (allowExternalIdp)
 * 2. The specific IDP is activated for the organization
 *
 */
export async function validateIDPLinkingPermissions({
  serviceConfig,
  userOrganizationId,
  idpId,
}: {
  serviceConfig: ServiceConfig;
  userOrganizationId: string;
  idpId: string;
}): Promise<boolean> {
  // Check organization login settings
  const loginSettings = await getLoginSettings({ serviceConfig, organization: userOrganizationId });

  if (!loginSettings?.allowExternalIdp) {
    return false;
  }

  // Check if the IDP is activated for the organization and allows linking
  const activeIDPs = await getActiveIdentityProviders({ serviceConfig, orgId: userOrganizationId, linking_allowed: true });

  const isIDPActive = activeIDPs.identityProviders?.some((idp) => idp.id === idpId);

  if (!isIDPActive) {
    return false;
  }

  return true;
}
