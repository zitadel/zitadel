import { generateFiles } from 'fumadocs-openapi';
import { createOpenAPI } from 'fumadocs-openapi/server';
import { writeFileSync, mkdirSync, readdirSync, lstatSync, readFileSync, existsSync } from 'fs';
import { join, dirname, basename } from 'path';
import { fileURLToPath } from 'url';
import { glob } from 'glob';
import yaml from 'js-yaml';

// Suppress "Generated: ..." logs to avoid Vercel log limits
const originalLog = console.log;
console.log = (...args) => {
  if (args.length > 0 && typeof args[0] === 'string' && args[0].startsWith('Generated: ')) {
    return;
  }
  originalLog(...args);
};

const __dirname = dirname(fileURLToPath(import.meta.url));
const ROOT_DIR = join(__dirname, '..');
const OPENAPI_ROOT = join(ROOT_DIR, 'openapi');
const CONTENT_ROOT = join(ROOT_DIR, 'content');
const CONTENT_VERSIONS_ROOT = join(ROOT_DIR, 'content');

async function generateVersionApiDocs(version: string) {
  const sourceRoot = join(OPENAPI_ROOT, version);
  if (!existsSync(sourceRoot)) return;

  const outputRoot = version === 'latest'
    ? join(CONTENT_ROOT, 'reference/api')
    : join(CONTENT_VERSIONS_ROOT, `${version}/reference/api`);

  console.log(`Generating API docs for version: ${version} -> ${outputRoot}`);
  mkdirSync(outputRoot, { recursive: true });

  const specs = await glob('**/*.openapi.yaml', { cwd: sourceRoot });
  const services: string[] = [];

  for (const specPath of specs) {
    const fullPath = join(sourceRoot, specPath);
    const content = readFileSync(fullPath, 'utf8');
    const doc = yaml.load(content) as any;

    if (!doc.paths || Object.keys(doc.paths).length === 0) continue;

    let service = basename(specPath, '.openapi.yaml');
    if (service.endsWith('_service')) {
      service = service.slice(0, -8);
    }

    // For services in subdirectories (like resource/userschema), 
    // we want a unique but readable name.
    const relDir = dirname(specPath);
    const folderPrefix = relDir !== '.' && !relDir.startsWith('openapi/zitadel')
      ? relDir.split('/').pop() + '-'
      : '';
    const uniqueService = folderPrefix + service;

    const outputDir = join(outputRoot, uniqueService);
    services.push(uniqueService);

    const api = createOpenAPI({
      input: [fullPath],
    });

    await generateFiles({
      input: api,
      output: outputDir,
      includeDescription: true,
    });

    const indexPath = join(outputDir, 'index.mdx');
    const title = uniqueService.split('-').map(s => s.charAt(0).toUpperCase() + s.slice(1)).join(' ');
    const indexContent = `---
title: ${title} API
---

API Reference for ${title}
`;
    writeFileSync(indexPath, indexContent);
  }

  // Generate meta.json
  const meta = {
    title: "APIs",
    pages: services.sort()
  };

  writeFileSync(
    join(outputRoot, 'meta.json'),
    JSON.stringify(meta, null, 2)
  );
}

async function run() {
  if (!existsSync(OPENAPI_ROOT)) {
    console.error('OpenAPI root not found. Run generate-buf.mjs first.');
    process.exit(1);
  }

  const versions = readdirSync(OPENAPI_ROOT).filter(f => lstatSync(join(OPENAPI_ROOT, f)).isDirectory());

  for (const version of versions) {
    await generateVersionApiDocs(version);
  }
}

run().catch(err => {
  console.error(err);
  process.exit(1);
});
