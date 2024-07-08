package command

import (
	"context"
	"io"
	"testing"

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

func TestAddAPIConfig(t *testing.T) {
	type fields struct {
		idGenerator id_generator.Generator
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
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "PROJE-XHsKt", "Errors.Invalid.Argument"),
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
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "PROJE-F7g21", "Errors.Invalid.Argument"),
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
				CreateErr: zerrors.ThrowNotFound(nil, "PROJE-Sf2gb", "Errors.Project.NotFound"),
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
						"clientID",
						"",
						domain.APIAuthMethodTypePrivateKeyJWT,
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{}
			id_generator.SetGenerator(tt.fields.idGenerator)
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
				), tt.args.filter, tt.want)
		})
	}
}

func TestCommandSide_AddAPIApplication(t *testing.T) {
	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id_generator.Generator
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
			name: "no aggregate id, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				apiApp:        &domain.APIApp{},
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
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "create api app basic, ok",
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
						project.NewAPIConfigAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"client1",
							"secret",
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
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.APIApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:              "app1",
					AppName:            "app",
					ClientID:           "client1",
					ClientSecretString: "secret",
					AuthMethodType:     domain.APIAuthMethodTypeBasic,
					State:              domain.AppStateActive,
				},
			},
		},
		{
			name: "create api app basic old ID format, ok",
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
						project.NewAPIConfigAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"client1@project1",
							"secret",
							domain.APIAuthMethodTypeBasic),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "app1", "client1@project1"),
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
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.APIApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:              "app1",
					AppName:            "app",
					ClientID:           "client1@project1",
					ClientSecretString: "secret",
					AuthMethodType:     domain.APIAuthMethodTypeBasic,
					State:              domain.AppStateActive,
				},
			},
		},
		{
			name: "create api app jwt, ok",
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
						project.NewAPIConfigAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"client1",
							"",
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
					ClientID:       "client1",
					AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
					State:          domain.AppStateActive,
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
			got, err := r.AddAPIApplication(tt.args.ctx, tt.args.apiApp, tt.args.resourceOwner)
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
		eventstore func(t *testing.T) *eventstore.Eventstore
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
				eventstore: expectEventstore(),
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
							project.NewAPIConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"client1@project",
								"",
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
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "change api app, ok",
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
							project.NewAPIConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"client1@project",
								"secret",
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
				eventstore:      tt.fields.eventstore(t),
				newHashedSecret: mockHashedSecret("secret"),
				defaultSecretGenerators: &SecretGenerators{
					ClientSecret: emptyConfig,
				},
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
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		appID         string
		projectID     string
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
							project.NewAPIConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"client1@project",
								"secret",
								domain.APIAuthMethodTypeBasic),
						),
					),
					expectPush(
						project.NewAPIConfigSecretChangedEvent(context.Background(),
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
				want: &domain.APIApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:              "app1",
					AppName:            "app",
					ClientID:           "client1@project",
					ClientSecretString: "secret",
					AuthMethodType:     domain.APIAuthMethodTypeBasic,
					State:              domain.AppStateActive,
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
			got, err := r.ChangeAPIApplicationSecret(tt.args.ctx, tt.args.projectID, tt.args.appID, tt.args.resourceOwner)
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

func TestCommands_VerifyAPIClientSecret(t *testing.T) {
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
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-DFnbf", "Errors.Project.App.NotExisting"),
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
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-Bf3fw", "Errors.Project.App.IsNotAPI"),
		},
		{
			name: "no secret set",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(
						project.NewApplicationAddedEvent(context.Background(), &agg.Aggregate, "appID", "appName"),
					),
					eventFromEventPusher(
						project.NewAPIConfigAddedEvent(context.Background(), &agg.Aggregate, "appID", "clientID", "", domain.APIAuthMethodTypePrivateKeyJWT),
					),
				),
			),
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-D3t5g", "Errors.Project.App.APIConfigInvalid"),
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
						project.NewAPIConfigAddedEvent(context.Background(), &agg.Aggregate, "appID", "clientID", hashedSecret, domain.APIAuthMethodTypePrivateKeyJWT),
					),
				),
				expectPush(
					project.NewAPIConfigSecretCheckSucceededEvent(context.Background(), &agg.Aggregate, "appID"),
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
						project.NewAPIConfigAddedEvent(context.Background(), &agg.Aggregate, "appID", "clientID", hashedSecret, domain.APIAuthMethodTypePrivateKeyJWT),
					),
				),
				expectPush(
					project.NewAPIConfigSecretCheckFailedEvent(context.Background(), &agg.Aggregate, "appID"),
				),
			),
			wantErr: zerrors.ThrowInvalidArgument(err, "COMMAND-SADfg", "Errors.Project.App.ClientSecretInvalid"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:   tt.eventstore(t),
				secretHasher: hasher,
			}
			err := c.VerifyAPIClientSecret(context.Background(), "projectID", "appID", tt.secret)
			c.jobs.Wait()
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
