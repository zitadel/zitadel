import { AdminService } from "../generated/zitadel/admin_pb.js";
import { AuthService } from "../generated/zitadel/auth_pb.js";
import { ManagementService } from "../generated/zitadel/management_pb.js";
import { SystemService } from "../generated/zitadel/system_pb.js";

import { createClientFor } from "../client.js";

export const createAdminServiceClient = createClientFor(AdminService);
export const createAuthServiceClient = createClientFor(AuthService);
export const createManagementServiceClient = createClientFor(ManagementService);
export const createSystemServiceClient = createClientFor(SystemService);
