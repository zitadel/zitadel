package handler

import (
	"context"

	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore"
)

func Start(ctx context.Context, es *eventstore.Eventstore, dbConf types.SQL) error {
	sqlClient, err := dbConf.Start()
	if err != nil {
		return err
	}

	NewOrgHandler(ctx, es, sqlClient)
	return nil
}
