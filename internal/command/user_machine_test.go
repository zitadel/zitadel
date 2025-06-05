package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddMachine(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx     context.Context
		machine *Machine
		state   *domain.UserState
		check   PermissionCheck
		options func(*Commands) []addMachineOption
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
			name: "user invalid, invalid argument error name",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
					Username: "username",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user invalid, invalid argument error username",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
					Name: "name",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "org policy not found, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(),
					expectFilter(),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
					Name:     "name",
					Username: "username",
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "add machine, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						user.NewMachineAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"name",
							"description",
							true,
							domain.OIDCTokenTypeBearer,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
					Description: "description",
					Name:        "name",
					Username:    "username",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "add machine - custom id, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("optionalID1", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						user.NewMachineAddedEvent(context.Background(),
							&user.NewAggregate("optionalID1", "org1").Aggregate,
							"username",
							"name",
							"description",
							true,
							domain.OIDCTokenTypeBearer,
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "optionalID1",
						ResourceOwner: "org1",
					},
					Description: "description",
					Name:        "name",
					Username:    "username",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "with username fallback to given username",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "aggregateID"),
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("aggregateID", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						user.NewMachineAddedEvent(context.Background(),
							&user.NewAggregate("aggregateID", "org1").Aggregate,
							"username",
							"name",
							"",
							true,
							domain.OIDCTokenTypeBearer,
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
					Name:     "name",
					Username: "username",
				},
				options: func(commands *Commands) []addMachineOption {
					return []addMachineOption{
						AddMachineWithUsernameToIDFallback(),
					}
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "with username fallback to generated id",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "aggregateID"),
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("aggregateID", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						user.NewMachineAddedEvent(context.Background(),
							&user.NewAggregate("aggregateID", "org1").Aggregate,
							"aggregateID",
							"name",
							"",
							true,
							domain.OIDCTokenTypeBearer,
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
					Name: "name",
				},
				options: func(commands *Commands) []addMachineOption {
					return []addMachineOption{
						AddMachineWithUsernameToIDFallback(),
					}
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "with username fallback to given id",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("aggregateID", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						user.NewMachineAddedEvent(context.Background(),
							&user.NewAggregate("aggregateID", "org1").Aggregate,
							"aggregateID",
							"name",
							"",
							true,
							domain.OIDCTokenTypeBearer,
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
						AggregateID:   "aggregateID",
					},
					Name: "name",
				},
				options: func(commands *Commands) []addMachineOption {
					return []addMachineOption{
						AddMachineWithUsernameToIDFallback(),
					}
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "with succeeding permission check, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						user.NewMachineAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"name",
							"description",
							true,
							domain.OIDCTokenTypeBearer,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
					Description: "description",
					Name:        "name",
					Username:    "username",
				},
				check: func(resourceOwner, aggregateID string) error {
					return nil
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "with failing permission check, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
					Description: "description",
					Name:        "name",
					Username:    "username",
				},
				check: func(resourceOwner, aggregateID string) error {
					return zerrors.ThrowPermissionDenied(nil, "", "")
				},
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "add machine, ok + deactive state",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						user.NewMachineAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"name",
							"description",
							true,
							domain.OIDCTokenTypeBearer,
						),
						user.NewUserDeactivatedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
					Description: "description",
					Name:        "name",
					Username:    "username",
				},
				state: func() *domain.UserState {
					state := domain.UserStateInactive
					return &state
				}(),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "add machine, ok + locked state",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						user.NewMachineAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"name",
							"description",
							true,
							domain.OIDCTokenTypeBearer,
						),
						user.NewUserLockedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
					Description: "description",
					Name:        "name",
					Username:    "username",
				},
				state: func() *domain.UserState {
					state := domain.UserStateLocked
					return &state
				}(),
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
				eventstore:      tt.fields.eventstore,
				idGenerator:     tt.fields.idGenerator,
				checkPermission: newMockPermissionCheckAllowed(),
			}
			var options []addMachineOption
			if tt.args.options != nil {
				options = tt.args.options(r)
			}
			got, err := r.AddMachine(tt.args.ctx, tt.args.machine, tt.args.state, tt.args.check, options...)
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

func TestCommandSide_ChangeMachine(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx     context.Context
		machine *Machine
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
			name: "user invalid, invalid argument error name",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
					Username: "username",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user invalid, invalid argument error username",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
					Name: "username",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
						AggregateID:   "user1",
					},
					Name:     "name",
					Username: "username",
				},
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
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"name",
								"description",
								true,
								domain.OIDCTokenTypeBearer,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
						AggregateID:   "user1",
					},
					Username:    "username",
					Name:        "name",
					Description: "description",
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "change machine, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"name",
								"description",
								true,
								domain.OIDCTokenTypeBearer,
							),
						),
					),
					expectPush(
						newMachineChangedEvent(context.Background(), "user1", "org1", "name1", "description1"),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				machine: &Machine{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
						AggregateID:   "user1",
					},
					Name:        "name1",
					Description: "description1",
				},
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
			got, err := r.ChangeMachine(tt.args.ctx, tt.args.machine)
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

func newMachineChangedEvent(ctx context.Context, userID, resourceOwner, name, description string) *user.MachineChangedEvent {
	event := user.NewMachineChangedEvent(ctx,
		&user.NewAggregate(userID, resourceOwner).Aggregate,
		[]user.MachineChanges{
			user.ChangeName(name),
			user.ChangeDescription(description),
		},
	)
	return event
}
