package postgres

import (
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type intent struct {
	eventstore.PushIntent

	sequence uint32
}

func makeIntents(pushIntents []eventstore.PushIntent) []*intent {
	res := make([]*intent, len(pushIntents))

	for i, pushIntent := range pushIntents {
		res[i] = &intent{PushIntent: pushIntent}
	}

	return res
}

func intentByAggregate(intents []*intent, aggregate *eventstore.Aggregate) *intent {
	for _, intent := range intents {
		if intent.Aggregate().Equals(aggregate) {
			return intent
		}
	}
	logging.WithFields("instance", aggregate.Instance, "owner", aggregate.Owner, "type", aggregate.Type, "id", aggregate.ID).Panic("no intent found")
	return nil
}

func checkSequences(intents []*intent) bool {
	for _, intent := range intents {
		if !eventstore.CheckSequence(intent.sequence, intent.CurrentSequence()) {
			return false
		}
	}
	return true
}
