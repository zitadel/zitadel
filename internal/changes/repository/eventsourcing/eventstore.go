package eventsourcing

import (
	es_int "github.com/caos/zitadel/internal/eventstore"
)

type ChangesEventstore struct {
	es_int.Eventstore
}

type ChangesConfig struct {
	es_int.Eventstore
}

func StartChanges(conf ChangesConfig) (*ChangesEventstore, error) {
	return &ChangesEventstore{
		Eventstore: conf.Eventstore,
	}, nil
}
