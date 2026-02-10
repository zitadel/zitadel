import { generateFiles } from 'fumadocs-openapi';
import { createOpenAPI } from 'fumadocs-openapi/server';
import { writeFileSync, mkdirSync, readdirSync, lstatSync, readFileSync, existsSync } from 'fs';
import path, { join, dirname, basename } from 'path';
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

/**
 * Workaround for fumadocs-openapi bug: when a schema has both root-level
 * `required` fields AND `oneOf`/`anyOf`, fumadocs' `...item` spread overwrites
 * `schema.required` with `item.required`, losing root-level required fields.
 *
 * Fix: merge parent `required` into each variant and remove from parent.
 */
function fixSchemaRequired(schema: any): void {
  if (!schema || typeof schema !== 'object') return;

  const union = schema.oneOf ?? schema.anyOf;
  if (Array.isArray(union) && Array.isArray(schema.required) && schema.required.length > 0) {
    for (const variant of union) {
      if (typeof variant !== 'object' || variant === null) continue;
      variant.required = [...new Set([...schema.required, ...(variant.required ?? [])])];
    }
    delete schema.required;
  }

  // Recurse into nested schemas
  for (const prop of Object.values(schema.properties ?? {})) {
    fixSchemaRequired(prop);
  }
  for (const item of (schema.oneOf ?? schema.anyOf ?? schema.allOf ?? [])) {
    fixSchemaRequired(item);
  }
  if (schema.items) fixSchemaRequired(schema.items);
}

function fixOneOfRequired(doc: any): void {
  for (const schema of Object.values(doc.components?.schemas ?? {})) {
    fixSchemaRequired(schema);
  }
  for (const pathItem of Object.values(doc.paths ?? {})) {
    for (const op of Object.values(pathItem as Record<string, any>)) {
      if (!op || typeof op !== 'object') continue;
      const bodySchema = (op as any)?.requestBody?.content?.['application/json']?.schema;
      if (bodySchema) fixSchemaRequired(bodySchema);
    }
  }
}

