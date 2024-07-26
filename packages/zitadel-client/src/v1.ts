import { createClientFor } from "./helpers";

import { AdminService } from "@zitadel/proto/zitadel/admin_connect";
import { AuthService } from "@zitadel/proto/zitadel/auth_connect";
import { ManagementService } from "@zitadel/proto/zitadel/management_connect";
import { SystemService } from "@zitadel/proto/zitadel/system_connect";

export const createAdminServiceClient = createClientFor(AdminService);
export const createAuthServiceClient = createClientFor(AuthService);
export const createManagementServiceClient = createClientFor(ManagementService);
export const createSystemServiceClient = createClientFor(SystemService);
