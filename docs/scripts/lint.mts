import { register } from 'node:module';
import { pathToFileURL } from 'node:url';
import { existsSync } from 'node:fs';
import { join, dirname, resolve } from 'node:path';

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

  const scanned = await scanURLs({
    preset: 'next',
    populate: {
      'docs/[[...slug]]': pages.map((page: any) => {
        return {
          value: {
            slug: page.slugs,
          },
          hashes: getHeadings(page),
        };
      }),
    },
  });

  const files = await getFiles();

  const linkErrors = await validateFiles(files, {
      scanned,
      markdown: {
        components: {
          Card: { attributes: ['href'] },
        },
      },
      checkRelativePaths: 'as-url',
    });

  printErrors(linkErrors, false);

  const imageErrors = checkImages(files);

  if (linkErrors.length === 0 && !imageErrors) {
    console.log('\n✅ All checks passed: No broken links or images found.');
  }

  if (linkErrors.length > 0 || imageErrors) {
    process.exit(1);
  }
}

function checkImages(files: FileObject[]): boolean {
  let hasErrors = false;

  for (const file of files) {
    const content = file.content;
    const imageRegex = /!\[.*?\]\((.*?)\)|<img.*?src=["'](.*?)["']/g;
    let match;

    while ((match = imageRegex.exec(content)) !== null) {
      let imagePath = match[1] || match[2];

      if (match[1]) {
        // Markdown link: remove title part of the link (e.g. /img.png "title")
        imagePath = imagePath.trim().split(/\s+/)[0];
      }
      
      // Ignore external links
      if (imagePath.startsWith('http') || imagePath.startsWith('https') || imagePath.startsWith('data:')) {
        continue;
      }

      // Remove query parameters or anchors
      imagePath = imagePath.split('?')[0].split('#')[0];
      
      // Decode URI components (e.g. %20 -> space)
      imagePath = decodeURIComponent(imagePath);

      let fullPath;
      if (imagePath.startsWith('/')) {
        // Absolute path relative to public folder
        fullPath = join(PUBLIC_ROOT, imagePath);
      } else {
        // Relative path relative to the markdown file
        fullPath = resolve(dirname(file.path), imagePath);
      }

      if (!existsSync(fullPath)) {
        console.error(`Broken image link in ${file.path}: ${imagePath}`);
        hasErrors = true;
      }
    }
  }

  return hasErrors;
}

function getHeadings({ data }: any): string[] {
  return data.toc.map((item: any) => item.url.slice(1));
}

import { readFile } from 'node:fs/promises';

async function getFiles() {
  const pages = [...source.getPages(), ...versionSource.getPages()];
  const promises = pages.map(
    async (page: any): Promise<FileObject> => ({
      path: page.file?.path || page.absolutePath, 
      content: await readFile(page.file?.path || page.absolutePath, 'utf8'), 
      url: page.url,
      data: page.data,
    }),
  );
  return Promise.all(promises);
}

await checkLinks();
