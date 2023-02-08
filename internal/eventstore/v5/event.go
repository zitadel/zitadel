package eventstore

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/service"
)

func NewCommand(typ string, version uint8, editor *Editor, agg *Aggregate) *Command {
	return &Command{
		typ:       typ,
		version:   version,
		editor:    editor,
		aggregate: agg,
	}
}

// Command represents an intent to store an action
type Command struct {
	typ       string     `json:"-"`
	version   uint8      `json:"-"`
	editor    *Editor    `json:"-"`
	aggregate *Aggregate `json:"-"`
}

// Type represents the action made
func (cmd *Command) Type() string {
	return cmd.typ
}

// Version describes the object definition at the time the event was created.
// The version counts up as soon as the action definition changes
func (cmd *Command) Version() uint8 {
	return cmd.version
}

// Editor represents who created the event
func (cmd *Command) Editor() *Editor {
	return cmd.editor
}

// Aggregate represents a definition of an object
func (cmd *Command) Aggregate() *Aggregate {
	return cmd.aggregate
}

func (cmd *Command) internal() {}

// Event represents an action made in the past
type Event struct {
	// Type represents the action made
	Type string `json:"-"`
	// Version describes the object definition at the time the event was created.
	// The version counts up as soon as the action definition changes
	Version uint8 `json:"-"`
	// Sequence is an up counting number representing the current revision of an object
	Sequence uint64 `json:"-"`
	// Payload represents the data added or changed
	Payload []byte `json:"-"`
	// Editor represents who created the event
	Editor Editor `json:"-"`
	// Aggregate represents the object the event belongs to
	Aggregate Aggregate `json:"-"`
}

func (e *Event) UnmarshalPayload(target any) error {
	if len(e.Payload) == 0 {
		return nil
	}
	return json.Unmarshal(e.Payload, target)
}

type Editor struct {
	// UserID represents the creator of the events
	// The field is required
	UserID string
	// Service represents which API service was used to create the event
	Service string
}

// NewEditorFromCtx creates an editor from the ctx.
// The user is grabbed from [authz.CtxData]
// The service is grabbed from [service]
func NewEditorFromCtx(ctx context.Context) *Editor {
	return &Editor{
		UserID:  authz.GetCtxData(ctx).UserID,
		Service: service.FromContext(ctx),
	}
}

// Aggregate represents a definition of an object
type Aggregate struct {
	// ID represents the identityfier of the aggregate
	ID string
	// Type represents the object type of the aggregate
	Type string
	// Owner represents the organization which owns the aggregate
	Owner string
	// InstanceID represents the instance the object is associated with
	InstanceID string
	sequence   uint64
}

// NewAggregate creates an aggregate
// The id is the object id
// The typ is the object type
// The ctx is used to grab the owner from [authz.CtxData]
// and instance id from [authz.Instance]
func NewAggregate(ctx context.Context, id, typ string) *Aggregate {
	return &Aggregate{
		ID:         id,
		Type:       typ,
		Owner:      authz.GetCtxData(ctx).OrgID,
		InstanceID: authz.GetInstance(ctx).InstanceID(),
	}
}
