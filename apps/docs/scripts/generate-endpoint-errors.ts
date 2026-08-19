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
// This *is* wired into the docs build's `generate` chain — the last step of
// generate-api-reference.ts's `--only-fix` branch calls it, so every normal
// docs build (and every Vercel deploy) re-applies whatever tracing JSON is
// committed under endpoint-error-tracing/, with no extra step. Deliberate:
// it's what makes "commit the tracing JSON, the pages update themselves" true.
// You only need to run this file directly when iterating locally without
// wanting to run the full docs generation pipeline.
//
// Unlike generate-error-reference.ts (which scans 2000+ backend Go files and
// is genuinely too slow to run on every build), this only reads the small
// set of already-traced JSON files plus the .proto files — cheap enough to
// run on every generate.
//
// ── The standard: what a v2 service needs for this to pick it up ────────
// discoverServices() below scans proto/zitadel/*/v2/*_service.proto and
// derives everything by reading source — nothing is hand-typed, so there is
// no service-name string anyone can get wrong (see discoverServices()'s own
// comment for why that matters). But that scan only works if a service
// follows the same layout every existing v2 service already follows. This is
// the exact, complete contract — not a summary of it:
//
// 1. Proto file at `proto/zitadel/<category>/v2/<anything>_service.proto`
//    (must end in `_service.proto` — that's how a real service is told apart
//    from a message-only file like metadata.proto, which has none). Must
//    contain `package zitadel.<category>.v2;` and `service <Name> { ... }`.
//    Each RPC declared as `rpc <OperationId>(...)`. This is already standard
//    ZITADEL proto style — nothing new to learn to satisfy it.
// 2. Go handlers in package `internal/api/grpc/<category>/v2` — any file(s),
//    the tracer scans the whole package. For each RPC you want covered,
//    there must be a method (any receiver, but `*Server` by convention)
//    named EXACTLY the RPC's name, e.g. `rpc CreateSession(...)` needs a
//    `func (s *Server) CreateSession(...)` somewhere in that package.
// 3. Docs content dir `apps/docs/content/reference/api/<category>/` must
//    exist with the generated per-operation MDX files already in it. This
//    comes from the normal docs generation pipeline once the proto is wired
//    into it the same way every other service is — not something to
//    hand-create.
// 4. The <category> path segment (the directory name right after
//    proto/zitadel/, internal/api/grpc/, and content/reference/api/) must be
//    IDENTICAL across all three. This is the one rule with no error message
//    if you get it wrong — a mismatch just makes the service quietly not
//    appear. Run `pnpm generate:endpoint-trace:category` and check the
//    "not fully wired up yet" section it prints (via
//    listUnwiredProtoServices() below) if a category you expect is missing;
//    it names exactly which of #2/#3 it can't find.
// 5. (Optional, cosmetic only) A `description: "..."` inside the proto's own
//    openapiv2_swagger info block, before the `service` line, shows in the
//    category picker. No description there just means the CLI shows none —
//    never invented.
//
// None of this is new process to adopt — it's already how every existing v2
// service in this repo is laid out. The contract only exists so tooling can
// rely on it holding, not to add extra steps to writing a new API.
//
// To actually trace a category's operations once it satisfies the above:
// 1. Run `go run ./internal/tools/errortrace --proto <service>.proto` (see
//    internal/tools/errortrace/main.go) to compute every zerrors.Throw*(...)
//    site reachable from each RPC's handler — mechanical, no manual code
//    reading required. It writes { "<OperationId>": { "handler": "file:line",
//    "errors": [{ "id", "file", "line", "reasoning" }] } }; drop that into
//    the service's tracing subdirectory (this script merges every *.json
//    file it finds there — one file per operation or batched, either works).
//    `pnpm generate:endpoint-trace:category` does this step for you, for a
//    whole category at once, picked from a menu.
// 2. Re-run this script (or let the CLI above do it). It resolves each
//    traced (id, file, line) against the existing error catalog for the
//    why-text and example response — if a newly-added error ID shows up as
//    "unmatched", the catalog is stale: re-run `pnpm generate:error-reference`
//    first, since that's what scans internal/**/*.go for every
//    zerrors.Throw*() site in the first place.
// ─────────────────────────────────────────────────────────────────────────

