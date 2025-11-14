package domain

import (
	"time"
)

type ObjectDetails struct {
	Sequence uint64
	// EventDate is the date of the last event that changed the object
	EventDate time.Time
	// CreationDate is the date of the first event that created the object
	CreationDate  time.Time
	ResourceOwner string
	// ID is the Aggregate ID of the resource
	ID string
}
