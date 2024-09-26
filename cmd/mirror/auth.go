package mirror

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
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

func copyAuth(ctx context.Context, config *Migration) error {
	sourceClient, err := database.Connect(config.Source, false, dialect.DBPurposeQuery)
	if err != nil {
		return fmt.Errorf("unable to connect to source database: %w", err)
	}
	defer sourceClient.Close()

	destClient, err := database.Connect(config.Destination, false, dialect.DBPurposeEventPusher)
	if err != nil {
		return fmt.Errorf("unable to connect to destination database: %w", err)
	}
	defer destClient.Close()

	return copyAuthRequests(ctx, sourceClient, destClient)
}

func copyAuthRequests(ctx context.Context, source, dest *database.DB) error {
	start := time.Now()

	sourceConn, err := source.Conn(ctx)
	if err != nil {
		return fmt.Errorf("unable to acquire connection: %w", err)
	}
	defer sourceConn.Close()

	r, w := io.Pipe()
	errs := make(chan error, 1)

	go func() {
		err = sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()
			_, err := conn.PgConn().CopyTo(ctx, w, "COPY (SELECT id, regexp_replace(request::TEXT, '\\\\u0000', '', 'g')::JSON request, code, request_type, creation_date, change_date, instance_id FROM auth.auth_requests "+instanceClause()+") TO STDOUT")
			w.Close()
			return err
		})
		errs <- err
	}()

	destConn, err := dest.Conn(ctx)
	if err != nil {
		return fmt.Errorf("unable to acquire connection: %w", err)
	}
	defer destConn.Close()

	var affected int64
	err = destConn.Raw(func(driverConn interface{}) error {
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
	if err != nil {
		return fmt.Errorf("unable to copy auth requests to destination: %w", err)
	}
	if err = <-errs; err != nil {
		return fmt.Errorf("unable to copy auth requests from source: %w", err)
	}
	logging.WithFields("took", time.Since(start), "count", affected).Info("auth requests migrated")
	return nil
}
