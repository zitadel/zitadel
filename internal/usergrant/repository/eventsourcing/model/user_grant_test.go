package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/usergrant/model"
	"reflect"
	"testing"
)

func TestAppendAddGrantEvent(t *testing.T) {
	type args struct {
		grant *UserGrant
		data  *UserGrant
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *UserGrant
	}{
		{
			name: "append add grant event",
			args: args{
				grant: &UserGrant{},
				data:  &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}},
				event: &es_models.Event{},
			},
			result: &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.data != nil {
				data, _ := json.Marshal(tt.args.data)
				tt.args.event.Data = data
			}
			tt.args.grant.appendAddGrantEvent(tt.args.event)
			if tt.args.grant.UserID != tt.result.UserID {
				t.Errorf("got wrong result grantID: actual %v expected: %v", tt.args.grant.UserID, tt.result.UserID)
			}
			if !reflect.DeepEqual(tt.args.grant.RoleKeys, tt.result.RoleKeys) {
				t.Errorf("got wrong result grantID: actual %v expected: %v", tt.args.grant.RoleKeys, tt.result.RoleKeys)
			}
		})
	}
}

func TestAppendChangeGrantEvent(t *testing.T) {
	type args struct {
		existing *UserGrant
		grant    *UserGrant
		event    *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *UserGrant
	}{
		{
			name: "append change grant event",
			args: args{
				existing: &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}},
				grant:    &UserGrant{ProjectID: "GrantID", RoleKeys: []string{"KeyChanged"}},
				event:    &es_models.Event{},
			},
			result: &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"KeyChanged"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.grant != nil {
				data, _ := json.Marshal(tt.args.grant)
				tt.args.event.Data = data
			}
			tt.args.existing.appendChangeGrantEvent(tt.args.event)

			if !reflect.DeepEqual(tt.args.grant.RoleKeys, tt.result.RoleKeys) {
				t.Errorf("got wrong result grantID: actual %v expected: %v", tt.args.grant.RoleKeys, tt.args.grant.RoleKeys)
			}
		})
	}
}

func TestAppendGrantStateEvent(t *testing.T) {
	type args struct {
		existing *UserGrant
		grant    *UserGrantID
		event    *es_models.Event
		state    model.UserGrantState
	}
	tests := []struct {
		name   string
		args   args
		result *UserGrant
	}{
		{
			name: "append deactivate grant event",
			args: args{
				existing: &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}},
				grant:    &UserGrantID{GrantID: "GrantID"},
				event:    &es_models.Event{},
				state:    model.USERGRANTSTATE_INACTIVE,
			},
			result: &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}, State: int32(model.USERGRANTSTATE_INACTIVE)},
		},
		{
			name: "append reactivate grant event",
			args: args{
				existing: &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}},
				grant:    &UserGrantID{GrantID: "GrantID"},
				event:    &es_models.Event{},
				state:    model.USERGRANTSTATE_ACTIVE,
			},
			result: &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}, State: int32(model.USERGRANTSTATE_ACTIVE)},
		},
		{
			name: "append remove grant event",
			args: args{
				existing: &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}},
				grant:    &UserGrantID{GrantID: "GrantID"},
				event:    &es_models.Event{},
				state:    model.USERGRANTSTATE_REMOVED,
			},
			result: &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}, State: int32(model.USERGRANTSTATE_REMOVED)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.existing.appendGrantStateEvent(tt.args.state)
			if tt.args.existing.State != tt.result.State {
				t.Errorf("got wrong result: actual: %v, expected: %v ", tt.result.State, tt.args.existing.State)
			}
		})
	}
}
