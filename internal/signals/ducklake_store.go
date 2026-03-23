//go:build cgo

// DuckLakeStore requires CGO because github.com/duckdb/duckdb-go links the
// DuckDB C library via cgo. Builds with CGO_ENABLED=0 use the stub in
// ducklake_store_nocgo.go instead.

// PREVIEW: Identity Signals is a preview feature. APIs, storage format,
// and configuration may change between releases without notice.

package signals

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	_ "github.com/duckdb/duckdb-go/v2" // DuckDB driver
)

// DuckLakeStore implements [SignalSink] and [SignalReader] backed by DuckLake.
// Catalog metadata is stored in PostgreSQL (via the ducklake extension),
// while signal data is written as Parquet files to the configured data path
// (local filesystem or S3).
type DuckLakeStore struct {
	mu     sync.RWMutex
	db     *sql.DB // DuckDB in-process connection (single-writer)
	pgDSN  string
	dlCfg  DuckLakeConfig
	closed bool
}

// NewDuckLakeStore creates a DuckLake-backed signal store.
// pgDSN is the PostgreSQL libpq key-value connection string used for the
// DuckLake catalog. The postgres+ducklake extensions are installed/loaded
// and the catalog is attached.
func NewDuckLakeStore(pgDSN string, dlCfg DuckLakeConfig) (*DuckLakeStore, error) {
	if pgDSN == "" {
		return nil, fmt.Errorf("ducklake: pgDSN must not be empty")
	}

	db, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, fmt.Errorf("ducklake: open duckdb: %w", err)
	}
	// DuckDB is single-writer; limit the pool to one connection to avoid
	// "database is locked" errors from concurrent writes.
	db.SetMaxOpenConns(1)

	// When ExtensionDirectory is configured, tell DuckDB to use that path for
	// INSTALL/LOAD instead of the default (~/.duckdb/extensions).
	// In container images extensions are pre-downloaded there so no internet
	// access is needed at runtime; INSTALL becomes a no-op.
	if dlCfg.ExtensionDirectory != "" {
		stmt := fmt.Sprintf("SET extension_directory='%s'", escapeSQLString(dlCfg.ExtensionDirectory))
		if _, err := db.Exec(stmt); err != nil {
			db.Close()
			return nil, fmt.Errorf("ducklake: set extension_directory: %w", err)
		}
	}

	installStmts := []string{
		"INSTALL ducklake", "LOAD ducklake",
		"INSTALL postgres", "LOAD postgres",
	}
	for _, stmt := range installStmts {
		if _, err := db.Exec(stmt); err != nil {
			db.Close()
			return nil, fmt.Errorf("ducklake: %s: %w", stmt, err)
		}
	}

	if dlCfg.Backend == ArchiveBackendS3 && dlCfg.S3.Endpoint != "" {
		if err := configureS3(db, dlCfg.S3); err != nil {
			db.Close()
			return nil, fmt.Errorf("ducklake: configure s3: %w", err)
		}
	}

	metadataSchema := dlCfg.MetadataSchema
	if metadataSchema == "" {
		metadataSchema = "signals"
	}
	attachSQL := fmt.Sprintf(
		"ATTACH 'ducklake:postgres:%s' AS signals (DATA_PATH '%s', METADATA_SCHEMA '%s')",
		escapeSQLString(pgDSN), escapeSQLString(dlCfg.DataPath), escapeSQLString(metadataSchema),
	)
	if _, err := db.Exec(attachSQL); err != nil {
		db.Close()
		if strings.Contains(err.Error(), "permission denied for schema") {
			return nil, fmt.Errorf("ducklake: attach catalog: %w\n\n"+
				"The DuckLake extension needs CREATE ON SCHEMA %s.\n"+
				"This schema is created automatically by 'zitadel init'.\n"+
				"If you skipped init, run as superuser on the ZITADEL database:\n"+
				"  CREATE SCHEMA IF NOT EXISTS %s;\n"+
				"  GRANT ALL ON ALL TABLES IN SCHEMA %s TO <your_zitadel_user>;",
				err, metadataSchema, metadataSchema, metadataSchema)
		}
		// Do not log attachSQL — it contains the PostgreSQL DSN with credentials.
		return nil, fmt.Errorf("ducklake: attach catalog failed: %w", err)
	}

	if err := createSignalsTable(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("ducklake: create table: %w", err)
	}

	return &DuckLakeStore{
		db:    db,
		pgDSN: pgDSN,
		dlCfg: dlCfg,
	}, nil
}

