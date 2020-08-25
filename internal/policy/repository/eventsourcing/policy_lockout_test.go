package eventsourcing

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/api/authz"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	policy_model "github.com/caos/zitadel/internal/policy/model"
)

func TestGetPasswordLockoutPolicyQuery(t *testing.T) {
	type args struct {
		recourceOwner string
		sequence      uint64
	}
	type res struct {
		filterLen int
		wantErr   bool
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "Get password lockout policy query ok",
			args: args{
				recourceOwner: "org",
				sequence:      14,
			},
			res: res{
				filterLen: 3,
			},
		},
		{
			name: "Get password lockout policy query, no org",
			args: args{
				sequence: 1,
			},
			res: res{
				filterLen: 3,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := PasswordLockoutPolicyQuery(tt.args.recourceOwner, tt.args.sequence)
			if !tt.res.wantErr && query == nil {
				t.Errorf("query should not be nil")
			}
			if !tt.res.wantErr && len(query.Filters) != tt.res.filterLen {
				t.Errorf("got wrong filter len: expected: %v, actual: %v ", tt.res.filterLen, len(query.Filters))
			}
		})
	}
}

func TestPasswordLockoutPolicyAggregate(t *testing.T) {

	type args struct {
		ctx        context.Context
		aggCreator *models.AggregateCreator
		policy     *PasswordLockoutPolicy
	}
	type res struct {
		eventLen int
		aggType  models.AggregateType
		wantErr  bool
		errFunc  func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create aggregate",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				aggCreator: models.NewAggregateCreator("Test"),
				policy:     &PasswordLockoutPolicy{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Description: "Test"},
			},
			res: res{
				eventLen: 0,
				aggType:  policy_model.PasswordLockoutPolicyAggregate,
			},
		},
		{
			name: "policy nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen: 0,
				aggType:  policy_model.PasswordLockoutPolicyAggregate,
				wantErr:  true,
				errFunc:  caos_errs.IsPreconditionFailed,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordLockoutPolicyAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.policy)

			if !tt.res.wantErr && agg == nil {
				t.Errorf("agg should not be nil")
			}
			if !tt.res.wantErr && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestPasswordLockoutPolicyCreateAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		newPolicy  *PasswordLockoutPolicy
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		wantErr   bool
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "policy update aggregate ok",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				newPolicy:  &PasswordLockoutPolicy{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Description: "PolicyName", State: int32(policy_model.PolicyStateActive)},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: policy_model.PasswordLockoutPolicyAdded,
			},
		},
		{
			name: "new policy nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				newPolicy:  nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: policy_model.PasswordLockoutPolicyAdded,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordLockoutPolicyCreateAggregate(tt.args.aggCreator, tt.args.newPolicy)(tt.args.ctx)

			if !tt.res.wantErr && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if !tt.res.wantErr && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if !tt.res.wantErr && agg.Events[0].Data == nil {
				t.Errorf("should have data in event")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestPasswordLockoutPolicyUpdateAggregate(t *testing.T) {
	type args struct {
		ctx            context.Context
		existingPolicy *PasswordLockoutPolicy
		newPolicy      *PasswordLockoutPolicy
		aggCreator     *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		wantErr   bool
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "policy update aggregate ok",
			args: args{
				ctx:            authz.NewMockContext("orgID", "userID"),
				existingPolicy: &PasswordLockoutPolicy{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Description: "PolicyName", State: int32(policy_model.PolicyStateActive)},
				newPolicy:      &PasswordLockoutPolicy{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Description: "PolicyName_Changed", State: int32(policy_model.PolicyStateActive)},
				aggCreator:     models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: policy_model.PasswordLockoutPolicyChanged,
			},
		},
		{
			name: "existing policy nil",
			args: args{
				ctx:            authz.NewMockContext("orgID", "userID"),
				existingPolicy: nil,
				aggCreator:     models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: policy_model.PasswordLockoutPolicyChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "new policy nil",
			args: args{
				ctx:            authz.NewMockContext("orgID", "userID"),
				existingPolicy: &PasswordLockoutPolicy{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Description: "ProjectName", State: int32(policy_model.PolicyStateActive)},
				newPolicy:      nil,
				aggCreator:     models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: policy_model.PasswordLockoutPolicyChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordLockoutPolicyUpdateAggregate(tt.args.aggCreator, tt.args.existingPolicy, tt.args.newPolicy)(tt.args.ctx)

			if !tt.res.wantErr && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if !tt.res.wantErr && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if !tt.res.wantErr && agg.Events[0].Data == nil {
				t.Errorf("should have data in event")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
