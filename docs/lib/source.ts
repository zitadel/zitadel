import { docs, versions } from '../.source/server';
import { type InferPageType, loader } from 'fumadocs-core/source';
import { lucideIconsPlugin } from 'fumadocs-core/source/lucide-icons';

// See https://fumadocs.dev/docs/headless/source-api for more info
export const source = loader({
  baseUrl: '/',
  source: docs.toFumadocsSource(),
  plugins: [lucideIconsPlugin()],
});

export const versionSource = loader({
  baseUrl: '/',
  source: versions.toFumadocsSource(),
  plugins: [lucideIconsPlugin()],
});

export function getPage(slugs: string[] | undefined) {
  const safeSlugs = slugs || [];
  // If the first slug matches a known version pattern (e.g., starts with 'v' and is in our list), use versionSource
  // For simplicity, we check if the page exists in versionSource first if it looks like a version
  if (safeSlugs.length > 0 && safeSlugs[0].startsWith('v')) {
    const page = versionSource.getPage(safeSlugs);
    if (page) return { page, source: versionSource };
  }

  return { page: source.getPage(safeSlugs), source: source };
}

export function getPageImage(page: InferPageType<typeof source>) {
  const segments = [...page.slugs, 'image.png'];

  return {
    segments,
    url: `/og/docs/${segments.join('/')}`,
  };
}

export async function getLLMText(page: InferPageType<typeof source>) {
  const processed = await page.data.getText('processed');

  return `# ${page.data.title}

${processed}`;
}
