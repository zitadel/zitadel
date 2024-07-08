package mirror

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"io"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	db "github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
	"github.com/zitadel/zitadel/internal/id_generator"
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
		Run: func(cmd *cobra.Command, args []string) {
			config := mustNewMigrationConfig(viper.GetViper())
			copyEventstore(cmd.Context(), config)
		},
	}

	cmd.Flags().BoolVar(&shouldReplace, "replace", false, "allow delete unique constraints of defined instances before copy")
	cmd.Flags().BoolVar(&shouldIgnorePrevious, "ignore-previous", false, "ignores previous migrations of the events table")

	return cmd
}

func copyEventstore(ctx context.Context, config *Migration) {
	sourceClient, err := db.Connect(config.Source, false, dialect.DBPurposeQuery)
	logging.OnError(err).Fatal("unable to connect to source database")
	defer sourceClient.Close()

	destClient, err := db.Connect(config.Destination, false, dialect.DBPurposeEventPusher)
	logging.OnError(err).Fatal("unable to connect to destination database")
	defer destClient.Close()

	copyEvents(ctx, sourceClient, destClient, config.EventBulkSize)
	copyUniqueConstraints(ctx, sourceClient, destClient)
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

func copyEvents(ctx context.Context, source, dest *db.DB, bulkSize uint32) {
	start := time.Now()
	reader, writer := io.Pipe()

	migrationID, err := id_generator.Next()
	logging.OnError(err).Fatal("unable to generate migration id")

	sourceConn, err := source.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire source connection")

	destConn, err := dest.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire dest connection")

	sourceES := eventstore.NewEventstoreFromOne(postgres.New(source, &postgres.Config{
		MaxRetries: 3,
	}))
	destinationES := eventstore.NewEventstoreFromOne(postgres.New(dest, &postgres.Config{
		MaxRetries: 3,
	}))

	previousMigration, err := queryLastSuccessfulMigration(ctx, destinationES, source.DatabaseName())
	logging.OnError(err).Fatal("unable to query latest successful migration")

	maxPosition, err := writeMigrationStart(ctx, sourceES, migrationID, dest.DatabaseName())
	logging.OnError(err).Fatal("unable to write migration started event")

	logging.WithFields("from", previousMigration.Position, "to", maxPosition).Info("start event migration")

	nextPos := make(chan bool, 1)
	pos := make(chan float64, 1)
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
		errs <- err
	}()

	// generate next position for
	go func() {
		defer close(pos)
		for range nextPos {
			var position float64
			err := dest.QueryRowContext(
				ctx,
				func(row *sql.Row) error {
					return row.Scan(&position)
				},
				positionQuery(dest),
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
			return zerrors.ThrowUnknown(err, "MIGRA-DTHi7", "unable to copy events into destination")
		}

		return nil
	})

	close(errs)
	writeCopyEventsDone(ctx, destinationES, migrationID, source.DatabaseName(), maxPosition, errs)

	logging.WithFields("took", time.Since(start), "count", eventCount).Info("events migrated")
}

func writeCopyEventsDone(ctx context.Context, es *eventstore.EventStore, id, source string, position float64, errs <-chan error) {
	joinedErrs := make([]error, 0, len(errs))
	for err := range errs {
		joinedErrs = append(joinedErrs, err)
	}
	err := errors.Join(joinedErrs...)

	if err != nil {
		logging.WithError(err).Error("unable to mirror events")
		err := writeMigrationFailed(ctx, es, id, source, err)
		logging.OnError(err).Fatal("unable to write failed event")
		return
	}

	err = writeMigrationSucceeded(ctx, es, id, source, position)
	logging.OnError(err).Fatal("unable to write failed event")
}

func copyUniqueConstraints(ctx context.Context, source, dest *db.DB) {
	start := time.Now()
	reader, writer := io.Pipe()
	errs := make(chan error, 1)

	sourceConn, err := source.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire source connection")

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
	logging.OnError(err).Fatal("unable to acquire dest connection")

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
	logging.OnError(err).Fatal("unable to copy unique constraints to destination")
	logging.OnError(<-errs).Fatal("unable to copy unique constraints from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("unique constraints migrated")
}
