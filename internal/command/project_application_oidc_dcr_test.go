package command

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_EnsureDCRProject(t *testing.T) {
	t.Parallel()

	dcrProjectAdded := func(projectID string) eventstore.Event {
		return eventFromEventPusher(project.NewProjectAddedEvent(
			context.Background(),
			&project.NewAggregate(projectID, "org1").Aggregate,
			DCRProjectName,
			false,
			false,
			false,
			domain.PrivateLabelingSettingUnspecified,
		))
	}

	t.Run("missing resource owner, invalid argument", func(t *testing.T) {
		t.Parallel()
		c := &Commands{
			eventstore:      expectEventstore()(t),
			checkPermission: newMockPermissionCheckNotAllowed(),
		}
		_, err := c.EnsureDCRProject(authz.WithInstanceID(context.Background(), "instanceID"), "")
		assert.True(t, zerrors.IsErrorInvalidArgument(err))
	})

	t.Run("existing project resolved from the eventstore", func(t *testing.T) {
		t.Parallel()
		c := &Commands{
			eventstore: expectEventstore(
				expectFilter(dcrProjectAdded("existing")),
			)(t),
			checkPermission: newMockPermissionCheckNotAllowed(),
		}
		projectID, err := c.EnsureDCRProject(authz.WithInstanceID(context.Background(), "instanceID"), "org1")
		assert.NoError(t, err)
		assert.Equal(t, "existing", projectID)
	})

	t.Run("project created when none exists, without permission check", func(t *testing.T) {
		t.Parallel()
		c := &Commands{
			eventstore: expectEventstore(
				expectFilter(), // lookup: no DCR project yet
				expectFilter(), // project write model for the new id
				expectPush(
					project.NewProjectAddedEvent(
						context.Background(),
						&project.NewAggregate("project1", "org1").Aggregate,
						DCRProjectName,
						false,
						false,
						false,
						domain.PrivateLabelingSettingUnspecified,
					),
				),
			)(t),
			idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "project1"),
			// Provisioning the DCR project must not depend on project.create.
			checkPermission: newMockPermissionCheckNotAllowed(),
		}
		c.setMilestonesCompletedForTest("instanceID")
		projectID, err := c.EnsureDCRProject(authz.WithInstanceID(context.Background(), "instanceID"), "org1")
		assert.NoError(t, err)
		assert.Equal(t, "project1", projectID)
	})

	t.Run("concurrent creation resolves the winner from the eventstore", func(t *testing.T) {
		t.Parallel()
		c := &Commands{
			eventstore: expectEventstore(
				expectFilter(), // lookup: no DCR project yet
				expectFilter(), // project write model for the new id
				expectPushFailed(
					zerrors.ThrowAlreadyExists(nil, "id", "project name already taken"),
					project.NewProjectAddedEvent(
						context.Background(),
						&project.NewAggregate("project1", "org1").Aggregate,
						DCRProjectName,
						false,
						false,
						false,
						domain.PrivateLabelingSettingUnspecified,
					),
				),
				expectFilter(dcrProjectAdded("winner")), // re-resolve: the racing winner
			)(t),
			idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "project1"),
			checkPermission: newMockPermissionCheckNotAllowed(),
		}
		c.setMilestonesCompletedForTest("instanceID")
		projectID, err := c.EnsureDCRProject(authz.WithInstanceID(context.Background(), "instanceID"), "org1")
		assert.NoError(t, err)
		assert.Equal(t, "winner", projectID)
	})
}

