package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestCommandSide_AddMachine(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsPreconditionFailed,
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
						[]*repository.Event{
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
						},
						uniqueConstraintsFromEventConstraint(user.NewAddUsernameUniqueConstraint("username", "org1", true)),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			got, err := r.AddMachine(tt.args.ctx, tt.args.machine)
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsNotFound,
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
				err: caos_errs.IsPreconditionFailed,
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
						[]*repository.Event{
							eventFromEventPusher(
								newMachineChangedEvent(context.Background(), "user1", "org1", "name1", "description1"),
							),
						},
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
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func newMachineChangedEvent(ctx context.Context, userID, resourceOwner, name, description string) *user.MachineChangedEvent {
	event, _ := user.NewMachineChangedEvent(ctx,
		&user.NewAggregate(userID, resourceOwner).Aggregate,
		[]user.MachineChanges{
			user.ChangeName(name),
			user.ChangeDescription(description),
		},
	)
	return event
}
