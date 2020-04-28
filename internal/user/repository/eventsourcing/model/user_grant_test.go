package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"testing"
)

func TestAppendAddGrantEvent(t *testing.T) {
	type args struct {
		user  *User
		role  *UserGrant
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *User
	}{
		{
			name: "append add grant event",
			args: args{
				user:  &User{},
				role:  &UserGrant{GrantID: "GrantID", RoleKeys: []string{"Key"}},
				event: &es_models.Event{},
			},
			result: &User{Grants: []*UserGrant{&UserGrant{GrantID: "GrantID", RoleKeys: []string{"Key"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.role != nil {
				data, _ := json.Marshal(tt.args.role)
				tt.args.event.Data = data
			}
			tt.args.user.appendAddGrantEvent(tt.args.event)
			if len(tt.args.user.Grants) != 1 {
				t.Errorf("got wrong result should have one grant actual: %v ", len(tt.args.user.Grants))
			}
			if tt.args.user.Grants[0] == tt.result.Grants[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Grants[0], tt.args.user.Grants[0])
			}
		})
	}
}

func TestAppendChangeGrantEvent(t *testing.T) {
	type args struct {
		user  *User
		grant *UserGrant
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *User
	}{
		{
			name: "append change grant event",
			args: args{
				user:  &User{Grants: []*UserGrant{&UserGrant{ProjectID: "ProjectID", RoleKeys: []string{"Key"}}}},
				grant: &UserGrant{ProjectID: "GrantID", RoleKeys: []string{"KeyChanged"}},
				event: &es_models.Event{},
			},
			result: &User{Grants: []*UserGrant{&UserGrant{GrantID: "GrantID", RoleKeys: []string{"KeyChanged"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.grant != nil {
				data, _ := json.Marshal(tt.args.grant)
				tt.args.event.Data = data
			}
			tt.args.user.appendChangeGrantEvent(tt.args.event)
			if len(tt.args.user.Grants) != 1 {
				t.Errorf("got wrong result should have one grant actual: %v ", len(tt.args.user.Grants))
			}
			if tt.args.user.Grants[0] == tt.result.Grants[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Grants[0], tt.args.user.Grants[0])
			}
		})
	}
}

func TestAppendRemoveGrantEvent(t *testing.T) {
	type args struct {
		user  *User
		grant *UserGrant
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *User
	}{
		{
			name: "append remove role event",
			args: args{
				user:  &User{Grants: []*UserGrant{&UserGrant{GrantID: "GrantID", RoleKeys: []string{"Key"}}}},
				grant: &UserGrant{GrantID: "GrantID"},
				event: &es_models.Event{},
			},
			result: &User{Grants: []*UserGrant{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.grant != nil {
				data, _ := json.Marshal(tt.args.grant)
				tt.args.event.Data = data
			}
			tt.args.user.appendRemoveGrantEvent(tt.args.event)
			if len(tt.args.user.Grants) != 0 {
				t.Errorf("got wrong result should have no grant actual: %v ", len(tt.args.user.Grants))
			}
		})
	}
}

func TestAppendGrantStateEvent(t *testing.T) {
	type args struct {
		user  *User
		grant *UserGrantID
		event *es_models.Event
		state model.UserGrantState
	}
	tests := []struct {
		name   string
		args   args
		result *User
	}{
		{
			name: "append deactivate grant event",
			args: args{
				user:  &User{Grants: []*UserGrant{&UserGrant{ProjectID: "ProjectID", RoleKeys: []string{"Key"}}}},
				grant: &UserGrantID{GrantID: "GrantID"},
				event: &es_models.Event{},
				state: model.USERGRANTSTATE_INACTIVE,
			},
			result: &User{Grants: []*UserGrant{&UserGrant{ProjectID: "ProjectID", RoleKeys: []string{"Key"}, State: int32(model.USERGRANTSTATE_INACTIVE)}}},
		},
		{
			name: "append reactivate grant event",
			args: args{
				user:  &User{Grants: []*UserGrant{&UserGrant{ProjectID: "ProjectID", RoleKeys: []string{"Key"}}}},
				grant: &UserGrantID{GrantID: "GrantID"},
				event: &es_models.Event{},
				state: model.USERGRANTSTATE_ACTIVE,
			},
			result: &User{Grants: []*UserGrant{&UserGrant{ProjectID: "ProjectID", RoleKeys: []string{"Key"}, State: int32(model.USERGRANTSTATE_ACTIVE)}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.grant != nil {
				data, _ := json.Marshal(tt.args.grant)
				tt.args.event.Data = data
			}
			tt.args.user.appendGrantStateEvent(tt.args.event, tt.args.state)
			if len(tt.args.user.Grants) != 1 {
				t.Errorf("got wrong result should have one grant actual: %v ", len(tt.args.user.Grants))
			}
			if tt.args.user.Grants[0] == tt.result.Grants[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Grants[0], tt.args.user.Grants[0])
			}
		})
	}
}
