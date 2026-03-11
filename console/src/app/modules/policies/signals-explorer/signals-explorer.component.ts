import { animate, state, style, transition, trigger } from '@angular/animations';
import { CommonModule } from '@angular/common';
import { Component, inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatChipsModule } from '@angular/material/chips';
import { TranslateModule } from '@ngx-translate/core';
import { GrpcService } from 'src/app/services/grpc.service';
import { ToastService } from 'src/app/services/toast.service';

import { Signal, SignalFilters, AggregationBucket, Finding } from 'src/app/proto/generated/zitadel/signal/v2/signal_pb';
import { ListSignalsRequest, AggregateSignalsRequest } from 'src/app/proto/generated/zitadel/signal/v2/signal_service_pb';
import { ListQuery } from 'src/app/proto/generated/zitadel/object/v2/object_pb';

interface TimeRange {
  label: string;
  value: string;
  bucket: string;
}

interface BreakdownRow {
  key: string;
  count: number;
  pct: number;
}

type Tab = 'findings' | 'overview' | 'logs';

@Component({
  selector: 'cnsl-signals-explorer',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    TranslateModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    MatTableModule,
    MatInputModule,
    MatFormFieldModule,
    MatTooltipModule,
    MatChipsModule,
  ],
  templateUrl: './signals-explorer.component.html',
  styleUrls: ['./signals-explorer.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('void', style({ height: '0', opacity: '0', overflow: 'hidden' })),
      state('*', style({ height: '*', opacity: '1' })),
      transition('void <=> *', animate('200ms ease-in-out')),
    ]),
  ],
})
export class SignalsExplorerComponent implements OnInit {
  private readonly grpc = inject(GrpcService);
  private readonly fb = inject(FormBuilder);
  private readonly toast = inject(ToastService);

  // Navigation
  activeTab: Tab = 'findings';

  // Loading
  loading = false;

  // List data (logs tab)
  signals: Signal.AsObject[] = [];
  totalCount = 0;
  offset = 0;
  limit = 50;

  // Chart
  chartBuckets: AggregationBucket.AsObject[] = [];
  chartLoading = false;
  chartPath = '';
  chartMaxCount = 0;
  chartWidth = 960;
  chartHeight = 160;

  // Summary metrics
  streamCounts: AggregationBucket.AsObject[] = [];
  outcomeCounts: AggregationBucket.AsObject[] = [];
  streams: string[] = [];

  // Breakdown aggregations (overview tab)
  topOperations: BreakdownRow[] = [];
  topResources: BreakdownRow[] = [];
  topIPs: BreakdownRow[] = [];
  topCountries: BreakdownRow[] = [];
  topUsers: BreakdownRow[] = [];

  // Findings tab data
  findingsTopRules: BreakdownRow[] = [];
  findingsTopUsers: BreakdownRow[] = [];
  findingsBlockedCount = 0;
  findingsChallengedCount = 0;
  findingsTotalCount = 0;
  findingsRecentSignals: Signal.AsObject[] = [];
  findingsLoading = false;

  // Expanded row (logs tab)
  expandedSignal: Signal.AsObject | null = null;

  filterForm: FormGroup = this.fb.group({
    stream: [''],
    outcome: [''],
    operation: [''],
    ip: [''],
    country: [''],
    user_id: [''],
    payload: [''],
    trace_id: [''],
    span_id: [''],
  });

  displayedColumns = ['createdAt', 'stream', 'resource', 'operation', 'outcome', 'ip', 'userId', 'findings', 'expand'];

  timeRanges: TimeRange[] = [
    { label: '1h', value: '1 hour', bucket: '1 minute' },
    { label: '6h', value: '6 hours', bucket: '5 minutes' },
    { label: '24h', value: '24 hours', bucket: '30 minutes' },
    { label: '7d', value: '7 days', bucket: '3 hours' },
    { label: '30d', value: '30 days', bucket: '12 hours' },
  ];
  selectedTimeRange: TimeRange = this.timeRanges[2];

  ngOnInit(): void {
    this.refresh();
    this.loadFindings();
  }

