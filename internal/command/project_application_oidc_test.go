package command

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestAddOIDCApp(t *testing.T) {
	type fields struct {
		idGenerator id.Generator
	}
	type args struct {
		app    *addOIDCApp
		filter preparation.FilterToQueryReducer
	}

	ctx := context.Background()
	agg := project.NewAggregate("test", "test")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   Want
	}{
		{
			name:   "invalid appID",
			fields: fields{},
			args: args{
				app: &addOIDCApp{
					AddApp: AddApp{
						Aggregate: *agg,
						ID:        "",
						Name:      "name",
					},
					GrantTypes:      []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ResponseTypes:   []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					Version:         domain.OIDCVersionV1,
					ApplicationType: domain.OIDCApplicationTypeWeb,
					AuthMethodType:  domain.OIDCAuthMethodTypeNone,
					AccessTokenType: domain.OIDCTokenTypeBearer,
				},
			},
			want: Want{
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "PROJE-NnavI", "Errors.Invalid.Argument"),
			},
		},
		{
			name:   "invalid name",
			fields: fields{},
			args: args{
				app: &addOIDCApp{
					AddApp: AddApp{
						Aggregate: *agg,
						ID:        "id",
						Name:      "",
					},
					GrantTypes:      []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ResponseTypes:   []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					Version:         domain.OIDCVersionV1,
					ApplicationType: domain.OIDCApplicationTypeWeb,
					AuthMethodType:  domain.OIDCAuthMethodTypeNone,
					AccessTokenType: domain.OIDCTokenTypeBearer,
				},
			},
			want: Want{
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "PROJE-Fef31", "Errors.Invalid.Argument"),
			},
		},
		{
			name:   "project doesn't exist",
			fields: fields{},
			args: args{
				app: &addOIDCApp{
					AddApp: AddApp{
						Aggregate: *agg,
						ID:        "id",
						Name:      "name",
					},
					GrantTypes:      []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ResponseTypes:   []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					Version:         domain.OIDCVersionV1,
					ApplicationType: domain.OIDCApplicationTypeWeb,
					AuthMethodType:  domain.OIDCAuthMethodTypeNone,
					AccessTokenType: domain.OIDCTokenTypeBearer,
				},
				filter: NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).
					Filter(),
			},
			want: Want{
				CreateErr: zerrors.ThrowNotFound(nil, "PROJE-6swVG", ""),
			},
		},
		{
			name: "correct, using uris with whitespaces",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "clientID"),
			},
			args: args{
				app: &addOIDCApp{
					AddApp: AddApp{
						Aggregate: *agg,
						ID:        "id",
						Name:      "name",
					},
					RedirectUris:           []string{" https://test.ch "},
					PostLogoutRedirectUris: []string{" https://test.ch/logout "},
					GrantTypes:             []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ResponseTypes:          []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					Version:                domain.OIDCVersionV1,
					AdditionalOrigins:      []string{" https://sub.test.ch "},
					ApplicationType:        domain.OIDCApplicationTypeWeb,
					AuthMethodType:         domain.OIDCAuthMethodTypeNone,
					AccessTokenType:        domain.OIDCTokenTypeBearer,
				},
				filter: NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							project.NewProjectAddedEvent(
								ctx,
								&agg.Aggregate,
								"project",
								false,
								false,
								false,
								domain.PrivateLabelingSettingUnspecified,
							),
						}, nil
					}).
					Filter(),
			},
			want: Want{
				Commands: []eventstore.Command{
					project.NewApplicationAddedEvent(ctx, &agg.Aggregate,
						"id",
						"name",
					),
					project.NewOIDCConfigAddedEvent(ctx, &agg.Aggregate,
						domain.OIDCVersionV1,
						"id",
						"clientID",
						"",
						[]string{"https://test.ch"},
						[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
						[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
						domain.OIDCApplicationTypeWeb,
						domain.OIDCAuthMethodTypeNone,
						[]string{"https://test.ch/logout"},
						false,
						domain.OIDCTokenTypeBearer,
						false,
						false,
						false,
						0,
						[]string{"https://sub.test.ch"},
						false,
						"",
						domain.LoginVersionUnspecified,
						"",
					),
				},
			},
		},
		{
			name: "correct",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "clientID"),
			},
			args: args{
				app: &addOIDCApp{
					AddApp: AddApp{
						Aggregate: *agg,
						ID:        "id",
						Name:      "name",
					},
					GrantTypes:    []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ResponseTypes: []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					Version:       domain.OIDCVersionV1,

					ApplicationType: domain.OIDCApplicationTypeWeb,
					AuthMethodType:  domain.OIDCAuthMethodTypeNone,
					AccessTokenType: domain.OIDCTokenTypeBearer,
				},
				filter: NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							project.NewProjectAddedEvent(
								ctx,
								&agg.Aggregate,
								"project",
								false,
								false,
								false,
								domain.PrivateLabelingSettingUnspecified,
							),
						}, nil
					}).
					Filter(),
			},
			want: Want{
				Commands: []eventstore.Command{
					project.NewApplicationAddedEvent(ctx, &agg.Aggregate,
						"id",
						"name",
					),
					project.NewOIDCConfigAddedEvent(ctx, &agg.Aggregate,
						domain.OIDCVersionV1,
						"id",
						"clientID",
						"",
						nil,
						[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
						[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
						domain.OIDCApplicationTypeWeb,
						domain.OIDCAuthMethodTypeNone,
						nil,
						false,
						domain.OIDCTokenTypeBearer,
						false,
						false,
						false,
						0,
						nil,
						false,
						"",
						domain.LoginVersionUnspecified,
						"",
					),
				},
			},
		},
		{
			name: "correct with old ID format",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "clientID@project"),
			},
			args: args{
				app: &addOIDCApp{
					AddApp: AddApp{
						Aggregate: *agg,
						ID:        "id",
						Name:      "name",
					},
					GrantTypes:    []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ResponseTypes: []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					Version:       domain.OIDCVersionV1,

					ApplicationType: domain.OIDCApplicationTypeWeb,
					AuthMethodType:  domain.OIDCAuthMethodTypeNone,
					AccessTokenType: domain.OIDCTokenTypeBearer,
				},
				filter: NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							project.NewProjectAddedEvent(
								ctx,
								&agg.Aggregate,
								"project",
								false,
								false,
								false,
								domain.PrivateLabelingSettingUnspecified,
							),
						}, nil
					}).
					Filter(),
			},
			want: Want{
				Commands: []eventstore.Command{
					project.NewApplicationAddedEvent(ctx, &agg.Aggregate,
						"id",
						"name",
					),
					project.NewOIDCConfigAddedEvent(ctx, &agg.Aggregate,
						domain.OIDCVersionV1,
						"id",
						"clientID@project",
						"",
						nil,
						[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
						[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
						domain.OIDCApplicationTypeWeb,
						domain.OIDCAuthMethodTypeNone,
						nil,
						false,
						domain.OIDCTokenTypeBearer,
						false,
						false,
						false,
						0,
						nil,
						false,
						"",
						domain.LoginVersionUnspecified,
						"",
					),
				},
			},
		},
		{
			name: "with secret",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "clientID"),
			},
			args: args{
				app: &addOIDCApp{
					AddApp: AddApp{
						Aggregate: *agg,
						ID:        "id",
						Name:      "name",
					},
					GrantTypes:    []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ResponseTypes: []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					Version:       domain.OIDCVersionV1,

					ApplicationType: domain.OIDCApplicationTypeWeb,
					AuthMethodType:  domain.OIDCAuthMethodTypeBasic,
					AccessTokenType: domain.OIDCTokenTypeBearer,
				},
				filter: NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							project.NewProjectAddedEvent(
								ctx,
								&agg.Aggregate,
								"project",
								false,
								false,
								false,
								domain.PrivateLabelingSettingUnspecified,
							),
						}, nil
					}).
					Filter(),
			},
			want: Want{
				Commands: []eventstore.Command{
					project.NewApplicationAddedEvent(ctx, &agg.Aggregate,
						"id",
						"name",
					),
					project.NewOIDCConfigAddedEvent(ctx, &agg.Aggregate,
						domain.OIDCVersionV1,
						"id",
						"clientID",
						"secret",
						nil,
						[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
						[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
						domain.OIDCApplicationTypeWeb,
						domain.OIDCAuthMethodTypeBasic,
						nil,
						false,
						domain.OIDCTokenTypeBearer,
						false,
						false,
						false,
						0,
						nil,
						false,
						"",
						domain.LoginVersionUnspecified,
						"",
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Commands{
				idGenerator:     tt.fields.idGenerator,
				newHashedSecret: mockHashedSecret("secret"),
				defaultSecretGenerators: &SecretGenerators{
					ClientSecret: emptyConfig,
				},
			}
			AssertValidation(t,
				context.Background(),
				c.AddOIDCAppCommand(
					tt.args.app,
				), tt.args.filter, tt.want)
		})
	}
}

func TestCommandSide_AddOIDCApplication(t *testing.T) {
	t.Parallel()

	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx           context.Context
		oidcApp       *domain.OIDCApp
		resourceOwner string
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
			name: "no aggregate id, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				oidcApp:       &domain.OIDCApp{},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:   "app1",
					AppName: "app",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid app, invalid argument error",
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
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:   "app1",
					AppName: "",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "create oidc app basic using whitespaces in uris, ok",
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
							"app",
						),
						project.NewOIDCConfigAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							domain.OIDCVersionV1,
							"app1",
							"client1",
							"secret",
							[]string{"https://test.ch"},
							[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
							[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
							domain.OIDCApplicationTypeWeb,
							domain.OIDCAuthMethodTypePost,
							[]string{"https://test.ch/logout"},
							true,
							domain.OIDCTokenTypeBearer,
							true,
							true,
							true,
							time.Second*1,
							[]string{"https://sub.test.ch"},
							true,
							"https://test.ch/backchannel",
							domain.LoginVersion2,
							"https://login.test.ch",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "app1", "client1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppName:                  "app",
					AuthMethodType:           gu.Ptr(domain.OIDCAuthMethodTypePost),
					OIDCVersion:              gu.Ptr(domain.OIDCVersionV1),
					RedirectUris:             []string{" https://test.ch "},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          gu.Ptr(domain.OIDCApplicationTypeWeb),
					PostLogoutRedirectUris:   []string{" https://test.ch/logout "},
					DevMode:                  gu.Ptr(true),
					AccessTokenType:          gu.Ptr(domain.OIDCTokenTypeBearer),
					AccessTokenRoleAssertion: gu.Ptr(true),
					IDTokenRoleAssertion:     gu.Ptr(true),
					IDTokenUserinfoAssertion: gu.Ptr(true),
					ClockSkew:                gu.Ptr(time.Second * 1),
					AdditionalOrigins:        []string{" https://sub.test.ch "},
					SkipNativeAppSuccessPage: gu.Ptr(true),
					BackChannelLogoutURI:     gu.Ptr(" https://test.ch/backchannel "),
					LoginVersion:             gu.Ptr(domain.LoginVersion2),
					LoginBaseURI:             gu.Ptr(" https://login.test.ch "),
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:                    "app1",
					AppName:                  "app",
					ClientID:                 "client1",
					ClientSecretString:       "secret",
					AuthMethodType:           gu.Ptr(domain.OIDCAuthMethodTypePost),
					OIDCVersion:              gu.Ptr(domain.OIDCVersionV1),
					RedirectUris:             []string{"https://test.ch"},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          gu.Ptr(domain.OIDCApplicationTypeWeb),
					PostLogoutRedirectUris:   []string{"https://test.ch/logout"},
					DevMode:                  gu.Ptr(true),
					AccessTokenType:          gu.Ptr(domain.OIDCTokenTypeBearer),
					AccessTokenRoleAssertion: gu.Ptr(true),
					IDTokenRoleAssertion:     gu.Ptr(true),
					IDTokenUserinfoAssertion: gu.Ptr(true),
					ClockSkew:                gu.Ptr(time.Second * 1),
					AdditionalOrigins:        []string{"https://sub.test.ch"},
					SkipNativeAppSuccessPage: gu.Ptr(true),
					BackChannelLogoutURI:     gu.Ptr("https://test.ch/backchannel"),
					LoginVersion:             gu.Ptr(domain.LoginVersion2),
					LoginBaseURI:             gu.Ptr("https://login.test.ch"),
					State:                    domain.AppStateActive,
					Compliance:               &domain.Compliance{},
				},
			},
		},
		{
			name: "create oidc app basic, ok",
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
							"app",
						),
						project.NewOIDCConfigAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							domain.OIDCVersionV1,
							"app1",
							"client1",
							"secret",
							[]string{"https://test.ch"},
							[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
							[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
							domain.OIDCApplicationTypeWeb,
							domain.OIDCAuthMethodTypePost,
							[]string{"https://test.ch/logout"},
							true,
							domain.OIDCTokenTypeBearer,
							true,
							true,
							true,
							time.Second*1,
							[]string{"https://sub.test.ch"},
							true,
							"https://test.ch/backchannel",
							domain.LoginVersion2,
							"https://login.test.ch",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "app1", "client1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppName:                  "app",
					AuthMethodType:           gu.Ptr(domain.OIDCAuthMethodTypePost),
					OIDCVersion:              gu.Ptr(domain.OIDCVersionV1),
					RedirectUris:             []string{"https://test.ch"},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          gu.Ptr(domain.OIDCApplicationTypeWeb),
					PostLogoutRedirectUris:   []string{"https://test.ch/logout"},
					DevMode:                  gu.Ptr(true),
					AccessTokenType:          gu.Ptr(domain.OIDCTokenTypeBearer),
					AccessTokenRoleAssertion: gu.Ptr(true),
					IDTokenRoleAssertion:     gu.Ptr(true),
					IDTokenUserinfoAssertion: gu.Ptr(true),
					ClockSkew:                gu.Ptr(time.Second * 1),
					AdditionalOrigins:        []string{"https://sub.test.ch"},
					SkipNativeAppSuccessPage: gu.Ptr(true),
					BackChannelLogoutURI:     gu.Ptr("https://test.ch/backchannel"),
					LoginVersion:             gu.Ptr(domain.LoginVersion2),
					LoginBaseURI:             gu.Ptr("https://login.test.ch"),
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:                    "app1",
					AppName:                  "app",
					ClientID:                 "client1",
					ClientSecretString:       "secret",
					AuthMethodType:           gu.Ptr(domain.OIDCAuthMethodTypePost),
					OIDCVersion:              gu.Ptr(domain.OIDCVersionV1),
					RedirectUris:             []string{"https://test.ch"},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          gu.Ptr(domain.OIDCApplicationTypeWeb),
					PostLogoutRedirectUris:   []string{"https://test.ch/logout"},
					DevMode:                  gu.Ptr(true),
					AccessTokenType:          gu.Ptr(domain.OIDCTokenTypeBearer),
					AccessTokenRoleAssertion: gu.Ptr(true),
					IDTokenRoleAssertion:     gu.Ptr(true),
					IDTokenUserinfoAssertion: gu.Ptr(true),
					ClockSkew:                gu.Ptr(time.Second * 1),
					AdditionalOrigins:        []string{"https://sub.test.ch"},
					SkipNativeAppSuccessPage: gu.Ptr(true),
					BackChannelLogoutURI:     gu.Ptr("https://test.ch/backchannel"),
					LoginVersion:             gu.Ptr(domain.LoginVersion2),
					LoginBaseURI:             gu.Ptr("https://login.test.ch"),
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
				checkPermission: newMockPermissionCheckAllowed(),
			}
			c.setMilestonesCompletedForTest("instanceID")
			got, err := c.AddOIDCApplication(tt.args.ctx, tt.args.oidcApp, tt.args.resourceOwner)
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

func TestCommandSide_ChangeOIDCApplication(t *testing.T) {
	t.Parallel()
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		oidcApp       *domain.OIDCApp
		resourceOwner string
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
			name: "invalid app, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID: "app1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing appid, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:          "",
					AuthMethodType: gu.Ptr(domain.OIDCAuthMethodTypePost),
					GrantTypes:     []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ResponseTypes:  []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing aggregateid, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "",
					},
					AppID:          "appid",
					AuthMethodType: gu.Ptr(domain.OIDCAuthMethodTypePost),
					GrantTypes:     []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ResponseTypes:  []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "app not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:          "app1",
					AuthMethodType: gu.Ptr(domain.OIDCAuthMethodTypePost),
					GrantTypes:     []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ResponseTypes:  []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
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
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
							),
						),
						eventFromEventPusher(
							project.NewOIDCConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								domain.OIDCVersionV1,
								"app1",
								"client1@project",
								"secret",
								[]string{"https://test.ch"},
								[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
								[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
								domain.OIDCApplicationTypeWeb,
								domain.OIDCAuthMethodTypePost,
								[]string{"https://test.ch/logout"},
								true,
								domain.OIDCTokenTypeBearer,
								true,
								true,
								true,
								time.Second*1,
								[]string{"https://sub.test.ch"},
								true,
								"https://test.ch/backchannel",
								domain.LoginVersion2,
								"https://login.test.ch",
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:                    "app1",
					AppName:                  "app",
					AuthMethodType:           gu.Ptr(domain.OIDCAuthMethodTypePost),
					OIDCVersion:              gu.Ptr(domain.OIDCVersionV1),
					RedirectUris:             []string{"https://test.ch"},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          gu.Ptr(domain.OIDCApplicationTypeWeb),
					PostLogoutRedirectUris:   []string{"https://test.ch/logout"},
					DevMode:                  gu.Ptr(true),
					AccessTokenType:          gu.Ptr(domain.OIDCTokenTypeBearer),
					AccessTokenRoleAssertion: gu.Ptr(true),
					IDTokenRoleAssertion:     gu.Ptr(true),
					IDTokenUserinfoAssertion: gu.Ptr(true),
					ClockSkew:                gu.Ptr(time.Second * 1),
					AdditionalOrigins:        []string{"https://sub.test.ch"},
					SkipNativeAppSuccessPage: gu.Ptr(true),
					BackChannelLogoutURI:     gu.Ptr("https://test.ch/backchannel"),
					LoginVersion:             gu.Ptr(domain.LoginVersion2),
					LoginBaseURI:             gu.Ptr("https://login.test.ch"),
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "no changes whitespaces are ignored, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
							),
						),
						eventFromEventPusher(
							project.NewOIDCConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								domain.OIDCVersionV1,
								"app1",
								"client1@project",
								"secret",
								[]string{"https://test.ch"},
								[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
								[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
								domain.OIDCApplicationTypeWeb,
								domain.OIDCAuthMethodTypePost,
								[]string{"https://test.ch/logout"},
								true,
								domain.OIDCTokenTypeBearer,
								true,
								true,
								true,
								time.Second*1,
								[]string{"https://sub.test.ch"},
								true,
								"https://test.ch/backchannel",
								domain.LoginVersion2,
								"https://login.test.ch",
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:                    "app1",
					AppName:                  "app",
					AuthMethodType:           gu.Ptr(domain.OIDCAuthMethodTypePost),
					OIDCVersion:              gu.Ptr(domain.OIDCVersionV1),
					RedirectUris:             []string{"https://test.ch "},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          gu.Ptr(domain.OIDCApplicationTypeWeb),
					PostLogoutRedirectUris:   []string{" https://test.ch/logout"},
					DevMode:                  gu.Ptr(true),
					AccessTokenType:          gu.Ptr(domain.OIDCTokenTypeBearer),
					AccessTokenRoleAssertion: gu.Ptr(true),
					IDTokenRoleAssertion:     gu.Ptr(true),
					IDTokenUserinfoAssertion: gu.Ptr(true),
					ClockSkew:                gu.Ptr(time.Second * 1),
					AdditionalOrigins:        []string{" https://sub.test.ch "},
					SkipNativeAppSuccessPage: gu.Ptr(true),
					BackChannelLogoutURI:     gu.Ptr(" https://test.ch/backchannel "),
					LoginVersion:             gu.Ptr(domain.LoginVersion2),
					LoginBaseURI:             gu.Ptr(" https://login.test.ch "),
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "partial change oidc app, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
							),
						),
						eventFromEventPusher(
							project.NewOIDCConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								domain.OIDCVersionV1,
								"app1",
								"client1@project",
								"secret",
								[]string{"https://test.ch"},
								[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
								[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
								domain.OIDCApplicationTypeWeb,
								domain.OIDCAuthMethodTypePost,
								[]string{"https://test.ch/logout"},
								false,
								domain.OIDCTokenTypeBearer,
								true,
								true,
								true,
								time.Second*1,
								[]string{"https://sub.test.ch"},
								true,
								"https://test.ch/backchannel",
								domain.LoginVersion1,
								"",
							),
						),
					),
					expectFilter(),
					expectPush(
						newOIDCAppChangedEvent(context.Background(),
							"app1",
							"project1",
							"org1"),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:          "app1",
					AppName:        "app",
					AuthMethodType: gu.Ptr(domain.OIDCAuthMethodTypeBasic),
					GrantTypes:     []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ResponseTypes:  []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:                    "app1",
					ClientID:                 "client1@project",
					AppName:                  "app",
					AuthMethodType:           gu.Ptr(domain.OIDCAuthMethodTypeBasic),
					OIDCVersion:              gu.Ptr(domain.OIDCVersionV1),
					RedirectUris:             []string{"https://test.ch"},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          gu.Ptr(domain.OIDCApplicationTypeWeb),
					PostLogoutRedirectUris:   []string{"https://test.ch/logout"},
					DevMode:                  gu.Ptr(false),
					AccessTokenType:          gu.Ptr(domain.OIDCTokenTypeBearer),
					AccessTokenRoleAssertion: gu.Ptr(true),
					IDTokenRoleAssertion:     gu.Ptr(true),
					IDTokenUserinfoAssertion: gu.Ptr(true),
					ClockSkew:                gu.Ptr(time.Second * 1),
					AdditionalOrigins:        []string{"https://sub.test.ch"},
					SkipNativeAppSuccessPage: gu.Ptr(true),
					BackChannelLogoutURI:     gu.Ptr("https://test.ch/backchannel"),
					LoginVersion:             gu.Ptr(domain.LoginVersion1),
					LoginBaseURI:             gu.Ptr(""),
					Compliance:               &domain.Compliance{},
					State:                    domain.AppStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: newMockPermissionCheckAllowed(),
			}
			got, err := r.UpdateOIDCApplication(tt.args.ctx, tt.args.oidcApp, tt.args.resourceOwner)
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

func TestCommandSide_ChangeOIDCApplicationSecret(t *testing.T) {
	t.Parallel()

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		appID         string
		projectID     string
		resourceOwner string
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
			name: "no projectid, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "no appid, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "app not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "change secret, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
							),
						),
						eventFromEventPusher(
							project.NewOIDCConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								domain.OIDCVersionV1,
								"app1",
								"client1@project",
								"secret",
								[]string{"https://test.ch"},
								[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
								[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
								domain.OIDCApplicationTypeWeb,
								domain.OIDCAuthMethodTypePost,
								[]string{"https://test.ch/logout"},
								true,
								domain.OIDCTokenTypeBearer,
								true,
								true,
								true,
								time.Second*1,
								[]string{"https://sub.test.ch"},
								false,
								"",
								domain.LoginVersionUnspecified,
								"",
							),
						),
					),
					expectPush(
						project.NewOIDCConfigSecretChangedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"secret",
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:                    "app1",
					AppName:                  "app",
					ClientID:                 "client1@project",
					ClientSecretString:       "secret",
					AuthMethodType:           gu.Ptr(domain.OIDCAuthMethodTypePost),
					OIDCVersion:              gu.Ptr(domain.OIDCVersionV1),
					RedirectUris:             []string{"https://test.ch"},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          gu.Ptr(domain.OIDCApplicationTypeWeb),
					PostLogoutRedirectUris:   []string{"https://test.ch/logout"},
					DevMode:                  gu.Ptr(true),
					AccessTokenType:          gu.Ptr(domain.OIDCTokenTypeBearer),
					AccessTokenRoleAssertion: gu.Ptr(true),
					IDTokenRoleAssertion:     gu.Ptr(true),
					IDTokenUserinfoAssertion: gu.Ptr(true),
					ClockSkew:                gu.Ptr(time.Second * 1),
					AdditionalOrigins:        []string{"https://sub.test.ch"},
					SkipNativeAppSuccessPage: gu.Ptr(false),
					BackChannelLogoutURI:     gu.Ptr(""),
					LoginVersion:             gu.Ptr(domain.LoginVersionUnspecified),
					LoginBaseURI:             gu.Ptr(""),
					State:                    domain.AppStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				newHashedSecret: mockHashedSecret("secret"),
				defaultSecretGenerators: &SecretGenerators{
					ClientSecret: emptyConfig,
				},
				checkPermission: newMockPermissionCheckAllowed(),
			}
			got, err := r.ChangeOIDCApplicationSecret(tt.args.ctx, tt.args.projectID, tt.args.appID, tt.args.resourceOwner)
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

func newOIDCAppChangedEvent(ctx context.Context, appID, projectID, resourceOwner string) *project.OIDCConfigChangedEvent {
	changes := []project.OIDCConfigChanges{
		project.ChangeAuthMethodType(domain.OIDCAuthMethodTypeBasic),
	}
	event, _ := project.NewOIDCConfigChangedEvent(ctx,
		&project.NewAggregate(projectID, resourceOwner).Aggregate,
		appID,
		changes,
	)
	return event
}
