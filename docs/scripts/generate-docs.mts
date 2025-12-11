import { generateFiles } from 'fumadocs-openapi';
import { createOpenAPI } from 'fumadocs-openapi/server';
import { writeFileSync } from 'fs';
import { join } from 'path';

const services = [
  'action',
  'application',
  'authorization',
  'feature',
  'group',
  'idp',
  'instance',
  'internal_permission',
  'oidc',
  'org',
  'project',
  'saml',
  'session',
  'settings',
  'user',
  'webkey'
];

const generateServiceDocs = (service: string, filename?: string) => {
  const name = filename || `${service}_service`;
  const api = createOpenAPI({
    input: [`./openapi/zitadel/${service}/v2/${name}.openapi.yaml`],
  });

  void generateFiles({
    input: api,
    output: `./content/docs/references/api/${service}`,
    includeDescription: true,
  });
};

services.forEach(service => generateServiceDocs(service));

const meta = {
  title: "API Reference",
  pages: services
};

writeFileSync(
  join(process.cwd(), 'content/docs/references/api/meta.json'),
  JSON.stringify(meta, null, 2)
);
