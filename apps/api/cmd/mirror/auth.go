package mirror

import (
	"context"
	_ "embed"
	"io"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

func authCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "mirrors the auth requests table from one database to another",
		Long: `mirrors the auth requests table from one database to another
ZITADEL needs to be initialized and set up with the --for-mirror flag
Only auth requests are mirrored`,
		Run: func(cmd *cobra.Command, args []string) {
			config := mustNewMigrationConfig(viper.GetViper())
			copyAuth(cmd.Context(), config)
		},
	}

	cmd.Flags().BoolVar(&shouldReplace, "replace", false, "allow delete auth requests of defined instances before copy")

	return cmd
}

func copyAuth(ctx context.Context, config *Migration) {
	sourceClient, err := database.Connect(config.Source, false)
	logging.OnError(err).Fatal("unable to connect to source database")
	defer sourceClient.Close()

	destClient, err := database.Connect(config.Destination, false)
	logging.OnError(err).Fatal("unable to connect to destination database")
	defer destClient.Close()

	copyAuthRequests(ctx, sourceClient, destClient, config.MaxAuthRequestAge)
}

func copyAuthRequests(ctx context.Context, source, dest *database.DB, maxAuthRequestAge time.Duration) {
	start := time.Now()

	logging.Info("creating index on auth.auth_requests.change_date to speed up copy in source database")
	_, err := source.ExecContext(ctx, "CREATE INDEX CONCURRENTLY IF NOT EXISTS auth_requests_change_date ON auth.auth_requests (change_date)")
	logging.OnError(err).Fatal("unable to create index on auth.auth_requests.change_date")

	sourceConn, err := source.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire connection")
	defer sourceConn.Close()

	r, w := io.Pipe()
	errs := make(chan error, 1)

	go func() {
		err = sourceConn.Raw(func(driverConn any) error {
			conn := driverConn.(*stdlib.Conn).Conn()
			_, err := conn.PgConn().CopyTo(ctx, w, "COPY (SELECT id, regexp_replace(request::TEXT, '\\\\u0000', '', 'g')::JSON request, code, request_type, creation_date, change_date, instance_id FROM auth.auth_requests "+instanceClause()+" AND change_date > NOW() - INTERVAL '"+strconv.FormatFloat(maxAuthRequestAge.Seconds(), 'f', -1, 64)+" seconds') TO STDOUT")
			w.Close()
			return err
		})
		errs <- err
	}()

	destConn, err := dest.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire connection")
	defer destConn.Close()

	var affected int64
	err = destConn.Raw(func(driverConn any) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		if shouldReplace {
			_, err := conn.Exec(ctx, "DELETE FROM auth.auth_requests "+instanceClause())
			if err != nil {
				return err
			}
		}

		tag, err := conn.PgConn().CopyFrom(ctx, r, "COPY auth.auth_requests FROM STDIN")
		affected = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy auth requests to destination")
	logging.OnError(<-errs).Fatal("unable to copy auth requests from source")
	logging.WithFields("took", time.Since(start), "count", affected).Info("auth requests migrated")
}
