package migrate

import (
	"context"
	_ "embed"
	"io"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

func eventstoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "eventstore",
		Short: "migrates the eventstore of an instance from one database to another",
		Long: `migrates the eventstore of an instance from one database to another
ZITADEL needs to be initialized
Migrate only copies events2 and unique constraints`,
		Run: func(cmd *cobra.Command, args []string) {
			config := mustNewMigrationConfig(viper.GetViper())
			copyEventstore(cmd.Context(), config)
		},
	}
}

func copyEventstore(ctx context.Context, config *Migration) {
	sourceClient, err := database.Connect(config.Source, false, false)
	logging.OnError(err).Fatal("unable to connect to source database")
	defer sourceClient.Close()

	destClient, err := database.Connect(config.Destination, false, true)
	logging.OnError(err).Fatal("unable to connect to destination database")
	defer destClient.Close()

	copyEvents(ctx, sourceClient, destClient)
	copyUniqueConstraints(ctx, sourceClient, destClient)
}

func positionQuery(db *database.DB) string {
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

func copyEvents(ctx context.Context, source, dest *database.DB) {
	start := time.Now()
	reader, writer := io.Pipe()
	errs := make(chan error, 1)

	sourceConn, err := source.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire source connection")

	// get position
	pos := make(chan float64)

	go func() {
		position := strconv.FormatFloat(<-pos, 'E', -1, 64)
		err := sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()
			// TODO: sql injection
			_, err := conn.PgConn().CopyTo(ctx, writer, "COPY (SELECT instance_id, aggregate_type, aggregate_id, event_type, sequence, revision, created_at,  regexp_replace(payload::TEXT, '\\\\u0000', '', 'g')::JSON payload, creator, owner, (SELECT "+position+"::DECIMAL) AS position, row_number() OVER (PARTITION BY instance_id ORDER BY position, in_tx_order) AS in_tx_order FROM eventstore.events2 "+instanceClause()+" ORDER BY instance_id, position, in_tx_order) TO STDOUT")
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
		tx, err := conn.Begin(ctx)
		if err != nil {
			return err
		}
		row := tx.QueryRow(ctx, positionQuery(dest))
		var position float64
		if err := row.Scan(&position); err != nil {
			return err
		}
		_ = tx.Commit(ctx)
		pos <- position

		tag, err := conn.PgConn().CopyFrom(ctx, reader, "COPY eventstore.events2 FROM STDIN")
		eventCount = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy events to destination")
	logging.OnError(<-errs).Fatal("unable to copy events from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("events migrated")
}

func copyUniqueConstraints(ctx context.Context, source, dest *database.DB) {
	start := time.Now()
	reader, writer := io.Pipe()
	errs := make(chan error, 1)

	sourceConn, err := source.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire source connection")

	go func() {
		err := sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()
			// TODO: sql injection
			_, err := conn.PgConn().CopyTo(ctx, writer, "COPY (SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints "+instanceClause()+") TO stdout")
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

		tag, err := conn.PgConn().CopyFrom(ctx, reader, "COPY eventstore.unique_constraints FROM stdin")
		eventCount = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy unique constraints to destination")
	logging.OnError(<-errs).Fatal("unable to copy unique constraints from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("unique constraints migrated")
}
