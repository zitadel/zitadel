package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_AddExecution(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx           context.Context
		addExecution  *Execution
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
				addExecution:  &Execution{},
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
						execution.NewAddedEvent(context.Background(),
							execution.NewAggregate("id1", "org1"),
							"name",
							domain.ExecutionTypeWebhook,
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
				addExecution: &Execution{
					Name:          "name",
					URL:           "https://example.com",
					ExecutionType: domain.ExecutionTypeWebhook,
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
						execution.NewAddedEvent(context.Background(),
							execution.NewAggregate("id1", "org1"),
							"name",
							domain.ExecutionTypeWebhook,
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
				addExecution: &Execution{
					Name:          "name",
					ExecutionType: domain.ExecutionTypeWebhook,
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
						execution.NewAddedEvent(context.Background(),
							execution.NewAggregate("id1", "org1"),
							"name",
							domain.ExecutionTypeWebhook,
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
				addExecution: &Execution{
					Name:             "name",
					ExecutionType:    domain.ExecutionTypeWebhook,
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
			details, err := c.AddExecution(tt.args.ctx, tt.args.addExecution, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, tt.args.addExecution.AggregateID)
				assert.Equal(t, tt.res.details, details)
			}
		})
	}
}

func TestCommands_ChangeExecution(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx             context.Context
		changeExecution *Execution
		resourceOwner   string
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
				changeExecution: &Execution{
					Name: "name",
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
				changeExecution: &Execution{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "id1",
					},
					Name: "name",
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
							execution.NewAddedEvent(context.Background(),
								execution.NewAggregate("id1", "org1"),
								"name",
								domain.ExecutionTypeWebhook,
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
				changeExecution: &Execution{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "id1",
					},
					ExecutionType: domain.ExecutionTypeWebhook,
					URL:           "https://example.com",
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
			"push ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewAddedEvent(context.Background(),
								execution.NewAggregate("id1", "org1"),
								"name",
								domain.ExecutionTypeWebhook,
								"https://example.com",
								0,
								false,
								false,
							),
						),
					),
					expectPush(
						execution.NewChangedEvent(context.Background(),
							execution.NewAggregate("id1", "org1"),
							[]execution.Changes{
								execution.ChangeURL("https://example2.com"),
							},
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				changeExecution: &Execution{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "id1",
					},
					URL:           "https://example2.com",
					ExecutionType: domain.ExecutionTypeWebhook,
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
							execution.NewAddedEvent(context.Background(),
								execution.NewAggregate("id1", "org1"),
								"name",
								domain.ExecutionTypeWebhook,
								"https://example.com",
								0,
								false,
								false,
							),
						),
					),
					expectPush(
						execution.NewChangedEvent(context.Background(),
							execution.NewAggregate("id1", "org1"),
							[]execution.Changes{
								execution.ChangeURL("https://example2.com"),
								execution.ChangeExecutionType(domain.ExecutionTypeRequestResponse),
								execution.ChangeTimeout(time.Second),
								execution.ChangeAsync(true),
								execution.ChangeInterruptOnError(true),
							},
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				changeExecution: &Execution{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "id1",
					},
					URL:              "https://example2.com",
					ExecutionType:    domain.ExecutionTypeRequestResponse,
					Timeout:          time.Second,
					Async:            true,
					InterruptOnError: true,
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
			details, err := c.ChangeExecution(tt.args.ctx, tt.args.changeExecution, tt.args.resourceOwner)
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

func TestCommands_DeleteExecution(t *testing.T) {
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
							execution.NewAddedEvent(context.Background(),
								execution.NewAggregate("id1", "org1"),
								"name",
								domain.ExecutionTypeWebhook,
								"https://example.com",
								0,
								false,
								false,
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("id1", "org1"),
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
			details, err := c.DeleteExecution(tt.args.ctx, tt.args.id, tt.args.resourceOwner)
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
