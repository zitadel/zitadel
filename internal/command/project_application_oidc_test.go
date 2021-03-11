package command

import (
	"context"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/id"
	id_mock "github.com/caos/zitadel/internal/id/mock"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCommandSide_AddOIDCApplication(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		idGenerator     id.Generator
		secretGenerator crypto.Generator
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsPreconditionFailed,
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
								"project"),
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
				err: caos_errs.IsErrorInvalidArgument,
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
								"project"),
						),
					),
					expectPush(
						[]*repository.Event{
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
									time.Minute*1),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewAddApplicationUniqueConstraint("app", "project1")),
					),
				),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "app1", "client1"),
				secretGenerator: GetMockSecretGenerator(t),
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
					ClockSkew:                time.Minute * 1,
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
					ClockSkew:                time.Minute * 1,
					State:                    domain.AppStateActive,
					Compliance:               &domain.Compliance{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:                 tt.fields.eventstore,
				idGenerator:                tt.fields.idGenerator,
				applicationSecretGenerator: tt.fields.secretGenerator,
			}
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
				err: caos_errs.IsErrorInvalidArgument,
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
					AppName:        "app",
					AuthMethodType: domain.OIDCAuthMethodTypePost,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
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
					AppName:        "app",
					AuthMethodType: domain.OIDCAuthMethodTypePost,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
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
					AppID:   "app1",
					AppName: "app",
				},
				resourceOwner: "org1",
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
								time.Minute*1),
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
					ClockSkew:                time.Minute * 1,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
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
								time.Minute*1),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newOIDCAppChangedEvent(context.Background(),
									"app1",
									"project1",
									"org1"),
							),
						},
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
					ClockSkew:                time.Minute * 2,
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
					ClockSkew:                time.Minute * 2,
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
		eventstore      *eventstore.Eventstore
		secretGenerator crypto.Generator
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsNotFound,
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
								time.Minute*1),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								project.NewOIDCConfigSecretChangedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"app1",
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("a"),
									}),
							),
						},
					),
				),
				secretGenerator: GetMockSecretGenerator(t),
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
					ClockSkew:                time.Minute * 1,
					State:                    domain.AppStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:                 tt.fields.eventstore,
				applicationSecretGenerator: tt.fields.secretGenerator,
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
		project.ChangeClockSkew(time.Minute * 2),
	}
	event, _ := project.NewOIDCConfigChangedEvent(ctx,
		&project.NewAggregate(projectID, resourceOwner).Aggregate,
		appID,
		changes,
	)
	return event
}
