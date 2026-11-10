package execution

import (
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/execution/target"
)

const (
	QueueName = "execution"
)

type Request struct {
	Aggregate   *eventstore.Aggregate `json:"aggregate"`
	Sequence    uint64                `json:"sequence"`
	EventType   eventstore.EventType  `json:"eventType"`
	CreatedAt   time.Time             `json:"createdAt"`
	UserID      string                `json:"userID"`
	EventData   []byte                `json:"eventData"`
	TargetsData []byte                `json:"targetsData"`
}

func NewRequest(e eventstore.Event, targets []target.Target) (*Request, error) {
	targetsData, err := json.Marshal(targets)
	if err != nil {
		return nil, err
	}
	// The underlying BaseEvent omits the data when using json.Marshal so we have to unmarshal it manually.
	var eventData json.RawMessage
	err = e.Unmarshal(&eventData)
	if err != nil {
		return nil, err
	}

	return &Request{
		Aggregate:   e.Aggregate(),
		Sequence:    e.Sequence(),
		EventType:   e.Type(),
		CreatedAt:   e.CreatedAt(),
		UserID:      e.Creator(),
		EventData:   eventData,
		TargetsData: targetsData,
	}, nil
}

func (e *Request) Kind() string {
	return "execution_request"
}

func ContextInfoFromRequest(e *Request) *ContextInfoEvent {
	return &ContextInfoEvent{
		AggregateID:   e.Aggregate.ID,
		AggregateType: string(e.Aggregate.Type),
		ResourceOwner: e.Aggregate.ResourceOwner,
		InstanceID:    e.Aggregate.InstanceID,
		Version:       string(e.Aggregate.Version),
		Sequence:      e.Sequence,
		EventType:     string(e.EventType),
		CreatedAt:     e.CreatedAt.Format(time.RFC3339Nano),
		UserID:        e.UserID,
		EventPayload:  e.EventData,
	}
}

type ContextInfoEvent struct {
	AggregateID   string          `json:"aggregateID,omitempty"`
	AggregateType string          `json:"aggregateType,omitempty"`
	ResourceOwner string          `json:"resourceOwner,omitempty"`
	InstanceID    string          `json:"instanceID,omitempty"`
	Version       string          `json:"version,omitempty"`
	Sequence      uint64          `json:"sequence,omitempty"`
	EventType     string          `json:"event_type,omitempty"`
	CreatedAt     string          `json:"created_at,omitempty"`
	UserID        string          `json:"userID,omitempty"`
	EventPayload  json.RawMessage `json:"event_payload,omitempty"`
}

func (c *ContextInfoEvent) GetHTTPRequestBody() []byte {
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return data
}

func (c *ContextInfoEvent) SetHTTPResponseBody(resp []byte) error {
	// response is irrelevant and will not be unmarshaled
	return nil
}

func (c *ContextInfoEvent) GetContent() any {
	return c.EventPayload
}
