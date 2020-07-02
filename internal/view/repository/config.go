package repository

import (
	"database/sql"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/errors"
	"github.com/jinzhu/gorm"
)

type ViewConfig struct {
	SQL *types.SQL
}

func Start(conf ViewConfig) (*sql.DB, *gorm.DB, error) {
	sqlClient, err := sql.Open("postgres", conf.SQL.ConnectionString())
	if err != nil {
		return nil, nil, errors.ThrowPreconditionFailed(err, "SQL-9qBtr", "unable to open database connection")
	}

	client, err := gorm.Open("postgres", sqlClient)
	if err != nil {
		return nil, nil, err
	}
	return sqlClient, client, nil
}
