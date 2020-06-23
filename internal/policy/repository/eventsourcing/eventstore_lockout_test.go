package eventsourcing

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/api/auth"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
	"github.com/golang/mock/gomock"
)

func TestGetPasswordLockoutPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *PolicyEventstore
		policy *model.PasswordLockoutPolicy
	}
	type res struct {
		policy  *model.PasswordLockoutPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "policy from events, ok",
			args: args{
				es:     GetMockGetPasswordLockoutPolicyOK(ctrl),
				policy: &model.PasswordLockoutPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				policy: &model.PasswordLockoutPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
		},
		{
			name: "policy from events, no events",
			args: args{
				es:     GetMockGetPasswordLockoutPolicyNoEvents(ctrl),
				policy: &model.PasswordLockoutPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 2}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.GetPasswordLockoutPolicy(nil, tt.args.policy.AggregateID)

			if !tt.res.wantErr && result.AggregateID != tt.res.policy.AggregateID {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.policy.AggregateID, result.AggregateID)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreatePasswordLockoutPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *PolicyEventstore
		ctx    context.Context
		policy *model.PasswordLockoutPolicy
	}
	type res struct {
		policy  *model.PasswordLockoutPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create policy, ok",
			args: args{
				es:     GetMockPasswordLockoutPolicyNoEvents(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				policy: &model.PasswordLockoutPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID1", Sequence: 2}, Description: "Name"},
			},
			res: res{
				policy: &model.PasswordLockoutPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID1", Sequence: 2}, Description: "Name"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.CreatePasswordLockoutPolicy(tt.args.ctx, tt.args.policy)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.Description != tt.res.policy.Description {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.policy.Description, result.Description)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUpdatePasswordLockoutPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *PolicyEventstore
		ctx context.Context
		new *model.PasswordLockoutPolicy
	}
	type res struct {
		policy  *model.PasswordLockoutPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "update policy, ok",
			args: args{
				es:  GetMockPasswordLockoutPolicy(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.PasswordLockoutPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Description: "NameNew"},
			},
			res: res{
				policy: &model.PasswordLockoutPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Description: "NameNew"},
			},
		},
		{
			name: "existing policy not found",
			args: args{
				es:  GetMockPasswordLockoutPolicyNoEvents(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.PasswordLockoutPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Description: "NameNew"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UpdatePasswordLockoutPolicy(tt.args.ctx, tt.args.new)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.Description != tt.res.policy.Description {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.policy.Description, result.Description)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
