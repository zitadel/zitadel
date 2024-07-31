package command

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
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
		err error
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
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-03j8f", "Errors.IDMissing"),
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
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-Ek9s", "Errors.User.ExternalIDP.MinimumExternalIDPNeeded"),
			},
		},
		{
			name: "userID doesnt match aggregate id, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(
								context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"userName",
								"firstName",
								"lastName",
								"nickName",
								"displayName",
								language.German,
								domain.GenderFemale,
								"email@Address.ch",
								false,
							),
						),
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
							AggregateID: "user2",
						},
						IDPConfigID:    "config1",
						ExternalUserID: "externaluser1",
					},
				},
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-33M0g", "Errors.IDMissing"),
			},
		},
		{
			name: "invalid external idp, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(
								context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"userName",
								"firstName",
								"lastName",
								"nickName",
								"displayName",
								language.German,
								domain.GenderFemale,
								"email@Address.ch",
								false,
							),
						),
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
						IDPConfigID:    "",
						ExternalUserID: "externaluser1",
					},
				},
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-6m9Kd", "Errors.User.ExternalIDP.Invalid"),
			},
		},
		{
			name: "config not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(
								context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"userName",
								"firstName",
								"lastName",
								"nickName",
								"displayName",
								language.German,
								domain.GenderFemale,
								"email@Address.ch",
								false,
							),
						),
					),
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
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-as02jin", "Errors.IDPConfig.NotExisting"),
			},
		},
		{
			name: "linking not allowed, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(
								context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"userName",
								"firstName",
								"lastName",
								"nickName",
								"displayName",
								language.German,
								domain.GenderFemale,
								"email@Address.ch",
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
								"name",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeUnspecified,
								true,
							),
						),
						eventFromEventPusher(
							org.NewIDPOIDCConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"clientID",
								"config1",
								"issuer",
								"authEndpoint",
								"tokenEndpoint",
								nil,
								domain.OIDCMappingFieldUnspecified,
								domain.OIDCMappingFieldUnspecified,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
								"name",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeUnspecified,
								true,
							),
						),
						eventFromEventPusher(
							org.NewIDPOIDCConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"clientID",
								"config1",
								"issuer",
								"authEndpoint",
								"tokenEndpoint",
								nil,
								domain.OIDCMappingFieldUnspecified,
								domain.OIDCMappingFieldUnspecified,
							),
						),
						eventFromEventPusher(
							func() eventstore.Command {
								e, _ := org.NewOIDCIDPChangedEvent(context.Background(),
									&org.NewAggregate("org1").Aggregate,
									"config1",
									[]idp.OIDCIDPChanges{
										idp.ChangeOIDCOptions(idp.OptionChanges{IsLinkingAllowed: gu.Ptr(false)}),
									},
								)
								return e
							}(),
						),
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
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sfee2", "Errors.ExternalIDP.LinkingNotAllowed"),
			},
		},
		{
			name: "add external idp org config, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(
								context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"userName",
								"firstName",
								"lastName",
								"nickName",
								"displayName",
								language.German,
								domain.GenderFemale,
								"email@Address.ch",
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
								"name",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeUnspecified,
								true,
							),
						),
						eventFromEventPusher(
							org.NewIDPOIDCConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"clientID",
								"config1",
								"issuer",
								"authEndpoint",
								"tokenEndpoint",
								nil,
								domain.OIDCMappingFieldUnspecified,
								domain.OIDCMappingFieldUnspecified,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
								"name",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeUnspecified,
								true,
							),
						),
						eventFromEventPusher(
							org.NewIDPOIDCConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"clientID",
								"config1",
								"issuer",
								"authEndpoint",
								"tokenEndpoint",
								nil,
								domain.OIDCMappingFieldUnspecified,
								domain.OIDCMappingFieldUnspecified,
							),
						),
					),
					expectPush(
						user.NewUserIDPLinkAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"config1",
							"name",
							"externaluser1",
						),
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
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(
								context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"userName",
								"firstName",
								"lastName",
								"nickName",
								"displayName",
								language.German,
								domain.GenderFemale,
								"email@Address.ch",
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusherWithInstanceID("instance1",
							instance.NewIDPConfigAddedEvent(context.Background(),
								&instance.NewAggregate("instance1").Aggregate,
								"config1",
								"name",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeUnspecified,
								true,
							),
						),
						eventFromEventPusherWithInstanceID("instance1",
							instance.NewIDPOIDCConfigAddedEvent(context.Background(),
								&instance.NewAggregate("instance1").Aggregate,
								"clientID",
								"config1",
								"issuer",
								"authEndpoint",
								"tokenEndpoint",
								nil,
								domain.OIDCMappingFieldUnspecified,
								domain.OIDCMappingFieldUnspecified,
							),
						),
					),
					expectFilter(
						eventFromEventPusherWithInstanceID("instance1",
							instance.NewIDPConfigAddedEvent(context.Background(),
								&instance.NewAggregate("instance1").Aggregate,
								"config1",
								"name",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeUnspecified,
								true,
							),
						),
						eventFromEventPusherWithInstanceID("instance1",
							instance.NewIDPOIDCConfigAddedEvent(context.Background(),
								&instance.NewAggregate("instance1").Aggregate,
								"clientID",
								"config1",
								"issuer",
								"authEndpoint",
								"tokenEndpoint",
								nil,
								domain.OIDCMappingFieldUnspecified,
								domain.OIDCMappingFieldUnspecified,
							),
						),
					),
					expectPush(
						user.NewUserIDPLinkAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"config1",
							"name",
							"externaluser1",
						),
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
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}

