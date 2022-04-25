package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/repository/features"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/id/mock"
	"github.com/caos/zitadel/internal/repository/action"
	"github.com/caos/zitadel/internal/repository/org"
)

func TestCommands_AddAction(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx           context.Context
		addAction     *domain.Action
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
				ctx: context.Background(),
				addAction: &domain.Action{
					Script: "test()",
				},
				resourceOwner: "org1",
			},
			res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			"no additional allowed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								e, _ := org.NewFeaturesSetEvent(context.Background(),
									&org.NewAggregate("org1").Aggregate,
									[]features.FeaturesChanges{
										features.ChangeMaxActions(1),
									},
								)
								return e
							}(),
						),
					),
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
								"name",
								"name() {};",
								0,
								false,
							),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				addAction: &domain.Action{
					Name:   "name",
					Script: "name() {};",
				},
				resourceOwner: "org1",
			},
			res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			"unique constraint failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								e, _ := org.NewFeaturesSetEvent(context.Background(),
									&org.NewAggregate("org1").Aggregate,
									[]features.FeaturesChanges{
										features.ChangeMaxActions(2),
									},
								)
								return e
							}(),
						),
					),
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
								"name",
								"name() {};",
								0,
								false,
							),
						),
					),
					expectPushFailed(
						errors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						[]*repository.Event{
							eventFromEventPusher(
								action.NewAddedEvent(context.Background(),
									&action.NewAggregate("id1", "org1").Aggregate,
									"name",
									"name() {};",
									0,
									false,
								),
							),
						},
						uniqueConstraintsFromEventConstraint(action.NewAddActionNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator: mock.ExpectID(t, "id1"),
			},
			args{
				ctx: context.Background(),
				addAction: &domain.Action{
					Name:   "name",
					Script: "name() {};",
				},
				resourceOwner: "org1",
			},
			res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			"push ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								e, _ := org.NewFeaturesSetEvent(context.Background(),
									&org.NewAggregate("org1").Aggregate,
									[]features.FeaturesChanges{
										features.ChangeMaxActions(2),
									},
								)
								return e
							}(),
						),
					),
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
								"name",
								"name() {};",
								0,
								false,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								action.NewAddedEvent(context.Background(),
									&action.NewAggregate("id2", "org1").Aggregate,
									"name2",
									"name2() {};",
									0,
									false,
								),
							),
						},
						uniqueConstraintsFromEventConstraint(action.NewAddActionNameUniqueConstraint("name2", "org1")),
					),
				),
				idGenerator: mock.ExpectID(t, "id2"),
			},
			args{
				ctx: context.Background(),
				addAction: &domain.Action{
					Name:   "name2",
					Script: "name2() {};",
				},
				resourceOwner: "org1",
			},
			res{
				id: "id2",
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
			id, details, err := c.AddAction(tt.args.ctx, tt.args.addAction, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.details, details)
			}
		})
	}
}

