import { register } from 'node:module';
import { pathToFileURL } from 'node:url';
import { join, dirname, resolve } from 'node:path';
import { writeFileSync, existsSync } from 'node:fs';
import { readFile, stat } from 'node:fs/promises';

register('./scripts/loader.mjs', pathToFileURL('./'));
register('fumadocs-mdx/node/loader', import.meta.url);

import {
  type FileObject,
  printErrors,
  scanURLs,
  validateFiles,
} from 'next-validate-link';

// Dynamic import to ensure loader is registered before importing source
const { source, versionSource } = await import('../lib/source');

const PUBLIC_ROOT = resolve('public');

async function checkLinks() {
  const pages = [...source.getPages(), ...versionSource.getPages()];
  console.log(`Total pages found: ${pages.length}`);
  if (pages.length > 0) {
    console.log(`First page slug: ${pages[0].slugs}`);
  }

  const files = await getFiles();

  // Load redirects and add them to the scanned URLs
  let redirects: any[] = [];
  if (existsSync('redirects.json')) {
    redirects = JSON.parse(await readFile('redirects.json', 'utf-8'));
  }

  const scanned = {
    urls: new Map<string, { hashes: string[] }>([
      // Load redirects first
      ...redirects.map((r: any): [string, { hashes: string[] }] => [r.source, { hashes: [] }]),
      // Load files second, so they overwrite redirects with actual headings
      ...files.map((f): [string, { hashes: string[] }] => [f.url!, { hashes: getHeadings(f) }]),
    ]),
    fallbackUrls: [],
  };

  writeFileSync('scanned-urls.json', JSON.stringify(Array.from(scanned.urls.keys()), null, 2));

  console.log(`Manually populated URLs count: ${scanned.urls.size}`);


  const linkErrors = await validateFiles(files, {
    scanned,
    baseUrl: '/docs',
    markdown: {
      components: {
        Card: { attributes: ['href'] },
      },
    },
    checkRelativePaths: 'as-url',
  });

  printErrors(linkErrors, false);

  const imageErrors = await checkImages(files);

  if (linkErrors.length === 0 && !imageErrors) {
    console.log('\n✅ All checks passed: No broken links or images found.');
  }

  if (linkErrors.length > 0 || imageErrors) {
    process.exit(1);
  }
}

async function checkImages(files: FileObject[]): Promise<boolean> {
  let hasErrors = false;

  await Promise.all(files.map(async (file) => {
    const content = file.content;
    const imageRegex = /!\[.*?\]\((.*?)\)|<img.*?src=[{"'](.*?)["'}]/g;
    const matches = Array.from(content.matchAll(imageRegex));

    await Promise.all(matches.map(async (match) => {
      let imagePath = match[1] || match[2];

      if (match[1]) {
        // Markdown link: remove title part of the link (e.g. /img.png "title")
        imagePath = imagePath.trim().split(/\s+/)[0];
      }

      // Ignore external links
      if (imagePath.startsWith('http') || imagePath.startsWith('https') || imagePath.startsWith('data:')) {
        return;
      }

      // Remove query parameters or anchors
      imagePath = imagePath.split('?')[0].split('#')[0];

      // Decode URI components (e.g. %20 -> space)
      imagePath = decodeURIComponent(imagePath);

      let fullPath;
      if (imagePath.startsWith('/')) {
        // Enforce /docs prefix for absolute paths
        if (!imagePath.startsWith('/docs/')) {
          console.error(`Broken image link in ${file.path}: ${imagePath} (must start with /docs/)`);
          hasErrors = true;
          return;
        }

        // Absolute path relative to public folder (strip /docs prefix)
        const relativeToPublic = imagePath.slice(5); // Remove '/docs'
        fullPath = join(PUBLIC_ROOT, relativeToPublic);
      } else {
        // Relative path relative to the markdown file
        fullPath = resolve(dirname(file.path), imagePath);
      }

      try {
        await stat(fullPath);
      } catch {
        console.error(`Broken image link in ${file.path}: ${imagePath}`);
        hasErrors = true;
      }
    }));
  }));

  return hasErrors;
}

function getHeadings({ data, content }: any): string[] {
  const headings = new Set<string>();

  if (data.toc && data.toc.length > 0) {
    data.toc.forEach((item: any) => headings.add(item.url.slice(1)));
  }
  if (data.structuredData?.headings) {
    data.structuredData.headings.forEach((h: any) => headings.add(h.id));
  }

  // Fallback: parse content directly
  if (content) {
    const lines = content.split('\n');
    let inCodeBlock = false;
    for (const line of lines) {
      if (line.trim().startsWith('```')) {
        inCodeBlock = !inCodeBlock;
        continue;
      }
      if (inCodeBlock) continue;

      const match = line.match(/^(#{1,6})\s+(.+)$/);
      if (match) {
        const title = match[2].trim();
        const slug = title
          .toLowerCase()
          .replace(/[^\w\s-]/g, '')
          .replace(/_/g, '-')
          .replace(/\s+/g, '-')
          .replace(/-+/g, '-');
        headings.add(slug);
      }
    }
  }

  return Array.from(headings);
}

async function getFiles() {
  const pages = [...source.getPages(), ...versionSource.getPages()];
  const promises = pages.map(
    async (page: any): Promise<FileObject> => ({
      path: page.file?.path || page.absolutePath,
      content: await readFile(page.file?.path || page.absolutePath, 'utf8'),
      url: page.url === '/' ? '/docs' : `/docs${page.url.startsWith('/') ? page.url : '/' + page.url}`,
      data: page.data,
    }),
  );
  return Promise.all(promises);
}

await checkLinks();
