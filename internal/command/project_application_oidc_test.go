package command

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/passwap"
	"github.com/zitadel/passwap/bcrypt"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestAddOIDCApp(t *testing.T) {
	type fields struct {
		idGenerator id_generator.Generator
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
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Commands{
				newHashedSecret: mockHashedSecret("secret"),
				defaultSecretGenerators: &SecretGenerators{
					ClientSecret: emptyConfig,
				},
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			AssertValidation(t,
				context.Background(),
				c.AddOIDCAppCommand(
					tt.args.app,
				), tt.args.filter, tt.want)
		})
	}
}

func TestCommandSide_AddOIDCApplication(t *testing.T) {
	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id_generator.Generator
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
				ctx:           context.Background(),
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
				ctx: context.Background(),
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
				ctx: context.Background(),
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
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "app1", "client1"),
			},
			args: args{
				ctx: context.Background(),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppName:                  "app",
					AuthMethodType:           domain.OIDCAuthMethodTypePost,
					OIDCVersion:              domain.OIDCVersionV1,
					RedirectUris:             []string{" https://test.ch "},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          domain.OIDCApplicationTypeWeb,
					PostLogoutRedirectUris:   []string{" https://test.ch/logout "},
					DevMode:                  true,
					AccessTokenType:          domain.OIDCTokenTypeBearer,
					AccessTokenRoleAssertion: true,
					IDTokenRoleAssertion:     true,
					IDTokenUserinfoAssertion: true,
					ClockSkew:                time.Second * 1,
					AdditionalOrigins:        []string{" https://sub.test.ch "},
					SkipNativeAppSuccessPage: true,
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
					AuthMethodType:           domain.OIDCAuthMethodTypePost,
					OIDCVersion:              domain.OIDCVersionV1,
					RedirectUris:             []string{"https://test.ch"},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          domain.OIDCApplicationTypeWeb,
					PostLogoutRedirectUris:   []string{"https://test.ch/logout"},
					DevMode:                  true,
					AccessTokenType:          domain.OIDCTokenTypeBearer,
					AccessTokenRoleAssertion: true,
					IDTokenRoleAssertion:     true,
					IDTokenUserinfoAssertion: true,
					ClockSkew:                time.Second * 1,
					AdditionalOrigins:        []string{"https://sub.test.ch"},
					SkipNativeAppSuccessPage: true,
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
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "app1", "client1"),
			},
			args: args{
				ctx: context.Background(),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppName:                  "app",
					AuthMethodType:           domain.OIDCAuthMethodTypePost,
					OIDCVersion:              domain.OIDCVersionV1,
					RedirectUris:             []string{"https://test.ch"},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          domain.OIDCApplicationTypeWeb,
					PostLogoutRedirectUris:   []string{"https://test.ch/logout"},
					DevMode:                  true,
					AccessTokenType:          domain.OIDCTokenTypeBearer,
					AccessTokenRoleAssertion: true,
					IDTokenRoleAssertion:     true,
					IDTokenUserinfoAssertion: true,
					ClockSkew:                time.Second * 1,
					AdditionalOrigins:        []string{"https://sub.test.ch"},
					SkipNativeAppSuccessPage: true,
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
					AuthMethodType:           domain.OIDCAuthMethodTypePost,
					OIDCVersion:              domain.OIDCVersionV1,
					RedirectUris:             []string{"https://test.ch"},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          domain.OIDCApplicationTypeWeb,
					PostLogoutRedirectUris:   []string{"https://test.ch/logout"},
					DevMode:                  true,
					AccessTokenType:          domain.OIDCTokenTypeBearer,
					AccessTokenRoleAssertion: true,
					IDTokenRoleAssertion:     true,
					IDTokenUserinfoAssertion: true,
					ClockSkew:                time.Second * 1,
					AdditionalOrigins:        []string{"https://sub.test.ch"},
					SkipNativeAppSuccessPage: true,
					State:                    domain.AppStateActive,
					Compliance:               &domain.Compliance{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				newHashedSecret: mockHashedSecret("secret"),
				defaultSecretGenerators: &SecretGenerators{
					ClientSecret: emptyConfig,
				},
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := r.AddOIDCApplication(tt.args.ctx, tt.args.oidcApp, tt.args.resourceOwner)
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
	type fields struct {
		eventstore *eventstore.Eventstore
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
				eventstore: eventstoreExpect(
					t,
				),
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:          "",
					AuthMethodType: domain.OIDCAuthMethodTypePost,
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				oidcApp: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "",
					},
					AppID:          "appid",
					AuthMethodType: domain.OIDCAuthMethodTypePost,
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
				eventstore: eventstoreExpect(
					t,
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
					AuthMethodType: domain.OIDCAuthMethodTypePost,
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
				eventstore: eventstoreExpect(
					t,
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
							),
						),
					),
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
					AuthMethodType:           domain.OIDCAuthMethodTypePost,
					OIDCVersion:              domain.OIDCVersionV1,
					RedirectUris:             []string{"https://test.ch"},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          domain.OIDCApplicationTypeWeb,
					PostLogoutRedirectUris:   []string{"https://test.ch/logout"},
					DevMode:                  true,
					AccessTokenType:          domain.OIDCTokenTypeBearer,
					AccessTokenRoleAssertion: true,
					IDTokenRoleAssertion:     true,
					IDTokenUserinfoAssertion: true,
					ClockSkew:                time.Second * 1,
					AdditionalOrigins:        []string{"https://sub.test.ch"},
					SkipNativeAppSuccessPage: true,
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
				eventstore: eventstoreExpect(
					t,
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
							),
						),
					),
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
					AuthMethodType:           domain.OIDCAuthMethodTypePost,
					OIDCVersion:              domain.OIDCVersionV1,
					RedirectUris:             []string{"https://test.ch "},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          domain.OIDCApplicationTypeWeb,
					PostLogoutRedirectUris:   []string{" https://test.ch/logout"},
					DevMode:                  true,
					AccessTokenType:          domain.OIDCTokenTypeBearer,
					AccessTokenRoleAssertion: true,
					IDTokenRoleAssertion:     true,
					IDTokenUserinfoAssertion: true,
					ClockSkew:                time.Second * 1,
					AdditionalOrigins:        []string{" https://sub.test.ch "},
					SkipNativeAppSuccessPage: true,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "change oidc app, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
							),
						),
					),
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
					AppID:                    "app1",
					AppName:                  "app",
					AuthMethodType:           domain.OIDCAuthMethodTypePost,
					OIDCVersion:              domain.OIDCVersionV1,
					RedirectUris:             []string{" https://test-change.ch "},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          domain.OIDCApplicationTypeWeb,
					PostLogoutRedirectUris:   []string{" https://test-change.ch/logout "},
					DevMode:                  true,
					AccessTokenType:          domain.OIDCTokenTypeJWT,
					AccessTokenRoleAssertion: false,
					IDTokenRoleAssertion:     false,
					IDTokenUserinfoAssertion: false,
					ClockSkew:                time.Second * 2,
					AdditionalOrigins:        []string{"https://sub.test.ch"},
					SkipNativeAppSuccessPage: true,
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
					AuthMethodType:           domain.OIDCAuthMethodTypePost,
					OIDCVersion:              domain.OIDCVersionV1,
					RedirectUris:             []string{"https://test-change.ch"},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          domain.OIDCApplicationTypeWeb,
					PostLogoutRedirectUris:   []string{"https://test-change.ch/logout"},
					DevMode:                  true,
					AccessTokenType:          domain.OIDCTokenTypeJWT,
					AccessTokenRoleAssertion: false,
					IDTokenRoleAssertion:     false,
					IDTokenUserinfoAssertion: false,
					ClockSkew:                time.Second * 2,
					AdditionalOrigins:        []string{"https://sub.test.ch"},
					SkipNativeAppSuccessPage: true,
					Compliance:               &domain.Compliance{},
					State:                    domain.AppStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeOIDCApplication(tt.args.ctx, tt.args.oidcApp, tt.args.resourceOwner)
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
					AuthMethodType:           domain.OIDCAuthMethodTypePost,
					OIDCVersion:              domain.OIDCVersionV1,
					RedirectUris:             []string{"https://test.ch"},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
					ApplicationType:          domain.OIDCApplicationTypeWeb,
					PostLogoutRedirectUris:   []string{"https://test.ch/logout"},
					DevMode:                  true,
					AccessTokenType:          domain.OIDCTokenTypeBearer,
					AccessTokenRoleAssertion: true,
					IDTokenRoleAssertion:     true,
					IDTokenUserinfoAssertion: true,
					ClockSkew:                time.Second * 1,
					AdditionalOrigins:        []string{"https://sub.test.ch"},
					SkipNativeAppSuccessPage: false,
					State:                    domain.AppStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				newHashedSecret: mockHashedSecret("secret"),
				defaultSecretGenerators: &SecretGenerators{
					ClientSecret: emptyConfig,
				},
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
		project.ChangeRedirectURIs([]string{"https://test-change.ch"}),
		project.ChangePostLogoutRedirectURIs([]string{"https://test-change.ch/logout"}),
		project.ChangeDevMode(true),
		project.ChangeAccessTokenType(domain.OIDCTokenTypeJWT),
		project.ChangeAccessTokenRoleAssertion(false),
		project.ChangeIDTokenRoleAssertion(false),
		project.ChangeIDTokenUserinfoAssertion(false),
		project.ChangeClockSkew(time.Second * 2),
	}
	event, _ := project.NewOIDCConfigChangedEvent(ctx,
		&project.NewAggregate(projectID, resourceOwner).Aggregate,
		appID,
		changes,
	)
	return event
}

