// Suppress MaxListenersExceededWarning from fumadocs-mdx internal concurrency
if (typeof process !== 'undefined') {
  process.setMaxListeners(30);
}

import {
  defineConfig,
  defineDocs,
  frontmatterSchema,
  metaSchema,
} from 'fumadocs-mdx/config';
import { z } from 'zod';
// @ts-ignore
import remarkHeadingId from 'remark-heading-id';

// You can customise Zod schemas for frontmatter and `meta.json` here
// see https://fumadocs.dev/docs/mdx/collections
export const docs = defineDocs({
  dir: 'content',
  docs: {
    schema: frontmatterSchema.extend({
      sidebar_label: z.string().optional(),
    }),
    files: ['**/*.md', '**/*.mdx', '!v*/**/*', '!**/_*'], // Exclude versioned folders at root and partials
    postprocess: {
      includeProcessedMarkdown: true,
    },
  },
  meta: {
    schema: metaSchema,
    files: ['**/meta.json', '!v*/**'],
  },
});

export const versions = defineDocs({
  dir: 'content',
  docs: {
    schema: frontmatterSchema.extend({
      sidebar_label: z.string().optional(),
    }),
    files: ['v*/**/*.md', 'v*/**/*.mdx', '!**/_*'], // Include only versioned folders from content
    postprocess: {
      includeProcessedMarkdown: true,
    },
  },
  meta: {
    schema: metaSchema,
    files: ['v*/meta.json', 'v*/**/meta.json'],
  },
});

import { readFileSync, existsSync } from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const findThemePath = () => {
  const possiblePaths = [
    path.resolve(__dirname, '../../packages/theme/shiki-theme.json'),
    path.resolve(__dirname, '../../../packages/theme/shiki-theme.json'),
    path.resolve(process.cwd(), '../../packages/theme/shiki-theme.json'),
  ];
  for (const p of possiblePaths) {
    if (existsSync(p)) return p;
  }
  return possiblePaths[0];
};

const shikiTheme = JSON.parse(readFileSync(findThemePath(), 'utf-8'));

export default defineConfig({
  mdxOptions: {
    remarkPlugins: [[remarkHeadingId, { defaults: true }]],
    rehypeCodeOptions: {
      themes: {
        light: shikiTheme,
        dark: shikiTheme,
      },
      langs: ['json', 'yaml', 'bash', 'sh', 'shell', 'http', 'nginx', 'dockerfile', 'go', 'python', 'javascript', 'typescript', 'tsx', 'jsx', 'css', 'html', 'csharp', 'java', 'xml', 'sql', 'properties', 'ini', 'diff', 'markdown', 'mdx'],
      // Map unknown languages to text or similar
      langAlias: {
        'env': 'bash',
        'curl': 'bash',
        'caddy': 'nginx',
        'conf': 'nginx',
        'HTTP': 'http',
        'JSON': 'json',
      },
    },
  },
});
