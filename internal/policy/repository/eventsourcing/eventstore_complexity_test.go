package eventsourcing

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/caos/zitadel/internal/api/authz"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
)

func TestGetPasswordComplexityPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *PolicyEventstore
		policy *model.PasswordComplexityPolicy
	}
	type res struct {
		policy  *model.PasswordComplexityPolicy
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
				es:     GetMockGetPasswordComplexityPolicyOK(ctrl),
				policy: &model.PasswordComplexityPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				policy: &model.PasswordComplexityPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
		},
		{
			name: "policy from events, no events",
			args: args{
				es:     GetMockGetPasswordComplexityPolicyNoEvents(ctrl),
				policy: &model.PasswordComplexityPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 2}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.GetPasswordComplexityPolicy(nil, tt.args.policy.AggregateID)

			if !tt.res.wantErr && result.AggregateID != tt.res.policy.AggregateID {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.policy.AggregateID, result.AggregateID)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreatePasswordComplexityPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *PolicyEventstore
		ctx    context.Context
		policy *model.PasswordComplexityPolicy
	}
	type res struct {
		policy  *model.PasswordComplexityPolicy
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
				es:     GetMockPasswordComplexityPolicyNoEvents(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				policy: &model.PasswordComplexityPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID1", Sequence: 2}, Description: "Name"},
			},
			res: res{
				policy: &model.PasswordComplexityPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID1", Sequence: 2}, Description: "Name"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.CreatePasswordComplexityPolicy(tt.args.ctx, tt.args.policy)

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

func TestUpdatePasswordComplexityPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *PolicyEventstore
		ctx context.Context
		new *model.PasswordComplexityPolicy
	}
	type res struct {
		policy  *model.PasswordComplexityPolicy
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
				es:  GetMockPasswordComplexityPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				new: &model.PasswordComplexityPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Description: "NameNew"},
			},
			res: res{
				policy: &model.PasswordComplexityPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Description: "NameNew"},
			},
		},
		{
			name: "existing policy not found",
			args: args{
				es:  GetMockPasswordComplexityPolicyNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				new: &model.PasswordComplexityPolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Description: "NameNew"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UpdatePasswordComplexityPolicy(tt.args.ctx, tt.args.new)

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
