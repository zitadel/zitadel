import fs from 'fs';
import { join, dirname, resolve } from 'path';
import { spawnSync } from 'child_process';
import { fileURLToPath } from 'url';
import os from 'os';
import semver from 'semver';

const __dirname = dirname(fileURLToPath(import.meta.url));
const ROOT_DIR = join(__dirname, '..');
const PROTO_DIR = join(ROOT_DIR, '../proto');
const OPENAPI_DIR = join(ROOT_DIR, 'openapi');

async function run() {
  if (!fs.existsSync(PROTO_DIR)) {
    console.error(`Proto directory not found at ${PROTO_DIR}`);
    process.exit(1);
  }

  const versions = fs.readdirSync(PROTO_DIR).filter(file => {
      const stats = fs.lstatSync(join(PROTO_DIR, file));
      if (!stats.isDirectory()) return false;
      return file === 'latest' || file === 'vTest' || (file.startsWith('v') && semver.valid(file.includes('.') ? file : file + '.0.0'));
  });

  console.log(`Found versioned protos: ${versions.join(', ')}`);

  const baseTempDir = fs.mkdtempSync(join(os.tmpdir(), 'zitadel-buf-'));
  console.log(`Using base temp dir: ${baseTempDir}`);

  try {
    for (const version of versions) {
      const protoPath = resolve(join(PROTO_DIR, version));
      const outputPath = resolve(join(OPENAPI_DIR, version));

      console.log(`\n--- Generating OpenAPI specs for ${version} ---`);
      
      fs.rmSync(outputPath, { recursive: true, force: true });
      fs.mkdirSync(outputPath, { recursive: true });

      const versionTempDir = join(baseTempDir, version);
      fs.mkdirSync(versionTempDir, { recursive: true });

      console.log(`Copying protos to temp dir...`);
      fs.cpSync(protoPath, versionTempDir, { recursive: true });

      const templatePath = resolve(join(ROOT_DIR, 'buf.gen.yaml'));

      const result = spawnSync('npx', [
        '@bufbuild/buf', 'generate',
        '--template', templatePath,
        '--output', outputPath
      ], {
        cwd: versionTempDir,
        stdio: 'inherit'
      });

      if (result.status !== 0) {
        console.error(`Failed to generate OpenAPI for ${version}`);
      } else {
        console.log(`Successfully generated OpenAPI for ${version}`);
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
