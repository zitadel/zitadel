package eventsourcing

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/api/authz"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"
)

func TestUserGrantAddedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		grant      *model.UserGrant
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
			name: "usergrant added ok",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				grant:      &model.UserGrant{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", ProjectID: "ProjectID"},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.UserGrantAdded,
			},
		},
		{
			name: "grant nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				grant:      nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggregates, err := UserGrantAddedAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.grant)

			if tt.res.errFunc == nil && len(aggregates[0].Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(aggregates[0].Events))
			}
			if tt.res.errFunc == nil && aggregates[0].Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, aggregates[0].Events[0].Type.String())
			}
			if tt.res.errFunc == nil && aggregates[0].Events[0].Data == nil {
				t.Errorf("should have data in event")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUserGrantChangedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.UserGrant
		new        *model.UserGrant
		cascade    bool
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change project grant",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserID:     "UserID",
					ProjectID:  "ProjectID",
					RoleKeys:   []string{"Key"}},
				new: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserID:     "UserID",
					ProjectID:  "ProjectID",
					RoleKeys:   []string{"KeyChanged"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserGrantChanged},
			},
		},
		{
			name: "change project grant cascade",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserID:     "UserID",
					ProjectID:  "ProjectID",
					RoleKeys:   []string{"Key"}},
				new: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserID:     "UserID",
					ProjectID:  "ProjectID",
					RoleKeys:   []string{"KeyChanged"},
				},
				cascade:    true,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserGrantCascadeChanged},
			},
		},
		{
			name: "existing grant nil",
			args: args{
				ctx:      authz.NewMockContext("orgID", "userID"),
				existing: nil,
				new: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserID:     "UserID",
					ProjectID:  "ProjectID",
					RoleKeys:   []string{"Key"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "grant nil",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserID:     "UserID",
					ProjectID:  "ProjectID",
					RoleKeys:   []string{"Key"}},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserGrantChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new, tt.args.cascade)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
				if tt.res.errFunc == nil && agg.Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
			}

			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUserGrantRemovedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.UserGrant
		new        *model.UserGrant
		cascade    bool
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove app",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserID:     "UserID",
					ProjectID:  "ProjectID",
					RoleKeys:   []string{"Key"}},
				new: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserGrantRemoved},
			},
		},
		{
			name: "remove app cascade",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserID:     "UserID",
					ProjectID:  "ProjectID",
					RoleKeys:   []string{"Key"}},
				new: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
				},
				cascade:    true,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserGrantCascadeRemoved},
			},
		},
		{
			name: "existing project nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "grant nil",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserID:     "UserID",
					ProjectID:  "ProjectID",
					RoleKeys:   []string{"Key"}},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggregates, err := UserGrantRemovedAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.existing, tt.args.new, tt.args.cascade)

			if tt.res.errFunc == nil && len(aggregates[0].Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(aggregates[0].Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && aggregates[0].Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], aggregates[0].Events[i].Type.String())
				}
			}

			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUserGrantDeactivatedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.UserGrant
		new        *model.UserGrant
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "deactivate project grant",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
				},
				new: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserGrantDeactivated},
			},
		},
		{
			name: "existing grant nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "grant nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.UserGrant{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserGrantDeactivatedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
			}

			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUserGrantReactivatedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.UserGrant
		new        *model.UserGrant
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "reactivate project grant",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
				},
				new: &model.UserGrant{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserGrantReactivated},
			},
		},
		{
			name: "existing grant nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "grant nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.UserGrant{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserGrantReactivatedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
			}

			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
