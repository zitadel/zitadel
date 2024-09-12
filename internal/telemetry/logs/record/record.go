package record

import (
	"github.com/sirupsen/logrus"
)

type Stream string

const (
	StreamActivity Stream = "activity"
)

type BaseStreamRecord struct {
	Version           string
	Stream            Stream
	ZITADELInstanceID string
}

type StreamRecord interface {
	Base() *BaseStreamRecord
	Fields() logrus.Fields
}
