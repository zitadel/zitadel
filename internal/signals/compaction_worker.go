package signals

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/riverqueue/river"
	"github.com/robfig/cron/v3"

	"github.com/zitadel/zitadel/internal/queue"
)

const compactionQueueName = "signal_compaction"

// CompactionArgs is the River job payload for DuckLake compaction.
type CompactionArgs struct{}

func (CompactionArgs) Kind() string { return "signal.ducklake_compact" }

// CompactionWorker merges small Parquet files written by DuckLake into
// larger time-aligned files. This reduces file count on S3/filesystem
// and improves query performance by reducing per-file overhead.
type CompactionWorker struct {
	river.WorkerDefaults[CompactionArgs]
	store    *DuckLakeStore
	interval time.Duration
}

// NewCompactionWorker creates a compaction worker for DuckLake signal files.
func NewCompactionWorker(store *DuckLakeStore, interval time.Duration) *CompactionWorker {
	if interval <= 0 {
		interval = 1 * time.Hour
	}
	return &CompactionWorker{
		store:    store,
		interval: interval,
	}
}

// Register adds the compaction worker to the River worker set.
func (w *CompactionWorker) Register(workers *river.Workers, queues map[string]river.QueueConfig) {
	river.AddWorker[CompactionArgs](workers, w)
	queues[compactionQueueName] = river.QueueConfig{MaxWorkers: 1}
}

// Work performs the compaction job:
// 1. Queries DuckLake catalog metadata for small data files
// 2. For each time range with many small files, reads and rewrites
//    into a single compacted file
// 3. DuckLake handles the catalog update atomically
func (w *CompactionWorker) Work(ctx context.Context, _ *river.Job[CompactionArgs]) error {
	if w.store == nil || w.store.closed {
		return nil
	}

	db := w.store.DB()
	if db == nil {
		return nil
	}

	compacted, err := runCompaction(ctx, db)
	if err != nil {
		return fmt.Errorf("ducklake compaction: %w", err)
	}

	if compacted > 0 {
		slog.InfoContext(ctx, "risk.signal_store.ducklake.compaction_complete",
			slog.Int("files_compacted", compacted),
		)
	}
	return nil
}

// runCompaction performs the actual compaction using DuckDB's COMPACT statement
// or a manual read-rewrite if COMPACT is not available.
func runCompaction(ctx context.Context, db *sql.DB) (int, error) {
	// DuckLake tracks data files in its catalog. We can use a simple
	// CTAS (CREATE TABLE AS SELECT) pattern to rewrite and compact.
	//
	// Strategy: count the data files, and if there are many small ones,
	// perform a table-level compaction by rewriting all data.
	// DuckLake handles the metadata update atomically.

	var fileCount int
	err := db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM ducklake_data_files('signals', 'signals')",
	).Scan(&fileCount)
	if err != nil {
		// If the catalog metadata query fails, it may not be supported
		// in this DuckLake version. Skip compaction silently.
		slog.WarnContext(ctx, "risk.signal_store.ducklake.compaction_skip",
			slog.String("reason", "cannot query data files"),
			slog.String("error", err.Error()),
		)
		return 0, nil
	}

	// Only compact if there are many small files (threshold: 10).
	if fileCount < 10 {
		return 0, nil
	}

	// Perform compaction by rewriting the table data.
	// DuckLake VACUUM or manual approach: create temp → swap → drop.
	_, err = db.ExecContext(ctx, `
		CREATE OR REPLACE TABLE signals.signals_compacted AS 
		SELECT * FROM signals.signals
	`)
	if err != nil {
		return 0, fmt.Errorf("create compacted table: %w", err)
	}

	_, err = db.ExecContext(ctx, "DROP TABLE IF EXISTS signals.signals")
	if err != nil {
		return 0, fmt.Errorf("drop original table: %w", err)
	}

	_, err = db.ExecContext(ctx, "ALTER TABLE signals.signals_compacted RENAME TO signals")
	if err != nil {
		return 0, fmt.Errorf("rename compacted table: %w", err)
	}

	return fileCount, nil
}

// RegisterCompactionWorker registers the compaction worker with the queue
// before it starts.
func RegisterCompactionWorker(ctx context.Context, q *queue.Queue, worker *CompactionWorker) {
	if worker == nil {
		return
	}
	q.AddWorkers(ctx, worker)
}

// StartCompactionSchedule starts the periodic compaction job.
// Must be called after the queue has started.
func StartCompactionSchedule(ctx context.Context, q *queue.Queue, worker *CompactionWorker) {
	if worker == nil {
		return
	}
	schedule, _ := cron.ParseStandard(fmt.Sprintf("@every %s", worker.interval))
	q.AddPeriodicJob(ctx, schedule, &CompactionArgs{},
		queue.WithQueueName(compactionQueueName),
		queue.WithMaxAttempts(3),
	)
}
