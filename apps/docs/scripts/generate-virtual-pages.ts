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

function traverse(items: readonly any[], baseDir: string) {
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
      const filePath = path.join(baseDir, cleanSlug, 'index.mdx');

      // --- CHANGE: Always overwrite the file ---
      console.log(`[+] Generating Virtual Page: ${path.relative(CONTENT_DIR, filePath)}`);
      ensureDirectoryExistence(filePath);

      const content = generateContent(
        item.link.title || item.label || 'Overview',
        item.link.description || ''
      );

      fs.writeFileSync(filePath, content);
    }

    if (item.items) {
      traverse(item.items, baseDir);
    }
  });
}

const VERSIONS_FILE = path.join(CONTENT_DIR, 'versions.json');

console.log('--- Scanning Sidebar for Virtual Pages ---');

// 1. Generate for Latest (Root Content)
console.log(`\nProcessing: Latest`);
traverse(guidesSidebar, CONTENT_DIR);
traverse(apisSidebar, CONTENT_DIR);
traverse(legalSidebar, CONTENT_DIR);

// 2. Generate for Versions
if (fs.existsSync(VERSIONS_FILE)) {
  const versions = JSON.parse(fs.readFileSync(VERSIONS_FILE, 'utf8'));

  versions.forEach((v: any) => {
    if (!v.param || v.param === 'latest' || !v.param.startsWith('v')) return;

    const versionDir = path.join(CONTENT_DIR, v.param);
    if (!fs.existsSync(versionDir)) return;

    console.log(`\nProcessing: ${v.param}`);
    // Assuming versions share the same structure/sidebar as they are mocked from this branch
    // independent of the "source" (git tag), if we want virtual pages we must generate them.
    traverse(guidesSidebar, versionDir);
    traverse(apisSidebar, versionDir);
    traverse(legalSidebar, versionDir);
  });
}

console.log('--- Done ---');