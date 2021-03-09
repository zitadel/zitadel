package domain

import "time"

type ObjectDetails struct {
	Sequence      uint64
	ChangeDate    time.Time
	ResourceOwner string
}