async function generateVersionApiDocs(version: string) {
  const sourceRoot = join(OPENAPI_ROOT, version);
  if (!existsSync(sourceRoot)) return;

  const outputRoot = version === 'latest'
    ? join(CONTENT_ROOT, 'reference/api')
    : join(CONTENT_VERSIONS_ROOT, `${version}/reference/api`);

  console.log(`Generating API docs for version: ${version} -> ${outputRoot}`);
  mkdirSync(outputRoot, { recursive: true });

  const specs = await glob('**/*.openapi.json', { cwd: sourceRoot });
  const services: string[] = [];

  for (const specPath of specs) {
    const fullPath = join(sourceRoot, specPath);
    const content = readFileSync(fullPath, 'utf8');
    const doc = JSON.parse(content) as any;

    if (!doc.paths || Object.keys(doc.paths).length === 0) continue;

    // Fix required fields lost during fumadocs oneOf rendering
    fixOneOfRequired(doc);
    writeFileSync(fullPath, JSON.stringify(doc, null, 2));

    let service = basename(specPath, '.openapi.json');
    if (service.endsWith('_service')) {
      service = service.slice(0, -8);
    }

    // For services in subdirectories (like resource/userschema), 
    // we want a unique but readable name.
    const relDir = dirname(specPath);
    const folderPrefix = relDir !== '.' && !relDir.startsWith('zitadel')
      ? relDir.split('/').pop() + '-'
      : '';
    const uniqueService = folderPrefix + service;

    const outputDir = join(outputRoot, uniqueService);
    services.push(uniqueService);

    const api = createOpenAPI({
      input: [path.relative(ROOT_DIR, fullPath)],
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
    pages: services
      .filter(s => !s.includes('beta') && !s.includes('alpha'))
      .sort()
  };

  writeFileSync(
    join(outputRoot, 'meta.json'),
    JSON.stringify(meta, null, 2)
  );
}

async function fixAllGeneratedLinks() {
  console.log('Post-processing: Fixing API links...');
  const files = await glob('**/*.{md,mdx}', { cwd: CONTENT_ROOT });
  const fileIndex = new Map<string, Map<string, string>>(); // version -> (Name -> absoluteURL)

  const getVersion = (filePath: string) => {
    const parts = filePath.split(path.sep);
    if (parts[0].startsWith('v4.')) return parts[0];
    return 'latest';
  };

  const getUrl = (filePath: string) => {
    return '/' + filePath.replace(/\.(md|mdx)$/, '').split(path.sep).join('/');
  };

  // Build index
  for (const file of files) {
    const version = getVersion(file);
    if (!fileIndex.has(version)) fileIndex.set(version, new Map());
    const versionMap = fileIndex.get(version)!;

    const name = basename(file.replace(/\.md$/, ''), '.mdx');
    const parts = name.split('.');

    const operation = parts[parts.length - 1];
    const service = parts.length > 1 ? parts[parts.length - 2] : null;
    const url = getUrl(file);

    // Prioritize non-beta/non-alpha in the index
    const isBeta = name.includes('beta') || name.includes('alpha');
    const existing = versionMap.get(operation.toLowerCase());
    if (!existing || (!isBeta && existing.includes('beta'))) {
      versionMap.set(operation.toLowerCase(), url);
    }

    if (service) {
      const serviceOp = `${service}.${operation}`.toLowerCase();
      const existingServiceOp = versionMap.get(serviceOp);
      if (!existingServiceOp || (!isBeta && existingServiceOp.includes('beta'))) {
        versionMap.set(serviceOp, url);
      }
    }
    versionMap.set(name.toLowerCase(), url);
  }

  let fixedCount = 0;
  for (const file of files) {
    const fullPath = join(CONTENT_ROOT, file);
    const content = readFileSync(fullPath, 'utf-8');
    const version = getVersion(file);
    const versionMap = fileIndex.get(version);

    // [Text](apis/resources/service_name/file_name.api.mdx)
    const linkRegex = /\[([^\]]+)\]\(([\/]?apis\/resources\/([^\/]+)\/([^\/)]+)\.api\.mdx)\)/g;

    let modified = false;
    let newContent = content.replace(linkRegex, (match, text, fullLink, serviceSlug, fileSlug) => {
      const toPascalCase = (s: string) => s.split(/[-_]/).map(p => p.charAt(0).toUpperCase() + p.slice(1)).join('');

      // Extraction strategy: fileSlug is "user-service-get-user-by-id"
      const serviceIndex = fileSlug.lastIndexOf('service-');
      const operationKebab = serviceIndex !== -1 ? fileSlug.slice(serviceIndex + 8) : fileSlug;
      const operation = toPascalCase(operationKebab);

      let targetUrl = versionMap?.get(operation.toLowerCase());

      if (!targetUrl && serviceSlug) {
        const serviceMatch = serviceSlug.match(/^([a-z_]+)_service/);
        if (serviceMatch) {
          const service = toPascalCase(serviceMatch[1]) + 'Service';
          targetUrl = versionMap?.get(`${service}.${operation}`.toLowerCase());
        }
      }

      if (targetUrl) {
        modified = true;
        fixedCount++;
        return `[${text}](${targetUrl})`;
      }
      return match;
    });

    // Fix v2beta links that were likely in the source proto comments
    // and also remove potential double /docs prefixing from source comments
    const internalLinkRegex = /\[([^\]]+)\]\(([\/]?docs\/)?reference\/api\/([^\/]+)\/([^\s)]+)\)/g;
    newContent = newContent.replace(internalLinkRegex, (match, text, docsPrefix, service, fileSlug) => {
      let targetFileSlug = fileSlug;
      let isV2Beta = fileSlug.includes('.v2beta.');

      if (isV2Beta) {
        const v2Target = fileSlug.replace('.v2beta.', '.v2.');
        // Only rename to v2 if the v2 target actually exists in our index
        if (versionMap?.has(v2Target.toLowerCase())) {
          targetFileSlug = v2Target;
        } else {
          // If no v2 target, and we excluded v2beta, this link is dead.
          // Return plain text instead of a broken link to satisfy the link checker.
          modified = true;
          return text;
        }
      }

      const targetUrl = `/reference/api/${service}/${targetFileSlug}`;
      if (docsPrefix || isV2Beta) {
        modified = true;
        return `[${text}](${targetUrl})`;
      }
      return match;
    });

    if (modified) {
      writeFileSync(fullPath, newContent);
    }
  }
  console.log(`Post-processing: Fixed ${fixedCount} links.`);
}

async function run() {
  const args = process.argv.slice(2);
  const onlyGenerate = args.includes('--only-generate');
  const onlyFix = args.includes('--only-fix');
  // Default to both if no specific flag is set
  const runAll = !onlyGenerate && !onlyFix;

  if (runAll || onlyGenerate) {
    if (!existsSync(OPENAPI_ROOT)) {
      console.error('OpenAPI root not found. Run generate-buf.mjs first.');
      process.exit(1);
    }

    const versions = readdirSync(OPENAPI_ROOT).filter(f => lstatSync(join(OPENAPI_ROOT, f)).isDirectory() && f !== 'zitadel');

    for (const version of versions) {
      await generateVersionApiDocs(version);
    }
  }

  if (runAll || onlyFix) {
    await fixAllGeneratedLinks();
  }
}

run().catch(err => {
  console.error(err);
  process.exit(1);
});
