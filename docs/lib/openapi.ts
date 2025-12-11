import { createOpenAPI } from 'fumadocs-openapi/server';

export const openapi = createOpenAPI({
  input: [
    './openapi/zitadel/user/v2/user_service.openapi.yaml',
    './openapi/zitadel/instance/v2/instance_service.openapi.yaml',
    './openapi/zitadel/application/v2/application_service.openapi.yaml',
    './openapi/zitadel/action/v2/action_service.openapi.yaml',
    './openapi/zitadel/authorization/v2/authorization_service.openapi.yaml',
    './openapi/zitadel/feature/v2/feature_service.openapi.yaml',
    './openapi/zitadel/group/v2/group_service.openapi.yaml',
    './openapi/zitadel/idp/v2/idp_service.openapi.yaml',
    './openapi/zitadel/oidc/v2/oidc_service.openapi.yaml',
    './openapi/zitadel/org/v2/org_service.openapi.yaml',
    './openapi/zitadel/project/v2/project_service.openapi.yaml',
    './openapi/zitadel/saml/v2/saml_service.openapi.yaml',
    './openapi/zitadel/session/v2/session_service.openapi.yaml',
    './openapi/zitadel/settings/v2/settings_service.openapi.yaml',
    './openapi/zitadel/webkey/v2/webkey_service.openapi.yaml',
    './openapi/zitadel/internal_permission/v2/internal_permission_service.openapi.yaml',
  ],
});
