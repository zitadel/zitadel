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

func authCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "mirrors the auth requests table between databases, or between a database and files",
		Long: `mirrors the auth requests table between databases, or between a database and files
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
	switch {
	case isSrcFile:
		destClient, err := db.Connect(config.Destination, false, dialect.DBPurposeEventPusher)
		logging.OnError(err).Fatal("unable to connect to destination database")
		defer destClient.Close()

		copyAuthRequestsFromFile(ctx, destClient, "auth.auth_requests.csv")
	case isDestFile:
		sourceClient, err := db.Connect(config.Source, false, dialect.DBPurposeQuery)
		logging.OnError(err).Fatal("unable to connect to source database")
		defer sourceClient.Close()

		copyAuthRequestsToFile(ctx, sourceClient, "auth.auth_requests.csv")
	default:
		sourceClient, err := db.Connect(config.Source, false, dialect.DBPurposeQuery)
		logging.OnError(err).Fatal("unable to connect to source database")
		defer sourceClient.Close()

		destClient, err := db.Connect(config.Destination, false, dialect.DBPurposeEventPusher)
		logging.OnError(err).Fatal("unable to connect to destination database")
		defer destClient.Close()

		copyAuthRequestsDB(ctx, sourceClient, destClient)
	}
}

func copyAuthRequestsDB(ctx context.Context, source, dest *db.DB) {
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
			stmt.WriteString(`COPY (SELECT id, 
							regexp_replace(request::TEXT, '\\\\u0000', '', 'g')::JSON 
							request, code, request_type, creation_date, change_date, instance_id 
							FROM auth.auth_requests `)
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
		if shouldReplace {
			stmt.WriteString("DELETE FROM auth.auth_requests ")
			stmt.WriteString(instanceClause())

			_, err := conn.Exec(ctx, stmt.String())
			if err != nil {
				return err
			}
		}

		stmt.Reset()
		stmt.WriteString("COPY auth.auth_requests FROM STDIN")

		tag, err := conn.PgConn().CopyFrom(ctx, reader, stmt.String())
		eventCount = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy auth requests to destination")
	logging.OnError(<-errs).Fatal("unable to copy auth requests from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("auth requests migrated")
}

func copyAuthRequestsFromFile(ctx context.Context, dest *db.DB, fileName string) {
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
		stmt.WriteString("DELETE FROM auth.auth_requests ")
		stmt.WriteString(instanceClause())

		_, err := conn.Exec(ctx, stmt.String())
		if err != nil {
			return err
		}

		stmt.Reset()
		stmt.WriteString(`COPY auth.auth_requests
						 (id, request, code, request_type,
						 creation_date, change_date, instance_id)
						 FROM STDIN (DELIMITER ',')`)

		tag, err := conn.PgConn().CopyFrom(ctx, reader, stmt.String())
		eventCount = tag.RowsAffected()

		return err
	})
	logging.OnError(err).Fatal("unable to copy auth requests to destination")
	logging.OnError(<-errs).Fatal("unable to copy auth requests from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("auth requests migrated from " + fileName)
}

func copyAuthRequestsToFile(ctx context.Context, src *db.DB, fileName string) {
	start := time.Now()

	sourceConn, err := src.Conn(ctx)
	logging.OnError(err).Fatal("unable to acquire source connection")
	defer sourceConn.Close()

	reader, writer := io.Pipe()
	errs := make(chan error, 1)

	var eventCount int64
	go func() {
		err = sourceConn.Raw(func(driverConn interface{}) error {
			conn := driverConn.(*stdlib.Conn).Conn()

			var stmt database.Statement
			stmt.WriteString(`COPY (SELECT id, 
							regexp_replace(request::TEXT, '\\\\u0000', '', 'g')::JSON 
							request, code, request_type, creation_date, change_date, 
							instance_id FROM auth.auth_requests `)
			stmt.WriteString(instanceClause())
			stmt.WriteString(") TO STDOUT (DELIMITER ',')")

			tag, err := conn.PgConn().CopyTo(ctx, writer, stmt.String())
			eventCount = tag.RowsAffected()
			writer.Close()

			return err
		})
		errs <- err
	}()

	destFile, err := os.OpenFile(filePath+fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	logging.OnError(err).Fatal("unable to open destination file")
	defer destFile.Close()

	_, err = io.Copy(destFile, reader)
	logging.OnError(err).Fatal("unable to copy auth requests to file")
	logging.OnError(<-errs).Fatal("unable to copy auth requests from source")
	logging.WithFields("took", time.Since(start), "count", eventCount).Info("auth requests copied to " + fileName)
}
