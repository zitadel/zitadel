package signals

// FilterType describes how a field is matched in WHERE clauses.
type FilterType string

const (
	// FilterExact matches column = ?.
	FilterExact FilterType = "exact"
	// FilterSubstring matches column ILIKE '%value%'.
	FilterSubstring FilterType = "substring"
	// FilterTraceCorrelated matches directly OR via trace_id correlation.
	FilterTraceCorrelated FilterType = "trace_correlated"
	// FilterBoolean matches column = CAST(? AS BOOLEAN).
	FilterBoolean FilterType = "boolean"
)

// SignalFieldDef describes a single column in the signals table
// for filtering and aggregation purposes.
type SignalFieldDef struct {
	// Column is the DuckDB column name (e.g. "user_id").
	Column string
	// Label is the canonical user-facing label per TERMINOLOGY.md.
	Label string
	// Filter determines how this field is matched in WHERE clauses.
	Filter FilterType
	// Groupable indicates whether this field can be used in GROUP BY.
	Groupable bool
}

// SignalFields is the authoritative registry of all queryable signal
// columns. Labels follow TERMINOLOGY.md conventions.
var SignalFields = []SignalFieldDef{
	{Column: "user_id", Label: "User", Filter: FilterTraceCorrelated, Groupable: true},
	{Column: "caller_id", Label: "Service Account", Filter: FilterExact, Groupable: true},
	{Column: "session_id", Label: "Session", Filter: FilterTraceCorrelated, Groupable: true},
	{Column: "fingerprint_id", Label: "Device", Filter: FilterExact, Groupable: true},
	{Column: "operation", Label: "Operation", Filter: FilterSubstring, Groupable: true},
	{Column: "stream", Label: "Stream", Filter: FilterExact, Groupable: true},
	{Column: "resource", Label: "Resource", Filter: FilterExact, Groupable: true},
	{Column: "outcome", Label: "Outcome", Filter: FilterExact, Groupable: true},
	{Column: "ip", Label: "IP Address", Filter: FilterExact, Groupable: true},
	{Column: "user_agent", Label: "User Agent", Filter: FilterSubstring, Groupable: true},
	{Column: "org_id", Label: "Organization", Filter: FilterTraceCorrelated, Groupable: true},
	{Column: "project_id", Label: "Project", Filter: FilterExact, Groupable: true},
	{Column: "client_id", Label: "Application", Filter: FilterTraceCorrelated, Groupable: true},
	{Column: "accept_language", Label: "Language", Filter: FilterExact, Groupable: true},
	{Column: "country", Label: "Country", Filter: FilterExact, Groupable: true},
	{Column: "forwarded_chain", Label: "Forwarded Chain", Filter: FilterSubstring, Groupable: false},
	{Column: "referer", Label: "Referer", Filter: FilterSubstring, Groupable: true},
	{Column: "sec_fetch_site", Label: "Fetch Site", Filter: FilterExact, Groupable: true},
	{Column: "is_https", Label: "HTTPS", Filter: FilterBoolean, Groupable: true},
	{Column: "payload", Label: "Payload", Filter: FilterSubstring, Groupable: false},
	{Column: "trace_id", Label: "Trace", Filter: FilterExact, Groupable: false},
	{Column: "span_id", Label: "Span", Filter: FilterExact, Groupable: false},
}

var fieldIndex map[string]*SignalFieldDef

func init() {
	fieldIndex = make(map[string]*SignalFieldDef, len(SignalFields))
	for i := range SignalFields {
		fieldIndex[SignalFields[i].Column] = &SignalFields[i]
	}
}

// FieldByColumn returns the field definition for a column name, or nil.
func FieldByColumn(col string) *SignalFieldDef {
	return fieldIndex[col]
}

// GroupableFields returns a map of column names that can be used in GROUP BY.
// The key and value are both the column name (for backward compatibility
// with the existing validateGroupBy pattern).
func GroupableFields() map[string]string {
	m := make(map[string]string, len(SignalFields))
	for _, f := range SignalFields {
		if f.Groupable {
			m[f.Column] = f.Column
		}
	}
	return m
}
