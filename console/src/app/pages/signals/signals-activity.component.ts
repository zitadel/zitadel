import { animate, state, style, transition, trigger } from '@angular/animations';
import { CommonModule } from '@angular/common';
import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { GrpcService } from 'src/app/services/grpc.service';
import { ToastService } from 'src/app/services/toast.service';
import { entityTypeFields, buildProtoFilters } from './signal-fields';

import type { Signal, AggregationBucket } from '@zitadel/proto/zitadel/signal/v2/signal_pb.js';

interface TimeRange {
  label: string;
  value: string;
  bucket: string;
}

interface TimelineEntry {
  signal: Signal;
  timeLabel: string;
  isFirstInGroup: boolean;
  groupLabel: string;
  traceColor: string; // color assigned to this trace (empty if singleton or no trace)
  hasTrace: boolean; // true if signal has a valid (non-zero) trace_id
}

interface TraceGroup {
  traceId: string;
  traceColor: string;
  signals: TimelineEntry[];
  // Merged context: all unique non-empty values across the trace per field
  merged: Record<string, string[]>;
  // Trace summary
  operations: string[];
  streams: string[];
  successCount: number;
  failureCount: number;
  firstTime: string;
  lastTime: string;
  durationMs: number | null;
}

type TimelineItem =
  | { type: 'signal'; entry: TimelineEntry; dateGroup: string; isFirstInGroup: boolean }
  | { type: 'trace'; group: TraceGroup; dateGroup: string; isFirstInGroup: boolean };

