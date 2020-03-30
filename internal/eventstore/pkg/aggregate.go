package pkg

import (
	"context"

	structpb "github.com/golang/protobuf/ptypes/struct"

	es_api "github.com/caos/citadel/eventstore/api/grpc"
	"github.com/caos/utils/proto"
	"github.com/golang/protobuf/ptypes"
)

type AggregateType string

type Aggregate struct {
	id              string
	typ             string
	currentSequence uint64
	modifierService string
	events          []*es_api.EventRequest
}

func NewAggregate(ctx context.Context, id string, typ AggregateType, currentSequence uint64) *Aggregate {
	return &Aggregate{
		id:              id,
		typ:             string(typ),
		currentSequence: currentSequence,
		modifierService: editorService(ctx),
		events:          make([]*es_api.EventRequest, 0, 5),
	}
}

func (a *Aggregate) AppendEvent(ctx context.Context, typ string, payload interface{}) *Aggregate {
	a.events = append(a.events, &es_api.EventRequest{
		CreationDate:    ptypes.TimestampNow(),
		Type:            typ,
		ModifierService: a.modifierService,
		ModifierUser:    editorUser(ctx),
		ModifierTenant:  editorOrg(ctx),
		ResourceOwner:   resourceOwner(ctx),
		Data:            toPBStruct(payload),
	})

	return a
}

func (a *Aggregate) ToAPI() *es_api.AggregateRequest {
	return &es_api.AggregateRequest{
		Type:           a.typ,
		Id:             a.id,
		LatestSequence: a.currentSequence,
		Events:         a.events,
	}
}

func toPBStruct(payload interface{}) *structpb.Struct {
	data, _ := proto.ToPBStruct(payload)
	return data
}

func editorUser(ctx context.Context) string {
	return "userID"
}

func editorOrg(ctx context.Context) string {
	return "orgID"
}

func resourceOwner(ctx context.Context) string {
	return "resourceOwner"
}

func editorService(ctx context.Context) string {
	return "service"
}
