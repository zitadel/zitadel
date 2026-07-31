// Emits apps/docs/components/EndpointErrors/data.json and stamps a one-line
// <EndpointErrors /> stub into each covered operation's generated MDX page.
//
// Two sources feed the data, merged per (service, operationId, status):
//   1. Response codes the .proto file already declares via grpc-gateway
//      openapiv2 annotations (file-level defaults + per-RPC additions) —
//      authoritative, but invisible in the generated docs today because the
//      OpenAPI generator tool reads a different annotation vocabulary and
//      silently drops them. Parsed directly from the .proto source here.
//   2. Traced (id, file, line) error call sites per endpoint
//      (endpoint-error-tracing/<service>/*.json), cross-referenced against
//      the error catalog from generate-error-reference.ts for the why-text.
//
// Run: node apps/docs/scripts/generate-endpoint-errors.ts
// (standalone, like generate-error-reference.ts — not wired into the docs
// build's `generate` chain; re-run manually after adding more traced output)
//
// ── Extending this to a new service ─────────────────────────────────────
// 1. Add one entry to SERVICES below: the service's full proto package name
//    (matches the generated MDX filename prefix, e.g. `<service>.<Op>.mdx`),
//    its content directory, its .proto file, and a fresh subdirectory under
//    endpoint-error-tracing/ for that service's traced output. That's the
//    only code change needed — everything past this point already reads
//    generically from whatever's in SERVICES.
// 2. For each operation you want covered, trace its errors: open its gRPC/
//    Connect handler (internal/api/grpc/<service>/...), follow every call
//    into the command/query layer and shared helpers, and note every
//    zerrors.Throw*(...) site actually reachable from that specific
//    operation (not just every throw in the files it happens to touch).
//    Write the result as { "<OperationId>": { "handler": "file:line",
//    "errors": [{ "id", "file", "line", "reasoning" }] } } into that
//    service's tracing subdirectory (one file per operation or batched,
//    either works — this script merges every *.json file it finds there).
// 3. Re-run this script. It resolves each traced (id, file, line) against
//    the existing error catalog for the why-text and example response, so
//    no explanation needs to be hand-written twice.
// ─────────────────────────────────────────────────────────────────────────

