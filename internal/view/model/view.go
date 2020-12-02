package model

import (
	"time"
)

type View struct {
	Database                 string
	ViewName                 string
	CurrentSequence          uint64
	EventTimestamp           time.Time
	LastSuccessfulSpoolerRun time.Time
}
