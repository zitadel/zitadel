import { CommonModule } from '@angular/common';
import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { GrpcService } from 'src/app/services/grpc.service';
import { ToastService } from 'src/app/services/toast.service';
import { groupableFields, suggestableFieldKeys, filterableFields, buildProtoFilters, filterLabelMap } from './signal-fields';

import type { MessageInitShape } from '@bufbuild/protobuf';
import type { AggregationBucket } from '@zitadel/proto/zitadel/signal/v2/signal_pb.js';
import { SignalFiltersSchema } from '@zitadel/proto/zitadel/signal/v2/signal_pb.js';

interface TimeRange {
  label: string;
  value: string;
  bucket: string; // default resolution for this range
}

interface Resolution {
  label: string;
  value: string;
  minutes: number;
}

interface BreakdownRow {
  key: string;
  count: number;
  pct: number;
}

interface FilterChip {
  key: string;
  label: string;
  value: string;
}

interface ChartSeries {
  key: string;
  color: string;
  path: string;
  fillPath: string;
  bars: { x: number; y: number; w: number; h: number; idx: number }[];
}

@Component({
  selector: 'cnsl-signals-query',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    TranslateModule,
    MatButtonModule,
    MatIconModule,
    MatMenuModule,
    MatProgressSpinnerModule,
    MatTooltipModule,
  ],
  templateUrl: './signals-query.component.html',
  styleUrls: ['./signals.component.scss'],
})
export class SignalsQueryComponent implements OnInit, OnDestroy {
  private readonly grpc = inject(GrpcService);
  private readonly fb = inject(FormBuilder);
  private readonly toast = inject(ToastService);
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);

  private alive = true;

  signalsAvailable = true;

  // X-axis timeline labels
  xAxisLabels: { text: string; pct: number }[] = [];
  // Y-axis tick values
  yAxisTicks: number[] = [];
  // Y-axis unit label (e.g. "ms", "count")
  yAxisLabel = '';
  // Scaled max used for chart rendering (matches top y-axis tick)
  chartScaleMax = 1;

  // Chart
  chartBuckets: AggregationBucket[] = [];
  chartLoading = false;
  chartPath = '';
  chartMaxCount = 0;
  chartWidth = 960;
  chartHeight = 200;
  chartType: 'line' | 'bar' = 'line';
  chartBars: { x: number; y: number; w: number; h: number; idx: number }[] = [];

  // Multi-series chart
  chartSeries: ChartSeries[] = [];
  readonly seriesColors = ['#6366f1', '#22c55e', '#f59e0b', '#ef4444', '#06b6d4'];

  // Resolution options — each has minutes equivalent for filtering
  readonly resolutions: Resolution[] = [
    { label: '1 min', value: '1 minute', minutes: 1 },
    { label: '5 min', value: '5 minutes', minutes: 5 },
    { label: '15 min', value: '15 minutes', minutes: 15 },
    { label: '30 min', value: '30 minutes', minutes: 30 },
    { label: '1 hour', value: '1 hour', minutes: 60 },
    { label: '6 hours', value: '6 hours', minutes: 360 },
    { label: '12 hours', value: '12 hours', minutes: 720 },
    { label: '1 day', value: '1 day', minutes: 1440 },
  ];
  selectedResolution: Resolution | null = null; // null = auto (from time range)

  // Map time range values to total minutes
  private readonly timeRangeMinutes: Record<string, number> = {
    '1 hour': 60,
    '6 hours': 360,
    '24 hours': 1440,
    '7 days': 10080,
    '30 days': 43200,
  };

  // Map bucket labels to milliseconds
  private readonly bucketMs: Record<string, number> = {
    '1 minute': 60_000,
    '5 minutes': 300_000,
    '15 minutes': 900_000,
    '30 minutes': 1_800_000,
    '1 hour': 3_600_000,
    '3 hours': 10_800_000,
    '6 hours': 21_600_000,
    '12 hours': 43_200_000,
  };

  // Boundaries of the currently displayed time window (used to build full x-axis grid)
  private chartWindowStart = 0;
  private chartWindowEnd = 0;

  get availableResolutions(): Resolution[] {
    const rangeMinutes = this.timeRangeMinutes[this.selectedTimeRange.value] ?? 1440;
    // Only show resolutions where the range has at least 3 buckets
    return this.resolutions.filter((r) => rangeMinutes / r.minutes >= 3);
  }

  // Summary
  streamCounts: AggregationBucket[] = [];
  outcomeCounts: AggregationBucket[] = [];
  streams: string[] = [];

  // Data source & metric selectors
  selectedSource = '';
  selectedMetric = 'count';

  private readonly floatMetrics = new Set(['avg', 'sum', 'p50', 'p95', 'p99']);

  get isFloatMetric(): boolean {
    return this.floatMetrics.has(this.selectedMetric);
  }

  /** Extract the numeric value from a bucket, using `value` for float metrics. */
  bv(b: AggregationBucket): number {
    return this.isFloatMetric ? Number(b.value ?? b.count) : Number(b.count);
  }

  readonly availableSources = [
    { key: 'requests', label: 'Requests' },
    { key: 'events', label: 'Events' },
  ];

  readonly metrics = [
    { key: 'count', label: 'Count' },
    { key: 'distinct_count', label: 'Unique Users' },
    { key: 'avg', label: 'Avg Duration (ms)' },
    { key: 'sum', label: 'Total Duration (ms)' },
    { key: 'p50', label: 'p50 Latency (ms)' },
    { key: 'p95', label: 'p95 Latency (ms)' },
    { key: 'p99', label: 'p99 Latency (ms)' },
  ];

  // Group-by + breakdown
  activeGroupBys: string[] = ['operation'];
  dimensionSearch = '';
  breakdownSearch = '';
  breakdownPage = 0;
  primaryBreakdown: BreakdownRow[] = [];
  dimensionCounts: Record<string, number> = {};

  readonly dimensions = groupableFields().map(f => ({ key: f.key, label: f.label }));

  // Filters
  filterForm: FormGroup = this.fb.group({
    stream: [''],
    outcome: [''],
    ...Object.fromEntries(filterableFields().map(f => [f.key, ['']])),
  });

  pendingFilterKey = '';
  pendingFilterLabel = '';
  filterSuggestions: { key: string; count: number }[] = [];
  filterSuggestionsLoading = false;
  filterInputValue = '';

  private readonly suggestableFields = new Set(suggestableFieldKeys());

  timeRanges: TimeRange[] = [
    { label: 'Last 1h', value: '1 hour', bucket: '5 minutes' },
    { label: 'Last 6h', value: '6 hours', bucket: '5 minutes' },
    { label: 'Last 24h', value: '24 hours', bucket: '30 minutes' },
    { label: 'Last 7d', value: '7 days', bucket: '3 hours' },
    { label: 'Last 30d', value: '30 days', bucket: '12 hours' },
  ];
  selectedTimeRange: TimeRange = this.timeRanges[1];

  ngOnInit(): void {
    this.restoreFromUrl();
    this.refresh();
  }

  ngOnDestroy(): void {
    this.alive = false;
  }

  refresh(): void {
    this.syncUrl();
    this.loadChart();
    this.loadDimensions();
    this.loadBreakdowns();
  }

  // --- Selectors ---

  get selectedSourceLabel(): string {
    if (!this.selectedSource) return 'All Signals';
    return this.availableSources.find((s) => s.key === this.selectedSource)?.label ?? this.selectedSource;
  }

  get selectedMetricLabel(): string {
    return this.metrics.find((m) => m.key === this.selectedMetric)?.label ?? this.selectedMetric;
  }

  get selectedGroupBy(): string {
    return this.activeGroupBys[0] || 'operation';
  }

  get selectedGroupLabel(): string {
    return this.dimensionLabel(this.selectedGroupBy);
  }

  selectSource(key: string): void {
    this.selectedSource = key;
    this.filterForm.patchValue({ stream: key });
    this.refresh();
  }

  selectMetric(key: string): void {
    this.selectedMetric = key;
    this.refresh();
  }

  selectTimeRange(range: TimeRange): void {
    this.selectedTimeRange = range;
    this.selectedResolution = null; // reset to auto
    this.refresh();
  }

  get activeResolution(): Resolution {
    if (this.selectedResolution) return this.selectedResolution;
    return (
      this.resolutions.find((r) => r.value === this.selectedTimeRange.bucket) ?? {
        label: this.selectedTimeRange.bucket,
        value: this.selectedTimeRange.bucket,
        minutes: 0,
      }
    );
  }

  get activeResolutionLabel(): string {
    return this.activeResolution.label;
  }

  selectResolution(res: Resolution): void {
    this.selectedResolution = res;
    this.loadChart();
  }

  // --- Group-by ---

  dimensionLabel(key: string): string {
    return this.dimensions.find((d) => d.key === key)?.label ?? key;
  }

  filteredDimensions() {
    const search = this.dimensionSearch.toLowerCase();
    return this.dimensions
      .filter((d) => !this.activeGroupBys.includes(d.key))
      .filter((d) => !search || d.label.toLowerCase().includes(search) || d.key.includes(search));
  }

  addGroupBy(key: string): void {
    if (!this.activeGroupBys.includes(key)) {
      this.activeGroupBys = [key];
      this.breakdownPage = 0;
      this.loadBreakdowns();
      this.loadChart();
    }
  }

  removeGroupBy(key: string): void {
    this.activeGroupBys = this.activeGroupBys.filter((k) => k !== key);
    this.primaryBreakdown = [];
    this.loadChart();
  }

  // --- Filters ---

  private readonly _filterLabelMap = filterLabelMap();

  activeFilterChips(): FilterChip[] {
    const chips: FilterChip[] = [];
    const vals = this.filterForm.value;
    for (const [key, val] of Object.entries(vals)) {
      if (val && key !== 'stream') {
        chips.push({ key, label: this._filterLabelMap[key] ?? key, value: val as string });
      }
    }
    return chips;
  }

  openFilterInput(key: string): void {
    this.pendingFilterKey = key;
    this.pendingFilterLabel = this.dimensions.find((d) => d.key === key)?.label ?? key;
    this.filterSuggestions = [];
    this.filterInputValue = '';
    if (this.suggestableFields.has(key)) {
      this.loadFilterSuggestions(key);
    }
  }

  applyPendingFilter(value: string): void {
    if (value && this.pendingFilterKey) {
      this.filterForm.patchValue({ [this.pendingFilterKey]: value.trim() });
      this.pendingFilterKey = '';
      this.pendingFilterLabel = '';
      this.filterSuggestions = [];
      this.filterInputValue = '';
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
      this.refresh();
    }
  }

  clearFilter(key: string): void {
    this.filterForm.patchValue({ [key]: '' });
    this.refresh();
  }

  // --- Breakdown search & pagination ---

  filteredBreakdown(): BreakdownRow[] {
    let rows = this.primaryBreakdown;
    if (this.breakdownSearch) {
      const s = this.breakdownSearch.toLowerCase();
      rows = rows.filter((r) => r.key.toLowerCase().includes(s));
    }
    const pageSize = 10;
    const start = this.breakdownPage * pageSize;
    return rows.slice(start, start + pageSize);
  }

  get breakdownTotalPages(): number {
    return Math.ceil(this.primaryBreakdown.length / 10);
  }

  breakdownPrev(): void {
    if (this.breakdownPage > 0) this.breakdownPage--;
  }

  breakdownNext(): void {
    if (this.breakdownPage < this.breakdownTotalPages - 1) this.breakdownPage++;
  }

  // --- Navigation ---

  private readonly validDrillFields = new Set([
    'operation',
    'ip',
    'country',
    'user_id',
    'org_id',
    'project_id',
    'client_id',
    'resource',
    'user_agent',
    'referer',
    'stream',
    'outcome',
  ]);

  drillDownToLogs(field: string, value: string): void {
    if (!this.validDrillFields.has(field)) return;
    this.router.navigate(['/signals/logs'], { queryParams: { [field]: value } });
  }

  /** Returns the series color for a breakdown row based on its index. */
  breakdownColor(idx: number): string {
    return this.seriesColors[idx % this.seriesColors.length];
  }

  trackByKey(_i: number, row: BreakdownRow): string {
    return row.key;
  }

  // --- Data loading ---

  private restoreFromUrl(): void {
    const p = this.route.snapshot.queryParams;
    // Source / metric
    if (p['source'] && this.availableSources.some((s) => s.key === p['source'])) {
      this.selectedSource = p['source'];
    }
    if (p['metric'] && this.metrics.some((m) => m.key === p['metric'])) {
      this.selectedMetric = p['metric'];
    }
    // Group by
    if (p['groupBy'] && this.dimensions.some((d) => d.key === p['groupBy'])) {
      this.activeGroupBys = [p['groupBy']];
    }
    // Time range
    if (p['time']) {
      const tr = this.timeRanges.find((r) => r.value === p['time']);
      if (tr) this.selectedTimeRange = tr;
    }
    // Filters
    const patchable: Record<string, string> = {};
    for (const key of Object.keys(this.filterForm.controls)) {
      if (p[key]) patchable[key] = p[key];
    }
    if (this.selectedSource) patchable['stream'] = this.selectedSource;
    if (Object.keys(patchable).length) {
      this.filterForm.patchValue(patchable);
    }
  }

  private syncUrl(): void {
    const params: Record<string, string> = {};
    if (this.selectedSource) params['source'] = this.selectedSource;
    if (this.selectedMetric !== 'count') params['metric'] = this.selectedMetric;
    if (this.activeGroupBys.length && this.activeGroupBys[0] !== 'operation') params['groupBy'] = this.activeGroupBys[0];
    if (this.selectedTimeRange !== this.timeRanges[1]) params['time'] = this.selectedTimeRange.value;
    const f = this.filterForm.value;
    for (const [key, val] of Object.entries(f)) {
      if (val && key !== 'stream') params[key] = val as string;
    }
    this.router.navigate([], { queryParams: params, queryParamsHandling: 'replace', replaceUrl: true });
  }

  private buildFilters(excludeField?: string): MessageInitShape<typeof SignalFiltersSchema> {
    return buildProtoFilters(this.filterForm.value, excludeField);
  }

  loadChart(): void {
    if (!this.grpc.signal) return;
    this.chartLoading = true;
    const secondaryGroupBy = this.activeGroupBys.length > 0 ? this.activeGroupBys[0] : '';
    this.grpc.signal
      .aggregateSignals({
        filters: this.buildFilters(),
        groupBy: 'time_bucket',
        metric: this.selectedMetric,
        timeBucket: this.activeResolution.value,
        secondaryGroupBy,
        limit: 5,
      })
      .then(
        (resp) => {
          if (!this.alive) return;
          this.chartBuckets = this.fillTimeGrid(resp.buckets ?? []);
          this.buildChartPath();
          this.chartLoading = false;
        },
        (err) => {
          this.chartLoading = false;
          this.handleApiError(err);
        },
      );
  }

  /**
   * Expands sparse API results into a complete time-series grid covering the
   * full selected time window. Buckets missing from the response are inserted
   * with count = 0 so the chart always spans the full window without gaps.
   */
  private fillTimeGrid(buckets: AggregationBucket[]): AggregationBucket[] {
    const bucketInterval = this.bucketMs[this.activeResolution.value] ?? 300_000;
    const rangeMs = (this.timeRangeMinutes[this.selectedTimeRange.value] ?? 360) * 60_000;
    const now = Date.now();
    // Align end to the nearest bucket boundary so the grid is stable across refreshes
    this.chartWindowEnd = Math.floor(now / bucketInterval) * bucketInterval;
    this.chartWindowStart = this.chartWindowEnd - rangeMs;

    // Build lookup from ISO key → buckets (may be multi-series)
    const hasSeries = buckets.some((b) => b.series);
    const byKey = new Map<string, AggregationBucket[]>();
    for (const b of buckets) {
      const list = byKey.get(b.key) ?? [];
      list.push(b);
      byKey.set(b.key, list);
    }

    // Discover the set of series present (for multi-series fill)
    const seriesSet = hasSeries ? [...new Set(buckets.map((b) => b.series).filter(Boolean))] : [''];

    const result: AggregationBucket[] = [];
    for (let t = this.chartWindowStart; t <= this.chartWindowEnd; t += bucketInterval) {
      const isoKey = new Date(t).toISOString().replace(/\.\d{3}Z$/, 'Z');
      if (hasSeries) {
        for (const s of seriesSet) {
          const existing = byKey.get(isoKey)?.find((b) => b.series === s);
          result.push(existing ?? ({ key: isoKey, count: BigInt(0), value: 0, series: s ?? '' } as AggregationBucket));
        }
      } else {
        const existing = byKey.get(isoKey)?.[0];
        result.push(existing ?? ({ key: isoKey, count: BigInt(0), value: 0, series: '' } as AggregationBucket));
      }
    }
    return result;
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
    for (const dim of this.dimensions) {
      this.grpc.signal
        .aggregateSignals({ filters: this.buildFilters(), groupBy: dim.key, metric: 'count', timeBucket: '' })
        .then(
          (resp) => {
            if (!this.alive) return;
            const buckets = (resp.buckets ?? []).filter((b) => b.key);
            this.dimensionCounts[dim.key] = buckets.length;
          },
          (err) => this.handleApiError(err),
        );
    }
  }

  loadBreakdowns(): void {
    if (!this.grpc.signal || this.activeGroupBys.length === 0) {
      this.primaryBreakdown = [];
      return;
    }
    const groupBy = this.activeGroupBys[0];
    this.grpc.signal
      .aggregateSignals({ filters: this.buildFilters(), groupBy, metric: this.selectedMetric, timeBucket: '' })
      .then(
        (resp) => {
          if (!this.alive) return;
          const buckets = resp.buckets ?? [];
          const total = buckets.reduce((s, b) => s + this.bv(b), 0) || 1;
          this.primaryBreakdown = buckets
            .filter((b) => b.key)
            .map((b) => ({ key: b.key, count: this.bv(b), pct: (this.bv(b) / total) * 100 }));
          this.breakdownPage = 0;
        },
        (err) => this.handleApiError(err),
      );
  }

  buildChartPath(): void {
    if (this.chartBuckets.length === 0) {
      this.chartPath = '';
      this.chartMaxCount = 0;
      this.chartScaleMax = 1;
      this.chartBars = [];
      this.chartSeries = [];
      this.xAxisLabels = [];
      this.yAxisTicks = [];
      this.yAxisLabel = '';
      return;
    }

    // Compute raw max from data
    this.chartMaxCount = Math.max(...this.chartBuckets.map((b) => this.bv(b)), 1);

    // Build y-axis ticks FIRST so the chart scale aligns with the labels
    this.buildYAxisTicks();

    // Check if we have multi-series data (buckets with non-empty series field)
    const hasSeries = this.chartBuckets.some((b) => b.series);

    if (hasSeries) {
      this.buildMultiSeriesChart();
    } else {
      this.buildSingleSeriesChart();
    }
    this.buildXAxisLabels();
  }

  private buildSingleSeriesChart(): void {
    this.chartSeries = [];
    const padding = 8;
    const w = this.chartWidth - padding * 2;
    const h = this.chartHeight - padding * 2;
    const max = this.chartScaleMax;
    const step = w / Math.max(this.chartBuckets.length - 1, 1);
    const points = this.chartBuckets.map((b, i) => {
      const x = padding + i * step;
      const y = padding + h - (this.bv(b) / max) * h;
      return `${x},${y}`;
    });
    this.chartPath = 'M' + points.join(' L');

    const barGap = 1;
    const barW = Math.max(w / this.chartBuckets.length - barGap, 1);
    this.chartBars = this.chartBuckets.map((b, i) => {
      const val = this.bv(b);
      const barH = (val / max) * h;
      return {
        x: padding + i * (barW + barGap),
        y: padding + h - barH,
        w: barW,
        h: barH,
        idx: i,
      };
    });
  }

  private buildMultiSeriesChart(): void {
    this.chartPath = '';
    this.chartBars = [];

    // Group buckets by series key
    const seriesMap = new Map<string, AggregationBucket[]>();
    for (const b of this.chartBuckets) {
      const key = b.series || '(other)';
      if (!seriesMap.has(key)) seriesMap.set(key, []);
      seriesMap.get(key)!.push(b);
    }

    // Collect all unique time keys — sort chronologically
    const allKeys = [...new Set(this.chartBuckets.map((b) => b.key))];
    allKeys.sort((a, b) => {
      const ta = new Date(a).getTime();
      const tb = new Date(b).getTime();
      if (!isNaN(ta) && !isNaN(tb)) return ta - tb;
      return a < b ? -1 : a > b ? 1 : 0;
    });

    if (allKeys.length === 0) {
      this.chartSeries = [];
      return;
    }

    // Find global max across all series
    this.chartMaxCount = Math.max(...this.chartBuckets.map((b) => this.bv(b)), 1);

    const padding = 8;
    const w = this.chartWidth - padding * 2;
    const h = this.chartHeight - padding * 2;
    const max = this.chartScaleMax;
    const seriesCount = seriesMap.size;

    this.chartSeries = [];
    let seriesIdx = 0;
    for (const [key, buckets] of seriesMap) {
      const color = this.seriesColors[seriesIdx % this.seriesColors.length];

      // Map this series' buckets by time key for fast lookup
      const byKey = new Map(buckets.map((b) => [b.key, this.bv(b)]));

      // Line path
      const step = w / Math.max(allKeys.length - 1, 1);
      const points = allKeys.map((tk, i) => {
        const val = byKey.get(tk) ?? 0;
        const x = padding + i * step;
        const y = padding + h - (val / max) * h;
        return `${x},${y}`;
      });
      const path = 'M' + points.join(' L');
      const fillPath = path + ` L${padding + (allKeys.length - 1) * step},${padding + h} L${padding},${padding + h} Z`;

      // Bar data — subdivide each time slot by series
      const barGap = 1;
      const slotW = Math.max(w / allKeys.length - barGap, 1);
      const subBarW = Math.max(slotW / seriesCount - 1, 1);
      const bars = allKeys.map((tk, i) => {
        const val = byKey.get(tk) ?? 0;
        const barH = (val / max) * h;
        return {
          x: padding + i * (slotW + barGap) + seriesIdx * (subBarW + 1),
          y: padding + h - barH,
          w: subBarW,
          h: barH,
          idx: i,
        };
      });

      this.chartSeries.push({ key, color, path, fillPath, bars });
      seriesIdx++;
    }
  }

  trackByBar(_i: number, bar: { idx: number }): number {
    return bar.idx;
  }

  getChartFillPath(): string {
    if (!this.chartPath || this.chartSeries.length > 0) return '';
    const padding = 8;
    const h = this.chartHeight - padding;
    return this.chartPath + ` L${this.chartWidth - padding},${h} L${padding},${h} Z`;
  }

  get isMultiSeries(): boolean {
    return this.chartSeries.length > 0;
  }

  get metricTotal(): number {
    return this.streamCounts.reduce((s, b) => s + Number(b.count), 0);
  }

  toNumber(val: bigint | number | string): number {
    return Number(val);
  }

  getDimensionCount(buckets: AggregationBucket[], key: string): number {
    return Number(buckets.find((b) => b.key === key)?.count ?? 0);
  }

  /** Shorten a fully-qualified operation name like /zitadel.user.v2.UserService/ListUsers → ListUsers */
  shortName(name: string): string {
    if (!name) return '';
    if (/^\d{1,3}(\.\d{1,3}){3}$/.test(name) || name.includes(':')) return name;
    const slashParts = name.split('/');
    if (slashParts.length >= 3) return slashParts[slashParts.length - 1];
    const dotParts = name.split('.');
    if (dotParts.length > 3) return dotParts.slice(-3).join('.');
    return name;
  }

  private buildXAxisLabels(): void {
    if (this.chartWindowStart === 0 || this.chartBuckets.length < 2) {
      this.xAxisLabels = [];
      return;
    }
    const windowMs = this.chartWindowEnd - this.chartWindowStart;
    const showDate = windowMs > 24 * 60 * 60 * 1000;

    // Pick ~6 evenly spaced labels across the full window
    const targetLabels = 6;
    const stepMs = windowMs / (targetLabels - 1);
    this.xAxisLabels = [];
    for (let i = 0; i < targetLabels; i++) {
      const t = this.chartWindowStart + i * stepMs;
      const d = new Date(t);
      const pct = (i / (targetLabels - 1)) * 100;
      const text = showDate
        ? d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' }) +
          ' ' +
          d.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' })
        : d.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' });
      this.xAxisLabels.push({ text, pct });
    }
  }

  private buildYAxisTicks(): void {
    if (this.chartMaxCount <= 0) {
      this.yAxisTicks = [];
      this.chartScaleMax = 1;
      this.yAxisLabel = '';
      return;
    }
    // Set y-axis unit label based on metric
    const metricUnits: Record<string, string> = {
      count: 'count',
      distinct_count: 'users',
      avg: 'ms',
      sum: 'ms',
      p50: 'ms',
      p95: 'ms',
      p99: 'ms',
    };
    this.yAxisLabel = metricUnits[this.selectedMetric] ?? '';

    // Nice round ticks — compute a ceiling max so bars align with labels
    const max = this.chartMaxCount;
    const tickCount = 4;
    const rawStep = max / tickCount;
    const magnitude = Math.pow(10, Math.floor(Math.log10(rawStep)));
    const niceStep = Math.ceil(rawStep / magnitude) * magnitude;
    const niceMax = niceStep * tickCount;
    // Use the nice max as the chart scale so bars/lines match the y-axis
    this.chartScaleMax = niceMax;
    this.yAxisTicks = [];
    for (let v = niceMax; v >= 0; v -= niceStep) {
      this.yAxisTicks.push(v);
    }
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
