'use client';

import { useState } from 'react';
import Link from 'next/link';
import { ChevronRight } from 'lucide-react';
import rawData from './data.json';

type ExampleResponse = {
  httpStatus: number;
  grpcCode: number;
  body: { code: number; message: string; details: { '@type': string; id: string; message: string }[] };
};
type Cause = { key: string; id: string; message: string; why: string; example: ExampleResponse };
type StatusGroup = { status: number; statusText: string; description?: string; causes: Cause[] };
const data = rawData as Record<string, Record<string, StatusGroup[]>>;

export default function EndpointErrors({ service, operationId }: { service: string; operationId: string }) {
  const groups = data[service]?.[operationId];
  const [selected, setSelected] = useState<Record<number, string>>({});

  // No entry at all means this operation hasn't been checked — stay silent
  // rather than expose that as a public-facing gap. An explicit empty array
  // means it *was* checked and confidently found nothing, which is worth
  // saying rather than looking identical to "never checked".
  if (!groups) return null;
  if (groups.length === 0) {
    return (
      <div className="not-prose rounded-xl border bg-fd-card p-3 text-sm text-fd-muted-foreground">
        No specific error causes are currently documented for this operation.
      </div>
    );
  }

  return (
    <div className="not-prose flex flex-col gap-3">
      <p className="text-sm text-fd-muted-foreground">
        Every error below comes back in the shape described on the{' '}
        <Link href="/apis/errors" className="underline decoration-fd-border hover:text-fd-primary">
          error reference
        </Link>
        .
      </p>
      <div className="not-prose flex flex-col rounded-xl border shadow-md overflow-hidden bg-fd-card text-fd-card-foreground">
        <div className="divide-y divide-fd-border">
          {groups.map((g) => {
            const selectedKey = selected[g.status] ?? g.causes[0]?.key;
            const cause = g.causes.find((c) => c.key === selectedKey);
            return (
              // <details>/<summary> instead of a hand-rolled button + state:
              // free keyboard handling, free ARIA expanded state, and —
              // unlike the button version — text inside stays findable via
              // Ctrl-F even while collapsed.
              <details key={g.status} className="group scroll-m-20">
                <summary className="flex cursor-pointer list-none items-center py-2 px-3 text-fd-foreground font-medium [&::-webkit-details-marker]:hidden">
                  <span className="flex flex-1 items-center gap-1 font-mono">
                    <ChevronRight className="size-3.5 shrink-0 text-fd-muted-foreground transition-transform group-open:rotate-90" />
                    {g.status} {g.statusText}
                  </span>
                  <code className="text-xs text-fd-muted-foreground">
                    {g.causes.length > 0 ? `${g.causes.length} cause${g.causes.length === 1 ? '' : 's'}` : 'no cause identified'}
                  </code>
                </summary>
                <div className="overflow-hidden ps-[1.125rem] pe-3 pb-3 flex flex-col gap-2">
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
              </details>
            );
          })}
        </div>
      </div>
    </div>
  );
}
