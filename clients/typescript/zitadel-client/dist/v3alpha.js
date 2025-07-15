import {
  createClientFor
} from "./chunk-27KHKGT3.js";

// src/v3alpha.ts
import { ZITADELUsers } from "@zitadel/proto/zitadel/resources/user/v3alpha/user_service_pb.js";
import { ZITADELUserSchemas } from "@zitadel/proto/zitadel/resources/userschema/v3alpha/user_schema_service_pb.js";
var createUserSchemaServiceClient = createClientFor(ZITADELUserSchemas);
var createUserServiceClient = createClientFor(ZITADELUsers);
export {
  createUserSchemaServiceClient,
  createUserServiceClient
};
//# sourceMappingURL=v3alpha.js.map