import fs from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const ROOT_DIR = join(__dirname, '..');
const PROTO_DIR = join(ROOT_DIR, '../proto');
const CONTENT_VERSIONS_DIR = join(ROOT_DIR, 'content/versions');
const CONTENT_API_DIR = join(ROOT_DIR, 'content/reference/api');
const PUBLIC_DIR = join(ROOT_DIR, 'public/docs');
const OPENAPI_DIR = join(ROOT_DIR, 'openapi');
const VERSIONS_FILE = join(ROOT_DIR, 'content/versions.json');
const NEXT_DIR = join(ROOT_DIR, '.next');

const targets = [
  NEXT_DIR,
  OPENAPI_DIR,
  join(ROOT_DIR, '.source'),
  join(ROOT_DIR, '.temp'),
  VERSIONS_FILE,
  CONTENT_API_DIR,
  CONTENT_VERSIONS_DIR,
];

async function run() {
  console.log('Cleaning up generated artifacts...');

  for (const target of targets) {
    if (fs.existsSync(target)) {
      console.log(`Removing ${target}...`);
      fs.rmSync(target, { recursive: true, force: true });
    }
  }

  // Versioned content and protos
  const cleanVersioned = (dir, pattern) => {
    if (!fs.existsSync(dir)) return;
    const files = fs.readdirSync(dir);
    for (const file of files) {
      if (file === 'latest' || (file.startsWith('v') && (file === 'vTest' || !isNaN(parseInt(file[1]))))) {
        const p = join(dir, file);
        console.log(`Removing ${p}...`);
        fs.rmSync(p, { recursive: true, force: true });
      }
    }
  };

  cleanVersioned(PROTO_DIR);
  cleanVersioned(PUBLIC_DIR);

  console.log('Cleanup complete.');
}

run().catch(err => {
  console.error(err);
  process.exit(1);
});
