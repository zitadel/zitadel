package eventstore

import (
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
)

type EventType = models.EventType

type Event struct {
	//ID is set by eventstore
	creationDate     time.Time
	typ              EventType
	previousSequence uint64
	data             []byte
	modifierService  string
	modifierTenant   string
	modifierUser     string
	resourceOwner    string
	aggregateType    AggregateType
	aggregateID      string
	aggregateVersion Version
}

func eventData(i interface{}) ([]byte, error) {
	switch v := i.(type) {
	case []byte:
		return v, nil
	case map[string]interface{}, interface{}:
		bytes, err := json.Marshal(v)
		if err != nil {
			return nil, errors.ThrowInvalidArgument(err, "MODEL-s2fgE", "unable to marshal data")
		}
		return bytes, nil
	case nil:
		return nil, nil
	}

	return nil, errors.ThrowInvalidArgument(nil, "MODEL-y7Lg5", "data is not valid")
}
