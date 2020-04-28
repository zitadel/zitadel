package view

import (
	"database/sql"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/errors"
	"github.com/jinzhu/gorm"
)

type ViewDB struct {
	SQL  *sql.DB
	GORM *gorm.DB
}

type ViewConfig struct {
	SQL *types.SQL
}

func Start(conf ViewConfig) (*ViewDB, error) {
	sqlClient, err := sql.Open("postgres", conf.SQL.ConnectionString())
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "SQL-9qBtr", "unable to open database connection")
	}

	client, err := gorm.Open("postgres", sqlClient)
	if err != nil {
		return nil, err
	}
	return &ViewDB{SQL: sqlClient, GORM: client}, nil
}
