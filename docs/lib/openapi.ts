import { createOpenAPI } from 'fumadocs-openapi/server';

import { globSync } from 'glob';
import { join } from 'path';

const isDev = process.env.NODE_ENV === 'development';
const globPattern = isDev ? 'openapi/latest/**/*.openapi.json' : 'openapi/**/*.openapi.json';

console.time('OpenAPI_Init');
console.log(`[OpenAPI] Initializing... Mode: ${isDev ? 'DEV (Refined)' : 'PROD'}. Pattern: ${globPattern}`);

export const openapi = createOpenAPI({
  input: globSync(globPattern, { cwd: process.cwd(), absolute: true }),
});
console.timeEnd('OpenAPI_Init');
