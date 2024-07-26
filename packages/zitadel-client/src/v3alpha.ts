import { createClientFor } from "./helpers";
import { UserSchemaService } from "@zitadel/proto/zitadel/user/schema/v3alpha/user_schema_service_connect";
import { UserService } from "@zitadel/proto/zitadel/user/v3alpha/user_service_connect";
import { ActionService } from "@zitadel/proto/zitadel/action/v3alpha/action_service_connect";

export const createUserSchemaServiceClient = createClientFor(UserSchemaService);
export const createUserServiceClient = createClientFor(UserService);
export const createActionServiceClient = createClientFor(ActionService);
