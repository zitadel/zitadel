// Scans the Go backend for zerrors.Throw<Kind>[f]("ID", "message") call
// sites and emits apps/docs/components/ErrorReference/data.json: a
// 3-level tree (subsystem -> gRPC kind -> cluster) that the ErrorReference
// component renders. Regenerate after backend error call sites change:
//   node apps/docs/scripts/generate-error-reference.ts
//
// This is intentionally NOT wired into the Nx `generate` chain (which
// `dev`/`build`/`lint` all depend on) — scanning 2000+ backend Go files on
// every docs build would cache-invalidate the whole site on unrelated
// backend PRs. Output (data.json + overrides) is committed instead, the
// same way apps/docs/content/apis/benchmarks/**/output.json is committed
// despite being externally generated.

import { readFileSync, writeFileSync, mkdirSync, readdirSync, existsSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import { execSync } from 'child_process';
import { globSync } from 'glob';
import yaml from 'js-yaml';
import { GRPC_STATUS, THROW_KIND_TO_GRPC_STATUS } from '../lib/grpc-status';

const __dirname = dirname(fileURLToPath(import.meta.url));
const REPO_ROOT = join(__dirname, '../../..');
const OUTPUT_PATH = join(__dirname, '../components/ErrorReference/data.json');
const OVERRIDES_DIR = join(__dirname, 'error-reference-overrides');

// ---------------------------------------------------------------------------
// Types (mirrored, informally, by apps/docs/components/ErrorReference/index.tsx)
// ---------------------------------------------------------------------------

interface RawEntry {
  id: string;
  kind: string; // GRPC_STATUS key
  message: string;
  messageType: 'i18n' | 'literal';
  i18nKey?: string;
  i18nResolved?: boolean;
  file: string;
  line: number;
}

interface ErrorLocation {
  file: string;
  line: number;
}

interface ErrorIdEntry {
  id: string;
  locations: ErrorLocation[];
  i18nKey?: string;
}

interface ExampleResponse {
  httpStatus: number;
  grpcCode: number;
  body: {
    code: number;
    message: string;
    details: [{ '@type': string; id: string; message: string }];
  };
}

interface ErrorCluster {
  key: string;
  message: string;
  messageType: 'i18n' | 'literal';
  ids: ErrorIdEntry[];
  why: string;
  whySource: 'curated' | 'i18n-decomposition' | 'defensive-guard-template' | 'generic-fallback';
  example: ExampleResponse;
  collidesWith: { key: string; message: string; subsystem: string }[];
}

interface KindBucket {
  kind: string;
  clusters: ErrorCluster[];
}

interface Subsystem {
  id: string;
  name: string;
  summary: string;
  kinds: KindBucket[];
}

interface ErrorReferenceData {
  meta: {
    generatedAt: string;
    sourceRef: string;
    totalCallSites: number;
    totalIds: number;
    totalClusters: number;
    totalCollisions: number;
  };
  subsystems: Subsystem[];
}

// ---------------------------------------------------------------------------
// Stage: subsystem mapping (directory prefix -> subsystem id)
// ---------------------------------------------------------------------------

const SUBSYSTEMS: { id: string; name: string; summary: string }[] = [
  {
    id: 'command',
    name: 'Command layer',
    summary:
      'internal/command is the write side of ZITADEL’s CQRS/event-sourcing core: nearly every state-changing gRPC/Connect call resolves to a Commands method covering users, orgs, instances, projects, IDPs, policies, providers, quotas, actions/executions, sessions, and groups. Errors cluster into three layers mapped to distinct gRPC codes: InvalidArgument for cheap syntactic validation before any eventstore read, NotFound/AlreadyExists for aggregate-existence checks against replayed write-model state, and PreconditionFailed for state-conditional rules discovered mid-command (not changed, already active, lockouts).',
  },
  {
    id: 'query',
    name: 'Query layer & legacy read models',
    summary:
      'internal/query is the current read side (Postgres projections), alongside the legacy internal/auth/internal/user/internal/view read-model stack it is gradually replacing. Errors here are dominated by NotFound (a projection row doesn’t exist or hasn’t caught up yet) and Internal (SQL statement construction / scan failures).',
  },
  {
    id: 'api',
    name: 'API surface, authn/authz, and identity federation',
    summary:
      'internal/api (gRPC/Connect handlers, REST gateway, login UI), internal/authz (token/permission checks), internal/idp (external identity provider integrations), and internal/webauthn sit at the transport/identity boundary. Errors here span request-shape validation, auth/session failures, and third-party IdP response handling.',
  },
  {
    id: 'storage',
    name: 'Storage/eventstore infrastructure & the new v3 backend',
    summary:
      'internal/eventstore (event append/query, Postgres-backed), internal/static (asset storage), and backend/v3 (the newer hexagonal-architecture, relational-storage-backed rewrite) — infrastructure errors: connection failures, constraint violations, unmarshal guards on stored event payloads.',
  },
  {
    id: 'scim',
    name: 'SCIM provisioning',
    summary:
      'internal/api/scim implements SCIM 2.0 user/group provisioning. Errors mirror the SCIM spec’s own error model (invalid filters, invalid paths, bulk-operation limits) more than ZITADEL’s general gRPC error shape.',
  },
  {
    id: 'notification',
    name: 'Notification channels',
    summary:
      'internal/notification renders and sends emails/SMS (templates, SMTP/Twilio/webhook providers). Errors are mostly provider-configuration and template-rendering failures.',
  },
  {
    id: 'crypto',
    name: 'Cryptography & key management',
    summary:
      'internal/crypto covers encryption, hashing, and key/secret generation used throughout the rest of the backend. Errors here are almost entirely configuration/algorithm mismatches, not runtime request data.',
  },
  {
    id: 'domain',
    name: 'Domain value objects, legacy event payloads, and misc infra',
    summary:
      'internal/domain (value-object validation), internal/repository (event definitions/unmarshalling for most aggregates), internal/execution, and other infra that doesn’t fit the groups above. A large share of this group is Internal-coded defensive guards inside event-unmarshalling code that only fire on event-schema corruption, not from ordinary API usage.',
  },
];

function subsystemForPath(relPath: string): string {
  if (relPath.startsWith('internal/api/scim/')) return 'scim';
  if (relPath.startsWith('internal/notification/')) return 'notification';
  if (relPath.startsWith('internal/crypto/')) return 'crypto';
  if (
    relPath.startsWith('internal/eventstore/') ||
    relPath.startsWith('internal/static/') ||
    relPath.startsWith('backend/v3/')
  )
    return 'storage';
  if (
    relPath.startsWith('internal/api/') ||
    relPath.startsWith('internal/authz/') ||
    relPath.startsWith('internal/idp/') ||
    relPath.startsWith('internal/webauthn/')
  )
    return 'api';
  if (
    relPath.startsWith('internal/query/') ||
    relPath.startsWith('internal/auth/') ||
    relPath.startsWith('internal/user/') ||
    relPath.startsWith('internal/view/') ||
    relPath.includes('/repository/view/') // legacy view-model repos for org/project/iam/user aggregates
  )
    return 'query';
  if (relPath.startsWith('internal/command/')) return 'command';
  // Event-payload unmarshal guards for these three aggregates are documented
  // (ERRORS.md) as conceptually part of the command layer, sharing its ID
  // prefixes (ORG-/INSTANCE-/PROJECT-), unlike other internal/repository/*.
  if (
    relPath.startsWith('internal/repository/org/') ||
    relPath.startsWith('internal/repository/instance/') ||
    relPath.startsWith('internal/repository/project/')
  )
    return 'command';
  return 'domain';
}

// ---------------------------------------------------------------------------
// Stage: scan .go files for zerrors.Throw*() call sites
// ---------------------------------------------------------------------------

const KIND_ALIASES: Record<string, keyof typeof GRPC_STATUS> = {
  ...THROW_KIND_TO_GRPC_STATUS,
  Error: 'UNKNOWN', // ThrowError always creates KindUnknown (see internal/zerrors/zerror.go)
};

const ID_SHAPE = /^[A-Za-z][A-Za-z0-9]*-[A-Za-z0-9]+$/;
const CALL_RE = /zerrors\.Throw([A-Z][A-Za-z]*?)(f)?\(/g;

// Reads up to 2 top-level (paren-depth-1) double-quoted string literals
// starting right after a Throw*( call's opening paren.
function extractTopLevelStrings(content: string, startIdx: number): string[] {
  const strings: string[] = [];
  let depth = 1;
  let i = startIdx;
  const limit = Math.min(content.length, startIdx + 4000);
  while (i < limit && depth > 0 && strings.length < 2) {
    const ch = content[i];
    if (ch === '"' && depth === 1) {
      let j = i + 1;
      let raw = '';
      while (j < content.length && content[j] !== '"') {
        if (content[j] === '\\') {
          raw += content[j] + content[j + 1];
          j += 2;
        } else {
          raw += content[j];
          j++;
        }
      }
      try {
        strings.push(JSON.parse('"' + raw + '"'));
      } catch {
        strings.push(raw);
      }
      i = j + 1;
      continue;
    }
    if (ch === '(') depth++;
    else if (ch === ')') depth--;
    i++;
  }
  return strings;
}

function lineOf(content: string, idx: number): number {
  let line = 1;
  for (let i = 0; i < idx; i++) if (content[i] === '\n') line++;
  return line;
}

function scanFile(absPath: string, relPath: string): RawEntry[] {
  const content = readFileSync(absPath, 'utf8');
  const entries: RawEntry[] = [];
  let m: RegExpExecArray | null;
  CALL_RE.lastIndex = 0;
  while ((m = CALL_RE.exec(content))) {
    const kindRaw = m[1];
    const grpcKey = KIND_ALIASES[kindRaw];
    if (!grpcKey) continue; // not a real Throw<Kind> (defensive)
    const strings = extractTopLevelStrings(content, CALL_RE.lastIndex);
    if (strings.length < 2 || !ID_SHAPE.test(strings[0])) continue;
    const [id, message] = strings;
    const messageType: 'i18n' | 'literal' = /^[A-Za-z][A-Za-z0-9]*(\.[A-Za-z][A-Za-z0-9]*)+$/.test(message)
      ? 'i18n'
      : 'literal';
    entries.push({
      id,
      kind: grpcKey,
      message,
      messageType,
      i18nKey: messageType === 'i18n' ? message : undefined,
      file: relPath,
      line: lineOf(content, m.index),
    });
  }
  return entries;
}

// ---------------------------------------------------------------------------
// Stage: i18n resolution
// ---------------------------------------------------------------------------

function loadYaml(relPath: string): any {
  try {
    return yaml.load(readFileSync(join(REPO_ROOT, relPath), 'utf8'));
  } catch {
    return {};
  }
}

const I18N_BUNDLES = [
  loadYaml('internal/static/i18n/en.yaml'),
  loadYaml('internal/api/ui/login/static/i18n/en.yaml'),
  loadYaml('internal/notification/static/i18n/en.yaml'),
];

function resolveI18n(key: string): string | undefined {
  for (const bundle of I18N_BUNDLES) {
    let node: any = bundle;
    for (const segment of key.split('.')) {
      if (node == null || typeof node !== 'object') {
        node = undefined;
        break;
      }
      node = node[segment];
    }
    if (typeof node === 'string') return node;
  }
  return undefined;
}

// ---------------------------------------------------------------------------
// Stage: "why" generation (mechanical baseline — curated overrides applied later)
// ---------------------------------------------------------------------------

const GUARD_PATTERNS = [/^unable to unmarshal /i, /^reduce\.wrong\.event\.type/i, /^cannot scan type /i, /^could not unmarshal /i];

const KIND_WHY_VERB: Record<string, (subject: string) => string> = {
  NOT_FOUND: (s) => `${s} could not be found.`,
  ALREADY_EXISTS: (s) => `${s} already exists.`,
  INVALID_ARGUMENT: (s) => `The request was rejected because ${s.toLowerCase()} is missing or invalid.`,
  FAILED_PRECONDITION: (s) => `The request was rejected because of ${s.toLowerCase()} — the current state doesn’t allow this operation.`,
  PERMISSION_DENIED: (s) => `The caller is authenticated but not authorized: ${s.toLowerCase()}.`,
  UNAUTHENTICATED: (s) => `The caller could not be authenticated: ${s.toLowerCase()}.`,
  RESOURCE_EXHAUSTED: (s) => `A limit was exceeded: ${s.toLowerCase()}.`,
  UNIMPLEMENTED: (s) => `${s} is not supported on this code path.`,
  INTERNAL: (s) => `Internal failure related to ${s.toLowerCase()}.`,
};

function humanizeI18nKey(key: string): string {
  const segments = key.split('.').filter((s) => s.toLowerCase() !== 'errors');
  const words = segments
    .map((seg) => seg.replace(/([a-z0-9])([A-Z])/g, '$1 $2'))
    .join(' ')
    .trim();
  return words.length > 0 ? `The ${words.toLowerCase()}` : 'This condition';
}

function mechanicalWhy(kind: string, message: string, messageType: 'i18n' | 'literal'): {
  why: string;
  whySource: ErrorCluster['whySource'];
} {
  if (GUARD_PATTERNS.some((re) => re.test(message))) {
    return {
      why: 'Internal defensive guard around event/projection unmarshalling — not expected to occur from ordinary API usage. Indicates event-schema corruption or a routing bug if it does.',
      whySource: 'defensive-guard-template',
    };
  }
  if (messageType === 'i18n') {
    const subject = humanizeI18nKey(message);
    const verb = KIND_WHY_VERB[kind];
    return {
      why: verb ? verb(subject) : `${subject} (${kind}).`,
      whySource: 'i18n-decomposition',
    };
  }
  return {
    why: `${kind} error: “${message}” — see the linked source location(s) for the exact trigger condition.`,
    whySource: 'generic-fallback',
  };
}

// ---------------------------------------------------------------------------
// Main pipeline
// ---------------------------------------------------------------------------

function main() {
  const files = globSync(['internal/**/*.go', 'backend/**/*.go'], {
    cwd: REPO_ROOT,
    ignore: ['**/*_test.go'],
  });

  const rawEntries: RawEntry[] = [];
  for (const relPath of files) {
    const abs = join(REPO_ROOT, relPath);
    rawEntries.push(...scanFile(abs, relPath));
  }

  // Resolve i18n text used for clustering/display, without discarding the raw key.
  for (const e of rawEntries) {
    if (e.messageType === 'i18n') {
      const resolved = resolveI18n(e.message);
      e.i18nResolved = resolved !== undefined;
      if (resolved !== undefined) e.message = resolved;
    }
  }

  // Cluster by (subsystem, kind, trimmed message).
  const clusterMap = new Map<
    string,
    { subsystem: string; kind: string; message: string; messageType: 'i18n' | 'literal'; entries: RawEntry[] }
  >();
  for (const e of rawEntries) {
    const subsystem = subsystemForPath(e.file);
    const message = e.message.trim();
    const key = `${subsystem}|${e.kind}|${message}`;
    if (!clusterMap.has(key)) {
      clusterMap.set(key, { subsystem, kind: e.kind, message, messageType: e.messageType, entries: [] });
    }
    clusterMap.get(key)!.entries.push(e);
  }

  // Collision detection: which cluster keys does each raw id appear in?
  const idToClusterKeys = new Map<string, Set<string>>();
  for (const [key, c] of clusterMap) {
    for (const e of c.entries) {
      if (!idToClusterKeys.has(e.id)) idToClusterKeys.set(e.id, new Set());
      idToClusterKeys.get(e.id)!.add(key);
    }
  }

  // Load curated overrides (clusterKey -> { why }).
  const overrides = new Map<string, { why: string }>();
  if (existsSync(OVERRIDES_DIR)) {
    for (const f of readdirSync(OVERRIDES_DIR).filter((f: string) => f.endsWith('.json'))) {
      const parsed = JSON.parse(readFileSync(join(OVERRIDES_DIR, f), 'utf8'));
      for (const [key, val] of Object.entries<{ why: string }>(parsed)) overrides.set(key, val);
    }
  }

  const subsystems: Subsystem[] = SUBSYSTEMS.map((s) => ({ ...s, kinds: [] }));
  const subsystemIndex = new Map(subsystems.map((s) => [s.id, s]));
  let totalClusters = 0;
  let totalCollisions = 0;

  for (const [key, c] of clusterMap) {
    totalClusters++;
    const idMap = new Map<string, ErrorIdEntry>();
    for (const e of c.entries) {
      if (!idMap.has(e.id)) idMap.set(e.id, { id: e.id, locations: [], i18nKey: e.i18nResolved === false ? e.i18nKey : undefined });
      idMap.get(e.id)!.locations.push({ file: e.file, line: e.line });
    }
    const ids = [...idMap.values()].sort((a, b) => a.id.localeCompare(b.id));

    // A cluster key already bakes in (subsystem, kind, message), so the same
    // id turning up under a different key is only a *real* collision — i.e.
    // the id means something genuinely different depending on call site — if
    // the message actually differs. Without this check, an id copy-pasted
    // verbatim into files that happen to sit in two different subsystems
    // would be flagged as colliding with itself.
    const collidesWith: ErrorCluster['collidesWith'] = [];
    const seenCollisionKeys = new Set<string>();
    for (const idEntry of ids) {
      for (const otherKey of idToClusterKeys.get(idEntry.id)!) {
        if (otherKey === key || seenCollisionKeys.has(otherKey)) continue;
        const other = clusterMap.get(otherKey)!;
        if (other.message === c.message) continue;
        seenCollisionKeys.add(otherKey);
        collidesWith.push({ key: otherKey, message: other.message, subsystem: other.subsystem });
      }
    }
    if (collidesWith.length > 0) totalCollisions++;

    const override = overrides.get(key);
    const mech = mechanicalWhy(c.kind, c.message, c.messageType);
    const why = override?.why ?? mech.why;
    const whySource: ErrorCluster['whySource'] = override ? 'curated' : mech.whySource;

    const status = GRPC_STATUS[c.kind];
    const repId = ids[0].id;
    const example: ExampleResponse = {
      httpStatus: status.httpStatus,
      grpcCode: status.code,
      body: {
        code: status.code,
        message: c.message,
        details: [{ '@type': 'type.googleapis.com/zitadel.v1.ErrorDetail', id: repId, message: c.message }],
      },
    };

    const cluster: ErrorCluster = {
      key,
      message: c.message,
      messageType: c.messageType,
      ids,
      why,
      whySource,
      example,
      collidesWith,
    };

    const subsystem = subsystemIndex.get(c.subsystem)!;
    let bucket = subsystem.kinds.find((k) => k.kind === c.kind);
    if (!bucket) {
      bucket = { kind: c.kind, clusters: [] };
      subsystem.kinds.push(bucket);
    }
    bucket.clusters.push(cluster);
  }

  // Sort clusters within each kind bucket by size descending (most common first),
  // then sort kind buckets within each subsystem by total id count descending.
  for (const s of subsystems) {
    for (const bucket of s.kinds) {
      bucket.clusters.sort((a, b) => b.ids.length - a.ids.length || a.message.localeCompare(b.message));
    }
    s.kinds.sort((a, b) => {
      const aCount = a.clusters.reduce((n, c) => n + c.ids.length, 0);
      const bCount = b.clusters.reduce((n, c) => n + c.ids.length, 0);
      return bCount - aCount;
    });
  }
  subsystems.sort((a, b) => {
    const count = (s: Subsystem) => s.kinds.reduce((n, k) => n + k.clusters.reduce((m, c) => m + c.ids.length, 0), 0);
    return count(b) - count(a);
  });

  let sourceRef = 'main';
  try {
    sourceRef = execSync('git rev-parse HEAD', { cwd: REPO_ROOT }).toString().trim();
  } catch {
    // best-effort only
  }

  const data: ErrorReferenceData = {
    meta: {
      generatedAt: new Date().toISOString(),
      sourceRef,
      totalCallSites: rawEntries.length,
      totalIds: idToClusterKeys.size,
      totalClusters,
      totalCollisions,
    },
    subsystems,
  };

  mkdirSync(dirname(OUTPUT_PATH), { recursive: true });
  writeFileSync(OUTPUT_PATH, JSON.stringify(data));

  console.log(`[error-reference] ${files.length} files scanned`);
  console.log(`[error-reference] ${rawEntries.length} call sites, ${idToClusterKeys.size} distinct IDs`);
  console.log(`[error-reference] ${totalClusters} clusters, ${totalCollisions} with a real collision`);
  console.log(`[error-reference] ${overrides.size} curated overrides applied`);
  console.log(`[error-reference] wrote ${OUTPUT_PATH}`);
}

main();