func TestCommandSide_RemoveUserIDPLink(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		checkPermission domain.PermissionCheck
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
				checkPermission: newMockPermissionCheckAllowed(),
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
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "aggregate id missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				link: &domain.UserIDPLink{
					IDPConfigID:    "config1",
					ExternalUserID: "externaluser1",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
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
				checkPermission: newMockPermissionCheckAllowed(),
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
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "external idp not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
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
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "remove external idp, permission error",
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
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
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
				err: zerrors.IsPermissionDenied,
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
						user.NewUserIDPLinkRemovedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"config1",
							"externaluser1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
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
				eventstore:      tt.fields.eventstore,
				checkPermission: tt.fields.checkPermission,
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsPreconditionFailed,
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
						user.NewUserIDPCheckSucceededEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:                  "request1",
								UserAgentID:         "useragent1",
								SelectedIDPConfigID: "config1",
							},
						),
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

func TestCommandSide_MigrateUserIDP(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx         context.Context
		userID      string
		orgID       string
		idpConfigID string
		previousID  string
		newID       string
	}
	type res struct {
		err error
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:         context.Background(),
				userID:      "",
				orgID:       "org1",
				idpConfigID: "idpConfig1",
				previousID:  "previousID",
				newID:       "newID",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-Sn3l1", "Errors.IDMissing"),
			},
		},
		{
			name: "idp link not active, precondition failed error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"idpConfig1",
								"displayName",
								"externalUserID",
							),
						),
					),
				),
			},
			args: args{
				ctx:         context.Background(),
				userID:      "user1",
				orgID:       "org1",
				idpConfigID: "idpConfig1",
				previousID:  "previousID",
				newID:       "newID",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-KJH2o", "Errors.User.ExternalIDP.NotFound"),
			},
		},
		{
			name: "external login check, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"idpConfig1",
								"displayName",
								"previousID",
							),
						),
					),
					expectPush(
						user.NewUserIDPExternalIDMigratedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"idpConfig1",
							"previousID",
							"newID",
						),
					),
				),
			},
			args: args{
				ctx:         context.Background(),
				userID:      "user1",
				orgID:       "org1",
				idpConfigID: "idpConfig1",
				previousID:  "previousID",
				newID:       "newID",
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			err := r.MigrateUserIDP(tt.args.ctx, tt.args.userID, tt.args.orgID, tt.args.idpConfigID, tt.args.previousID, tt.args.newID)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}

func TestCommandSide_UpdateUserIDPLinkUsername(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx            context.Context
		userID         string
		orgID          string
		idpConfigID    string
		externalUserID string
		newUsername    string
	}
	type res struct {
		err error
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:            context.Background(),
				userID:         "",
				orgID:          "org1",
				idpConfigID:    "idpConfig1",
				externalUserID: "externalUserID",
				newUsername:    "newUsername",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-SFegz", "Errors.IDMissing"),
			},
		},
		{
			name: "idp link not active, precondition failed error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"idpConfig1",
								"displayName",
								"externalUserID",
							),
						),
					),
				),
			},
			args: args{
				ctx:            context.Background(),
				userID:         "user1",
				orgID:          "org1",
				idpConfigID:    "idpConfig1",
				externalUserID: "otherID",
				newUsername:    "newUsername",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-DGhre", "Errors.User.ExternalIDP.NotFound"),
			},
		},
		{
			name: "external username not changed, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"idpConfig1",
								"displayName",
								"externalUserID",
							),
						),
					),
				),
			},
			args: args{
				ctx:            context.Background(),
				userID:         "user1",
				orgID:          "org1",
				idpConfigID:    "idpConfig1",
				externalUserID: "externalUserID",
				newUsername:    "displayName",
			},
			res: res{},
		},
		{
			name: "external username update, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"idpConfig1",
								"displayName",
								"externalUserID",
							),
						),
					),
					expectPush(
						user.NewUserIDPExternalUsernameEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"idpConfig1",
							"externalUserID",
							"newUsername",
						),
					),
				),
			},
			args: args{
				ctx:            context.Background(),
				userID:         "user1",
				orgID:          "org1",
				idpConfigID:    "idpConfig1",
				externalUserID: "externalUserID",
				newUsername:    "newUsername",
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			err := r.UpdateUserIDPLinkUsername(tt.args.ctx, tt.args.userID, tt.args.orgID, tt.args.idpConfigID, tt.args.externalUserID, tt.args.newUsername)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}
