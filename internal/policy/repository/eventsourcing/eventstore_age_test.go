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

func TestGetPasswordAgePolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *PolicyEventstore
		policy *model.PasswordAgePolicy
	}
	type res struct {
		policy  *model.PasswordAgePolicy
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
				es:     GetMockGetPasswordAgePolicyOK(ctrl),
				policy: &model.PasswordAgePolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				policy: &model.PasswordAgePolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
		},
		{
			name: "policy from events, no events",
			args: args{
				es:     GetMockGetPasswordAgePolicyNoEvents(ctrl),
				policy: &model.PasswordAgePolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 2}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.GetPasswordAgePolicy(nil, tt.args.policy.AggregateID)

			if !tt.res.wantErr && result.AggregateID != tt.res.policy.AggregateID {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.policy.AggregateID, result.AggregateID)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreatePasswordAgePolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *PolicyEventstore
		ctx    context.Context
		policy *model.PasswordAgePolicy
	}
	type res struct {
		policy  *model.PasswordAgePolicy
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
				es:     GetMockPasswordAgePolicyNoEvents(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				policy: &model.PasswordAgePolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID1", Sequence: 2}, Description: "Name"},
			},
			res: res{
				policy: &model.PasswordAgePolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID1", Sequence: 2}, Description: "Name"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.CreatePasswordAgePolicy(tt.args.ctx, tt.args.policy)

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

func TestUpdatePasswordAgePolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *PolicyEventstore
		ctx context.Context
		new *model.PasswordAgePolicy
	}
	type res struct {
		policy  *model.PasswordAgePolicy
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
				es:  GetMockPasswordAgePolicy(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.PasswordAgePolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Description: "NameNew"},
			},
			res: res{
				policy: &model.PasswordAgePolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Description: "NameNew"},
			},
		},
		{
			name: "existing policy not found",
			args: args{
				es:  GetMockPasswordAgePolicyNoEvents(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.PasswordAgePolicy{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Description: "NameNew"},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UpdatePasswordAgePolicy(tt.args.ctx, tt.args.new)

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
