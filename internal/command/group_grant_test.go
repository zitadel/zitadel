package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/groupgrant"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func addNewGroupGrantPreConditionEvents(groupID, projectID, orgID string, roleKeys []string) []eventstore.Event {
	events := []eventstore.Event{
		eventFromEventPusher(addNewGroupEvent(groupID, orgID)),
		eventFromEventPusher(
			project.NewProjectAddedEvent(context.Background(),
				&project.NewAggregate(projectID, orgID).Aggregate,
				"project",
				false,
				false,
				false,
				domain.PrivateLabelingSettingUnspecified,
			),
		),
	}
	for _, roleKey := range roleKeys {
		events = append(events, eventFromEventPusher(
			project.NewRoleAddedEvent(context.Background(),
				&project.NewAggregate(projectID, orgID).Aggregate,
				roleKey,
				roleKey,
				"",
			),
		))
	}
	return events
}

func addNewGroupGrantCrossOrgPreConditionEvents(groupID, projectID, projectOrgID, groupOrgID, grantID string, grantRoleKeys []string) []eventstore.Event {
	events := []eventstore.Event{
		eventFromEventPusher(addNewGroupEvent(groupID, groupOrgID)),
		eventFromEventPusher(
			project.NewProjectAddedEvent(context.Background(),
				&project.NewAggregate(projectID, projectOrgID).Aggregate,
				"project",
				false,
				false,
				false,
				domain.PrivateLabelingSettingUnspecified,
			),
		),
	}
	for _, roleKey := range grantRoleKeys {
		events = append(events, eventFromEventPusher(
			project.NewRoleAddedEvent(context.Background(),
				&project.NewAggregate(projectID, projectOrgID).Aggregate,
				roleKey,
				roleKey,
				"",
			),
		))
	}
	events = append(events, eventFromEventPusher(
		project.NewGrantAddedEvent(context.Background(),
			&project.NewAggregate(projectID, projectOrgID).Aggregate,
			grantID,
			groupOrgID,
			grantRoleKeys,
		),
	))
	return events
}

func addNewGroupGrantAddedEvent(grantID, groupID, projectID, orgID string, roleKeys []string) *groupgrant.GroupGrantAddedEvent {
	return groupgrant.NewGroupGrantAddedEvent(context.Background(),
		&groupgrant.NewAggregate(grantID, orgID).Aggregate,
		groupID,
		projectID,
		"",
		roleKeys,
	)
}

