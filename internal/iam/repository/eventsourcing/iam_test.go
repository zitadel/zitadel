package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"testing"
)

func TestSetUpStartedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		iam        *model.Iam
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "setupstarted aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				iam:        &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamSetupStarted,
			},
		},
		{
			name: "iam nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				iam:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamSetupStarted,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IamSetupStartedAggregate(tt.args.aggCreator, tt.args.iam)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestSetUpDoneAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.Iam
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "setup done aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamSetupDone,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamSetupDone,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IamSetupDoneAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestGlobalOrgAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.Iam
		orgID      string
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "global org set aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				orgID:      "orgID",
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.GlobalOrgSet,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   nil,
				orgID:      "orgID",
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "global org empty",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IamSetGlobalOrgAggregate(tt.args.aggCreator, tt.args.existing, tt.args.orgID)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestIamProjectAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.Iam
		projectID  string
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "iam project id set aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				projectID:  "projectID",
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamProjectSet,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   nil,
				projectID:  "projectID",
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "project id empty",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IamSetIamProjectAggregate(tt.args.aggCreator, tt.args.existing, tt.args.projectID)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestIamMemberAddedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.Iam
		new        *model.IamMember
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "iammember added ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				new:        &model.IamMember{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, UserID: "UserID", Roles: []string{"Roles"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberAdded,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberAdded,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberAdded,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IamMemberAddedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc == nil && agg.Events[0].Data == nil {
				t.Errorf("should have data in event")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestIamMemberChangedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.Iam
		new        *model.IamMember
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
			name: "iammember changed ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				new:        &model.IamMember{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, UserID: "UserID", Roles: []string{"Roles"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberChanged,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IamMemberChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestIamMemberRemovedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.Iam
		new        *model.IamMember
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
			name: "iammember removed ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				new:        &model.IamMember{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, UserID: "UserID", Roles: []string{"Roles"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberRemoved,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberRemoved,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberRemoved,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IamMemberRemovedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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
