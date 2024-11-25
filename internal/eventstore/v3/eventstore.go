package eventstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	dialect.RegisterAfterConnect(registerEventstoreTypes)
}

var (
	// pushPlaceholderFmt defines how data are inserted into the events table
	pushPlaceholderFmt string
	// uniqueConstraintPlaceholderFmt defines the format of the unique constraint error returned from the database
	uniqueConstraintPlaceholderFmt string

	_ eventstore.Pusher = (*Eventstore)(nil)
)

type Eventstore struct {
	client *database.DB
}

func registerEventstoreTypes(ctx context.Context, conn *pgx.Conn) error {
	types, err := conn.LoadTypes(ctx, []string{
		"eventstore._command",
		"eventstore.command",
		"eventstore.push",
	})
	if err != nil {
		logging.WithError(err).Debug("unable to load types")
		return nil
	}
	m := conn.TypeMap()

	m.RegisterTypes(types)
	dialect.RegisterDefaultPgTypeVariants[command](m, "eventstore.command", "eventstore._command")

	return nil
}

// Client implements the [eventstore.Pusher]
func (es *Eventstore) Client() *database.DB {
	return es.client
}

func NewEventstore(client *database.DB) *Eventstore {
	switch client.Type() {
	case "cockroach":
		pushPlaceholderFmt = "($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, hlc_to_timestamp(cluster_logical_timestamp()), cluster_logical_timestamp(), $%d)"
		uniqueConstraintPlaceholderFmt = "('%s', '%s', '%s')"
	case "postgres":
		pushPlaceholderFmt = "($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp()), $%d)"
		uniqueConstraintPlaceholderFmt = "(%s, %s, %s)"
	}

	return &Eventstore{client: client}
}

func (es *Eventstore) Health(ctx context.Context) error {
	return es.client.PingContext(ctx)
}

var errTypesNotFound = errors.New("types not found")

func checkExecutionPlan(ctx context.Context, conn *sql.Conn) error {
	return conn.Raw(func(driverConn any) error {
		conn, ok := driverConn.(*stdlib.Conn)
		if !ok {
			return errTypesNotFound
		}

		var cmd *command
		if _, ok := conn.Conn().TypeMap().TypeForValue(cmd); ok {
			return nil
		}
		return registerEventstoreTypes(ctx, conn.Conn())
	})
}

func (es *Eventstore) pushTx(ctx context.Context, client database.ContextQueryExecuter) (tx database.Tx, deferrable func(err error) error, err error) {
	var beginner database.Beginner
	switch c := client.(type) {
	case database.Tx:
		return c, nil, nil
	case database.Client:
		beginner = c
	default:
		beginner = es.client
	}
	isolationLevel := sql.LevelReadCommitted
	// cockroach requires serializable to execute the push function
	// because we use [cluster_logical_timestamp()](https://www.cockroachlabs.com/docs/stable/functions-and-operators#system-info-functions)
	if es.client.Type() == "cockroach" {
		isolationLevel = sql.LevelSerializable
	}
	tx, err = beginner.BeginTx(ctx, &sql.TxOptions{
		Isolation: isolationLevel,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, nil, err
	}
	return tx, func(err error) error { return database.CloseTransaction(tx, err) }, nil
}