  switchTab(tab: Tab): void {
    this.activeTab = tab;
    if (tab === 'logs' && this.signals.length === 0) {
      this.search();
    }
    if (tab === 'findings') {
      this.loadFindings();
    }
  }

  refresh(): void {
    this.loadChart();
    this.loadDimensions();
    this.loadBreakdowns();
    if (this.activeTab === 'logs') {
      this.search();
    }
    if (this.activeTab === 'findings') {
      this.loadFindings();
    }
  }

  selectTimeRange(range: TimeRange): void {
    this.selectedTimeRange = range;
    this.offset = 0;
    this.refresh();
  }

  toggleStream(stream: string): void {
    const current = this.filterForm.get('stream')?.value;
    this.filterForm.patchValue({ stream: current === stream ? '' : stream });
    this.offset = 0;
    this.refresh();
  }

  toggleOutcome(outcome: string): void {
    const current = this.filterForm.get('outcome')?.value;
    this.filterForm.patchValue({ outcome: current === outcome ? '' : outcome });
    this.offset = 0;
    this.refresh();
  }

  toggleRow(signal: Signal.AsObject, event: MouseEvent): void {
    event.stopPropagation();
    this.expandedSignal = this.expandedSignal === signal ? null : signal;
  }

  drillDown(field: string, value: string): void {
    this.filterForm.patchValue({ [field]: value });
    this.activeTab = 'logs';
    this.offset = 0;
    this.refresh();
    this.search();
  }

  // Filters helper for all aggregate calls
  private buildFilters(): SignalFilters {
    const f = this.filterForm.value;
    const filters = new SignalFilters();
    if (f.stream) filters.setStream(f.stream);
    if (f.outcome) filters.setOutcome(f.outcome);
    if (f.operation) filters.setOperation(f.operation);
    if (f.ip) filters.setIp(f.ip);
    if (f.country) filters.setCountry(f.country);
    if (f.user_id) filters.setUserId(f.user_id);
    if (f.payload) filters.setPayload(f.payload);
    if (f.trace_id) filters.setTraceId(f.trace_id);
    if (f.span_id) filters.setSpanId(f.span_id);
    return filters;
  }

  search(): void {
    if (!this.grpc.signal) return;
    this.loading = true;

    const query = new ListQuery();
    query.setOffset(this.offset);
    query.setLimit(this.limit);

    const req = new ListSignalsRequest();
    req.setQuery(query);
    req.setFilters(this.buildFilters());

    this.grpc.signal.listSignals(req, null).then(
      (resp) => {
        this.signals = resp.getSignalsList().map((s) => s.toObject());
        this.totalCount = resp.getDetails()?.getTotalResult() ?? 0;
        this.loading = false;
      },
      (err) => {
        this.toast.showError(err);
        this.loading = false;
      },
    );
  }

  loadChart(): void {
    if (!this.grpc.signal) return;
    this.chartLoading = true;

    const req = new AggregateSignalsRequest();
    req.setFilters(this.buildFilters());
    req.setGroupBy('time_bucket');
    req.setMetric('count');
    req.setTimeBucket(this.selectedTimeRange.bucket);

    this.grpc.signal.aggregateSignals(req, null).then(
      (resp) => {
        this.chartBuckets = resp.getBucketsList().map((b) => b.toObject());
        this.buildChartPath();
        this.chartLoading = false;
      },
      (err) => {
        this.toast.showError(err);
        this.chartLoading = false;
      },
    );
  }

  loadDimensions(): void {
    if (!this.grpc.signal) return;

    // Stream counts
    const streamReq = new AggregateSignalsRequest();
    streamReq.setFilters(this.buildFilters());
    streamReq.setGroupBy('stream');
    streamReq.setMetric('count');
    this.grpc.signal.aggregateSignals(streamReq, null).then((resp) => {
      this.streamCounts = resp.getBucketsList().map((b) => b.toObject());
      this.streams = this.streamCounts.map((b) => b.key).filter((k) => k);
    });

    // Outcome counts
    const outcomeReq = new AggregateSignalsRequest();
    outcomeReq.setFilters(this.buildFilters());
    outcomeReq.setGroupBy('outcome');
    outcomeReq.setMetric('count');
    this.grpc.signal.aggregateSignals(outcomeReq, null).then((resp) => {
      this.outcomeCounts = resp.getBucketsList().map((b) => b.toObject());
    });
  }

