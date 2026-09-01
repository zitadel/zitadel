// Wraps the Go call-graph tracer (internal/tools/errortrace) — from "a Go
// handler file, a proto service, a list of categories, or a git diff" to
// "the docs page shows it", in one command. Four ways to say what to trace,
// pick at most one:
//
//   --go-file <path>        the service that owns this Go handler file
//                             pnpm generate:endpoint-trace -- --go-file internal/api/grpc/user/v2/metadata.go
//   --proto <path>           the service declared in this .proto file
//                             pnpm generate:endpoint-trace -- --proto proto/zitadel/user/v2/user_service.proto
//   --only <cat1,cat2>       just these discovered categories
//                             pnpm generate:endpoint-trace -- --only user,feature
//   --changed-since <ref>    whatever categories changed between <ref> and HEAD
//                             pnpm generate:endpoint-trace -- --changed-since origin/main
//   (no flags)                every discovered category
//
// Whichever mode, each traced service's tracing JSON goes to that service's
// tracingDir (apps/docs/scripts/endpoint-error-tracing/<category>/), then
// generate-endpoint-errors.ts runs once at the end so the docs data picks it
// up immediately.
//
// Used directly for local one-off tracing, and by the two GitHub Actions
// workflows: endpoint-error-tables-backfill.yml (no flags — everything not
// covered yet) and endpoint-error-tables-watch.yml (--changed-since the PR's
// base — only what that PR touched).
import { execFileSync } from 'child_process';
import { existsSync, mkdirSync, writeFileSync } from 'fs';
import { basename, dirname, isAbsolute, join, relative, resolve } from 'path';
import { fileURLToPath } from 'url';
import { discoverTraceableServices, type ServiceConfig } from './generate-endpoint-errors';

const __dirname = dirname(fileURLToPath(import.meta.url));
const DOCS_ROOT = join(__dirname, '..');
const REPO_ROOT = join(DOCS_ROOT, '../..');

function fail(message: string): never {
  console.error(`[trace-endpoint-errors] ${message}`);
  process.exit(1);
}

// path.relative() rather than a raw string-prefix check — resolve() and
// join() both produce backslash-separated paths on Windows, so comparing
// against a hand-appended '/' (as a plain prefix match would) silently never
// matches there. relative() is the portable way to ask "is childPath inside
// parentDir".
function isInside(parentDir: string, childPath: string): boolean {
  const rel = relative(parentDir, childPath);
  return rel !== '' && !rel.startsWith('..') && !isAbsolute(rel);
}

function parseArgs(argv: string[]) {
  const get = (flag: string) => {
    const i = argv.indexOf(flag);
    return i === -1 ? undefined : argv[i + 1];
  };
  return { goFile: get('--go-file'), proto: get('--proto'), only: get('--only'), changedSince: get('--changed-since') };
}

// Runs the Go tracer against one service's proto file and writes its output.
// Returns the operation count found (0 if the tracer found nothing, in
// which case nothing is written).
function traceService(svc: ServiceConfig, outFile: string): number {
  console.log(`[trace-endpoint-errors] tracing ${svc.category} (${svc.service})...`);
  const output = execFileSync('go', ['run', './internal/tools/errortrace', '--proto', svc.protoFile], {
    cwd: REPO_ROOT,
    stdio: ['ignore', 'pipe', 'inherit'],
  }).toString();

  const traced = JSON.parse(output);
  const opCount = Object.keys(traced).length;
  if (opCount === 0) {
    console.error(`[trace-endpoint-errors] ${svc.category}: tracer found no operations — skipping write`);
    return 0;
  }
  if (!existsSync(svc.tracingDir)) mkdirSync(svc.tracingDir, { recursive: true });
  writeFileSync(outFile, output);
  console.log(`[trace-endpoint-errors] wrote ${opCount} operation(s) to ${outFile}`);
  return opCount;
}

