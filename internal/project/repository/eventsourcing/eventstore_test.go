package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
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
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
			},
		},
		{
			name: "project from events, no events",
			args: args{
				es:      GetMockProjectByIDNoEvents(ctrl),
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
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
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ProjectByID(nil, tt.args.project.ID)

			if !tt.res.wantErr && result.ID != tt.res.project.ID {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.project.ID, result.ID)
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
			name: "create project, ok",
			args: args{
				es:      GetMockManipulateProject(ctrl),
				ctx:     auth.NewMockContext("orgID", "userID"),
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name"},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name"},
			},
		},
		{
			name: "create project no name",
			args: args{
				es:      GetMockManipulateProject(ctrl),
				ctx:     auth.NewMockContext("orgID", "userID"),
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.CreateProject(tt.args.ctx, tt.args.project)

			if !tt.res.wantErr && result.ID == "" {
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
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "NameNew"},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "NameNew"},
			},
		},
		{
			name: "update project no name",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: ""},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "NameNew"},
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

			if !tt.res.wantErr && result.ID == "" {
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
		es       *ProjectEventstore
		ctx      context.Context
		existing *model.Project
		new      *model.Project
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
				es:       GetMockManipulateProject(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name", State: model.PROJECTSTATE_ACTIVE},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "NameNew", State: model.PROJECTSTATE_INACTIVE},
			},
		},
		{
			name: "deactivate project with inactive state",
			args: args{
				es:       GetMockManipulateInactiveProject(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name", State: model.PROJECTSTATE_INACTIVE},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:       GetMockManipulateProjectNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name", State: model.PROJECTSTATE_ACTIVE},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.DeactivateProject(tt.args.ctx, tt.args.existing.ID)

			if !tt.res.wantErr && result.ID == "" {
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
		es       *ProjectEventstore
		ctx      context.Context
		existing *model.Project
		new      *model.Project
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
				es:       GetMockManipulateInactiveProject(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name", State: model.PROJECTSTATE_INACTIVE},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "NameNew", State: model.PROJECTSTATE_ACTIVE},
			},
		},
		{
			name: "reactivate project with inactive state",
			args: args{
				es:       GetMockManipulateProject(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name", State: model.PROJECTSTATE_ACTIVE},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:       GetMockManipulateProjectNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name", State: model.PROJECTSTATE_ACTIVE},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ReactivateProject(tt.args.ctx, tt.args.existing.ID)

			if !tt.res.wantErr && result.ID == "" {
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
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID"},
			},
			res: res{
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"Role"}},
			},
		},
		{
			name: "project member from events, no events",
			args: args{
				es:     GetMockProjectByIDNoEvents(ctrl),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID"},
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
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
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
			if !tt.res.wantErr && result.ID != tt.res.member.ID {
				t.Errorf("got wrong result id: expected: %v, actual: %v ", tt.res.member.ID, result.ID)
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
				ctx:    auth.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
			res: res{
				result: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
		},
		{
			name: "no userid",
			args: args{
				es:     GetMockManipulateProject(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Roles: []string{"Roles"}},
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
				ctx:    auth.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID"},
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
				ctx:    auth.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
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
				ctx:    auth.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
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

			if !tt.res.wantErr && result.ID == "" {
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
				ctx:    auth.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"ChangeRoles"}},
			},
			res: res{
				result: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
		},
		{
			name: "no userid",
			args: args{
				es:     GetMockManipulateProject(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Roles: []string{"ChangeRoles"}},
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
				ctx:    auth.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID"},
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
				ctx:    auth.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
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
				ctx:    auth.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"ChangeRoles"}},
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

			if !tt.res.wantErr && result.ID == "" {
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
		es       *ProjectEventstore
		ctx      context.Context
		existing *model.Project
		member   *model.ProjectMember
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
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{
					ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					Name:       "Name",
					Members:    []*model.ProjectMember{&model.ProjectMember{UserID: "UserID", Roles: []string{"Roles"}}},
				},
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID"},
			},
			res: res{
				result: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
		},
		{
			name: "no userid",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{
					ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					Name:       "Name",
					Members:    []*model.ProjectMember{&model.ProjectMember{UserID: "UserID", Roles: []string{"Roles"}}},
				},
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Roles: []string{"ChangeRoles"}},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{
					ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					Name:       "Name",
				},
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
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
				ctx:    auth.NewMockContext("orgID", "userID"),
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"ChangeRoles"}},
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
			name: "add project role, ok",
			args: args{
				es:   GetMockManipulateProject(ctrl),
				ctx:  auth.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
			},
			res: res{
				result: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
			},
		},
		{
			name: "no key",
			args: args{
				es:   GetMockManipulateProject(ctrl),
				ctx:  auth.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, DisplayName: "DisplayName", Group: "Group"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "role already existing",
			args: args{
				es:   GetMockManipulateProjectWithRole(ctrl),
				ctx:  auth.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:   GetMockManipulateProjectNoEvents(ctrl),
				ctx:  auth.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddProjectRole(tt.args.ctx, tt.args.role)

			if !tt.res.wantErr && result.ID == "" {
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
				ctx:  auth.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Key: "Key", DisplayName: "DisplayNameChanged", Group: "Group"},
			},
			res: res{
				result: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Key: "Key", DisplayName: "DisplayNameChanged", Group: "Group"},
			},
		},
		{
			name: "no key",
			args: args{
				es:   GetMockManipulateProject(ctrl),
				ctx:  auth.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, DisplayName: "DisplayName", Group: "Group"},
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
				ctx:  auth.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
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
				ctx:  auth.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
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

			if !tt.res.wantErr && result.ID == "" {
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
				ctx:  auth.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Key: "Key"},
			},
		},
		{
			name: "no key",
			args: args{
				es:   GetMockManipulateProject(ctrl),
				ctx:  auth.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, DisplayName: "DisplayName", Group: "Group"},
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
				ctx:  auth.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
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
				ctx:  auth.NewMockContext("orgID", "userID"),
				role: &model.ProjectRole{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Key: "Key", DisplayName: "DisplayName", Group: "Group"},
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
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, AppID: "AppID"},
			},
			res: res{
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID:      "AppID",
					Name:       "Name",
					OIDCConfig: &model.OIDCConfig{ClientID: "ClientID"}},
			},
		},
		{
			name: "no events for project",
			args: args{
				es:  GetMockProjectByIDNoEvents(ctrl),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, AppID: "AppID"},
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
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
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
			result, err := tt.args.es.ApplicationByIDs(nil, tt.args.app.ID, tt.args.app.AppID)
			if !tt.res.wantErr && result.ID != tt.res.app.ID {
				t.Errorf("got wrong result id: expected: %v, actual: %v ", tt.res.app.ID, result.ID)
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
		wantErr bool
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
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_CODE},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_AUTHORIZATION_CODE},
					},
				},
			},
			res: res{
				result: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					Name: "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_CODE},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_AUTHORIZATION_CODE},
					},
				},
			},
		},
		{
			name: "invalid app",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_CODE},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_AUTHORIZATION_CODE},
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
			result, err := tt.args.es.AddApplication(tt.args.ctx, tt.args.app)

			if !tt.res.wantErr && result.AppID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.OIDCConfig == nil && result.OIDCConfig.ClientSecretString == "" {
				t.Errorf("result has no client secret")
			}
			if !tt.res.wantErr && result.OIDCConfig == nil && result.OIDCConfig.ClientID == "" {
				t.Errorf("result has no clientid")
			}
			if !tt.res.wantErr && result.Name != tt.res.result.Name {
				t.Errorf("got wrong result key: expected: %v, actual: %v ", tt.res.result.Name, result.Name)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
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
				es:  GetMockManipulateProjectWithOIDCApp(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
					Name:  "NameChanged",
				},
			},
			res: res{
				result: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
					Name:  "NameChanged",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_CODE},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_AUTHORIZATION_CODE},
					},
				},
			},
		},
		{
			name: "invalid app",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_CODE},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_AUTHORIZATION_CODE},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_CODE},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_AUTHORIZATION_CODE},
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

			if !tt.res.wantErr && result.ID == "" {
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
				es:  GetMockManipulateProjectWithOIDCApp(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
				},
			},
		},
		{
			name: "no appID",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
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
				es:  GetMockManipulateProjectWithOIDCApp(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
				},
			},
			res: res{
				result: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					State: model.APPSTATE_INACTIVE,
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_CODE},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_AUTHORIZATION_CODE},
					},
				},
			},
		},
		{
			name: "no app id",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_CODE},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_AUTHORIZATION_CODE},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_CODE},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_AUTHORIZATION_CODE},
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
			result, err := tt.args.es.DeactivateApplication(tt.args.ctx, tt.args.app.ID, tt.args.app.AppID)

			if !tt.res.wantErr && result.ID == "" {
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
				es:  GetMockManipulateProjectWithOIDCApp(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
				},
			},
			res: res{
				result: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					State: model.APPSTATE_ACTIVE,
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_CODE},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_AUTHORIZATION_CODE},
					},
				},
			},
		},
		{
			name: "no app id",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_CODE},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_AUTHORIZATION_CODE},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				app: &model.Application{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					AppID: "AppID",
					Name:  "Name",
					OIDCConfig: &model.OIDCConfig{
						ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_CODE},
						GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_AUTHORIZATION_CODE},
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
			result, err := tt.args.es.ReactivateApplication(tt.args.ctx, tt.args.app.ID, tt.args.app.AppID)

			if !tt.res.wantErr && result.ID == "" {
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
				es:  GetMockManipulateProjectWithOIDCApp(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot:    es_models.ObjectRoot{ID: "ID", Sequence: 0},
					AppID:         "AppID",
					ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_ID_TOKEN},
					GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_IMPLICIT},
				},
			},
			res: res{
				result: &model.OIDCConfig{
					ObjectRoot:    es_models.ObjectRoot{ID: "ID", Sequence: 0},
					AppID:         "AppID",
					ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_ID_TOKEN},
					GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_IMPLICIT},
				},
			},
		},
		{
			name: "invalid config",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot:    es_models.ObjectRoot{ID: "ID", Sequence: 0},
					AppID:         "AppID",
					ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_ID_TOKEN},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot:    es_models.ObjectRoot{ID: "ID", Sequence: 0},
					AppID:         "AppID",
					ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_ID_TOKEN},
					GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_IMPLICIT},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot:    es_models.ObjectRoot{ID: "ID", Sequence: 0},
					AppID:         "AppID",
					ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_ID_TOKEN},
					GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_IMPLICIT},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot:    es_models.ObjectRoot{ID: "ID", Sequence: 0},
					AppID:         "AppID",
					ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_ID_TOKEN},
					GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_IMPLICIT},
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

			if !tt.res.wantErr && result.ID == "" {
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
				es:  GetMockManipulateProjectWithOIDCApp(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 0},
					AppID:      "AppID",
				},
			},
			res: res{
				result: &model.OIDCConfig{
					ObjectRoot:    es_models.ObjectRoot{ID: "ID", Sequence: 0},
					AppID:         "AppID",
					ResponseTypes: []model.OIDCResponseType{model.OIDCRESPONSETYPE_ID_TOKEN},
					GrantTypes:    []model.OIDCGrantType{model.OIDCGRANTTYPE_IMPLICIT},
				},
			},
		},
		{
			name: "no appID",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 0},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 0},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 0},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				config: &model.OIDCConfig{
					ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 0},
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
			result, err := tt.args.es.ChangeOIDCConfigSecret(tt.args.ctx, tt.args.config.ID, tt.args.config.AppID)

			if !tt.res.wantErr && result.ID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.AppID != tt.res.result.AppID {
				t.Errorf("got wrong result AppID: expected: %v, actual: %v ", tt.res.result.AppID, result.AppID)
			}
			if !tt.res.wantErr && result.ClientSecretString == "" {
				t.Errorf("got wrong result should habe client secret")
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
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, GrantID: "GrantID"},
			},
			res: res{
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID:      "GrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"Key"},
				},
			},
		},
		{
			name: "no events for project",
			args: args{
				es:    GetMockProjectByIDNoEvents(ctrl),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, GrantID: "GrantID"},
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
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
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
			result, err := tt.args.es.ProjectGrantByIDs(nil, tt.args.grant.ID, tt.args.grant.GrantID)
			if !tt.res.wantErr && result.ID != tt.res.grant.ID {
				t.Errorf("got wrong result id: expected: %v, actual: %v ", tt.res.grant.ID, result.ID)
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
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID:      "GrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"Key"},
				},
			},
			res: res{
				result: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID:      "GrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"Key"},
				},
			},
		},
		{
			name: "invalid grant",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID: "GrantID",
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
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID:      "GrantID",
					GrantedOrgID: "GrantedOrgID",
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
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID:      "GrantID",
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

func TestChangeProjectGrant(t *testing.T) {
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
			name: "change grant, ok",
			args: args{
				es:  GetMockManipulateProjectWithGrant(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID:      "GrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"Key Changed"},
				},
			},
			res: res{
				result: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID:      "GrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"Key Changed"},
				},
			},
		},
		{
			name: "invalid grant",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID:  "GrantID",
					RoleKeys: []string{"Key Changed"},
				},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID:      "GrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"Key Changed"},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID:      "GrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"Key Changed"},
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
			result, err := tt.args.es.ChangeProjectGrant(tt.args.ctx, tt.args.grant)

			if !tt.res.wantErr && result.ID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.GrantID != tt.res.result.GrantID {
				t.Errorf("got wrong result GrantID: expected: %v, actual: %v ", tt.res.result.GrantID, result.GrantID)
			}
			if !tt.res.wantErr && !reflect.DeepEqual(result.RoleKeys, tt.res.result.RoleKeys) {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.result.RoleKeys, result.GrantID)
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
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID: "GrantID",
				},
			},
		},
		{
			name: "no grantID",
			args: args{
				es:    GetMockManipulateProject(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID: "GrantID",
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
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID: "GrantID",
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
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID: "GrantID",
				},
			},
			res: res{
				result: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID: "GrantID",
					State:   model.PROJECTGRANTSTATE_INACTIVE,
				},
			},
		},
		{
			name: "no grant id",
			args: args{
				es:    GetMockManipulateProject(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID: "GrantID",
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
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID: "GrantID",
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
			result, err := tt.args.es.DeactivateProjectGrant(tt.args.ctx, tt.args.grant.ID, tt.args.grant.GrantID)

			if !tt.res.wantErr && result.ID == "" {
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
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID: "GrantID",
				},
			},
			res: res{
				result: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID: "GrantID",
					State:   model.PROJECTGRANTSTATE_ACTIVE,
				},
			},
		},
		{
			name: "no grant id",
			args: args{
				es:    GetMockManipulateProject(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
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
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID: "GrantID",
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
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.ProjectGrant{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					GrantID: "GrantID",
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
			result, err := tt.args.es.ReactivateProjectGrant(tt.args.ctx, tt.args.grant.ID, tt.args.grant.GrantID)

			if !tt.res.wantErr && result.ID == "" {
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
