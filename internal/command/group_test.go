package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddGroup(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx           context.Context
		group         *domain.Group
		resourceOwner string
		ownerID       string
	}
	type res struct {
		want *domain.Group
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid group, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				group:         &domain.Group{},
				resourceOwner: "org1",
				ownerID:       "user1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "org with resourceowner empty",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				group: &domain.Group{
					Name:        "group",
					Description: "group description",
				},
				resourceOwner: "",
				ownerID:       "user1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "org with owner empty",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				group: &domain.Group{
					Name:        "group",
					Description: "group description",
				},
				resourceOwner: "org1",
				ownerID:       "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "org with group owner, error already exists",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectPushFailed(zerrors.ThrowAlreadyExists(nil, "ERROR", "internl"),
						group.NewGroupAddedEvent(
							context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							"group",
							"group description",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "group1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				group: &domain.Group{
					Name:        "group",
					Description: "group description",
				},
				resourceOwner: "org1",
				ownerID:       "user1",
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "org with group owner, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectPush(
						group.NewGroupAddedEvent(
							context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							"group",
							"group description",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "group1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				group: &domain.Group{
					Name:        "group",
					Description: "group description",
				},
				resourceOwner: "org1",
				ownerID:       "user1",
			},
			res: res{
				want: &domain.Group{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
						AggregateID:   "group1",
					},
					Name:        "group",
					Description: "group description",
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
			c.setMilestonesCompletedForTest("instanceID")
			got, err := c.AddGroup(tt.args.ctx, tt.args.group, tt.args.resourceOwner, tt.args.ownerID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ChangeGroup(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		group         *domain.Group
		resourceOwner string
	}
	type res struct {
		want *domain.Group
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid group, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "group1",
					},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid group empty aggregateid, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					Name: "group",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "group not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "group1",
					},
					Name: "group change",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "group removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group",
								"group description"),
						),
						eventFromEventPusher(
							group.NewGroupRemovedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group"),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "group1",
					},
					Name: "group change",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group",
								"group description"),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "group1",
					},
					Name:        "group",
					Description: "group description",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "group change with name and unique constraints, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group",
								"group deacription"),
						),
					),
					expectPush(
						newGroupChangedEvent(context.Background(),
							"group1",
							"org1",
							"group",
							"group-new",
							"group description",
							"group new description",
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "group1",
					},
					Name:        "group-new",
					Description: "group description",
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Group{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "group1",
						ResourceOwner: "org1",
					},
					Name:        "group-new",
					Description: "group description",
				},
			},
		},
		{
			name: "group change without name and unique constraints, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group",
								"group description"),
						),
					),
					expectPush(
						newGroupChangedEvent(context.Background(),
							"group1",
							"org1",
							"",
							"",
							"group description",
							"group description",
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				group: &domain.Group{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "group1",
					},
					Name:        "group",
					Description: "group description",
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Group{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "group1",
						ResourceOwner: "org1",
					},
					Name:        "group",
					Description: "group description",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeGroup(tt.args.ctx, tt.args.group, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_DeactivateGroup(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		groupID       string
		resourceOwner string
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid group id, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid resourceowner, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "group1",
				resourceOwner: "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "group not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "group1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "group removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group",
								"group description"),
						),
						eventFromEventPusher(
							group.NewGroupRemovedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group"),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "group1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "group already inactive, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group",
								"group description"),
						),
						eventFromEventPusher(
							group.NewGroupDeactivatedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "group1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "group deactivate, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group",
								"group description"),
						),
					),
					expectPush(
						group.NewGroupDeactivatedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "group1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.DeactivateGroup(tt.args.ctx, tt.args.groupID, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ReactivateGroup(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		groupID       string
		resourceOwner string
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid group id, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid resourceowner, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "group1",
				resourceOwner: "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "group not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "group1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "group removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group", "group description"),
						),
						eventFromEventPusher(
							group.NewGroupRemovedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group"),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "group1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "group not inactive, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group", "group description"),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "group1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "group reactivate, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group", "group description"),
						),
						eventFromEventPusher(
							group.NewGroupDeactivatedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate),
						),
					),
					expectPush(
						group.NewGroupReactivatedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "group1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ReactivateGroup(tt.args.ctx, tt.args.groupID, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_RemoveGroup(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		groupID       string
		resourceOwner string
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid group id, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid resourceowner, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "group1",
				resourceOwner: "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "group not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "group1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "group removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group", "group description"),
						),
						eventFromEventPusher(
							group.NewGroupRemovedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group"),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "group1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "group remove, without entityConstraints, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							group.NewGroupAddedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"group", "group description"),
						),
					),
					// no saml application events
					expectFilter(),
					expectPush(
						group.NewGroupRemovedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							"group"),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				groupID:       "group1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveGroup(tt.args.ctx, tt.args.groupID, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func newGroupChangedEvent(ctx context.Context, groupID, resourceOwner, oldName, newName, oldDescription, newDescription string) *group.GroupChangeEvent {
	changes := []group.GroupChanges{}
	if newName != "" {
		changes = append(changes, group.ChangeName(newName), group.ChangeDescription(newDescription))
	}
	event, _ := group.NewGroupChangeEvent(ctx,
		&group.NewAggregate(groupID, resourceOwner).Aggregate,
		oldName,
		oldDescription,
		changes,
	)
	return event
}

func TestAddGroup(t *testing.T) {
	type args struct {
		a           *group.Aggregate
		name        string
		owner       string
		description string
	}

	ctx := context.Background()
	agg := group.NewAggregate("test", "test")

	tests := []struct {
		name string
		args args
		want Want
	}{
		{
			name: "invalid name",
			args: args{
				a:           agg,
				name:        "",
				owner:       "owner",
				description: "",
			},
			want: Want{
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "PROJE-C01yo", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid description",
			args: args{
				a:           agg,
				name:        "name",
				owner:       "owner",
				description: "",
			},
			want: Want{
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "PROJE-AO52V", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid owner",
			args: args{
				a:           agg,
				name:        "name",
				owner:       "",
				description: "description",
			},
			want: Want{
				ValidationErr: zerrors.ThrowPreconditionFailed(nil, "PROJE-hzxwo", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "correct",
			args: args{
				a:           agg,
				name:        "ZITADEL",
				owner:       "CAOS AG",
				description: "Zitadel Group",
			},
			want: Want{
				Commands: []eventstore.Command{
					group.NewGroupAddedEvent(ctx, &agg.Aggregate,
						"ZITADEL",
						"ZITADEL DEFAULT",
					),
					group.NewGroupMemberAddedEvent(ctx, &agg.Aggregate,
						"CAOS AG",
						"role1"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertValidation(t, context.Background(), AddGroupCommand(tt.args.a, tt.args.name, tt.args.owner, tt.args.description), nil, tt.want)
		})
	}
}