func TestCommandSide_AddDynamicOIDCClient(t *testing.T) {
	t.Parallel()

	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx           context.Context
		projectID     string
		resourceOwner string
		oidcApp       *domain.OIDCApp
	}
	type res struct {
		want *domain.OIDCApp
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "missing project id, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				projectID:     "",
				resourceOwner: "org1",
				oidcApp:       &domain.OIDCApp{AppName: "app"},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid app, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				projectID:     "project1",
				resourceOwner: "org1",
				oidcApp:       &domain.OIDCApp{AppName: ""},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, precondition failed error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				projectID:     "project1",
				resourceOwner: "org1",
				oidcApp: &domain.OIDCApp{
					AppName:        "app",
					AuthMethodType: gu.Ptr(domain.OIDCAuthMethodTypeNone),
					RedirectUris:   []string{"https://client.example.com/callback"},
					ResponseTypes:  []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:     []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "confidential client (client_secret_basic), secret returned, no permission required",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingUnspecified),
						),
					),
					expectFilter(),
					expectPush(
						project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app (app1)",
						),
						project.NewOIDCConfigAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							domain.OIDCVersionV1,
							"app1",
							"client1",
							"secret",
							[]string{"https://client.example.com/callback"},
							[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
							[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
							domain.OIDCApplicationTypeWeb,
							domain.OIDCAuthMethodTypeBasic,
							[]string{},
							false,
							domain.OIDCTokenTypeBearer,
							false,
							false,
							false,
							0,
							[]string{},
							false,
							"",
							domain.LoginVersionUnspecified,
							"",
						),
						project.NewOIDCConfigRegistrationTokenChangedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"secret",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "app1", "client1"),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				projectID:     "project1",
				resourceOwner: "org1",
				oidcApp: &domain.OIDCApp{
					AppName:         "app",
					AuthMethodType:  gu.Ptr(domain.OIDCAuthMethodTypeBasic),
					OIDCVersion:     gu.Ptr(domain.OIDCVersionV1),
					RedirectUris:    []string{"https://client.example.com/callback"},
					ResponseTypes:   []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:      []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType: gu.Ptr(domain.OIDCApplicationTypeWeb),
					AccessTokenType: gu.Ptr(domain.OIDCTokenTypeBearer),
				},
			},
			res: res{
				want: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:                    "app1",
					AppName:                  "app (app1)",
					ClientID:                 "client1",
					ClientSecretString:       "secret",
					RegistrationAccessToken:  "secret",
					AuthMethodType:           gu.Ptr(domain.OIDCAuthMethodTypeBasic),
					OIDCVersion:              gu.Ptr(domain.OIDCVersionV1),
					RedirectUris:             []string{"https://client.example.com/callback"},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          gu.Ptr(domain.OIDCApplicationTypeWeb),
					AccessTokenType:          gu.Ptr(domain.OIDCTokenTypeBearer),
					DevMode:                  gu.Ptr(false),
					AccessTokenRoleAssertion: gu.Ptr(false),
					IDTokenRoleAssertion:     gu.Ptr(false),
					IDTokenUserinfoAssertion: gu.Ptr(false),
					ClockSkew:                gu.Ptr(time.Duration(0)),
					SkipNativeAppSuccessPage: gu.Ptr(false),
					BackChannelLogoutURI:     gu.Ptr(""),
					LoginVersion:             gu.Ptr(domain.LoginVersionUnspecified),
					LoginBaseURI:             gu.Ptr(""),
					State:                    domain.AppStateActive,
					Compliance:               &domain.Compliance{},
				},
			},
		},
		{
			name: "public client (none), no secret returned, no permission required",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingUnspecified),
						),
					),
					expectFilter(),
					expectPush(
						project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"DCR Client app1",
						),
						project.NewOIDCConfigAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							domain.OIDCVersionV1,
							"app1",
							"client1",
							"",
							[]string{"https://client.example.com/callback"},
							[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
							[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
							domain.OIDCApplicationTypeWeb,
							domain.OIDCAuthMethodTypeNone,
							[]string{},
							false,
							domain.OIDCTokenTypeBearer,
							false,
							false,
							false,
							0,
							[]string{},
							false,
							"",
							domain.LoginVersionUnspecified,
							"",
						),
						project.NewOIDCConfigRegistrationTokenChangedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"secret",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "app1", "client1"),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				projectID:     "project1",
				resourceOwner: "org1",
				oidcApp: &domain.OIDCApp{
					AppName:         "",
					AuthMethodType:  gu.Ptr(domain.OIDCAuthMethodTypeNone),
					OIDCVersion:     gu.Ptr(domain.OIDCVersionV1),
					RedirectUris:    []string{"https://client.example.com/callback"},
					ResponseTypes:   []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:      []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType: gu.Ptr(domain.OIDCApplicationTypeWeb),
					AccessTokenType: gu.Ptr(domain.OIDCTokenTypeBearer),
				},
			},
			res: res{
				want: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:                    "app1",
					AppName:                  "DCR Client app1",
					ClientID:                 "client1",
					RegistrationAccessToken:  "secret",
					AuthMethodType:           gu.Ptr(domain.OIDCAuthMethodTypeNone),
					OIDCVersion:              gu.Ptr(domain.OIDCVersionV1),
					RedirectUris:             []string{"https://client.example.com/callback"},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          gu.Ptr(domain.OIDCApplicationTypeWeb),
					AccessTokenType:          gu.Ptr(domain.OIDCTokenTypeBearer),
					DevMode:                  gu.Ptr(false),
					AccessTokenRoleAssertion: gu.Ptr(false),
					IDTokenRoleAssertion:     gu.Ptr(false),
					IDTokenUserinfoAssertion: gu.Ptr(false),
					ClockSkew:                gu.Ptr(time.Duration(0)),
					SkipNativeAppSuccessPage: gu.Ptr(false),
					BackChannelLogoutURI:     gu.Ptr(""),
					LoginVersion:             gu.Ptr(domain.LoginVersionUnspecified),
					LoginBaseURI:             gu.Ptr(""),
					State:                    domain.AppStateActive,
					Compliance:               &domain.Compliance{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				idGenerator:     tt.fields.idGenerator,
				newHashedSecret: mockHashedSecret("secret"),
				defaultSecretGenerators: &SecretGenerators{
					ClientSecret: emptyConfig,
				},
				// Dynamic client registration must NOT depend on app.write: deny all
				// permissions and assert the client is still created.
				checkPermission: newMockPermissionCheckNotAllowed(),
				ipLookupFunction: func(_ string) ([]net.IP, error) {
					return []net.IP{net.ParseIP("1.2.3.4")}, nil
				},
			}
			c.setMilestonesCompletedForTest("instanceID")
			got, err := c.AddDynamicOIDCClient(tt.args.ctx, tt.args.projectID, tt.args.resourceOwner, tt.args.oidcApp)
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

func TestCommandSide_VerifyDynamicClientRegistrationToken(t *testing.T) {
	t.Parallel()

	regToken := func() *project.OIDCConfigRegistrationTokenChangedEvent {
		return project.NewOIDCConfigRegistrationTokenChangedEvent(context.Background(),
			&project.NewAggregate("project1", "org1").Aggregate, "app1", "$plain$$secret")
	}

	tests := []struct {
		name       string
		eventstore func(*testing.T) *eventstore.Eventstore
		secret     string
		wantErr    func(error) bool
	}{
		{
			name:       "valid token",
			eventstore: expectEventstore(expectFilter(eventFromEventPusher(regToken()))),
			secret:     "secret",
			wantErr:    nil,
		},
		{
			name:       "wrong secret is unauthenticated",
			eventstore: expectEventstore(expectFilter(eventFromEventPusher(regToken()))),
			secret:     "wrong",
			wantErr:    zerrors.IsUnauthenticated,
		},
		{
			name:       "no token set is unauthenticated",
			eventstore: expectEventstore(expectFilter()),
			secret:     "secret",
			wantErr:    zerrors.IsUnauthenticated,
		},
		{
			name:       "empty secret is unauthenticated",
			eventstore: expectEventstore(),
			secret:     "",
			wantErr:    zerrors.IsUnauthenticated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Commands{
				eventstore:   tt.eventstore(t),
				secretHasher: mockPasswordHasher(""),
			}
			err := c.VerifyDynamicClientRegistrationToken(context.Background(), "project1", "app1", "org1", tt.secret)
			if tt.wantErr == nil {
				assert.NoError(t, err)
				return
			}
			assert.True(t, tt.wantErr(err), "got wrong err: %v", err)
		})
	}
}

