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
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/target"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_AddTarget(t *testing.T) {
	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id_generator.Generator
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
			"no resourceowner, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				add:           &AddTarget{},
				resourceOwner: "",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no name, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				add:           &AddTarget{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no timeout, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				add: &AddTarget{
					Name: "name",
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no Endpoint, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				add: &AddTarget{
					Name:     "name",
					Timeout:  time.Second,
					Endpoint: "",
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no parsable Endpoint, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				add: &AddTarget{
					Name:     "name",
					Timeout:  time.Second,
					Endpoint: "://",
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"unique constraint failed, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						target.NewAddedEvent(context.Background(),
							target.NewAggregate("id1", "instance"),
							"name",
							domain.TargetTypeWebhook,
							"https://example.com",
							time.Second,
							false,
						),
					),
				),
				idGenerator: mock.ExpectID(t, "id1"),
			},
			args{
				ctx: context.Background(),
				add: &AddTarget{
					Name:       "name",
					Endpoint:   "https://example.com",
					Timeout:    time.Second,
					TargetType: domain.TargetTypeWebhook,
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"already existing",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target", "instance"),
						),
					),
				),
				idGenerator: mock.ExpectID(t, "id1"),
			},
			args{
				ctx: context.Background(),
				add: &AddTarget{
					Name:       "name",
					TargetType: domain.TargetTypeWebhook,
					Timeout:    time.Second,
					Endpoint:   "https://example.com",
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			"push ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						targetAddEvent("id1", "instance"),
					),
				),
				idGenerator: mock.ExpectID(t, "id1"),
			},
			args{
				ctx: context.Background(),
				add: &AddTarget{
					Name:       "name",
					TargetType: domain.TargetTypeWebhook,
					Timeout:    time.Second,
					Endpoint:   "https://example.com",
				},
				resourceOwner: "instance",
			},
			res{
				id: "id1",
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push full ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						func() eventstore.Command {
							event := targetAddEvent("id1", "instance")
							event.InterruptOnError = true
							return event
						}(),
					),
				),
				idGenerator: mock.ExpectID(t, "id1"),
			},
			args{
				ctx: context.Background(),
				add: &AddTarget{
					Name:             "name",
					TargetType:       domain.TargetTypeWebhook,
					Endpoint:         "https://example.com",
					Timeout:          time.Second,
					InterruptOnError: true,
				},
				resourceOwner: "instance",
			},
			res{
				id: "id1",
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
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
		eventstore func(t *testing.T) *eventstore.Eventstore
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
			"resourceowner missing, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				change:        &ChangeTarget{},
				resourceOwner: "",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"id missing, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				change:        &ChangeTarget{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"name empty, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				change: &ChangeTarget{
					Name: gu.Ptr(""),
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"timeout empty, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				change: &ChangeTarget{
					Timeout: gu.Ptr(time.Duration(0)),
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"Endpoint empty, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				change: &ChangeTarget{
					Endpoint: gu.Ptr(""),
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"Endpoint not parsable, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				change: &ChangeTarget{
					Endpoint: gu.Ptr("://"),
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"not found, error",
			fields{
				eventstore: expectEventstore(
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
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"no changes",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target", "instance"),
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
					TargetType: gu.Ptr(domain.TargetTypeWebhook),
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"unique constraint failed, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target", "instance"),
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						target.NewChangedEvent(context.Background(),
							target.NewAggregate("id1", "instance"),
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
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"push ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("id1", "instance"),
						),
					),
					expectPush(
						target.NewChangedEvent(context.Background(),
							target.NewAggregate("id1", "instance"),
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
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push full ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("id1", "instance"),
						),
					),
					expectPush(
						target.NewChangedEvent(context.Background(),
							target.NewAggregate("id1", "instance"),
							[]target.Changes{
								target.ChangeName("name", "name2"),
								target.ChangeEndpoint("https://example2.com"),
								target.ChangeTargetType(domain.TargetTypeCall),
								target.ChangeTimeout(10 * time.Second),
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
					Endpoint:         gu.Ptr("https://example2.com"),
					TargetType:       gu.Ptr(domain.TargetTypeCall),
					Timeout:          gu.Ptr(10 * time.Second),
					InterruptOnError: gu.Ptr(true),
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
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
		eventstore func(t *testing.T) *eventstore.Eventstore
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
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				id:            "",
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"not found, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"remove ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("id1", "instance"),
						),
					),
					expectPush(
						target.NewRemovedEvent(context.Background(),
							target.NewAggregate("id1", "instance"),
							"name",
						),
					),
				),
			},
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
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
