package readmodel

import (
	"time"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type readModel struct {
	CreationDate time.Time
	ChangeDate   time.Time
	Instance     string
	Owner        string
	Sequence     uint32
}

func (rm *readModel) reduce(event *eventstore.Event[eventstore.StoragePayload]) {
	if rm.CreationDate.IsZero() {
		rm.CreationDate = event.CreatedAt
	}
	if rm.Instance == "" {
		rm.Instance = event.Aggregate.Instance
	}
	if rm.Owner == "" {
		rm.Owner = event.Aggregate.Owner
	}

	rm.ChangeDate = event.CreatedAt
	rm.Sequence = event.Sequence
}
