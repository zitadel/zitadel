import { animate, state, style, transition, trigger } from '@angular/animations';
import { CommonModule } from '@angular/common';
import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSortModule, Sort } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { GrpcService } from 'src/app/services/grpc.service';
import { ToastService } from 'src/app/services/toast.service';
import { SIGNAL_FIELDS, suggestableFieldKeys, filterLabelMap, filterableFields, buildProtoFilters } from './signal-fields';

import type { MessageInitShape } from '@bufbuild/protobuf';
import type { Signal, AggregationBucket } from '@zitadel/proto/zitadel/signal/v2/signal_pb.js';
import { SignalFiltersSchema } from '@zitadel/proto/zitadel/signal/v2/signal_pb.js';

interface TimeRange {
  label: string;
  value: string;
  bucket: string;
}

@Component({
  selector: 'cnsl-signals-logs',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    TranslateModule,
    MatButtonModule,
    MatIconModule,
    MatMenuModule,
    MatProgressSpinnerModule,
    MatSortModule,
    MatTableModule,
    MatTooltipModule,
  ],
  templateUrl: './signals-logs.component.html',
  styleUrls: ['./signals.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('void', style({ height: '0', opacity: '0', overflow: 'hidden' })),
      state('*', style({ height: '*', opacity: '1' })),
      transition('void <=> *', animate('200ms ease-in-out')),
    ]),
  ],
})
export class SignalsLogsComponent implements OnInit, OnDestroy {
  private readonly grpc = inject(GrpcService);
  private readonly fb = inject(FormBuilder);
  private readonly toast = inject(ToastService);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);

  private alive = true;

  signalsAvailable = true;
  loading = false;
  signals: Signal[] = [];
  sortedSignals: Signal[] = [];
  totalCount = 0;
  offset = 0;
  limit = 50;

  // Stream/outcome counts
  streamCounts: AggregationBucket[] = [];
  outcomeCounts: AggregationBucket[] = [];
  streams: string[] = [];
  dimensionCounts: Record<string, number> = {};

  expandedSignals = new Set<Signal>();
  highlightedSignal: Signal | null = null;

  // Filter prompt state
  pendingFilterKey = '';
  pendingFilterLabel = '';
  filterSuggestions: { key: string; count: number }[] = [];
  filterSuggestionsLoading = false;
  filterInputValue = '';

  private readonly suggestableFields = suggestableFieldKeys();
  readonly filterFieldDefs = filterableFields();

  filterForm: FormGroup = this.fb.group({
    stream: [''],
    outcome: [''],
    ...Object.fromEntries(filterableFields().map(f => [f.key, ['']])),
  });

  private readonly _filterLabelMap = filterLabelMap();

  displayedColumns = ['createdAt', 'stream', 'userId', 'operation', 'outcome', 'ip', 'expand'];

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
    // Restore filters from URL
    const patchable: Record<string, string> = {};
    for (const key of Object.keys(this.filterForm.controls)) {
      if (params[key]) patchable[key] = params[key];
    }
    if (Object.keys(patchable).length) {
      this.filterForm.patchValue(patchable);
    }
    // Restore time range
    if (params['time']) {
      const tr = this.timeRanges.find((r) => r.value === params['time']);
      if (tr) this.selectedTimeRange = tr;
    }
    // Capture highlight hint from Activity "View in Logs"
    if (params['highlight']) {
      this._highlightOp = params['highlight'];
      this._highlightTs = params['highlight_ts'] ? Number(params['highlight_ts']) : 0;
    }
    this.refresh();
  }

  ngOnDestroy(): void {
    this.alive = false;
  }

  private _highlightOp = '';
  private _highlightTs = 0;

  refresh(): void {
    this.syncUrl();
    this.loadDimensions();
    this.search();
  }

  selectTimeRange(range: TimeRange): void {
    this.selectedTimeRange = range;
    this.offset = 0;
    this.refresh();
  }

  selectStream(stream: string): void {
    this.filterForm.patchValue({ stream });
    this.offset = 0;
    this.refresh();
  }

  selectOutcome(outcome: string): void {
    this.filterForm.patchValue({ outcome });
    this.offset = 0;
    this.refresh();
  }

  toggleRow(signal: Signal, event: MouseEvent): void {
    event.stopPropagation();
    if (this.expandedSignals.has(signal)) {
      this.expandedSignals.delete(signal);
    } else {
      this.expandedSignals.add(signal);
    }
  }

  sortData(sort: Sort): void {
    if (!sort.active || sort.direction === '') {
      this.sortedSignals = [...this.signals];
      return;
    }
    const dir = sort.direction === 'asc' ? 1 : -1;
    this.sortedSignals = [...this.signals].sort((a, b) => {
      switch (sort.active) {
        case 'createdAt': {
          const ta = Number(a.createdAt?.seconds ?? 0);
          const tb = Number(b.createdAt?.seconds ?? 0);
          return (ta - tb) * dir;
        }
        case 'stream':
          return (a.stream ?? '').localeCompare(b.stream ?? '') * dir;
        case 'operation':
          return (a.operation ?? '').localeCompare(b.operation ?? '') * dir;
        case 'outcome':
          return (a.outcome ?? '').localeCompare(b.outcome ?? '') * dir;
        case 'ip':
          return (a.ip ?? '').localeCompare(b.ip ?? '') * dir;
        case 'userId':
          return (a.userId ?? '').localeCompare(b.userId ?? '') * dir;
        default:
          return 0;
      }
    });
  }

  hasActiveFilters(): boolean {
    const f = this.filterForm.value;
    return Object.entries(f).some(([k, v]) => !!v && k !== 'stream' && k !== 'outcome');
  }

  activeFilterChips(): { key: string; label: string; value: string }[] {
    const f = this.filterForm.value;
    const chips: { key: string; label: string; value: string }[] = [];
    for (const [key, value] of Object.entries(f)) {
      if (value && key !== 'stream' && key !== 'outcome') {
        chips.push({ key, label: this._filterLabelMap[key] || key, value: value as string });
      }
    }
    return chips;
  }

  openFilterInput(key: string): void {
    this.pendingFilterKey = key;
    this.pendingFilterLabel = this._filterLabelMap[key] || key;
    this.filterSuggestions = [];
    this.filterInputValue = '';
    if (this.suggestableFields.has(key)) {
      this.loadFilterSuggestions(key);
    }
  }

  applyPendingFilter(value: string): void {
    if (value && this.pendingFilterKey) {
      const trimmed = value.trim().substring(0, 512);
      this.filterForm.patchValue({ [this.pendingFilterKey]: trimmed });
      this.pendingFilterKey = '';
      this.pendingFilterLabel = '';
      this.filterSuggestions = [];
      this.filterInputValue = '';
      this.offset = 0;
      this.refresh();
    }
  }

  cancelPendingFilter(): void {
    this.pendingFilterKey = '';
    this.pendingFilterLabel = '';
    this.filterSuggestions = [];
    this.filterInputValue = '';
  }

  loadFilterSuggestions(field: string): void {
    if (!this.grpc.signal) return;
    this.filterSuggestionsLoading = true;
    this.grpc.signal
      .aggregateSignals({
        filters: this.buildFilters(field),
        groupBy: field,
        metric: 'count',
        timeBucket: '',
      })
      .then(
        (resp) => {
          if (!this.alive) return;
          this.filterSuggestions = (resp.buckets ?? [])
            .filter((b) => b.key)
            .map((b) => ({ key: b.key, count: Number(b.count) }))
            .slice(0, 50);
          this.filterSuggestionsLoading = false;
        },
        (err) => {
          this.filterSuggestionsLoading = false;
          this.handleApiError(err);
        },
      );
  }

  filteredSuggestions(): { key: string; count: number }[] {
    if (!this.filterInputValue) return this.filterSuggestions;
    const search = this.filterInputValue.toLowerCase();
    return this.filterSuggestions.filter((s) => s.key.toLowerCase().includes(search));
  }

  selectFilterSuggestion(value: string): void {
    if (this.pendingFilterKey) {
      this.filterForm.patchValue({ [this.pendingFilterKey]: value });
      this.pendingFilterKey = '';
      this.pendingFilterLabel = '';
      this.filterSuggestions = [];
      this.filterInputValue = '';
      this.offset = 0;
      this.refresh();
    }
  }

  clearFilter(key: string): void {
    this.filterForm.patchValue({ [key]: '' });
    this.offset = 0;
    this.refresh();
  }

  resetFilters(): void {
    const stream = this.filterForm.get('stream')?.value;
    const outcome = this.filterForm.get('outcome')?.value;
    this.filterForm.reset();
    this.filterForm.patchValue({ stream, outcome });
    this.offset = 0;
    this.refresh();
  }

  correlate(signal: Signal): void {
    if (!signal.traceId) return;
    this.filterForm.reset();
    this.filterForm.patchValue({ trace_id: signal.traceId });
    this.offset = 0;
    this.refresh();
  }

  addFilter(filterKey: string, value: string): void {
    if (!value || value === '—') return;
    this.filterForm.patchValue({ [filterKey]: value });
    this.offset = 0;
    this.refresh();
  }

  private readonly validActivityFields = new Set(['user_id', 'client_id', 'org_id', 'trace_id', 'session_id', 'ip']);

  viewActivity(entityType: string, entityValue: string): void {
    if (!this.validActivityFields.has(entityType)) return;
    if (!entityValue || entityValue === '—') return;
    this.router.navigate(['/signals/activity'], { queryParams: { [entityType]: entityValue } });
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

  copyToClipboard(value: string, event: MouseEvent): void {
    event.stopPropagation();
    if (!value || value === '—') return;
    navigator.clipboard.writeText(value).then(() => {
      this.toast.showInfo('Copied to clipboard');
    });
  }

  private buildFilters(excludeField?: string): MessageInitShape<typeof SignalFiltersSchema> {
    return buildProtoFilters(this.filterForm.value, excludeField);
  }

  private syncUrl(): void {
    const params: Record<string, string> = {};
    const f = this.filterForm.value;
    for (const [key, val] of Object.entries(f)) {
      if (val) params[key] = val as string;
    }
    if (this.selectedTimeRange !== this.timeRanges[2]) params['time'] = this.selectedTimeRange.value;
    this.router.navigate([], { queryParams: params, queryParamsHandling: 'replace', replaceUrl: true });
  }

  search(): void {
    if (!this.grpc.signal) return;
    this.loading = true;
    this.grpc.signal
      .listSignals({
        query: { offset: BigInt(this.offset), limit: this.limit, asc: false },
        filters: this.buildFilters(),
      })
      .then(
        (resp) => {
          if (!this.alive) return;
          this.signals = resp.signals ?? [];
          this.sortedSignals = [...this.signals];
          this.totalCount = Number(resp.details?.totalResult ?? 0);
          this.loading = false;
          this.tryHighlight();
        },
        (err) => {
          this.loading = false;
          this.handleApiError(err);
        },
      );
  }

  /** Auto-expand and highlight the row matching the hint from Activity → "View in Logs" */
  private tryHighlight(): void {
    if (!this._highlightOp) return;
    let best: Signal | null = null;
    let bestDiff = Infinity;
    for (const s of this.sortedSignals) {
      if (s.operation !== this._highlightOp) continue;
      if (this._highlightTs) {
        const sTs = Number(s.createdAt?.seconds ?? 0) * 1000;
        const diff = Math.abs(sTs - this._highlightTs);
        if (diff < bestDiff) {
          bestDiff = diff;
          best = s;
        }
      } else {
        best = s;
        break;
      }
    }
    if (best) {
      this.expandedSignals.add(best);
      this.highlightedSignal = best;
    }
    // Only highlight on first load
    this._highlightOp = '';
    this._highlightTs = 0;
  }

  loadDimensions(): void {
    if (!this.grpc.signal) return;
    this.grpc.signal
      .aggregateSignals({ filters: this.buildFilters(), groupBy: 'stream', metric: 'count', timeBucket: '' })
      .then(
        (resp) => {
          if (!this.alive) return;
          this.streamCounts = resp.buckets ?? [];
          this.streams = this.streamCounts.map((b) => b.key).filter((k) => k);
        },
        (err) => this.handleApiError(err),
      );
    this.grpc.signal
      .aggregateSignals({ filters: this.buildFilters(), groupBy: 'outcome', metric: 'count', timeBucket: '' })
      .then(
        (resp) => {
          if (!this.alive) return;
          this.outcomeCounts = resp.buckets ?? [];
        },
        (err) => this.handleApiError(err),
      );
    for (const field of this.filterFieldDefs.filter(f => f.suggestable).map(f => f.key)) {
      this.grpc.signal
        .aggregateSignals({ filters: this.buildFilters(), groupBy: field, metric: 'count', timeBucket: '' })
        .then(
          (resp) => {
            if (!this.alive) return;
            this.dimensionCounts[field] = (resp.buckets ?? []).filter((b) => b.key).length;
          },
          (err) => this.handleApiError(err),
        );
    }
  }

  getDimensionCount(buckets: AggregationBucket[], key: string): number {
    return Number(buckets.find((b) => b.key === key)?.count ?? 0);
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

  toMillis(ts: any): number | null {
    if (!ts?.seconds) return null;
    return Number(ts.seconds) * 1000;
  }

  extractError(signal: Signal): string | null {
    if (!signal.payload) return null;
    try {
      const obj = JSON.parse(signal.payload);
      return obj?.error || obj?.Error || obj?.message || null;
    } catch {
      return null;
    }
  }

  nextPage(): void {
    this.offset += this.limit;
    this.search();
  }

  prevPage(): void {
    this.offset = Math.max(0, this.offset - this.limit);
    this.search();
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
