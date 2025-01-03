import { createClientFor } from "./helpers";

import { AdminService } from "@zitadel/proto/zitadel/admin_pb";
import { AuthService } from "@zitadel/proto/zitadel/auth_pb";
import { ManagementService } from "@zitadel/proto/zitadel/management_pb";
import { SystemService } from "@zitadel/proto/zitadel/system_pb";

export const createAdminServiceClient = createClientFor(AdminService);
export const createAuthServiceClient = createClientFor(AuthService);
export const createManagementServiceClient = createClientFor(ManagementService);
export const createSystemServiceClient = createClientFor(SystemService);