func TestCommands_VerifyOIDCClientSecret(t *testing.T) {
	hasher := &crypto.Hasher{
		Swapper: passwap.NewSwapper(bcrypt.New(bcrypt.MinCost)),
	}
	hashedSecret, err := hasher.Hash("secret")
	require.NoError(t, err)
	agg := project.NewAggregate("projectID", "orgID")

	tests := []struct {
		name       string
		secret     string
		eventstore func(*testing.T) *eventstore.Eventstore
		wantErr    error
	}{
		{
			name: "filter error",
			eventstore: expectEventstore(
				expectFilterError(io.ErrClosedPipe),
			),
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "app not exists",
			eventstore: expectEventstore(
				expectFilter(),
			),
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-D8hba", "Errors.Project.App.NotExisting"),
		},
		{
			name: "wrong app type",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(
						project.NewApplicationAddedEvent(context.Background(), &agg.Aggregate, "appID", "appName"),
					),
				),
			),
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-BHgn2", "Errors.Project.App.IsNotOIDC"),
		},
		{
			name: "no secret set",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(
						project.NewApplicationAddedEvent(context.Background(), &agg.Aggregate, "appID", "appName"),
					),
					eventFromEventPusher(
						project.NewOIDCConfigAddedEvent(context.Background(),
							&agg.Aggregate,
							domain.OIDCVersionV1,
							"appID",
							"client1@project",
							"",
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
						),
					),
				),
			),
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-D6hba", "Errors.Project.App.OIDCConfigInvalid"),
		},
		{
			name:   "check succeeded",
			secret: "secret",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(
						project.NewApplicationAddedEvent(context.Background(), &agg.Aggregate, "appID", "appName"),
					),
					eventFromEventPusher(
						project.NewOIDCConfigAddedEvent(context.Background(),
							&agg.Aggregate,
							domain.OIDCVersionV1,
							"appID",
							"client1@project",
							hashedSecret,
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
						),
					),
				),
				expectPush(
					project.NewOIDCConfigSecretCheckSucceededEvent(context.Background(), &agg.Aggregate, "appID"),
				),
			),
		},
		{
			name:   "check failed",
			secret: "wrong!",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(
						project.NewApplicationAddedEvent(context.Background(), &agg.Aggregate, "appID", "appName"),
					),
					eventFromEventPusher(
						project.NewOIDCConfigAddedEvent(context.Background(),
							&agg.Aggregate,
							domain.OIDCVersionV1,
							"appID",
							"client1@project",
							hashedSecret,
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
						),
					),
				),
				expectPush(
					project.NewOIDCConfigSecretCheckFailedEvent(context.Background(), &agg.Aggregate, "appID"),
				),
			),
			wantErr: zerrors.ThrowInvalidArgument(err, "COMMAND-Bz542", "Errors.Project.App.ClientSecretInvalid"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:   tt.eventstore(t),
				secretHasher: hasher,
			}
			err := c.VerifyOIDCClientSecret(context.Background(), "projectID", "appID", tt.secret)
			c.jobs.Wait()
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