func TestCommands_AddGroupGrant(t *testing.T) {
	t.Parallel()

	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
		idGenerator     id.Generator
	}
	type args struct {
		ctx   context.Context
		grant *AddGroupGrant
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr func(error) bool
	}{
		{
			name: "missing role keys, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				grant: &AddGroupGrant{
					GroupID:   "group1",
					ProjectID: "project1",
				},
			},
			wantErr: zerrors.IsErrorInvalidArgument,
		},
		{
			name: "group not found, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				grant: &AddGroupGrant{
					GroupID:   "group1",
					ProjectID: "project1",
					RoleKeys:  []string{"role1"},
				},
			},
			wantErr: zerrors.IsPreconditionFailed,
		},
		{
			name: "missing permission, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(addNewGroupEvent("group1", "org1")),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx: context.Background(),
				grant: &AddGroupGrant{
					GroupID:   "group1",
					ProjectID: "project1",
					RoleKeys:  []string{"role1"},
				},
			},
			wantErr: zerrors.IsPermissionDenied,
		},
		{
			name: "role does not exist on project, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(addNewGroupEvent("group1", "org1")),
					),
					expectFilter(
						addNewGroupGrantPreConditionEvents("group1", "project1", "org1", []string{"role1"})...,
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				grant: &AddGroupGrant{
					GroupID:   "group1",
					ProjectID: "project1",
					RoleKeys:  []string{"role1", "missing-role"},
				},
			},
			wantErr: zerrors.IsPreconditionFailed,
		},
		{
			name: "group grant added, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(addNewGroupEvent("group1", "org1")),
					),
					expectFilter(
						addNewGroupGrantPreConditionEvents("group1", "project1", "org1", []string{"role1", "role2"})...,
					),
					expectPush(
						addNewGroupGrantAddedEvent("grant1", "group1", "project1", "org1", []string{"role1", "role2"}),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "grant1"),
			},
			args: args{
				ctx: context.Background(),
				grant: &AddGroupGrant{
					GroupID:   "group1",
					ProjectID: "project1",
					RoleKeys:  []string{"role1", "role2"},
				},
			},
			want: &domain.ObjectDetails{
				ID:            "grant1",
				ResourceOwner: "org1",
			},
		},
		{
			name: "cross-org project grant added, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(addNewGroupEvent("group1", "org1")),
					),
					expectFilter(
						addNewGroupGrantCrossOrgPreConditionEvents("group1", "project1", "org2", "org1", "projectgrant1", []string{"role1"})...,
					),
					expectPush(
						groupgrant.NewGroupGrantAddedEvent(context.Background(),
							&groupgrant.NewAggregate("grant1", "org1").Aggregate,
							"group1",
							"project1",
							"projectgrant1",
							[]string{"role1"},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "grant1"),
			},
			args: args{
				ctx: context.Background(),
				grant: &AddGroupGrant{
					GroupID:        "group1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"role1"},
				},
			},
			want: &domain.ObjectDetails{
				ID:            "grant1",
				ResourceOwner: "org1",
			},
		},
		{
			name: "cross-org project without grant, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(addNewGroupEvent("group1", "org1")),
					),
					expectFilter(
						eventFromEventPusher(addNewGroupEvent("group1", "org1")),
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org2").Aggregate,
								"project",
								false,
								false,
								false,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org2").Aggregate,
								"role1",
								"role1",
								"",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				grant: &AddGroupGrant{
					GroupID:        "group1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"role1"},
				},
			},
			wantErr: zerrors.IsPreconditionFailed,
		},
		{
			name: "cross-org grant missing requested role, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(addNewGroupEvent("group1", "org1")),
					),
					expectFilter(
						addNewGroupGrantCrossOrgPreConditionEvents("group1", "project1", "org2", "org1", "projectgrant1", []string{"role1"})...,
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				grant: &AddGroupGrant{
					GroupID:        "group1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"role2"},
				},
			},
			wantErr: zerrors.IsPreconditionFailed,
		},
		{
			name: "cross-org grant changed roles, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(addNewGroupEvent("group1", "org1")),
					),
					expectFilter(
						append(
							addNewGroupGrantCrossOrgPreConditionEvents("group1", "project1", "org2", "org1", "projectgrant1", []string{"role1"}),
							eventFromEventPusher(
								project.NewRoleAddedEvent(context.Background(),
									&project.NewAggregate("project1", "org2").Aggregate,
									"role2",
									"role2",
									"",
								),
							),
							eventFromEventPusher(
								project.NewGrantChangedEvent(context.Background(),
									&project.NewAggregate("project1", "org2").Aggregate,
									"projectgrant1",
									[]string{"role1", "role2"},
								),
							),
						)...,
					),
					expectPush(
						groupgrant.NewGroupGrantAddedEvent(context.Background(),
							&groupgrant.NewAggregate("grant1", "org1").Aggregate,
							"group1",
							"project1",
							"projectgrant1",
							[]string{"role2"},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "grant1"),
			},
			args: args{
				ctx: context.Background(),
				grant: &AddGroupGrant{
					GroupID:        "group1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"role2"},
				},
			},
			want: &domain.ObjectDetails{
				ID:            "grant1",
				ResourceOwner: "org1",
			},
		},
		{
			name: "cross-org grant removed, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(addNewGroupEvent("group1", "org1")),
					),
					expectFilter(
						append(
							addNewGroupGrantCrossOrgPreConditionEvents("group1", "project1", "org2", "org1", "projectgrant1", []string{"role1"}),
							eventFromEventPusher(
								project.NewGrantRemovedEvent(context.Background(),
									&project.NewAggregate("project1", "org2").Aggregate,
									"projectgrant1",
									"org1",
								),
							),
						)...,
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				grant: &AddGroupGrant{
					GroupID:        "group1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"role1"},
				},
			},
			wantErr: zerrors.IsPreconditionFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
				idGenerator:     tt.fields.idGenerator,
			}
			got, err := c.AddGroupGrant(tt.args.ctx, tt.args.grant)
			if tt.wantErr != nil {
				require.True(t, tt.wantErr(err))
				return
			}
			require.NoError(t, err)
			assertObjectDetails(t, tt.want, got)
		})
	}
}

