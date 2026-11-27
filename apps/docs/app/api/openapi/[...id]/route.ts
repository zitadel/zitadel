import { openapi, openapiDocumentIds } from '@/lib/openapi';

// Serves the OpenAPI documents that API reference pages load in the browser,
// prerendered at build time like the pages.
export const dynamic = 'force-static';
export const dynamicParams = false;

export async function GET(_request: Request, context: { params: Promise<{ id: string[] }> }) {
  const { id } = await context.params;
  return Response.json((await openapi.getSchema(id.join('/'))).bundled);
}

export function generateStaticParams() {
  return openapiDocumentIds.map((id) => ({ id: id.split('/') }));
}
