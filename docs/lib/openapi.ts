import { createOpenAPI } from 'fumadocs-openapi/server';

import { globSync } from 'glob';
import { join } from 'path';

export const openapi = createOpenAPI({
  input: globSync('openapi/**/*.openapi.json', { cwd: process.cwd(), absolute: true }),
});
