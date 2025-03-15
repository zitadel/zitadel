package event

import (
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/storage/eventstore"
)

type store struct {
	es *eventstore.Eventstore
}

func Store(client database.Executor) *store {
	return &store{
		es: eventstore.New(client),
	}
}
