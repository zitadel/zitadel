package eventsourcing

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/caos/zitadel/internal/api/authz"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
)

func TestProjectByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es      *ProjectEventstore
		project *model.Project
	}
	type res struct {
		project *model.Project
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "project from events, ok",
			args: args{
				es:      GetMockProjectByIDOK(ctrl),
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
		},
		{
			name: "project from events, no events",
			args: args{
				es:      GetMockProjectByIDNoEvents(ctrl),
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "project from events, no id",
			args: args{
				es:      GetMockProjectByIDNoEvents(ctrl),
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ProjectByID(nil, tt.args.project.AggregateID)

			if !tt.res.wantErr && result.AggregateID != tt.res.project.AggregateID {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.project.AggregateID, result.AggregateID)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreateProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es      *ProjectEventstore
		ctx     context.Context
		project *model.Project
		global  bool
	}
	type res struct {
		project *model.Project
		role    string
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create project, ok",
			args: args{
				es:      GetMockManipulateProject(ctrl),
				ctx:     authz.NewMockContext("orgID", "userID"),
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "Name"},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "Name"},
				role:    projectOwnerRole,
			},
		},
		{
			name: "create global project, ok",
			args: args{
				es:      GetMockManipulateProject(ctrl),
				ctx:     authz.NewMockContext("orgID", "userID"),
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "Name"},
				global:  true,
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "Name"},
				role:    projectOwnerGlobalRole,
			},
		},
		{
			name: "create project no name",
			args: args{
				es:      GetMockManipulateProject(ctrl),
				ctx:     authz.NewMockContext("orgID", "userID"),
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.CreateProject(tt.args.ctx, tt.args.project, tt.args.global)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.Name != tt.res.project.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.project.Name, result.Name)
			}
			if !tt.res.wantErr && result.Members[0].Roles[0] != tt.res.role {
				t.Errorf("got wrong result role: expected: %v, actual: %v ", tt.res.role, result.Members[0].Roles[0])
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUpdateProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *ProjectEventstore
		ctx context.Context
		new *model.Project
	}
	type res struct {
		project *model.Project
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "update project, ok",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				new: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "NameNew"},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "NameNew"},
			},
		},
		{
			name: "update project no name",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				new: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: ""},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				new: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "NameNew"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UpdateProject(tt.args.ctx, tt.args.new)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.Name != tt.res.project.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.project.Name, result.Name)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestDeactivateProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es              *ProjectEventstore
		ctx             context.Context
		existingProject *model.Project
		newProject      *model.Project
	}
	type res struct {
		project *model.Project
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "deactivate project, ok",
			args: args{
				es:              GetMockManipulateProject(ctrl),
				ctx:             authz.NewMockContext("orgID", "userID"),
				existingProject: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "Name", State: model.ProjectStateActive},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "NameNew", State: model.ProjectStateInactive},
			},
		},
		{
			name: "deactivate project with inactive state",
			args: args{
				es:              GetMockManipulateInactiveProject(ctrl),
				ctx:             authz.NewMockContext("orgID", "userID"),
				existingProject: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "Name", State: model.ProjectStateInactive},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:              GetMockManipulateProjectNoEvents(ctrl),
				ctx:             authz.NewMockContext("orgID", "userID"),
				existingProject: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "Name", State: model.ProjectStateActive},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.DeactivateProject(tt.args.ctx, tt.args.existingProject.AggregateID)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.State != tt.res.project.State {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.project.State, result.State)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestReactivateProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es             *ProjectEventstore
		ctx            context.Context
		existingPolicy *model.Project
		newPolicy      *model.Project
	}
	type res struct {
		project *model.Project
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "reactivate project, ok",
			args: args{
				es:             GetMockManipulateInactiveProject(ctrl),
				ctx:            authz.NewMockContext("orgID", "userID"),
				existingPolicy: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "Name", State: model.ProjectStateInactive},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "NameNew", State: model.ProjectStateActive},
			},
		},
		{
			name: "reactivate project with inactive state",
			args: args{
				es:             GetMockManipulateProject(ctrl),
				ctx:            authz.NewMockContext("orgID", "userID"),
				existingPolicy: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "Name", State: model.ProjectStateActive},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:             GetMockManipulateProjectNoEvents(ctrl),
				ctx:            authz.NewMockContext("orgID", "userID"),
				existingPolicy: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Name: "Name", State: model.ProjectStateActive},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ReactivateProject(tt.args.ctx, tt.args.existingPolicy.AggregateID)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.State != tt.res.project.State {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.project.State, result.State)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRemoveProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es             *ProjectEventstore
		ctx            context.Context
		existingPolicy *model.Project
	}
	type res struct {
		result  *model.Project
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove project, ok",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				existingPolicy: &model.Project{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					Name:       "Name",
					Members: []*model.ProjectMember{
						{UserID: "UserID", Roles: []string{"Roles"}},
					},
				},
			},
			res: res{
				result: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
		},
		{
			name: "no projectid",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				existingPolicy: &model.Project{
					ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Name:       "Name",
					Members: []*model.ProjectMember{
						{UserID: "UserID", Roles: []string{"Roles"}},
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "project not existing",
			args: args{
				es:             GetMockManipulateProject(ctrl),
				ctx:            authz.NewMockContext("orgID", "userID"),
				existingPolicy: &model.Project{},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:             GetMockManipulateProjectNoEvents(ctrl),
				ctx:            authz.NewMockContext("orgID", "userID"),
				existingPolicy: &model.Project{ObjectRoot: es_models.ObjectRoot{AggregateID: "OtherAggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveProject(tt.args.ctx, tt.args.existingPolicy)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestProjectMemberByIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *ProjectEventstore
		member *model.ProjectMember
	}
	type res struct {
		member  *model.ProjectMember
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "project member from events, ok",
			args: args{
				es:     GetMockProjectMemberByIDsOK(ctrl),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID"},
			},
			res: res{
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Role"}},
			},
		},
		{
			name: "project member from events, no events",
			args: args{
				es:     GetMockProjectByIDNoEvents(ctrl),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "project member from events, no id",
			args: args{
				es:     GetMockProjectByIDNoEvents(ctrl),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: func(err error) bool {
					return caos_errs.IsPreconditionFailed(err)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ProjectMemberByIDs(nil, tt.args.member)
			if !tt.res.wantErr && result.AggregateID != tt.res.member.AggregateID {
				t.Errorf("got wrong result id: expected: %v, actual: %v ", tt.res.member.AggregateID, result.AggregateID)
			}
			if !tt.res.wantErr && result.UserID != tt.res.member.UserID {
				t.Errorf("got wrong result userid: expected: %v, actual: %v ", tt.res.member.UserID, result.UserID)
			}
			if !tt.res.wantErr && len(result.Roles) != len(tt.res.member.Roles) {
				t.Errorf("got wrong result roles: expected: %v, actual: %v ", tt.res.member.Roles, result.Roles)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddProjectMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *ProjectEventstore
		ctx    context.Context
		member *model.ProjectMember
	}
	type res struct {
		result  *model.ProjectMember
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add project member, ok",
			args: args{
				es:     GetMockManipulateProject(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
			res: res{
				result: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
		},
		{
			name: "no userid",
			args: args{
				es:     GetMockManipulateProject(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Roles: []string{"Roles"}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "no roles",
			args: args{
				es:     GetMockManipulateProject(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member already existing",
			args: args{
				es:     GetMockManipulateProjectWithMember(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:     GetMockManipulateProjectNoEvents(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddProjectMember(tt.args.ctx, tt.args.member)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.UserID != tt.res.result.UserID {
				t.Errorf("got wrong result userid: expected: %v, actual: %v ", tt.res.result.UserID, result.UserID)
			}
			if !tt.res.wantErr && len(result.Roles) != len(tt.res.result.Roles) {
				t.Errorf("got wrong result roles: expected: %v, actual: %v ", tt.res.result.Roles, result.Roles)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeProjectMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *ProjectEventstore
		ctx    context.Context
		member *model.ProjectMember
	}
	type res struct {
		result  *model.ProjectMember
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add project member, ok",
			args: args{
				es:     GetMockManipulateProjectWithMember(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"ChangeRoles"}},
			},
			res: res{
				result: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
		},
		{
			name: "no userid",
			args: args{
				es:     GetMockManipulateProject(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Roles: []string{"ChangeRoles"}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "no roles",
			args: args{
				es:     GetMockManipulateProject(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member not existing",
			args: args{
				es:     GetMockManipulateProject(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:     GetMockManipulateProjectNoEvents(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"ChangeRoles"}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeProjectMember(tt.args.ctx, tt.args.member)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.UserID != tt.res.result.UserID {
				t.Errorf("got wrong result userid: expected: %v, actual: %v ", tt.res.result.UserID, result.UserID)
			}
			if !tt.res.wantErr && len(result.Roles) != len(tt.res.result.Roles) {
				t.Errorf("got wrong result roles: expected: %v, actual: %v ", tt.res.result.Roles, result.Roles)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRemoveProjectMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es              *ProjectEventstore
		ctx             context.Context
		existingProject *model.Project
		member          *model.ProjectMember
	}
	type res struct {
		result  *model.ProjectMember
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove project member, ok",
			args: args{
				es:  GetMockManipulateProjectWithMember(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				existingProject: &model.Project{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					Name:       "Name",
					Members: []*model.ProjectMember{
						{UserID: "UserID", Roles: []string{"Roles"}},
					},
				},
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID"},
			},
			res: res{
				result: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
		},
		{
			name: "no userid",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				existingProject: &model.Project{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					Name:       "Name",
					Members: []*model.ProjectMember{
						{UserID: "UserID", Roles: []string{"Roles"}},
					},
				},
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Roles: []string{"ChangeRoles"}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member not existing",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				existingProject: &model.Project{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					Name:       "Name",
				},
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:     GetMockManipulateProjectNoEvents(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"ChangeRoles"}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveProjectMember(tt.args.ctx, tt.args.member)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddProjectRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *ProjectEventstore
		ctx   context.Context
		roles []*model.ProjectRole
	}
	type res struct {
		result  *model.ProjectRole
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add project role, ok",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				roles: []*model.ProjectRole{
					{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
				},
			},
			res: res{
				result: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
			},
		},
		{
			name: "no key",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				roles: []*model.ProjectRole{
					{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, DisplayName: "DisplayName", Group: "Group"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "role already existing",
			args: args{
				es:  GetMockManipulateProjectWithRole(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				roles: []*model.ProjectRole{
					{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				roles: []*model.ProjectRole{
					{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddProjectRoles(tt.args.ctx, tt.args.roles...)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.Key != tt.res.result.Key {
				t.Errorf("got wrong result key: expected: %v, actual: %v ", tt.res.result.Key, result.Key)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeProjectRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *ProjectEventstore
		ctx  context.Context
		role *model.ProjectRole
	}
	type res struct {
		result  *model.ProjectRole
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change project role, ok",
			args: args{
				es:   GetMockManipulateProjectWithRole(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Key: "Key", DisplayName: "DisplayNameChanged", Group: "Group"},
			},
			res: res{
				result: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Key: "Key", DisplayName: "DisplayNameChanged", Group: "Group"},
			},
		},
		{
			name: "no key",
			args: args{
				es:   GetMockManipulateProject(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, DisplayName: "DisplayName", Group: "Group"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "role not existing",
			args: args{
				es:   GetMockManipulateProject(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:   GetMockManipulateProjectNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeProjectRole(tt.args.ctx, tt.args.role)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.Key != tt.res.result.Key {
				t.Errorf("got wrong result key: expected: %v, actual: %v ", tt.res.result.Key, result.Key)
			}
			if !tt.res.wantErr && result.DisplayName != tt.res.result.DisplayName {
				t.Errorf("got wrong result displayName: expected: %v, actual: %v ", tt.res.result.DisplayName, result.DisplayName)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRemoveProjectRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *ProjectEventstore
		ctx  context.Context
		role *model.ProjectRole
	}
	type res struct {
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove project role, ok",
			args: args{
				es:   GetMockManipulateProjectWithRole(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Key: "Key"},
			},
		},
		{
			name: "no key",
			args: args{
				es:   GetMockManipulateProject(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, DisplayName: "DisplayName", Group: "Group"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "role not existing",
			args: args{
				es:   GetMockManipulateProject(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:   GetMockManipulateProjectNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveProjectRole(tt.args.ctx, tt.args.role)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestApplicationByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *ProjectEventstore
		app *model.Application
	}
	type res struct {
		app     *model.Application
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get app",
			args: args{
				es:  GetMockProjectAppsByIDsOK(ctrl),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, AppID: "AppID"},
			},
			res: res{
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID:      "AppID",
					Name:       "Name",
					OIDCConfig: &model.OIDCConfig{ClientID: "ClientID"}},
			},
		},
		{
			name: "no events for project",
			args: args{
				es:  GetMockProjectByIDNoEvents(ctrl),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, AppID: "AppID"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "app has no id",
			args: args{
				es:  GetMockProjectByIDNoEvents(ctrl),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: func(err error) bool {
					return caos_errs.IsPreconditionFailed(err)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ApplicationByIDs(nil, tt.args.app.AggregateID, tt.args.app.AppID)
			if !tt.res.wantErr && result.AggregateID != tt.res.app.AggregateID {
				t.Errorf("got wrong result id: expected: %v, actual: %v ", tt.res.app.AggregateID, result.AggregateID)
			}
			if !tt.res.wantErr && result.AppID != tt.res.app.AppID {
				t.Errorf("got wrong result appid: expected: %v, actual: %v ", tt.res.app.AppID, result.AppID)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *ProjectEventstore
		ctx context.Context
		app *model.Application
	}
	type res struct {
		result  *model.Application
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add app, ok",
			args: args{
				es:  GetMockManipulateProjectWithPw(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeCode},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
					},
				},
			},
			res: res{
				result: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					Name: "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeCode},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
					},
				},
			},
		},
		{
			name: "add app (none), ok",
			args: args{
				es:  GetMockManipulateProjectWithPw(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes:  []model.OIDCResponseType{model.OIDCResponseTypeCode},
						GrantTypes:     []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
						AuthMethodType: model.OIDCAuthMethodTypeNone,
					},
				},
			},
			res: res{
				result: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					Name: "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes:  []model.OIDCResponseType{model.OIDCResponseTypeCode},
						GrantTypes:     []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
						AuthMethodType: model.OIDCAuthMethodTypeNone,
					},
				},
			},
		},
		{
			name: "invalid app",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeCode},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
					},
				},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddApplication(tt.args.ctx, tt.args.app)
			if tt.res.errFunc == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if tt.res.result != nil && result.AppID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.result != nil && (tt.res.result.OIDCConfig.AuthMethodType != model.OIDCAuthMethodTypeNone && result.OIDCConfig.ClientSecretString == "") {
				t.Errorf("result has no client secret")
			}
			if tt.res.result != nil && result.OIDCConfig.ClientID == "" {
				t.Errorf("result has no clientid")
			}
			if tt.res.result != nil && tt.res.result.Name != result.Name {
				t.Errorf("got wrong result key: expected: %v, actual: %v ", tt.res.result.Name, result.Name)
			}
		})
	}
}

func TestChangeApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *ProjectEventstore
		ctx context.Context
		app *model.Application
	}
	type res struct {
		result  *model.Application
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change app, ok",
			args: args{
				es:  GetMockManipulateProjectWithOIDCApp(ctrl, model.OIDCAuthMethodTypeBasic),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "NameChanged",
				},
			},
			res: res{
				result: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "NameChanged",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeCode},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
					},
				},
			},
		},
		{
			name: "invalid app",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app not existing",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeCode},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeCode},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeApplication(tt.args.ctx, tt.args.app)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.AppID != tt.res.result.AppID {
				t.Errorf("got wrong result AppID: expected: %v, actual: %v ", tt.res.result.AppID, result.AppID)
			}
			if !tt.res.wantErr && result.Name != tt.res.result.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.result.Name, result.Name)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRemoveApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *ProjectEventstore
		ctx context.Context
		app *model.Application
	}
	type res struct {
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove app, ok",
			args: args{
				es:  GetMockManipulateProjectWithOIDCApp(ctrl, model.OIDCAuthMethodTypeBasic),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
				},
			},
		},
		{
			name: "no appID",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app not existing",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveApplication(tt.args.ctx, tt.args.app)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestDeactivateApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *ProjectEventstore
		ctx context.Context
		app *model.Application
	}
	type res struct {
		result  *model.Application
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "deactivate, ok",
			args: args{
				es:  GetMockManipulateProjectWithOIDCApp(ctrl, model.OIDCAuthMethodTypeBasic),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
				},
			},
			res: res{
				result: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					State: model.AppStateInactive,
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeCode},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
					},
				},
			},
		},
		{
			name: "no app id",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app not existing",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeCode},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeCode},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.DeactivateApplication(tt.args.ctx, tt.args.app.AggregateID, tt.args.app.AppID)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.AppID != tt.res.result.AppID {
				t.Errorf("got wrong result AppID: expected: %v, actual: %v ", tt.res.result.AppID, result.AppID)
			}
			if !tt.res.wantErr && result.State != tt.res.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.res.result.State, result.State)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestReactivateApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *ProjectEventstore
		ctx context.Context
		app *model.Application
	}
	type res struct {
		result  *model.Application
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "reactivate, ok",
			args: args{
				es:  GetMockManipulateProjectWithOIDCApp(ctrl, model.OIDCAuthMethodTypeBasic),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
				},
			},
			res: res{
				result: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					State: model.AppStateActive,
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeCode},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
					},
				},
			},
		},
		{
			name: "no app id",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app not existing",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeCode},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeCode},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ReactivateApplication(tt.args.ctx, tt.args.app.AggregateID, tt.args.app.AppID)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.AppID != tt.res.result.AppID {
				t.Errorf("got wrong result AppID: expected: %v, actual: %v ", tt.res.result.AppID, result.AppID)
			}
			if !tt.res.wantErr && result.State != tt.res.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.res.result.State, result.State)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeOIDCConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *ProjectEventstore
		ctx    context.Context
		config *model.OIDCConfig
	}
	type res struct {
		result  *model.OIDCConfig
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change oidc config, ok",
			args: args{
				es:  GetMockManipulateProjectWithOIDCApp(ctrl, model.OIDCAuthMethodTypeBasic),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot:    es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AppID:         "AppID",
					ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeIDToken},
					GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeImplicit},
				},
			},
			res: res{
				result: &model.OIDCConfig{
					ObjectRoot:    es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AppID:         "AppID",
					ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeIDToken},
					GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeImplicit},
				},
			},
		},
		{
			name: "invalid config",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot:    es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AppID:         "AppID",
					ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeIDToken},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app is not oidc",
			args: args{
				es:  GetMockManipulateProjectWithSAMLApp(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot:    es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AppID:         "AppID",
					ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeIDToken},
					GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeImplicit},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app not existing",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot:    es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AppID:         "AppID",
					ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeIDToken},
					GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeImplicit},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot:    es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AppID:         "AppID",
					ResponseTypes: []model.OIDCResponseType{model.OIDCResponseTypeIDToken},
					GrantTypes:    []model.OIDCGrantType{model.OIDCGrantTypeImplicit},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeOIDCConfig(tt.args.ctx, tt.args.config)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.AppID != tt.res.result.AppID {
				t.Errorf("got wrong result AppID: expected: %v, actual: %v ", tt.res.result.AppID, result.AppID)
			}
			if !tt.res.wantErr && result.ResponseTypes[0] != tt.res.result.ResponseTypes[0] {
				t.Errorf("got wrong result responsetype: expected: %v, actual: %v ", tt.res.result.ResponseTypes[0], result.ResponseTypes[0])
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeOIDCConfigSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *ProjectEventstore
		ctx    context.Context
		config *model.OIDCConfig
	}
	type res struct {
		result  *model.OIDCConfig
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change oidc config secret, ok",
			args: args{
				es:  GetMockManipulateProjectWithOIDCApp(ctrl, model.OIDCAuthMethodTypeBasic),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AppID:      "AppID",
				},
			},
			res: res{
				result: &model.OIDCConfig{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AppID:          "AppID",
					ResponseTypes:  []model.OIDCResponseType{model.OIDCResponseTypeCode},
					GrantTypes:     []model.OIDCGrantType{model.OIDCGrantTypeAuthorizationCode},
					AuthMethodType: model.OIDCAuthMethodTypeBasic,
				},
			},
		},
		{
			name: "auth method none, error",
			args: args{
				es:  GetMockManipulateProjectWithOIDCApp(ctrl, model.OIDCAuthMethodTypeNone),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AppID:      "AppID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "no appID",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AppID:      "AppID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app is not oidc",
			args: args{
				es:  GetMockManipulateProjectWithSAMLApp(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AppID:      "AppID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app not existing",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AppID:      "AppID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AppID:      "AppID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeOIDCConfigSecret(tt.args.ctx, tt.args.config.AggregateID, tt.args.config.AppID)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.AppID != tt.res.result.AppID {
				t.Errorf("got wrong result AppID: expected: %v, actual: %v ", tt.res.result.AppID, result.AppID)
			}
			if !tt.res.wantErr && result.ClientSecretString == "" {
				t.Errorf("got wrong result must have client secret")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestProjectGrantByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *ProjectEventstore
		grant *model.ProjectGrant
	}
	type res struct {
		grant   *model.ProjectGrant
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get grant",
			args: args{
				es:    GetMockProjectGrantByIDsOK(ctrl),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1}, GrantID: "ProjectGrantID"},
			},
			res: res{
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID:      "ProjectGrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"Key"},
				},
			},
		},
		{
			name: "no events for project",
			args: args{
				es:    GetMockProjectByIDNoEvents(ctrl),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1}, GrantID: "ProjectGrantID"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "grant has no id",
			args: args{
				es:    GetMockProjectByIDNoEvents(ctrl),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: func(err error) bool {
					return caos_errs.IsPreconditionFailed(err)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ProjectGrantByIDs(nil, tt.args.grant.AggregateID, tt.args.grant.GrantID)
			if !tt.res.wantErr && result.AggregateID != tt.res.grant.AggregateID {
				t.Errorf("got wrong result id: expected: %v, actual: %v ", tt.res.grant.AggregateID, result.AggregateID)
			}
			if !tt.res.wantErr && result.GrantID != tt.res.grant.GrantID {
				t.Errorf("got wrong result grantid: expected: %v, actual: %v ", tt.res.grant.GrantID, result.GrantID)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddProjectGrant(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *ProjectEventstore
		ctx   context.Context
		grant *model.ProjectGrant
	}
	type res struct {
		result  *model.ProjectGrant
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add grant, ok",
			args: args{
				es:  GetMockManipulateProjectWithRole(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID:      "ProjectGrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"Key"},
				},
			},
			res: res{
				result: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID:      "ProjectGrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"Key"},
				},
			},
		},
		{
			name: "invalid grant",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "grant for org already exists",
			args: args{
				es:  GetMockManipulateProjectWithGrant(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID:      "ProjectGrantID",
					GrantedOrgID: "GrantedOrgID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "role not existing on project",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID:      "ProjectGrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"Key"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID:      "ProjectGrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"Key"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddProjectGrant(tt.args.ctx, tt.args.grant)

			if !tt.res.wantErr && result.GrantID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRemoveProjectGrant(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *ProjectEventstore
		ctx   context.Context
		grant *model.ProjectGrant
	}
	type res struct {
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove app, ok",
			args: args{
				es:  GetMockManipulateProjectWithGrant(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
				},
			},
		},
		{
			name: "no grantID",
			args: args{
				es:    GetMockManipulateProject(ctrl),
				ctx:   authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "grant not existing",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveProjectGrant(tt.args.ctx, tt.args.grant)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestDeactivateProjectGrant(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *ProjectEventstore
		ctx   context.Context
		grant *model.ProjectGrant
	}
	type res struct {
		result  *model.ProjectGrant
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "deactivate, ok",
			args: args{
				es:  GetMockManipulateProjectWithGrant(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
				},
			},
			res: res{
				result: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
					State:   model.ProjectGrantStateInactive,
				},
			},
		},
		{
			name: "no grant id",
			args: args{
				es:    GetMockManipulateProject(ctrl),
				ctx:   authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "grant not existing",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.DeactivateProjectGrant(tt.args.ctx, tt.args.grant.AggregateID, tt.args.grant.GrantID)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.GrantID != tt.res.result.GrantID {
				t.Errorf("got wrong result AppID: expected: %v, actual: %v ", tt.res.result.GrantID, result.GrantID)
			}
			if !tt.res.wantErr && result.State != tt.res.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.res.result.State, result.State)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestReactivateProjectGrant(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *ProjectEventstore
		ctx   context.Context
		grant *model.ProjectGrant
	}
	type res struct {
		result  *model.ProjectGrant
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "reactivate, ok",
			args: args{
				es:  GetMockManipulateProjectWithGrant(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
				},
			},
			res: res{
				result: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
					State:   model.ProjectGrantStateActive,
				},
			},
		},
		{
			name: "no grant id",
			args: args{
				es:    GetMockManipulateProject(ctrl),
				ctx:   authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "grant not existing",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ReactivateProjectGrant(tt.args.ctx, tt.args.grant.AggregateID, tt.args.grant.GrantID)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.GrantID != tt.res.result.GrantID {
				t.Errorf("got wrong result AppID: expected: %v, actual: %v ", tt.res.result.GrantID, result.GrantID)
			}
			if !tt.res.wantErr && result.State != tt.res.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.res.result.State, result.State)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestProjectGrantMemberByIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *ProjectEventstore
		member *model.ProjectGrantMember
	}
	type res struct {
		member  *model.ProjectGrantMember
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "projectgrant  member from events, ok",
			args: args{
				es:     GetMockProjectGrantMemberByIDsOK(ctrl),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1}, GrantID: "ProjectGrantID", UserID: "UserID"},
			},
			res: res{
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1}, GrantID: "ProjectGrantID", UserID: "UserID", Roles: []string{"Role"}},
			},
		},
		{
			name: "no project events",
			args: args{
				es:     GetMockProjectByIDNoEvents(ctrl),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1}, GrantID: "ProjectGrantID", UserID: "UserID"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "member id missing",
			args: args{
				es:     GetMockProjectByIDNoEvents(ctrl),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: func(err error) bool {
					return caos_errs.IsPreconditionFailed(err)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ProjectGrantMemberByIDs(nil, tt.args.member)
			if !tt.res.wantErr && result.AggregateID != tt.res.member.AggregateID {
				t.Errorf("got wrong result id: expected: %v, actual: %v ", tt.res.member.AggregateID, result.AggregateID)
			}
			if !tt.res.wantErr && result.UserID != tt.res.member.UserID {
				t.Errorf("got wrong result userid: expected: %v, actual: %v ", tt.res.member.UserID, result.UserID)
			}
			if !tt.res.wantErr && len(result.Roles) != len(tt.res.member.Roles) {
				t.Errorf("got wrong result roles: expected: %v, actual: %v ", tt.res.member.Roles, result.Roles)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddProjectGrantMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *ProjectEventstore
		ctx    context.Context
		member *model.ProjectGrantMember
	}
	type res struct {
		result  *model.ProjectGrantMember
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add project grant member",
			args: args{
				es:  GetMockManipulateProjectWithGrantExistingRole(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
					UserID:  "UserID",
					Roles:   []string{"Role"},
				},
			},
			res: res{
				result: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
					UserID:  "UserID",
					Roles:   []string{"Role"},
				},
			},
		},
		{
			name: "invalid member",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					Roles: []string{"Role"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
					UserID:  "UserID",
					Roles:   []string{"Role"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "member already existing",
			args: args{
				es:  GetMockManipulateProjectWithGrantMember(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
					UserID:  "UserID",
					Roles:   []string{"Role"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorAlreadyExists,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddProjectGrantMember(tt.args.ctx, tt.args.member)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.GrantID != tt.res.result.GrantID {
				t.Errorf("got wrong result ProjectGrantID: expected: %v, actual: %v ", tt.res.result.GrantID, result.GrantID)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeProjectGrantMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *ProjectEventstore
		ctx    context.Context
		member *model.ProjectGrantMember
	}
	type res struct {
		result  *model.ProjectGrantMember
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change project grant member",
			args: args{
				es:  GetMockManipulateProjectWithGrantMember(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
					UserID:  "UserID",
					Roles:   []string{"RoleChanged"},
				},
			},
			res: res{
				result: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
					UserID:  "UserID",
					Roles:   []string{"RoleChanged"},
				},
			},
		},
		{
			name: "invalid member",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					Roles: []string{"Role"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
					UserID:  "UserID",
					Roles:   []string{"Role"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "user not member of grant",
			args: args{
				es:  GetMockManipulateProjectWithGrant(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
					UserID:  "UserID",
					Roles:   []string{"Role"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeProjectGrantMember(tt.args.ctx, tt.args.member)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.GrantID != tt.res.result.GrantID {
				t.Errorf("got wrong result ProjectGrantID: expected: %v, actual: %v ", tt.res.result.GrantID, result.GrantID)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRemoveProjectGrantMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *ProjectEventstore
		ctx    context.Context
		member *model.ProjectGrantMember
	}
	type res struct {
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove project grant member",
			args: args{
				es:  GetMockManipulateProjectWithGrantMember(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
					UserID:  "UserID",
					Roles:   []string{"RoleChanged"},
				},
			},
		},
		{
			name: "invalid member",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					Roles: []string{"Role"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateProjectNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
					UserID:  "UserID",
					Roles:   []string{"Role"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "user not member of grant",
			args: args{
				es:  GetMockManipulateProjectWithGrant(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				member: &model.ProjectGrantMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1},
					GrantID: "ProjectGrantID",
					UserID:  "UserID",
					Roles:   []string{"Role"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveProjectGrantMember(tt.args.ctx, tt.args.member)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
func TestChangesProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es           *ProjectEventstore
		id           string
		lastSequence uint64
		limit        uint64
	}
	type res struct {
		changes *model.ProjectChanges
		project *model.Project
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "changes from events, ok",
			args: args{
				es:           GetMockChangesProjectOK(ctrl),
				id:           "1",
				lastSequence: 0,
				limit:        0,
			},
			res: res{
				changes: &model.ProjectChanges{
					Changes: []*model.ProjectChange{
						{EventType: "", Sequence: 1, ModifierId: ""},
					},
					LastSequence: 1,
				},
				project: &model.Project{Name: "MusterProject"},
			},
		},
		{
			name: "changes from events, no events",
			args: args{
				es:           GetMockChangesProjectNoEvents(ctrl),
				id:           "2",
				lastSequence: 0,
				limit:        0,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ProjectChanges(nil, tt.args.id, tt.args.lastSequence, tt.args.limit, false)

			project := &model.Project{}
			if result != nil && len(result.Changes) > 0 {
				b, err := json.Marshal(result.Changes[0].Data)
				json.Unmarshal(b, project)
				if err != nil {
				}
			}
			if !tt.res.wantErr && result.LastSequence != tt.res.changes.LastSequence && project.Name != tt.res.project.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.changes.LastSequence, result.LastSequence)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangesApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es           *ProjectEventstore
		id           string
		secID        string
		lastSequence uint64
		limit        uint64
	}
	type res struct {
		changes *model.ApplicationChanges
		app     *model.Application
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "changes from events, ok",
			args: args{
				es:           GetMockChangesApplicationOK(ctrl),
				id:           "1",
				secID:        "AppId",
				lastSequence: 0,
				limit:        0,
			},
			res: res{
				changes: &model.ApplicationChanges{
					Changes: []*model.ApplicationChange{
						{EventType: "", Sequence: 1, ModifierId: ""},
					},
					LastSequence: 1,
				},
				app: &model.Application{Name: "MusterApp", AppID: "AppId", Type: 3},
			},
		},
		{
			name: "changes from events, no events",
			args: args{
				es:           GetMockChangesApplicationNoEvents(ctrl),
				id:           "2",
				secID:        "2",
				lastSequence: 0,
				limit:        0,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ApplicationChanges(nil, tt.args.id, tt.args.secID, tt.args.lastSequence, tt.args.limit, false)

			app := &model.Application{}
			if result != nil && len(result.Changes) > 0 {
				b, err := json.Marshal(result.Changes[0].Data)
				json.Unmarshal(b, app)
				if err != nil {
				}
			}
			if !tt.res.wantErr && result.LastSequence != tt.res.changes.LastSequence && app.Name != tt.res.app.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.changes.LastSequence, result.LastSequence)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
