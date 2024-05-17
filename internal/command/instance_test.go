package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func projectAddedEvents(ctx context.Context, instanceID, orgID, id, owner string, externalSecure bool) []eventstore.Command {
	events := []eventstore.Command{
		project.NewProjectAddedEvent(ctx,
			&project.NewAggregate(id, orgID).Aggregate,
			"ZITADEL",
			false,
			false,
			false,
			domain.PrivateLabelingSettingUnspecified,
		),
		project.NewProjectMemberAddedEvent(ctx,
			&project.NewAggregate(id, orgID).Aggregate,
			owner,
			domain.RoleProjectOwner,
		),
		instance.NewIAMProjectSetEvent(ctx,
			&instance.NewAggregate(instanceID).Aggregate,
			id,
		),
	}
	events = append(events, apiAppEvents(ctx, orgID, id, "mgmt-id", "Management-API")...)
	events = append(events, apiAppEvents(ctx, orgID, id, "admin-id", "Admin-API")...)
	events = append(events, apiAppEvents(ctx, orgID, id, "auth-id", "Auth-API")...)

	consoleAppID := "console-id"
	consoleClientID := "clientID@zitadel"
	events = append(events, oidcAppEvents(ctx, orgID, id, consoleAppID, "Console", consoleClientID, externalSecure)...)
	events = append(events,
		instance.NewIAMConsoleSetEvent(ctx,
			&instance.NewAggregate(instanceID).Aggregate,
			&consoleClientID,
			&consoleAppID,
		),
	)
	return events
}

func apiAppEvents(ctx context.Context, orgID, projectID, id, name string) []eventstore.Command {
	return []eventstore.Command{
		project.NewApplicationAddedEvent(
			ctx,
			&project.NewAggregate(projectID, orgID).Aggregate,
			id,
			name,
		),
		project.NewAPIConfigAddedEvent(ctx,
			&project.NewAggregate(projectID, orgID).Aggregate,
			id,
			"clientID@zitadel",
			"",
			domain.APIAuthMethodTypePrivateKeyJWT,
		),
	}
}

