import { generateFiles } from 'fumadocs-openapi';
import { createOpenAPI } from 'fumadocs-openapi/server';
import { writeFileSync, mkdirSync } from 'fs';
import { join } from 'path';

// Suppress "Generated: ..." logs to avoid Vercel log limits
const originalLog = console.log;
console.log = (...args) => {
  if (args.length > 0 && typeof args[0] === 'string' && args[0].startsWith('Generated: ')) {
    return;
  }
  originalLog(...args);
};

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

const generateServiceDocs = (service: string, filename?: string, version: string = 'v2') => {
  const name = filename || `${service}_service`;
  const api = createOpenAPI({
    input: [`./openapi/zitadel/${service}/${version}/${name}.openapi.yaml`],
  });

  void generateFiles({
    input: api,
    output: `./content/reference/api/${service}`,
    includeDescription: true,
  });

  const indexPath = join(process.cwd(), `./content/reference/api/${service}/index.mdx`);
  const content = `---
title: ${service.charAt(0).toUpperCase() + service.slice(1)} API
---

API Reference for ${service}
`;
  writeFileSync(indexPath, content);
};

services.forEach(service => generateServiceDocs(service));
generateServiceDocs('org', undefined, 'v2beta');

const generateUserSchemaDocs = () => {
  const api = createOpenAPI({
    input: ['./openapi/zitadel/resources/userschema/v3alpha/user_schema_service.openapi.yaml'],
  });

  void generateFiles({
    input: api,
    output: './content/reference/api/user_schema',
    includeDescription: true,
  });

  const indexPath = join(process.cwd(), './content/reference/api/user_schema/index.mdx');
  const content = `---
title: User Schema API
---

API Reference for User Schema
`;
  writeFileSync(indexPath, content);
};
generateUserSchemaDocs();

const meta = {
  title: "APIs",
  pages: [...services, 'user_schema']
};

mkdirSync(join(process.cwd(), 'content/reference/api'), { recursive: true });

writeFileSync(
  join(process.cwd(), 'content/reference/api/meta.json'),
  JSON.stringify(meta, null, 2)
);

const v1Services = [
  'admin',
  'auth',
  'management',
  'system'
];

const generateV1ServiceDocs = (service: string) => {
  const api = createOpenAPI({
    input: [`./openapi/zitadel/${service}.openapi.yaml`],
  });

  void generateFiles({
    input: api,
    output: `./content/reference/api-v1/${service}`,
    includeDescription: true,
  });

  const indexPath = join(process.cwd(), `./content/reference/api-v1/${service}/index.mdx`);
  const content = `---
title: ${service.charAt(0).toUpperCase() + service.slice(1)} API
---

API Reference for ${service}
`;
  writeFileSync(indexPath, content);
};

v1Services.forEach(service => generateV1ServiceDocs(service));

const v1Meta = {
  title: "API v1",
  pages: v1Services
};

mkdirSync(join(process.cwd(), 'content/reference/api-v1'), { recursive: true });

writeFileSync(
  join(process.cwd(), 'content/reference/api-v1/meta.json'),
  JSON.stringify(v1Meta, null, 2)
);

