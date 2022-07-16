package cockroach

import (
	"database/sql"

	//sql import
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/zitadel/zitadel/internal/database/dialect"
)

func init() {
	config := &Config{}
	dialect.Register(config, config, true)
}

func connect(config interface{}) (*sql.DB, error) {
	// config := new(Config)

	// // mapstructure.DecoderConfig

	// if err := v.Unmarshal(config); err != nil {
	// 	return nil, err
	// }

	// client, err := sql.Open("pgx", config.String())
	// if err != nil {
	// 	return nil, err
	// }

	// client.SetMaxOpenConns(int(config.MaxOpenConns))
	// client.SetConnMaxLifetime(config.MaxConnLifetime)
	// client.SetConnMaxIdleTime(config.MaxConnIdleTime)

	// if err := client.Ping(); err != nil {
	// 	return nil, errors.ThrowPreconditionFailed(err, "POSTG-LaDzr", "Errors.Database.Connection.Failed")
	// }

	// return client, nil

	return nil, nil
}
