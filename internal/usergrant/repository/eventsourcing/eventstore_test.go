package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/usergrant/model"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *UserGrantEventStore
		grant *model.UserGrant
	}
	type res struct {
		grant   *model.UserGrant
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user from events, ok",
			args: args{
				es:    GetMockUserGrantByIDOK(ctrl),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", ProjectID: "ProjectID", RoleKeys: []string{"Key"}},
			},
		},
		{
			name: "no events found",
			args: args{
				es:    GetMockUserGrantByIDNoEvents(ctrl),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "no id",
			args: args{
				es:    GetMockUserGrantByIDOK(ctrl),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "removed state",
			args: args{
				es:    GetMockUserGrantByIDRemoved(ctrl),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UserGrantByID(nil, tt.args.grant.AggregateID)

			if !tt.res.wantErr && result.AggregateID != tt.res.grant.AggregateID {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.grant.AggregateID, result.AggregateID)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddUserGrant(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *UserGrantEventStore
		ctx   context.Context
		grant *model.UserGrant
	}
	type res struct {
		result  *model.UserGrant
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
				es:  GetMockManipulateUserGrant(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					ProjectID: "ProjectID",
					UserID:    "UserID",
					RoleKeys:  []string{"Key"},
				},
			},
			res: res{
				result: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					ProjectID: "ProjectID",
					UserID:    "UserID",
					RoleKeys:  []string{"Key"},
				},
			},
		},
		{
			name: "invalid grant",
			args: args{
				es:    GetMockManipulateUserGrant(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddUserGrant(tt.args.ctx, tt.args.grant)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.UserID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeUserGrant(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *UserGrantEventStore
		ctx   context.Context
		grant *model.UserGrant
	}
	type res struct {
		result  *model.UserGrant
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
				es:  GetMockManipulateUserGrant(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					RoleKeys: []string{"KeyChanged"},
				},
			},
			res: res{
				result: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					RoleKeys: []string{"KeyChanged"},
				},
			},
		},
		{
			name: "invalid grant",
			args: args{
				es:    GetMockManipulateUserGrant(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				grant: nil,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:  GetMockManipulateUserGrantNoEvents(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					RoleKeys: []string{"KeyChanged"},
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
			result, err := tt.args.es.ChangeUserGrant(tt.args.ctx, tt.args.grant)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && !reflect.DeepEqual(result.RoleKeys, tt.res.result.RoleKeys) {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.result.RoleKeys, result.RoleKeys)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRemoveUserGrant(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *UserGrantEventStore
		ctx   context.Context
		grant *model.UserGrant
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
			name: "remove grant, ok",
			args: args{
				es:    GetMockManipulateUserGrant(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
		},
		{
			name: "no grantID",
			args: args{
				es:    GetMockManipulateUserGrant(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing grant not found",
			args: args{
				es:    GetMockManipulateUserGrantNoEvents(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveUserGrant(tt.args.ctx, tt.args.grant.AggregateID)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestDeactivateUserGrant(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *UserGrantEventStore
		ctx   context.Context
		grant *model.UserGrant
	}
	type res struct {
		result  *model.UserGrant
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
				es:    GetMockManipulateUserGrant(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				result: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					ProjectID: "ProjectID",
					State:     model.USERGRANTSTATE_INACTIVE,
				},
			},
		},
		{
			name: "no grant id",
			args: args{
				es:    GetMockManipulateUserGrant(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "grant not existing",
			args: args{
				es:    GetMockManipulateUserGrantNoEvents(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.DeactivateUserGrant(tt.args.ctx, tt.args.grant.AggregateID)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.ProjectID != tt.res.result.ProjectID {
				t.Errorf("got wrong result AppID: expected: %v, actual: %v ", tt.res.result.ProjectID, result.ProjectID)
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

func TestReactivateUserGrant(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *UserGrantEventStore
		ctx   context.Context
		grant *model.UserGrant
	}
	type res struct {
		result  *model.UserGrant
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
				es:    GetMockManipulateUserGrant(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				result: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					State: model.USERGRANTSTATE_ACTIVE,
				},
			},
		},
		{
			name: "no grant id",
			args: args{
				es:    GetMockManipulateUserGrant(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "grant not existing",
			args: args{
				es:    GetMockManipulateUserGrantNoEvents(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				grant: &model.UserGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ReactivateUserGrant(tt.args.ctx, tt.args.grant.AggregateID)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
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
