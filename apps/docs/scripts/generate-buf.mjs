import fs from 'fs';
import { join, dirname, resolve } from 'path';
import { spawnSync } from 'child_process';
import { fileURLToPath } from 'url';
import os from 'os';

const __dirname = dirname(fileURLToPath(import.meta.url));
const ROOT_DIR = join(__dirname, '..');
const PROTO_DIR = join(ROOT_DIR, '../../proto');
const OPENAPI_DIR = join(ROOT_DIR, 'openapi');
const VERSIONS_FILE = join(ROOT_DIR, 'content/versions.json');
const REPO_URL = 'https://github.com/zitadel/zitadel.git';

async function run() {
  if (!fs.existsSync(VERSIONS_FILE)) {
      console.error('versions.json not found. Run fetch-docs.mjs first.');
      process.exit(1);
  }

  const versions = JSON.parse(fs.readFileSync(VERSIONS_FILE, 'utf8'));
  console.log(`Processing ${versions.length} versions from versions.json`);

  const baseTempDir = fs.mkdtempSync(join(os.tmpdir(), 'zitadel-buf-'));
  // Use a subdirectory for local generation to avoid pollution
  const localGenDir = join(baseTempDir, 'local'); 
  fs.mkdirSync(localGenDir, { recursive: true });

  const templatePath = resolve(join(ROOT_DIR, 'buf.gen.yaml'));

  try {
    for (const v of versions) {
      if (v.type === 'external') continue; // Skip archive links

      const label = v.param;
      const outputPath = resolve(join(OPENAPI_DIR, label));
      
      console.log(`\n--- Generating OpenAPI specs for ${label} ---`);
      
      fs.rmSync(outputPath, { recursive: true, force: true });
      fs.mkdirSync(outputPath, { recursive: true });

      // Determine buf input based on refType
      let bufInput;
      if (v.refType === 'local') {
          // Point to local proto directory (repo root/proto)
          bufInput = PROTO_DIR; // e.g. /home/ffo/git/zitadel/zitadel/proto
      } else {
          // Construct git input
          // Format: <url>#<refType>=<ref>,subdir=<subdir>
          // Example: https://github.com/zitadel/zitadel.git#tag=v4.11.0,subdir=proto
          const refPart = v.refType === 'branch' ? `branch=${v.ref}` : `tag=${v.ref}`;
          bufInput = `${REPO_URL}#${refPart},subdir=proto`;
      }
      
      console.log(`Using input: ${bufInput}`);

      const result = spawnSync('npx', [
        '@bufbuild/buf', 'generate',
        bufInput,
        '--template', templatePath,
        '--output', outputPath
      ], {
        cwd: localGenDir, // Run in temp dir
        stdio: 'inherit',
        env: {
            ...process.env,
            BUF_TOKEN: '851d3e2519b882d9e6da46eadec5e3ccc6a966dae0ce5e78dd285d9f912e35fd'
        }
      });

      if (result.status !== 0) {
        console.error(`Failed to generate OpenAPI for ${label}`);
        process.exit(1);
      } else {
        console.log(`Successfully generated OpenAPI for ${label}`);
      }
    }
  } finally {
    fs.rmSync(baseTempDir, { recursive: true, force: true });
  }
}

run().catch(err => {
  console.error(err);
  process.exit(1);
});