  loadBreakdowns(): void {
    if (!this.grpc.signal) return;
    const fields = [
      { groupBy: 'operation', target: 'topOperations' as const },
      { groupBy: 'resource', target: 'topResources' as const },
      { groupBy: 'ip', target: 'topIPs' as const },
      { groupBy: 'country', target: 'topCountries' as const },
      { groupBy: 'user_id', target: 'topUsers' as const },
    ];
    for (const f of fields) {
      const req = new AggregateSignalsRequest();
      req.setFilters(this.buildFilters());
      req.setGroupBy(f.groupBy);
      req.setMetric('count');
      this.grpc.signal.aggregateSignals(req, null).then((resp) => {
        const buckets = resp.getBucketsList().map((b) => b.toObject());
        const total = buckets.reduce((s, b) => s + b.count, 0) || 1;
        this[f.target] = buckets
          .filter((b) => b.key)
          .slice(0, 10)
          .map((b) => ({ key: b.key, count: b.count, pct: (b.count / total) * 100 }));
      });
    }
  }

  loadFindings(): void {
    if (!this.grpc.signal) return;
    this.findingsLoading = true;

    // Top rules: aggregate detection stream by resource (which is "rule:<id>")
    const ruleReq = new AggregateSignalsRequest();
    const ruleFilters = this.buildFilters();
    ruleFilters.setStream('detection');
    ruleReq.setFilters(ruleFilters);
    ruleReq.setGroupBy('resource');
    ruleReq.setMetric('count');
    this.grpc.signal.aggregateSignals(ruleReq, null).then((resp) => {
      const buckets = resp.getBucketsList().map((b) => b.toObject());
      const total = buckets.reduce((s, b) => s + b.count, 0) || 1;
      this.findingsTotalCount = total;
      this.findingsTopRules = buckets
        .filter((b) => b.key)
        .slice(0, 10)
        .map((b) => ({ key: b.key, count: b.count, pct: (b.count / total) * 100 }));
    });

    // Top users affected by detection findings
    const userReq = new AggregateSignalsRequest();
    const userFilters = this.buildFilters();
    userFilters.setStream('detection');
    userReq.setFilters(userFilters);
    userReq.setGroupBy('user_id');
    userReq.setMetric('count');
    this.grpc.signal.aggregateSignals(userReq, null).then((resp) => {
      const buckets = resp.getBucketsList().map((b) => b.toObject());
      const total = buckets.reduce((s, b) => s + b.count, 0) || 1;
      this.findingsTopUsers = buckets
        .filter((b) => b.key)
        .slice(0, 10)
        .map((b) => ({ key: b.key, count: b.count, pct: (b.count / total) * 100 }));
    });

    // Blocked vs challenged counts
    const blockedReq = new AggregateSignalsRequest();
    const blockedFilters = this.buildFilters();
    blockedFilters.setStream('detection');
    blockedFilters.setOutcome('blocked');
    blockedReq.setFilters(blockedFilters);
    blockedReq.setGroupBy('outcome');
    blockedReq.setMetric('count');
    this.grpc.signal.aggregateSignals(blockedReq, null).then((resp) => {
      this.findingsBlockedCount = resp.getBucketsList().reduce((s, b) => s + b.toObject().count, 0);
    });

    const challengedReq = new AggregateSignalsRequest();
    const challengedFilters = this.buildFilters();
    challengedFilters.setStream('detection');
    challengedFilters.setOutcome('challenged');
    challengedReq.setFilters(challengedFilters);
    challengedReq.setGroupBy('outcome');
    challengedReq.setMetric('count');
    this.grpc.signal.aggregateSignals(challengedReq, null).then((resp) => {
      this.findingsChallengedCount = resp.getBucketsList().reduce((s, b) => s + b.toObject().count, 0);
    });

    // Recent signals with findings (detection stream, most recent)
    const recentQuery = new ListQuery();
    recentQuery.setOffset(0);
    recentQuery.setLimit(20);
    const recentReq = new ListSignalsRequest();
    const recentFilters = this.buildFilters();
    recentFilters.setStream('detection');
    recentReq.setQuery(recentQuery);
    recentReq.setFilters(recentFilters);
    this.grpc.signal.listSignals(recentReq, null).then(
      (resp) => {
        this.findingsRecentSignals = resp.getSignalsList().map((s) => s.toObject());
        this.findingsLoading = false;
      },
      (err) => {
        this.toast.showError(err);
        this.findingsLoading = false;
      },
    );
  }

