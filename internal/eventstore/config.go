package eventstore

import (
	"database/sql"

	"github.com/caos/zitadel/internal/config/types"
	z_sql "github.com/caos/zitadel/internal/eventstore/repository/sql"
)

func Start(sqlClient *sql.DB) (*Eventstore, error) {
	return NewEventstore(z_sql.NewCRDB(sqlClient)), nil
}

func StartWithUser(baseConfig types.SQLBase, userConfig types.SQLUser) (*Eventstore, error) {
	sqlClient, err := userConfig.Start(baseConfig)
	if err != nil {
		return nil, err
	}

	return NewEventstore(z_sql.NewCRDB(sqlClient)), nil
}
