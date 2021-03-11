package command

import (
	"context"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
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
		//{
		//	name: "no aggregate id, invalid argument error",
		//	fields: fields{
		//		eventstore: eventstoreExpect(
		//			t,
		//		),
		//	},
		//	args: args{
		//		ctx:           context.Background(),
		//		oidcApp:        &domain.OIDCApp{},
		//		resourceOwner: "org1",
		//	},
		//	res: res{
		//		err: caos_errs.IsErrorInvalidArgument,
		//	},
		//},
		//{
		//	name: "project not existing, not found error",
		//	fields: fields{
		//		eventstore: eventstoreExpect(
		//			t,
		//			expectFilter(),
		//		),
		//	},
		//	args: args{
		//		ctx: context.Background(),
		//		oidcApp: &domain.OIDCApp{
		//			ObjectRoot: models.ObjectRoot{
		//				AggregateID: "project1",
		//			},
		//			AppID:   "app1",
		//			AppName: "app",
		//		},
		//		resourceOwner: "org1",
		//	},
		//	res: res{
		//		err: caos_errs.IsPreconditionFailed,
		//	},
		//},
		//{
		//	name: "invalid app, invalid argument error",
		//	fields: fields{
		//		eventstore: eventstoreExpect(
		//			t,
		//			expectFilter(
		//				eventFromEventPusher(
		//					project.NewProjectAddedEvent(context.Background(),
		//						&project.NewAggregate("project1", "org1").Aggregate,
		//						"project"),
		//				),
		//			),
		//		),
		//	},
		//	args: args{
		//		ctx: context.Background(),
		//		oidcApp: &domain.OIDCApp{
		//			ObjectRoot: models.ObjectRoot{
		//				AggregateID: "project1",
		//			},
		//			AppID:   "app1",
		//			AppName: "",
		//		},
		//		resourceOwner: "org1",
		//	},
		//	res: res{
		//		err: caos_errs.IsErrorInvalidArgument,
		//	},
		//},
		{
			name: "create api app basic, ok",
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
									[]string{},
									[]domain.OIDCResponseType{domain.OIDCResponseTypeIDTokenToken},
									[]domain.OIDCGrantType{domain.OIDCGrantTypeImplicit},
									domain.OIDCApplicationTypeNative,
									domain.OIDCAuthMethodTypeBasic,
									nil,
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
					AuthMethodType:           domain.OIDCAuthMethodTypeBasic,
					OIDCVersion:              domain.OIDCVersionV1,
					RedirectUris:             []string{},
					ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeIDTokenToken},
					GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeImplicit},
					ApplicationType:          domain.OIDCApplicationTypeNative,
					PostLogoutRedirectUris:   []string{},
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
					AppID:              "app1",
					AppName:            "app",
					ClientID:           "client1@project",
					ClientSecretString: "a",
					AuthMethodType:     domain.OIDCAuthMethodTypeBasic,
					State:              domain.AppStateActive,
				},
			},
		},
		{
			name: "create api app jwt, ok",
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
								project.NewAPIConfigAddedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"app1",
									"client1@project",
									nil,
									domain.APIAuthMethodTypePrivateKeyJWT),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewAddApplicationUniqueConstraint("app", "project1")),
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
					AppName:        "app",
					AuthMethodType: domain.OIDCAuthMethodTypePost,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.OIDCApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:          "app1",
					AppName:        "app",
					ClientID:       "client1@project",
					AuthMethodType: domain.OIDCAuthMethodTypePost,
					State:          domain.AppStateActive,
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

