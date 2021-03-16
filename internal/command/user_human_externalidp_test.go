package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestCommandSide_BulkAddExternalIDPs(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
		externalIDPs  []*domain.ExternalIDP
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "missing userid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "",
				externalIDPs: []*domain.ExternalIDP{
					{
						IDPConfigID:    "config1",
						ExternalUserID: "externaluser1",
					},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "no external idps, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "userID doesnt match aggregate id, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				externalIDPs: []*domain.ExternalIDP{
					{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "user2",
						},
						IDPConfigID:    "config1",
						ExternalUserID: "externaluser1",
					},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid external idp, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				externalIDPs: []*domain.ExternalIDP{
					{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "user1",
						},
						IDPConfigID:    "",
						ExternalUserID: "externaluser1",
					},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "config not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				externalIDPs: []*domain.ExternalIDP{
					{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "user1",
						},
						IDPConfigID:    "config1",
						ExternalUserID: "externaluser1",
					},
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "add external idp org config, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"name",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeUnspecified,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanExternalIDPAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"config1",
									"name",
									"externaluser1",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(user.NewAddExternalIDPUniqueConstraint("config1", "externaluser1")),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				externalIDPs: []*domain.ExternalIDP{
					{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "user1",
						},
						IDPConfigID:    "config1",
						DisplayName:    "name",
						ExternalUserID: "externaluser1",
					},
				},
			},
			res: res{},
		},
		{
			name: "add external idp iam config, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							iam.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"name",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeUnspecified,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanExternalIDPAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"config1",
									"name",
									"externaluser1",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(user.NewAddExternalIDPUniqueConstraint("config1", "externaluser1")),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				externalIDPs: []*domain.ExternalIDP{
					{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "user1",
						},
						IDPConfigID:    "config1",
						DisplayName:    "name",
						ExternalUserID: "externaluser1",
					},
				},
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			err := r.BulkAddedHumanExternalIDP(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, tt.args.externalIDPs)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_RemoveExternalIDP(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx         context.Context
		externalIDP *domain.ExternalIDP
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
			name: "invalid idp, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				externalIDP: &domain.ExternalIDP{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					IDPConfigID:    "",
					ExternalUserID: "externaluser1",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "aggregate id missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				externalIDP: &domain.ExternalIDP{
					IDPConfigID:    "config1",
					ExternalUserID: "externaluser1",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "user removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanExternalIDPAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"config1",
								"name",
								"externaluser1",
							),
						),
						eventFromEventPusher(
							user.NewUserRemovedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				externalIDP: &domain.ExternalIDP{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					IDPConfigID:    "config1",
					ExternalUserID: "externaluser1",
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "external idp not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				externalIDP: &domain.ExternalIDP{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					IDPConfigID:    "config1",
					ExternalUserID: "externaluser1",
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "remove external idp, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanExternalIDPAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"config1",
								"name",
								"externaluser1",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanExternalIDPRemovedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"config1",
									"externaluser1",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(user.NewRemoveExternalIDPUniqueConstraint("config1", "externaluser1")),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				externalIDP: &domain.ExternalIDP{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					IDPConfigID:    "config1",
					ExternalUserID: "externaluser1",
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
			got, err := r.RemoveHumanExternalIDP(tt.args.ctx, tt.args.externalIDP)
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
