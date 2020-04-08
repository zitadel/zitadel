package eventsourcing

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/api/auth"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
)

func TestProjectByIDQuery(t *testing.T) {
	type args struct {
		id       string
		sequence uint64
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
			name: "project by id query ok",
			args: args{
				id:       "ID",
				sequence: 1,
			},
			res: res{
				filterLen: 3,
			},
		},
		{
			name: "project by id query, no id",
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
			query, err := ProjectByIDQuery(tt.args.id, tt.args.sequence)
			if !tt.res.wantErr && query == nil {
				t.Errorf("query should not be nil")
			}
			if !tt.res.wantErr && len(query.Filters) != tt.res.filterLen {
				t.Errorf("got wrong filter len: expected: %v, actual: %v ", tt.res.filterLen, len(query.Filters))
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestProjectQuery(t *testing.T) {
	type args struct {
		sequence uint64
	}
	type res struct {
		filterLen int
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "project query ok",
			args: args{
				sequence: 1,
			},
			res: res{
				filterLen: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := ProjectQuery(tt.args.sequence)
			if query == nil {
				t.Errorf("query should not be nil")
			}
			if len(query.Filters) != tt.res.filterLen {
				t.Errorf("got wrong filter len: expected: %v, actual: %v ", tt.res.filterLen, len(query.Filters))
			}
		})
	}
}

func TestProjectAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		aggCreator *models.AggregateCreator
		id         string
		sequence   uint64
	}
	type res struct {
		eventLen int
		aggType  models.AggregateType
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "project update aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				aggCreator: models.NewAggregateCreator("Test"),
				id:         "ID",
				sequence:   1,
			},
			res: res{
				eventLen: 0,
				aggType:  model.ProjectAggregate,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, _ := ProjectAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.id, tt.args.sequence)

			if agg == nil {
				t.Errorf("agg should not be nil")
			}
			if len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
		})
	}
}

func TestProjectCreateAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		new        *Project
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
			name: "project update aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				new:        &Project{ObjectRoot: models.ObjectRoot{ID: "ID"}, Name: "ProjectName", State: int32(model.Active)},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectAdded,
			},
		},
		{
			name: "new project nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectAdded,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectCreateAggregate(tt.args.aggCreator, tt.args.new)(tt.args.ctx)

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

func TestProjectUpdateAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *Project
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
			name: "project update aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{ID: "ID"}, Name: "ProjectName", State: int32(model.Active)},
				new:        &Project{ObjectRoot: models.ObjectRoot{ID: "ID"}, Name: "ProjectName_Changed", State: int32(model.Active)},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectChanged,
			},
		},
		{
			name: "existing project nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "new project nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{ID: "ID"}, Name: "ProjectName", State: int32(model.Active)},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectUpdateAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectDeactivateAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
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
			name: "project deactivate aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{ID: "ID"}, Name: "ProjectName", State: int32(model.Active)},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectDeactivated,
			},
		},
		{
			name: "existing project nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectDeactivated,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectDeactivateAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

			if !tt.res.wantErr && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if !tt.res.wantErr && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestProjectReactivateAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
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
			name: "project reactivate aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{ID: "ID"}, Name: "ProjectName", State: int32(model.Inactive)},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectReactivated,
			},
		},
		{
			name: "existing project nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectReactivated,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectReactivateAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

			if !tt.res.wantErr && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if !tt.res.wantErr && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
