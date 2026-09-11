import { createSearchAPI } from 'fumadocs-core/search/server';
import { source, versionSource } from '@/lib/source';

// 1. Process latest docs
const latestPages = source.getPages().map((page) => ({
  title: page.data.title,
  description: page.data.description,
  url: page.url,
  id: page.url,
  structuredData: page.data.structuredData,
  tag: 'latest',
}));

// 2. Process versioned docs
const versionedPages = versionSource.getPages().map((page) => {
  // Extract version from URL (e.g., /v4.17/guides/... -> v4.17)
  const match = page.url.match(/^\/(v\d+\.\d+)/);
  
  return {
    title: page.data.title,
    description: page.data.description,
    url: page.url,
    id: page.url,
    structuredData: page.data.structuredData,
    tag: match ? match[1] : 'latest',
  };
});

// 3. Export combined search index
export const { GET } = createSearchAPI('advanced', {
  indexes: [...latestPages, ...versionedPages],
});