import fs from 'node:fs/promises';
import path from 'node:path';
import yaml from 'js-yaml';

const V1_BASE_PATHS: Record<string, string> = {
  'admin.openapi.json': '/admin/v1',
  'auth.openapi.json': '/auth/v1',
  'management.openapi.json': '/management/v1',
  'system.openapi.json': '/system/v1',
};


async function walk(dir: string): Promise<string[]> {
  let files: string[] = [];
  const list = await fs.readdir(dir);
  for (const file of list) {
    const filepath = path.join(dir, file);
    const stat = await fs.stat(filepath);
    if (stat.isDirectory()) {
      files = files.concat(await walk(filepath));
    } else {
      if (file.endsWith('.openapi.json')) {
        files.push(filepath);
      }
    }
  }
  return files;
}

function removeAdditionalPropertiesFalse(obj: any) {
  if (typeof obj !== 'object' || obj === null) return;

  if (Array.isArray(obj)) {
    for (const item of obj) {
      removeAdditionalPropertiesFalse(item);
    }
    return;
  }

  if (obj.additionalProperties === false) {
    delete obj.additionalProperties;
  }

  for (const key in obj) {
    removeAdditionalPropertiesFalse(obj[key]);
  }
}

async function cleanOpenApi() {
  const openApiDir = path.join(process.cwd(), 'openapi');
  try {
    await fs.access(openApiDir);
  } catch {
    console.log('No openapi directory found.');
    return;
  }

  const files = await walk(openApiDir);

  for (const file of files) {
    const content = await fs.readFile(file, 'utf-8');
    const filename = path.basename(file);

    let doc: any;
    try {
      doc = JSON.parse(content);
    } catch (e) {
      console.error(`Error parsing JSON for ${file}:`, e);
      continue;
    }

    // 1. Inject servers block for v1 APIs
    if (filename in V1_BASE_PATHS) {
      doc.servers = [
        {
          url: V1_BASE_PATHS[filename],
          description: "Zitadel " + filename.split('.')[0] + " API v1"
        }
      ];
    }

    // 2. Remove additionalProperties: false
    removeAdditionalPropertiesFalse(doc);

    await fs.writeFile(file, JSON.stringify(doc, null, 2));
  }
}

cleanOpenApi();
