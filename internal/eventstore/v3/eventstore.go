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

	new_db "github.com/zitadel/zitadel/backend/v3/storage/database"
	new_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	dialect.RegisterAfterConnect(RegisterEventstoreTypes)
}

var (
	// pushPlaceholderFmt defines how data are inserted into the events table
	pushPlaceholderFmt = "($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp()), $%d)"
	// uniqueConstraintPlaceholderFmt defines the format of the unique constraint error returned from the database
	uniqueConstraintPlaceholderFmt = "(%s, %s, %s)"

	_ eventstore.Pusher = (*Eventstore)(nil)
)

type Eventstore struct {
	client new_db.Pool
	queue  eventstore.ExecutionQueue
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

	command2Type = &pgtype.Type{
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
				{
					Name: "written_by_v3",
					Type: &pgtype.Type{
						Name:  "bool",
						OID:   pgtype.BoolOID,
						Codec: pgtype.BoolCodec{},
					},
				},
			},
		},
	}
	command2ArrayCodec = &pgtype.Type{
		Codec: &pgtype.ArrayCodec{
			ElementType: command2Type,
		},
	}
)

var typeMu sync.Mutex

func RegisterEventstoreTypes(ctx context.Context, conn *pgx.Conn) error {
	// conn.TypeMap is not thread safe
	typeMu.Lock()
	defer typeMu.Unlock()

	registerCommandType[command](ctx, conn, commandType, commandArrayCodec, "command")
	registerCommandType[command2](ctx, conn, command2Type, command2ArrayCodec, "command2")

	return nil
}

func registerCommandType[T interface{ command | command2 }](ctx context.Context, conn *pgx.Conn, typeCodec, arrayCodec *pgtype.Type, name string) {
	if typeCodec.OID == 0 || arrayCodec.OID == 0 {
		err := conn.QueryRow(ctx, "select oid, typarray from pg_type where typname = $1 and typnamespace = (select oid from pg_namespace where nspname = $2)", name, "eventstore").
			Scan(&typeCodec.OID, &arrayCodec.OID)
		if err != nil {
			logging.WithError(err).Debug("failed to get oid for command type")
			return
		}
		if typeCodec.OID == 0 || arrayCodec.OID == 0 {
			logging.Debug("oid for command type not found")
			return
		}
	}

	conn.TypeMap().RegisterTypes([]*pgtype.Type{
		{
			Name:  "eventstore." + name,
			Codec: typeCodec.Codec,
			OID:   typeCodec.OID,
		},
		{
			Name:  name,
			Codec: typeCodec.Codec,
			OID:   typeCodec.OID,
		},
		{
			Name:  "eventstore._" + name,
			Codec: arrayCodec.Codec,
			OID:   arrayCodec.OID,
		},
		{
			Name:  "_" + name,
			Codec: arrayCodec.Codec,
			OID:   arrayCodec.OID,
		},
	})

	dialect.RegisterDefaultPgTypeVariants[T](conn.TypeMap(), "eventstore."+name, "eventstore._"+name)
	dialect.RegisterDefaultPgTypeVariants[T](conn.TypeMap(), name, "_"+name)
}

func NewEventstore(client *database.DB, opts ...EventstoreOption) *Eventstore {
	es := &Eventstore{
		client: new_sql.SQLPool(client.DB),
	}
	for _, opt := range opts {
		opt(es)
	}
	return es
}

func (es *Eventstore) Health(ctx context.Context) error {
	return es.client.Ping(ctx)
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

func (es *Eventstore) pushTx(ctx context.Context, client new_db.QueryExecutor) (tx new_db.Transaction, deferrable func(err error) error, err error) {
	tx, ok := client.(new_db.Transaction)
	if ok {
		return tx, nil, nil
	}
	beginner, ok := client.(new_db.Beginner)
	if !ok {
		beginner = es.client
	}

	tx, err = beginner.Begin(ctx, &new_db.TransactionOptions{
		IsolationLevel: new_db.IsolationLevelReadCommitted,
		AccessMode:     new_db.AccessModeReadWrite,
	})
	if err != nil {
		return nil, nil, err
	}
	return tx, func(err error) error { return tx.End(ctx, err) }, nil
}

type EventstoreOption func(*Eventstore)

func WithExecutionQueueOption(queue eventstore.ExecutionQueue) EventstoreOption {
	return func(es *Eventstore) {
		es.queue = queue
	}
}
