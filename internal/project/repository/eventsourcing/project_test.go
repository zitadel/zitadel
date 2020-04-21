package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/crypto"
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
				id:       "AggregateID",
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
		project    *Project
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
				ctx:        auth.NewMockContext("orgID", "userID"),
				aggCreator: models.NewAggregateCreator("Test"),
				project:    &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				eventLen: 0,
				aggType:  model.ProjectAggregate,
			},
		},
		{
			name: "project nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen: 0,
				aggType:  model.ProjectAggregate,
				wantErr:  true,
				errFunc:  caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.project)

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
				new:        &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
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
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName_Changed", State: int32(model.PROJECTSTATE_ACTIVE)},
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
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
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
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
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
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_INACTIVE)},
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

func TestProjectMemberAddedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *ProjectMember
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
			name: "projectmember added ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        &ProjectMember{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, UserID: "UserID", Roles: []string{"Roles"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectMemberAdded,
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
				eventType: model.ProjectMemberAdded,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectMemberAdded,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectMemberAddedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectMemberChangedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *ProjectMember
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
			name: "projectmember changed ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        &ProjectMember{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, UserID: "UserID", Roles: []string{"Roles"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectMemberChanged,
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
				eventType: model.ProjectMemberChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectMemberChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectMemberChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectMemberRemovedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *ProjectMember
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
			name: "projectmember removed ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        &ProjectMember{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, UserID: "UserID", Roles: []string{"Roles"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectMemberRemoved,
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
				eventType: model.ProjectMemberRemoved,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectMemberRemoved,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectMemberRemovedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectRoleAddedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *ProjectRole
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
			name: "projectrole added ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        &ProjectRole{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Key: "Key"},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectRoleAdded,
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
				eventType: model.ProjectRoleAdded,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectRoleAdded,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectRoleAddedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectRoleChangedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *ProjectRole
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
			name: "projectmember changed ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        &ProjectRole{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Key: "Key"},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectRoleChanged,
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
				eventType: model.ProjectRoleChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectRoleChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectRoleChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectRoleRemovedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *ProjectRole
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
			name: "projectrole changed ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        &ProjectRole{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Key: "Key"},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectRoleRemoved,
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
				eventType: model.ProjectRoleRemoved,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectRoleRemoved,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectRoleRemovedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectAppAddedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *Application
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
			name: "add oidc application",
			args: args{
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new: &Application{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppId",
					Name:       "Name",
					OIDCConfig: &OIDCConfig{AppID: "AppID", ClientID: "ClientID"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.ApplicationAdded, model.OIDCConfigAdded},
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
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
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
			agg, err := ApplicationAddedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectAppChangedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *Application
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
			name: "change app",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &Project{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					Name:       "ProjectName",
					State:      int32(model.PROJECTSTATE_ACTIVE),
					Applications: []*Application{
						&Application{AppID: "AppID", Name: "Name"},
					}},
				new: &Application{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppId",
					Name:       "NameChanged",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.ApplicationChanged},
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
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
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
			agg, err := ApplicationChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectAppRemovedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *Application
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
			name: "remove app",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &Project{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					Name:       "ProjectName",
					State:      int32(model.PROJECTSTATE_ACTIVE),
					Applications: []*Application{
						&Application{AppID: "AppID", Name: "Name"},
					}},
				new: &Application{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppId",
					Name:       "Name",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.ApplicationRemoved},
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
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
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
			agg, err := ApplicationRemovedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectAppDeactivatedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *Application
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
			name: "deactivate app",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &Project{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					Name:       "ProjectName",
					State:      int32(model.PROJECTSTATE_ACTIVE),
					Applications: []*Application{
						&Application{AppID: "AppID", Name: "Name"},
					}},
				new: &Application{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppId",
					Name:       "Name",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.ApplicationDeactivated},
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
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
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
			agg, err := ApplicationDeactivatedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectAppReactivatedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *Application
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
			name: "deactivate app",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &Project{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					Name:       "ProjectName",
					State:      int32(model.PROJECTSTATE_ACTIVE),
					Applications: []*Application{
						&Application{AppID: "AppID", Name: "Name"},
					}},
				new: &Application{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppId",
					Name:       "Name",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.ApplicationReactivated},
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
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
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
			agg, err := ApplicationReactivatedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestOIDCConfigchangAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *OIDCConfig
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
			name: "deactivate app",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &Project{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					Name:       "ProjectName",
					State:      int32(model.PROJECTSTATE_ACTIVE),
					Applications: []*Application{
						&Application{AppID: "AppID", Name: "Name", OIDCConfig: &OIDCConfig{AppID: "AppID", AuthMethodType: 1}},
					}},
				new: &OIDCConfig{
					ObjectRoot:     models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:          "AppID",
					AuthMethodType: 2,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.OIDCConfigChanged},
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
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
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
			agg, err := OIDCConfigChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestOIDCConfigSecretChangeAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *OIDCConfig
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
			name: "change client secret",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &Project{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					Name:       "ProjectName",
					State:      int32(model.PROJECTSTATE_ACTIVE),
					Applications: []*Application{
						&Application{AppID: "AppID", Name: "Name", OIDCConfig: &OIDCConfig{AppID: "AppID", AuthMethodType: 1}},
					}},
				new: &OIDCConfig{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:        "AppID",
					ClientSecret: &crypto.CryptoValue{},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.OIDCConfigSecretChanged},
			},
		},
		{
			name: "existing project nil",
			args: args{
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: nil,
				new: &OIDCConfig{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:      "AppID",
				},
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
			agg, err := OIDCConfigSecretChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new.AppID, tt.args.new.ClientSecret)(tt.args.ctx)

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

func TestProjectGrantAddedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *ProjectGrant
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
			name: "projectgrant added ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        &ProjectGrant{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, GrantID: "GrantID", GrantedOrgID: "OrgID"},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectGrantAdded,
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
				eventType: model.ProjectGrantAdded,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "grant nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectGrantAdded,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectGrantAddedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectGrantChangedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *ProjectGrant
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
			name: "change project grant",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &Project{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Name:       "ProjectName",
					State:      int32(model.PROJECTSTATE_ACTIVE),
					Grants: []*ProjectGrant{
						&ProjectGrant{GrantID: "GrantID", GrantedOrgID: "GrantedOrgID", RoleKeys: []string{"Key"}},
					}},
				new: &ProjectGrant{
					ObjectRoot:   models.ObjectRoot{AggregateID: "ID"},
					GrantID:      "GrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"KeyChanged"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.ProjectGrantChanged},
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
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "projectgrant nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
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
			agg, err := ProjectGrantChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectGrantRemovedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *ProjectGrant
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
			name: "remove app",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &Project{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Name:       "ProjectName",
					State:      int32(model.PROJECTSTATE_ACTIVE),
					Grants: []*ProjectGrant{
						&ProjectGrant{GrantID: "GrantID", GrantedOrgID: "GrantedOrgID"},
					}},
				new: &ProjectGrant{
					ObjectRoot:   models.ObjectRoot{AggregateID: "ID"},
					GrantID:      "GrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"KeyChanged"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.ProjectGrantRemoved},
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
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "projectgrant nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
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
			agg, err := ProjectGrantRemovedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectGrantDeactivatedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *ProjectGrant
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
			name: "deactivate project grant",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &Project{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Name:       "ProjectName",
					State:      int32(model.PROJECTSTATE_ACTIVE),
					Grants: []*ProjectGrant{
						&ProjectGrant{GrantID: "GrantID", GrantedOrgID: "GrantedOrgID"},
					}},
				new: &ProjectGrant{
					ObjectRoot:   models.ObjectRoot{AggregateID: "ID"},
					GrantID:      "GrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"KeyChanged"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.ProjectGrantDeactivated},
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
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "grant nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
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
			agg, err := ProjectGrantDeactivatedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectGrantReactivatedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *ProjectGrant
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
			name: "reactivate project grant",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &Project{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Name:       "ProjectName",
					State:      int32(model.PROJECTSTATE_INACTIVE),
					Grants: []*ProjectGrant{
						&ProjectGrant{GrantID: "GrantID", GrantedOrgID: "GrantedOrgID"},
					}},
				new: &ProjectGrant{
					ObjectRoot:   models.ObjectRoot{AggregateID: "ID"},
					GrantID:      "GrantID",
					GrantedOrgID: "GrantedOrgID",
					RoleKeys:     []string{"KeyChanged"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.ProjectGrantReactivated},
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
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "grant nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_INACTIVE)},
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
			agg, err := ProjectGrantReactivatedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectGrantMemberAddedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *ProjectGrantMember
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
			name: "project grant member added ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        &ProjectGrantMember{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, GrantID: "GrantID", UserID: "UserID", Roles: []string{"Roles"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectGrantMemberAdded,
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
				eventType: model.ProjectGrantMemberAdded,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectGrantMemberAdded,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectGrantMemberAddedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectGrantMemberChangedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *ProjectGrantMember
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
			name: "project grant member changed ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        &ProjectGrantMember{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", Roles: []string{"RolesChanged"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectGrantMemberChanged,
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
				eventType: model.ProjectGrantMemberChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectGrantMemberChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectGrantMemberChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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

func TestProjectGrantMemberRemovedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *Project
		new        *ProjectGrantMember
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
			name: "project grant member removed ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        &ProjectGrantMember{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, UserID: "UserID", Roles: []string{"Roles"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectGrantMemberRemoved,
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
				eventType: model.ProjectGrantMemberRemoved,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &Project{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}, Name: "ProjectName", State: int32(model.PROJECTSTATE_ACTIVE)},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.ProjectGrantMemberRemoved,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProjectGrantMemberRemovedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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