func configureS3(db *sql.DB, s3 ArchiveS3Config) error {
	stmts := []string{
		"INSTALL httpfs",
		"LOAD httpfs",
	}
	if s3.AccessKey != "" {
		stmts = append(stmts,
			fmt.Sprintf("SET s3_access_key_id='%s'", escapeSQLString(s3.AccessKey)),
			fmt.Sprintf("SET s3_secret_access_key='%s'", escapeSQLString(s3.SecretKey)),
		)
	}
	if s3.Endpoint != "" {
		stmts = append(stmts,
			fmt.Sprintf("SET s3_endpoint='%s'", escapeSQLString(s3.Endpoint)),
		)
	}
	if !s3.UseSSL {
		stmts = append(stmts, "SET s3_use_ssl=false")
	}
	for i, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			// Never log the statement — it may contain S3 credentials.
			return fmt.Errorf("s3 configuration step %d failed: %w", i, err)
		}
	}
	return nil
}

const createSignalsTableSQL = `
CREATE TABLE IF NOT EXISTS signals.signals (
	instance_id      VARCHAR NOT NULL,
	user_id          VARCHAR NOT NULL,
	caller_id        VARCHAR NOT NULL,
	session_id       VARCHAR NOT NULL,
	fingerprint_id   VARCHAR NOT NULL,
	operation        VARCHAR NOT NULL,
	stream           VARCHAR NOT NULL,
	resource         VARCHAR NOT NULL,
	outcome          VARCHAR NOT NULL,
	created_at       TIMESTAMP NOT NULL,
	ip               VARCHAR NOT NULL,
	user_agent       VARCHAR NOT NULL,
	org_id           VARCHAR NOT NULL DEFAULT '',
	project_id       VARCHAR NOT NULL DEFAULT '',
	client_id        VARCHAR NOT NULL DEFAULT '',
	accept_language  VARCHAR NOT NULL,
	country          VARCHAR NOT NULL,
	forwarded_chain  VARCHAR NOT NULL,
	referer          VARCHAR NOT NULL,
	sec_fetch_site   VARCHAR NOT NULL,
	is_https         BOOLEAN NOT NULL,
	findings         VARCHAR NOT NULL,
	payload          VARCHAR NOT NULL DEFAULT '',
	trace_id         VARCHAR NOT NULL DEFAULT '',
	span_id          VARCHAR NOT NULL DEFAULT '',
	duration_ms      BIGINT NOT NULL DEFAULT 0
)
`

func createSignalsTable(db *sql.DB) error {
	if _, err := db.Exec(createSignalsTableSQL); err != nil {
		return err
	}
	// Idempotent migration: add duration_ms for existing tables.
	_, _ = db.Exec("ALTER TABLE signals.signals ADD COLUMN IF NOT EXISTS duration_ms BIGINT NOT NULL DEFAULT 0")
	return nil
}

