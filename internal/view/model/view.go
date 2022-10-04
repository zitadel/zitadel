package model

import (
	"time"
)

type View struct {
	Database                 string
	ViewName                 string
	EventTimestamp           time.Time
	EventID                  string
	LastSuccessfulSpoolerRun time.Time
}
