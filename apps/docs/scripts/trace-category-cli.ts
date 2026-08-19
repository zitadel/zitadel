// Interactive picker: pnpm generate:endpoint-trace:category
//
// Lists every API category discovered by generate-endpoint-errors.ts's
// discoverServices() — the same categorization already used across the docs
// site under apps/docs/content/reference/api/, nothing invented here — lets
// you pick one by number, then traces every RPC in that whole service in one
// pass via the Go call-graph tracer (internal/tools/errortrace --proto) and
// regenerates the docs data.
//
// For tracing a single .go file or a single service by exact path instead of
// picking from a menu, use trace-endpoint-errors.ts directly.
import { createInterface } from 'readline/promises';
import { execFileSync } from 'child_process';
import { existsSync, mkdirSync, writeFileSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import { discoverServices, listUnwiredProtoServices, type ServiceConfig } from './generate-endpoint-errors';

const __dirname = dirname(fileURLToPath(import.meta.url));
const DOCS_ROOT = join(__dirname, '..');
const REPO_ROOT = join(DOCS_ROOT, '../..');

function printMenu(services: ServiceConfig[]) {
  console.log('\nZITADEL API v2 categories (from apps/docs/content/reference/api/):\n');
  services.forEach((s, i) => {
    console.log(`  ${i + 1}. ${s.category} — ${s.service}`);
    if (s.description) console.log(`     ${s.description}`);
    const preview = s.rpcNames.slice(0, 3).join(', ') + (s.rpcNames.length > 3 ? ', ...' : '');
    console.log(`     ${s.rpcNames.length} operation(s): ${preview}`);
    console.log();
  });

  const unwired = listUnwiredProtoServices();
  if (unwired.length > 0) {
    console.log('Not shown above — proto service found, but not fully wired up yet:\n');
    for (const u of unwired) {
      console.log(`  - ${u.category} (${u.service}): missing ${u.missing.join(' and ')}`);
    }
    console.log();
  }
}

async function pickService(rl: ReturnType<typeof createInterface>, services: ServiceConfig[]): Promise<ServiceConfig | undefined> {
  while (true) {
    const answer = (await rl.question(`Pick a category (1-${services.length}), or "q" to quit: `)).trim();
    if (answer.toLowerCase() === 'q') return undefined;
    const n = Number(answer);
    if (Number.isInteger(n) && n >= 1 && n <= services.length) return services[n - 1];
    console.log(`Not a valid number 1-${services.length} — try again.`);
  }
}

async function main() {
  const services = discoverServices();
  if (services.length === 0) {
    console.error('[trace-category-cli] no API categories discovered — is proto/zitadel present?');
    process.exit(1);
  }
  printMenu(services);

  // A single question() call, deliberately — reading a second line after the
  // stream has already delivered its last one is a known readline/promises
  // pitfall: once stdin (piped or, at EOF, even a closed TTY session) has no
  // more data, a *second* question() issued after the first one resolves can
  // hang forever instead of ever settling. One prompt sidesteps it entirely.
  const rl = createInterface({ input: process.stdin, output: process.stdout });
  let svc: ServiceConfig | undefined;
  try {
    svc = await pickService(rl, services);
  } finally {
    rl.close();
  }
  if (!svc) {
    console.log('Cancelled.');
    return;
  }
  console.log(`\nSelected: ${svc.category} (${svc.service}) — tracing all ${svc.rpcNames.length} operation(s)...`);

  console.log(`[trace-category-cli] go run ./internal/tools/errortrace --proto ${svc.protoFile}`);
  const output = execFileSync('go', ['run', './internal/tools/errortrace', '--proto', svc.protoFile], {
    cwd: REPO_ROOT,
    stdio: ['ignore', 'pipe', 'inherit'],
  }).toString();

  const traced = JSON.parse(output);
  const opCount = Object.keys(traced).length;
  if (opCount === 0) {
    console.error('[trace-category-cli] the tracer found no operations for this service — nothing written.');
    process.exit(1);
  }

  if (!existsSync(svc.tracingDir)) mkdirSync(svc.tracingDir, { recursive: true });
  const outFile = join(svc.tracingDir, `${svc.category}-ops.json`);
  writeFileSync(outFile, output);
  console.log(`[trace-category-cli] wrote ${opCount} operation(s) to ${outFile}`);

  console.log('[trace-category-cli] regenerating docs data...');
  execFileSync('pnpm', ['exec', 'tsx', 'scripts/generate-endpoint-errors.ts'], { cwd: DOCS_ROOT, stdio: 'inherit' });

  console.log(`\nDone. ${svc.category} (${svc.service}): ${opCount} operation(s) traced.`);
}

main();
