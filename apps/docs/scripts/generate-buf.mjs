import fs from 'fs';
import { glob } from 'glob';
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
  // Simple p-limit implementation to avoid adding dependency
  const pLimit = (concurrency) => {
    const queue = [];
    let active = 0;

    const next = () => {
      active--;
      if (queue.length > 0) {
        queue.shift()();
      }
    };

    const run = async (fn) => {
      if (active >= concurrency) {
        await new Promise((resolve) => queue.push(resolve));
      }
      active++;
      try {
        return await fn();
      } finally {
        next();
      }
    };

    return run;
  };
  
  const limit = pLimit(2);

  await Promise.all(versions.map(v => limit(async () => {
    if (v.type === 'external') return; // Skip archive links

    const label = v.param;
    const outputPath = resolve(join(OPENAPI_DIR, label));
      
    console.log(`\n--- Generating OpenAPI specs for ${label} ---`);
      
    fs.rmSync(outputPath, { recursive: true, force: true });
    fs.mkdirSync(outputPath, { recursive: true });

    // Create a unique temp dir for this specific generation task to avoid conflicts
    const taskTempDir = join(baseTempDir, label);
    fs.mkdirSync(taskTempDir, { recursive: true });

    // Determine buf input based on refType
    let bufInput;
    if (v.refType === 'local') {
        // Point to local proto directory (repo root/proto)
        bufInput = PROTO_DIR; 
    } else {
        const refPart = v.refType === 'branch' ? `branch=${v.ref}` : `tag=${v.ref}`;
        bufInput = `${REPO_URL}#${refPart},subdir=proto`;
    }
      
    console.log(`Using input for ${label}: ${bufInput}`);


    // Use spawn (async) instead of spawnSync to avoid blocking the event loop
    await new Promise(async (resolvePromise, rejectPromise) => {
        // Dynamic discovery of excluded paths
        const getExcludedPaths = async () => {
             const patterns = ['**/v3alpha'];
             const excluded = new Set();
             
             if (v.refType === 'local') {
                 // For local, use glob against PROTO_DIR
                 // patterns are relative to PROTO_DIR, but glob needs to search inside PROTO_DIR
                 // The result of glob will be paths relative to PROTO_DIR if cwd is PROTO_DIR
                 try {
                    const files = await glob(patterns, { cwd: PROTO_DIR, nodir: false });
                    // We only want directories that contain v3alpha, or files?
                    // buf --exclude-path expects directories or files.
                    // If we found directories ending in v3alpha, usage is fine.
                    // glob '**/v3alpha' matches directories or files ending in v3alpha.
                    // Ensure we pass absolute paths for local input to avoid CWD resolution issues
                    files.forEach(f => excluded.add(resolve(PROTO_DIR, f)));
                 } catch (e) {
                     console.warn('Failed to glob excluded paths locally', e);
                 }
             } else {
                 // For remote (tag/branch), use git ls-tree
                 // We need to run git ls-tree -r <ref> --name-only
                 // and filter for lines matching /v3alpha/
                 // We run this in ROOT_DIR (where .git presumably is available via parent)
                 // The paths returned by git ls-tree are relative to repo root.
                 // buf input is relative to 'proto' subdir.
                 // So we need to strip 'proto/' prefix.
                 try {
                     const ref = v.ref;
                     const gitCmd = spawnSync('git', ['ls-tree', '-r', '--name-only', ref], { cwd: ROOT_DIR, encoding: 'utf-8' });
                     if (gitCmd.error) throw gitCmd.error;
                     const files = gitCmd.stdout.split('\n');
                     files.forEach(f => {
                         if (f.includes('v3alpha') && f.startsWith('proto/')) {
                             // strip 'proto/' (6 chars)
                             excluded.add(f.substring(6));
                         }
                     });
                 } catch (e) {
                     console.warn(`Failed to list tree for ${v.param}`, e);
                 }
             }
             return Array.from(excluded);
        };

        const excludedPaths = await getExcludedPaths();
        if (excludedPaths.length > 0) {
            console.log(`Excluding ${excludedPaths.length} paths matching v3alpha`);
        }

        import('child_process').then(({ spawn }) => {
            const args = [
              '@bufbuild/buf', 'generate',
              bufInput,
              '--template', templatePath,
              '--output', outputPath
            ];

            // Add excluded paths
            for (const excludedPath of excludedPaths) {
                args.push('--exclude-path', excludedPath);
            }

            const child = spawn('npx', args, {
              cwd: taskTempDir, 
              stdio: 'inherit',
              env: {
                  ...process.env,
                  BUF_TOKEN: '851d3e2519b882d9e6da46eadec5e3ccc6a966dae0ce5e78dd285d9f912e35fd'
              }
            });

            child.on('close', (code) => {
                if (code !== 0) {
                    rejectPromise(new Error(`Failed to generate OpenAPI for ${label} (exit code ${code})`));
                } else {
                    console.log(`Successfully generated OpenAPI for ${label}`);
                    resolvePromise();
                }
            });
              
            child.on('error', (err) => {
                rejectPromise(err);
            });
        });
    });
  })));
  } finally {
    fs.rmSync(baseTempDir, { recursive: true, force: true });
  }
}

run().catch(err => {
  console.error(err);
  process.exit(1);
});