const insertSignalSQL = `
INSERT INTO signals.signals (
	instance_id, user_id, caller_id, session_id, fingerprint_id,
	operation, stream, resource, outcome, created_at,
	ip, user_agent, org_id, project_id, client_id,
	accept_language, country, forwarded_chain,
	referer, sec_fetch_site, is_https, findings, payload,
	trace_id, span_id, duration_ms
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

// WriteBatch inserts a batch of signals. Called by the [Emitter] debouncer.
func (s *DuckLakeStore) WriteBatch(ctx context.Context, signals []RecordedSignal) error {
	if len(signals) == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return fmt.Errorf("ducklake: store closed")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ducklake: begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, insertSignalSQL)
	if err != nil {
		return fmt.Errorf("ducklake: prepare: %w", err)
	}
	defer stmt.Close()

	for _, sig := range signals {
		findingsJSON, merr := json.Marshal(sig.Findings)
		if merr != nil || findingsJSON == nil {
			findingsJSON = []byte("[]")
		}
		_, err = stmt.ExecContext(ctx,
			sig.InstanceID,
			sig.UserID,
			sig.CallerID,
			sig.SessionID,
			sig.FingerprintID,
			sig.Operation,
			string(sig.Stream),
			sig.Resource,
			string(sig.Outcome),
			sig.Timestamp.UTC(),
			sig.IP,
			sig.UserAgent,
			sig.OrgID,
			sig.ProjectID,
			sig.ClientID,
			sig.AcceptLanguage,
			sig.Country,
			strings.Join(sig.ForwardedChain, ","),
			sig.Referer,
			sig.SecFetchSite,
			sig.IsHTTPS,
			string(findingsJSON),
			sig.Payload,
			sig.TraceID,
			sig.SpanID,
			sig.DurationMs,
		)
		if err != nil {
			return fmt.Errorf("ducklake: insert: %w", err)
		}
	}

	return tx.Commit()
}

// SearchSignals queries signals with arbitrary filters for the Signals API.
func (s *DuckLakeStore) SearchSignals(ctx context.Context, filters SignalFilters, offset, limit int) ([]RecordedSignal, int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return nil, 0, fmt.Errorf("ducklake: store closed")
	}

	where, args := filtersToSQL(filters)

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM signals.signals WHERE %s", where)
	var total int64
	if err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("ducklake: count: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT instance_id, user_id, caller_id, session_id, fingerprint_id,
		       operation, stream, resource, outcome, created_at,
		       ip, user_agent, org_id, project_id, client_id,
		       accept_language, country, forwarded_chain,
		       referer, sec_fetch_site, is_https, findings, payload,
		       trace_id, span_id, duration_ms
		FROM signals.signals
		WHERE %s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, where)
	args = append(args, limit, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	results, err := scanSignals(rows)
	return results, total, err
}

// AggregateSignals runs an aggregation query over signals.
func (s *DuckLakeStore) AggregateSignals(ctx context.Context, filters SignalFilters, req AggregateRequest) ([]AggregationBucket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return nil, fmt.Errorf("ducklake: store closed")
	}

	where, args := filtersToSQL(filters)

	var groupExpr, selectExpr string
	switch req.GroupBy {
	case "time_bucket":
		interval := req.TimeBucket
		if interval == "" {
			interval = "1 hour"
		}
		if !isAllowedInterval(interval) {
			return nil, fmt.Errorf("ducklake: unsupported time_bucket interval: %q", interval)
		}
		groupExpr = fmt.Sprintf("time_bucket(INTERVAL '%s', created_at)", interval)
		// Format the timestamp as ISO 8601 so the frontend can parse it
		// reliably with new Date(). DuckDB's strftime(ts, format) syntax.
		selectExpr = fmt.Sprintf("strftime(%s, '%%Y-%%m-%%dT%%H:%%M:%%SZ') AS bucket_key", groupExpr)
	default:
		col, err := validateGroupBy(req.GroupBy)
		if err != nil {
			return nil, fmt.Errorf("ducklake: %w", err)
		}
		groupExpr = col
		selectExpr = col + " AS bucket_key"
	}

	metricExpr := "COUNT(*)"
	// isFloatMetric indicates the result is a float (scanned into Value, not Count).
	isFloatMetric := false
	switch req.Metric {
	case "", "count":
		metricExpr = "COUNT(*)"
	case "distinct_count":
		metricExpr = fmt.Sprintf("COUNT(DISTINCT %s)", groupExpr)
	case "avg":
		metricExpr = "AVG(duration_ms)"
		isFloatMetric = true
	case "sum":
		metricExpr = "SUM(duration_ms)"
		isFloatMetric = true
	case "p50":
		metricExpr = "QUANTILE_CONT(duration_ms, 0.5)"
		isFloatMetric = true
	case "p95":
		metricExpr = "QUANTILE_CONT(duration_ms, 0.95)"
		isFloatMetric = true
	case "p99":
		metricExpr = "QUANTILE_CONT(duration_ms, 0.99)"
		isFloatMetric = true
	default:
		return nil, fmt.Errorf("ducklake: unsupported metric: %q", req.Metric)
	}

	// For non-time_bucket groupings, exclude empty/null keys
	emptyFilter := ""
	if req.GroupBy != "time_bucket" {
		if fd := FieldByColumn(req.GroupBy); fd != nil && fd.Filter == FilterBoolean {
			emptyFilter = fmt.Sprintf(" AND %s IS NOT NULL", groupExpr)
		} else {
			emptyFilter = fmt.Sprintf(" AND %s != ''", groupExpr)
		}
	}

	// Multi-dimensional: secondary group-by produces per-series time buckets
	if req.SecondaryGroupBy != "" {
		secondaryCol, err := validateGroupBy(req.SecondaryGroupBy)
		if err != nil {
			return nil, fmt.Errorf("ducklake: unsupported secondary_group_by: %w", err)
		}
		limit := req.Limit
		if limit <= 0 {
			limit = 5
		}
		if limit > 20 {
			limit = 20
		}
		// Two-step approach to avoid double-WHERE parameter issues:
		// 1. Find top N series values
		secondaryNullFilter := fmt.Sprintf("%s IS NOT NULL AND %s != ''", secondaryCol, secondaryCol)
		if fd := FieldByColumn(req.SecondaryGroupBy); fd != nil && fd.Filter == FilterBoolean {
			secondaryNullFilter = fmt.Sprintf("%s IS NOT NULL", secondaryCol)
		}
		topQuery := fmt.Sprintf(`
			SELECT %s AS series_key
			FROM signals.signals
			WHERE %s AND %s
			GROUP BY %s
			ORDER BY COUNT(*) DESC
			LIMIT %d
		`, secondaryCol, where, secondaryNullFilter, secondaryCol, limit)

		topRows, err := s.db.QueryContext(ctx, topQuery, args...)
		if err != nil {
			return nil, fmt.Errorf("ducklake: aggregate multi (top): %w", err)
		}
		var topValues []string
		for topRows.Next() {
			var v string
			if err := topRows.Scan(&v); err != nil {
				topRows.Close()
				return nil, err
			}
			topValues = append(topValues, v)
		}
		topRows.Close()
		if err := topRows.Err(); err != nil {
			return nil, err
		}

		if len(topValues) == 0 {
			return nil, nil
		}

		// 2. Query time buckets for each series value individually and merge
		var allBuckets []AggregationBucket
		for _, sv := range topValues {
			seriesArgs := append(append([]any{}, args...), sv)
			seriesQuery := fmt.Sprintf(`
				SELECT %s, %s AS value
				FROM signals.signals
				WHERE %s%s AND %s = ?
				GROUP BY %s
				ORDER BY bucket_key ASC
				LIMIT 1000
			`, selectExpr, metricExpr, where, emptyFilter, secondaryCol, groupExpr)

			sRows, err := s.db.QueryContext(ctx, seriesQuery, seriesArgs...)
			if err != nil {
				return nil, fmt.Errorf("ducklake: aggregate multi (series %q): %w", sv, err)
			}
			for sRows.Next() {
				var b AggregationBucket
				if isFloatMetric {
					if err := sRows.Scan(&b.Key, &b.Value); err != nil {
						sRows.Close()
						return nil, err
					}
					b.Count = int64(b.Value)
				} else {
					if err := sRows.Scan(&b.Key, &b.Count); err != nil {
						sRows.Close()
						return nil, err
					}
					b.Value = float64(b.Count)
				}
				b.Series = sv
				allBuckets = append(allBuckets, b)
			}
			sRows.Close()
			if err := sRows.Err(); err != nil {
				return nil, err
			}
		}
		return allBuckets, nil
	}

	// Determine sort order: chronological for time_bucket, by count for dimensions
	orderClause := "value DESC"
	if req.GroupBy == "time_bucket" {
		orderClause = "bucket_key ASC"
	}

	query := fmt.Sprintf(`
		SELECT %s, %s AS value
		FROM signals.signals
		WHERE %s%s
		GROUP BY %s
		ORDER BY %s
		LIMIT 100
	`, selectExpr, metricExpr, where, emptyFilter, groupExpr, orderClause)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ducklake: aggregate: %w", err)
	}
	defer rows.Close()

	var buckets []AggregationBucket
	for rows.Next() {
		var b AggregationBucket
		if isFloatMetric {
			if err := rows.Scan(&b.Key, &b.Value); err != nil {
				return nil, err
			}
			b.Count = int64(b.Value)
		} else {
			if err := rows.Scan(&b.Key, &b.Count); err != nil {
				return nil, err
			}
			b.Value = float64(b.Count)
		}
		buckets = append(buckets, b)
	}
	return buckets, rows.Err()
}

// PruneStream deletes signals older than the retention duration for the given stream.
// When instanceID is empty, signals are pruned across all instances — this is
// intentional for the retention worker which runs globally. The instanceID
// parameter is kept for future per-instance retention policies.
// Returns the number of rows deleted.
func (s *DuckLakeStore) PruneStream(ctx context.Context, instanceID string, stream SignalStream, retention time.Duration) (int64, error) {
	// DELETE is a write operation; acquire exclusive lock for DuckDB single-writer.
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return 0, fmt.Errorf("ducklake: store closed")
	}

	cutoff := time.Now().UTC().Add(-retention)
	var clauses []string
	var args []any

	if instanceID != "" {
		clauses = append(clauses, "instance_id = ?")
		args = append(args, instanceID)
	}
	clauses = append(clauses, "stream = ?", "created_at < ?")
	args = append(args, string(stream), cutoff)

	query := fmt.Sprintf("DELETE FROM signals.signals WHERE %s", strings.Join(clauses, " AND "))
	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("ducklake: prune: %w", err)
	}
	return result.RowsAffected()
}

// Compact merges small Parquet files into larger ones. It holds an
// exclusive lock to prevent concurrent reads/writes during the table
// swap. The operation runs inside a single DuckDB transaction so a
// crash mid-compaction cannot lose data.
//
// A 5-minute timeout is enforced to prevent indefinite lock holding
// that would block all query API calls.
func (s *DuckLakeStore) Compact(ctx context.Context) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return 0, fmt.Errorf("ducklake: store closed")
	}

	// Enforce a maximum duration so compaction cannot block reads indefinitely.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	var fileCount int
	// 'signals' refers to the DuckDB catalog alias (always "signals",
	// see ATTACH ... AS signals) and the table name (signals.signals).
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM ducklake_data_files('signals', 'signals')",
	).Scan(&fileCount)
	if err != nil {
		return 0, nil // not critical — skip this cycle
	}
	threshold := s.dlCfg.CompactionThreshold
	if threshold <= 0 {
		threshold = 10
	}
	if fileCount < threshold {
		return 0, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("ducklake: compact begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.ExecContext(ctx, `CREATE OR REPLACE TABLE signals.signals_compacted AS SELECT * FROM signals.signals`); err != nil {
		return 0, fmt.Errorf("compact: create compacted table: %w", err)
	}
	if _, err = tx.ExecContext(ctx, "DROP TABLE IF EXISTS signals.signals"); err != nil {
		return 0, fmt.Errorf("compact: drop original table: %w", err)
	}
	if _, err = tx.ExecContext(ctx, "ALTER TABLE signals.signals_compacted RENAME TO signals"); err != nil {
		return 0, fmt.Errorf("compact: rename compacted table: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("compact: commit: %w", err)
	}
	return fileCount, nil
}

// Close detaches the DuckLake catalog and closes the DuckDB connection.
func (s *DuckLakeStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	_, _ = s.db.Exec("DETACH signals")
	return s.db.Close()
}

// LogInfo logs the store configuration at startup.
func (s *DuckLakeStore) LogInfo(ctx context.Context) {
	slog.InfoContext(ctx, "identity_signals.ducklake.started",
		slog.String("data_path", s.dlCfg.DataPath),
		slog.String("backend", string(s.dlCfg.Backend)),
	)
}

func scanSignals(rows *sql.Rows) ([]RecordedSignal, error) {
	var results []RecordedSignal
	for rows.Next() {
		var (
			rs             RecordedSignal
			stream         string
			outcome        string
			forwardedChain string
			findingsJSON   string
		)
		if err := rows.Scan(
			&rs.InstanceID, &rs.UserID, &rs.CallerID, &rs.SessionID,
			&rs.FingerprintID, &rs.Operation, &stream, &rs.Resource,
			&outcome, &rs.Timestamp, &rs.IP, &rs.UserAgent,
			&rs.OrgID, &rs.ProjectID, &rs.ClientID,
			&rs.AcceptLanguage, &rs.Country, &forwardedChain,
			&rs.Referer, &rs.SecFetchSite, &rs.IsHTTPS, &findingsJSON,
			&rs.Payload, &rs.TraceID, &rs.SpanID, &rs.DurationMs,
		); err != nil {
			return nil, err
		}
		rs.Stream = SignalStream(stream)
		rs.Outcome = Outcome(outcome)
		if forwardedChain != "" {
			rs.ForwardedChain = strings.Split(forwardedChain, ",")
		}
		if findingsJSON != "" && findingsJSON != "[]" {
			_ = json.Unmarshal([]byte(findingsJSON), &rs.Findings)
		}
		results = append(results, rs)
	}
	return results, rows.Err()
}
