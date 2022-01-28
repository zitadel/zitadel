package eventstore

import (
	"github.com/caos/zitadel/internal/config/types"
	typesV2 "github.com/caos/zitadel/internal/config/v2/types"
	"github.com/caos/zitadel/internal/eventstore/repository/sql"
)

func Start(sqlConfig types.SQL) (*Eventstore, error) {
	sqlClient, err := sqlConfig.Start()
	if err != nil {
		return nil, err
	}

	return NewEventstore(sql.NewCRDB(sqlClient)), nil
}

func StartV2(sqlConfig typesV2.SQL) (*Eventstore, error) {
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

func StartWithUser2(baseConfig types.SQLBase2, userConfig types.SQLUser2) (*Eventstore, error) {
	sqlClient, err := userConfig.Start2(baseConfig)
	if err != nil {
		return nil, err
	}

	return NewEventstore(sql.NewCRDB(sqlClient)), nil
}
