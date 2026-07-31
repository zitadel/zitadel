'use client';

import { useState, type ReactNode } from 'react';
import Link from 'next/link';
import { ChevronRight } from 'lucide-react';
import { cn } from '@/utils/cn';
import { GRPC_STATUS } from '@/lib/grpc-status';
import rawData from './data.json';

type ExampleResponse = {
  httpStatus: number;
  grpcCode: number;
  body: { code: number; message: string; details: { '@type': string; id: string; message: string }[] };
};
type Cause = { key: string; id: string; message: string; why: string; example: ExampleResponse };
type StatusGroup = { status: number; statusText: string; description?: string; causes: Cause[] };
const data = rawData as Record<string, Record<string, StatusGroup[]>>;

const GENERIC_TEMPLATE = `{
  "code": "<gRPC status code (0-16)>",
  "message": "<Human-readable error> (<Error ID>)",
  "details": [
    {
      "@type": "type.googleapis.com/zitadel.v1.ErrorDetail",
      "id": "<Error ID — maps to a source code location>",
      "message": "<Internal error key or human-readable description>"
    }
  ]
}`;

const FIELD_MEANINGS = [
  { field: 'code', what: 'The gRPC status code, mapped to an HTTP status.', example: '3 = INVALID_ARGUMENT = HTTP 400' },
  { field: 'message', what: 'What went wrong, plus the error ID in parentheses.', example: '"Errors.User.Metadata.KeyEmpty (USER-Drght)"' },
  { field: 'details[].id', what: 'The error ID from the source code — the exact thing to search for.', example: 'USER-Drght' },
  {
    field: 'details[].message',
    what: 'The diagnostic text: either a translated sentence or an internal-only key.',
    example: '"Errors.User.Metadata.KeyEmpty"',
  },
];

// A representative subset of the codes that actually show up across this
// API's errors (verified against components/ErrorReference/data.json),
// not the full gRPC code list — see gRPC Status Codes for that.
const COMMON_CODE_NAMES = ['OK', 'INVALID_ARGUMENT', 'NOT_FOUND', 'ALREADY_EXISTS', 'PERMISSION_DENIED', 'FAILED_PRECONDITION', 'UNAUTHENTICATED', 'INTERNAL'];

function Table({ head, rows }: { head: string[]; rows: ReactNode[][] }) {
  return (
    <div className="overflow-x-auto rounded-md border border-fd-border">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-fd-border bg-fd-secondary/50">
            {head.map((h) => (
              <th key={h} className="px-3 py-2 text-start font-medium text-fd-foreground">
                {h}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-fd-border">
          {rows.map((row, i) => (
            <tr key={i}>
              {row.map((cell, j) => (
                <td key={j} className="px-3 py-2 align-top text-fd-muted-foreground">
                  {cell}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function ErrorBodySchema() {
  const commonCodes = COMMON_CODE_NAMES.map((name) => GRPC_STATUS[name]).filter(Boolean);
  return (
    <div className="not-prose flex flex-col rounded-xl border shadow-md overflow-hidden bg-fd-card text-fd-card-foreground p-3 gap-3">
      <p className="text-sm text-fd-foreground">Every error on this page comes back in this shape:</p>
      <pre className="p-2 overflow-x-auto text-xs rounded-md bg-fd-secondary text-fd-secondary-foreground">
        <code>{GENERIC_TEMPLATE}</code>
      </pre>

      <p className="text-sm font-medium text-fd-foreground">What each field means</p>
      <Table
        head={['Field', 'What it is', 'Example']}
        rows={FIELD_MEANINGS.map((f) => [<code key="f">{f.field}</code>, f.what, <code key="e">{f.example}</code>])}
      />

      <p className="text-sm font-medium text-fd-foreground">Common codes you&apos;ll see</p>
      <Table
        head={['code', 'gRPC name', 'HTTP', 'Meaning']}
        rows={commonCodes.map((s) => [<code key="c">{s.code}</code>, s.name, <code key="h">{s.httpStatus}</code>, s.httpText])}
      />
      <p className="text-xs text-fd-muted-foreground">
        Full mapping: <Link href="/apis/statuscodes" className="underline decoration-fd-border hover:text-fd-primary">gRPC Status Codes</Link>
      </p>
    </div>
  );
}

export default function EndpointErrors({ service, operationId }: { service: string; operationId: string }) {
  const groups = data[service]?.[operationId];
  const [open, setOpen] = useState<number | null>(null);
  const [selected, setSelected] = useState<Record<number, string>>({});

  if (!groups || groups.length === 0) return null;

  return (
    <div className="not-prose flex flex-col gap-3">
      <ErrorBodySchema />
      <div className="not-prose flex flex-col rounded-xl border shadow-md overflow-hidden bg-fd-card text-fd-card-foreground">
      <div className="divide-y divide-fd-border">
        {groups.map((g) => {
          const isOpen = open === g.status;
          const selectedKey = selected[g.status] ?? g.causes[0]?.key;
          const cause = g.causes.find((c) => c.key === selectedKey);
          return (
            <div key={g.status} className="scroll-m-20">
              <h3 className="not-prose flex items-center py-2 px-3 text-fd-foreground font-medium">
                <button
                  type="button"
                  onClick={() => setOpen(isOpen ? null : g.status)}
                  className="flex flex-1 items-center gap-1 text-start group/accordion focus-visible:outline-none font-mono"
                >
                  <ChevronRight
                    className={cn(
                      'size-3.5 text-fd-muted-foreground shrink-0 transition-transform group-focus-visible/accordion:text-fd-primary',
                      isOpen && 'rotate-90',
                    )}
                  />
                  {g.status} {g.statusText}
                </button>
                <p className="text-fd-muted-foreground not-prose">
                  <code className="text-xs">
                    {g.causes.length > 0 ? `${g.causes.length} cause${g.causes.length === 1 ? '' : 's'}` : 'no cause identified'}
                  </code>
                </p>
              </h3>
              {isOpen && (
                <div className="overflow-hidden ps-4.5 pe-3 pb-3 flex flex-col gap-2">
                  {g.description && <p className="text-sm text-fd-muted-foreground">{g.description}</p>}
                  {g.causes.length > 0 && (
                    <>
                      <select
                        value={selectedKey}
                        onChange={(e) => setSelected((prev) => ({ ...prev, [g.status]: e.target.value }))}
                        className="p-2 bg-transparent text-sm font-medium rounded-md border border-fd-border hover:bg-fd-accent/50 focus:outline-none focus:ring-2 focus:ring-fd-ring text-fd-foreground w-full"
                      >
                        {g.causes.map((c) => (
                          <option key={c.key} value={c.key}>
                            {c.message} ({c.id})
                          </option>
                        ))}
                      </select>
                      {cause && (
                        <div className="text-sm flex flex-col gap-2">
                          <div className="flex flex-col gap-1">
                            <code className="text-xs text-fd-muted-foreground">{cause.id}</code>
                            <p className="text-fd-foreground">{cause.why}</p>
                          </div>
                          <div className="flex flex-col gap-1">
                            <p className="text-xs font-medium text-fd-muted-foreground">
                              Example response — a real call to this endpoint that hits this cause comes back exactly like this:
                            </p>
                            <pre className="p-2 overflow-x-auto text-xs rounded-md bg-fd-secondary text-fd-secondary-foreground">
                              <code>{JSON.stringify(cause.example.body, null, 2)}</code>
                            </pre>
                          </div>
                        </div>
                      )}
                    </>
                  )}
                </div>
              )}
            </div>
          );
        })}
      </div>
      </div>
    </div>
  );
}
