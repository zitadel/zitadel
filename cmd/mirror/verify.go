package mirror

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	cryptoDatabase "github.com/zitadel/zitadel/internal/crypto/database"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/query/projection"
)

func verifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "counts if source and dest have the same amount of entries",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				logging.OnError(cmd.Context(), err).Error("zitadel mirror verify command failed")
			}()
			config, shutdown, err := newMigrationConfig(cmd, viper.GetViper())
			if err != nil {
				return err
			}
			defer func() {
				err = errors.Join(err, shutdown(cmd.Context()))
			}()
			verifyMigration(cmd.Context(), config)
			return nil
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

func verifyMigration(ctx context.Context, config *Migration) {
	sourceClient, err := database.Connect(config.Source, false)
	logging.OnError(ctx, err).Fatal("unable to connect to source database")
	defer sourceClient.Close()

	destClient, err := database.Connect(config.Destination, false)
	logging.OnError(ctx, err).Fatal("unable to connect to destination database")
	defer destClient.Close()

	for _, schema := range schemas {
		for _, table := range append(getTables(ctx, destClient, schema), getViews(ctx, destClient, schema)...) {
			sourceCount := countEntries(ctx, sourceClient, table)
			destCount := countEntries(ctx, destClient, table)

			logger := logging.FromCtx(ctx).With("table", table, "dest", destCount, "source", sourceCount)
			if sourceCount == destCount {
				logger.DebugContext(ctx, "equal count")
				continue
			}
			logger.InfoContext(ctx, "unequal count", "diff", destCount-sourceCount)
		}
	}
}

func getTables(ctx context.Context, dest *database.DB, schema string) (tables []string) {
	err := dest.QueryContext(
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
	logging.OnError(ctx, err).Fatal("unable to query tables")
	return tables
}

func getViews(ctx context.Context, dest *database.DB, schema string) (tables []string) {
	err := dest.QueryContext(
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
	logging.OnError(ctx, err).Fatal("unable to query views")
	return tables
}

func countEntries(ctx context.Context, client *database.DB, table string) (count int) {
	instanceClause := instanceClause()
	noInstanceIDColumn := []string{
		projection.InstanceProjectionTable,
		projection.SystemFeatureTable,
		cryptoDatabase.EncryptionKeysTable,
	}
	if slices.Contains(noInstanceIDColumn, table) {
		instanceClause = ""
	}

	err := client.QueryRowContext(
		ctx,
		func(r *sql.Row) error {
			return r.Scan(&count)
		},
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", table, instanceClause),
	)
	logging.OnError(ctx, err).Fatal("unable to count")

	return count
}
