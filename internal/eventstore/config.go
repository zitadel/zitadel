package eventstore

import (
	"github.com/zitadel/zitadel/internal/config/types"
	"github.com/zitadel/zitadel/internal/eventstore/repository/sql"
)

func Start(sqlConfig types.SQL) (*Eventstore, error) {
	sqlClient, err := sqlConfig.Start()
	if err != nil {
		return nil, err
	}

	return NewEventstore(sql.NewCRDB(sqlClient)), nil
}

func StartWithUser(baseConfig types.SQLBase, userConfig types.SQLUser) (*Eventstore, error) {
	sqlClient, err := userConfig.Start(baseConfig)
	if err != nil {
		return nil, err
	}

	return NewEventstore(sql.NewCRDB(sqlClient)), nil
}
