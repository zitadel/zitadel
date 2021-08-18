package handler

import "github.com/caos/zitadel/internal/eventstore"

type EventstoreHandler interface {
	Handler
}

type eventstoreHandler struct {
	es eventstore.Eventstore
}
