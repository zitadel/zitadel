import { createClientFor } from "./helpers.js";

import { AdminService } from "@zitadel/proto/zitadel/admin_pb.js";
import { AuthService } from "@zitadel/proto/zitadel/auth_pb.js";
import { ManagementService } from "@zitadel/proto/zitadel/management_pb.js";
import { SystemService } from "@zitadel/proto/zitadel/system_pb.js";

export const createAdminServiceClient: ReturnType<typeof createClientFor<typeof AdminService>> =
  createClientFor(AdminService);
export const createAuthServiceClient: ReturnType<typeof createClientFor<typeof AuthService>> = createClientFor(AuthService);
export const createManagementServiceClient: ReturnType<typeof createClientFor<typeof ManagementService>> =
  createClientFor(ManagementService);
export const createSystemServiceClient: ReturnType<typeof createClientFor<typeof SystemService>> =
  createClientFor(SystemService);
