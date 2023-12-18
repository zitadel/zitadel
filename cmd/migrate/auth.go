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

func authCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "migrates the auth requests table from one database to another",
		Long: `migrates the auth requests table from one database to another
ZITADEL needs to be initialized
Migrations only copies auth requests`,
		Run: func(cmd *cobra.Command, args []string) {
			config := mustNewMigrationConfig(viper.GetViper())
			copyAuth(cmd.Context(), config)
		},
	}
}

func copyAuth(ctx context.Context, config *Migration) {
	sourceClient, err := database.Connect(config.Source, false, false)
	logging.OnError(err).Fatal("unable to connect to source database")
	defer sourceClient.Close()

	destClient, err := database.Connect(config.Destination, false, true)
	logging.OnError(err).Fatal("unable to connect to destination database")
	defer destClient.Close()

	copyAuthRequests(ctx, sourceClient, destClient)
}

func copyAuthRequests(ctx context.Context, source, dest *database.DB) {
	start := time.Now()

	sourceConn, err := source.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire connection")
	defer sourceConn.Close()

	r, w := io.Pipe()
	errs := make(chan error, 1)

	go func() {
		err = sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()
			_, err := conn.PgConn().CopyTo(ctx, w, "COPY (SELECT * FROM auth.auth_requests "+instanceClause()+") TO stdout")
			w.Close()
			return err
		})
		errs <- err
	}()

	destConn, err := dest.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire connection")
	defer destConn.Close()

	var eventCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		tag, err := conn.PgConn().CopyFrom(ctx, r, "COPY auth.auth_requests FROM stdin")
		eventCount = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy auth requests to destination")
	logging.OnError(<-errs).Fatal("unable to copy auth requests from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("auth requests migrated")
}
