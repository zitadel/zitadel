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

// DuckLakeStore implements [Store] and [SignalSink] backed by DuckLake.
// Catalog metadata is stored in PostgreSQL (via the ducklake extension),
// while signal data is written as Parquet files to the configured data path
// (local filesystem or S3).
type DuckLakeStore struct {
	mu     sync.RWMutex
	db     *sql.DB // DuckDB in-process connection (single-writer)
	cfg    SnapshotConfig
	pgDSN  string
	dlCfg  DuckLakeConfig
	closed bool
}

// NewDuckLakeStore creates a DuckLake-backed signal store.
// pgDSN is the PostgreSQL libpq key-value connection string used for the
// DuckLake catalog. The postgres+ducklake extensions are installed/loaded
// and the catalog is attached.
func NewDuckLakeStore(pgDSN string, cfg SnapshotConfig, dlCfg DuckLakeConfig) (*DuckLakeStore, error) {
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

	// Install ducklake + postgres extensions.
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

	// Configure S3 credentials if needed.
	if dlCfg.Backend == ArchiveBackendS3 && dlCfg.S3.Endpoint != "" {
		if err := configureS3(db, dlCfg.S3); err != nil {
			db.Close()
			return nil, fmt.Errorf("ducklake: configure s3: %w", err)
		}
	}

	// Attach the DuckLake catalog via PostgreSQL.
	attachSQL := fmt.Sprintf(
		"ATTACH 'ducklake:postgres:%s' AS signals (DATA_PATH '%s')",
		pgDSN, dlCfg.DataPath,
	)
	if _, err := db.Exec(attachSQL); err != nil {
		db.Close()
		if strings.Contains(err.Error(), "permission denied for schema") {
			return nil, fmt.Errorf("ducklake: attach catalog: %w\n\nThe DuckLake postgres extension needs CREATE ON SCHEMA public.\n"+
				"Run as a superuser on the ZITADEL database:\n"+
				"  GRANT CREATE ON SCHEMA public TO <your_zitadel_user>;", err)
		}
		return nil, fmt.Errorf("ducklake: attach catalog: %w", err)
	}

	// Create the signals table if it doesn't exist.
	if err := createSignalsTable(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("ducklake: create table: %w", err)
	}

	return &DuckLakeStore{
		db:    db,
		cfg:   cfg,
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
			fmt.Sprintf("SET s3_access_key_id='%s'", s3.AccessKey),
			fmt.Sprintf("SET s3_secret_access_key='%s'", s3.SecretKey),
		)
	}
	if s3.Endpoint != "" {
		stmts = append(stmts,
			fmt.Sprintf("SET s3_endpoint='%s'", s3.Endpoint),
		)
	}
	if !s3.UseSSL {
		stmts = append(stmts, "SET s3_use_ssl=false")
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("%s: %w", stmt, err)
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
	accept_language  VARCHAR NOT NULL,
	country          VARCHAR NOT NULL,
	forwarded_chain  VARCHAR NOT NULL,
	referer          VARCHAR NOT NULL,
	sec_fetch_site   VARCHAR NOT NULL,
	is_https         BOOLEAN NOT NULL,
	findings         VARCHAR NOT NULL,
	payload          VARCHAR NOT NULL DEFAULT '',
	trace_id         VARCHAR NOT NULL DEFAULT '',
	span_id          VARCHAR NOT NULL DEFAULT ''
)
`

func createSignalsTable(db *sql.DB) error {
	_, err := db.Exec(createSignalsTableSQL)
	return err
}

// Save inserts a single signal with findings. Called by the risk engine
// after evaluation.
func (s *DuckLakeStore) Save(ctx context.Context, signal Signal, findings []RecordedFinding, _ SnapshotConfig) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return fmt.Errorf("ducklake: store closed")
	}

	findingsJSON, err := json.Marshal(findings)
	if err != nil {
		findingsJSON = []byte("[]")
	}

	_, err = s.db.ExecContext(ctx, insertSignalSQL,
		signal.InstanceID,
		signal.UserID,
		signal.CallerID,
		signal.SessionID,
		signal.FingerprintID,
		signal.Operation,
		string(signal.Stream),
		signal.Resource,
		string(signal.Outcome),
		signal.Timestamp.UTC(),
		signal.IP,
		signal.UserAgent,
		signal.AcceptLanguage,
		signal.Country,
		strings.Join(signal.ForwardedChain, ","),
		signal.Referer,
		signal.SecFetchSite,
		signal.IsHTTPS,
		string(findingsJSON),
		signal.Payload,
		signal.TraceID,
		signal.SpanID,
	)
	return err
}

const insertSignalSQL = `
INSERT INTO signals.signals (
	instance_id, user_id, caller_id, session_id, fingerprint_id,
	operation, stream, resource, outcome, created_at,
	ip, user_agent, accept_language, country, forwarded_chain,
	referer, sec_fetch_site, is_https, findings, payload,
	trace_id, span_id
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

// WriteBatch inserts a batch of signals. Called by the [Emitter] debouncer.
func (s *DuckLakeStore) WriteBatch(ctx context.Context, signals []RecordedSignal) error {
	if len(signals) == 0 {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
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
		findingsJSON, err := json.Marshal(sig.Findings)
		if err != nil {
			findingsJSON = []byte("[]")
		}
		if findingsJSON == nil {
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
		)
		if err != nil {
			return fmt.Errorf("ducklake: insert: %w", err)
		}
	}

	return tx.Commit()
}

// Snapshot returns recent signals for the user and session associated with
// the given signal. Used by the risk engine to build the RiskContext.
func (s *DuckLakeStore) Snapshot(ctx context.Context, signal Signal, cfg SnapshotConfig) (Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return Snapshot{}, fmt.Errorf("ducklake: store closed")
	}

	cfg = effectiveSnapshotConfig(cfg, s.cfg)
	cutoff := signalCutoff(signal.Timestamp, cfg.HistoryWindow, cfg.ContextChangeWindow)

	var snapshot Snapshot
	var err error

	if signal.UserID != "" {
		snapshot.UserSignals, err = s.querySignals(ctx,
			"instance_id = ? AND user_id = ? AND created_at >= ?",
			cfg.MaxSignalsPerUser,
			signal.InstanceID, signal.UserID, cutoff,
		)
		if err != nil {
			return Snapshot{}, fmt.Errorf("ducklake: user signals: %w", err)
		}
	}
	if signal.SessionID != "" {
		snapshot.SessionSignals, err = s.querySignals(ctx,
			"instance_id = ? AND session_id = ? AND created_at >= ?",
			cfg.MaxSignalsPerSession,
			signal.InstanceID, signal.SessionID, cutoff,
		)
		if err != nil {
			return Snapshot{}, fmt.Errorf("ducklake: session signals: %w", err)
		}
	}
	return snapshot, nil
}

