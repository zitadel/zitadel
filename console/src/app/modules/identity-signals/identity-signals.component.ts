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

import type { MessageInitShape } from '@bufbuild/protobuf';
import type { Signal, AggregationBucket } from '@zitadel/proto/zitadel/signal/v2/signal_pb.js';
import { SignalFiltersSchema } from '@zitadel/proto/zitadel/signal/v2/signal_pb.js';

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

type Tab = 'overview' | 'logs';

@Component({
  selector: 'cnsl-identity-signals',
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
  templateUrl: './identity-signals.component.html',
  styleUrls: ['./identity-signals.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('void', style({ height: '0', opacity: '0', overflow: 'hidden' })),
      state('*', style({ height: '*', opacity: '1' })),
      transition('void <=> *', animate('200ms ease-in-out')),
    ]),
  ],
})
export class IdentitySignalsComponent implements OnInit {
  private readonly grpc = inject(GrpcService);
  private readonly fb = inject(FormBuilder);
  private readonly toast = inject(ToastService);

  // Navigation
  activeTab: Tab = 'overview';

  // Loading
  loading = false;

  // List data (logs tab)
  signals: Signal[] = [];
  totalCount = 0;
  offset = 0;
  limit = 50;

  // Chart
  chartBuckets: AggregationBucket[] = [];
  chartLoading = false;
  chartPath = '';
  chartMaxCount = 0;
  chartWidth = 960;
  chartHeight = 160;

  // Summary metrics
  streamCounts: AggregationBucket[] = [];
  outcomeCounts: AggregationBucket[] = [];
  streams: string[] = [];

  // Breakdown aggregations (overview tab)
  topOperations: BreakdownRow[] = [];
  topResources: BreakdownRow[] = [];
  topIPs: BreakdownRow[] = [];
  topCountries: BreakdownRow[] = [];
  topUsers: BreakdownRow[] = [];
  topOrgs: BreakdownRow[] = [];
  topProjects: BreakdownRow[] = [];
  topClients: BreakdownRow[] = [];

  // Expanded row (logs tab)
  expandedSignal: Signal | null = null;

  filterForm: FormGroup = this.fb.group({
    stream: [''],
    outcome: [''],
    operation: [''],
    ip: [''],
    country: [''],
    user_id: [''],
    org_id: [''],
    project_id: [''],
    client_id: [''],
    payload: [''],
    trace_id: [''],
    span_id: [''],
  });

  displayedColumns = ['createdAt', 'stream', 'resource', 'operation', 'outcome', 'ip', 'userId', 'expand'];

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
  }

  switchTab(tab: Tab): void {
    this.activeTab = tab;
    if (tab === 'logs' && this.signals.length === 0) {
      this.search();
    }
  }

  refresh(): void {
    this.loadChart();
    this.loadDimensions();
    this.loadBreakdowns();
    if (this.activeTab === 'logs') {
      this.search();
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

  toggleRow(signal: Signal, event: MouseEvent): void {
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

  private buildFilters(): MessageInitShape<typeof SignalFiltersSchema> {
    const f = this.filterForm.value;
    const filters: Record<string, string> = {};
    if (f.stream) filters['stream'] = f.stream;
    if (f.outcome) filters['outcome'] = f.outcome;
    if (f.operation) filters['operation'] = f.operation;
    if (f.ip) filters['ip'] = f.ip;
    if (f.country) filters['country'] = f.country;
    if (f.user_id) filters['userId'] = f.user_id;
    if (f.org_id) filters['orgId'] = f.org_id;
    if (f.project_id) filters['projectId'] = f.project_id;
    if (f.client_id) filters['clientId'] = f.client_id;
    if (f.payload) filters['payload'] = f.payload;
    if (f.trace_id) filters['traceId'] = f.trace_id;
    if (f.span_id) filters['spanId'] = f.span_id;
    return filters;
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
          this.signals = resp.signals ?? [];
          this.totalCount = Number(resp.details?.totalResult ?? 0);
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
    this.grpc.signal
      .aggregateSignals({
        filters: this.buildFilters(),
        groupBy: 'time_bucket',
        metric: 'count',
        timeBucket: this.selectedTimeRange.bucket,
      })
      .then(
        (resp) => {
          this.chartBuckets = resp.buckets ?? [];
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
    this.grpc.signal
      .aggregateSignals({
        filters: this.buildFilters(),
        groupBy: 'stream',
        metric: 'count',
        timeBucket: '',
      })
      .then((resp) => {
        this.streamCounts = resp.buckets ?? [];
        this.streams = this.streamCounts.map((b) => b.key).filter((k) => k);
      });

    // Outcome counts
    this.grpc.signal
      .aggregateSignals({
        filters: this.buildFilters(),
        groupBy: 'outcome',
        metric: 'count',
        timeBucket: '',
      })
      .then((resp) => {
        this.outcomeCounts = resp.buckets ?? [];
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
      { groupBy: 'org_id', target: 'topOrgs' as const },
      { groupBy: 'project_id', target: 'topProjects' as const },
      { groupBy: 'client_id', target: 'topClients' as const },
    ];
    for (const f of fields) {
      this.grpc.signal
        .aggregateSignals({
          filters: this.buildFilters(),
          groupBy: f.groupBy,
          metric: 'count',
          timeBucket: '',
        })
        .then((resp) => {
          const buckets = resp.buckets ?? [];
          const total = buckets.reduce((s, b) => s + Number(b.count), 0) || 1;
          this[f.target] = buckets
            .filter((b) => b.key)
            .slice(0, 10)
            .map((b) => ({ key: b.key, count: Number(b.count), pct: (Number(b.count) / total) * 100 }));
        });
    }
  }

  buildChartPath(): void {
    if (this.chartBuckets.length === 0) {
      this.chartPath = '';
      this.chartMaxCount = 0;
      return;
    }
    this.chartMaxCount = Math.max(...this.chartBuckets.map((b) => Number(b.count)), 1);
    const padding = 8;
    const w = this.chartWidth - padding * 2;
    const h = this.chartHeight - padding * 2;
    const step = w / Math.max(this.chartBuckets.length - 1, 1);
    const points = this.chartBuckets.map((b, i) => {
      const x = padding + i * step;
      const y = padding + h - (Number(b.count) / this.chartMaxCount) * h;
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
    return this.streamCounts.reduce((s, b) => s + Number(b.count), 0);
  }

  get metricFailures(): number {
    return Number(this.outcomeCounts.find((b) => b.key === 'failure')?.count ?? 0);
  }

  get metricSuccessRate(): number {
    const total = this.metricTotal;
    if (total === 0) return 100;
    const success = Number(this.outcomeCounts.find((b) => b.key === 'success')?.count ?? 0);
    return Math.round((success / total) * 1000) / 10;
  }

  get metricUniqueStreams(): number {
    return this.streams.length;
  }

  getDimensionCount(buckets: AggregationBucket[], key: string): number {
    return Number(buckets.find((b) => b.key === key)?.count ?? 0);
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

  correlate(signal: Signal): void {
    if (!signal.traceId) return;
    this.filterForm.reset();
    this.filterForm.patchValue({ trace_id: signal.traceId });
    this.activeTab = 'logs';
    this.offset = 0;
    this.refresh();
    this.search();
  }
}
