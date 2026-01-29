import { createOpenAPI } from 'fumadocs-openapi/server';
import { createAPIPage } from 'fumadocs-openapi/ui';
import path from 'path';
import client from './api-page.client';

export function APIPage({ document, operations }: { document: string; operations: any[] }) {


  if (!path.isAbsolute(document) && !document.startsWith('openapi/')) {
    throw new Error(`APIPage document path MUST start with 'openapi/'. Received: ${document}`);
  }

  const absoluteDocument = path.isAbsolute(document)
    ? document
    : path.join(process.cwd(), 'openapi', document.slice('openapi/'.length));

  const openapi = createOpenAPI({
    input: [absoluteDocument],
  });

  const InnerAPIPage = createAPIPage(openapi, {
    client,
  });

  return <InnerAPIPage document={absoluteDocument} operations={operations} />;
}

