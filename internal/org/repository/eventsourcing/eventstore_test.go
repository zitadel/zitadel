package eventsourcing

import (
	"context"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/errors"
	es_mock "github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/golang/mock/gomock"
)

type testOrgEventstore struct {
	OrgEventstore
	mockEventstore *es_mock.MockEventstore
}

func newTestEventstore(t *testing.T) *testOrgEventstore {
	mock := mockEventstore(t)
	return &testOrgEventstore{OrgEventstore: OrgEventstore{Eventstore: mock}, mockEventstore: mock}
}

func (es *testOrgEventstore) expectFilterEvents(events []*es_models.Event, err error) *testOrgEventstore {
	es.mockEventstore.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, err)
	return es
}

func (es *testOrgEventstore) expectPushEvents(startSequence uint64, err error) *testOrgEventstore {
	es.mockEventstore.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return err
		})
	return es
}

func (es *testOrgEventstore) expectAggregateCreator() *testOrgEventstore {
	es.mockEventstore.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	return es
}

func mockEventstore(t *testing.T) *es_mock.MockEventstore {
	ctrl := gomock.NewController(t)
	e := es_mock.NewMockEventstore(ctrl)

	return e
}

func TestOrgEventstore_OrgByID(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		ctx context.Context
		org *org_model.Org
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no input org",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name:   "no aggregate id in input org",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4}},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name:   "no events found success",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{}, nil)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}},
			},
			res: res{
				expectedSequence: 4,
				isErr:            nil,
			},
		},
		{
			name:   "filter fail",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{}, errors.ThrowInternal(nil, "EVENT-SAa1O", "message"))},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "new events found and added success",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{
				{Sequence: 6},
			}, nil)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor", ChangeDate: time.Now(), CreationDate: time.Now()}},
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Eventstore.OrgByID(tt.args.ctx, tt.args.org)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_DeactivateOrg(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		ctx context.Context
		org *org_model.Org
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no input org",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name:   "push failed",
			fields: fields{Eventstore: newTestEventstore(t).expectAggregateCreator().expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test"))},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name:   "push failed",
			fields: fields{Eventstore: newTestEventstore(t).expectAggregateCreator().expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test"))},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name:   "push correct",
			fields: fields{Eventstore: newTestEventstore(t).expectAggregateCreator().expectPushEvents(6, nil)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}},
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
		{
			name:   "org already inactive error",
			fields: fields{Eventstore: newTestEventstore(t).expectAggregateCreator().expectPushEvents(6, nil)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}, State: org_model.ORGSTATE_INACTIVE},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Eventstore.DeactivateOrg(tt.args.ctx, tt.args.org)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_ReactivateOrg(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		ctx context.Context
		org *org_model.Org
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no input org",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name:   "push failed",
			fields: fields{Eventstore: newTestEventstore(t).expectAggregateCreator().expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test"))},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}, State: org_model.ORGSTATE_INACTIVE},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name:   "push failed",
			fields: fields{Eventstore: newTestEventstore(t).expectAggregateCreator().expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test"))},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}, State: org_model.ORGSTATE_INACTIVE},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name:   "push correct",
			fields: fields{Eventstore: newTestEventstore(t).expectAggregateCreator().expectPushEvents(6, nil)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}, State: org_model.ORGSTATE_INACTIVE},
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
		{
			name:   "org already inactive error",
			fields: fields{Eventstore: newTestEventstore(t).expectAggregateCreator().expectPushEvents(6, nil)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}, State: org_model.ORGSTATE_ACTIVE},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Eventstore.ReactivateOrg(tt.args.ctx, tt.args.org)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_OrgMemberByIDs(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		ctx context.Context
		org *org_model.Org
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no input member",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name:   "no aggregate id in input org",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4}},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name:   "no events found success",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{}, nil)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}},
			},
			res: res{
				expectedSequence: 4,
				isErr:            nil,
			},
		},
		{
			name:   "filter fail",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{}, errors.ThrowInternal(nil, "EVENT-SAa1O", "message"))},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "new events found and added success",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{
				{Sequence: 6},
			}, nil)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor", ChangeDate: time.Now(), CreationDate: time.Now()}},
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Eventstore.OrgByID(tt.args.ctx, tt.args.org)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}
