import { ZITADELUserSchemas } from "@zitadel/proto/zitadel/resources/userschema/v3alpha/user_schema_service_pb.js";
import { createClientFor } from "./helpers.js";

export const createUserSchemaServiceClient: ReturnType<typeof createClientFor<typeof ZITADELUserSchemas>> =
  createClientFor(ZITADELUserSchemas);
