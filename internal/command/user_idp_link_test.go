package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestCommandSide_BulkAddUserIDPLinks(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
		links         []*domain.UserIDPLink
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
				links: []*domain.UserIDPLink{
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
				links: []*domain.UserIDPLink{
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
				links: []*domain.UserIDPLink{
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
				links: []*domain.UserIDPLink{
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
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewUserIDPLinkAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"config1",
									"name",
									"externaluser1",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(user.NewAddUserIDPLinkUniqueConstraint("config1", "externaluser1")),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				links: []*domain.UserIDPLink{
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
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewUserIDPLinkAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"config1",
									"name",
									"externaluser1",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(user.NewAddUserIDPLinkUniqueConstraint("config1", "externaluser1")),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				links: []*domain.UserIDPLink{
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
			err := r.BulkAddedUserIDPLinks(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, tt.args.links)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_RemoveUserIDPLink(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx  context.Context
		link *domain.UserIDPLink
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
				link: &domain.UserIDPLink{
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
				link: &domain.UserIDPLink{
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
							user.NewUserIDPLinkAddedEvent(context.Background(),
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
								nil,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				link: &domain.UserIDPLink{
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
				link: &domain.UserIDPLink{
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
							user.NewUserIDPLinkAddedEvent(context.Background(),
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
								user.NewUserIDPLinkRemovedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"config1",
									"externaluser1",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(user.NewRemoveUserIDPLinkUniqueConstraint("config1", "externaluser1")),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				link: &domain.UserIDPLink{
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
			got, err := r.RemoveUserIDPLink(tt.args.ctx, tt.args.link)
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

func TestCommandSide_ExternalLoginCheck(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx         context.Context
		orgID       string
		userID      string
		authRequest *domain.AuthRequest
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
			name: "userid missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "",
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
							user.NewUserIDPLinkAddedEvent(context.Background(),
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
								nil,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "external login check, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewUserIDPCheckSucceededEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									&user.AuthRequestInfo{
										ID:                  "request1",
										UserAgentID:         "useragent1",
										SelectedIDPConfigID: "config1",
									},
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				authRequest: &domain.AuthRequest{
					ID:                  "request1",
					AgentID:             "useragent1",
					SelectedIDPConfigID: "config1",
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
			err := r.UserIDPLoginChecked(tt.args.ctx, tt.args.orgID, tt.args.userID, tt.args.authRequest)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
