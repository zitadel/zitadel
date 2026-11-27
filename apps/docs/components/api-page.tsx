'use client';
import { preload } from 'react-dom';
import { createOpenAPIPage } from 'fumadocs-openapi/ui';

// The Next.js basePath, see next.config.mjs.
const documentUrl = (id: string) => `/docs/api/openapi/${id}`;

export const OpenAPIPage = createOpenAPIPage({
  async loadDocument(id) {
    if (typeof window === 'undefined') {
      // the browser starts downloading the document with the page
      preload(documentUrl(id), { as: 'fetch', crossOrigin: 'anonymous' });
      const { openapi } = await import('@/lib/openapi');
      return (await openapi.getSchema(id)).bundled;
    }

    const response = await fetch(documentUrl(id));
    if (!response.ok) {
      throw new Error(`Failed to load OpenAPI document "${id}": ${response.status}`);
    }
    return response.json();
  },
});