func TestCommandSide_UpdateDynamicOIDCClient(t *testing.T) {
	t.Parallel()

	existingClient := func() []eventstore.Event {
		return []eventstore.Event{
			eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
				&project.NewAggregate("project1", "org1").Aggregate, "app1", "DCR Client app1")),
			eventFromEventPusher(project.NewOIDCConfigAddedEvent(context.Background(),
				&project.NewAggregate("project1", "org1").Aggregate,
				domain.OIDCVersionV1, "app1", "client1", "",
				[]string{"https://client.example.com/callback"},
				[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
				[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
				domain.OIDCApplicationTypeWeb, domain.OIDCAuthMethodTypeNone,
				[]string{}, false, domain.OIDCTokenTypeBearer,
				false, false, false, 0, []string{}, false, "",
				domain.LoginVersionUnspecified, "")),
		}
	}
	sameMetadata := &domain.OIDCApp{
		ObjectRoot:      models.ObjectRoot{AggregateID: "project1"},
		AppID:           "app1",
		RedirectUris:    []string{"https://client.example.com/callback"},
		ResponseTypes:   []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
		GrantTypes:      []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
		ApplicationType: gu.Ptr(domain.OIDCApplicationTypeWeb),
		AuthMethodType:  gu.Ptr(domain.OIDCAuthMethodTypeNone),
		OIDCVersion:     gu.Ptr(domain.OIDCVersionV1),
		AccessTokenType: gu.Ptr(domain.OIDCTokenTypeBearer),
	}

	t.Run("unchanged metadata still rotates the token, no permission required", func(t *testing.T) {
		t.Parallel()
		c := &Commands{
			eventstore: expectEventstore(
				expectFilter(existingClient()...),
				expectFilter(),
				expectPush(
					project.NewOIDCConfigRegistrationTokenChangedEvent(context.Background(),
						&project.NewAggregate("project1", "org1").Aggregate, "app1", "secret"),
				),
			)(t),
			newHashedSecret: mockHashedSecret("secret"),
			checkPermission: newMockPermissionCheckNotAllowed(),
		}
		got, err := c.UpdateDynamicOIDCClient(context.Background(), sameMetadata, "org1")
		assert.NoError(t, err)
		assert.Equal(t, "client1", got.ClientID)
		assert.Equal(t, "secret", got.RegistrationAccessToken)
	})

	t.Run("not existing client is not found", func(t *testing.T) {
		t.Parallel()
		c := &Commands{
			eventstore:      expectEventstore(expectFilter())(t),
			checkPermission: newMockPermissionCheckNotAllowed(),
		}
		_, err := c.UpdateDynamicOIDCClient(context.Background(), sameMetadata, "org1")
		assert.True(t, zerrors.IsNotFound(err), "got wrong err: %v", err)
	})
}

