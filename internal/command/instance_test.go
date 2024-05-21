package command

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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

func projectClientIDs() []string {
	return []string{"clientID", "clientID", "clientID", "clientID"}
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

func orgFilters(ctx context.Context, orgID string, machine, human bool) []expect {
	orgAgg := org.NewAggregate(orgID)

	filters := []expect{
		expectFilter(),
		expectFilter(
			org.NewOrgAddedEvent(ctx, &orgAgg.Aggregate, ""),
		),
	}
	if machine {
		filters = append(filters, machineFilters(true)...)
		filters = append(filters, adminMemberFilters(orgID, "USER")...)
	}
	if human {
		filters = append(filters, humanFilters()...)
		filters = append(filters, adminMemberFilters(orgID, "USER")...)
	}

	return append(filters,
		projectFilters()...,
	)
}

func orgEvents(ctx context.Context, instanceID, orgID, name, projectID, defaultDomain string, externalSecure bool, machine, human bool) []eventstore.Command {
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
	if machine {
		machineID := "USER-MACHINE"
		events = append(events, machineEvents(ctx, instanceID, orgID, machineID, "PAT")...)
		owner = machineID
	}
	if human {
		userID := "USER"
		events = append(events, humanEvents(ctx, instanceID, orgID, userID)...)
		owner = userID
	}

	events = append(events, projectAddedEvents(ctx, instanceID, orgID, projectID, owner, externalSecure)...)
	return events
}

func orgIDs() []string {
	return slices.Concat([]string{"USER-MACHINE", "PAT", "USER"}, projectClientIDs())
}

func generatedDomainEvents(ctx context.Context, instanceID, defaultDomain string) []eventstore.Command {
	instanceAgg := instance.NewAggregate(instanceID)
	return []eventstore.Command{
		instance.NewDomainAddedEvent(ctx, &instanceAgg.Aggregate, defaultDomain, true),
		instance.NewDomainPrimarySetEvent(ctx, &instanceAgg.Aggregate, defaultDomain),
	}
}

func domainFilters() []expect {
	return []expect{}
}

func humanFilters() []expect {
	return []expect{
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
	}
}

func machineFilters(pat bool) []expect {
	filters := []expect{
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
	}
	if pat {
		filters = append(filters,
			expectFilter(),
			expectFilter(),
		)
	}
	return filters
}

func projectFilters() []expect {
	return []expect{
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
	}
}

func adminMemberFilters(orgID, userID string) []expect {
	return []expect{
		expectFilter(
			addHumanEvent(context.Background(), orgID, userID),
		),
		expectFilter(),
		expectFilter(
			addHumanEvent(context.Background(), orgID, userID),
		),
		expectFilter(),
	}
}

func humanEvents(ctx context.Context, instanceID, orgID, userID string) []eventstore.Command {
	agg := user.NewAggregate(userID, orgID)
	instanceAgg := instance.NewAggregate(instanceID)
	orgAgg := org.NewAggregate(orgID)
	return []eventstore.Command{
		addHumanEvent(ctx, orgID, userID),
		user.NewHumanEmailVerifiedEvent(ctx, &agg.Aggregate),
		org.NewMemberAddedEvent(ctx, &orgAgg.Aggregate, userID, domain.RoleOrgOwner),
		instance.NewMemberAddedEvent(ctx, &instanceAgg.Aggregate, userID, domain.RoleIAMOwner),
	}
}

func addHumanEvent(ctx context.Context, orgID, userID string) *user.HumanAddedEvent {
	agg := user.NewAggregate(userID, orgID)
	return func() *user.HumanAddedEvent {
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
	}()
}

// machineEvents all events from setup to create the machine user, machinekey can't be tested here, as the public key is not provided and as such the value in the event can't be expected
func machineEvents(ctx context.Context, instanceID, orgID, userID, patID string) []eventstore.Command {
	agg := user.NewAggregate(userID, orgID)
	instanceAgg := instance.NewAggregate(instanceID)
	orgAgg := org.NewAggregate(orgID)
	events := []eventstore.Command{addMachineEvent(ctx, orgID, userID)}
	if patID != "" {
		events = append(events,
			user.NewPersonalAccessTokenAddedEvent(
				ctx,
				&agg.Aggregate,
				patID,
				time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
				nil,
			),
		)
	}
	return append(events,
		org.NewMemberAddedEvent(ctx, &orgAgg.Aggregate, userID, domain.RoleOrgOwner),
		instance.NewMemberAddedEvent(ctx, &instanceAgg.Aggregate, userID, domain.RoleIAMOwner),
	)
}

func addMachineEvent(ctx context.Context, orgID, userID string) *user.MachineAddedEvent {
	agg := user.NewAggregate(userID, orgID)
	return user.NewMachineAddedEvent(ctx,
		&agg.Aggregate,
		"zitadel-admin-machine",
		"ZITADEL-machine",
		"Admin",
		false,
		domain.OIDCTokenTypeBearer,
	)
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
		eventstore  func(t *testing.T) *eventstore.Eventstore
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
				eventstore: expectEventstore(
					slices.Concat(
						projectFilters(),
						[]expect{expectPush(
							projectAddedEvents(context.Background(),
								"INSTANCE",
								"ORG",
								"PROJECT",
								"owner",
								false,
							)...,
						),
						},
					)...,
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, projectClientIDs()...),
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
				eventstore:  tt.fields.eventstore(t),
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
		keyAlgorithm       crypto.EncryptionAlgorithm
	}
	type args struct {
		instanceAgg *instance.Aggregate
		orgAgg      *org.Aggregate
		machine     *AddMachine
		human       *AddHuman
	}
	type res struct {
		owner      string
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
					slices.Concat(
						humanFilters(),
						adminMemberFilters("ORG", "USER"),
						[]expect{
							expectPush(
								humanEvents(context.Background(),
									"INSTANCE",
									"ORG",
									"USER",
								)...,
							),
						},
					)...,
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
				owner:      "USER",
				pat:        false,
				machineKey: false,
				err:        nil,
			},
		},
		{
			name: "machine, ok",
			fields: fields{
				eventstore: expectEventstore(
					slices.Concat(
						machineFilters(true),
						adminMemberFilters("ORG", "USER-MACHINE"),
						[]expect{
							expectPush(
								machineEvents(context.Background(),
									"INSTANCE",
									"ORG",
									"USER-MACHINE",
									"PAT",
								)...,
							),
						},
					)...,
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "USER-MACHINE", "PAT"),
				roles: []authz.RoleMapping{
					{Role: domain.RoleOrgOwner, Permissions: []string{""}},
					{Role: domain.RoleIAMOwner, Permissions: []string{""}},
				},
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				instanceAgg: instance.NewAggregate("INSTANCE"),
				orgAgg:      org.NewAggregate("ORG"),
				machine: &AddMachine{
					Machine: &Machine{
						Username:        "zitadel-admin-machine",
						Name:            "ZITADEL-machine",
						Description:     "Admin",
						AccessTokenType: domain.OIDCTokenTypeBearer,
					},
					Pat: &AddPat{
						ExpirationDate: time.Time{},
						Scopes:         nil,
					},
					/* not predictable with the key value in the events
					MachineKey: &AddMachineKey{
						Type:           domain.AuthNKeyTypeJSON,
						ExpirationDate: time.Time{},
					},
					*/
				},
			},
			res: res{
				owner:      "USER-MACHINE",
				pat:        true,
				machineKey: false,
				err:        nil,
			},
		},
		{
			name: "human and machine, ok",
			fields: fields{
				eventstore: expectEventstore(
					slices.Concat(
						machineFilters(true),
						adminMemberFilters("ORG", "USER-MACHINE"),
						humanFilters(),
						adminMemberFilters("ORG", "USER"),
						[]expect{
							expectPush(
								slices.Concat(
									machineEvents(context.Background(),
										"INSTANCE",
										"ORG",
										"USER-MACHINE",
										"PAT",
									),
									humanEvents(context.Background(),
										"INSTANCE",
										"ORG",
										"USER",
									),
								)...,
							),
						},
					)...,
				),
				userPasswordHasher: mockPasswordHasher("x"),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "USER-MACHINE", "PAT", "USER"),
				roles: []authz.RoleMapping{
					{Role: domain.RoleOrgOwner, Permissions: []string{""}},
					{Role: domain.RoleIAMOwner, Permissions: []string{""}},
				},
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				instanceAgg: instance.NewAggregate("INSTANCE"),
				orgAgg:      org.NewAggregate("ORG"),
				machine: &AddMachine{
					Machine: &Machine{
						Username:        "zitadel-admin-machine",
						Name:            "ZITADEL-machine",
						Description:     "Admin",
						AccessTokenType: domain.OIDCTokenTypeBearer,
					},
					Pat: &AddPat{
						ExpirationDate: time.Time{},
						Scopes:         nil,
					},
					/* not predictable with the key value in the events
					MachineKey: &AddMachineKey{
						Type:           domain.AuthNKeyTypeJSON,
						ExpirationDate: time.Time{},
					},
					*/
				},
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
				owner:      "USER",
				pat:        true,
				machineKey: false,
				err:        nil,
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
				keyAlgorithm:       tt.fields.keyAlgorithm,
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
				assert.Equal(t, owner, tt.res.owner)
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

func TestCommandSide_setupDefaultOrg(t *testing.T) {
	type fields struct {
		eventstore         func(t *testing.T) *eventstore.Eventstore
		idGenerator        id.Generator
		userPasswordHasher *crypto.Hasher
		roles              []authz.RoleMapping
		keyAlgorithm       crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx          context.Context
		instanceAgg  *instance.Aggregate
		instanceName string
		orgName      string
		machine      *AddMachine
		human        *AddHuman
		ids          ZitadelConfig
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
			name: "human and machine, ok",
			fields: fields{
				eventstore: expectEventstore(
					slices.Concat(
						orgFilters(context.Background(),
							"ORG",
							true,
							true,
						),
						[]expect{
							expectPush(
								slices.Concat(
									orgEvents(context.Background(),
										"INSTANCE",
										"ORG",
										"ZITADEL",
										"PROJECT",
										"zitadel.domain",
										false,
										true,
										true,
									),
								)...,
							),
						},
					)...,
				),
				userPasswordHasher: mockPasswordHasher("x"),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, orgIDs()...),
				roles: []authz.RoleMapping{
					{Role: domain.RoleOrgOwner, Permissions: []string{""}},
					{Role: domain.RoleIAMOwner, Permissions: []string{""}},
				},
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:         authz.WithRequestedDomain(context.Background(), "DOMAIN"),
				instanceAgg: instance.NewAggregate("INSTANCE"),
				orgName:     "ZITADEL",
				machine: &AddMachine{
					Machine: &Machine{
						Username:        "zitadel-admin-machine",
						Name:            "ZITADEL-machine",
						Description:     "Admin",
						AccessTokenType: domain.OIDCTokenTypeBearer,
					},
					Pat: &AddPat{
						ExpirationDate: time.Time{},
						Scopes:         nil,
					},
					/* not predictable with the key value in the events
					MachineKey: &AddMachineKey{
						Type:           domain.AuthNKeyTypeJSON,
						ExpirationDate: time.Time{},
					},
					*/
				},
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
				pat:        true,
				machineKey: false,
				err:        nil,
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
				keyAlgorithm:       tt.fields.keyAlgorithm,
			}
			validations := make([]preparation.Validation, 0)
			pat, mk, err := setupDefaultOrg(tt.args.ctx, r, &validations, tt.args.instanceAgg, tt.args.orgName, tt.args.machine, tt.args.human, tt.args.ids)
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
