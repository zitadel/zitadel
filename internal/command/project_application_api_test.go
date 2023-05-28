package command

import (
	"context"
	"testing"

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

func TestAddAPIConfig(t *testing.T) {
	type fields struct {
		idGenerator id.Generator
	}
	type args struct {
		a      *project.Aggregate
		appID  string
		name   string
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
				a:     agg,
				appID: "",
				name:  "name",
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-XHsKt", "Errors.Invalid.Argument"),
			},
		},
		{
			name:   "invalid name",
			fields: fields{},
			args: args{
				a:     agg,
				appID: "appID",
				name:  "",
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-F7g21", "Errors.Invalid.Argument"),
			},
		},
		{
			name:   "project not exists",
			fields: fields{},
			args: args{
				a:     agg,
				appID: "id",
				name:  "name",
				filter: NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).
					Filter(),
			},
			want: Want{
				CreateErr: errors.ThrowNotFound(nil, "PROJE-Sf2gb", "Errors.Project.NotFound"),
			},
		},
		{
			name: "correct without client secret",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "clientID"),
			},
			args: args{
				a:     agg,
				appID: "appID",
				name:  "name",
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
					project.NewApplicationAddedEvent(
						ctx,
						&agg.Aggregate,
						"appID",
						"name",
					),
					project.NewAPIConfigAddedEvent(ctx, &agg.Aggregate,
						"appID",
						"clientID@project",
						nil,
						domain.APIAuthMethodTypePrivateKeyJWT,
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				idGenerator: tt.fields.idGenerator,
			}
			AssertValidation(t,
				context.Background(),
				c.AddAPIAppCommand(
					&addAPIApp{
						AddApp: AddApp{
							Aggregate: *tt.args.a,
							ID:        tt.args.appID,
							Name:      tt.args.name,
						},
						AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
					},
					nil,
				), tt.args.filter, tt.want)
		})
	}
}

func TestCommandSide_AddAPIApplication(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx             context.Context
		apiApp          *domain.APIApp
		resourceOwner   string
		secretGenerator crypto.Generator
	}
	type res struct {
		want *domain.APIApp
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
				apiApp:        &domain.APIApp{},
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
				apiApp: &domain.APIApp{
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
				apiApp: &domain.APIApp{
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
			name: "create api app basic, ok",
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
						project.NewAPIConfigAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"client1@project",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
							domain.APIAuthMethodTypeBasic),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "app1", "client1"),
			},
			args: args{
				ctx: context.Background(),
				apiApp: &domain.APIApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppName:        "app",
					AuthMethodType: domain.APIAuthMethodTypeBasic,
				},
				resourceOwner:   "org1",
				secretGenerator: GetMockSecretGenerator(t),
			},
			res: res{
				want: &domain.APIApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:              "app1",
					AppName:            "app",
					ClientID:           "client1@project",
					ClientSecretString: "a",
					AuthMethodType:     domain.APIAuthMethodTypeBasic,
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
						project.NewAPIConfigAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"client1@project",
							nil,
							domain.APIAuthMethodTypePrivateKeyJWT),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "app1", "client1"),
			},
			args: args{
				ctx: context.Background(),
				apiApp: &domain.APIApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppName:        "app",
					AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.APIApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:          "app1",
					AppName:        "app",
					ClientID:       "client1@project",
					AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
					State:          domain.AppStateActive,
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
			got, err := r.AddAPIApplication(tt.args.ctx, tt.args.apiApp, tt.args.resourceOwner, tt.args.secretGenerator)
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

func TestCommandSide_ChangeAPIApplication(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		apiApp        *domain.APIApp
		resourceOwner string
	}
	type res struct {
		want *domain.APIApp
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "missing appid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				apiApp: &domain.APIApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:          "",
					AppName:        "app",
					AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
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
				apiApp: &domain.APIApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "",
					},
					AppID:          "appid",
					AppName:        "app",
					AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
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
				apiApp: &domain.APIApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:   "app1",
					AppName: "app",
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
							project.NewAPIConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"client1@project",
								nil,
								domain.APIAuthMethodTypePrivateKeyJWT),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				apiApp: &domain.APIApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:          "app1",
					AppName:        "app",
					AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			name: "change api app, ok",
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
							project.NewAPIConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"client1@project",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								domain.APIAuthMethodTypeBasic),
						),
					),
					expectPush(
						newAPIAppChangedEvent(context.Background(),
							"app1",
							"project1",
							"org1",
							domain.APIAuthMethodTypePrivateKeyJWT),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				apiApp: &domain.APIApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:          "app1",
					AppName:        "app",
					AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.APIApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:          "app1",
					AppName:        "app",
					ClientID:       "client1@project",
					AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
					State:          domain.AppStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeAPIApplication(tt.args.ctx, tt.args.apiApp, tt.args.resourceOwner)
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

func TestCommandSide_ChangeAPIApplicationSecret(t *testing.T) {
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
		want *domain.APIApp
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
							project.NewAPIConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"client1@project",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								domain.APIAuthMethodTypeBasic),
						),
					),
					expectPush(
						project.NewAPIConfigSecretChangedEvent(context.Background(),
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
				want: &domain.APIApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:              "app1",
					AppName:            "app",
					ClientID:           "client1@project",
					ClientSecretString: "a",
					AuthMethodType:     domain.APIAuthMethodTypeBasic,
					State:              domain.AppStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeAPIApplicationSecret(tt.args.ctx, tt.args.projectID, tt.args.appID, tt.args.resourceOwner, tt.args.secretGenerator)
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

func newAPIAppChangedEvent(ctx context.Context, appID, projectID, resourceOwner string, authMethodType domain.APIAuthMethodType) *project.APIConfigChangedEvent {
	changes := []project.APIConfigChanges{
		project.ChangeAPIAuthMethodType(authMethodType),
	}
	event, _ := project.NewAPIConfigChangedEvent(ctx,
		&project.NewAggregate(projectID, resourceOwner).Aggregate,
		appID,
		changes,
	)
	return event
}
