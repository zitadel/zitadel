/**
 * Signal field registry — single source of truth for all queryable
 * signal columns. Labels follow TERMINOLOGY.md conventions.
 *
 * The backend has an identical registry in internal/signals/fields.go.
 * Keep both in sync when adding new fields.
 */

export type SignalFilterType = 'exact' | 'substring' | 'trace_correlated' | 'boolean';

export interface SignalFieldDef {
  /** DuckDB column name (e.g. "user_id") */
  key: string;
  /** Canonical user-facing label per TERMINOLOGY.md */
  label: string;
  /** How this field is matched in WHERE clauses */
  filterType: SignalFilterType;
  /** Whether this field can be used in GROUP BY aggregations */
  groupable: boolean;
  /** Whether autocomplete suggestions should be loaded for this field */
  suggestable: boolean;
  /** Material icon name for entity-type usage (optional) */
  icon?: string;
}

/**
 * All queryable signal fields. Order determines display order in
 * filter dropdowns and breakdown tabs.
 */
export const SIGNAL_FIELDS: SignalFieldDef[] = [
  { key: 'user_id', label: 'User', filterType: 'trace_correlated', groupable: true, suggestable: true, icon: 'person' },
  { key: 'caller_id', label: 'Service Account', filterType: 'exact', groupable: true, suggestable: true, icon: 'smart_toy' },
  { key: 'session_id', label: 'Session', filterType: 'trace_correlated', groupable: true, suggestable: true, icon: 'key' },
  { key: 'fingerprint_id', label: 'Device', filterType: 'exact', groupable: true, suggestable: true, icon: 'fingerprint' },
  { key: 'operation', label: 'Operation', filterType: 'substring', groupable: true, suggestable: true, icon: 'api' },
  { key: 'stream', label: 'Stream', filterType: 'exact', groupable: true, suggestable: false },
  { key: 'resource', label: 'Resource', filterType: 'exact', groupable: true, suggestable: true },
  { key: 'outcome', label: 'Outcome', filterType: 'exact', groupable: true, suggestable: false },
  { key: 'ip', label: 'IP Address', filterType: 'exact', groupable: true, suggestable: true, icon: 'language' },
  { key: 'user_agent', label: 'User Agent', filterType: 'substring', groupable: true, suggestable: true },
  { key: 'org_id', label: 'Organization', filterType: 'trace_correlated', groupable: true, suggestable: true, icon: 'business' },
  { key: 'project_id', label: 'Project', filterType: 'exact', groupable: true, suggestable: true },
  { key: 'client_id', label: 'Application', filterType: 'trace_correlated', groupable: true, suggestable: true, icon: 'devices' },
  { key: 'accept_language', label: 'Language', filterType: 'exact', groupable: true, suggestable: true },
  { key: 'country', label: 'Country', filterType: 'exact', groupable: true, suggestable: true },
  { key: 'forwarded_chain', label: 'Forwarded Chain', filterType: 'substring', groupable: false, suggestable: false },
  { key: 'referer', label: 'Referer', filterType: 'substring', groupable: true, suggestable: true },
  { key: 'sec_fetch_site', label: 'Fetch Site', filterType: 'exact', groupable: true, suggestable: false },
  { key: 'is_https', label: 'HTTPS', filterType: 'boolean', groupable: true, suggestable: false },
  { key: 'payload', label: 'Payload', filterType: 'substring', groupable: false, suggestable: false },
  { key: 'trace_id', label: 'Trace', filterType: 'exact', groupable: false, suggestable: false, icon: 'link' },
  { key: 'span_id', label: 'Span', filterType: 'exact', groupable: false, suggestable: false },
];

/** Map of column → field definition for O(1) lookup */
const fieldIndex = new Map<string, SignalFieldDef>();
for (const f of SIGNAL_FIELDS) {
  fieldIndex.set(f.key, f);
}

/** Get field definition by column name */
export function fieldByKey(key: string): SignalFieldDef | undefined {
  return fieldIndex.get(key);
}

/** Get canonical label for a column */
export function fieldLabel(key: string): string {
  return fieldIndex.get(key)?.label ?? key;
}

/** All fields that can appear as filter chips (excludes stream/outcome which are dropdowns) */
export function filterableFields(): SignalFieldDef[] {
  return SIGNAL_FIELDS.filter(f => f.key !== 'stream' && f.key !== 'outcome');
}

/** All fields that can be used as GROUP BY dimensions */
export function groupableFields(): SignalFieldDef[] {
  return SIGNAL_FIELDS.filter(f => f.groupable);
}

/** All fields that support autocomplete suggestions */
export function suggestableFieldKeys(): Set<string> {
  return new Set(SIGNAL_FIELDS.filter(f => f.suggestable).map(f => f.key));
}

/** Fields suitable as entity types in the activity view (have an icon) */
export function entityTypeFields(): SignalFieldDef[] {
  return SIGNAL_FIELDS.filter(f => !!f.icon);
}

/** Build a filter label map (key→label) for backward compat */
export function filterLabelMap(): Record<string, string> {
  const m: Record<string, string> = {};
  for (const f of SIGNAL_FIELDS) {
    m[f.key] = f.label;
  }
  return m;
}

/**
 * Convert a snake_case DuckDB column name to the camelCase proto
 * field name used by connect-es (e.g. "user_id" → "userId").
 */
function snakeToCamel(s: string): string {
  return s.replace(/_([a-z])/g, (_, c) => c.toUpperCase());
}

/**
 * Build a proto SignalFilters object from form values.
 * Form keys are snake_case (matching DuckDB columns);
 * proto fields are camelCase (generated by bufbuild).
 *
 * @param formValues - FormGroup.value object
 * @param excludeField - optional field to skip (for suggestion queries)
 */
export function buildProtoFilters(
  formValues: Record<string, string>,
  excludeField?: string,
): Record<string, string> {
  const filters: Record<string, string> = {};
  for (const [key, val] of Object.entries(formValues)) {
    if (!val || key === excludeField) continue;
    filters[snakeToCamel(key)] = val;
  }
  return filters;
}
