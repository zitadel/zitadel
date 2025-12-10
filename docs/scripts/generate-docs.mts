import { generateFiles } from 'fumadocs-openapi';
import { createOpenAPI } from 'fumadocs-openapi/server';

const userOpenapi = createOpenAPI({
  input: ['./openapi/zitadel/user/v2/user_service.openapi.yaml'],
});

void generateFiles({
  input: userOpenapi,
  output: './content/docs/api/user',
  includeDescription: true,
});

const instanceOpenapi = createOpenAPI({
  input: ['./openapi/zitadel/instance/v2/instance_service.openapi.yaml'],
});

void generateFiles({
  input: instanceOpenapi,
  output: './content/docs/api/instance',
  includeDescription: true,
});

const applicationOpenapi = createOpenAPI({
  input: ['./openapi/zitadel/application/v2/application_service.openapi.yaml'],
});

void generateFiles({
  input: applicationOpenapi,
  output: './content/docs/api/application',
  includeDescription: true,
});
