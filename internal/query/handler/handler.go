package handler

import (
	"context"

	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore"
)

func StartWithUser(ctx context.Context, es *eventstore.Eventstore, baseConfig types.SQLBase, userConfig types.SQLUser) error {
	sqlClient, err := userConfig.Start(baseConfig)
	if err != nil {
		return err
	}

	NewOrgHandler(ctx, es, sqlClient)
	return nil
}