  buildChartPath(): void {
    if (this.chartBuckets.length === 0) {
      this.chartPath = '';
      this.chartMaxCount = 0;
      return;
    }
    this.chartMaxCount = Math.max(...this.chartBuckets.map((b) => b.count), 1);
    const padding = 8;
    const w = this.chartWidth - padding * 2;
    const h = this.chartHeight - padding * 2;
    const step = w / Math.max(this.chartBuckets.length - 1, 1);
    const points = this.chartBuckets.map((b, i) => {
      const x = padding + i * step;
      const y = padding + h - (b.count / this.chartMaxCount) * h;
      return `${x},${y}`;
    });
    this.chartPath = 'M' + points.join(' L');
  }

  getChartFillPath(): string {
    if (!this.chartPath) return '';
    const padding = 8;
    const h = this.chartHeight - padding;
    return this.chartPath + ` L${this.chartWidth - padding},${h} L${padding},${h} Z`;
  }

  get metricTotal(): number {
    return this.streamCounts.reduce((s, b) => s + b.count, 0);
  }

  get metricFailures(): number {
    return this.outcomeCounts.find((b) => b.key === 'failure')?.count ?? 0;
  }

  get metricSuccessRate(): number {
    const total = this.metricTotal;
    if (total === 0) return 100;
    const success = this.outcomeCounts.find((b) => b.key === 'success')?.count ?? 0;
    return Math.round((success / total) * 1000) / 10;
  }

  get metricUniqueStreams(): number {
    return this.streams.length;
  }

  getDimensionCount(buckets: AggregationBucket.AsObject[], key: string): number {
    return buckets.find((b) => b.key === key)?.count ?? 0;
  }

  findingsCount(signal: Signal.AsObject): number {
    return signal.findingsList?.length ?? 0;
  }

  findingBlocks = (f: Finding.AsObject): boolean => f.block;

  hasBlockFindings(signal: Signal.AsObject): boolean {
    return signal.findingsList?.some((f) => f.block) ?? false;
  }

  nextPage(): void {
    this.offset += this.limit;
    this.search();
  }

  prevPage(): void {
    this.offset = Math.max(0, this.offset - this.limit);
    this.search();
  }

  resetFilters(): void {
    this.filterForm.reset();
    this.offset = 0;
    this.refresh();
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

  trackByKey(_i: number, row: BreakdownRow): string {
    return row.key;
  }

  /** Extract error message from JSON payload, if present. */
  extractError(signal: Signal.AsObject): string | null {
    if (!signal.payload) return null;
    try {
      const obj = JSON.parse(signal.payload);
      return obj?.error || obj?.Error || obj?.message || null;
    } catch {
      return null;
    }
  }

  /** Correlate by trace ID: filter the explorer to show all signals sharing the same OTEL trace. */
  correlate(signal: Signal.AsObject): void {
    if (!signal.traceId) return;
    this.filterForm.reset();
    this.filterForm.patchValue({ trace_id: signal.traceId });
    this.activeTab = 'logs';
    this.offset = 0;
    this.refresh();
    this.search();
  }

  /** Drill from a finding badge to the detection-stream signal(s) sharing the same trace. */
  drillToFinding(signal: Signal.AsObject, event: MouseEvent): void {
    event.stopPropagation();
    if (!signal.traceId) return;
    this.filterForm.reset();
    this.filterForm.patchValue({ trace_id: signal.traceId, stream: 'detection' });
    this.activeTab = 'logs';
    this.offset = 0;
    this.refresh();
    this.search();
  }
}
