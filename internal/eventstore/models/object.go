package models

import (
	"time"
)

type ObjectRoot struct {
	ID           string    `json:"-"`
	Sequence     uint64    `json:"-"`
	CreationDate time.Time `json:"-"`
	ChangeDate   time.Time `json:"-"`
}

func (o *ObjectRoot) AppendEvent(event *Event) {
	if o.ID == "" {
		o.ID = event.AggregateID
	}

	o.ChangeDate = event.CreationDate
	if event.PreviousSequence.Int64 == 0 {
		o.CreationDate = o.ChangeDate
	}

	o.Sequence = event.Sequence
}
