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
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor-org", ChangeDate: time.Now(), CreationDate: time.Now()}},
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
		ctx   context.Context
		orgID string
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
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "push failed",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil).
				expectAggregateCreator().
				expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test"))},
			args: args{
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "push correct",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil).
				expectAggregateCreator().
				expectPushEvents(6, nil)},
			args: args{
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
		{
			name: "org already inactive error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgInactiveEvent()}, nil).
				expectAggregateCreator().
				expectPushEvents(6, nil)},
			args: args{
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Eventstore.DeactivateOrg(tt.args.ctx, tt.args.orgID)
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
		ctx   context.Context
		orgID string
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
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "push failed",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgInactiveEvent()}, nil).
				expectAggregateCreator().
				expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test"))},
			args: args{
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "push correct",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgInactiveEvent()}, nil).
				expectAggregateCreator().
				expectPushEvents(6, nil)},
			args: args{
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
		{
			name: "org already active error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil).
				expectAggregateCreator().
				expectPushEvents(6, nil)},
			args: args{
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Eventstore.ReactivateOrg(tt.args.ctx, tt.args.orgID)
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
		ctx    context.Context
		member *org_model.OrgMember
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
				ctx:    auth.NewMockContext("user", "org"),
				member: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name:   "no aggregate id in input member",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx:    auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4}, UserID: "asdf"},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name:   "no aggregate id in input member",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx:    auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "asdf"}},
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
				ctx:    auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "plants"}, UserID: "banana"},
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
				ctx:    auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "plants"}, UserID: "banana"},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "new events found and added success",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{
				{Sequence: 6, Data: []byte("{\"userId\": \"banana\", \"roles\": [\"bananaa\"]}"), Type: org_model.OrgMemberChanged},
			}, nil)},
			args: args{
				ctx:    auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "plants", ChangeDate: time.Now(), CreationDate: time.Now()}, UserID: "banana"},
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
		{
			name: "not member of org error",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{
				{Sequence: 6, Data: []byte("{\"userId\": \"banana\", \"roles\": [\"bananaa\"]}"), Type: org_model.OrgMemberAdded},
				{Sequence: 7, Data: []byte("{\"userId\": \"apple\"}"), Type: org_model.OrgMemberRemoved},
			}, nil)},
			args: args{
				ctx:    auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "plants", ChangeDate: time.Now(), CreationDate: time.Now()}, UserID: "apple"},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Eventstore.OrgMemberByIDs(tt.args.ctx, tt.args.member)
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

