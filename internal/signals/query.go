package signals

import (
	"fmt"
	"strings"
	"time"
)

// allowedIntervals is a strict allowlist for time_bucket intervals
// to prevent SQL injection through the interval parameter.
var allowedIntervals = map[string]bool{
	"1 minute":   true,
	"5 minutes":  true,
	"10 minutes": true,
	"15 minutes": true,
	"30 minutes": true,
	"1 hour":     true,
	"3 hours":    true,
	"6 hours":    true,
	"12 hours":   true,
	"1 day":      true,
	"1 week":     true,
	"1 month":    true,
}

func isAllowedInterval(interval string) bool {
	return allowedIntervals[interval]
}

// filtersToSQL builds a WHERE clause from SignalFilters using
// parameterised queries (?-placeholders). The instance_id filter is
// always included as the first clause to enforce tenant isolation.
//
// Field filter behavior is driven by the SignalFields registry:
//   - FilterExact: column = ?
//   - FilterSubstring: column ILIKE '%value%'
//   - FilterTraceCorrelated: (column = ? OR trace_id IN (...))
//
// Unknown field names in Fields are silently ignored.
func filtersToSQL(f SignalFilters) (string, []any) {
	var clauses []string
	var args []any

	clauses = append(clauses, "instance_id = ?")
	args = append(args, f.InstanceID)

	for col, val := range f.Fields {
		if val == "" {
			continue
		}
		fd := FieldByColumn(col)
		if fd == nil {
			continue
		}
		switch fd.Filter {
		case FilterTraceCorrelated:
			c, a := traceCorrelationClause(col, val, f.InstanceID, f.After, f.Before)
			clauses = append(clauses, c)
			args = append(args, a...)
		case FilterSubstring:
			clauses = append(clauses, col+" ILIKE ?")
			args = append(args, "%"+val+"%")
		case FilterExact:
			clauses = append(clauses, col+" = ?")
			args = append(args, val)
		case FilterBoolean:
			clauses = append(clauses, col+" = CAST(? AS BOOLEAN)")
			args = append(args, val)
		}
	}

	if f.After != nil {
		clauses = append(clauses, "created_at >= ?")
		args = append(args, f.After.UTC())
	}
	if f.Before != nil {
		clauses = append(clauses, "created_at < ?")
		args = append(args, f.Before.UTC())
	}

	return strings.Join(clauses, " AND "), args
}

// traceCorrelationClause builds a compound filter that matches a field
// directly OR via trace_id correlation. The subquery finds trace_ids
// associated with the entity and the outer clause includes any signal
// sharing one of those trace_ids.
//
// Time bounds are passed into the subquery to prevent full table scans.
func traceCorrelationClause(field, value, instanceID string, after, before *time.Time) (string, []any) {
	var subClauses []string
	var subArgs []any

	subClauses = append(subClauses, "instance_id = ?")
	subArgs = append(subArgs, instanceID)

	subClauses = append(subClauses, field+" = ?")
	subArgs = append(subArgs, value)

	subClauses = append(subClauses, "trace_id != ''")

	if after != nil {
		subClauses = append(subClauses, "created_at >= ?")
		subArgs = append(subArgs, after.UTC())
	}
	if before != nil {
		subClauses = append(subClauses, "created_at < ?")
		subArgs = append(subArgs, before.UTC())
	}

	subWhere := strings.Join(subClauses, " AND ")
	clause := fmt.Sprintf(
		"(%s = ? OR (trace_id != '' AND trace_id IN ("+
			"SELECT DISTINCT trace_id FROM signals.signals "+
			"WHERE %s"+
			")))",
		field, subWhere,
	)

	// First arg is for the outer "field = ?", rest are subquery args
	args := append([]any{value}, subArgs...)
	return clause, args
}

func escapeSQLString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// validateGroupBy checks that the group_by field is in the registry
// and marked as groupable. Returns the SQL column name or an error.
func validateGroupBy(field string) (string, error) {
	if field == "time_bucket" {
		return "time_bucket", nil
	}
	fd := FieldByColumn(field)
	if fd == nil || !fd.Groupable {
		return "", fmt.Errorf("unsupported group_by field: %q", field)
	}
	return fd.Column, nil
}
