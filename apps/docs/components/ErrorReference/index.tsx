'use client';

import { useEffect, useMemo, useState } from 'react';
import { ChevronRight, TriangleAlert } from 'lucide-react';
import { cn } from '@/utils/cn';

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
  body: { code: number; message: string; details: [{ '@type': string; id: string; message: string }] };
}
type WhySource = 'curated' | 'i18n-decomposition' | 'defensive-guard-template' | 'generic-fallback';
interface ErrorCluster {
  key: string;
  message: string;
  messageType: 'i18n' | 'literal';
  ids: ErrorIdEntry[];
  why: string;
  whySource: WhySource;
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

const WHY_SOURCE_LABEL: Record<WhySource, string> = {
  curated: 'Researched from source',
  'i18n-decomposition': 'Derived from message key',
  'defensive-guard-template': 'Internal guard',
  'generic-fallback': 'Auto-generated',
};

const WHY_SOURCE_STYLE: Record<WhySource, string> = {
  curated: 'bg-fd-accent text-fd-accent-foreground',
  'i18n-decomposition': 'border border-fd-border text-fd-muted-foreground',
  'defensive-guard-template': 'border border-fd-border text-fd-muted-foreground',
  'generic-fallback': 'border border-dashed border-fd-border text-fd-muted-foreground',
};

const MAX_LOCATIONS_SHOWN = 5;
const MAX_SEARCH_RESULTS = 150;

function clusterIdCount(cluster: ErrorCluster) {
  return cluster.ids.length;
}
function kindIdCount(kind: KindBucket) {
  return kind.clusters.reduce((n, c) => n + clusterIdCount(c), 0);
}
function subsystemIdCount(s: Subsystem) {
  return s.kinds.reduce((n, k) => n + kindIdCount(k), 0);
}

function matchesQuery(cluster: ErrorCluster, q: string) {
  if (cluster.message.toLowerCase().includes(q)) return true;
  if (cluster.why.toLowerCase().includes(q)) return true;
  return cluster.ids.some((e) => e.id.toLowerCase().includes(q));
}

export default function ErrorReference() {
  const [data, setData] = useState<ErrorReferenceData | null>(null);
  const [loadError, setLoadError] = useState(false);
  const [query, setQuery] = useState('');
  const [expanded, setExpanded] = useState<Set<string>>(new Set());

  useEffect(() => {
    let cancelled = false;
    fetch('/api/error-reference')
      .then((res) => {
        if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
        return res.json();
      })
      .then((json) => {
        if (!cancelled) setData(json as ErrorReferenceData);
      })
      .catch(() => {
        if (!cancelled) setLoadError(true);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  const trimmedQuery = query.trim().toLowerCase();

  const searchResults = useMemo(() => {
    if (!data || !trimmedQuery) return null;
    const results: { subsystem: Subsystem; kind: KindBucket; cluster: ErrorCluster }[] = [];
    for (const subsystem of data.subsystems) {
      for (const kind of subsystem.kinds) {
        for (const cluster of kind.clusters) {
          if (matchesQuery(cluster, trimmedQuery)) {
            results.push({ subsystem, kind, cluster });
            if (results.length >= MAX_SEARCH_RESULTS) return results;
          }
        }
      }
    }
    return results;
  }, [data, trimmedQuery]);

  function toggle(key: string) {
    setExpanded((prev) => {
      const next = new Set(prev);
      if (next.has(key)) next.delete(key);
      else next.add(key);
      return next;
    });
  }

  if (loadError) {
    return <p className="text-sm text-fd-muted-foreground">Couldn&apos;t load the error reference. Try reloading the page.</p>;
  }
  if (!data) {
    return <p className="text-sm text-fd-muted-foreground">Loading error reference…</p>;
  }

  return (
    <div className="not-prose flex flex-col gap-4">
      <div className="flex flex-col gap-2">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder={`Search ${data.meta.totalIds.toLocaleString()} error IDs or messages, e.g. "COMMAND-2M0fs" or "not found"…`}
          className="w-full px-3 py-2 text-sm bg-transparent border rounded-md border-fd-border text-fd-foreground focus:outline-none focus:ring-2 focus:ring-fd-accent"
        />
        <p className="text-xs text-fd-muted-foreground">
          {data.meta.totalIds.toLocaleString()} error IDs, grouped into {data.meta.totalClusters.toLocaleString()}{' '}
          distinct causes across {data.subsystems.length} subsystems. {data.meta.totalCollisions} of those causes
          share an ID string with at least one other cause — see the warning on those entries.
        </p>
      </div>

      {searchResults ? (
        <SearchResults results={searchResults} truncated={searchResults.length >= MAX_SEARCH_RESULTS} />
      ) : (
        <div className="flex flex-col gap-2">
          {data.subsystems.map((subsystem) => (
            <SubsystemSection
              key={subsystem.id}
              subsystem={subsystem}
              expanded={expanded}
              toggle={toggle}
            />
          ))}
        </div>
      )}
    </div>
  );
}

function SearchResults({
  results,
  truncated,
}: {
  results: { subsystem: Subsystem; kind: KindBucket; cluster: ErrorCluster }[];
  truncated: boolean;
}) {
  if (results.length === 0) {
    return <p className="text-sm text-fd-muted-foreground">No matches.</p>;
  }
  return (
    <div className="flex flex-col gap-2">
      {truncated && (
        <p className="text-xs text-fd-muted-foreground">
          Showing the first {MAX_SEARCH_RESULTS} matches — refine your search for more specific results.
        </p>
      )}
      {results.map(({ subsystem, kind, cluster }) => (
        <div key={cluster.key} className="border rounded-md border-fd-border">
          <div className="px-3 py-2 text-xs text-fd-muted-foreground border-b border-fd-border">
            {subsystem.name} <ChevronRight className="inline size-3" /> {kind.kind}
          </div>
          <ClusterDetail cluster={cluster} />
        </div>
      ))}
    </div>
  );
}

function SubsystemSection({
  subsystem,
  expanded,
  toggle,
}: {
  subsystem: Subsystem;
  expanded: Set<string>;
  toggle: (key: string) => void;
}) {
  const open = expanded.has(subsystem.id);
  return (
    <div className="border rounded-md border-fd-border">
      <button
        type="button"
        onClick={() => toggle(subsystem.id)}
        aria-expanded={open}
        className="flex items-start w-full gap-2 px-3 py-2 text-left"
      >
        <ChevronRight className={cn('size-4 shrink-0 mt-0.5 transition-transform text-fd-muted-foreground', open && 'rotate-90')} />
        <div className="flex flex-col gap-1">
          <span className="text-sm font-medium text-fd-foreground">
            {subsystem.name} <span className="font-normal text-fd-muted-foreground">({subsystemIdCount(subsystem)} IDs)</span>
          </span>
          {open && <span className="text-xs text-fd-muted-foreground">{subsystem.summary}</span>}
        </div>
      </button>
      {open && (
        <div className="flex flex-col gap-2 px-3 pb-3 pl-9">
          {subsystem.kinds.map((kind) => (
            <KindSection
              key={kind.kind}
              parentKey={subsystem.id}
              kind={kind}
              expanded={expanded}
              toggle={toggle}
            />
          ))}
        </div>
      )}
    </div>
  );
}

function KindSection({
  parentKey,
  kind,
  expanded,
  toggle,
}: {
  parentKey: string;
  kind: KindBucket;
  expanded: Set<string>;
  toggle: (key: string) => void;
}) {
  const key = `${parentKey}::${kind.kind}`;
  const open = expanded.has(key);
  return (
    <div className="border rounded-md border-fd-border/60">
      <button type="button" onClick={() => toggle(key)} aria-expanded={open} className="flex items-center w-full gap-2 px-3 py-1.5 text-left">
        <ChevronRight className={cn('size-3.5 shrink-0 transition-transform text-fd-muted-foreground', open && 'rotate-90')} />
        <span className="font-mono text-xs text-fd-foreground">{kind.kind}</span>
        <span className="text-xs text-fd-muted-foreground">
          {kind.clusters.length} cause{kind.clusters.length === 1 ? '' : 's'}, {kindIdCount(kind)} IDs
        </span>
      </button>
      {open && (
        <div className="flex flex-col gap-1.5 px-3 pb-2 ml-5">
          {kind.clusters.map((cluster) => (
            <ClusterRow key={cluster.key} cluster={cluster} expanded={expanded} toggle={toggle} />
          ))}
        </div>
      )}
    </div>
  );
}

function ClusterRow({
  cluster,
  expanded,
  toggle,
}: {
  cluster: ErrorCluster;
  expanded: Set<string>;
  toggle: (key: string) => void;
}) {
  const open = expanded.has(cluster.key);
  return (
    <div className="border rounded-md border-fd-border/60">
      <button type="button" onClick={() => toggle(cluster.key)} aria-expanded={open} className="flex items-start w-full gap-2 px-2 py-1.5 text-left">
        <ChevronRight className={cn('size-3.5 shrink-0 mt-0.5 transition-transform text-fd-muted-foreground', open && 'rotate-90')} />
        <span className="flex-1 font-mono text-xs break-all text-fd-foreground">{cluster.message || '(empty message)'}</span>
        {cluster.collidesWith.length > 0 && <TriangleAlert className="size-3.5 shrink-0 text-amber-500" />}
        <span className="text-xs shrink-0 text-fd-muted-foreground">{cluster.ids.length}</span>
      </button>
      {open && <ClusterDetail cluster={cluster} />}
    </div>
  );
}

function ClusterDetail({ cluster }: { cluster: ErrorCluster }) {
  return (
    <div className="flex flex-col gap-3 px-3 pb-3 text-xs">
      <div className="flex items-start justify-between gap-2">
        <p className="text-fd-foreground">{cluster.why}</p>
        <span className={cn('shrink-0 rounded px-1.5 py-0.5 text-[10px] font-medium whitespace-nowrap', WHY_SOURCE_STYLE[cluster.whySource])}>
          {WHY_SOURCE_LABEL[cluster.whySource]}
        </span>
      </div>

      {cluster.collidesWith.length > 0 && (
        <div className="flex flex-col gap-1 p-2 border rounded-md border-amber-500/40 bg-amber-500/10">
          <p className="flex items-center gap-1 font-medium text-amber-600 dark:text-amber-400">
            <TriangleAlert className="size-3.5" /> Same ID string used elsewhere with a different meaning
          </p>
          <p className="text-fd-muted-foreground">
            Don&apos;t disambiguate by ID alone — a client also needs the message. Also thrown for:
          </p>
          <ul className="list-disc list-inside text-fd-muted-foreground">
            {cluster.collidesWith.map((c) => (
              <li key={c.key}>
                <span className="font-mono">&ldquo;{c.message}&rdquo;</span> ({c.subsystem})
              </li>
            ))}
          </ul>
        </div>
      )}

      <div className="flex flex-col gap-1.5">
        <p className="font-medium text-fd-muted-foreground">{cluster.ids.length === 1 ? 'ID' : `IDs (${cluster.ids.length})`}</p>
        <div className="flex flex-col gap-2">
          {cluster.ids.map((idEntry) => (
            <div key={idEntry.id} className="flex flex-col gap-0.5">
              <code className="font-mono text-fd-foreground">{idEntry.id}</code>
              <div className="flex flex-wrap gap-x-3 gap-y-0.5 pl-3 text-fd-muted-foreground">
                {idEntry.locations.slice(0, MAX_LOCATIONS_SHOWN).map((loc, i) => (
                  <span key={i} className="font-mono">
                    {loc.file}:{loc.line}
                  </span>
                ))}
                {idEntry.locations.length > MAX_LOCATIONS_SHOWN && <span>+{idEntry.locations.length - MAX_LOCATIONS_SHOWN} more</span>}
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="flex flex-col gap-1">
        <p className="font-medium text-fd-muted-foreground">
          Example response &mdash; HTTP {cluster.example.httpStatus}, gRPC {cluster.example.grpcCode}
        </p>
        <pre className="p-2 overflow-x-auto rounded-md bg-fd-secondary text-fd-secondary-foreground">
          <code>{JSON.stringify(cluster.example.body, null, 2)}</code>
        </pre>
      </div>
    </div>
  );
}
