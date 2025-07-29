package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	settings "github.com/zitadel/zitadel/internal/repository/organization_settings"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_SetSettingsOrganization(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx      context.Context
		settings *SetOrganizationSettings
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
			name: "org id missing, invalid argument error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				settings: &SetOrganizationSettings{
					OrganizationID:              "",
					OrganizationScopedUsernames: boolPtr(true),
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "org not found, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				settings: &SetOrganizationSettings{
					OrganizationID:              "org1",
					OrganizationScopedUsernames: boolPtr(true),
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "settings already existing, no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterPreOrganizationSettings("org1", true, true, true),
					expectFilterOrgDomainPolicy(false),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				settings: &SetOrganizationSettings{
					OrganizationID:              "org1",
					OrganizationScopedUsernames: boolPtr(true),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "org1",
				},
			},
		},
		{
			name: "settings set, new",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterPreOrganizationSettings("org1", true, false, false),
					expectFilterOrgDomainPolicy(false),
					expectFilterOrganizationScopedUsernames(false, "username1", "username2", "username3"),
					expectPush(
						settings.NewOrganizationSettingsAddedEvent(context.Background(),
							&settings.NewAggregate("org1", "org1").Aggregate,
							[]string{"username1", "username2", "username3"},
							true,
							false,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				settings: &SetOrganizationSettings{
					OrganizationID:              "org1",
					OrganizationScopedUsernames: boolPtr(true),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "org1",
				},
			},
		},
		{
			name: "settings set, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterPreOrganizationSettings("org1", true, false, true),
					expectFilterOrgDomainPolicy(false),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx: context.Background(),
				settings: &SetOrganizationSettings{
					OrganizationID:              "org1",
					OrganizationScopedUsernames: boolPtr(true),
				},
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "settings set, changed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterPreOrganizationSettings("org1", true, true, false),
					expectFilterOrgDomainPolicy(false),
					expectFilterOrganizationScopedUsernames(false, "username1", "username2", "username3"),
					expectPush(
						settings.NewOrganizationSettingsAddedEvent(context.Background(),
							&settings.NewAggregate("org1", "org1").Aggregate,
							[]string{"username1", "username2", "username3"},
							true,
							false,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				settings: &SetOrganizationSettings{
					OrganizationID:              "org1",
					OrganizationScopedUsernames: boolPtr(true),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "org1",
				},
			},
		},
		{
			name: "settings not set, not existing",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterPreOrganizationSettings("org1", true, false, false),
					expectFilterOrgDomainPolicy(false),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				settings: &SetOrganizationSettings{
					OrganizationID:              "org1",
					OrganizationScopedUsernames: boolPtr(false),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "org1",
				},
			},
		},
		{
			name: "settings set, changed, usernameMustBeDomain set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterPreOrganizationSettings("org1", true, true, false),
					expectFilterOrgDomainPolicy(true),
					expectPush(
						settings.NewOrganizationSettingsAddedEvent(context.Background(),
							&settings.NewAggregate("org1", "org1").Aggregate,
							[]string{"username1", "username2", "username3"},
							true,
							true,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				settings: &SetOrganizationSettings{
					OrganizationID:              "org1",
					OrganizationScopedUsernames: boolPtr(true),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.SetOrganizationSettings(tt.args.ctx, tt.args.settings)
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

func TestCommandSide_DeleteSettingsOrganization(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx   context.Context
		orgID string
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
			name: "org id missing, invalid argument error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "settings delete, no change",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "org1",
				},
			},
		},
		{
			name: "settings delete, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterOrganizationSettings("org1", true, true),
					expectFilterOrgDomainPolicy(false),
					expectFilterOrganizationScopedUsernames(false, "username1", "username2", "username3"),
					expectPush(
						settings.NewOrganizationSettingsRemovedEvent(context.Background(),
							&settings.NewAggregate("org1", "org1").Aggregate,
							[]string{"username1", "username2", "username3"},
							false,
							true,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "org1",
				},
			},
		},
		{
			name: "settings delete, unset, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterOrganizationSettings("org1", true, false),
					expectFilterOrgDomainPolicy(false),
					expectPush(
						settings.NewOrganizationSettingsRemovedEvent(context.Background(),
							&settings.NewAggregate("org1", "org1").Aggregate,
							[]string{"username1", "username2", "username3"},
							false,
							false,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "org1",
				},
			},
		},
		{
			name: "settings delete, unset, usernameMustBeDomain set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterOrganizationSettings("org1", true, false),
					expectFilterOrgDomainPolicy(true),
					expectFilterOrganizationScopedUsernames(true, "username1", "username2", "username3"),
					expectPush(
						settings.NewOrganizationSettingsRemovedEvent(context.Background(),
							&settings.NewAggregate("org1", "org1").Aggregate,
							[]string{"username1", "username2", "username3"},
							true,
							true,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "org1",
				},
			},
		},
		{
			name: "settings delete, set, usernameMustBeDomain set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterOrganizationSettings("org1", true, true),
					expectFilterOrgDomainPolicy(true),
					expectPush(
						settings.NewOrganizationSettingsRemovedEvent(context.Background(),
							&settings.NewAggregate("org1", "org1").Aggregate,
							[]string{"username1", "username2", "username3"},
							true,
							true,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "org1",
				},
			},
		},
		{
			name: "settings delete, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterOrganizationSettings("org1", true, true),
					expectFilterOrgDomainPolicy(true),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.DeleteOrganizationSettings(tt.args.ctx, tt.args.orgID)
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

func expectFilterPreOrganizationSettings(orgID string, orgExisting, settingExisting, orgScopedUsernames bool) expect {
	var events []eventstore.Event
	events = append(events,
		expectFilterPreOrganizationSettingsEvents(context.Background(), orgID, orgExisting)...,
	)
	events = append(events,
		expectFilterOrganizationSettingsEvents(context.Background(), orgID, settingExisting, orgScopedUsernames)...,
	)
	return expectFilter(
		events...,
	)
}

func expectFilterPreOrganizationSettingsEvents(ctx context.Context, orgID string, orgExisting bool) []eventstore.Event {
	var events []eventstore.Event
	if orgExisting {
		events = append(events,
			eventFromEventPusher(
				org.NewOrgAddedEvent(ctx,
					&org.NewAggregate(orgID).Aggregate,
					"org",
				),
			),
		)
	}
	return events
}

func expectFilterOrganizationSettings(orgID string, settingExisting, orgScopedUsernames bool) expect {
	return expectFilter(
		expectFilterOrganizationSettingsEvents(context.Background(), orgID, settingExisting, orgScopedUsernames)...,
	)
}

func expectFilterOrganizationSettingsEvents(ctx context.Context, orgID string, settingExisting, orgScopedUsernames bool) []eventstore.Event {
	var events []eventstore.Event
	if settingExisting {
		events = append(events,
			eventFromEventPusher(
				settings.NewOrganizationSettingsAddedEvent(ctx,
					&settings.NewAggregate(orgID, orgID).Aggregate,
					[]string{},
					orgScopedUsernames,
					!orgScopedUsernames,
				),
			),
		)
	}
	return events
}

func expectFilterOrgDomainPolicy(userLoginMustBeDomain bool) expect {
	return expectFilter(
		eventFromEventPusher(
			org.NewDomainPolicyAddedEvent(context.Background(),
				&org.NewAggregate("org1").Aggregate,
				userLoginMustBeDomain, false, false,
			),
		),
	)
}

func expectFilterOrganizationScopedUsernames(userMustBeDomain bool, usernames ...string) expect {
	events := make([]eventstore.Event, len(usernames))
	for i, username := range usernames {
		events[i] = eventFromEventPusher(
			user.NewHumanAddedEvent(context.Background(),
				&user.NewAggregate(username, "org1").Aggregate,
				username,
				"firstname",
				"lastname",
				"nickname",
				"displayname",
				language.German,
				domain.GenderUnspecified,
				"email@test.ch",
				userMustBeDomain,
			),
		)
	}
	return expectFilter(
		events...,
	)
}
