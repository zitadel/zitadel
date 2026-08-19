// Wraps the Go call-graph tracer (internal/tools/errortrace) so a dev can go
// from "a Go handler file or a proto service file" to "the docs page shows
// it" in one command — no manual code reading, no hand-written tracing JSON:
//
//   pnpm generate:endpoint-trace -- --go-file internal/api/grpc/user/v2/metadata.go
//   pnpm generate:endpoint-trace -- --proto proto/zitadel/user/v2/user_service.proto
//
// It runs the tracer, matches the input to a service discovered by
// generate-endpoint-errors.ts's discoverServices() (by goPackageDir/protoFile
// — nothing hand-typed to match against), writes the result into that
// service's tracingDir, then re-runs generate-endpoint-errors.ts so the docs
// page picks it up immediately.
//
// If --go-file/--proto doesn't match any discovered service, this stops with
// an error rather than guessing — that means the category's proto file, docs
// content directory, or Go handler package doesn't exist yet (or the given
// path is outside all of them); see discoverServices() in
// generate-endpoint-errors.ts for exactly what it checks.
//
// For a guided, "pick a category from a list" alternative to this, see
// trace-category-cli.ts (pnpm generate:endpoint-trace:category).
import { execFileSync } from 'child_process';
import { existsSync, mkdirSync, writeFileSync } from 'fs';
import { basename, dirname, isAbsolute, join, relative, resolve } from 'path';
import { fileURLToPath } from 'url';
import { discoverServices } from './generate-endpoint-errors';

const SERVICES = discoverServices();

const __dirname = dirname(fileURLToPath(import.meta.url));
const DOCS_ROOT = join(__dirname, '..');
const REPO_ROOT = join(DOCS_ROOT, '../..');

function parseArgs(argv: string[]) {
  const get = (flag: string) => {
    const i = argv.indexOf(flag);
    return i === -1 ? undefined : argv[i + 1];
  };
  return { goFile: get('--go-file'), proto: get('--proto'), goPackage: get('--go-package') };
}

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

function main() {
  const { goFile, proto, goPackage } = parseArgs(process.argv.slice(2));
  if ((goFile == null) === (proto == null)) {
    fail('pass exactly one of --go-file or --proto');
  }

  const goFileAbs = goFile ? resolve(process.cwd(), goFile) : undefined;
  const protoAbs = proto ? resolve(process.cwd(), proto) : undefined;

  const svc = goFileAbs
    ? SERVICES.find((s) => isInside(s.goPackageDir, goFileAbs))
    : SERVICES.find((s) => s.protoFile === protoAbs);
  if (!svc) {
    fail(
      `${goFileAbs ?? protoAbs} doesn't match any discovered service. That means its proto file, docs content ` +
        `directory, or Go handler package doesn't exist yet — see discoverServices() in generate-endpoint-errors.ts.`,
    );
  }

  const traceArgs = ['run', './internal/tools/errortrace'];
  if (goFileAbs) traceArgs.push('--go-file', goFileAbs);
  if (protoAbs) traceArgs.push('--proto', protoAbs);
  if (goPackage) traceArgs.push('--go-package', goPackage);

  console.log(`[trace-endpoint-errors] go ${traceArgs.join(' ')}`);
  const output = execFileSync('go', traceArgs, { cwd: REPO_ROOT, stdio: ['ignore', 'pipe', 'inherit'] }).toString();

  const traced = JSON.parse(output);
  const opCount = Object.keys(traced).length;
  if (opCount === 0) {
    fail('the tracer found no operations — check --go-file/--proto point at an actual handler/service file');
  }

  if (!existsSync(svc.tracingDir)) mkdirSync(svc.tracingDir, { recursive: true });
  const base = goFileAbs ? basename(goFileAbs, '.go') : basename(protoAbs!, '.proto');
  const outFile = join(svc.tracingDir, `${base}-ops.json`);
  writeFileSync(outFile, output);
  console.log(`[trace-endpoint-errors] wrote ${opCount} operation(s) to ${outFile}`);

  console.log('[trace-endpoint-errors] regenerating docs data...');
  execFileSync('pnpm', ['exec', 'tsx', 'scripts/generate-endpoint-errors.ts'], { cwd: DOCS_ROOT, stdio: 'inherit' });
}

main();
