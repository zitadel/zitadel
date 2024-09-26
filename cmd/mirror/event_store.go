package mirror

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	db "github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/v2/database"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/eventstore/postgres"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var shouldIgnorePrevious bool

func eventstoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "eventstore",
		Short: "mirrors the eventstore of an instance from one database to another",
		Long: `mirrors the eventstore of an instance from one database to another
ZITADEL needs to be initialized and set up with the --for-mirror flag
Migrate only copies events2 and unique constraints`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := mustNewMigrationConfig(viper.GetViper())
			return copyEventstore(cmd.Context(), config)
		},
	}

	cmd.Flags().BoolVar(&shouldReplace, "replace", false, "allow delete unique constraints of defined instances before copy")
	cmd.Flags().BoolVar(&shouldIgnorePrevious, "ignore-previous", false, "ignores previous migrations of the events table")

	return cmd
}

func copyEventstore(ctx context.Context, config *Migration) error {
	sourceClient, err := db.Connect(config.Source, false, dialect.DBPurposeEventPusher)
	if err != nil {
		return fmt.Errorf("unable to connect to source database: %w", err)
	}
	defer sourceClient.Close()

	destClient, err := db.Connect(config.Destination, false, dialect.DBPurposeEventPusher)
	if err != nil {
		return fmt.Errorf("unable to connect to destination database: %w", err)
	}
	defer destClient.Close()

	if err = copyEvents(ctx, sourceClient, destClient, config.EventBulkSize); err != nil {
		return fmt.Errorf("unable to copy events: %w", err)
	}
	if err = copyUniqueConstraints(ctx, sourceClient, destClient); err != nil {
		return fmt.Errorf("unable to copy unique constraints: %w", err)
	}
	return nil
}

func positionQuery(db *db.DB) string {
	switch db.Type() {
	case "postgres":
		return "SELECT EXTRACT(EPOCH FROM clock_timestamp())"
	case "cockroach":
		return "SELECT cluster_logical_timestamp()"
	default:
		logging.WithFields("db_type", db.Type()).Fatal("database type not recognized")
		return ""
	}
}

func copyEvents(ctx context.Context, source, dest *db.DB, bulkSize uint32) error {
	start := time.Now()
	reader, writer := io.Pipe()
	migrationID, err := id.SonyFlakeGenerator().Next()
	if err != nil {
		return fmt.Errorf("unable to generate migration id: %w", err)
	}
	sourceConn, err := source.Conn(ctx)
	if err != nil {
		return fmt.Errorf("unable to acquire source connection: %w", err)
	}
	destConn, err := dest.Conn(ctx)
	if err != nil {
		return fmt.Errorf("unable to acquire dest connection: %w", err)
	}
	sourceES := eventstore.NewEventstoreFromOne(postgres.New(source, &postgres.Config{
		MaxRetries: 3,
	}))
	destinationES := eventstore.NewEventstoreFromOne(postgres.New(dest, &postgres.Config{
		MaxRetries: 3,
	}))

	previousMigration, err := queryLastSuccessfulMigration(ctx, destinationES, source.DatabaseName())
	if err != nil {
		return fmt.Errorf("unable to query latest successful migration: %w", err)
	}
	maxPosition, err := writeMigrationStart(ctx, sourceES, migrationID, dest.DatabaseName())
	if err != nil {
		return fmt.Errorf("unable to write migration start: %w", err)

	}
	logging.WithFields("from", previousMigration.Position, "to", maxPosition).Info("start event migration")

	nextPos := make(chan bool, 1)
	pos := make(chan float64, 1)
	errs := make(chan error, 3)

	go func() {
		goErr := sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()
			nextPos <- true
			var i uint32
			for position := range pos {
				var stmt database.Statement
				stmt.WriteString("COPY (SELECT instance_id, aggregate_type, aggregate_id, event_type, sequence, revision, created_at, regexp_replace(payload::TEXT, '\\\\u0000', '', 'g')::JSON payload, creator, owner, ")
				stmt.WriteArg(position)
				stmt.WriteString(" position, row_number() OVER (PARTITION BY instance_id ORDER BY position, in_tx_order) AS in_tx_order FROM eventstore.events2 ")
				stmt.WriteString(instanceClause())
				stmt.WriteString(" AND ")
				database.NewNumberAtMost(maxPosition).Write(&stmt, "position")
				stmt.WriteString(" AND ")
				database.NewNumberGreater(previousMigration.Position).Write(&stmt, "position")
				stmt.WriteString(" ORDER BY instance_id, position, in_tx_order")
				stmt.WriteString(" LIMIT ")
				stmt.WriteArg(bulkSize)
				stmt.WriteString(" OFFSET ")
				stmt.WriteArg(bulkSize * i)
				stmt.WriteString(") TO STDOUT")

				// Copy does not allow args so we use we replace the args in the statement
				tag, err := conn.PgConn().CopyTo(ctx, writer, stmt.Debug())
				if err != nil {
					return zerrors.ThrowUnknownf(err, "MIGRA-KTuSq", "unable to copy events from source during iteration %d", i)
				}
				if tag.RowsAffected() < int64(bulkSize) {
					return nil
				}

				nextPos <- true
				i++
			}
			return nil
		})
		writer.Close()
		close(nextPos)
		errs <- goErr
	}()

	// generate next position for
	go func() {
		defer close(pos)
		for range nextPos {
			var position float64
			goErr := dest.QueryRowContext(
				ctx,
				func(row *sql.Row) error {
					return row.Scan(&position)
				},
				positionQuery(dest),
			)
			if goErr != nil {
				errs <- zerrors.ThrowUnknown(err, "MIGRA-kMyPH", "unable to query next position")
				return
			}
			pos <- position
		}
	}()
	var eventCount int64
	errs <- destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		tag, cbErr := conn.PgConn().CopyFrom(ctx, reader, "COPY eventstore.events2 FROM STDIN")
		eventCount = tag.RowsAffected()
		if cbErr != nil {
			return zerrors.ThrowUnknown(cbErr, "MIGRA-DTHi7", "unable to copy events into destination")
		}

		return nil
	})
	close(errs)
	if err = writeCopyEventsDone(ctx, destinationES, migrationID, source.DatabaseName(), maxPosition, errs); err != nil {
		return fmt.Errorf("unable to write migration done: %w", err)
	}
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("events migrated")
	return nil
}

