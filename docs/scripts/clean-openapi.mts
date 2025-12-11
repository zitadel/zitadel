import fs from 'node:fs/promises';
import path from 'node:path';

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
    const lines = content.split('\n');
    
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
