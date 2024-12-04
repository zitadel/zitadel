package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mitchellh/mapstructure"
	"github.com/zitadel/logging"

	_ "github.com/zitadel/zitadel/internal/database/cockroach"
	"github.com/zitadel/zitadel/internal/database/dialect"
	_ "github.com/zitadel/zitadel/internal/database/postgres"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ContextQuerier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type ContextExecuter interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type ContextQueryExecuter interface {
	ContextQuerier
	ContextExecuter
}

type Client interface {
	ContextQueryExecuter
	Beginner
	Conn(ctx context.Context) (*sql.Conn, error)
}

type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Tx interface {
	ContextQueryExecuter
	Commit() error
	Rollback() error
}

var (
	_ Client = (*sql.DB)(nil)
	_ Tx     = (*sql.Tx)(nil)
)

func CloseTransaction(tx Tx, err error) error {
	if err != nil {
		rollbackErr := tx.Rollback()
		logging.OnError(rollbackErr).Error("failed to rollback transaction")
		return err
	}

	commitErr := tx.Commit()
	logging.OnError(commitErr).Error("failed to commit transaction")
	return commitErr
}

type Config struct {
	Dialects                   map[string]interface{} `mapstructure:",remain"`
	EventPushConnRatio         float64
	ProjectionSpoolerConnRatio float64
	connector                  dialect.Connector
}

func (c *Config) SetConnector(connector dialect.Connector) {
	c.connector = connector
}

type DB struct {
	*sql.DB
	dialect.Database
	Pool *pgxpool.Pool
}

func (db *DB) Query(scan func(*sql.Rows) error, query string, args ...any) error {
	return db.QueryContext(context.Background(), scan, query, args...)
}

func (db *DB) QueryContext(ctx context.Context, scan func(rows *sql.Rows) error, query string, args ...any) (err error) {
	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := rows.Close()
		logging.OnError(closeErr).Info("rows.Close failed")
	}()

	if err = scan(rows); err != nil {
		return err
	}
	return rows.Err()
}

func (db *DB) QueryRow(scan func(*sql.Row) error, query string, args ...any) (err error) {
	return db.QueryRowContext(context.Background(), scan, query, args...)
}

func (db *DB) QueryRowContext(ctx context.Context, scan func(row *sql.Row) error, query string, args ...any) (err error) {
	row := db.DB.QueryRowContext(ctx, query, args...)
	logging.OnError(row.Err()).Error("unexpected query error")

	err = scan(row)
	if err != nil {
		return err
	}
	return row.Err()
}

func QueryJSONObject[T any](ctx context.Context, db *DB, query string, args ...any) (*T, error) {
	var data []byte
	err := db.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&data)
	}, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DATAB-Oath6", "Errors.Internal")
	}
	obj := new(T)
	if err = json.Unmarshal(data, obj); err != nil {
		return nil, zerrors.ThrowInternal(err, "DATAB-Vohs6", "Errors.Internal")
	}
	return obj, nil
}

func Connect(config Config, useAdmin bool, purpose dialect.DBPurpose) (*DB, error) {
	client, pool, err := config.connector.Connect(useAdmin, config.EventPushConnRatio, config.ProjectionSpoolerConnRatio, purpose)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(); err != nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "DATAB-0pIWD", "Errors.Database.Connection.Failed")
	}

	return &DB{
		DB:       client,
		Database: config.connector,
		Pool:     pool,
	}, nil
}

func DecodeHook(from, to reflect.Value) (_ interface{}, err error) {
	if to.Type() != reflect.TypeOf(Config{}) {
		return from.Interface(), nil
	}

	config := new(Config)
	if err = mapstructure.Decode(from.Interface(), config); err != nil {
		return nil, err
	}

	configuredDialect := dialect.SelectByConfig(config.Dialects)
	configs := make([]interface{}, 0, len(config.Dialects)-1)

	for name, dialectConfig := range config.Dialects {
		if !configuredDialect.Matcher.MatchName(name) {
			continue
		}

		configs = append(configs, dialectConfig)
	}

	config.connector, err = configuredDialect.Matcher.Decode(configs)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c Config) DatabaseName() string {
	return c.connector.DatabaseName()
}

func (c Config) Username() string {
	return c.connector.Username()
}

func (c Config) Password() string {
	return c.connector.Password()
}

func (c Config) Type() string {
	return c.connector.Type()
}

func EscapeLikeWildcards(value string) string {
	value = strings.ReplaceAll(value, "%", "\\%")
	value = strings.ReplaceAll(value, "_", "\\_")
	return value
}
