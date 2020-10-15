package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"testing"
)

func TestPasswordLockoutPolicyAddedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.Org
		new        *iam_es_model.PasswordLockoutPolicy
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		wantErr    bool
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add password lockout polciy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.Org{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
				},
				new: &iam_es_model.PasswordLockoutPolicy{
					ObjectRoot:  models.ObjectRoot{AggregateID: "AggregateID"},
					MaxAttempts: 10,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.PasswordLockoutPolicyAdded},
			},
		},
		{
			name: "existing org nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "password lockout policy config nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Org{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordLockoutPolicyAddedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if !tt.res.wantErr && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
				if agg.Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
			}
		})
	}
}

func TestPasswordLockoutPolicyChangedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.Org
		new        *iam_es_model.PasswordLockoutPolicy
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		wantErr    bool
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change password lockout policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.Org{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					PasswordLockoutPolicy: &iam_es_model.PasswordLockoutPolicy{
						ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					},
				},
				new: &iam_es_model.PasswordLockoutPolicy{
					ObjectRoot:          models.ObjectRoot{AggregateID: "AggregateID"},
					MaxAttempts:         5,
					ShowLockOutFailures: true,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.PasswordLockoutPolicyChanged},
			},
		},
		{
			name: "existing org nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "password lockout policy config nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Org{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordLockoutPolicyChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if !tt.res.wantErr && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
				if agg.Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
			}
		})
	}
}

func TestPasswordLockoutPolicyRemovedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.Org
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		wantErr    bool
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove password lockout policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.Org{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					PasswordLockoutPolicy: &iam_es_model.PasswordLockoutPolicy{
						ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.PasswordLockoutPolicyRemoved},
			},
		},
		{
			name: "existing org nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordLockoutPolicyRemovedAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
			}
		})
	}
}
