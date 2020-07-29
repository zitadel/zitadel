package model

import (
	"time"
)

type View struct {
	Database         string
	ViewName         string
	CurrentSequence  uint64
	CurrentTimestamp time.Time
}
