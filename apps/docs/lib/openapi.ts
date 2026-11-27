import { createOpenAPI } from 'fumadocs-openapi/server';
import { globSync } from 'glob';
import path from 'path';

// Document IDs are the spec paths relative to the docs app, e.g.
// `openapi/latest/zitadel/user/v2/user_service.openapi.json`. The generated API
// reference pages (scripts/generate-api-reference.ts) refer to documents by the same IDs.
export const openapiDocumentIds = globSync('openapi/**/*.openapi.json', {
  cwd: process.cwd(),
  posix: true,
}).sort();

export const openapi = createOpenAPI({
  input: Object.fromEntries(
    openapiDocumentIds.map((id) => [
      id,
      // Keep `'openapi'` as a literal segment: Next.js traces this path into the server
      // functions, and a bare `process.cwd()` join would pull in the whole app directory.
      path.join(process.cwd(), 'openapi', id.slice('openapi/'.length)),
    ]),
  ),
  // Some documents are close to 1 MB and each version has hundreds of operation pages,
  // so pages don't embed them. The browser loads each one once per session from
  // app/api/openapi/[...id] (see `loadDocument` in components/api-page.tsx).
  includeDocument: false,
});
