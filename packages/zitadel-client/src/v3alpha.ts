import { UserSchemaService } from "@zitadel/proto/zitadel/user/schema/v3alpha/user_schema_service_pb";
import { UserService } from "@zitadel/proto/zitadel/user/v3alpha/user_service_pb";
import { createClientFor } from "./helpers";

export const createUserSchemaServiceClient = createClientFor(UserSchemaService);
export const createUserServiceClient = createClientFor(UserService);