//
//func TestCommandSide_ChangeOIDCApplication(t *testing.T) {
//	type fields struct {
//		eventstore *eventstore.Eventstore
//	}
//	type args struct {
//		ctx           context.Context
//		apiApp        *domain.APIApp
//		resourceOwner string
//	}
//	type res struct {
//		want *domain.APIApp
//		err  func(error) bool
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		res    res
//	}{
//		{
//			name: "invalid app, invalid argument error",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//				),
//			},
//			args: args{
//				ctx: context.Background(),
//				apiApp: &domain.APIApp{
//					ObjectRoot: models.ObjectRoot{
//						AggregateID: "project1",
//					},
//					AppID: "app1",
//				},
//				resourceOwner: "org1",
//			},
//			res: res{
//				err: caos_errs.IsErrorInvalidArgument,
//			},
//		},
//		{
//			name: "missing appid, invalid argument error",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//				),
//			},
//			args: args{
//				ctx: context.Background(),
//				apiApp: &domain.APIApp{
//					ObjectRoot: models.ObjectRoot{
//						AggregateID: "project1",
//					},
//					AppID:          "",
//					AppName:        "app",
//					AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
//				},
//				resourceOwner: "org1",
//			},
//			res: res{
//				err: caos_errs.IsErrorInvalidArgument,
//			},
//		},
//		{
//			name: "missing aggregateid, invalid argument error",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//				),
//			},
//			args: args{
//				ctx: context.Background(),
//				apiApp: &domain.APIApp{
//					ObjectRoot: models.ObjectRoot{
//						AggregateID: "",
//					},
//					AppID:          "appid",
//					AppName:        "app",
//					AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
//				},
//				resourceOwner: "org1",
//			},
//			res: res{
//				err: caos_errs.IsErrorInvalidArgument,
//			},
//		},
//		{
//			name: "app not existing, not found error",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//					expectFilter(),
//				),
//			},
//			args: args{
//				ctx: context.Background(),
//				apiApp: &domain.APIApp{
//					ObjectRoot: models.ObjectRoot{
//						AggregateID: "project1",
//					},
//					AppID:   "app1",
//					AppName: "app",
//				},
//				resourceOwner: "org1",
//			},
//			res: res{
//				err: caos_errs.IsNotFound,
//			},
//		},
//		{
//			name: "no changes, precondition error",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//					expectFilter(
//						eventFromEventPusher(
//							project.NewApplicationAddedEvent(context.Background(),
//								&project.NewAggregate("project1", "org1").Aggregate,
//								"app1",
//								"app",
//							),
//						),
//						eventFromEventPusher(
//							project.NewAPIConfigAddedEvent(context.Background(),
//								&project.NewAggregate("project1", "org1").Aggregate,
//								"app1",
//								"client1@project",
//								nil,
//								domain.APIAuthMethodTypePrivateKeyJWT),
//						),
//					),
//				),
//			},
//			args: args{
//				ctx: context.Background(),
//				apiApp: &domain.APIApp{
//					ObjectRoot: models.ObjectRoot{
//						AggregateID: "project1",
//					},
//					AppID:          "app1",
//					AppName:        "app",
//					AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
//				},
//				resourceOwner: "org1",
//			},
//			res: res{
//				err: caos_errs.IsPreconditionFailed,
//			},
//		},
//		{
//			name: "change api app, ok",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//					expectFilter(
//						eventFromEventPusher(
//							project.NewApplicationAddedEvent(context.Background(),
//								&project.NewAggregate("project1", "org1").Aggregate,
//								"app1",
//								"app",
//							),
//						),
//						eventFromEventPusher(
//							project.NewAPIConfigAddedEvent(context.Background(),
//								&project.NewAggregate("project1", "org1").Aggregate,
//								"app1",
//								"client1@project",
//								&crypto.CryptoValue{
//									CryptoType: crypto.TypeEncryption,
//									Algorithm:  "enc",
//									KeyID:      "id",
//									Crypted:    []byte("a"),
//								},
//								domain.APIAuthMethodTypeBasic),
//						),
//					),
//					expectPush(
//						[]*repository.Event{
//							eventFromEventPusher(
//								newAPIAppChangedEvent(context.Background(),
//									"app1",
//									"project1",
//									"org1",
//									domain.APIAuthMethodTypePrivateKeyJWT),
//							),
//						},
//					),
//				),
//			},
//			args: args{
//				ctx: context.Background(),
//				apiApp: &domain.APIApp{
//					ObjectRoot: models.ObjectRoot{
//						AggregateID: "project1",
//					},
//					AppID:          "app1",
//					AppName:        "app",
//					AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
//				},
//				resourceOwner: "org1",
//			},
//			res: res{
//				want: &domain.APIApp{
//					ObjectRoot: models.ObjectRoot{
//						AggregateID:   "project1",
//						ResourceOwner: "org1",
//					},
//					AppID:          "app1",
//					AppName:        "app",
//					ClientID:       "client1@project",
//					AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
//					State:          domain.AppStateActive,
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &Commands{
//				eventstore: tt.fields.eventstore,
//			}
//			got, err := r.ChangeAPIApplication(tt.args.ctx, tt.args.apiApp, tt.args.resourceOwner)
//			if tt.res.err == nil {
//				assert.NoError(t, err)
//			}
//			if tt.res.err != nil && !tt.res.err(err) {
//				t.Errorf("got wrong err: %v ", err)
//			}
//			if tt.res.err == nil {
//				assert.Equal(t, tt.res.want, got)
//			}
//		})
//	}
//}
//
//func TestCommandSide_ChangeOIDCApplicationSecret(t *testing.T) {
//	type fields struct {
//		eventstore      *eventstore.Eventstore
//		secretGenerator crypto.Generator
//	}
//	type args struct {
//		ctx           context.Context
//		appID         string
//		projectID     string
//		resourceOwner string
//	}
//	type res struct {
//		want *domain.APIApp
//		err  func(error) bool
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		res    res
//	}{
//		{
//			name: "no projectid, invalid argument error",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//				),
//			},
//			args: args{
//				ctx:           context.Background(),
//				appID:         "app1",
//				resourceOwner: "org1",
//			},
//			res: res{
//				err: caos_errs.IsErrorInvalidArgument,
//			},
//		},
//		{
//			name: "no appid, invalid argument error",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//				),
//			},
//			args: args{
//				ctx:           context.Background(),
//				projectID:     "project1",
//				appID:         "",
//				resourceOwner: "org1",
//			},
//			res: res{
//				err: caos_errs.IsErrorInvalidArgument,
//			},
//		},
//		{
//			name: "app not existing, not found error",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//					expectFilter(),
//				),
//			},
//			args: args{
//				ctx:           context.Background(),
//				projectID:     "project1",
//				appID:         "app1",
//				resourceOwner: "org1",
//			},
//			res: res{
//				err: caos_errs.IsNotFound,
//			},
//		},
//		{
//			name: "change secret, ok",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//					expectFilter(
//						eventFromEventPusher(
//							project.NewApplicationAddedEvent(context.Background(),
//								&project.NewAggregate("project1", "org1").Aggregate,
//								"app1",
//								"app",
//							),
//						),
//						eventFromEventPusher(
//							project.NewAPIConfigAddedEvent(context.Background(),
//								&project.NewAggregate("project1", "org1").Aggregate,
//								"app1",
//								"client1@project",
//								&crypto.CryptoValue{
//									CryptoType: crypto.TypeEncryption,
//									Algorithm:  "enc",
//									KeyID:      "id",
//									Crypted:    []byte("a"),
//								},
//								domain.APIAuthMethodTypeBasic),
//						),
//					),
//					expectPush(
//						[]*repository.Event{
//							eventFromEventPusher(
//								project.NewAPIConfigSecretChangedEvent(context.Background(),
//									&project.NewAggregate("project1", "org1").Aggregate,
//									"app1",
//									&crypto.CryptoValue{
//										CryptoType: crypto.TypeEncryption,
//										Algorithm:  "enc",
//										KeyID:      "id",
//										Crypted:    []byte("a"),
//									}),
//							),
//						},
//					),
//				),
//				secretGenerator: GetMockSecretGenerator(t),
//			},
//			args: args{
//				ctx:           context.Background(),
//				projectID:     "project1",
//				appID:         "app1",
//				resourceOwner: "org1",
//			},
//			res: res{
//				want: &domain.APIApp{
//					ObjectRoot: models.ObjectRoot{
//						AggregateID:   "project1",
//						ResourceOwner: "org1",
//					},
//					AppID:              "app1",
//					AppName:            "app",
//					ClientID:           "client1@project",
//					ClientSecretString: "a",
//					AuthMethodType:     domain.APIAuthMethodTypeBasic,
//					State:              domain.AppStateActive,
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &Commands{
//				eventstore:                 tt.fields.eventstore,
//				applicationSecretGenerator: tt.fields.secretGenerator,
//			}
//			got, err := r.ChangeAPIApplicationSecret(tt.args.ctx, tt.args.projectID, tt.args.appID, tt.args.resourceOwner)
//			if tt.res.err == nil {
//				assert.NoError(t, err)
//			}
//			if tt.res.err != nil && !tt.res.err(err) {
//				t.Errorf("got wrong err: %v ", err)
//			}
//			if tt.res.err == nil {
//				assert.Equal(t, tt.res.want, got)
//			}
//		})
//	}
//}
//
//func newOIDCAppChangedEvent(ctx context.Context, appID, projectID, resourceOwner string, authMethodType domain.APIAuthMethodType) *project.APIConfigChangedEvent {
//	changes := []project.APIConfigChanges{
//		project.ChangeAPIAuthMethodType(authMethodType),
//	}
//	event, _ := project.NewAPIConfigChangedEvent(ctx,
//		&project.NewAggregate(projectID, resourceOwner).Aggregate,
//		appID,
//		changes,
//	)
//	return event
//}
