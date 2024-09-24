package record

import (
	"github.com/sirupsen/logrus"

	"github.com/zitadel/zitadel/pkg/streams"
)

type BaseStreamRecord struct {
	Version           string
	Stream            streams.Stream
	ZITADELInstanceID string
}

type StreamRecord interface {
	Base() *BaseStreamRecord
	Fields() logrus.Fields
}
