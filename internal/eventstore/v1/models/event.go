package models

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type EventType string

func (et EventType) String() string {
	return string(et)
}

var _ eventstore.Event = (*Event)(nil)

type Event struct {
	ID               string
	Seq              uint64
	Pos              decimal.Decimal
	CreationDate     time.Time
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
func (e *Event) Aggregate() *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            e.AggregateID,
		Type:          e.AggregateType,
		ResourceOwner: e.ResourceOwner,
		InstanceID:    e.InstanceID,
		// Version:       eventstore.Version(e.AggregateVersion),
	}
}

// CreatedAt implements [eventstore.Event]
func (e *Event) CreatedAt() time.Time {
	return e.CreationDate
}

// DataAsBytes implements [eventstore.Event]
func (e *Event) DataAsBytes() []byte {
	return e.Data
}

// Unmarshal implements [eventstore.Event]
func (e *Event) Unmarshal(ptr any) error {
	if len(e.Data) == 0 {
		return nil
	}
	return json.Unmarshal(e.Data, ptr)
}

// EditorService implements [eventstore.Event]
func (e *Event) EditorService() string {
	return e.Service
}

// Creator implements [eventstore.action]
func (e *Event) Creator() string {
	return e.User
}

// Sequence implements [eventstore.Event]
func (e *Event) Sequence() uint64 {
	return e.Seq
}

// Position implements [eventstore.Event]
func (e *Event) Position() decimal.Decimal {
	return e.Pos
}

// Type implements [eventstore.action]
func (e *Event) Type() eventstore.EventType {
	return e.Typ
}

// Type implements [eventstore.action]
func (e *Event) Revision() uint16 {
	return 0
}

func eventData(i interface{}) ([]byte, error) {
	switch v := i.(type) {
	case []byte:
		return v, nil
	case map[string]interface{}:
		bytes, err := json.Marshal(v)
		if err != nil {
			return nil, zerrors.ThrowInvalidArgument(err, "MODEL-s2fgE", "unable to marshal data")
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
			return nil, zerrors.ThrowInvalidArgument(nil, "MODEL-rjWdN", "data is not valid")
		}
		bytes, err := json.Marshal(v)
		if err != nil {
			return nil, zerrors.ThrowInvalidArgument(err, "MODEL-Y2OpM", "unable to marshal data")
		}
		return bytes, nil
	}
}

func (e *Event) Validate() error {
	if e == nil {
		return zerrors.ThrowPreconditionFailed(nil, "MODEL-oEAG4", "event is nil")
	}
	if string(e.Typ) == "" {
		return zerrors.ThrowPreconditionFailed(nil, "MODEL-R2sB0", "type not defined")
	}

	if e.AggregateID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "MODEL-A6WwL", "aggregate id not set")
	}
	if e.AggregateType == "" {
		return zerrors.ThrowPreconditionFailed(nil, "MODEL-EzdyK", "aggregate type not set")
	}
	if err := e.AggregateVersion.Validate(); err != nil {
		return zerrors.ThrowPreconditionFailed(err, "MODEL-KO71q", "version invalid")
	}
	if e.Service == "" {
		return zerrors.ThrowPreconditionFailed(nil, "MODEL-4Yqik", "editor service not set")
	}
	if e.User == "" {
		return zerrors.ThrowPreconditionFailed(nil, "MODEL-L3NHO", "editor user not set")
	}
	if e.ResourceOwner == "" {
		return zerrors.ThrowPreconditionFailed(nil, "MODEL-omFVT", "resource ow")
	}
	return nil
}
