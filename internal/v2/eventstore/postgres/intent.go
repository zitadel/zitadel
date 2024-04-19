package postgres

import (
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type intent struct {
	*eventstore.PushAggregate

	sequence uint32
}

func makeIntents(pushIntent *eventstore.PushIntent) []*intent {
	res := make([]*intent, len(pushIntent.Aggregates()))

	for i, aggregate := range pushIntent.Aggregates() {
		res[i] = &intent{PushAggregate: aggregate}
	}

	return res
}

func intentByAggregate(intents []*intent, aggregate *eventstore.Aggregate) *intent {
	for _, intent := range intents {
		if intent.PushAggregate.Aggregate().Equals(aggregate) {
			return intent
		}
	}
	logging.WithFields("instance", aggregate.Instance, "owner", aggregate.Owner, "type", aggregate.Type, "id", aggregate.ID).Panic("no intent found")
	return nil
}

func checkSequences(intents []*intent) bool {
	for _, intent := range intents {
		if !eventstore.CheckSequence(intent.sequence, intent.PushAggregate.CurrentSequence()) {
			return false
		}
	}
	return true
}
