import fs from 'node:fs';
import path from 'node:path';

// Adjust path if needed
import { guidesSidebar, apisSidebar, legalSidebar } from '../lib/sidebar-data';

const CONTENT_DIR = path.join(process.cwd(), 'content');

function ensureDirectoryExistence(filePath: string) {
  const dirname = path.dirname(filePath);
  if (fs.existsSync(dirname)) {
    return true;
  }
  ensureDirectoryExistence(dirname);
  fs.mkdirSync(dirname);
}

function generateContent(title: string, description: string) {
  const safeDescription = description.replace(/"/g, '\\"');

  return `---
title: "${title}"
description: "${safeDescription}"
---

{/* THIS FILE IS AUTO-GENERATED FROM SIDEBAR-DATA.
  ANY MANUAL CHANGES WILL BE OVERWRITTEN.
*/}

import { Cards } from 'fumadocs-ui/components/card';

<Cards />
`;
}

function traverse(items: readonly any[]) {
  items.forEach((item) => {
    if (
      item.type === 'category' &&
      item.link &&
      item.link.type === 'generated-index' &&
      item.link.slug
    ) {
      const slug = item.link.slug;

      // Remove leading slash or 'docs/' prefix
      const cleanSlug = slug.replace(/^\/|docs\/|\/$/g, '');
      const filePath = path.join(CONTENT_DIR, cleanSlug, 'index.mdx');

      // --- CHANGE: Always overwrite the file ---
      console.log(`[+] Generating Virtual Page: ${cleanSlug}`);
      ensureDirectoryExistence(filePath);

      const content = generateContent(
        item.link.title || item.label || 'Overview',
        item.link.description || ''
      );

      fs.writeFileSync(filePath, content);
    }

    if (item.items) {
      traverse(item.items);
    }
  });
}

console.log('--- Scanning Sidebar for Virtual Pages ---');
traverse(guidesSidebar);
traverse(apisSidebar);
traverse(legalSidebar);
console.log('--- Done ---');