func writeCopyEventsDone(ctx context.Context, es *eventstore.EventStore, id, source string, position decimal.Decimal, errs <-chan error) error {
	joinedErrs := make([]error, 0, len(errs))
	for err := range errs {
		joinedErrs = append(joinedErrs, err)
	}
	err := errors.Join(joinedErrs...)

	if err != nil {
		logging.WithError(err).Error("unable to mirror events")
		err = writeMigrationFailed(ctx, es, id, source, err)
		if err != nil {
			return fmt.Errorf("unable to write failed event: %w", err)
		}
		return nil
	}

	if err = writeMigrationSucceeded(ctx, es, id, source, position); err != nil {
		return fmt.Errorf("unable to write succeeded event: %w", err)
	}
	return nil
}

func copyUniqueConstraints(ctx context.Context, source, dest *db.DB) error {
	start := time.Now()
	reader, writer := io.Pipe()
	errs := make(chan error, 1)

	sourceConn, err := source.Conn(ctx)
	if err != nil {
		return fmt.Errorf("unable to acquire source connection: %w", err)
	}

	go func() {
		errs <- sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()
			var stmt database.Statement
			stmt.WriteString("COPY (SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints ")
			stmt.WriteString(instanceClause())
			stmt.WriteString(") TO stdout")

			_, err := conn.PgConn().CopyTo(ctx, writer, stmt.String())
			writer.Close()
			return err
		})
	}()

	destConn, err := dest.Conn(ctx)
	if err != nil {
		return fmt.Errorf("unable to acquire dest connection: %w", err)
	}

	var eventCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		if shouldReplace {
			var stmt database.Statement
			stmt.WriteString("DELETE FROM eventstore.unique_constraints ")
			stmt.WriteString(instanceClause())

			_, cbErr := conn.Exec(ctx, stmt.String())
			if cbErr != nil {
				return cbErr
			}
		}

		tag, cbErr := conn.PgConn().CopyFrom(ctx, reader, "COPY eventstore.unique_constraints FROM stdin")
		eventCount = tag.RowsAffected()
		return cbErr
	})
	if err != nil {
		return fmt.Errorf("unable to copy unique constraints to destination: %w", err)
	}
	if err = <-errs; err != nil {
		return fmt.Errorf("unable to copy unique constraints from source: %w", err)
	}
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("unique constraints migrated")
	return nil
}
