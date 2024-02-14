package command

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/target"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_AddTarget(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx           context.Context
		add           *AddTarget
		resourceOwner string
	}
	type res struct {
		id      string
		details *domain.ObjectDetails
		err     func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"no name, error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				add:           &AddTarget{},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"unique constraint failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						target.NewAddedEvent(context.Background(),
							target.NewAggregate("id1", "org1"),
							"name",
							domain.TargetTypeWebhook,
							"https://example.com",
							0,
							false,
							false,
						),
					),
				),
				idGenerator: mock.ExpectID(t, "id1"),
			},
			args{
				ctx: context.Background(),
				add: &AddTarget{
					Name:          "name",
					URL:           "https://example.com",
					ExecutionType: domain.TargetTypeWebhook,
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"push ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectPush(
						target.NewAddedEvent(context.Background(),
							target.NewAggregate("id1", "org1"),
							"name",
							domain.TargetTypeWebhook,
							"https://example.com",
							0,
							false,
							false,
						),
					),
				),
				idGenerator: mock.ExpectID(t, "id1"),
			},
			args{
				ctx: context.Background(),
				add: &AddTarget{
					Name:          "name",
					ExecutionType: domain.TargetTypeWebhook,
					URL:           "https://example.com",
				},
				resourceOwner: "org1",
			},
			res{
				id: "id1",
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push full ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectPush(
						target.NewAddedEvent(context.Background(),
							target.NewAggregate("id1", "org1"),
							"name",
							domain.TargetTypeWebhook,
							"https://example.com",
							time.Second,
							true,
							true,
						),
					),
				),
				idGenerator: mock.ExpectID(t, "id1"),
			},
			args{
				ctx: context.Background(),
				add: &AddTarget{
					Name:             "name",
					ExecutionType:    domain.TargetTypeWebhook,
					URL:              "https://example.com",
					Timeout:          time.Second,
					Async:            true,
					InterruptOnError: true,
				},
				resourceOwner: "org1",
			},
			res{
				id: "id1",
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			details, err := c.AddTarget(tt.args.ctx, tt.args.add, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, tt.args.add.AggregateID)
				assert.Equal(t, tt.res.details, details)
			}
		})
	}
}

func TestCommands_ChangeTarget(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		change        *ChangeTarget
		resourceOwner string
	}
	type res struct {
		details *domain.ObjectDetails
		err     func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"id missing, error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				change: &ChangeTarget{
					Name: stringPointer("name"),
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"not found, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args{
				ctx: context.Background(),
				change: &ChangeTarget{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "id1",
					},
					Name: gu.Ptr("name"),
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"no changes",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("id1", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								0,
								false,
								false,
							),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				change: &ChangeTarget{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "id1",
					},
					ExecutionType: gu.Ptr(domain.TargetTypeWebhook),
				},
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"unique constraint failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("id1", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								0,
								false,
								false,
							),
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						target.NewChangedEvent(context.Background(),
							target.NewAggregate("id1", "org1"),
							[]target.Changes{
								target.ChangeName("name", "name2"),
							},
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				change: &ChangeTarget{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "id1",
					},
					Name: gu.Ptr("name2"),
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"push ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("id1", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								0,
								false,
								false,
							),
						),
					),
					expectPush(
						target.NewChangedEvent(context.Background(),
							target.NewAggregate("id1", "org1"),
							[]target.Changes{
								target.ChangeName("name", "name2"),
							},
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				change: &ChangeTarget{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "id1",
					},
					Name: gu.Ptr("name2"),
				},
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push full ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("id1", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								0,
								false,
								false,
							),
						),
					),
					expectPush(
						target.NewChangedEvent(context.Background(),
							target.NewAggregate("id1", "org1"),
							[]target.Changes{
								target.ChangeName("name", "name2"),
								target.ChangeURL("https://example2.com"),
								target.ChangeExecutionType(domain.TargetTypeRequestResponse),
								target.ChangeTimeout(time.Second),
								target.ChangeAsync(true),
								target.ChangeInterruptOnError(true),
							},
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				change: &ChangeTarget{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "id1",
					},
					Name:             gu.Ptr("name2"),
					URL:              gu.Ptr("https://example2.com"),
					ExecutionType:    gu.Ptr(domain.TargetTypeRequestResponse),
					Timeout:          gu.Ptr(time.Second),
					Async:            gu.Ptr(true),
					InterruptOnError: gu.Ptr(true),
				},
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			details, err := c.ChangeTarget(tt.args.ctx, tt.args.change, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.details, details)
			}
		})
	}
}

func TestCommands_DeleteTarget(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		id            string
		resourceOwner string
	}
	type res struct {
		details *domain.ObjectDetails
		err     func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"id missing, error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				id:            "",
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"not found, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"remove ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("id1", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								0,
								false,
								false,
							),
						),
					),
					expectPush(
						target.NewRemovedEvent(context.Background(),
							target.NewAggregate("id1", "org1"),
							"name",
						),
					),
				),
			},
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			details, err := c.DeleteTarget(tt.args.ctx, tt.args.id, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.details, details)
			}
		})
	}
}
