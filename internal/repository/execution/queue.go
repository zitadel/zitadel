package execution

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
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
		CreatedAt:     e.CreatedAt.Format(time.RFC3339),
		UserID:        e.UserID,
		EventPayload:  e.EventData,
	}
}

type ContextInfoEvent struct {
	AggregateID   string `json:"aggregateID,omitempty"`
	AggregateType string `json:"aggregateType,omitempty"`
	ResourceOwner string `json:"resourceOwner,omitempty"`
	InstanceID    string `json:"instanceID,omitempty"`
	Version       string `json:"version,omitempty"`
	Sequence      uint64 `json:"sequence,omitempty"`
	EventType     string `json:"event_type,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UserID        string `json:"userID,omitempty"`
	EventPayload  []byte `json:"event_payload,omitempty"`
}

func (c *ContextInfoEvent) GetHTTPRequestBody() []byte {
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return data
}

func (c *ContextInfoEvent) SetHTTPResponseBody(resp []byte) error {
	// response is irrelevant and will not be unmarshalled
	return nil
}

func (c *ContextInfoEvent) GetContent() interface{} {
	return c.EventPayload
}

func (e *Request) WithLogFields(entry *logging.Entry) *logging.Entry {
	return entry.
		WithField("instanceID", e.Aggregate.InstanceID).
		WithField("resourceOwner", e.Aggregate.ResourceOwner).
		WithField("aggregateType", e.Aggregate.Type).
		WithField("aggregateVersion", e.Aggregate.Version).
		WithField("aggregateID", e.Aggregate.ID).
		WithField("sequence", e.Sequence).
		WithField("eventType", e.EventType)
}