func TestOrgEventstore_AddOrgMember(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		ctx    context.Context
		member *org_model.OrgMember
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
				ctx:    auth.NewMockContext("user", "org"),
				member: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name: "push failed",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
				}, nil).
				expectAggregateCreator().
				expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test"))},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "push correct",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
				}, nil).
				expectPushEvents(6, nil),
			},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
		{
			name: "member already exists error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						Type:     org_model.OrgMemberAdded,
						Data:     []byte(`{"userId": "hodor", "roles": ["master"]}`),
						Sequence: 6,
					},
				}, nil).
				expectPushEvents(0, errors.ThrowAlreadyExists(nil, "EVENT-yLTI6", "weiss nöd wie teste")),
			},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorAlreadyExists,
			},
		},
		{
			name: "member deleted success",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectPushEvents(10, nil).
				expectFilterEvents([]*es_models.Event{
					{
						Type:     org_model.OrgMemberAdded,
						Data:     []byte(`{"userId": "hodor", "roles": ["master"]}`),
						Sequence: 6,
					},
					{
						Type:     org_model.OrgMemberRemoved,
						Data:     []byte(`{"userId": "hodor"}`),
						Sequence: 10,
					},
				}, nil),
			},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 10,
				isErr:            nil,
			},
		},
		{
			name: "org not exists error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents(nil, nil).
				expectPushEvents(0, errors.ThrowAlreadyExists(nil, "EVENT-yLTI6", "weiss nöd wie teste")),
			},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorAlreadyExists,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Eventstore.AddOrgMember(tt.args.ctx, tt.args.member)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should not be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_ChangeOrgMember(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		isErr            func(error) bool
		expectedSequence uint64
	}
	type args struct {
		ctx    context.Context
		member *org_model.OrgMember
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
				ctx:    auth.NewMockContext("user", "org"),
				member: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name: "member not found error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgMemberAdded,
						Data:        []byte(`{"userId": "brudi", "roles": ["master of desaster"]}`),
						Sequence:    6,
					},
				}, nil),
			},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsNotFound,
			},
		},
		{
			name: "member found no changes error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgMemberAdded,
						Data:        []byte(`{"userId": "hodor", "roles": ["master"]}`),
						Sequence:    6,
					},
				}, nil),
			},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "push error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgMemberAdded,
						Data:        []byte(`{"userId": "hodor", "roles": ["master"]}`),
						Sequence:    6,
					},
				}, nil).
				expectPushEvents(0, errors.ThrowInternal(nil, "PEVENT-3wqa2", "test")),
			},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master of desaster"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "change success",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgMemberAdded,
						Data:        []byte(`{"userId": "hodor", "roles": ["master"]}`),
						Sequence:    6,
					},
				}, nil).
				expectPushEvents(7, nil),
			},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master of desaster"},
				},
			},
			res: res{
				expectedSequence: 7,
				isErr:            nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &OrgEventstore{
				Eventstore: tt.fields.Eventstore,
			}
			got, err := es.ChangeOrgMember(tt.args.ctx, tt.args.member)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should not be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_RemoveOrgMember(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		isErr func(error) bool
	}
	type args struct {
		ctx    context.Context
		member *org_model.OrgMember
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
				ctx:    auth.NewMockContext("user", "org"),
				member: nil,
			},
			res: res{
				isErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "member not found error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgMemberAdded,
						Data:        []byte(`{"userId": "brudi", "roles": ["master of desaster"]}`),
						Sequence:    6,
					},
				}, nil),
			},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master"},
				},
			},
			res: res{
				isErr: nil,
			},
		},
		{
			name: "push error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgMemberAdded,
						Data:        []byte(`{"userId": "hodor", "roles": ["master"]}`),
						Sequence:    6,
					},
				}, nil).
				expectPushEvents(0, errors.ThrowInternal(nil, "PEVENT-3wqa2", "test")),
			},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
				},
			},
			res: res{
				isErr: errors.IsInternal,
			},
		},
		{
			name: "remove success",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
					{
						AggregateID: "hodor-org",
						Type:        org_model.OrgMemberAdded,
						Data:        []byte(`{"userId": "hodor", "roles": ["master"]}`),
						Sequence:    6,
					},
				}, nil).
				expectPushEvents(7, nil),
			},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
				},
			},
			res: res{
				isErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &OrgEventstore{
				Eventstore: tt.fields.Eventstore,
			}
			err := es.RemoveOrgMember(tt.args.ctx, tt.args.member)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
		})
	}
}

func orgCreatedEvent() *es_models.Event {
	return &es_models.Event{
		AggregateID:      "hodor-org",
		AggregateType:    org_model.OrgAggregate,
		AggregateVersion: "v1",
		CreationDate:     time.Now().Add(-1 * time.Minute),
		Data:             []byte(`{"name": "hodor-org", "domain":"hodor.org"}`),
		EditorService:    "testsvc",
		EditorUser:       "testuser",
		ID:               "sdlfö4t23kj",
		ResourceOwner:    "hodor-org",
		Sequence:         32,
		Type:             org_model.OrgAdded,
	}
}

func orgInactiveEvent() *es_models.Event {
	return &es_models.Event{
		AggregateID:      "hodor-org",
		AggregateType:    org_model.OrgAggregate,
		AggregateVersion: "v1",
		CreationDate:     time.Now().Add(-1 * time.Minute),
		Data:             nil,
		EditorService:    "testsvc",
		EditorUser:       "testuser",
		ID:               "sdlfö4t23kj",
		ResourceOwner:    "hodor-org",
		Sequence:         52,
		Type:             org_model.OrgDeactivated,
	}
}