func (s *DuckLakeStore) querySignals(ctx context.Context, where string, limit int, args ...any) ([]RecordedSignal, error) {
	query := fmt.Sprintf(`
		SELECT instance_id, user_id, caller_id, session_id, fingerprint_id,
		       operation, stream, resource, outcome, created_at,
		       ip, user_agent, accept_language, country, forwarded_chain,
		       referer, sec_fetch_site, is_https, findings, payload,
		       trace_id, span_id
		FROM signals.signals
		WHERE %s
		ORDER BY created_at ASC
		LIMIT ?
	`, where)
	args = append(args, limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
			&rs.AcceptLanguage, &rs.Country, &forwardedChain,
			&rs.Referer, &rs.SecFetchSite, &rs.IsHTTPS, &findingsJSON,
			&rs.Payload, &rs.TraceID, &rs.SpanID,
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

// SearchSignals queries signals with arbitrary filters for the Signals API.
// Returns matching signals sorted by created_at descending with pagination.
func (s *DuckLakeStore) SearchSignals(ctx context.Context, filters SignalFilters, offset, limit int) ([]RecordedSignal, int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return nil, 0, fmt.Errorf("ducklake: store closed")
	}

	where, args := filters.toSQL()

	// Count total matches.
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM signals.signals WHERE %s", where)
	var total int64
	if err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("ducklake: count: %w", err)
	}

	// Query the page.
	query := fmt.Sprintf(`
		SELECT instance_id, user_id, caller_id, session_id, fingerprint_id,
		       operation, stream, resource, outcome, created_at,
		       ip, user_agent, accept_language, country, forwarded_chain,
		       referer, sec_fetch_site, is_https, findings, payload,
		       trace_id, span_id
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
			&rs.AcceptLanguage, &rs.Country, &forwardedChain,
			&rs.Referer, &rs.SecFetchSite, &rs.IsHTTPS, &findingsJSON,
			&rs.Payload, &rs.TraceID, &rs.SpanID,
		); err != nil {
			return nil, 0, err
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
	return results, total, rows.Err()
}

// AggregateSignals runs an aggregation query over signals.
func (s *DuckLakeStore) AggregateSignals(ctx context.Context, filters SignalFilters, agg AggregationRequest) ([]AggregationBucket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return nil, fmt.Errorf("ducklake: store closed")
	}

	where, args := filters.toSQL()

	var groupExpr, selectExpr string
	switch agg.GroupBy {
	case AggGroupByTimeBucket:
		interval := agg.TimeBucketInterval
		if interval == "" {
			interval = "1 hour"
		}
		groupExpr = fmt.Sprintf("time_bucket(INTERVAL '%s', created_at)", interval)
		selectExpr = groupExpr + " AS bucket_key"
	case AggGroupByField:
		groupExpr = agg.FieldName
		selectExpr = agg.FieldName + " AS bucket_key"
	default:
		return nil, fmt.Errorf("ducklake: unknown aggregation group_by: %s", agg.GroupBy)
	}

	var metricExpr string
	switch agg.Metric {
	case AggMetricCount:
		metricExpr = "COUNT(*)"
	case AggMetricDistinctCount:
		if agg.DistinctField == "" {
			return nil, fmt.Errorf("ducklake: distinct_field required for distinct_count metric")
		}
		metricExpr = fmt.Sprintf("COUNT(DISTINCT %s)", agg.DistinctField)
	default:
		metricExpr = "COUNT(*)"
	}

	query := fmt.Sprintf(`
		SELECT %s, %s AS value
		FROM signals.signals
		WHERE %s
		GROUP BY %s
		ORDER BY value DESC
		LIMIT ?
	`, selectExpr, metricExpr, where, groupExpr)

	topN := agg.TopN
	if topN <= 0 {
		topN = 100
	}
	args = append(args, topN)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ducklake: aggregate: %w", err)
	}
	defer rows.Close()

	var buckets []AggregationBucket
	for rows.Next() {
		var b AggregationBucket
		if err := rows.Scan(&b.Key, &b.Value); err != nil {
			return nil, err
		}
		buckets = append(buckets, b)
	}
	return buckets, rows.Err()
}

// FindingFilters defines query filters for finding-level searches.
type FindingFilters struct {
	SignalFilters        // embed signal-level filters for instance isolation and time range
	FindingName   string // exact match on finding name (e.g. "failure_burst")
	FindingSource string // exact match on finding source (e.g. "rule:_builtin_failure_burst")
	BlockOnly     bool   // only return blocking findings
	ChallengeOnly bool   // only return challenge findings
}

// FindingResult is a finding with its originating signal context.
type FindingResult struct {
	RecordedFinding
	// Signal context for correlation.
	SignalTimestamp time.Time
	UserID         string
	SessionID      string
	IP             string
	Operation      string
	Stream         SignalStream
	Outcome        Outcome
	TraceID        string
}

// SearchFindings queries findings across signals by unnesting the JSON
// findings column. Results are individual findings with signal context
// attached for correlation.
func (s *DuckLakeStore) SearchFindings(ctx context.Context, filters FindingFilters, offset, limit int) ([]FindingResult, int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return nil, 0, fmt.Errorf("ducklake: store closed")
	}

	where, args := filters.SignalFilters.toSQL()

	// Only include signals that have findings.
	where += " AND findings != '[]' AND findings != ''"

	// Base CTE that unnests findings from the JSON column.
	cte := fmt.Sprintf(`
		WITH finding_rows AS (
			SELECT
				s.created_at AS signal_timestamp,
				s.user_id,
				s.session_id,
				s.ip,
				s.operation,
				s.stream,
				s.outcome,
				s.trace_id,
				f.name,
				f.source,
				f.message,
				COALESCE(f.block, false) AS block,
				COALESCE(f.confidence, 0.0) AS confidence,
				COALESCE(f.challenge, false) AS challenge,
				COALESCE(f.challenge_type, '') AS challenge_type
			FROM signals.signals s,
				LATERAL (
					SELECT
						json_extract_string(j, '$.name') AS name,
						json_extract_string(j, '$.source') AS source,
						json_extract_string(j, '$.message') AS message,
						CAST(json_extract(j, '$.block') AS BOOLEAN) AS block,
						CAST(json_extract(j, '$.confidence') AS DOUBLE) AS confidence,
						CAST(json_extract(j, '$.challenge') AS BOOLEAN) AS challenge,
						json_extract_string(j, '$.challenge_type') AS challenge_type
					FROM unnest(from_json(s.findings, '["json"]')) AS t(j)
				) f
			WHERE %s
		)`, where)

	// Build finding-level filters.
	var findingClauses []string
	if filters.FindingName != "" {
		findingClauses = append(findingClauses, "name = ?")
		args = append(args, filters.FindingName)
	}
	if filters.FindingSource != "" {
		findingClauses = append(findingClauses, "source = ?")
		args = append(args, filters.FindingSource)
	}
	if filters.BlockOnly {
		findingClauses = append(findingClauses, "block = true")
	}
	if filters.ChallengeOnly {
		findingClauses = append(findingClauses, "challenge = true")
	}

	findingWhere := "1=1"
	if len(findingClauses) > 0 {
		findingWhere = strings.Join(findingClauses, " AND ")
	}

	// Count total.
	countQuery := fmt.Sprintf("%s SELECT COUNT(*) FROM finding_rows WHERE %s", cte, findingWhere)
	var total int64
	if err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("ducklake: count findings: %w", err)
	}

	// Query page.
	query := fmt.Sprintf(`%s
		SELECT signal_timestamp, user_id, session_id, ip, operation,
		       stream, outcome, trace_id,
		       name, source, message, block, confidence, challenge, challenge_type
		FROM finding_rows
		WHERE %s
		ORDER BY signal_timestamp DESC
		LIMIT ? OFFSET ?
	`, cte, findingWhere)
	pageArgs := make([]any, len(args))
	copy(pageArgs, args)
	pageArgs = append(pageArgs, limit, offset)

	rows, err := s.db.QueryContext(ctx, query, pageArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("ducklake: search findings: %w", err)
	}
	defer rows.Close()

	var results []FindingResult
	for rows.Next() {
		var (
			fr     FindingResult
			stream string
			outcome string
		)
		if err := rows.Scan(
			&fr.SignalTimestamp, &fr.UserID, &fr.SessionID, &fr.IP,
			&fr.Operation, &stream, &outcome, &fr.TraceID,
			&fr.Name, &fr.Source, &fr.Message, &fr.Block,
			&fr.Confidence, &fr.Challenge, &fr.ChallengeType,
		); err != nil {
			return nil, 0, err
		}
		fr.Stream = SignalStream(stream)
		fr.Outcome = Outcome(outcome)
		results = append(results, fr)
	}
	return results, total, rows.Err()
}

// AggregateFindings runs an aggregation query over unnested findings.
// GroupBy can be "name", "source", or "block" to group findings by those attributes.
func (s *DuckLakeStore) AggregateFindings(ctx context.Context, filters FindingFilters, groupBy string, topN int) ([]AggregationBucket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return nil, fmt.Errorf("ducklake: store closed")
	}

	where, args := filters.SignalFilters.toSQL()
	where += " AND findings != '[]' AND findings != ''"

	// Validate groupBy to prevent SQL injection.
	var groupExpr string
	switch groupBy {
	case "name":
		groupExpr = "f_name"
	case "source":
		groupExpr = "f_source"
	case "block":
		groupExpr = "CAST(f_block AS VARCHAR)"
	case "user_id":
		groupExpr = "s.user_id"
	case "session_id":
		groupExpr = "s.session_id"
	case "outcome":
		groupExpr = "s.outcome"
	default:
		return nil, fmt.Errorf("ducklake: unsupported finding group_by: %q", groupBy)
	}

	// Build finding-level filter clauses.
	var findingClauses []string
	if filters.FindingName != "" {
		findingClauses = append(findingClauses, "f_name = ?")
		args = append(args, filters.FindingName)
	}
	if filters.FindingSource != "" {
		findingClauses = append(findingClauses, "f_source = ?")
		args = append(args, filters.FindingSource)
	}
	if filters.BlockOnly {
		findingClauses = append(findingClauses, "f_block = true")
	}
	if filters.ChallengeOnly {
		findingClauses = append(findingClauses, "f_challenge = true")
	}

	findingWhere := ""
	if len(findingClauses) > 0 {
		findingWhere = " AND " + strings.Join(findingClauses, " AND ")
	}

	if topN <= 0 {
		topN = 100
	}

	query := fmt.Sprintf(`
		SELECT %s AS bucket_key, COUNT(*) AS value
		FROM signals.signals s,
			LATERAL (
				SELECT
					json_extract_string(j, '$.name') AS f_name,
					json_extract_string(j, '$.source') AS f_source,
					CAST(json_extract(j, '$.block') AS BOOLEAN) AS f_block,
					CAST(json_extract(j, '$.challenge') AS BOOLEAN) AS f_challenge
				FROM unnest(from_json(s.findings, '["json"]')) AS t(j)
			) f
		WHERE %s%s
		GROUP BY %s
		ORDER BY value DESC
		LIMIT ?
	`, groupExpr, where, findingWhere, groupExpr)
	args = append(args, topN)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ducklake: aggregate findings: %w", err)
	}
	defer rows.Close()

	var buckets []AggregationBucket
	for rows.Next() {
		var b AggregationBucket
		if err := rows.Scan(&b.Key, &b.Value); err != nil {
			return nil, err
		}
		buckets = append(buckets, b)
	}
	return buckets, rows.Err()
}

// AppendFindings merges additional findings into a signal row identified by
// instance_id + session_id + created_at. This is used by the async LLM path
// in observe mode: the signal is persisted before the model responds, and
// findings are appended after classification completes.
func (s *DuckLakeStore) AppendFindings(ctx context.Context, instanceID, sessionID string, createdAt time.Time, findings []RecordedFinding) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return fmt.Errorf("ducklake: store closed")
	}

	findingsJSON, err := json.Marshal(findings)
	if err != nil {
		return fmt.Errorf("ducklake: marshal findings: %w", err)
	}

	_, err = s.db.ExecContext(ctx, appendFindingsSQL,
		string(findingsJSON),
		instanceID,
		sessionID,
		createdAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("ducklake: append findings: %w", err)
	}
	return nil
}

const appendFindingsSQL = `
UPDATE signals.signals
SET findings = ?
WHERE instance_id = ? AND session_id = ? AND created_at = ?
  AND findings = '[]'
LIMIT 1
`

// DB returns the underlying DuckDB connection for advanced queries.
func (s *DuckLakeStore) DB() *sql.DB {
	return s.db
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

// SignalFilters defines query filters for the Signals API.
type SignalFilters struct {
	InstanceID string
	UserID     string
	SessionID  string
	IP         string
	Operation  string
	Stream     string
	Outcome    string
	Country    string
	Resource   string
	Payload    string
	TraceID    string
	SpanID     string
	After      *time.Time
	Before     *time.Time
}

func (f SignalFilters) toSQL() (string, []any) {
	var clauses []string
	var args []any

	// Instance isolation is always required.
	clauses = append(clauses, "instance_id = ?")
	args = append(args, f.InstanceID)

	if f.UserID != "" {
		clauses = append(clauses, "user_id = ?")
		args = append(args, f.UserID)
	}
	if f.SessionID != "" {
		clauses = append(clauses, "session_id = ?")
		args = append(args, f.SessionID)
	}
	if f.IP != "" {
		clauses = append(clauses, "ip = ?")
		args = append(args, f.IP)
	}
	if f.Operation != "" {
		clauses = append(clauses, "operation ILIKE ?")
		args = append(args, "%"+f.Operation+"%")
	}
	if f.Stream != "" {
		clauses = append(clauses, "stream = ?")
		args = append(args, f.Stream)
	}
	if f.Outcome != "" {
		clauses = append(clauses, "outcome = ?")
		args = append(args, f.Outcome)
	}
	if f.Country != "" {
		clauses = append(clauses, "country = ?")
		args = append(args, f.Country)
	}
	if f.Resource != "" {
		clauses = append(clauses, "resource = ?")
		args = append(args, f.Resource)
	}
	if f.Payload != "" {
		clauses = append(clauses, "payload ILIKE ?")
		args = append(args, "%"+f.Payload+"%")
	}
	if f.TraceID != "" {
		clauses = append(clauses, "trace_id = ?")
		args = append(args, f.TraceID)
	}
	if f.SpanID != "" {
		clauses = append(clauses, "span_id = ?")
		args = append(args, f.SpanID)
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

// AggregationRequest defines what aggregation to perform.
type AggregationRequest struct {
	GroupBy            AggGroupBy
	FieldName          string // for AggGroupByField
	TimeBucketInterval string // for AggGroupByTimeBucket, e.g. "1 hour", "1 day"
	Metric             AggMetric
	DistinctField      string // for AggMetricDistinctCount
	TopN               int    // max buckets returned (default: 100)
}

// AggGroupBy specifies the grouping dimension.
type AggGroupBy string

const (
	AggGroupByField      AggGroupBy = "field"
	AggGroupByTimeBucket AggGroupBy = "time_bucket"
)

// AggMetric specifies the metric to compute per bucket.
type AggMetric string

const (
	AggMetricCount         AggMetric = "count"
	AggMetricDistinctCount AggMetric = "distinct_count"
)

// AggregationBucket is a single result row from an aggregation query.
type AggregationBucket struct {
	Key   string
	Value int64
}

// LogInfo logs a summary of the store configuration.
func (s *DuckLakeStore) LogInfo(ctx context.Context) {
	slog.InfoContext(ctx, "risk.signal_store.ducklake.started",
		slog.String("data_path", s.dlCfg.DataPath),
		slog.String("backend", string(s.dlCfg.Backend)),
	)
}
