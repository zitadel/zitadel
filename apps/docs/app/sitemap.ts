import { MetadataRoute } from 'next';
import { source } from '@/lib/source';

export const revalidate = false;
export const dynamic = 'force-static';

export default function sitemap(): MetadataRoute.Sitemap {
  const baseUrl = 'https://zitadel.com/docs';

  const docsPages = source.getPages().map((page) => ({
    url: `${baseUrl}${page.url === '/' ? '' : page.url}`,
    changeFrequency: 'weekly' as const,
    priority: 0.8,
  }));

  return [
    {
      url: baseUrl,
      changeFrequency: 'daily' as const,
      priority: 1,
    },
    ...docsPages,
  ].filter(
    (item, index, self) => index === self.findIndex((t) => t.url === item.url),
  );
}
