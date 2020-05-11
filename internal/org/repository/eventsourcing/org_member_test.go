package eventsourcing

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

func TestOrgMemberAddedAggregate(t *testing.T) {
	type res struct {
		isErr      func(error) bool
		eventCount int
	}
	type args struct {
		aggCreator *es_models.AggregateCreator
		member     *OrgMember
		ctx        context.Context
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no member",
			args: args{
				aggCreator: es_models.NewAggregateCreator("test"),
				ctx:        auth.NewMockContext("org", "user"),
				member:     nil,
			},
			res: res{
				isErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "member added sucessfully",
			args: args{
				aggCreator: es_models.NewAggregateCreator("test"),
				ctx:        auth.NewMockContext("org", "user"),
				member: &OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "asdf", Sequence: 234},
				},
			},
			res: res{
				isErr:      nil,
				eventCount: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggregate, err := OrgMemberAddedAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.member)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got: %v", err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if tt.res.isErr == nil && aggregate == nil {
				t.Error("aggregate must not be nil")
			}
			if tt.res.isErr == nil && len(aggregate.Events) != tt.res.eventCount {
				t.Error("wrong amount of events")
			}
		})
	}
}

func TestOrgMemberChangedAggregate(t *testing.T) {
	type res struct {
		isErr      func(error) bool
		eventCount int
	}
	type args struct {
		aggCreator     *es_models.AggregateCreator
		existingMember *OrgMember
		member         *OrgMember
		ctx            context.Context
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no member",
			args: args{
				aggCreator:     es_models.NewAggregateCreator("test"),
				ctx:            auth.NewMockContext("org", "user"),
				member:         nil,
				existingMember: &OrgMember{},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "no existing member",
			args: args{
				aggCreator:     es_models.NewAggregateCreator("test"),
				ctx:            auth.NewMockContext("org", "user"),
				existingMember: nil,
				member:         &OrgMember{},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "no changes",
			args: args{
				aggCreator: es_models.NewAggregateCreator("test"),
				ctx:        auth.NewMockContext("org", "user"),
				member: &OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "asdf", Sequence: 234},
				},
				existingMember: &OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "asdf", Sequence: 234},
				},
			},
			res: res{
				isErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "with changes success",
			args: args{
				aggCreator: es_models.NewAggregateCreator("test"),
				ctx:        auth.NewMockContext("org", "user"),
				member: &OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "asdf", Sequence: 234},
					Roles:      []string{"asdf"},
				},
				existingMember: &OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "asdf", Sequence: 234},
					Roles:      []string{"asdf", "woeri"},
				},
			},
			res: res{
				isErr:      nil,
				eventCount: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggregateCreator := OrgMemberChangedAggregate(tt.args.aggCreator, tt.args.existingMember, tt.args.member)
			aggregate, err := aggregateCreator(tt.args.ctx)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got: %v", err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if tt.res.isErr == nil && aggregate == nil {
				t.Error("aggregate must not be nil")
			}
			if tt.res.isErr == nil && len(aggregate.Events) != tt.res.eventCount {
				t.Error("wrong amount of events")
			}
		})
	}
}

func TestOrgMemberRemovedAggregate(t *testing.T) {
	type res struct {
		isErr      func(error) bool
		eventCount int
	}
	type args struct {
		aggCreator *es_models.AggregateCreator
		member     *OrgMember
		ctx        context.Context
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no member",
			args: args{
				aggCreator: es_models.NewAggregateCreator("test"),
				ctx:        auth.NewMockContext("org", "user"),
				member:     nil,
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "member added sucessfully",
			args: args{
				aggCreator: es_models.NewAggregateCreator("test"),
				ctx:        auth.NewMockContext("org", "user"),
				member: &OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "asdf", Sequence: 234},
				},
			},
			res: res{
				isErr:      nil,
				eventCount: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggregateCreator := OrgMemberRemovedAggregate(tt.args.aggCreator, tt.args.member)
			aggregate, err := aggregateCreator(tt.args.ctx)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got: %v", err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if tt.res.isErr == nil && aggregate == nil {
				t.Error("aggregate must not be nil")
			}
			if tt.res.isErr == nil && len(aggregate.Events) != tt.res.eventCount {
				t.Error("wrong amount of events")
			}
		})
	}
}

func Test_addMemberValidation(t *testing.T) {
	type res struct {
		isErr            func(error) bool
		preivousSequence uint64
	}
	type args struct {
		aggregate *es_models.Aggregate
		events    []*es_models.Event
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no events",
			args: args{
				aggregate: &es_models.Aggregate{},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "only org events",
			args: args{
				aggregate: &es_models.Aggregate{},
				events: []*es_models.Event{
					{
						AggregateType: "org",
						Sequence:      13,
					},
					{
						AggregateType: "org",
						Sequence:      142,
					},
				},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "only user events",
			args: args{
				aggregate: &es_models.Aggregate{},
				events: []*es_models.Event{
					{
						AggregateType: "user",
						Sequence:      13,
					},
					{
						AggregateType: "user",
						Sequence:      142,
					},
				},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "user and org events success",
			args: args{
				aggregate: &es_models.Aggregate{},
				events: []*es_models.Event{
					{
						AggregateType: "user",
						Sequence:      13,
					},
					{
						AggregateType: "org",
						Sequence:      142,
					},
				},
			},
			res: res{
				isErr:            nil,
				preivousSequence: 142,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validaiton := addMemberValidation(tt.args.aggregate)
			err := validaiton(tt.args.events...)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got: %v", err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if tt.res.isErr == nil && tt.args.aggregate.PreviousSequence != tt.res.preivousSequence {
				t.Errorf("wrong previous sequence got: %d want %d", tt.args.aggregate.PreviousSequence, tt.res.preivousSequence)
			}
		})
	}
}
