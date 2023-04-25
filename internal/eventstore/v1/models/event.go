package models

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type EventType string

func (et EventType) String() string {
	return string(et)
}

// var _ eventstore.Event = (*Event)(nil)

type Event struct {
	ID               string
	Seq              uint64
	CreatedAt        time.Time
	Typ              eventstore.EventType
	PreviousSequence uint64
	Data             []byte

	AggregateID      string
	AggregateType    eventstore.AggregateType
	AggregateVersion eventstore.Version
	Service          string
	User             string
	ResourceOwner    string
	InstanceID       string
}

// Aggregate implements [eventstore.Event]
func (e *Event) Aggregate() eventstore.Aggregate {
	return eventstore.Aggregate{
		ID:            e.AggregateID,
		Type:          eventstore.AggregateType(e.AggregateType),
		ResourceOwner: e.ResourceOwner,
		InstanceID:    e.InstanceID,
		Version:       eventstore.Version(e.AggregateVersion),
	}
}

// CreationDate implements [eventstore.Event]
func (e *Event) CreationDate() time.Time {
	return e.CreatedAt
}

// DataAsBytes implements [eventstore.Event]
func (e *Event) DataAsBytes() []byte {
	return e.Data
}

// EditorService implements [eventstore.Event]
func (e *Event) EditorService() string {
	return e.Service
}

// EditorUser implements [eventstore.Event]
func (e *Event) EditorUser() string {
	return e.User
}

// PreviousAggregateSequence implements [eventstore.Event]
func (e *Event) PreviousAggregateSequence() uint64 {
	return e.PreviousSequence
}

// PreviousAggregateTypeSequence implements [eventstore.Event]
func (e *Event) PreviousAggregateTypeSequence() uint64 {
	return e.PreviousSequence
}

// Sequence implements [eventstore.Event]
func (e *Event) Sequence() uint64 {
	return e.Seq
}

// Type implements [eventstore.Event]
func (e *Event) Type() eventstore.EventType {
	return eventstore.EventType(e.Typ)
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
	if string(e.Typ) == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-R2sB0", "type not defined")
	}

	if e.AggregateID == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-A6WwL", "aggregate id not set")
	}
	if e.AggregateType == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-EzdyK", "aggregate type not set")
	}
	// TODO: check still needed?
	// if err := e.AggregateVersion.Validate(); err != nil {
	// 	return err
	// }

	if e.Service == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-4Yqik", "editor service not set")
	}
	if e.User == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-L3NHO", "editor user not set")
	}
	if e.ResourceOwner == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-omFVT", "resource ow")
	}
	return nil
}
