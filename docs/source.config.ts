import {
  defineConfig,
  defineDocs,
  frontmatterSchema,
  metaSchema,
} from 'fumadocs-mdx/config';
// @ts-ignore
import remarkHeadingId from 'remark-heading-id';

// You can customise Zod schemas for frontmatter and `meta.json` here
// see https://fumadocs.dev/docs/mdx/collections
export const docs = defineDocs({
  dir: 'content',
  docs: {
    schema: frontmatterSchema,
    files: ['**/!(_|v)*.md', '**/!(_|v)*.mdx'], // Exclude versioned folders
    postprocess: {
      includeProcessedMarkdown: true,
    },
  },
  meta: {
    schema: metaSchema,
    files: ['**/meta.json'],
  },
});

export const versions = defineDocs({
  dir: 'content',
  docs: {
    schema: frontmatterSchema,
    files: ['v*/!(_|v)*.md', 'v*/!(_|v)*.mdx', 'v*/**/!(_|v)*.md', 'v*/**/!(_|v)*.mdx'], // Include only versioned folders from content
    postprocess: {
      includeProcessedMarkdown: true,
    },
  },
  meta: {
    schema: metaSchema,
    files: ['v*/meta.json', 'v*/**/meta.json'],
  },
});

export default defineConfig({
  mdxOptions: {
    remarkPlugins: [[remarkHeadingId, { defaults: true }]],
    rehypeCodeOptions: {
      themes: {
        light: 'github-light',
        dark: 'github-dark',
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
