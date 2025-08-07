package models

import (
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type ObjectRoot struct {
	AggregateID   string    `json:"-"`
	Sequence      uint64    `json:"-"`
	ResourceOwner string    `json:"-"`
	InstanceID    string    `json:"-"`
	CreationDate  time.Time `json:"-"`
	ChangeDate    time.Time `json:"-"`
}

func (o *ObjectRoot) AppendEvent(event eventstore.Event) {
	if o.AggregateID == "" {
		o.AggregateID = event.Aggregate().ID
	} else if o.AggregateID != event.Aggregate().ID {
		return
	}
	if o.ResourceOwner == "" {
		o.ResourceOwner = event.Aggregate().ResourceOwner
	}
	if o.InstanceID == "" {
		o.InstanceID = event.Aggregate().InstanceID
	}

	o.ChangeDate = event.CreatedAt()
	if o.CreationDate.IsZero() {
		o.CreationDate = o.ChangeDate
	}

	o.Sequence = event.Sequence()
}
func (o *ObjectRoot) IsZero() bool {
	return o.AggregateID == ""
}