function regenerate() {
  console.log('[trace-endpoint-errors] regenerating docs data...');
  execFileSync('pnpm', ['exec', 'tsx', 'scripts/generate-endpoint-errors.ts'], { cwd: DOCS_ROOT, stdio: 'inherit' });
}

// Maps whatever changed between sinceRef and HEAD to the categories it
// touches — did it change a category's .proto file (or anything alongside
// it) or anything under its Go handler package.
function findAffectedCategories(sinceRef: string, services: ServiceConfig[]): Set<string> {
  const changed = execFileSync('git', ['diff', '--name-only', sinceRef, 'HEAD'], { cwd: REPO_ROOT })
    .toString()
    .split('\n')
    .map((l) => l.trim())
    .filter(Boolean)
    .map((f) => resolve(REPO_ROOT, f));

  const affected = new Set<string>();
  for (const file of changed) {
    for (const svc of services) {
      if (file === svc.protoFile || isInside(dirname(svc.protoFile), file) || isInside(svc.goPackageDir, file)) {
        affected.add(svc.category);
      }
    }
  }
  return affected;
}

function main() {
  const { goFile, proto, only, changedSince } = parseArgs(process.argv.slice(2));
  const modesGiven = [goFile, proto, only, changedSince].filter((v) => v !== undefined).length;
  if (modesGiven > 1) fail('pass at most one of --go-file, --proto, --only, --changed-since');

  const services = discoverTraceableServices();

  // Single-service mode: trace exactly the one service that owns this file.
  if (goFile || proto) {
    const goFileAbs = goFile ? resolve(process.cwd(), goFile) : undefined;
    const protoAbs = proto ? resolve(process.cwd(), proto) : undefined;
    const svc = goFileAbs
      ? services.find((s) => isInside(s.goPackageDir, goFileAbs))
      : services.find((s) => s.protoFile === protoAbs);
    if (!svc) {
      fail(
        `${goFileAbs ?? protoAbs} doesn't match any discovered service. That means its proto file or Go handler ` +
          `package doesn't exist yet — see discoverTraceableServices() in generate-endpoint-errors.ts.`,
      );
    }
    const base = goFileAbs ? basename(goFileAbs, '.go') : basename(protoAbs!, '.proto');
    const outFile = join(svc.tracingDir, `${base}-ops.json`);
    const opCount = traceService(svc, outFile);
    if (opCount === 0) fail('the tracer found no operations — check --go-file/--proto point at an actual handler/service file');
    regenerate();
    return;
  }

  // Category-list, changed-since, or (no flags) every-category mode.
  let selected: ServiceConfig[];
  if (changedSince) {
    const affected = findAffectedCategories(changedSince, services);
    // Printed unconditionally (even when empty) as a stable, parseable line
    // — endpoint-error-tables-watch.yml reads this back into $GITHUB_OUTPUT.
    console.log(`AFFECTED_CATEGORIES=${[...affected].join(',')}`);
    selected = services.filter((s) => affected.has(s.category));
  } else if (only) {
    const wanted = new Set(only.split(',').map((s) => s.trim()).filter(Boolean));
    const missing = [...wanted].filter((c) => !services.some((s) => s.category === c));
    if (missing.length > 0) fail(`unknown categor${missing.length === 1 ? 'y' : 'ies'}: ${missing.join(', ')}`);
    selected = services.filter((s) => wanted.has(s.category));
  } else {
    selected = services;
  }

  if (selected.length === 0) {
    console.log('[trace-endpoint-errors] nothing to trace.');
    return;
  }

  let totalOps = 0;
  for (const svc of selected) {
    const outFile = join(svc.tracingDir, `${svc.category}-ops.json`);
    totalOps += traceService(svc, outFile);
  }

  regenerate();
  console.log(
    `\n[trace-endpoint-errors] done. ${selected.length} categor${selected.length === 1 ? 'y' : 'ies'} traced, ${totalOps} operation(s) total.`,
  );
}

main();