import { readFileSync, writeFileSync, existsSync, readdirSync, mkdirSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import { GRPC_STATUS } from '../lib/grpc-status';

const __dirname = dirname(fileURLToPath(import.meta.url));
const DOCS_ROOT = join(__dirname, '..');
const REPO_ROOT = join(DOCS_ROOT, '../..');
const DATA_OUT = join(DOCS_ROOT, 'components/EndpointErrors/data.json');

const SERVICES = [
  {
    service: 'zitadel.user.v2.UserService',
    contentDir: join(DOCS_ROOT, 'content/reference/api/user'),
    protoFile: join(REPO_ROOT, 'proto/zitadel/user/v2/user_service.proto'),
    tracingDir: join(__dirname, 'endpoint-error-tracing/user-v2'),
  },
  // Add more services here — see "Extending this to a new service" above.
];

const STATUS_TEXT: Record<number, string> = Object.fromEntries(Object.values(GRPC_STATUS).map((s) => [s.httpStatus, s.httpText]));

interface ErrorCluster {
  key: string;
  message: string;
  why: string;
  ids: { id: string; locations: { file: string; line: number }[] }[];
  example: { httpStatus: number; grpcCode: number };
}
type ExampleResponse = {
  httpStatus: number;
  grpcCode: number;
  body: { code: number; message: string; details: { '@type': string; id: string; message: string }[] };
};
type Cause = { key: string; id: string; message: string; why: string; example: ExampleResponse };
type StatusGroup = { status: number; statusText: string; description?: string; causes: Cause[] };
type EndpointErrorsData = Record<string, Record<string, StatusGroup[]>>;

// Mirrors generate-error-reference.ts's example construction: the server
// always appends " (ID)" to the top-level message (see
// internal/api/grpc/gerrors/zitadel_errors.go), but uses the plain message
// with no suffix inside details[].message. Built per-cause with the exact ID
// this endpoint actually throws, not an arbitrary representative from the
// cluster, so the example matches what this specific endpoint returns.
function buildExample(cluster: ErrorCluster, id: string, message: string): ExampleResponse {
  return {
    httpStatus: cluster.example.httpStatus,
    grpcCode: cluster.example.grpcCode,
    body: {
      code: cluster.example.grpcCode,
      message: `${message} (${id})`,
      details: [{ '@type': 'type.googleapis.com/zitadel.v1.ErrorDetail', id, message }],
    },
  };
}

const clusters: ErrorCluster[] = JSON.parse(readFileSync(join(DOCS_ROOT, 'components/ErrorReference/data.json'), 'utf8')).subsystems.flatMap(
  (s: any) => s.kinds.flatMap((k: any) => k.clusters),
);

function findCluster(id: string, file: string, line: number): ErrorCluster | null {
  const norm = (f: string) => f.replace(/^\.?\/*/, '').trim();
  let fallback: ErrorCluster | null = null;
  for (const c of clusters) {
    for (const e of c.ids) {
      if (e.id !== id) continue;
      if (e.locations.some((l) => norm(l.file) === norm(file) && l.line === line)) return c;
      if (e.locations.some((l) => norm(l.file) === norm(file))) fallback = c;
    }
  }
  return fallback ?? clusters.find((c) => c.ids.some((e) => e.id === id)) ?? null;
}

// --- .proto openapiv2 annotation parsing -----------------------------------

function extractBalancedFrom(text: string, searchStart: number): string | null {
  const start = text.indexOf('{', searchStart);
  if (start === -1) return null;
  let depth = 0;
  for (let i = start; i < text.length; i++) {
    if (text[i] === '{') depth++;
    else if (text[i] === '}' && --depth === 0) return text.slice(start + 1, i);
  }
  return null;
}

function extractResponses(text: string): { status: number; description: string }[] {
  const out: { status: number; description: string }[] = [];
  let i = 0;
  while (true) {
    const idx = text.indexOf('responses:', i);
    if (idx === -1) break;
    const block = extractBalancedFrom(text, idx);
    const key = block?.match(/key:\s*"(\d+)"/);
    if (key) out.push({ status: +key[1], description: block!.match(/description:\s*"((?:[^"\\]|\\.)*)"/)?.[1].replace(/\\"/g, '"') ?? '' });
    i = idx + 10;
  }
  return out;
}

function parseDeclaredResponses(protoFile: string): Record<string, { status: number; description: string }[]> {
  if (!existsSync(protoFile)) return {};
  const text = readFileSync(protoFile, 'utf8');
  const serviceIdx = text.search(/\bservice\s+\w+\s*\{/);
  const defaults = extractResponses(extractBalancedFrom(text.slice(0, serviceIdx), text.indexOf('openapiv2_swagger) = ')) ?? '').filter(
    (r) => r.status >= 400,
  );

  const result: Record<string, { status: number; description: string }[]> = {};
  const rpcRe = /\brpc\s+(\w+)\s*\(/g;
  let m: RegExpExecArray | null;
  while ((m = rpcRe.exec(text))) {
    const body = extractBalancedFrom(text, m.index + m[0].length);
    if (body === null) continue;
    const specific = extractResponses(extractBalancedFrom(body, body.indexOf('openapiv2_operation) = ')) ?? '').filter((r) => r.status >= 400);
    const merged = new Map<number, { status: number; description: string }>();
    for (const r of [...defaults, ...specific]) merged.set(r.status, r);
    result[m[1]] = [...merged.values()].sort((a, b) => a.status - b.status);
  }
  return result;
}

// --- main --------------------------------------------------------------

function main() {
  const data: EndpointErrorsData = {};
  let opsUpdated = 0;
  let matched = 0;
  let unmatched = 0;

  for (const svc of SERVICES) {
    data[svc.service] = {};
    const declaredByOp = parseDeclaredResponses(svc.protoFile);
    // Scoped to this service's own tracing subdirectory — keeps two
    // services from silently colliding if they happen to both define an
    // operation with the same name.
    const traced: Record<string, { errors: { id: string; file: string; line: number }[] }> = {};
    if (existsSync(svc.tracingDir)) {
      for (const f of readdirSync(svc.tracingDir).filter((f) => f.endsWith('.json'))) {
        Object.assign(traced, JSON.parse(readFileSync(join(svc.tracingDir, f), 'utf8')));
      }
    }

    for (const [operationId, opData] of Object.entries(traced)) {
      const mdxPath = join(svc.contentDir, `${svc.service}.${operationId}.mdx`);
      if (!existsSync(mdxPath)) continue;

      const byStatus = new Map<number, StatusGroup>();
      for (const d of declaredByOp[operationId] ?? []) {
        byStatus.set(d.status, { status: d.status, statusText: STATUS_TEXT[d.status] ?? '', description: d.description, causes: [] });
      }
      for (const site of opData.errors) {
        const cluster = findCluster(site.id, site.file, site.line);
        if (!cluster) {
          unmatched++;
          continue;
        }
        matched++;
        const status = cluster.example.httpStatus;
        if (!byStatus.has(status)) byStatus.set(status, { status, statusText: STATUS_TEXT[status] ?? '', causes: [] });
        const group = byStatus.get(status)!;
        if (!group.causes.some((c) => c.key === cluster.key)) {
          const message = cluster.message.trim() || '(empty message)';
          group.causes.push({ key: cluster.key, id: site.id, message, why: cluster.why, example: buildExample(cluster, site.id, message) });
        }
      }

      const groups = [...byStatus.values()].sort((a, b) => a.status - b.status);
      for (const g of groups) g.causes.sort((a, b) => a.message.localeCompare(b.message));
      if (groups.length === 0) continue;

      data[svc.service][operationId] = groups;

      const original = readFileSync(mdxPath, 'utf8');
      const marker = '## Possible error responses';
      const base = original.includes(marker) ? original.slice(0, original.indexOf(marker)) : original;
      const stub = `## Possible error responses\n\n<EndpointErrors service="${svc.service}" operationId="${operationId}" />\n`;
      writeFileSync(mdxPath, base.replace(/\n*$/, '') + '\n\n' + stub);
      opsUpdated++;
    }
  }

  mkdirSync(dirname(DATA_OUT), { recursive: true });
  writeFileSync(DATA_OUT, JSON.stringify(data));
  console.log(`[endpoint-errors] ${opsUpdated} operation page(s) updated, ${matched} sites matched, ${unmatched} unmatched`);
}

main();
