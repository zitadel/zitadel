package eventstore

import (
	"time"
)

type Config struct {
	PushTimeout time.Duration
	MaxRetries  uint32

	Pusher   Pusher
	Querier  Querier
	Searcher Searcher
	Queue    ExecutionQueue

	// Autovacuum tunes PostgreSQL's autovacuum and autoanalyze behavior of the
	// events2 table.
	Autovacuum AutovacuumConfig
}

// AutovacuumConfig configures table-level autovacuum and autoanalyze tuning
// for the eventstore's events2 table.
//
// PostgreSQL's default autovacuum uses a percentage-based scale factor, so as
// events2 grows, ever more rows need to change before a vacuum or analyze is
// triggered. This lets the table's visibility map and the query planner's
// statistics grow stale, which degrades read performance on large instances.
// Enabling this config replaces the scale factors with static, row-count
// based thresholds, so vacuum and analyze keep running at a predictable
// cadence regardless of table size.
type AutovacuumConfig struct {
	// Enabled applies the thresholds below to the events2 table and disables
	// the percentage-based autovacuum scale factors. When disabled, the table
	// is reset to the cluster's default autovacuum settings.
	Enabled bool
	// VacuumInsertThreshold is the number of inserted rows that triggers an insert-only
	// autovacuum run, regardless of table size. It is also applied as autovacuum_vacuum_threshold.
	// Must be greater than 10000 to avoid thrashing the database with constant autovacuum activity.
	VacuumInsertThreshold uint32
	// AnalyzeThreshold is the number of changed rows that triggers an
	// autoanalyze run, regardless of table size. Must be greater than 10000
	// to avoid thrashing the database with constant autoanalyze activity.
	AnalyzeThreshold uint32
}
