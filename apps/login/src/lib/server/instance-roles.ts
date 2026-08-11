import { createLogger } from "@/lib/logger";
import { createServiceForHost } from "@/lib/service";
import { getIDPByID, getInstanceId, ServiceConfig } from "@/lib/zitadel";
import { Client } from "@zitadel/client";
import type { InstanceRolesInfo } from "@zitadel/proto/zitadel/idp/v2/idp_pb";
import { InternalPermissionService } from "@zitadel/proto/zitadel/internal_permission/v2/internal_permission_service_pb";

const logger = createLogger("instance-roles");

const ZITADEL_PROJECT_ROLES_CLAIM = "urn:zitadel:iam:org:project:roles";
const IAM_ROLE_PREFIX = "IAM_";

type SyncParams = {
  serviceConfig: ServiceConfig;
  intent: {
    idpInformation?:
      | {
          idpId: string;
          rawInformation?: unknown;
        }
      | undefined;
  };
  userId: string;
};

/**
 * Extracts the instance member roles a user should hold from the ZITADEL
 * project-roles claim ({role: {orgId: orgDomain}}), honoring only roles that
 * were granted in one of the organizations configured in the ZITADEL IdP's
 * instanceRolesInfo (matched on organization id AND domain) and that use an
 * instance role key (IAM_ prefix).
 *
 * Mirrors the filtering of the login v1 flow (external_provider_handler.go).
 */
export function instanceRolesFromClaim(
  rawInformation: unknown,
  rolesInfo: Pick<InstanceRolesInfo, "organizationId" | "organizationDomain">[],
): string[] {
  const raw = rawInformation as Record<string, unknown> | undefined | null;
  const claim = raw?.[ZITADEL_PROJECT_ROLES_CLAIM];
  if (!claim || typeof claim !== "object" || Array.isArray(claim)) {
    return [];
  }

  const roles: string[] = [];
  for (const [role, orgs] of Object.entries(claim as Record<string, unknown>)) {
    if (!role.startsWith(IAM_ROLE_PREFIX)) {
      continue;
    }
    if (!orgs || typeof orgs !== "object" || Array.isArray(orgs)) {
      continue;
    }
    const granted = Object.entries(orgs as Record<string, unknown>).some(
      ([orgId, orgDomain]) =>
        typeof orgDomain === "string" &&
        rolesInfo.some((info) => info.organizationId === orgId && info.organizationDomain === orgDomain),
    );
    if (granted) {
      roles.push(role);
    }
  }
  // The claim is a map, so sort to keep the resulting role set deterministic
  // (mirrors the login v1 flow).
  return roles.sort();
}

/**
 * Synchronizes instance member roles after a login via a ZITADEL identity
 * provider that has instanceRolesInfo configured (e.g. the support-access IdP):
 * roles from the token's project-roles claim are added to the user's instance
 * membership. Merge-only — roles the user already holds are never removed.
 *
 * Never throws: role synchronization must not block the login; a failed sync
 * is retried implicitly on the next login.
 */
export async function syncInstanceRolesFromIdpIntent({ serviceConfig, intent, userId }: SyncParams): Promise<void> {
  try {
    const idpId = intent.idpInformation?.idpId;
    if (!idpId || !userId) {
      return;
    }

    const idp = await getIDPByID({ serviceConfig, id: idpId });
    const config = idp?.config?.config;
    if (config?.case !== "zitadel") {
      return;
    }
    const rolesInfo = config.value.instanceRolesInfo;
    if (!rolesInfo?.length) {
      return;
    }

    // Instance-wide role grants may only originate from an instance-scoped IdP.
    // An organization-scoped provider must never confer instance membership, even
    // if its instanceRolesInfo is somehow populated (e.g. stale data): honoring it
    // would let an org owner escalate themselves to instance-level roles.
    // Mirrors the login v1 guard (external_provider_handler.go) and fails closed.
    const instanceId = await getInstanceId({ serviceConfig });
    if (!instanceId || idp?.details?.resourceOwner !== instanceId) {
      logger.warn("Instance role sync skipped: IdP with instanceRolesInfo is not instance-scoped", { idpId });
      return;
    }

    const grantedRoles = instanceRolesFromClaim(intent.idpInformation?.rawInformation, rolesInfo);
    if (!grantedRoles.length) {
      return;
    }

    const permissionService: Client<typeof InternalPermissionService> = await createServiceForHost(
      InternalPermissionService,
      serviceConfig,
    );
    const instanceResource = { resource: { case: "instance" as const, value: true } };

    // Note: the claim originates from an external instance whose role set may
    // differ. Unknown role keys are rejected by the API on write; the claim
    // filter above already restricts to instance (IAM_) role keys.
    const roles = grantedRoles;

    const { administrators } = await permissionService.listAdministrators(
      {
        filters: [
          { filter: { case: "inUserIdsFilter", value: { ids: [userId] } } },
          { filter: { case: "resource", value: { resource: { case: "instance", value: true } } } },
        ],
      },
      {},
    );
    const existing = administrators?.find((administrator) => administrator.resource?.case === "instance");

    if (!existing) {
      await permissionService.createAdministrator({ userId, resource: instanceResource, roles }, {});
      logger.info("Added instance administrator from ZITADEL IdP roles", { userId, roles });
      return;
    }

    // Merge-only: retain roles granted elsewhere, never remove any.
    const merged = Array.from(new Set([...existing.roles, ...roles]));
    if (merged.length === existing.roles.length) {
      return;
    }

    await permissionService.updateAdministrator({ userId, resource: instanceResource, roles: merged }, {});
    logger.info("Updated instance administrator from ZITADEL IdP roles", { userId, roles: merged });
  } catch (error) {
    logger.warn("Could not synchronize instance roles from IdP intent", { error });
  }
}
