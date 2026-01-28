import { MetadataRoute } from 'next';
import { source, versionSource } from '@/lib/source';

export default function sitemap(): MetadataRoute.Sitemap {
    const baseUrl = 'https://zitadel.com/docs';

    const docsPages = source.getPages().map((page) => ({
        url: `${baseUrl}${page.url === '/' ? '' : page.url}`,
        lastModified: new Date(),
        changeFrequency: 'weekly' as const,
        priority: 0.8,
    }));

    const versionPages = versionSource.getPages().map((page) => ({
        url: `${baseUrl}${page.url}`,
        lastModified: new Date(),
        changeFrequency: 'monthly' as const,
        priority: 0.3,
    }));

    return [
        {
            url: baseUrl,
            lastModified: new Date(),
            changeFrequency: 'daily' as const,
            priority: 1,
        },
        ...docsPages,
    ].filter((item, index, self) =>
        index === self.findIndex((t) => t.url === item.url)
    );
}
