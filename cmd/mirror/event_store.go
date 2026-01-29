package mirror

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"io"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
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
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				logging.OnError(cmd.Context(), err).ErrorContext(cmd.Context(), "zitadel mirror eventstore command failed")
			}()

			config, shutdown, err := newMigrationConfig(cmd.Context(), viper.GetViper())
			if err != nil {
				return err
			}
			// Set logger again to include changes from config
			cmd.SetContext(logging.NewCtx(cmd.Context(), logging.StreamRuntime))
			defer func() {
				err = errors.Join(err, shutdown(cmd.Context()))
			}()
			defer func() {
				if recErr, ok := recover().(error); ok {
					err = recErr
				}
			}()
			copyEventstore(cmd.Context(), config)
			return nil
		},
	}

	cmd.Flags().BoolVar(&shouldReplace, "replace", false, "allow delete unique constraints of defined instances before copy")
	cmd.Flags().BoolVar(&shouldIgnorePrevious, "ignore-previous", false, "ignores previous migrations of the events table")

	return cmd
}

func copyEventstore(ctx context.Context, config *Migration) {
	sourceClient, err := db.Connect(config.Source, false)
	panicOnError(ctx, err, "unable to connect to source database")
	defer sourceClient.Close()

	destClient, err := db.Connect(config.Destination, false)
	panicOnError(ctx, err, "unable to connect to destination database")
	defer destClient.Close()

	copyEvents(ctx, sourceClient, destClient, config.EventBulkSize)
	copyUniqueConstraints(ctx, sourceClient, destClient)
}

func positionQuery(db *db.DB) (string, error) {
	switch db.Type() {
	case dialect.DatabaseTypePostgres:
		return "SELECT EXTRACT(EPOCH FROM clock_timestamp())", nil
	case dialect.DatabaseTypeCockroach:
		return "SELECT cluster_logical_timestamp()", nil
	default:
		return "", errors.New("database type not recognized")
	}
}

func copyEvents(ctx context.Context, source, dest *db.DB, bulkSize uint32) {
	logging.Info(ctx, "starting to copy events")
	start := time.Now()
	reader, writer := io.Pipe()

	migrationID, err := id.SonyFlakeGenerator().Next()
	panicOnError(ctx, err, "unable to generate migration id")

	sourceConn, err := source.Conn(ctx)
	panicOnError(ctx, err, "unable to acquire source connection")

	destConn, err := dest.Conn(ctx)
	panicOnError(ctx, err, "unable to acquire dest connection")

	destinationES := eventstore.NewEventstoreFromOne(postgres.New(dest, &postgres.Config{
		MaxRetries: 3,
	}))

	previousMigration, err := queryLastSuccessfulMigration(ctx, destinationES, source.DatabaseName())
	panicOnError(ctx, err, "unable to query latest successful migration")

	var maxPosition decimal.Decimal
	err = source.QueryRowContext(ctx,
		func(row *sql.Row) error {
			return row.Scan(&maxPosition)
		},
		"SELECT MAX(position) FROM eventstore.events2 "+instanceClause(),
	)
	panicOnError(ctx, err, "unable to query max position from source")
	logging.Info(ctx, "start event migration", "from", previousMigration.Position, "to", maxPosition)

	nextPos := make(chan bool, 1)
	pos := make(chan decimal.Decimal, 1)
	errs := make(chan error, 3)

	go func() {
		err := sourceConn.Raw(func(driverConn interface{}) error {
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
				logging.Info(ctx, "batch of events copied", "batch_count", i)

				if tag.RowsAffected() < int64(bulkSize) {
					logging.Info(ctx, "last batch of events copied", "batch_count", i)
					return nil
				}

				nextPos <- true
				i++
			}
			return nil
		})
		writer.Close()
		close(nextPos)
		errs <- err
	}()

	// generate next position for
	go func() {
		defer close(pos)
		for range nextPos {
			var position decimal.Decimal
			query, err := positionQuery(dest)
			if err != nil {
				errs <- zerrors.ThrowUnknown(err, "MIGRA-Hy6t3", "unable to generate position query")
				return
			}

			err = dest.QueryRowContext(
				ctx,
				func(row *sql.Row) error {
					return row.Scan(&position)
				},
				query,
			)
			if err != nil {
				errs <- zerrors.ThrowUnknown(err, "MIGRA-kMyPH", "unable to query next position")
				return
			}
			pos <- position
		}
	}()

	var eventCount int64
	errs <- destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		tag, err := conn.PgConn().CopyFrom(ctx, reader, "COPY eventstore.events2 FROM STDIN")
		eventCount = tag.RowsAffected()
		if err != nil {
			pgErr := new(pgconn.PgError)
			errors.As(err, &pgErr)

			logging.WithError(ctx, err).ErrorContext(ctx, "unable to copy events into destination", "pg_err_details", pgErr.Detail)
			return zerrors.ThrowUnknown(err, "MIGRA-DTHi7", "unable to copy events into destination")
		}

		return nil
	})

	close(errs)
	writeCopyEventsDone(ctx, destinationES, migrationID, source.DatabaseName(), maxPosition, errs)

	logging.Info(ctx, "events migrated", "took", time.Since(start), "count", eventCount)
}

func writeCopyEventsDone(ctx context.Context, es *eventstore.EventStore, id, source string, position decimal.Decimal, errs <-chan error) {
	joinedErrs := make([]error, 0, len(errs))
	for err := range errs {
		joinedErrs = append(joinedErrs, err)
	}
	err := errors.Join(joinedErrs...)

	if err != nil {
		logging.WithError(ctx, err).Error("unable to mirror events")
		err := writeMigrationFailed(ctx, es, id, source, err)
		panicOnError(ctx, err, "unable to write failed event")
		return
	}

	err = writeMigrationSucceeded(ctx, es, id, source, position)
	panicOnError(ctx, err, "unable to write succeeded event")
}

func copyUniqueConstraints(ctx context.Context, source, dest *db.DB) {
	logging.Info(ctx, "starting to copy unique constraints")
	start := time.Now()
	reader, writer := io.Pipe()
	errs := make(chan error, 1)

	sourceConn, err := source.Conn(ctx)
	panicOnError(ctx, err, "unable to acquire source connection")

	go func() {
		err := sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()
			var stmt database.Statement
			stmt.WriteString("COPY (SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints ")
			stmt.WriteString(instanceClause())
			stmt.WriteString(") TO stdout")

			_, err := conn.PgConn().CopyTo(ctx, writer, stmt.String())
			writer.Close()
			return err
		})
		errs <- err
	}()

	destConn, err := dest.Conn(ctx)
	panicOnError(ctx, err, "unable to acquire dest connection")

	var eventCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		if shouldReplace {
			var stmt database.Statement
			stmt.WriteString("DELETE FROM eventstore.unique_constraints ")
			stmt.WriteString(instanceClause())

			_, err := conn.Exec(ctx, stmt.String())
			if err != nil {
				return err
			}
		}

		tag, err := conn.PgConn().CopyFrom(ctx, reader, "COPY eventstore.unique_constraints FROM stdin")
		eventCount = tag.RowsAffected()

		return err
	})
	panicOnError(ctx, err, "unable to copy unique constraints to destination")
	panicOnError(ctx, <-errs, "unable to copy unique constraints from source")
	logging.Info(ctx, "unique constraints migrated", "took", time.Since(start), "count", eventCount)
}
