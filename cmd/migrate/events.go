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

func eventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "migrates the events of an instance from one database to another",
		Long: `migrates the events of an instance from one database to another
ZITADEL needs to be initialized
Migrate only copies the events`,
		Run: func(cmd *cobra.Command, args []string) {
			config := mustNewMigrationConfig(viper.GetViper())
			events(cmd.Context(), config, instanceID)
		},
	}

	migrateEventsFlags(cmd)

	return cmd
}

func migrateEventsFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVar(&configPaths, "config", nil, "paths to config files")
}

func events(ctx context.Context, config *Migration, instanceID string) {
	start := time.Now()
	if instanceID == "" {
		logging.Fatal("no instance id set")
	}
	sourceClient, err := database.Connect(config.Source, false, false)
	logging.OnError(err).Fatal("unable to connect to source database")
	defer sourceClient.Close()

	destClient, err := database.Connect(config.Destination, false, true)
	logging.OnError(err).Fatal("unable to connect to destination database")
	defer destClient.Close()

	sourceConn, err := sourceClient.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire connection")

	r, w := io.Pipe()
	errs := make(chan error, 1)

	// get position
	pos := make(chan float64)

	go func() {
		position := strconv.FormatFloat(<-pos, 'E', -1, 64)
		err = sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn() // conn is a *pgx.Conn
			// Do pgx specific stuff with conn
			// TODO: sql injection
			_, err := conn.PgConn().CopyTo(ctx, w, "COPY (SELECT instance_id, aggregate_type, aggregate_id, event_type, sequence, revision, created_at, payload, creator, owner, (SELECT "+position+"::DECIMAL) AS position, row_number() OVER () AS in_tx_order FROM eventstore.events2 where instance_id = '"+instanceID+"' ORDER BY position, in_tx_order) TO stdout")
			w.Close()
			// TODO: unique constraints, assets
			return err
		})
		errs <- err
	}()

	destConn, err := destClient.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire connection")

	var eventCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()
		tx, err := conn.Begin(ctx)
		if err != nil {
			return err
		}
		row := tx.QueryRow(ctx, positionQuery(destClient))
		var position float64
		if err := row.Scan(&position); err != nil {
			return err
		}
		_ = tx.Commit(ctx)
		pos <- position

		tag, err := conn.PgConn().CopyFrom(ctx, r, "COPY eventstore.events2 FROM stdin")
		eventCount = tag.RowsAffected()
		// TODO: unique constraints, assets

		return err
	})
	logging.OnError(err).Fatal("unable to copy events to destination")
	logging.OnError(<-errs).Fatal("unable to copy events from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("events migrated")
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
