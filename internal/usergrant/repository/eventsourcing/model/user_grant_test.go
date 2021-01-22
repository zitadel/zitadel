package model

import (
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/usergrant/model"
)

func TestAppendGrantStateEvent(t *testing.T) {
	type args struct {
		grant   *UserGrant
		grantID *UserGrantID
		event   *es_models.Event
		state   model.UserGrantState
	}
	tests := []struct {
		name   string
		args   args
		result *UserGrant
	}{
		{
			name: "append deactivate grant event",
			args: args{
				grant:   &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}},
				grantID: &UserGrantID{GrantID: "ProjectGrantID"},
				event:   &es_models.Event{},
				state:   model.UserGrantStateInactive,
			},
			result: &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}, State: int32(model.UserGrantStateInactive)},
		},
		{
			name: "append reactivate grant event",
			args: args{
				grant:   &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}},
				grantID: &UserGrantID{GrantID: "ProjectGrantID"},
				event:   &es_models.Event{},
				state:   model.UserGrantStateActive,
			},
			result: &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}, State: int32(model.UserGrantStateActive)},
		},
		{
			name: "append remove grant event",
			args: args{
				grant:   &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}},
				grantID: &UserGrantID{GrantID: "ProjectGrantID"},
				event:   &es_models.Event{},
				state:   model.UserGrantStateRemoved,
			},
			result: &UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}, State: int32(model.UserGrantStateRemoved)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.grant.appendGrantStateEvent(tt.args.state)
			if tt.args.grant.State != tt.result.State {
				t.Errorf("got wrong result: actual: %v, expected: %v ", tt.result.State, tt.args.grant.State)
			}
		})
	}
}