func TestCommandSide_RemoveDynamicOIDCClient(t *testing.T) {
	t.Parallel()

	t.Run("remove existing client, no permission required", func(t *testing.T) {
		t.Parallel()
		c := &Commands{
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
						&project.NewAggregate("project1", "org1").Aggregate, "app1", "DCR Client app1")),
				),
				expectPush(
					project.NewApplicationRemovedEvent(context.Background(),
						&project.NewAggregate("project1", "org1").Aggregate, "app1", "DCR Client app1", ""),
				),
			)(t),
			checkPermission: newMockPermissionCheckNotAllowed(),
		}
		got, err := c.RemoveDynamicOIDCClient(context.Background(), "project1", "app1", "org1")
		assert.NoError(t, err)
		assertObjectDetails(t, &domain.ObjectDetails{ResourceOwner: "org1"}, got)
	})

	t.Run("not existing client is not found", func(t *testing.T) {
		t.Parallel()
		c := &Commands{
			eventstore:      expectEventstore(expectFilter())(t),
			checkPermission: newMockPermissionCheckNotAllowed(),
		}
		_, err := c.RemoveDynamicOIDCClient(context.Background(), "project1", "app1", "org1")
		assert.True(t, zerrors.IsNotFound(err), "got wrong err: %v", err)
	})
}
