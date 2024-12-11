package eventstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sync"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	dialect.RegisterAfterConnect(RegisterEventstoreTypes)
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

var (
	textType = &pgtype.Type{
		Name:  "text",
		OID:   pgtype.TextOID,
		Codec: pgtype.TextCodec{},
	}
	commandType = &pgtype.Type{
		Codec: &pgtype.CompositeCodec{
			Fields: []pgtype.CompositeCodecField{
				{
					Name: "instance_id",
					Type: textType,
				},
				{
					Name: "aggregate_type",
					Type: textType,
				},
				{
					Name: "aggregate_id",
					Type: textType,
				},
				{
					Name: "command_type",
					Type: textType,
				},
				{
					Name: "revision",
					Type: &pgtype.Type{
						Name:  "int2",
						OID:   pgtype.Int2OID,
						Codec: pgtype.Int2Codec{},
					},
				},
				{
					Name: "payload",
					Type: &pgtype.Type{
						Name: "jsonb",
						OID:  pgtype.JSONBOID,
						Codec: &pgtype.JSONBCodec{
							Marshal:   json.Marshal,
							Unmarshal: json.Unmarshal,
						},
					},
				},
				{
					Name: "creator",
					Type: textType,
				},
				{
					Name: "owner",
					Type: textType,
				},
			},
		},
	}
	commandArrayCodec = &pgtype.Type{
		Codec: &pgtype.ArrayCodec{
			ElementType: commandType,
		},
	}
)

var typeMu sync.Mutex

func RegisterEventstoreTypes(ctx context.Context, conn *pgx.Conn) error {
	// conn.TypeMap is not thread safe
	typeMu.Lock()
	defer typeMu.Unlock()

	m := conn.TypeMap()

	var cmd *command
	if _, ok := m.TypeForValue(cmd); ok {
		return nil
	}

	if commandType.OID == 0 || commandArrayCodec.OID == 0 {
		err := conn.QueryRow(ctx, "select oid, typarray from pg_type where typname = $1 and typnamespace = (select oid from pg_namespace where nspname = $2)", "command", "eventstore").
			Scan(&commandType.OID, &commandArrayCodec.OID)
		if err != nil {
			logging.WithError(err).Debug("failed to get oid for command type")
			return nil
		}
		if commandType.OID == 0 || commandArrayCodec.OID == 0 {
			logging.Debug("oid for command type not found")
			return nil
		}
	}

	m.RegisterTypes([]*pgtype.Type{
		{
			Name:  "eventstore.command",
			Codec: commandType.Codec,
			OID:   commandType.OID,
		},
		{
			Name:  "command",
			Codec: commandType.Codec,
			OID:   commandType.OID,
		},
		{
			Name:  "eventstore._command",
			Codec: commandArrayCodec.Codec,
			OID:   commandArrayCodec.OID,
		},
		{
			Name:  "_command",
			Codec: commandArrayCodec.Codec,
			OID:   commandArrayCodec.OID,
		},
	})
	dialect.RegisterDefaultPgTypeVariants[command](m, "eventstore.command", "eventstore._command")
	dialect.RegisterDefaultPgTypeVariants[command](m, "command", "_command")

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

func CheckExecutionPlan(ctx context.Context, conn *sql.Conn) error {
	return conn.Raw(func(driverConn any) error {
		if _, ok := driverConn.(sqlmock.SqlmockCommon); ok {
			return nil
		}
		conn, ok := driverConn.(*stdlib.Conn)
		if !ok {
			return errTypesNotFound
		}

		return RegisterEventstoreTypes(ctx, conn.Conn())
	})
}

func (es *Eventstore) pushTx(ctx context.Context, client database.ContextQueryExecuter) (tx database.Tx, deferrable func(err error) error, err error) {
	tx, ok := client.(database.Tx)
	if ok {
		return tx, nil, nil
	}
	beginner, ok := client.(database.Beginner)
	if !ok {
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
