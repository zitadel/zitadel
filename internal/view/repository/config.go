package repository

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/zitadel/zitadel/internal/config/types"
	"github.com/zitadel/zitadel/internal/errors"
)

type ViewConfig struct {
	SQL *types.SQL
}

func Start(conf ViewConfig) (*sql.DB, *gorm.DB, error) {
	sqlClient, err := conf.SQL.Start()
	if err != nil {
		return nil, nil, errors.ThrowPreconditionFailed(err, "SQL-9qBtr", "unable to open database connection")
	}

	client, err := gorm.Open("postgres", sqlClient)
	if err != nil {
		return nil, nil, err
	}
	return sqlClient, client, nil
}
