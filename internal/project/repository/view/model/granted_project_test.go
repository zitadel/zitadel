package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"github.com/lib/pq"
	"reflect"
	"testing"
)

func mockProjectData(project *es_model.Project) []byte {
	data, _ := json.Marshal(project)
	return data
}

func mockProjectGrantData(grant *es_model.ProjectGrant) []byte {
	data, _ := json.Marshal(grant)
	return data
}

func TestGrantedProjectAppendEvent(t *testing.T) {
	type args struct {
		event   *es_models.Event
		project *GrantedProjectView
	}
	tests := []struct {
		name   string
		args   args
		result *GrantedProjectView
	}{
		{
			name: "append added project event",
			args: args{
				event:   &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ProjectAdded, ResourceOwner: "OrgID", Data: mockProjectData(&es_model.Project{Name: "ProjectName"})},
				project: &GrantedProjectView{},
			},
			result: &GrantedProjectView{ProjectID: "AggregateID", ResourceOwner: "OrgID", OrgID: "OrgID", Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
		},
		{
			name: "append change project event",
			args: args{
				event:   &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ProjectChanged, ResourceOwner: "OrgID", Data: mockProjectData(&es_model.Project{Name: "ProjectNameChanged"})},
				project: &GrantedProjectView{ProjectID: "AggregateID", OrgID: "OrgID", ResourceOwner: "OrgID", Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
			},
			result: &GrantedProjectView{ProjectID: "AggregateID", ResourceOwner: "OrgID", OrgID: "OrgID", Name: "ProjectNameChanged", State: int32(model.PROJECTSTATE_ACTIVE)},
		},
		{
			name: "append project deactivate event",
			args: args{
				event:   &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ProjectDeactivated, ResourceOwner: "OrgID"},
				project: &GrantedProjectView{ProjectID: "AggregateID", ResourceOwner: "OrgID", OrgID: "OrgID", Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
			},
			result: &GrantedProjectView{ProjectID: "AggregateID", ResourceOwner: "OrgID", OrgID: "OrgID", Name: "ProjectName", State: int32(model.PROJECTSTATE_INACTIVE)},
		},
		{
			name: "append project reactivate event",
			args: args{
				event:   &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ProjectReactivated, ResourceOwner: "OrgID"},
				project: &GrantedProjectView{ProjectID: "AggregateID", ResourceOwner: "OrgID", OrgID: "OrgID", Name: "ProjectName", State: int32(model.PROJECTSTATE_INACTIVE)},
			},
			result: &GrantedProjectView{ProjectID: "AggregateID", ResourceOwner: "OrgID", OrgID: "OrgID", Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
		},
		{
			name: "append added project grant event",
			args: args{
				event:   &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ProjectGrantAdded, ResourceOwner: "OrgID", Data: mockProjectGrantData(&es_model.ProjectGrant{GrantID: "GrantID", GrantedOrgID: "GrantedOrgID", RoleKeys: pq.StringArray{"Role"}})},
				project: &GrantedProjectView{},
			},
			result: &GrantedProjectView{ProjectID: "AggregateID", ResourceOwner: "OrgID", OrgID: "GrantedOrgID", State: int32(model.PROJECTSTATE_ACTIVE), GrantedRoleKeys: pq.StringArray{"Role"}},
		},
		{
			name: "append change project grant event",
			args: args{
				event:   &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ProjectGrantChanged, ResourceOwner: "OrgID", Data: mockProjectGrantData(&es_model.ProjectGrant{GrantID: "GrantID", RoleKeys: pq.StringArray{"RoleChanged"}})},
				project: &GrantedProjectView{ProjectID: "AggregateID", ResourceOwner: "OrgID", OrgID: "GrantedOrgID", State: int32(model.PROJECTSTATE_ACTIVE), GrantedRoleKeys: pq.StringArray{"Role"}},
			},
			result: &GrantedProjectView{ProjectID: "AggregateID", ResourceOwner: "OrgID", OrgID: "GrantedOrgID", State: int32(model.PROJECTSTATE_ACTIVE), GrantedRoleKeys: pq.StringArray{"RoleChanged"}},
		},
		{
			name: "append deactivate project grant event",
			args: args{
				event:   &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ProjectGrantDeactivated, ResourceOwner: "OrgID", Data: mockProjectGrantData(&es_model.ProjectGrant{GrantID: "GrantID"})},
				project: &GrantedProjectView{ProjectID: "AggregateID", ResourceOwner: "OrgID", OrgID: "GrantedOrgID", State: int32(model.PROJECTSTATE_ACTIVE), GrantedRoleKeys: pq.StringArray{"Role"}},
			},
			result: &GrantedProjectView{ProjectID: "AggregateID", ResourceOwner: "OrgID", OrgID: "GrantedOrgID", State: int32(model.PROJECTSTATE_INACTIVE), GrantedRoleKeys: pq.StringArray{"Role"}},
		},
		{
			name: "append reactivate project grant event",
			args: args{
				event:   &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ProjectGrantReactivated, ResourceOwner: "OrgID", Data: mockProjectGrantData(&es_model.ProjectGrant{GrantID: "GrantID"})},
				project: &GrantedProjectView{ProjectID: "AggregateID", ResourceOwner: "OrgID", OrgID: "GrantedOrgID", State: int32(model.PROJECTSTATE_INACTIVE), GrantedRoleKeys: pq.StringArray{"Role"}},
			},
			result: &GrantedProjectView{ProjectID: "AggregateID", ResourceOwner: "OrgID", OrgID: "GrantedOrgID", State: int32(model.PROJECTSTATE_ACTIVE), GrantedRoleKeys: pq.StringArray{"Role"}},
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
