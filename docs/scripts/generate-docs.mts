import { generateFiles } from 'fumadocs-openapi';
import { createOpenAPI } from 'fumadocs-openapi/server';

const generateServiceDocs = (service: string, filename?: string) => {
  const name = filename || `${service}_service`;
  const api = createOpenAPI({
    input: [`./openapi/zitadel/${service}/v2/${name}.openapi.yaml`],
  });

  void generateFiles({
    input: api,
    output: `./content/docs/api/${service}`,
    includeDescription: true,
  });
};

generateServiceDocs('user');
generateServiceDocs('instance');
generateServiceDocs('application');
generateServiceDocs('action');
generateServiceDocs('authorization');
generateServiceDocs('feature');
generateServiceDocs('group');
generateServiceDocs('idp');
generateServiceDocs('oidc');
generateServiceDocs('org');
generateServiceDocs('project');
generateServiceDocs('saml');
generateServiceDocs('session');
generateServiceDocs('settings');
generateServiceDocs('webkey');
generateServiceDocs('internal_permission');
