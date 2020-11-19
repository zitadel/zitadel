package model

import (
	"encoding/json"
	"reflect"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/usergrant/model"
	es_model "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"
	"github.com/lib/pq"
)

func mockUserGrantData(grant *es_model.UserGrant) []byte {
	data, _ := json.Marshal(grant)
	return data
}

func TestUserAppendEvent(t *testing.T) {
	type args struct {
		event *es_models.Event
		grant *UserGrantView
	}
	tests := []struct {
		name   string
		args   args
		result *UserGrantView
	}{
		{
			name: "append added grant event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserGrantAdded, ResourceOwner: "OrgID", Data: mockUserGrantData(&es_model.UserGrant{UserID: "UserID", ProjectID: "ProjectID", RoleKeys: pq.StringArray{"Keys"}})},
				grant: &UserGrantView{},
			},
			result: &UserGrantView{ID: "AggregateID", ResourceOwner: "OrgID", UserID: "UserID", ProjectID: "ProjectID", RoleKeys: pq.StringArray{"Keys"}, State: int32(model.UserGrantStateActive)},
		},
		{
			name: "append change grant profile event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserGrantChanged, ResourceOwner: "OrgID", Data: mockUserGrantData(&es_model.UserGrant{RoleKeys: pq.StringArray{"KeysChanged"}})},
				grant: &UserGrantView{ID: "AggregateID", ResourceOwner: "OrgID", UserID: "UserID", ProjectID: "ProjectID", RoleKeys: pq.StringArray{"Keys"}, State: int32(model.UserGrantStateActive)},
			},
			result: &UserGrantView{ID: "AggregateID", ResourceOwner: "OrgID", UserID: "UserID", ProjectID: "ProjectID", RoleKeys: pq.StringArray{"KeysChanged"}, State: int32(model.UserGrantStateActive)},
		},
		{
			name: "append grant deactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserGrantDeactivated, ResourceOwner: "OrgID"},
				grant: &UserGrantView{ID: "AggregateID", ResourceOwner: "OrgID", UserID: "UserID", ProjectID: "ProjectID", RoleKeys: pq.StringArray{"Keys"}, State: int32(model.UserGrantStateActive)},
			},
			result: &UserGrantView{ID: "AggregateID", ResourceOwner: "OrgID", UserID: "UserID", ProjectID: "ProjectID", RoleKeys: pq.StringArray{"Keys"}, State: int32(model.UserGrantStateInactive)},
		},
		{
			name: "append grant reactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.UserGrantReactivated, ResourceOwner: "OrgID"},
				grant: &UserGrantView{ID: "AggregateID", ResourceOwner: "OrgID", UserID: "UserID", ProjectID: "ProjectID", RoleKeys: pq.StringArray{"Keys"}, State: int32(model.UserGrantStateInactive)},
			},
			result: &UserGrantView{ID: "AggregateID", ResourceOwner: "OrgID", UserID: "UserID", ProjectID: "ProjectID", RoleKeys: pq.StringArray{"Keys"}, State: int32(model.UserGrantStateActive)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.grant.AppendEvent(tt.args.event)
			if tt.args.grant.ID != tt.result.ID {
				t.Errorf("got wrong result ID: expected: %v, actual: %v ", tt.result.ID, tt.args.grant.ID)
			}
			if tt.args.grant.ResourceOwner != tt.result.ResourceOwner {
				t.Errorf("got wrong result ResourceOwner: expected: %v, actual: %v ", tt.result.ResourceOwner, tt.args.grant.ResourceOwner)
			}
			if tt.args.grant.UserID != tt.result.UserID {
				t.Errorf("got wrong result UserID: expected: %v, actual: %v ", tt.result.UserID, tt.args.grant.UserID)
			}
			if tt.args.grant.ProjectID != tt.result.ProjectID {
				t.Errorf("got wrong result ProjectID: expected: %v, actual: %v ", tt.result.ProjectID, tt.args.grant.ProjectID)
			}
			if !reflect.DeepEqual(tt.args.grant.RoleKeys, tt.result.RoleKeys) {
				t.Errorf("got wrong result RoleKeys: expected: %v, actual: %v ", tt.result.RoleKeys, tt.args.grant.RoleKeys)
			}
		})
	}
}
