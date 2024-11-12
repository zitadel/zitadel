package mirror

import (
	"bufio"
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"

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

func verifyMigration(ctx context.Context, config *Migration) {
	if isSrcFile || isDestFile {
		verifyFileMigration(ctx, config)
		return
	}
	verifyMigrationDB(ctx, config)
}

func verifyMigrationDB(ctx context.Context, config *Migration) {
	sourceClient, err := database.Connect(config.Source, false, dialect.DBPurposeQuery)
	logging.OnError(err).Fatal("unable to connect to source database")
	defer sourceClient.Close()

	destClient, err := database.Connect(config.Destination, false, dialect.DBPurposeEventPusher)
	logging.OnError(err).Fatal("unable to connect to destination database")
	defer destClient.Close()

	for _, schema := range schemas {
		for _, table := range append(getTables(ctx, destClient, schema), getViews(ctx, destClient, schema)...) {
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
}

func verifyFileMigration(ctx context.Context, config *Migration) {
	var skipCount int64
	var equalCount int64

	if isDestFile {
		sourceClient, err := database.Connect(config.Source, false, dialect.DBPurposeQuery)
		logging.OnError(err).Fatal("unable to connect to source database")
		defer sourceClient.Close()

		for _, schema := range schemas {
			for _, table := range append(getTables(ctx, sourceClient, schema), getViews(ctx, sourceClient, schema)...) {
				destCount := countEntriesFromFile(table)
				if destCount == 0 {
					skipCount++
					continue
				}

				sourceCount := countEntries(ctx, sourceClient, table)

				entry := logging.WithFields("table", table, "dest", destCount, "source", sourceCount)
				if sourceCount == destCount {
					equalCount++
					entry.Debug("equal count")
					continue
				}
				entry.WithField("diff", destCount-sourceCount).Info("unequal count")
			}
		}

		logging.WithFields("skip", skipCount, "equal", equalCount).Info("verification done for database to file migration")
	}

	if isSrcFile {
		destClient, err := database.Connect(config.Destination, false, dialect.DBPurposeQuery)
		logging.OnError(err).Fatal("unable to connect to destination database")
		defer destClient.Close()

		for _, schema := range schemas {
			for _, table := range append(getTables(ctx, destClient, schema), getViews(ctx, destClient, schema)...) {
				sourceCount := countEntriesFromFile(table)
				if sourceCount == 0 {
					skipCount++
					continue
				}

				destCount := countEntries(ctx, destClient, table)

				entry := logging.WithFields("table", table, "dest", destCount, "source", sourceCount)
				if sourceCount == destCount {
					equalCount++
					entry.Debug("equal count")
					continue
				}
				entry.WithField("diff", destCount-sourceCount).Info("unequal count")
			}
		}

		logging.WithFields("skip", skipCount, "equal", equalCount).Info("verification done for file to database migration")
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
	logging.WithFields("schema", schema).OnError(err).Fatal("unable to query tables")
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
	logging.WithFields("schema", schema).OnError(err).Fatal("unable to query views")
	return tables
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

func countEntriesFromFile(table string) (count int) {
	fileName := fmt.Sprintf("/%s.csv", table)

	if _, err := os.Stat(filePath + fileName); os.IsNotExist(err) {
		return count
	}

	srcFile, err := os.OpenFile(filePath+fileName, os.O_RDONLY, 0)
	logging.OnError(err).Fatal("unable to open file")
	defer srcFile.Close()

	scanner := bufio.NewScanner(srcFile)
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		logging.OnError(err).Fatal("unable to read file")
	}

	return count
}
