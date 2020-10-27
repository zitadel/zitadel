package eventstore

import (
	"time"

	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

//Event is the representation of a state change
type Event interface {
	//CheckPrevious ensures the event order if true
	// if false the previous sequence is not checked on push
	CheckPrevious() bool
	//EditorService must return the name of the service which creates the new event
	EditorService() string
	//EditorUser must return the id of the user who created the event
	EditorUser() string
	//Type must return an event type which should be unique in the aggregate
	Type() EventType
	//Data returns the payload of the event. It represent the changed fields by the event
	// valid types are:
	// * nil (no payload),
	// * json byte array
	// * struct which can be marshalled to json
	// * pointer to struct which can be marshalled to json
	Data() interface{}
	//MetaData returns all data saved on a event
	// It must not be set on push
	// The event mapper function must set this struct
	MetaData() *EventMetaData
}

func MetaDataFromRepo(event *repository.Event) *EventMetaData {
	return &EventMetaData{
		AggregateID:       event.AggregateID,
		AggregateType:     AggregateType(event.AggregateType),
		AggregateVersion:  Version(event.Version),
		PreviouseSequence: event.PreviousSequence,
		ResourceOwner:     event.ResourceOwner,
		Sequence:          event.Sequence,
		CreationDate:      event.CreationDate,
	}
}

type EventMetaData struct {
	AggregateID       string
	AggregateType     AggregateType
	ResourceOwner     string
	AggregateVersion  Version
	Sequence          uint64
	PreviouseSequence uint64
	CreationDate      time.Time
}
