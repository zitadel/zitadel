import { createClientFor } from "./helpers.js";

import { AdminService } from "@zitadel/proto/zitadel/admin_pb.js";
import { AuthService } from "@zitadel/proto/zitadel/auth_pb.js";
import { ManagementService } from "@zitadel/proto/zitadel/management_pb.js";
import { SystemService } from "@zitadel/proto/zitadel/system_pb.js";

export const createAdminServiceClient = createClientFor(AdminService);
export const createAuthServiceClient = createClientFor(AuthService);
export const createManagementServiceClient = createClientFor(ManagementService);
export const createSystemServiceClient = createClientFor(SystemService);
