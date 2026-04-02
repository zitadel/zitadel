import { createClientFor, type Client } from "@zitadel/client";
import { UserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { OrganizationService } from "@zitadel/proto/zitadel/org/v2/org_service_pb";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb";
import { SessionService } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { InternalPermissionService } from "@zitadel/proto/zitadel/internal_permission/v2/internal_permission_service_pb";
import { InstanceService } from "@zitadel/proto/zitadel/instance/v2/instance_service_pb";
import { getTransport } from "./transport";

/**
 * Service client factories.
 * Each returns a typed connectRPC client connected to the ZITADEL instance.
 *
 * Usage:
 *   const users = getUserService();
 *   const result = await users.listUsers({ ... });
 */

export function getUserService(): Client<typeof UserService> {
  return createClientFor(UserService)(getTransport());
}

export function getOrganizationService(): Client<typeof OrganizationService> {
  return createClientFor(OrganizationService)(getTransport());
}

export function getSettingsService(): Client<typeof SettingsService> {
  return createClientFor(SettingsService)(getTransport());
}

export function getSessionService(): Client<typeof SessionService> {
  return createClientFor(SessionService)(getTransport());
}

export function getInternalPermissionService(): Client<typeof InternalPermissionService> {
  return createClientFor(InternalPermissionService)(getTransport());
}

export function getInstanceService(): Client<typeof InstanceService> {
  return createClientFor(InstanceService)(getTransport());
}
