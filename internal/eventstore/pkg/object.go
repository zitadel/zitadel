package pkg

import (
	"time"

	es_api "github.com/caos/citadel/eventstore/api/grpc"
	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"
)

type ObjectRoot struct {
	ID           string    `json="-"`
	Sequence     uint64    `json="-"`
	CreationDate time.Time `json="-"`
	ChangeDate   time.Time `json="-"`
}

func (o *ObjectRoot) AppendEvent(event *es_api.EventResponse) {
	if o.ID == "" {
		o.ID = event.AggregateId
	}

	var err error
	o.ChangeDate, err = ptypes.Timestamp(event.CreationDate)
	logging.Log("MODEL-Oqzfc").OnError(err).Debug("unvalid event creationdate")
	if event.PreviousSequence == 0 {
		o.CreationDate = o.ChangeDate
	}

	o.Sequence = event.Sequence
}
