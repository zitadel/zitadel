package query

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
)

func (q *QuerySide) UserEvents(ctx context.Context, orgID, userID string, sequence uint64) ([]eventstore.EventReader, error) {
	query := NewUserEventSearchQuery(userID, orgID, sequence)
	return q.eventstore.FilterEvents(ctx, query)
}
