import { createOpenAPI } from 'fumadocs-openapi/server';
import { createAPIPage } from 'fumadocs-openapi/ui';
import client from './api-page.client';

export function APIPage({ document, operations }: { document: string; operations: any[] }) {
  const start = Date.now();

  const openapi = createOpenAPI({
    input: [document],
  });

  const InnerAPIPage = createAPIPage(openapi, {
    client,
  });


  return <InnerAPIPage document={document} operations={operations} />;
}

