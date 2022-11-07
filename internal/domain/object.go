package domain

import "time"

type ObjectDetails struct {
	Sequence      uint64
	EventDate     time.Time
	ResourceOwner string
}
