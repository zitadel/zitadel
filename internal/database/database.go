package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/zitadel/logging"

	_ "github.com/zitadel/zitadel/internal/database/cockroach"
	"github.com/zitadel/zitadel/internal/database/dialect"
	_ "github.com/zitadel/zitadel/internal/database/postgres"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Config struct {
	Dialects           map[string]interface{} `mapstructure:",remain"`
	EventPushConnRatio float32
	connector          dialect.Connector
}

func (c *Config) SetConnector(connector dialect.Connector) {
	c.connector = connector
}

type DB struct {
	*sql.DB
	dialect.Database
}

func (db *DB) Query(scan func(*sql.Rows) error, query string, args ...any) error {
	return db.QueryContext(context.Background(), scan, query, args...)
}

func (db *DB) QueryContext(ctx context.Context, scan func(rows *sql.Rows) error, query string, args ...any) (err error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			logging.OnError(rollbackErr).Info("commit of read only transaction failed")
			return
		}
		err = tx.Commit()
	}()

	rows, err := tx.QueryContext(ctx, query, args...)
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
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			logging.OnError(rollbackErr).Info("commit of read only transaction failed")
			return
		}
		err = tx.Commit()
	}()

	row := tx.QueryRowContext(ctx, query, args...)

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

const (
	zitadelAppName          = "zitadel"
	EventstorePusherAppName = "zitadel_es_pusher"
)

func Connect(config Config, useAdmin, isEventPusher bool) (*DB, error) {
	appName := zitadelAppName
	if isEventPusher {
		appName = EventstorePusherAppName
	}

	client, err := config.connector.Connect(useAdmin, isEventPusher, config.EventPushConnRatio, appName)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(); err != nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "DATAB-0pIWD", "Errors.Database.Connection.Failed")
	}

	return &DB{
		DB:       client,
		Database: config.connector,
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
