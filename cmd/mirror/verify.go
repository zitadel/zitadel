package mirror

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
)

func verifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "counts if source and dest have the same amount of entries",
		Run: func(cmd *cobra.Command, args []string) {
			config := mustNewMigrationConfig(viper.GetViper())
			verifyMigration(cmd.Context(), config)
		},
	}
}

var schemas = []string{
	"adminapi",
	"auth",
	"eventstore",
	"projections",
	"system",
}

func verifyMigration(ctx context.Context, config *Migration) error {
	sourceClient, err := database.Connect(config.Source, false, dialect.DBPurposeQuery)
	if err != nil {
		return fmt.Errorf("unable to connect to source database: %w", err)
	}
	defer sourceClient.Close()

	destClient, err := database.Connect(config.Destination, false, dialect.DBPurposeEventPusher)
	if err != nil {
		return fmt.Errorf("unable to connect to destination database: %w", err)
	}
	defer destClient.Close()

	for _, schema := range schemas {
		var tables []string
		tables, err = getTables(ctx, destClient, schema)
		if err != nil {
			return fmt.Errorf("unable to get tables: %w", err)
		}
		var views []string
		views, err = getViews(ctx, destClient, schema)
		if err != nil {
			return fmt.Errorf("unable to get views: %w", err)
		}
		for _, table := range append(tables, views...) {
			sourceCount := countEntries(ctx, sourceClient, table)
			destCount := countEntries(ctx, destClient, table)

			entry := logging.WithFields("table", table, "dest", destCount, "source", sourceCount)
			if sourceCount == destCount {
				entry.Debug("equal count")
				continue
			}
			entry.WithField("diff", destCount-sourceCount).Info("unequal count")
		}
	}
	return nil
}

func getTables(ctx context.Context, dest *database.DB, schema string) (tables []string, err error) {
	err = dest.QueryContext(
		ctx,
		func(r *sql.Rows) error {
			for r.Next() {
				var table string
				if err := r.Scan(&table); err != nil {
					return err
				}
				tables = append(tables, table)
			}
			return r.Err()
		},
		"SELECT CONCAT(schemaname, '.', tablename) FROM pg_tables WHERE schemaname = $1",
		schema,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to query tables: %w", err)
	}
	return tables, nil
}

func getViews(ctx context.Context, dest *database.DB, schema string) (tables []string, err error) {
	err = dest.QueryContext(
		ctx,
		func(r *sql.Rows) error {
			for r.Next() {
				var table string
				if err := r.Scan(&table); err != nil {
					return err
				}
				tables = append(tables, table)
			}
			return r.Err()
		},
		"SELECT CONCAT(schemaname, '.', viewname) FROM pg_views WHERE schemaname = $1",
		schema,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to query views in schema %s: %w", schema, err)
	}
	return tables, nil
}

func countEntries(ctx context.Context, client *database.DB, table string) (count int) {
	err := client.QueryRowContext(
		ctx,
		func(r *sql.Row) error {
			return r.Scan(&count)
		},
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", table, instanceClause()),
	)
	logging.WithFields("table", table, "db", client.DatabaseName()).OnError(err).Error("unable to count")

	return count
}
