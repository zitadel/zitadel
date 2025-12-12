import { register } from 'node:module';
import { pathToFileURL } from 'node:url';

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

async function checkLinks() {
  const pages = [...source.getPages(), ...versionSource.getPages()];
  
  const managementPage = pages.find(p => p.url === '/docs/references/api-v1/management');
  if (managementPage) {
    console.log('Lint: Management page found in source:', managementPage.url);
    console.log('Lint: Management page slugs:', managementPage.slugs);
  } else {
    console.log('Lint: Management page NOT found in source');
  }

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

  printErrors(
    await validateFiles(await getFiles(), {
      scanned,
      markdown: {
        components: {
          Card: { attributes: ['href'] },
        },
      },
      checkRelativePaths: 'as-url',
    }),
    true,
  );
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
