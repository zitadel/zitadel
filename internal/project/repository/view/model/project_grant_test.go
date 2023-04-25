package model

import (
	"encoding/json"
	"reflect"
	"testing"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/project/model"
	es_model "github.com/zitadel/zitadel/internal/project/repository/eventsourcing/model"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func mockProjectData(project *es_model.Project) []byte {
	data, _ := json.Marshal(project)
	return data
}

func mockProjectGrantData(grant *es_model.ProjectGrant) []byte {
	data, _ := json.Marshal(grant)
	return data
}

func TestProjectGrantAppendEvent(t *testing.T) {
	type args struct {
		event   *es_models.Event
		project *ProjectGrantView
	}
	tests := []struct {
		name   string
		args   args
		result *ProjectGrantView
	}{
		{
			name: "append added project grant event",
			args: args{
				event:   &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: project.GrantAddedType, ResourceOwner: "GrantedOrgID", Data: mockProjectGrantData(&es_model.ProjectGrant{GrantID: "ProjectGrantID", GrantedOrgID: "GrantedOrgID", RoleKeys: []string{"Role"}})},
				project: &ProjectGrantView{},
			},
			result: &ProjectGrantView{ProjectID: "AggregateID", ResourceOwner: "GrantedOrgID", OrgID: "GrantedOrgID", State: int32(model.ProjectStateActive), GrantedRoleKeys: []string{"Role"}},
		},
		{
			name: "append change project grant event",
			args: args{
				event:   &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: project.GrantChangedType, ResourceOwner: "GrantedOrgID", Data: mockProjectGrantData(&es_model.ProjectGrant{GrantID: "ProjectGrantID", RoleKeys: []string{"RoleChanged"}})},
				project: &ProjectGrantView{ProjectID: "AggregateID", ResourceOwner: "GrantedOrgID", OrgID: "GrantedOrgID", State: int32(model.ProjectStateActive), GrantedRoleKeys: []string{"Role"}},
			},
			result: &ProjectGrantView{ProjectID: "AggregateID", ResourceOwner: "GrantedOrgID", OrgID: "GrantedOrgID", State: int32(model.ProjectStateActive), GrantedRoleKeys: []string{"RoleChanged"}},
		},
		{
			name: "append deactivate project grant event",
			args: args{
				event:   &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: project.GrantDeactivatedType, ResourceOwner: "GrantedOrgID", Data: mockProjectGrantData(&es_model.ProjectGrant{GrantID: "ProjectGrantID"})},
				project: &ProjectGrantView{ProjectID: "AggregateID", ResourceOwner: "GrantedOrgID", OrgID: "GrantedOrgID", State: int32(model.ProjectStateActive), GrantedRoleKeys: []string{"Role"}},
			},
			result: &ProjectGrantView{ProjectID: "AggregateID", ResourceOwner: "GrantedOrgID", OrgID: "GrantedOrgID", State: int32(model.ProjectStateInactive), GrantedRoleKeys: []string{"Role"}},
		},
		{
			name: "append reactivate project grant event",
			args: args{
				event:   &es_models.Event{AggregateID: "AggregateID", Seq: 1, Typ: project.GrantReactivatedType, ResourceOwner: "GrantedOrgID", Data: mockProjectGrantData(&es_model.ProjectGrant{GrantID: "ProjectGrantID"})},
				project: &ProjectGrantView{ProjectID: "AggregateID", ResourceOwner: "GrantedOrgID", OrgID: "GrantedOrgID", State: int32(model.ProjectStateInactive), GrantedRoleKeys: []string{"Role"}},
			},
			result: &ProjectGrantView{ProjectID: "AggregateID", ResourceOwner: "GrantedOrgID", OrgID: "GrantedOrgID", State: int32(model.ProjectStateActive), GrantedRoleKeys: []string{"Role"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.project.AppendEvent(tt.args.event)
			if tt.args.project.ProjectID != tt.result.ProjectID {
				t.Errorf("got wrong result projectID: expected: %v, actual: %v ", tt.result.ProjectID, tt.args.project.ProjectID)
			}
			if tt.args.project.OrgID != tt.result.OrgID {
				t.Errorf("got wrong result orgID: expected: %v, actual: %v ", tt.result.OrgID, tt.args.project.OrgID)
			}
			if tt.args.project.ResourceOwner != tt.result.ResourceOwner {
				t.Errorf("got wrong result ResourceOwner: expected: %v, actual: %v ", tt.result.ResourceOwner, tt.args.project.ResourceOwner)
			}
			if tt.args.project.Name != tt.result.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.result.Name, tt.args.project.Name)
			}
			if tt.args.project.State != tt.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.result.State, tt.args.project.State)
			}
			if !reflect.DeepEqual(tt.args.project.GrantedRoleKeys, tt.result.GrantedRoleKeys) {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.result.GrantedRoleKeys, tt.args.project.GrantedRoleKeys)
			}
		})
	}
}
