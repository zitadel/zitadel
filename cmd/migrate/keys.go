package migrate

import (
	"context"
	_ "embed"
	"io"
	"time"

	"github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

func keysCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "keys",
		Short: "migrates the keys of the system from one database to another",
		Long: `migrates the keys of the system from one database to another
ZITADEL needs to be initialized
Migrations only copies the keys`,
		Run: func(cmd *cobra.Command, args []string) {
			config := mustNewMigrationConfig(viper.GetViper())
			keys(cmd.Context(), config)
		},
	}
}

func keys(ctx context.Context, config *Migration) {
	start := time.Now()
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

	go func() {
		err = sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()
			_, err := conn.PgConn().CopyTo(ctx, w, "COPY (SELECT id, key FROM system.encryption_keys) TO stdout")
			w.Close()
			return err
		})
		errs <- err
	}()

	destConn, err := destClient.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire connection")

	var eventCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		tag, err := conn.PgConn().CopyFrom(ctx, r, "COPY system.encryption_keys FROM stdin")
		eventCount = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy encryption keys to destination")
	logging.OnError(<-errs).Fatal("unable to copy encryption keys from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("encryption keys migrated")
}
