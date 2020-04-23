package models

import (
	"time"
)

type ObjectRoot struct {
	AggregateID  string    `json:"-"`
	Sequence     uint64    `json:"-"`
	CreationDate time.Time `json:"-"`
	ChangeDate   time.Time `json:"-"`
}

func (o *ObjectRoot) AppendEvent(event *Event) {
	if o.AggregateID == "" {
		o.AggregateID = event.AggregateID
	}

	o.ChangeDate = event.CreationDate
	if event.PreviousSequence == 0 {
		o.CreationDate = o.ChangeDate
	}

	o.Sequence = event.Sequence
}
