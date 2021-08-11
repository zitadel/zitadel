package view

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/keypair"
)

func KeyPairQuery(latestSequence uint64) *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(keypair.AggregateType).
		SequenceGreater(latestSequence).
		EventTypes(keypair.AddedEventType).
		Builder()
}
