package eventsourcing

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/api/authz"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

func TestLabelPolicyAddedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.Org
		new        *iam_es_model.LabelPolicy
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
			name: "add label policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.Org{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
				},
				new: &iam_es_model.LabelPolicy{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					PrimaryColor: "000000",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.LabelPolicyAdded},
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
			name: "label policy config nil",
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
			agg, err := LabelPolicyAddedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

			if !tt.res.wantErr && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if !tt.res.wantErr && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
				if !tt.res.wantErr && agg.Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
			}

			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestLabelPolicyChangedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.Org
		new        *iam_es_model.LabelPolicy
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
			name: "change label policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.Org{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					LabelPolicy: &iam_es_model.LabelPolicy{
						ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					},
				},
				new: &iam_es_model.LabelPolicy{
					ObjectRoot:     models.ObjectRoot{AggregateID: "AggregateID"},
					PrimaryColor:   "000000",
					SecondaryColor: "FFFFFF",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.LabelPolicyChanged},
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
			name: "label policy config nil",
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
			agg, err := LabelPolicyChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

			if !tt.res.wantErr && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if !tt.res.wantErr && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
				if !tt.res.wantErr && agg.Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
			}

			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
