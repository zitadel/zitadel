package mirror

import (
	"context"
	_ "embed"
	"io"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

func systemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "mirrors the system tables of ZITADEL from one database to another",
		Long: `mirrors the system tables of ZITADEL from one database to another
ZITADEL needs to be initialized
Only keys and assets are mirrored`,
		Run: func(cmd *cobra.Command, args []string) {
			config := mustNewMigrationConfig(viper.GetViper())
			copySystem(cmd.Context(), config)
		},
	}

	cmd.Flags().BoolVar(&shouldReplace, "replace", false, "allow delete ALL keys and assets of defined instances before copy")

	return cmd
}

func copySystem(ctx context.Context, config *Migration) {
	sourceClient, err := database.Connect(config.Source, false)
	logging.OnError(err).Fatal("unable to connect to source database")
	defer sourceClient.Close()

	destClient, err := database.Connect(config.Destination, false)
	logging.OnError(err).Fatal("unable to connect to destination database")
	defer destClient.Close()

	copyAssets(ctx, sourceClient, destClient)
	copyEncryptionKeys(ctx, sourceClient, destClient)
}

func copyAssets(ctx context.Context, source, dest *database.DB) {
	logging.Info("starting to copy assets")
	start := time.Now()

	sourceConn, err := source.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire source connection")
	defer sourceConn.Close()

	r, w := io.Pipe()
	errs := make(chan error, 1)

	go func() {
		err = sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()
			// ignore hash column because it's computed
			_, err := conn.PgConn().CopyTo(ctx, w, "COPY (SELECT instance_id, asset_type, resource_owner, name, content_type, data, updated_at FROM system.assets "+instanceClause()+") TO stdout")
			w.Close()
			return err
		})
		errs <- err
	}()

	destConn, err := dest.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire dest connection")
	defer destConn.Close()

	var assetCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		if shouldReplace {
			_, err := conn.Exec(ctx, "DELETE FROM system.assets "+instanceClause())
			if err != nil {
				return err
			}
		}

		tag, err := conn.PgConn().CopyFrom(ctx, r, "COPY system.assets (instance_id, asset_type, resource_owner, name, content_type, data, updated_at) FROM stdin")
		assetCount = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy assets to destination")
	logging.OnError(<-errs).Fatal("unable to copy assets from source")
	logging.WithFields("took", time.Since(start), "count", assetCount).Info("assets migrated")
}

func copyEncryptionKeys(ctx context.Context, source, dest *database.DB) {
	logging.Info("starting to copy encryption keys")
	start := time.Now()

	sourceConn, err := source.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire source connection")
	defer sourceConn.Close()

	r, w := io.Pipe()
	errs := make(chan error, 1)

	go func() {
		err = sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()
			// ignore hash column because it's computed
			_, err := conn.PgConn().CopyTo(ctx, w, "COPY system.encryption_keys TO stdout")
			w.Close()
			return err
		})
		errs <- err
	}()

	destConn, err := dest.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire dest connection")
	defer destConn.Close()

	var keyCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		if shouldReplace {
			_, err := conn.Exec(ctx, "TRUNCATE system.encryption_keys")
			if err != nil {
				return err
			}
		}

		tag, err := conn.PgConn().CopyFrom(ctx, r, "COPY system.encryption_keys FROM stdin")
		keyCount = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy encryption keys to destination")
	logging.OnError(<-errs).Fatal("unable to copy encryption keys from source")
	logging.WithFields("took", time.Since(start), "count", keyCount).Info("encryption keys migrated")
}
