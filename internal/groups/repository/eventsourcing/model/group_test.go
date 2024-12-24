package model

import (
	"encoding/json"
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/groups/model"
	"github.com/zitadel/zitadel/internal/repository/group"
)

func TestGroupFromEvents(t *testing.T) {
	type args struct {
		event []eventstore.Event
		group *Group
	}
	tests := []struct {
		name   string
		args   args
		result *Group
	}{
		{
			name: "group from events, ok",
			args: args{
				event: []eventstore.Event{
					&es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: group.GroupAddedType},
				},
				group: &Group{Name: "GroupName"},
			},
			result: &Group{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, State: int32(model.GroupStateActive), Name: "GroupName"},
		},
		{
			name: "group from events, nil group",
			args: args{
				event: []eventstore.Event{
					&es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: group.GroupAddedType},
				},
				group: nil,
			},
			result: &Group{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, State: int32(model.GroupStateActive)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.group != nil {
				data, _ := json.Marshal(tt.args.group)
				tt.args.event[0].(*es_models.Event).Data = data
			}
			result, _ := GroupFromEvents(tt.args.group, tt.args.event...)
			if result.Name != tt.result.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.result.Name, result.Name)
			}
		})
	}
}

func TestAppendEvent(t *testing.T) {
	type args struct {
		event *es_models.Event
		group *Group
	}
	tests := []struct {
		name   string
		args   args
		result *Group
	}{
		{
			name: "append added event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: group.GroupAddedType},
				group: &Group{Name: "GroupName"},
			},
			result: &Group{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, State: int32(model.GroupStateActive), Name: "GroupName"},
		},
		{
			name: "append change event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: group.GroupChangedType},
				group: &Group{Name: "GroupName"},
			},
			result: &Group{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, State: int32(model.GroupStateActive), Name: "GroupName"},
		},
		{
			name: "append deactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: group.GroupDeactivatedType},
			},
			result: &Group{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, State: int32(model.GroupStateInactive)},
		},
		{
			name: "append reactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: group.GroupReactivatedType},
			},
			result: &Group{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, State: int32(model.GroupStateActive)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.group != nil {
				data, _ := json.Marshal(tt.args.group)
				tt.args.event.Data = data
			}
			result := new(Group)
			result.AppendEvent(tt.args.event)
			if result.State != tt.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.result.State, result.State)
			}
			if result.Name != tt.result.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.result.Name, result.Name)
			}
			if result.ObjectRoot.AggregateID != tt.result.ObjectRoot.AggregateID {
				t.Errorf("got wrong result id: expected: %v, actual: %v ", tt.result.ObjectRoot.AggregateID, result.ObjectRoot.AggregateID)
			}
		})
	}
}

func TestAppendDeactivatedEvent(t *testing.T) {
	type args struct {
		group *Group
	}
	tests := []struct {
		name   string
		args   args
		result *Group
	}{
		{
			name: "append reactivate event",
			args: args{
				group: &Group{},
			},
			result: &Group{State: int32(model.GroupStateInactive)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.group.appendDeactivatedEvent()
			if tt.args.group.State != tt.result.State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.group)
			}
		})
	}
}

func TestAppendReactivatedEvent(t *testing.T) {
	type args struct {
		group *Group
	}
	tests := []struct {
		name   string
		args   args
		result *Group
	}{
		{
			name: "append reactivate event",
			args: args{
				group: &Group{},
			},
			result: &Group{State: int32(model.GroupStateActive)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.group.appendReactivatedEvent()
			if tt.args.group.State != tt.result.State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.group)
			}
		})
	}
}
