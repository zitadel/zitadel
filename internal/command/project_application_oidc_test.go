package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func TestAddOIDCApp(t *testing.T) {
	type fields struct {
		idGenerator id.Generator
	}
	type args struct {
		app             *addOIDCApp
		clientSecretAlg crypto.HashAlgorithm
		filter          preparation.FilterToQueryReducer
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
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-NnavI", "Errors.Invalid.Argument"),
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
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-Fef31", "Errors.Invalid.Argument"),
			},
		},
		{
			name:   "project not exists",
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
				CreateErr: errors.ThrowNotFound(nil, "PROJE-6swVG", ""),
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
						"clientID@project",
						nil,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Commands{
				idGenerator: tt.fields.idGenerator,
			}
			AssertValidation(t,
				context.Background(),
				c.AddOIDCAppCommand(
					tt.args.app,
					tt.args.clientSecretAlg,
				), tt.args.filter, tt.want)
		})
	}
}

func TestCommandSide_AddOIDCApplication(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx             context.Context
		oidcApp         *domain.OIDCApp
		resourceOwner   string
		secretGenerator crypto.Generator
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				oidcApp:       &domain.OIDCApp{},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, not found error",
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
					AppID:   "app1",
					AppName: "app",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid app, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "create oidc app basic, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
							"client1@project",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
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
				resourceOwner:   "org1",
				secretGenerator: GetMockSecretGenerator(t),
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
					ClientSecretString:       "a",
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
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			got, err := r.AddOIDCApplication(tt.args.ctx, tt.args.oidcApp, tt.args.resourceOwner, tt.args.secretGenerator)
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
				err: errors.IsErrorInvalidArgument,
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
				err: errors.IsErrorInvalidArgument,
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
				err: errors.IsErrorInvalidArgument,
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
				err: errors.IsNotFound,
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
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
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
				err: errors.IsPreconditionFailed,
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
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
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
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx             context.Context
		appID           string
		projectID       string
		resourceOwner   string
		secretGenerator crypto.Generator
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "no appid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "",
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
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
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsNotFound,
			},
		},
		{
			name: "change secret, ok",
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
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
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
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
						),
					),
				),
			},
			args: args{
				ctx:             context.Background(),
				projectID:       "project1",
				appID:           "app1",
				resourceOwner:   "org1",
				secretGenerator: GetMockSecretGenerator(t),
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
					ClientSecretString:       "a",
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
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeOIDCApplicationSecret(tt.args.ctx, tt.args.projectID, tt.args.appID, tt.args.resourceOwner, tt.args.secretGenerator)
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
