package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/settings"
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
					OrganizationID: "",
					UserUniqueness: boolPtr(true),
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
					OrganizationID: "org1",
					UserUniqueness: boolPtr(true),
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
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org",
							),
						),
						eventFromEventPusher(
							settings.NewOrganizationSettingsAddedEvent(context.Background(),
								&settings.NewAggregate("org1", "org1").Aggregate,
								true,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				settings: &SetOrganizationSettings{
					OrganizationID: "org1",
					UserUniqueness: boolPtr(true),
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
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org",
							),
						),
					),
					expectPush(
						settings.NewOrganizationSettingsAddedEvent(context.Background(),
							&settings.NewAggregate("org1", "org1").Aggregate,
							true,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				settings: &SetOrganizationSettings{
					OrganizationID: "org1",
					UserUniqueness: boolPtr(true),
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
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx: context.Background(),
				settings: &SetOrganizationSettings{
					OrganizationID: "org1",
					UserUniqueness: boolPtr(true),
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
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org",
							),
						),
						eventFromEventPusher(
							settings.NewOrganizationSettingsAddedEvent(context.Background(),
								&settings.NewAggregate("org1", "org1").Aggregate,
								false,
							),
						),
					),
					expectPush(
						settings.NewOrganizationSettingsAddedEvent(context.Background(),
							&settings.NewAggregate("org1", "org1").Aggregate,
							true,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				settings: &SetOrganizationSettings{
					OrganizationID: "org1",
					UserUniqueness: boolPtr(true),
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
			name: "settings not set, not changed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				settings: &SetOrganizationSettings{
					OrganizationID: "org1",
					UserUniqueness: boolPtr(false),
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
					expectFilter(
						eventFromEventPusher(
							settings.NewOrganizationSettingsAddedEvent(context.Background(),
								&settings.NewAggregate("org1", "org1").Aggregate,
								true,
							),
						),
					),
					expectPush(
						settings.NewOrganizationSettingsRemovedEvent(context.Background(),
							&settings.NewAggregate("org1", "org1").Aggregate,
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
					expectFilter(
						eventFromEventPusher(
							settings.NewOrganizationSettingsAddedEvent(context.Background(),
								&settings.NewAggregate("org1", "org1").Aggregate,
								true,
							),
						),
					),
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
