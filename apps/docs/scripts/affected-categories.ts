// Maps a list of changed repo-relative file paths (one per line, from
// `git diff --name-only`, on stdin) to the v2 API categories they touch —
// i.e. did this change a category's .proto service file or anything under
// its Go handler package. Prints a comma-separated category list (empty
// string if none matched), for `trace-all-categories.ts --only <list>`.
//
//   git diff --name-only origin/main...HEAD | pnpm exec tsx scripts/affected-categories.ts
import { relative, dirname, resolve, join } from 'path';
import { fileURLToPath } from 'url';
import { discoverServices } from './generate-endpoint-errors';

const REPO_ROOT = join(dirname(fileURLToPath(import.meta.url)), '../../..');

function isInside(dir: string, file: string): boolean {
  const rel = relative(dir, file);
  return rel !== '' && !rel.startsWith('..') && !rel.startsWith('/');
}

async function main() {
  const chunks: Buffer[] = [];
  for await (const chunk of process.stdin) chunks.push(chunk as Buffer);
  const changed = Buffer.concat(chunks)
    .toString('utf8')
    .split('\n')
    .map((l) => l.trim())
    .filter(Boolean)
    // git diff --name-only prints repo-relative paths; resolve them against
    // the repo root rather than process.cwd(), which may be apps/docs.
    .map((f) => resolve(REPO_ROOT, f));

  const services = discoverServices();
  const affected = new Set<string>();
  for (const file of changed) {
    for (const svc of services) {
      if (file === svc.protoFile || isInside(dirname(svc.protoFile), file) || isInside(svc.goPackageDir, file)) {
        affected.add(svc.category);
      }
    }
  }
  process.stdout.write([...affected].join(','));
}

main();