@Component({
  selector: 'cnsl-signals-activity',
  standalone: true,
  imports: [
    CommonModule,
    TranslateModule,
    MatButtonModule,
    MatIconModule,
    MatMenuModule,
    MatProgressSpinnerModule,
    MatTooltipModule,
  ],
  templateUrl: './signals-activity.component.html',
  styleUrls: ['./signals.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('void', style({ height: '0', opacity: '0', overflow: 'hidden' })),
      state('*', style({ height: '*', opacity: '1' })),
      transition('void <=> *', animate('200ms ease-in-out')),
    ]),
  ],
})
export class SignalsActivityComponent implements OnInit, OnDestroy {
  private readonly grpc = inject(GrpcService);
  private readonly toast = inject(ToastService);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);

  private alive = true;

  signalsAvailable = true;
  loading = false;
  signals: Signal[] = [];
  totalCount = 0;
  offset = 0;
  limit = 100;

  // The entity being traced
  entityType: 'user_id' | 'client_id' | 'org_id' | 'trace_id' = 'user_id';
  entityValue = '';
  searchInput = '';

  // Stats
  operationCounts: AggregationBucket[] = [];
  outcomeCounts: AggregationBucket[] = [];

  // Chronological timeline with trace color hints
  timeline: TimelineEntry[] = [];
  timelineItems: TimelineItem[] = [];
  expandedSignals = new Set<Signal>();
  expandedTraces = new Set<string>();

  // Trace color palette for grouping signals visually
  private readonly traceColors = ['#6366f1', '#22c55e', '#f59e0b', '#06b6d4', '#ef4444', '#8b5cf6', '#ec4899', '#14b8a6'];

  // Entity type display
  readonly entityTypes = entityTypeFields().map(f => ({ key: f.key, label: f.label, icon: f.icon! }));

  timeRanges: TimeRange[] = [
    { label: 'Last 1h', value: '1 hour', bucket: '1 minute' },
    { label: 'Last 6h', value: '6 hours', bucket: '5 minutes' },
    { label: 'Last 24h', value: '24 hours', bucket: '30 minutes' },
    { label: 'Last 7d', value: '7 days', bucket: '3 hours' },
    { label: 'Last 30d', value: '30 days', bucket: '12 hours' },
  ];
  selectedTimeRange: TimeRange = this.timeRanges[2];

  ngOnInit(): void {
    const params = this.route.snapshot.queryParams;
    // Detect entity type from query params
    for (const et of this.entityTypes) {
      if (params[et.key]) {
        this.entityType = et.key as any;
        this.entityValue = params[et.key];
        this.searchInput = params[et.key];
        break;
      }
    }
    // Restore time range
    if (params['time']) {
      const tr = this.timeRanges.find((r) => r.value === params['time']);
      if (tr) this.selectedTimeRange = tr;
    }
    if (this.entityValue) {
      this.refresh();
    }
  }

  ngOnDestroy(): void {
    this.alive = false;
  }

  get entityLabel(): string {
    return this.entityTypes.find((e) => e.key === this.entityType)?.label ?? this.entityType;
  }

  get entityIcon(): string {
    return this.entityTypes.find((e) => e.key === this.entityType)?.icon ?? 'person';
  }

  selectEntityType(key: string): void {
    this.entityType = key as any;
  }

  searchEntity(): void {
    const val = this.searchInput.trim();
    if (!val) return;
    this.entityValue = val;
    this.offset = 0;
    this.expandedSignals.clear();
    this.expandedTraces.clear();
    this.syncUrl();
    this.refresh();
  }

  selectTimeRange(range: TimeRange): void {
    this.selectedTimeRange = range;
    this.offset = 0;
    this.syncUrl();
    this.refresh();
  }

  private syncUrl(): void {
    const params: Record<string, string> = {};
    if (this.entityValue) params[this.entityType] = this.entityValue;
    if (this.selectedTimeRange !== this.timeRanges[2]) params['time'] = this.selectedTimeRange.value;
    this.router.navigate([], { queryParams: params, queryParamsHandling: 'replace', replaceUrl: true });
  }

  refresh(): void {
    this.loadTimeline();
    this.loadStats();
  }

  loadTimeline(): void {
    if (!this.grpc.signal || !this.entityValue) return;
    this.loading = true;

    const filters = buildProtoFilters({ [this.entityType]: this.entityValue });

    this.grpc.signal
      .listSignals({
        query: { offset: BigInt(this.offset), limit: this.limit, asc: false },
        filters,
      })
      .then(
        (resp) => {
          if (!this.alive) return;
          this.signals = resp.signals ?? [];
          this.totalCount = Number(resp.details?.totalResult ?? 0);
          this.buildTimeline();
          this.loading = false;
        },
        (err) => {
          this.loading = false;
          this.handleApiError(err);
        },
      );
  }

  loadStats(): void {
    if (!this.grpc.signal || !this.entityValue) return;

    const filters = buildProtoFilters({ [this.entityType]: this.entityValue });

    this.grpc.signal.aggregateSignals({ filters, groupBy: 'operation', metric: 'count', timeBucket: '' }).then(
      (resp) => {
        if (!this.alive) return;
        this.operationCounts = resp.buckets ?? [];
      },
      (err) => this.handleApiError(err),
    );

    this.grpc.signal.aggregateSignals({ filters, groupBy: 'outcome', metric: 'count', timeBucket: '' }).then(
      (resp) => {
        if (!this.alive) return;
        this.outcomeCounts = resp.buckets ?? [];
      },
      (err) => this.handleApiError(err),
    );
  }

  private buildTimeline(): void {
    this.timeline = [];
    let lastGroup = '';

    const zeroTrace = '00000000000000000000000000000000';

    // Count occurrences of each trace ID to find multi-signal traces.
    const traceCounts = new Map<string, number>();
    for (const s of this.signals) {
      if (s.traceId && s.traceId !== zeroTrace) {
        traceCounts.set(s.traceId, (traceCounts.get(s.traceId) ?? 0) + 1);
      }
    }

    // Assign colors only to traces that appear 2+ times on this page.
    const traceColorMap = new Map<string, string>();
    let colorIdx = 0;
    for (const [traceId, count] of traceCounts) {
      if (count >= 2) {
        traceColorMap.set(traceId, this.traceColors[colorIdx % this.traceColors.length]);
        colorIdx++;
      }
    }

    for (const s of this.signals) {
      const ts = this.toMillis(s.createdAt);
      const d = ts ? new Date(ts) : null;
      const timeLabel = d ? this.formatTime(d) : '—';
      const groupLabel = d ? this.formatDateGroup(d) : '—';
      const isFirstInGroup = groupLabel !== lastGroup;
      lastGroup = groupLabel;

      const hasTrace = !!(s.traceId && s.traceId !== zeroTrace);
      const traceColor = hasTrace ? (traceColorMap.get(s.traceId) ?? '') : '';

      this.timeline.push({ signal: s, timeLabel, isFirstInGroup, groupLabel, traceColor, hasTrace });
    }

    // Second pass: build timelineItems with trace grouping
    this.buildTimelineItems(traceColorMap);
  }

  private buildTimelineItems(traceColorMap: Map<string, string>): void {
    this.timelineItems = [];

    // Group timeline entries by traceId for multi-signal traces
    const traceEntries = new Map<string, TimelineEntry[]>();
    const zeroTrace = '00000000000000000000000000000000';

    for (const entry of this.timeline) {
      const tid = entry.signal.traceId;
      if (tid && tid !== zeroTrace && traceColorMap.has(tid)) {
        if (!traceEntries.has(tid)) {
          traceEntries.set(tid, []);
        }
        traceEntries.get(tid)!.push(entry);
      }
    }

    // Track which traces we've already emitted
    const emittedTraces = new Set<string>();
    let lastDateGroup = '';

    for (const entry of this.timeline) {
      const tid = entry.signal.traceId;
      const isMultiTrace = tid && tid !== zeroTrace && traceColorMap.has(tid);

      if (isMultiTrace) {
        if (emittedTraces.has(tid)) {
          // Skip — already emitted as part of a trace group
          continue;
        }
        emittedTraces.add(tid);

        const signals = traceEntries.get(tid)!;
        const traceColor = traceColorMap.get(tid) ?? '';
        const merged = this.computeMerged(signals.map(e => e.signal));

        // Compute trace summary
        const ops = new Set<string>();
        const streams = new Set<string>();
        let successCount = 0;
        let failureCount = 0;
        let firstMs: number | null = null;
        let lastMs: number | null = null;
        for (const e of signals) {
          if (e.signal.operation) ops.add(this.shortName(e.signal.operation));
          if (e.signal.stream) streams.add(e.signal.stream);
          if (e.signal.outcome === 'success') successCount++;
          if (e.signal.outcome === 'failure') failureCount++;
          const ms = this.toMillis(e.signal.createdAt);
          if (ms !== null) {
            if (firstMs === null || ms < firstMs) firstMs = ms;
            if (lastMs === null || ms > lastMs) lastMs = ms;
          }
        }
        const durationMs = firstMs !== null && lastMs !== null ? lastMs - firstMs : null;
        const firstTime = signals[signals.length - 1]?.timeLabel ?? '';
        const lastTime = signals[0]?.timeLabel ?? '';

        const dateGroup = entry.groupLabel;
        const isFirstInGroup = dateGroup !== lastDateGroup;
        lastDateGroup = dateGroup;

        this.timelineItems.push({
          type: 'trace',
          group: {
            traceId: tid, traceColor, signals, merged,
            operations: Array.from(ops),
            streams: Array.from(streams),
            successCount, failureCount,
            firstTime, lastTime, durationMs,
          },
          dateGroup,
          isFirstInGroup,
        });
      } else {
        const dateGroup = entry.groupLabel;
        const isFirstInGroup = dateGroup !== lastDateGroup;
        lastDateGroup = dateGroup;

        this.timelineItems.push({
          type: 'signal',
          entry,
          dateGroup,
          isFirstInGroup,
        });
      }
    }
  }

  private computeMerged(signals: Signal[]): Record<string, string[]> {
    const merged: Record<string, string[]> = {};
    const fields = ['userId', 'orgId', 'clientId', 'sessionId', 'projectId', 'ip', 'country', 'userAgent'] as const;
    for (const field of fields) {
      const unique = new Set<string>();
      for (const s of signals) {
        const v = (s as any)[field];
        if (v) unique.add(v);
      }
      if (unique.size > 0) {
        merged[field] = Array.from(unique);
      }
    }
    return merged;
  }

  formatDuration(ms: number | null): string {
    if (ms === null) return '';
    if (ms < 1000) return ms + 'ms';
    if (ms < 60000) return (ms / 1000).toFixed(1) + 's';
    if (ms < 3600000) return Math.floor(ms / 60000) + 'm ' + Math.floor((ms % 60000) / 1000) + 's';
    return Math.floor(ms / 3600000) + 'h ' + Math.floor((ms % 3600000) / 60000) + 'm';
  }

  /** Navigate to Activity filtered by a trace ID */
  viewTrace(traceId: string): void {
    this.entityType = 'trace_id';
    this.entityValue = traceId;
    this.searchInput = traceId;
    this.offset = 0;
    this.expandedSignals.clear();
    this.expandedTraces.clear();
    this.syncUrl();
    this.refresh();
  }

  private formatTime(d: Date): string {
    return d.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit', second: '2-digit' });
  }

  private formatDateGroup(d: Date): string {
    const now = new Date();
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const target = new Date(d.getFullYear(), d.getMonth(), d.getDate());
    const diffDays = Math.floor((today.getTime() - target.getTime()) / 86400000);

    if (diffDays === 0) return 'Today';
    if (diffDays === 1) return 'Yesterday';
    return d.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' });
  }

  toggleEntry(signal: Signal): void {
    if (this.expandedSignals.has(signal)) {
      this.expandedSignals.delete(signal);
    } else {
      this.expandedSignals.add(signal);
    }
  }

  toggleTrace(traceId: string): void {
    if (this.expandedTraces.has(traceId)) {
      this.expandedTraces.delete(traceId);
    } else {
      this.expandedTraces.add(traceId);
    }
  }

  hasAnyMerged(merged: Record<string, string[]>): boolean {
    return Object.keys(merged).length > 0;
  }

  /** Check if a child signal's field value is already shown in the merged root */
  isInMerged(merged: Record<string, string[]>, field: string, value: string): boolean {
    return !!(merged[field] && merged[field].includes(value));
  }

  copyToClipboard(value: string, event: MouseEvent): void {
    event.stopPropagation();
    if (!value || value === '—') return;
    navigator.clipboard.writeText(value).then(() => {
      this.toast.showInfo('Copied to clipboard');
    });
  }

  navigateToEntity(type: 'user' | 'org' | 'project', id: string): void {
    if (!id || id === '—') return;
    switch (type) {
      case 'user':
        this.router.navigate(['/users', id]);
        break;
      case 'org':
        this.router.navigate(['/orgs', id]);
        break;
      case 'project':
        this.router.navigate(['/projects', id]);
        break;
    }
  }

  openInLogs(): void {
    this.router.navigate(['/signals/logs'], { queryParams: { [this.entityType]: this.entityValue } });
  }

  viewSignalInLogs(signal: Signal): void {
    const params: Record<string, string> = {};
    if (signal.traceId) params['trace_id'] = signal.traceId;
    else if (signal.operation) params['operation'] = signal.operation;
    // Pass highlight hint so Logs can auto-expand this exact row
    if (signal.operation) params['highlight'] = signal.operation;
    const ts = this.toMillis(signal.createdAt);
    if (ts) params['highlight_ts'] = String(ts);
    this.router.navigate(['/signals/logs'], { queryParams: params });
  }

  openInExplorer(): void {
    this.router.navigate(['/signals/explore'], { queryParams: { [this.entityType]: this.entityValue } });
  }

  toMillis(ts: any): number | null {
    if (!ts?.seconds) return null;
    return Number(ts.seconds) * 1000;
  }

  get totalSignals(): number {
    return this.totalCount;
  }

  get successCount(): number {
    return Number(this.outcomeCounts.find((b) => b.key === 'success')?.count ?? 0);
  }

  get failureCount(): number {
    return Number(this.outcomeCounts.find((b) => b.key === 'failure')?.count ?? 0);
  }

  get topOperations(): { key: string; count: number }[] {
    return this.operationCounts
      .filter((b) => b.key)
      .slice(0, 5)
      .map((b) => ({ key: b.key, count: Number(b.count) }));
  }

  nextPage(): void {
    this.offset += this.limit;
    this.loadTimeline();
  }

  prevPage(): void {
    this.offset = Math.max(0, this.offset - this.limit);
    this.loadTimeline();
  }

  get hasNextPage(): boolean {
    return this.offset + this.limit < this.totalCount;
  }

  get hasPrevPage(): boolean {
    return this.offset > 0;
  }

  get currentPage(): number {
    return Math.floor(this.offset / this.limit) + 1;
  }

  get totalPages(): number {
    return Math.ceil(this.totalCount / this.limit) || 1;
  }

  shortName(name: string): string {
    if (!name) return '';
    if (/^\d{1,3}(\.\d{1,3}){3}$/.test(name) || name.includes(':')) return name;
    const slashParts = name.split('/');
    if (slashParts.length >= 3) return slashParts[slashParts.length - 1];
    const dotParts = name.split('.');
    if (dotParts.length > 3) return dotParts.slice(-3).join('.');
    return name;
  }

  private handleApiError(err: any): void {
    if (this.isServiceUnavailable(err)) {
      this.signalsAvailable = false;
      return;
    }
    this.toast.showError(err);
  }

  private isServiceUnavailable(err: any): boolean {
    const code = err?.code ?? err?.status;
    return code === 12 || code === 5;
  }
}
