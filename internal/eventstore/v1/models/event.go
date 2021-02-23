package models

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/caos/zitadel/internal/errors"
)

type EventType string

func (et EventType) String() string {
	return string(et)
}

type Event struct {
	ID               string
	Sequence         uint64
	CreationDate     time.Time
	Type             EventType
	PreviousSequence uint64
	Data             []byte

	AggregateID      string
	AggregateType    AggregateType
	AggregateVersion Version
	EditorService    string
	EditorUser       string
	ResourceOwner    string
}

func eventData(i interface{}) ([]byte, error) {
	switch v := i.(type) {
	case []byte:
		return v, nil
	case map[string]interface{}:
		bytes, err := json.Marshal(v)
		if err != nil {
			return nil, errors.ThrowInvalidArgument(err, "MODEL-s2fgE", "unable to marshal data")
		}
		return bytes, nil
	case nil:
		return nil, nil
	default:
		t := reflect.TypeOf(i)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			return nil, errors.ThrowInvalidArgument(nil, "MODEL-rjWdN", "data is not valid")
		}
		bytes, err := json.Marshal(v)
		if err != nil {
			return nil, errors.ThrowInvalidArgument(err, "MODEL-Y2OpM", "unable to marshal data")
		}
		return bytes, nil
	}
}

func (e *Event) Validate() error {
	if e == nil {
		return errors.ThrowPreconditionFailed(nil, "MODEL-oEAG4", "event is nil")
	}
	if string(e.Type) == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-R2sB0", "type not defined")
	}

	if e.AggregateID == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-A6WwL", "aggregate id not set")
	}
	if e.AggregateType == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-EzdyK", "aggregate type not set")
	}
	if err := e.AggregateVersion.Validate(); err != nil {
		return err
	}

	if e.EditorService == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-4Yqik", "editor service not set")
	}
	if e.EditorUser == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-L3NHO", "editor user not set")
	}
	if e.ResourceOwner == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-omFVT", "resource ow")
	}
	return nil
}