func oidcAppEvents(ctx context.Context, orgID, projectID, id, name, clientID string, externalSecure bool) []eventstore.Command {
	return []eventstore.Command{
		project.NewApplicationAddedEvent(
			ctx,
			&project.NewAggregate(projectID, orgID).Aggregate,
			id,
			name,
		),
		project.NewOIDCConfigAddedEvent(ctx,
			&project.NewAggregate(projectID, orgID).Aggregate,
			domain.OIDCVersionV1,
			id,
			clientID,
			"",
			[]string{},
			[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
			[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
			domain.OIDCApplicationTypeUserAgent,
			domain.OIDCAuthMethodTypeNone,
			[]string{},
			!externalSecure,
			domain.OIDCTokenTypeBearer,
			false,
			false,
			false,
			0,
			nil,
			false,
		),
	}
}

func orgEvents(ctx context.Context, instanceID, orgID, name, projectID, defaultDomain string, externalSecure bool) []eventstore.Command {
	instanceAgg := instance.NewAggregate(instanceID)
	orgAgg := org.NewAggregate(orgID)
	events := []eventstore.Command{
		org.NewOrgAddedEvent(ctx, &orgAgg.Aggregate, name),
		org.NewDomainAddedEvent(ctx, &orgAgg.Aggregate, defaultDomain),
		org.NewDomainVerifiedEvent(ctx, &orgAgg.Aggregate, defaultDomain),
		org.NewDomainPrimarySetEvent(ctx, &orgAgg.Aggregate, defaultDomain),
		instance.NewDefaultOrgSetEventEvent(ctx, &instanceAgg.Aggregate, orgID),
	}
	owner := ""
	events = append(events, projectAddedEvents(ctx, instanceID, orgID, projectID, owner, externalSecure)...)
	return events
}

func humanEvents(ctx context.Context, instanceID, orgID, userID string) []eventstore.Command {
	agg := user.NewAggregate(userID, orgID)
	instanceAgg := instance.NewAggregate(instanceID)
	orgAgg := org.NewAggregate(orgID)
	return []eventstore.Command{
		func() *user.HumanAddedEvent {
			event := user.NewHumanAddedEvent(
				ctx,
				&agg.Aggregate,
				"zitadel-admin",
				"ZITADEL",
				"Admin",
				"",
				"ZITADEL Admin",
				language.English,
				0,
				"admin@zitadel.test",
				false,
			)
			event.AddPasswordData("$plain$x$password", false)
			return event
		}(),
		user.NewHumanEmailVerifiedEvent(ctx, &agg.Aggregate),
		org.NewMemberAddedEvent(ctx, &orgAgg.Aggregate, userID, domain.RoleOrgOwner),
		instance.NewMemberAddedEvent(ctx, &instanceAgg.Aggregate, userID, domain.RoleIAMOwner),
	}
}

func testSetup(ctx context.Context, c *Commands, validations []preparation.Validation) error {
	//nolint:staticcheck
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validations...)
	if err != nil {
		return err
	}

	_, err = c.eventstore.Push(ctx, cmds...)
	return err
}

func TestCommandSide_setupMinimalInterfaces(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx         context.Context
		instanceAgg *instance.Aggregate
		orgAgg      *org.Aggregate
		owner       string
		ids         ZitadelConfig
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
			name: "create, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(),
					expectFilter(),
					expectFilter(),
					expectPush(
						projectAddedEvents(context.Background(),
							"INSTANCE",
							"ORG",
							"PROJECT",
							"owner",
							false,
						)...,
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "clientID", "clientID", "clientID", "clientID"),
			},
			args: args{
				ctx:         authz.WithInstanceID(context.Background(), "INSTANCE"),
				instanceAgg: instance.NewAggregate("INSTANCE"),
				orgAgg:      org.NewAggregate("ORG"),
				owner:       "owner",
				ids: ZitadelConfig{
					instanceID:   "INSTANCE",
					orgID:        "ORG",
					projectID:    "PROJECT",
					consoleAppID: "console-id",
					authAppID:    "auth-id",
					mgmtAppID:    "mgmt-id",
					adminAppID:   "admin-id",
				},
			},
			res: res{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			validations := make([]preparation.Validation, 0)
			setupMinimalInterfaces(r, &validations, tt.args.instanceAgg, tt.args.orgAgg, tt.args.owner, tt.args.ids)

			err := testSetup(tt.args.ctx, r, validations)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_setupAdmins(t *testing.T) {
	type fields struct {
		eventstore         func(t *testing.T) *eventstore.Eventstore
		idGenerator        id.Generator
		userPasswordHasher *crypto.Hasher
		roles              []authz.RoleMapping
	}
	type args struct {
		instanceAgg *instance.Aggregate
		orgAgg      *org.Aggregate
		machine     *AddMachine
		human       *AddHuman
	}
	type res struct {
		pat        bool
		machineKey bool
		err        func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "human, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						org.NewDomainPolicyAddedEvent(
							context.Background(),
							&org.NewAggregate("id").Aggregate,
							true,
							true,
							true,
						),
					),
					expectFilter(
						org.NewPasswordComplexityPolicyAddedEvent(
							context.Background(),
							&org.NewAggregate("id").Aggregate,
							2,
							false,
							false,
							false,
							false,
						),
					),
					expectFilter(
						user.NewHumanAddedEvent(
							context.Background(),
							&user.NewAggregate("USER", "ORG").Aggregate,
							"zitadel-admin",
							"ZITADEL",
							"Admin",
							"",
							"ZITADEL Admin",
							language.English,
							0,
							"admin@zitadel.test",
							false,
						),
					),
					expectFilter(),
					expectFilter(
						user.NewHumanAddedEvent(
							context.Background(),
							&user.NewAggregate("USER", "ORG").Aggregate,
							"zitadel-admin",
							"ZITADEL",
							"Admin",
							"",
							"ZITADEL Admin",
							language.English,
							0,
							"admin@zitadel.test",
							false,
						),
					),
					expectFilter(),
					expectPush(
						humanEvents(context.Background(),
							"INSTANCE",
							"ORG",
							"USER",
						)...,
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "USER"),
				userPasswordHasher: mockPasswordHasher("x"),
				roles: []authz.RoleMapping{
					{Role: domain.RoleOrgOwner, Permissions: []string{""}},
					{Role: domain.RoleIAMOwner, Permissions: []string{""}},
				},
			},
			args: args{
				instanceAgg: instance.NewAggregate("INSTANCE"),
				orgAgg:      org.NewAggregate("ORG"),
				human: &AddHuman{
					Username:  "zitadel-admin",
					FirstName: "ZITADEL",
					LastName:  "Admin",
					Email: Email{
						Address:  domain.EmailAddress("admin@zitadel.test"),
						Verified: true,
					},
					PreferredLanguage:      language.English,
					Password:               "password",
					PasswordChangeRequired: false,
				},
			},
			res: res{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:         tt.fields.eventstore(t),
				idGenerator:        tt.fields.idGenerator,
				zitadelRoles:       tt.fields.roles,
				userPasswordHasher: tt.fields.userPasswordHasher,
			}
			validations := make([]preparation.Validation, 0)
			owner, pat, mk, err := setupAdmins(r, &validations, tt.args.instanceAgg, tt.args.orgAgg, tt.args.machine, tt.args.human)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}

			err = testSetup(context.Background(), r, validations)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}

			if tt.res.err == nil {
				assert.NotEmpty(t, owner)
				if tt.res.pat {
					assert.NotNil(t, pat)
				}
				if tt.res.machineKey {
					assert.NotNil(t, mk)
				}
			}
		})
	}
}

func TestCommandSide_UpdateInstance(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx  context.Context
		name string
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
			name: "empty name, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "instance not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "INSTANCE_CHANGED",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "instance removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewInstanceAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
						eventFromEventPusher(
							instance.NewInstanceRemovedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
								nil,
							),
						),
					),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "INSTANCE_CHANGED",
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
							instance.NewInstanceAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
					),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "INSTANCE",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "instance change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewInstanceAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
					),
					expectPush(
						instance.NewInstanceChangedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"INSTANCE_CHANGED",
						),
					),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "INSTANCE_CHANGED",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.UpdateInstance(tt.args.ctx, tt.args.name)
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

func TestCommandSide_RemoveInstance(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
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
			name: "instance not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				instanceID: "INSTANCE",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "instance removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewInstanceAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
						eventFromEventPusher(
							instance.NewInstanceRemovedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
								nil,
							),
						),
					),
				),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				instanceID: "INSTANCE",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "instance remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewInstanceAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"instance.domain",
								true,
							),
						),
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"custom.domain",
								false,
							),
						),
					),
					expectPush(
						instance.NewInstanceRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"INSTANCE",
							[]string{
								"instance.domain",
								"custom.domain",
							},
						),
					),
				),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				instanceID: "INSTANCE",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveInstance(tt.args.ctx, tt.args.instanceID)
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
