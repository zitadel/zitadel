package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"slices"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/internal/database"
)

//go:embed sql/cleanup_events2_count.sql
var cleanupEvents2CountSQL string

//go:embed sql/cleanup_events2_delete.sql
var cleanupEvents2DeleteSQL string

type cleanupEvents2Config struct {
	Database database.Config
}

type cleanupSummary struct {
	EventType      string
	RowCount       int64
	AggregateCount int64
}

func cleanupEvents2Cmd() *cobra.Command {
	var (
		execute    bool
		olderThan  time.Duration
		batchSize  int
		maxBatches int
	)

	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "purge terminal auth-flow history from eventstore.events2",
		Long: `Purges terminal auth-flow aggregates from eventstore.events2 once their live state has
already been materialized elsewhere and the configured retention period has elapsed.

This only purges oidc_session aggregates after their current access token and current
refresh token are no longer usable. It intentionally excludes session and session_logout
aggregates because their events still participate in live token or logout validation.

Without --execute, the command performs a dry run and only prints what would be deleted.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if olderThan <= 0 {
				return fmt.Errorf("--older-than must be greater than 0")
			}
			if batchSize <= 0 {
				return fmt.Errorf("--batch-size must be greater than 0")
			}
			if maxBatches < 0 {
				return fmt.Errorf("--max-batches must be 0 or greater")
			}

			cfg, err := readCleanupEvents2Config(viper.GetViper())
			if err != nil {
				return fmt.Errorf("read config: %w", err)
			}

			db, err := database.Connect(cfg.Database, false)
			if err != nil {
				return fmt.Errorf("connect database: %w", err)
			}
			defer db.Close()

			cutoff := time.Now().UTC().Add(-olderThan)
			fmt.Fprintf(cmd.OutOrStdout(), "cutoff=%s execute=%t batch_size=%d max_batches=%d\n", cutoff.Format(time.RFC3339), execute, batchSize, maxBatches)

			if !execute {
				summaries, err := queryCleanupSummaries(cmd.Context(), db, cleanupEvents2CountSQL, cutoff)
				if err != nil {
					return err
				}
				printCleanupSummaries(cmd, "dry-run", summaries)
				return nil
			}

			var (
				allSummaries []cleanupSummary
				batchesRun   int
			)

			for maxBatches == 0 || batchesRun < maxBatches {
				summaries, err := queryCleanupSummaries(cmd.Context(), db, cleanupEvents2DeleteSQL, cutoff, batchSize)
				if err != nil {
					return err
				}
				if totalRows(summaries) == 0 {
					break
				}
				allSummaries = append(allSummaries, summaries...)
				batchesRun++
				printCleanupSummaries(cmd, fmt.Sprintf("batch=%d", batchesRun), summaries)
			}

			if batchesRun == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no eligible rows found")
				return nil
			}

			printCleanupSummaries(cmd, "deleted", mergeCleanupSummaries(allSummaries))
			return nil
		},
	}

	cmd.Flags().BoolVar(&execute, "execute", false, "delete eligible rows instead of performing a dry run")
	cmd.Flags().DurationVar(&olderThan, "older-than", 0, "only delete aggregates whose terminal event is older than this duration")
	cmd.Flags().IntVar(&batchSize, "batch-size", 1000, "maximum number of aggregates to delete per batch")
	cmd.Flags().IntVar(&maxBatches, "max-batches", 0, "maximum number of delete batches to run; 0 means until no eligible rows remain")

	return cmd
}

func readCleanupEvents2Config(v *viper.Viper) (*cleanupEvents2Config, error) {
	cfg := new(cleanupEvents2Config)
	if err := v.Unmarshal(cfg, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		database.DecodeHook(false),
	))); err != nil {
		return nil, err
	}
	return cfg, nil
}

func queryCleanupSummaries(ctx context.Context, db *database.DB, query string, args ...any) ([]cleanupSummary, error) {
	summaries := make([]cleanupSummary, 0)
	err := db.QueryContext(ctx, func(rows *sql.Rows) error {
		for rows.Next() {
			var summary cleanupSummary
			if err := rows.Scan(&summary.EventType, &summary.RowCount, &summary.AggregateCount); err != nil {
				return err
			}
			summaries = append(summaries, summary)
		}
		return rows.Err()
	}, query, args...)
	if err != nil {
		return nil, err
	}
	return summaries, nil
}

func printCleanupSummaries(cmd *cobra.Command, label string, summaries []cleanupSummary) {
	summaries = mergeCleanupSummaries(summaries)
	if len(summaries) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "%s rows=0 aggregates=0\n", label)
		return
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s rows=%d aggregates=%d\n", label, totalRows(summaries), totalAggregates(summaries))
	for _, summary := range summaries {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s rows=%d aggregates=%d\n", summary.EventType, summary.RowCount, summary.AggregateCount)
	}
}

func mergeCleanupSummaries(summaries []cleanupSummary) []cleanupSummary {
	merged := make(map[string]cleanupSummary, len(summaries))
	for _, summary := range summaries {
		current := merged[summary.EventType]
		current.EventType = summary.EventType
		current.RowCount += summary.RowCount
		current.AggregateCount += summary.AggregateCount
		merged[summary.EventType] = current
	}

	result := make([]cleanupSummary, 0, len(merged))
	for _, summary := range merged {
		result = append(result, summary)
	}
	slices.SortFunc(result, func(a, b cleanupSummary) int {
		switch {
		case a.EventType < b.EventType:
			return -1
		case a.EventType > b.EventType:
			return 1
		default:
			return 0
		}
	})
	return result
}

func totalRows(summaries []cleanupSummary) int64 {
	var total int64
	for _, summary := range summaries {
		total += summary.RowCount
	}
	return total
}

func totalAggregates(summaries []cleanupSummary) int64 {
	var total int64
	for _, summary := range summaries {
		total += summary.AggregateCount
	}
	return total
}
