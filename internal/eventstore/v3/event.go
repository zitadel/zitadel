package eventstore

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	_ eventstore.Event = (*event)(nil)
)

type command struct {
	InstanceID    string
	AggregateType string
	AggregateID   string
	CommandType   string
	Revision      uint16
	Payload       Payload
	Creator       string
	Owner         string
}

func (c *command) Aggregate() *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            c.AggregateID,
		Type:          eventstore.AggregateType(c.AggregateType),
		ResourceOwner: c.Owner,
		InstanceID:    c.InstanceID,
		Version:       eventstore.Version("v" + strconv.Itoa(int(c.Revision))),
	}
}

type event struct {
	command   *command
	createdAt time.Time
	sequence  uint64
	position  float64
}

// TODO: remove on v3
func commandToEventOld(sequence *latestSequence, cmd eventstore.Command) (_ *event, err error) {
	var payload Payload
	if cmd.Payload() != nil {
		payload, err = json.Marshal(cmd.Payload())
		if err != nil {
			logging.WithError(err).Warn("marshal payload failed")
			return nil, zerrors.ThrowInternal(err, "V3-MInPK", "Errors.Internal")
		}
	}
	return &event{
		command: &command{
			InstanceID:    sequence.aggregate.InstanceID,
			AggregateType: string(sequence.aggregate.Type),
			AggregateID:   sequence.aggregate.ID,
			CommandType:   string(cmd.Type()),
			Revision:      cmd.Revision(),
			Payload:       payload,
			Creator:       cmd.Creator(),
			Owner:         sequence.aggregate.ResourceOwner,
		},
		sequence: sequence.sequence,
	}, nil
}

func commandsToEvents(ctx context.Context, cmds []eventstore.Command) (_ []eventstore.Event, _ []*command, err error) {
	events := make([]eventstore.Event, len(cmds))
	commands := make([]*command, len(cmds))
	for i, cmd := range cmds {
		if cmd.Aggregate().InstanceID == "" {
			cmd.Aggregate().InstanceID = authz.GetInstance(ctx).InstanceID()
		}
		events[i], err = commandToEvent(cmd)
		if err != nil {
			return nil, nil, err
		}
		commands[i] = events[i].(*event).command
	}
	return events, commands, nil
}

func commandToEvent(cmd eventstore.Command) (_ eventstore.Event, err error) {
	var payload Payload
	if cmd.Payload() != nil {
		payload, err = json.Marshal(cmd.Payload())
		if err != nil {
			logging.WithError(err).Warn("marshal payload failed")
			return nil, zerrors.ThrowInternal(err, "V3-MInPK", "Errors.Internal")
		}
	}

	command := &command{
		InstanceID:    cmd.Aggregate().InstanceID,
		AggregateType: string(cmd.Aggregate().Type),
		AggregateID:   cmd.Aggregate().ID,
		CommandType:   string(cmd.Type()),
		Revision:      cmd.Revision(),
		Payload:       payload,
		Creator:       cmd.Creator(),
		Owner:         cmd.Aggregate().ResourceOwner,
	}

	return &event{
		command: command,
	}, nil
}

// CreationDate implements [eventstore.Event]
func (e *event) CreationDate() time.Time {
	return e.CreatedAt()
}

// EditorUser implements [eventstore.Event]
func (e *event) EditorUser() string {
	return e.Creator()
}

// Aggregate implements [eventstore.Event]
func (e *event) Aggregate() *eventstore.Aggregate {
	return e.command.Aggregate()
}

// Creator implements [eventstore.Event]
func (e *event) Creator() string {
	return e.command.Creator
}

// Revision implements [eventstore.Event]
func (e *event) Revision() uint16 {
	return e.command.Revision
}

// Type implements [eventstore.Event]
func (e *event) Type() eventstore.EventType {
	return eventstore.EventType(e.command.CommandType)
}

// CreatedAt implements [eventstore.Event]
func (e *event) CreatedAt() time.Time {
	return e.createdAt
}

// Sequence implements [eventstore.Event]
func (e *event) Sequence() uint64 {
	return e.sequence
}

// Sequence implements [eventstore.Event]
func (e *event) Position() float64 {
	return e.position
}

// Unmarshal implements [eventstore.Event]
func (e *event) Unmarshal(ptr any) error {
	if len(e.command.Payload) == 0 {
		return nil
	}
	if err := json.Unmarshal(e.command.Payload, ptr); err != nil {
		return zerrors.ThrowInternal(err, "V3-u8qVo", "Errors.Internal")
	}

	return nil
}

// DataAsBytes implements [eventstore.Event]
func (e *event) DataAsBytes() []byte {
	return e.command.Payload
}