func TestCommands_ChangeAction(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		changeAction  *domain.Action
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
				changeAction: &domain.Action{
					Name:   "name",
					Script: "name() {};",
				},
				resourceOwner: "org1",
			},
			res{
				err: errors.IsErrorInvalidArgument,
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
				changeAction: &domain.Action{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "id1",
					},
					Name:   "name",
					Script: "name() {};",
				},
				resourceOwner: "org1",
			},
			res{
				err: errors.IsNotFound,
			},
		},
		{
			"no changes, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
								"name",
								"name() {};",
								0,
								false,
							),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				changeAction: &domain.Action{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "id1",
					},
					Name:   "name",
					Script: "name() {};",
				},
				resourceOwner: "org1",
			},
			res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			"unique constraint failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
								"name",
								"name() {};",
								0,
								false,
							),
						),
					),
					expectPushFailed(
						errors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						[]*repository.Event{
							eventFromEventPusher(
								func() *action.ChangedEvent {
									event, _ := action.NewChangedEvent(context.Background(),
										&action.NewAggregate("id1", "org1").Aggregate,
										[]action.ActionChanges{
											action.ChangeName("name2", "name"),
											action.ChangeScript("name2() {};"),
										},
									)
									return event
								}(),
							),
						},
						uniqueConstraintsFromEventConstraint(action.NewRemoveActionNameUniqueConstraint("name", "org1")),
						uniqueConstraintsFromEventConstraint(action.NewAddActionNameUniqueConstraint("name2", "org1")),
					),
				),
			},
			args{
				ctx: context.Background(),
				changeAction: &domain.Action{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "id1",
					},
					Name:   "name2",
					Script: "name2() {};",
				},
				resourceOwner: "org1",
			},
			res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			"push ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
								"name",
								"name() {};",
								0,
								false,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								func() *action.ChangedEvent {
									event, _ := action.NewChangedEvent(context.Background(),
										&action.NewAggregate("id1", "org1").Aggregate,
										[]action.ActionChanges{
											action.ChangeName("name2", "name"),
											action.ChangeScript("name2() {};"),
										},
									)
									return event
								}(),
							),
						},
						uniqueConstraintsFromEventConstraint(action.NewRemoveActionNameUniqueConstraint("name", "org1")),
						uniqueConstraintsFromEventConstraint(action.NewAddActionNameUniqueConstraint("name2", "org1")),
					),
				),
			},
			args{
				ctx: context.Background(),
				changeAction: &domain.Action{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "id1",
					},
					Name:   "name2",
					Script: "name2() {};",
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
			details, err := c.ChangeAction(tt.args.ctx, tt.args.changeAction, tt.args.resourceOwner)
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

func TestCommands_DeactivateAction(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		actionID      string
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
				actionID:      "",
				resourceOwner: "org1",
			},
			res{
				err: errors.IsErrorInvalidArgument,
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
				actionID:      "id1",
				resourceOwner: "org1",
			},
			res{
				err: errors.IsNotFound,
			},
		},
		{
			"not active, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
								"name",
								"name() {};",
								0,
								false,
							),
						),
						eventFromEventPusher(
							action.NewDeactivatedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args{
				ctx:           context.Background(),
				actionID:      "id1",
				resourceOwner: "org1",
			},
			res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			"deactivate ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
								"name",
								"name() {};",
								0,
								false,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								action.NewDeactivatedEvent(context.Background(),
									&action.NewAggregate("id1", "org1").Aggregate,
								),
							),
						},
					),
				),
			},
			args{
				ctx:           context.Background(),
				actionID:      "id1",
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
			details, err := c.DeactivateAction(tt.args.ctx, tt.args.actionID, tt.args.resourceOwner)
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

func TestCommands_ReactivateAction(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		actionID      string
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
				actionID:      "",
				resourceOwner: "org1",
			},
			res{
				err: errors.IsErrorInvalidArgument,
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
				actionID:      "id1",
				resourceOwner: "org1",
			},
			res{
				err: errors.IsNotFound,
			},
		},
		{
			"not inactive, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
								"name",
								"name() {};",
								0,
								false,
							),
						),
					),
				),
			},
			args{
				ctx:           context.Background(),
				actionID:      "id1",
				resourceOwner: "org1",
			},
			res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			"no additional allowed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
								"name",
								"name() {};",
								0,
								false,
							),
						),
						eventFromEventPusher(
							action.NewDeactivatedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								e, _ := org.NewFeaturesSetEvent(context.Background(),
									&org.NewAggregate("org1").Aggregate,
									[]features.FeaturesChanges{
										features.ChangeMaxActions(1),
									},
								)
								return e
							}(),
						),
					),
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id2", "org1").Aggregate,
								"name2",
								"name2() {};",
								0,
								false,
							),
						),
					),
				),
			},
			args{
				ctx:           context.Background(),
				actionID:      "id1",
				resourceOwner: "org1",
			},
			res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			"reactivate ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
								"name",
								"name() {};",
								0,
								false,
							),
						),
						eventFromEventPusher(
							action.NewDeactivatedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							func() eventstore.Command {
								e, _ := org.NewFeaturesSetEvent(context.Background(),
									&org.NewAggregate("org1").Aggregate,
									[]features.FeaturesChanges{
										features.ChangeMaxActions(1),
									},
								)
								return e
							}(),
						),
					),
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
								"name",
								"name() {};",
								0,
								false,
							),
						),
						eventFromEventPusher(
							action.NewDeactivatedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								action.NewReactivatedEvent(context.Background(),
									&action.NewAggregate("id1", "org1").Aggregate,
								),
							),
						},
					),
				),
			},
			args{
				ctx:           context.Background(),
				actionID:      "id1",
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
			details, err := c.ReactivateAction(tt.args.ctx, tt.args.actionID, tt.args.resourceOwner)
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

func TestCommands_DeleteAction(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		id            string
		resourceOwner string
		flowTypes     []domain.FlowType
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
			"id or resourceOwner emtpy, error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				id:            "",
				resourceOwner: "",
			},
			res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			"action not found, error",
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
				err: errors.IsNotFound,
			},
		},
		{
			"remove ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
								"name",
								"name() {};",
								0,
								false,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								action.NewRemovedEvent(context.Background(),
									&action.NewAggregate("id1", "org1").Aggregate,
									"name",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(action.NewRemoveActionNameUniqueConstraint("name", "org1")),
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
		{
			"remove with used action ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("id1", "org1").Aggregate,
								"name",
								"name() {};",
								0,
								false,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								action.NewRemovedEvent(context.Background(),
									&action.NewAggregate("id1", "org1").Aggregate,
									"name",
								),
							),
							eventFromEventPusher(
								org.NewTriggerActionsCascadeRemovedEvent(context.Background(),
									&org.NewAggregate("org1").Aggregate,
									domain.FlowTypeExternalAuthentication,
									"id1",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(action.NewRemoveActionNameUniqueConstraint("name", "org1")),
					),
				),
			},
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "org1",
				flowTypes: []domain.FlowType{
					domain.FlowTypeExternalAuthentication,
				},
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
			details, err := c.DeleteAction(tt.args.ctx, tt.args.id, tt.args.resourceOwner, tt.args.flowTypes...)
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
