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

const searchAPI = createSearchAPI('advanced', {
  // https://docs.orama.com/docs/orama-js/supported-languages
  language: 'english',
  indexes: [
    ...indexPages(source, () => LATEST_VERSION),
    ...indexPages(versionSource, getVersionFromUrl),
  ],
});

// Results only change with a deploy, so the CDN answers repeated queries (cached per
// URL, i.e. per query and tag) instead of a function that may first have to rebuild
// the index on a cold start. A static index is not an option: one version alone
// serializes to over 20 MB of JSON.
export async function GET(request: Request) {
  const response = await searchAPI.GET(request);
  response.headers.set('Cache-Control', 'public, s-maxage=86400, stale-while-revalidate=604800');
  return response;
}
