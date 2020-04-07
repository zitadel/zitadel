package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	"github.com/golang/mock/gomock"
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
			result, err := tt.args.es.ProjectByID(nil, tt.args.project)

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
			name: "update project, ok",
			args: args{
				es:       GetMockManipulateProject(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name"},
				new:      &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "NameNew"},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "NameNew"},
			},
		},
		{
			name: "update project no name",
			args: args{
				es:       GetMockManipulateProject(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name"},
				new:      &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: ""},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UpdateProject(tt.args.ctx, tt.args.existing, tt.args.new)

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
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name", State: model.Active},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "NameNew", State: model.Inactive},
			},
		},
		{
			name: "deactivate project with inactive state",
			args: args{
				es:       GetMockManipulateProject(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name", State: model.Inactive},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.DeactivateProject(tt.args.ctx, tt.args.existing)

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
				es:       GetMockManipulateProject(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name", State: model.Inactive},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "NameNew", State: model.Active},
			},
		},
		{
			name: "reactivate project with inactive state",
			args: args{
				es:       GetMockManipulateProject(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name", State: model.Active},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ReactivateProject(tt.args.ctx, tt.args.existing)

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
			name: "add project member, ok",
			args: args{
				es:       GetMockManipulateProject(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name"},
				member:   &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
			res: res{
				result: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
		},
		{
			name: "no userid",
			args: args{
				es:       GetMockManipulateProject(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name"},
				member:   &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Roles: []string{"Roles"}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "no roles",
			args: args{
				es:       GetMockManipulateProject(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name"},
				member:   &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member already existing",
			args: args{
				es:       GetMockManipulateProject(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name", Members: []*model.ProjectMember{&model.ProjectMember{UserID: "UserID"}}},
				member:   &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorAlreadyExists,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddProjectMember(tt.args.ctx, tt.args.existing, tt.args.member)

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
			name: "add project member, ok",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{
					ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					Name:       "Name",
					Members:    []*model.ProjectMember{&model.ProjectMember{UserID: "UserID", Roles: []string{"Roles"}}},
				},
				member: &model.ProjectMember{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, UserID: "UserID", Roles: []string{"ChangeRoles"}},
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
			name: "no roles",
			args: args{
				es:  GetMockManipulateProject(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.Project{
					ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1},
					Name:       "Name",
					Members:    []*model.ProjectMember{&model.ProjectMember{UserID: "UserID", Roles: []string{"Roles"}}},
				},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeProjectMember(tt.args.ctx, tt.args.existing, tt.args.member)

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
			name: "add project member, ok",
			args: args{
				es:  GetMockManipulateProject(ctrl),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveProjectMember(tt.args.ctx, tt.args.existing, tt.args.member)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
