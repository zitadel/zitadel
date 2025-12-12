import fs from 'node:fs/promises';
import path from 'node:path';
import yaml from 'js-yaml';

const V1_BASE_PATHS: Record<string, string> = {
  'admin.openapi.yaml': '/admin/v1',
  'auth.openapi.yaml': '/auth/v1',
  'management.openapi.yaml': '/management/v1',
  'system.openapi.yaml': '/system/v1',
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
      if (file.endsWith('.openapi.yaml')) {
        files.push(filepath);
      }
    }
  }
  return files;
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
    
    // 1. Inject servers block for v1 APIs
    if (filename in V1_BASE_PATHS) {
      try {
        const doc = yaml.load(content) as any;
        if (doc) {
          doc.servers = [
            {
              url: V1_BASE_PATHS[filename],
              description: "ZITADEL " + filename.split('.')[0] + " API v1"
            }
          ];
          await fs.writeFile(file, yaml.dump(doc, { noRefs: true }));
          console.log(`Added servers block to ${file}`);
          // Re-read content for the next step
        }
      } catch (e) {
        console.error(`Error processing YAML for ${file}:`, e);
      }
    }

    // Re-read content in case it was modified above, or use the original content
    const currentContent = await fs.readFile(file, 'utf-8');
    const lines = currentContent.split('\n');
    
    // Filter out "additionalProperties: false" lines
    // We assume it's on its own line with indentation
    const newLines = lines.filter(line => !line.trim().match(/^additionalProperties:\s*false/));
    
    if (lines.length !== newLines.length) {
        await fs.writeFile(file, newLines.join('\n'));
        console.log(`Cleaned ${file}`);
    }
  }
}

cleanOpenApi();
