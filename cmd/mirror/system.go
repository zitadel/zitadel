package mirror

import (
	"context"
	_ "embed"
	"io"
	"os"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	db "github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
	"github.com/zitadel/zitadel/internal/v2/database"
)

func systemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "mirrors the system tables of ZITADEL between databases, or between a database and files",
		Long: `mirrors the system tables of ZITADEL between databases, or between a database and files
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
	switch {
	case isSrcFile:
		destClient, err := db.Connect(config.Destination, false, dialect.DBPurposeEventPusher)
		logging.OnError(err).Fatal("unable to connect to destination database")
		defer destClient.Close()

		copyAssetsFromFile(ctx, destClient, "system.assets.csv")
		copyEncryptionKeysFromFile(ctx, destClient, "system.encryption_keys.csv")
	case isDestFile:
		sourceClient, err := db.Connect(config.Source, false, dialect.DBPurposeQuery)
		logging.OnError(err).Fatal("unable to connect to source database")
		defer sourceClient.Close()

		copyAssetsToFile(ctx, sourceClient, "system.assets.csv")
		copyEncryptionKeysToFile(ctx, sourceClient, "system.encryption_keys.csv")
	default:
		sourceClient, err := db.Connect(config.Source, false, dialect.DBPurposeQuery)
		logging.OnError(err).Fatal("unable to connect to source database")
		defer sourceClient.Close()

		destClient, err := db.Connect(config.Destination, false, dialect.DBPurposeEventPusher)
		logging.OnError(err).Fatal("unable to connect to destination database")
		defer destClient.Close()

		copyAssetsDB(ctx, sourceClient, destClient)
		copyEncryptionKeysDB(ctx, sourceClient, destClient)
	}
}

func copyAssetsDB(ctx context.Context, source, dest *db.DB) {
	start := time.Now()

	sourceConn, err := source.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire source connection")
	defer sourceConn.Close()

	reader, writer := io.Pipe()
	errs := make(chan error, 1)

	go func() {
		err = sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()

			// ignore hash column because it's computed
			var stmt database.Statement
			stmt.WriteString(`COPY (SELECT instance_id, asset_type, 
							resource_owner, name, content_type, data, updated_at 
							FROM system.assets `)
			stmt.WriteString(instanceClause())
			stmt.WriteString(") TO STDOUT")

			_, err := conn.PgConn().CopyTo(ctx, writer, stmt.String())
			writer.Close()

			return err
		})
		errs <- err
	}()

	destConn, err := dest.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire destination connection")
	defer destConn.Close()

	var eventCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		var stmt database.Statement
		stmt.WriteString("DELETE FROM system.assets ")
		stmt.WriteString(instanceClause())

		if shouldReplace {
			_, err := conn.Exec(ctx, stmt.String())
			if err != nil {
				return err
			}
		}

		stmt.Reset()
		stmt.WriteString(`COPY system.assets 
						(instance_id, asset_type, resource_owner, 
						name, content_type, data, updated_at) 
						FROM STDIN`)

		tag, err := conn.PgConn().CopyFrom(ctx, reader, stmt.String())
		eventCount = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy assets to destination")
	logging.OnError(<-errs).Fatal("unable to copy assets from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("assets migrated")
}

func copyAssetsFromFile(ctx context.Context, dest *db.DB, fileName string) {
	start := time.Now()

	srcFile, err := os.OpenFile(filePath+fileName, os.O_RDONLY, 0)
	logging.OnError(err).Fatal("unable to open source file")
	defer srcFile.Close()

	reader, writer := io.Pipe()
	errs := make(chan error, 1)

	go func() {
		_, err := srcFile.WriteTo(writer)
		writer.Close()
		errs <- err
	}()

	destConn, err := dest.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire destination connection")
	defer destConn.Close()

	var eventCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		var stmt database.Statement
		stmt.WriteString("DELETE FROM system.assets ")
		stmt.WriteString(instanceClause())

		_, err := conn.Exec(ctx, stmt.String())
		if err != nil {
			return err
		}

		stmt.Reset()
		stmt.WriteString(`COPY system.assets 
						(instance_id, asset_type, resource_owner, 
						name, content_type, data, updated_at) 
						FROM STDIN (DELIMITER ',')`)

		tag, err := conn.PgConn().CopyFrom(ctx, reader, stmt.String())
		eventCount = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy assets to destination")
	logging.OnError(<-errs).Fatal("unable to copy assets from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("assets migrated from " + fileName)
}

func copyAssetsToFile(ctx context.Context, source *db.DB, fileName string) {
	start := time.Now()

	sourceConn, err := source.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire source connection")
	defer sourceConn.Close()

	reader, writer := io.Pipe()
	errs := make(chan error, 1)

	var eventCount int64
	go func() {
		err = sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()

			// ignore hash column because it's computed
			var stmt database.Statement
			stmt.WriteString(`COPY (SELECT instance_id, asset_type, 
							resource_owner, name, content_type, data, updated_at 
							FROM system.assets `)
			stmt.WriteString(instanceClause())
			stmt.WriteString(") TO STDOUT (DELIMITER ',')")

			tag, err := conn.PgConn().CopyTo(ctx, writer, stmt.String())
			eventCount = tag.RowsAffected()
			writer.Close()

			return err
		})
		errs <- err
	}()

	destFile, err := os.OpenFile(filePath+fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	logging.OnError(err).Fatal("unable to open destination file")
	defer destFile.Close()

	_, err = io.Copy(destFile, reader)
	logging.OnError(err).Fatal("unable to copy assets to destination")
	logging.OnError(<-errs).Fatal("unable to copy assets from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("assets copied to " + fileName)
}

func copyEncryptionKeysDB(ctx context.Context, source, dest *db.DB) {
	start := time.Now()

	sourceConn, err := source.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire source connection")
	defer sourceConn.Close()

	reader, writer := io.Pipe()
	errs := make(chan error, 1)

	go func() {
		err = sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()

			var stmt database.Statement
			stmt.WriteString("COPY system.encryption_keys TO STDOUT")

			_, err := conn.PgConn().CopyTo(ctx, writer, stmt.String())
			writer.Close()

			return err
		})
		errs <- err
	}()

	destConn, err := dest.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire destination connection")
	defer destConn.Close()

	var eventCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		var stmt database.Statement
		if shouldReplace {
			stmt.WriteString("TRUNCATE system.encryption_keys")

			_, err := conn.Exec(ctx, stmt.String())
			if err != nil {
				return err
			}
		}

		stmt.Reset()
		stmt.WriteString("COPY system.encryption_keys FROM STDIN")

		tag, err := conn.PgConn().CopyFrom(ctx, reader, stmt.String())
		eventCount = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy encryption keys to destination")
	logging.OnError(<-errs).Fatal("unable to copy encryption keys from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("encryption keys migrated")
}

func copyEncryptionKeysFromFile(ctx context.Context, dest *db.DB, fileName string) {
	start := time.Now()

	srcFile, err := os.OpenFile(filePath+fileName, os.O_RDONLY, 0)
	logging.OnError(err).Fatal("unable to open source file")
	defer srcFile.Close()

	reader, writer := io.Pipe()
	errs := make(chan error, 1)

	go func() {
		_, err := srcFile.WriteTo(writer)
		writer.Close()
		errs <- err
	}()

	destConn, err := dest.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire destination connection")
	defer destConn.Close()

	var eventCount int64
	err = destConn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		var stmt database.Statement
		stmt.WriteString("TRUNCATE system.encryption_keys")

		_, err := conn.Exec(ctx, stmt.String())
		if err != nil {
			return err
		}

		stmt.Reset()
		stmt.WriteString("COPY system.encryption_keys FROM STDIN (DELIMITER ',')")

		tag, err := conn.PgConn().CopyFrom(ctx, reader, stmt.String())
		eventCount = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy encryption keys to destination")
	logging.OnError(<-errs).Fatal("unable to copy encryption keys from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("encryption keys migrated from " + fileName)
}

func copyEncryptionKeysToFile(ctx context.Context, source *db.DB, fileName string) {
	start := time.Now()

	sourceConn, err := source.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire source connection")
	defer sourceConn.Close()

	reader, writer := io.Pipe()
	errs := make(chan error, 1)

	var eventCount int64
	go func() {
		err = sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()

			var stmt database.Statement
			stmt.WriteString("COPY system.encryption_keys TO STDOUT (DELIMITER ',')")

			tag, err := conn.PgConn().CopyTo(ctx, writer, stmt.String())
			eventCount = tag.RowsAffected()
			writer.Close()

			return err
		})
		errs <- err
	}()

	destFile, err := os.OpenFile(filePath+fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	logging.OnError(err).Fatal("unable to open destination file")
	defer destFile.Close()

	_, err = io.Copy(destFile, reader)
	logging.OnError(err).Fatal("unable to copy encryption keys to destination")
	logging.OnError(<-errs).Fatal("unable to copy encryption keys from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("encryption keys copied to " + fileName)
}
