import { getBreadcrumbItems } from 'fumadocs-core/breadcrumb';
import { createSearchAPI } from 'fumadocs-core/search/server';
import { source, versionSource } from '@/lib/source';
import { getVersionFromUrl, LATEST_VERSION } from '@/lib/versions';

function indexPages(loader: typeof source | typeof versionSource, getTag: (url: string) => string) {
  const tree = loader.getPageTree();

  return loader.getPages().map((page) => ({
    id: page.url,
    url: page.url,
    title: page.data.title,
    description: page.data.description,
    structuredData: page.data.structuredData,
    breadcrumbs: getBreadcrumbItems(page.url, tree)
      .map((item) => item.name)
      .filter((name): name is string => typeof name === 'string'),
    tag: getTag(page.url),
  }));
}

// 1. Force Next.js to cache this index statically at build time
export const revalidate = false; 

// 2. Export `staticGET` as `GET`
export const { staticGET: GET } = createSearchAPI('advanced', {
  language: 'english',
  indexes: [
    ...indexPages(source, () => LATEST_VERSION),
    ...indexPages(versionSource, getVersionFromUrl),
  ],
});