import { ZITADELUsers } from "@zitadel/proto/zitadel/resources/user/v3alpha/user_service_connect";
import { ZITADELUserSchemas } from "@zitadel/proto/zitadel/resources/userschema/v3alpha/user_schema_service_connect";
import { createClientFor } from "./helpers";

export const createUserSchemaServiceClient = createClientFor(ZITADELUserSchemas);
export const createUserServiceClient = createClientFor(ZITADELUsers);
