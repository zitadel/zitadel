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
