package migrate

import (
	"context"
	_ "embed"
	"io"

	"github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	instanceID string
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "migrates the events of an instance from one database to another",
		Long: `migrates the events of an instance from one database to another
ZITADEL needs to be initialized
Migrate only copies the events`,
		Run: func(cmd *cobra.Command, args []string) {
			config := MustNewConfig(viper.GetViper())
			migrate(cmd.Context(), config, instanceID)
		},
	}

	Flags(cmd)

	return cmd
}

func Flags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&instanceID, "instance", "", "id of the instance to migrate")
	cmd.PersistentFlags().StringArrayVar(&configPaths, "config", nil, "paths to config files")
}

func migrate(ctx context.Context, config *Config, instanceID string) {
	if instanceID == "" {
		logging.Fatal("no instance id set")
	}
	sourceClient, err := database.Connect(config.Source, false, false)
	logging.OnError(err).Fatal("unable to connect to source database")

	destClient, err := database.Connect(config.Destination, false, true)
	logging.OnError(err).Fatal("unable to connect to destination database")

	sourceConn, err := sourceClient.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire connection")

	r, w := io.Pipe()
	errs := make(chan error, 1)

	go func() {
		err = sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn() // conn is a *pgx.Conn
			// Do pgx specific stuff with conn

			_, err := conn.PgConn().CopyTo(ctx, w, "COPY (SELECT * FROM eventstore.events2 where instance_id = '"+instanceID+"' ORDER BY position, in_tx_order) TO stdout")
			w.Close()
			return err
		})
		errs <- err
	}()

	destConn, err := destClient.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire connection")

	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn() // conn is a *pgx.Conn
		// Do pgx specific stuff with conn
		_, err := conn.PgConn().CopyFrom(ctx, r, "COPY eventstore.events2 FROM stdin")
		// conn.CopyFrom(...)
		return err
	})
	logging.OnError(err).Fatal("unable to copy events to destination")
	logging.OnError(<-errs).Fatal("unable to copy events from source")
}
