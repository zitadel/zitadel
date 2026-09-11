import { getBreadcrumbItems } from 'fumadocs-core/breadcrumb';
import { createSearchAPI } from 'fumadocs-core/search/server';
import { source, versionSource } from '@/lib/source';
import { getVersionFromUrl, LATEST_VERSION } from '@/lib/versions';

// Both loaders share one index. Every entry is tagged with its docs version so the
// client can scope results to the version being viewed (see app/providers.tsx).
function indexPages(loader: typeof source | typeof versionSource, getTag: (url: string) => string) {
  const tree = loader.getPageTree();

  return loader.getPages().map((page) => ({
    id: page.url,
    url: page.url,
    title: page.data.title,
    description: page.data.description,
    structuredData: page.data.structuredData,
    // Sidebar path (e.g. "Deploy & Operate › Self-Hosted"), shown with each result.
    breadcrumbs: getBreadcrumbItems(page.url, tree)
      .map((item) => item.name)
      .filter((name): name is string => typeof name === 'string'),
    tag: getTag(page.url),
  }));
}

export const { GET } = createSearchAPI('advanced', {
  // https://docs.orama.com/docs/orama-js/supported-languages
  language: 'english',
  indexes: [
    ...indexPages(source, () => LATEST_VERSION),
    ...indexPages(versionSource, getVersionFromUrl),
  ],
});
