import { createOpenAPI } from 'fumadocs-openapi/server';

export const openapi = createOpenAPI({
  input: [
    './openapi/zitadel/user/v2/user_service.openapi.yaml',
    './openapi/zitadel/instance/v2/instance_service.openapi.yaml',
    './openapi/zitadel/application/v2/application_service.openapi.yaml',
  ],
});
