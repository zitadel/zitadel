package repository

import (
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/storage/eventstore"
)

type eventStore struct {
	es *eventstore.Eventstore
}

func events(client database.Executor) *eventStore {
	return &eventStore{
		es: eventstore.New(client),
	}
}
