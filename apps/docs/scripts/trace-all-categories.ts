// Non-interactive counterpart to trace-category-cli.ts, built for CI.
//
//   pnpm exec tsx scripts/trace-all-categories.ts                # every discovered category
//   pnpm exec tsx scripts/trace-all-categories.ts --only user,feature   # just these
//
// For each selected category: runs the Go call-graph tracer
// (internal/tools/errortrace --proto) and writes its output to that
// service's tracingDir as `${category}-ops.json` — the same file the
// interactive picker produces, so either path is safe to mix. Once every
// selected category is traced, regenerates the docs data
// (generate-endpoint-errors.ts) exactly once at the end.
//
// Used by two GitHub Actions workflows: a manual one-time backfill (no
// --only, traces everything not covered yet) and a per-PR watcher (--only
// scoped to whatever categories a PR's proto/handler changes touched).
import { execFileSync } from 'child_process';
import { existsSync, mkdirSync, writeFileSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import { discoverTraceableServices, type ServiceConfig } from './generate-endpoint-errors';

const __dirname = dirname(fileURLToPath(import.meta.url));
const DOCS_ROOT = join(__dirname, '..');
const REPO_ROOT = join(DOCS_ROOT, '../..');

function parseOnly(argv: string[]): Set<string> | null {
  const i = argv.indexOf('--only');
  if (i === -1) return null;
  const value = argv[i + 1];
  if (!value) {
    console.error('[trace-all-categories] --only requires a comma-separated list of categories');
    process.exit(1);
  }
  return new Set(value.split(',').map((s) => s.trim()).filter(Boolean));
}

function traceOne(svc: ServiceConfig): number {
  console.log(`[trace-all-categories] tracing ${svc.category} (${svc.service})...`);
  const output = execFileSync('go', ['run', './internal/tools/errortrace', '--proto', svc.protoFile], {
    cwd: REPO_ROOT,
    stdio: ['ignore', 'pipe', 'inherit'],
  }).toString();

  const traced = JSON.parse(output);
  const opCount = Object.keys(traced).length;
  if (opCount === 0) {
    console.error(`[trace-all-categories] ${svc.category}: tracer found no operations — skipping write`);
    return 0;
  }

  if (!existsSync(svc.tracingDir)) mkdirSync(svc.tracingDir, { recursive: true });
  const outFile = join(svc.tracingDir, `${svc.category}-ops.json`);
  writeFileSync(outFile, output);
  console.log(`[trace-all-categories] wrote ${opCount} operation(s) to ${outFile}`);
  return opCount;
}

function main() {
  const only = parseOnly(process.argv.slice(2));
  const services = discoverTraceableServices();
  const selected = only ? services.filter((s) => only.has(s.category)) : services;

  if (only) {
    const missing = [...only].filter((c) => !services.some((s) => s.category === c));
    if (missing.length > 0) {
      console.error(`[trace-all-categories] unknown categor${missing.length === 1 ? 'y' : 'ies'}: ${missing.join(', ')}`);
      process.exit(1);
    }
  }
  if (selected.length === 0) {
    console.log('[trace-all-categories] nothing to trace.');
    return;
  }

  let totalOps = 0;
  for (const svc of selected) {
    totalOps += traceOne(svc);
  }

  console.log('[trace-all-categories] regenerating docs data...');
  execFileSync('pnpm', ['exec', 'tsx', 'scripts/generate-endpoint-errors.ts'], { cwd: DOCS_ROOT, stdio: 'inherit' });

  console.log(`\n[trace-all-categories] done. ${selected.length} categor${selected.length === 1 ? 'y' : 'ies'} traced, ${totalOps} operation(s) total.`);
}

main();
