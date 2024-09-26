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

func systemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "mirrors the system tables of ZITADEL from one database to another",
		Long: `mirrors the system tables of ZITADEL from one database to another
ZITADEL needs to be initialized
Only keys and assets are mirrored`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := mustNewMigrationConfig(viper.GetViper())
			return copySystem(cmd.Context(), config)
		},
	}

	cmd.Flags().BoolVar(&shouldReplace, "replace", false, "allow delete ALL keys and assets of defined instances before copy")

	return cmd
}

func copySystem(ctx context.Context, config *Migration) error {
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

	if err = copyAssets(ctx, sourceClient, destClient); err != nil {
		return fmt.Errorf("unable to copy assets: %w", err)
	}
	if err = copyEncryptionKeys(ctx, sourceClient, destClient); err != nil {
		return fmt.Errorf("unable to copy encryption keys: %w", err)
	}
	return nil
}

func copyAssets(ctx context.Context, source, dest *database.DB) error {
	start := time.Now()

	sourceConn, err := source.Conn(ctx)
	if err != nil {
		return fmt.Errorf("unable to acquire source connection: %w", err)
	}
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
	if err != nil {
		return fmt.Errorf("unable to acquire dest connection: %w", err)
	}
	defer destConn.Close()

	var eventCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		if shouldReplace {
			_, err := conn.Exec(ctx, "DELETE FROM system.assets "+instanceClause())
			if err != nil {
				return err
			}
		}

		tag, err := conn.PgConn().CopyFrom(ctx, r, "COPY system.assets (instance_id, asset_type, resource_owner, name, content_type, data, updated_at) FROM stdin")
		eventCount = tag.RowsAffected()

		return err
	})
	if err != nil {
		return fmt.Errorf("unable to copy assets to destination: %w", err)
	}
	if err = <-errs; err != nil {
		return fmt.Errorf("unable to copy assets from source: %w", err)
	}
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("assets migrated")
	return nil
}

func copyEncryptionKeys(ctx context.Context, source, dest *database.DB) error {
	start := time.Now()

	sourceConn, err := source.Conn(ctx)
	if err != nil {
		return fmt.Errorf("unable to acquire source connection: %w", err)
	}
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
	if err != nil {
		return fmt.Errorf("unable to acquire dest connection: %w", err)
	}
	defer destConn.Close()

	var eventCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		if shouldReplace {
			_, err := conn.Exec(ctx, "TRUNCATE system.encryption_keys")
			if err != nil {
				return err
			}
		}

		tag, err := conn.PgConn().CopyFrom(ctx, r, "COPY system.encryption_keys FROM stdin")
		eventCount = tag.RowsAffected()

		return err
	})
	if err != nil {
		return fmt.Errorf("unable to copy encryption keys to destination: %w", err)
	}
	if err = <-errs; err != nil {
		return fmt.Errorf("unable to copy encryption keys from source: %w", err)
	}
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("encryption keys migrated")
	return nil
}
