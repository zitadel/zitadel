import { createClientFor } from "./helpers";

import { AdminService } from "@zitadel/proto/gen/zitadel/admin_connect";
import { AuthService } from "@zitadel/proto/gen/zitadel/auth_connect";
import { ManagementService } from "@zitadel/proto/gen/zitadel/management_connect";
import { SystemService } from "@zitadel/proto/gen/zitadel/system_connect";

export const createAdminServiceClient = createClientFor(AdminService);
export const createAuthServiceClient = createClientFor(AuthService);
export const createManagementServiceClient = createClientFor(ManagementService);
export const createSystemServiceClient = createClientFor(SystemService);