func TestCommands_ChangeGroupGrant(t *testing.T) {
	t.Parallel()

	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx      context.Context
		grantID  string
		roleKeys []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr func(error) bool
	}{
		{
			name: "grant not found, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:      context.Background(),
				grantID:  "grant1",
				roleKeys: []string{"role1"},
			},
			wantErr: zerrors.IsNotFound,
		},
		{
			name: "missing permission, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupGrantAddedEvent("grant1", "group1", "project1", "org1", []string{"role1"}),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:      context.Background(),
				grantID:  "grant1",
				roleKeys: []string{"role1", "role2"},
			},
			wantErr: zerrors.IsPermissionDenied,
		},
		{
			name: "roles unchanged, no events pushed, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupGrantAddedEvent("grant1", "group1", "project1", "org1", []string{"role1"}),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:      context.Background(),
				grantID:  "grant1",
				roleKeys: []string{"role1"},
			},
			want: &domain.ObjectDetails{
				ID:            "grant1",
				ResourceOwner: "org1",
			},
		},
		{
			name: "roles changed, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupGrantAddedEvent("grant1", "group1", "project1", "org1", []string{"role1"}),
						),
					),
					expectFilter(
						addNewGroupGrantPreConditionEvents("group1", "project1", "org1", []string{"role1", "role2"})...,
					),
					expectPush(
						groupgrant.NewGroupGrantChangedEvent(context.Background(),
							&groupgrant.NewAggregate("grant1", "org1").Aggregate,
							[]string{"role1", "role2"},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:      context.Background(),
				grantID:  "grant1",
				roleKeys: []string{"role1", "role2"},
			},
			want: &domain.ObjectDetails{
				ID:            "grant1",
				ResourceOwner: "org1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := c.ChangeGroupGrant(tt.args.ctx, tt.args.grantID, tt.args.roleKeys)
			if tt.wantErr != nil {
				require.True(t, tt.wantErr(err))
				return
			}
			require.NoError(t, err)
			assertObjectDetails(t, tt.want, got)
		})
	}
}

func TestCommands_RemoveGroupGrant(t *testing.T) {
	t.Parallel()

	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx     context.Context
		grantID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr func(error) bool
	}{
		{
			name: "grant not found, desired state achieved, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:     context.Background(),
				grantID: "grant1",
			},
			want: &domain.ObjectDetails{
				ID: "grant1",
			},
		},
		{
			name: "missing permission, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupGrantAddedEvent("grant1", "group1", "project1", "org1", []string{"role1"}),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:     context.Background(),
				grantID: "grant1",
			},
			wantErr: zerrors.IsPermissionDenied,
		},
		{
			name: "group grant removed, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupGrantAddedEvent("grant1", "group1", "project1", "org1", []string{"role1"}),
						),
					),
					expectPush(
						groupgrant.NewGroupGrantRemovedEvent(context.Background(),
							&groupgrant.NewAggregate("grant1", "org1").Aggregate,
							"group1",
							"project1",
							"",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:     context.Background(),
				grantID: "grant1",
			},
			want: &domain.ObjectDetails{
				ID:            "grant1",
				ResourceOwner: "org1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := c.RemoveGroupGrant(tt.args.ctx, tt.args.grantID)
			if tt.wantErr != nil {
				require.True(t, tt.wantErr(err))
				return
			}
			require.NoError(t, err)
			assertObjectDetails(t, tt.want, got)
		})
	}
}