import { readFileSync, writeFileSync, existsSync, readdirSync, mkdirSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import { GRPC_STATUS } from '../lib/grpc-status';

const __dirname = dirname(fileURLToPath(import.meta.url));
const DOCS_ROOT = join(__dirname, '..');
const REPO_ROOT = join(DOCS_ROOT, '../..');
const DATA_OUT = join(DOCS_ROOT, 'components/EndpointErrors/data.json');
const PROTO_ZITADEL_ROOT = join(REPO_ROOT, 'proto/zitadel');

export interface ServiceConfig {
  service: string; // full proto package + service name, e.g. 'zitadel.org.v2.OrganizationService'
  category: string; // proto/docs directory segment, e.g. 'org'
  description?: string; // from the proto's own openapiv2_swagger info block, if it has one — never invented
  rpcNames: string[]; // declaration order, straight from the .proto file
  contentDir: string;
  protoFile: string;
  tracingDir: string;
  goPackageDir: string;
}

interface Candidate extends ServiceConfig {
  contentDirExists: boolean;
  goPackageDirExists: boolean;
}

// Reads every proto/zitadel/*/v2/*_service.proto file and derives everything
// discoverServices()/listUnwiredProtoServices() need, regardless of whether
// its docs/Go side is wired up yet — the two public functions below just
// filter this differently. Kept as one shared scan so "is this category
// wired up" and "what's it missing" can never disagree with each other.
function scanCandidates(): Candidate[] {
  if (!existsSync(PROTO_ZITADEL_ROOT)) return [];
  const categories = readdirSync(PROTO_ZITADEL_ROOT, { withFileTypes: true })
    .filter((e) => e.isDirectory())
    .map((e) => e.name)
    .sort();

  const candidates: Candidate[] = [];
  for (const category of categories) {
    const v2Dir = join(PROTO_ZITADEL_ROOT, category, 'v2');
    if (!existsSync(v2Dir)) continue;
    const serviceFileName = readdirSync(v2Dir).find((f) => f.endsWith('_service.proto'));
    if (!serviceFileName) continue; // message-only dirs (metadata, filter, object, error, ...) — not a service

    const protoFile = join(v2Dir, serviceFileName);
    const text = readFileSync(protoFile, 'utf8');

    const serviceIdx = text.search(/\bservice\s+\w+\s*\{/);
    if (serviceIdx === -1) continue;
    const serviceName = text.slice(serviceIdx).match(/\bservice\s+(\w+)\s*\{/)?.[1];
    if (!serviceName) continue;

    const protoPackage = text.match(/\bpackage\s+([\w.]+)\s*;/)?.[1] ?? `zitadel.${category}.v2`;
    // Only search before the `service` keyword — same region parseDeclaredResponses
    // reads the swagger defaults from — so a per-RPC/per-field description
    // declared later in the file is never mistaken for the service's own.
    const description = text
      .slice(0, serviceIdx)
      .match(/\bdescription:\s*"((?:[^"\\]|\\.)*)"/)?.[1]
      .replace(/\\"/g, '"')
      .replace(/\\n/g, ' ')
      .trim();
    const rpcNames = [...text.matchAll(/\brpc\s+(\w+)\s*\(/g)].map((m) => m[1]);

    const contentDir = join(DOCS_ROOT, 'content/reference/api', category);
    const goPackageDir = join(REPO_ROOT, 'internal/api/grpc', category, 'v2');

    candidates.push({
      service: `${protoPackage}.${serviceName}`,
      category,
      description,
      rpcNames,
      contentDir,
      protoFile,
      tracingDir: join(__dirname, 'endpoint-error-tracing', `${category}-v2`),
      goPackageDir,
      contentDirExists: existsSync(contentDir),
      goPackageDirExists: existsSync(goPackageDir),
    });
  }
  return candidates;
}

// Discovers every v2 API service that's fully wired up (proto + docs content
// dir + Go handler package all exist), by reading source directly instead of
// a hand-typed list. This exists because a hand-typed SERVICES array is
// exactly how a real bug happened: an entry once read
// `service: 'zitadel.org.v2.OrgService'` — a guess — when the .proto file
// actually declares `service OrganizationService`. Reading the name from the
// source it's declared in makes that whole class of mistake impossible,
// rather than just "be more careful next time."
export function discoverServices(): ServiceConfig[] {
  return scanCandidates().filter((c) => c.contentDirExists && c.goPackageDirExists);
}

// The other half of discoverServices(): proto services that exist but aren't
// fully wired up yet, with exactly what's missing — rule #4 in the "standard"
// comment above is the one with no error message if violated (a <category>
// segment mismatch just makes a service silently not appear), so this makes
// that failure mode visible instead of silent. trace-category-cli.ts prints
// this below the main menu.
export function listUnwiredProtoServices(): { category: string; service: string; protoFile: string; missing: string[] }[] {
  return scanCandidates()
    .filter((c) => !c.contentDirExists || !c.goPackageDirExists)
    .map((c) => ({
      category: c.category,
      service: c.service,
      protoFile: c.protoFile,
      missing: [
        !c.contentDirExists && `docs content dir (${relativeToRepo(c.contentDir)})`,
        !c.goPackageDirExists && `Go handler package (${relativeToRepo(c.goPackageDir)})`,
      ].filter((x): x is string => Boolean(x)),
    }));
}

function relativeToRepo(absPath: string): string {
  return absPath.startsWith(REPO_ROOT + '/') ? absPath.slice(REPO_ROOT.length + 1) : absPath;
}

const SERVICES = discoverServices();

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

// A "GRPC-<CODE>" id (from internal/tools/errortrace's recordRawStatus) is a
// raw status.Errorf(codes.X, ...) call, not a zerrors.Throw* — it was never
// scanned into the error catalog, so there's no cluster to look up. Build
// one on the fly instead of reporting it as unmatched: the code name alone
// is enough to know the HTTP/gRPC status, and the traced site carries its
// own literal message (see errorSite.Message in main.go) since nothing else
// has it.
function syntheticGrpcCluster(id: string, message: string): ErrorCluster | null {
  if (!id.startsWith('GRPC-')) return null;
  const status = GRPC_STATUS[id.slice('GRPC-'.length)];
  if (!status) return null;
  return {
    key: `synthetic|${id}`,
    message,
    why: `This operation returns this status directly from its handler code, rather than through the usual error-catalog convention — so this explanation is mechanical, not researched.`,
    ids: [{ id, locations: [] }],
    example: { httpStatus: status.httpStatus, grpcCode: status.code },
  };
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
    const traced: Record<string, { errors: { id: string; file: string; line: number; message?: string }[] }> = {};
    if (existsSync(svc.tracingDir)) {
      for (const f of readdirSync(svc.tracingDir).filter((f) => f.endsWith('.json'))) {
        Object.assign(traced, JSON.parse(readFileSync(join(svc.tracingDir, f), 'utf8')));
      }
    }

    // Union declaredByOp into the operation list, but only for a service
    // that's actually been traced at all (traced has >=1 entry) — within an
    // already-traced service, an operation with proto-declared responses but
    // nothing in `traced` must still show its declared defaults rather than
    // being silently skipped. For a service nobody has traced yet, stay
    // fully empty rather than surface generic proto-declared placeholders
    // for every operation across every untouched service.
    const operationIds =
      Object.keys(traced).length > 0 ? new Set([...Object.keys(traced), ...Object.keys(declaredByOp)]) : new Set(Object.keys(traced));
    for (const operationId of operationIds) {
      const mdxPath = join(svc.contentDir, `${svc.service}.${operationId}.mdx`);
      if (!existsSync(mdxPath)) continue;

      const byStatus = new Map<number, StatusGroup>();
      for (const d of declaredByOp[operationId] ?? []) {
        byStatus.set(d.status, { status: d.status, statusText: STATUS_TEXT[d.status] ?? '', description: d.description, causes: [] });
      }
      for (const site of traced[operationId]?.errors ?? []) {
        // A real, hand-assigned zerrors ID can coincidentally start with
        // "GRPC-" too (confirmed: GRPC-vR9nC is a real catalog entry, not
        // one of ours) — so the real catalog always gets first look, and
        // the synthetic path is only a fallback for what it can't explain.
        const cluster = findCluster(site.id, site.file, site.line) ?? syntheticGrpcCluster(site.id, site.message ?? '');
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

      // A status declared generically at the proto/file level (not
      // substantiated by any traced cause) is a reasonable placeholder when
      // this operation was never traced — we genuinely don't know yet. But
      // if it *was* traced and the tracer confidently found zero reachable
      // causes for that status, keeping the declared-only group is actively
      // misleading, not just incomplete: e.g. SetOrganizationFeatures is a
      // one-line `status.Errorf(codes.Unimplemented, ...)` stub that never
      // touches zerrors at all, so it can only ever return UNIMPLEMENTED —
      // yet the proto's generic file-level 403/404 defaults would otherwise
      // show up here as if they were real possibilities for this endpoint.
      const wasTraced = operationId in traced;
      const byStatusFiltered = wasTraced ? new Map([...byStatus].filter(([, g]) => g.causes.length > 0)) : byStatus;
      const groups = [...byStatusFiltered.values()].sort((a, b) => a.status - b.status);
      for (const g of groups) g.causes.sort((a, b) => a.message.localeCompare(b.message));
      const marker = '## Possible error responses';
      const original = readFileSync(mdxPath, 'utf8');
      const base = original.includes(marker) ? original.slice(0, original.indexOf(marker)) : original;

      if (groups.length === 0 && !wasTraced) {
        // Never traced at all — stay fully silent rather than expose
        // internal tooling state ("not traced yet") on a public docs page.
        // Still strip a stale section from a previous run if one exists:
        // original only differs from base when the marker was actually
        // present, so this is a no-op write otherwise.
        if (original !== base) writeFileSync(mdxPath, base.replace(/\n*$/, '') + '\n');
        continue;
      }

      // Two ways to land here with an empty `groups`: either it genuinely
      // has causes, or it was traced and confidently found none (an empty
      // array here, as opposed to no entry at all, is what tells
      // EndpointErrors this is a real, checked result — not staleness).
      // Rendering a note either way is what actually resolves the
      // SetOrganizationFeatures confusion: before this, "no table" looked
      // identical whether an operation had been traced or not.
      data[svc.service][operationId] = groups;
      const stub = `## Possible error responses\n\n<EndpointErrors service="${svc.service}" operationId="${operationId}" />\n`;
      writeFileSync(mdxPath, base.replace(/\n*$/, '') + '\n\n' + stub);
      opsUpdated++;
    }
  }

  mkdirSync(dirname(DATA_OUT), { recursive: true });
  writeFileSync(DATA_OUT, JSON.stringify(data));
  console.log(`[endpoint-errors] ${opsUpdated} operation page(s) updated, ${matched} sites matched, ${unmatched} unmatched`);
}

// Only run the merge when this file is executed directly (`tsx
// generate-endpoint-errors.ts`) — not when another script imports
// discoverServices() from it (trace-endpoint-errors.ts, trace-category-cli.ts),
// which would otherwise silently trigger a full merge run as a side effect
// of just wanting the service list.
if (import.meta.url === `file://${process.argv[1]}`) main();
