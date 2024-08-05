import { UserSchemaService } from "@zitadel/proto/gen/zitadel/user/schema/v3alpha/user_schema_service_connect";
import { UserService } from "@zitadel/proto/gen/zitadel/user/v3alpha/user_service_connect";
import { createClientFor } from "./helpers";

export const createUserSchemaServiceClient = createClientFor(UserSchemaService);
export const createUserServiceClient = createClientFor(UserService);
