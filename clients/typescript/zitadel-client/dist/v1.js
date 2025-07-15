import {
  createClientFor
} from "./chunk-27KHKGT3.js";

// src/v1.ts
import { AdminService } from "@zitadel/proto/zitadel/admin_pb.js";
import { AuthService } from "@zitadel/proto/zitadel/auth_pb.js";
import { ManagementService } from "@zitadel/proto/zitadel/management_pb.js";
import { SystemService } from "@zitadel/proto/zitadel/system_pb.js";
var createAdminServiceClient = createClientFor(AdminService);
var createAuthServiceClient = createClientFor(AuthService);
var createManagementServiceClient = createClientFor(ManagementService);
var createSystemServiceClient = createClientFor(SystemService);
export {
  createAdminServiceClient,
  createAuthServiceClient,
  createManagementServiceClient,
  createSystemServiceClient
};
//# sourceMappingURL=v1.js.map