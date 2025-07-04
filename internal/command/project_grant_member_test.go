package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddProjectGrantMember(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		zitadelRoles    []authz.RoleMapping
		checkPermission domain.PermissionCheck
	}
	type args struct {
		member *AddProjectGrantMember
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
			name: "invalid member, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				member: &AddProjectGrantMember{
					ProjectID: "project1",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid roles, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				member: &AddProjectGrantMember{
					ProjectID: "project1",
					GrantID:   "projectgrant1",
					UserID:    "user1",
					Roles:     []string{"PROJECT_GRANT_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "PROJECT_GRANT_OWNER",
					},
				},
			},
			args: args{
				member: &AddProjectGrantMember{
					ProjectID: "project1",
					GrantID:   "projectgrant1",
					UserID:    "user1",
					Roles:     []string{"PROJECT_GRANT_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "member add uniqueconstraint err, already exists",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username1",
								"firstname1",
								"lastname1",
								"nickname1",
								"displayname1",
								language.German,
								domain.GenderMale,
								"email1",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"grantedorg1",
								[]string{"key1"},
							),
							),
						),
					),
					expectFilter(),
					expectPushFailed(zerrors.ThrowAlreadyExists(nil, "ERROR", "internal"),
						project.NewProjectGrantMemberAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"user1",
							"projectgrant1",
							[]string{"PROJECT_GRANT_OWNER"}...,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "PROJECT_GRANT_OWNER",
					},
				},
			},
			args: args{
				member: &AddProjectGrantMember{
					ProjectID: "project1",
					GrantID:   "projectgrant1",
					UserID:    "user1",
					Roles:     []string{"PROJECT_GRANT_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "member add, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username1",
								"firstname1",
								"lastname1",
								"nickname1",
								"displayname1",
								language.German,
								domain.GenderMale,
								"email1",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"grantedorg1",
								[]string{"key1"},
							),
							),
						),
					),
					expectFilter(),
					expectPush(
						project.NewProjectGrantMemberAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"user1",
							"projectgrant1",
							[]string{"PROJECT_GRANT_OWNER"}...,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "PROJECT_GRANT_OWNER",
					},
				},
			},
			args: args{
				member: &AddProjectGrantMember{
					ProjectID: "project1",
					GrantID:   "projectgrant1",
					UserID:    "user1",
					Roles:     []string{"PROJECT_GRANT_OWNER"},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "project1",
				},
			},
		},
		{
			name: "member add, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username1",
								"firstname1",
								"lastname1",
								"nickname1",
								"displayname1",
								language.German,
								domain.GenderMale,
								"email1",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"grantedorg1",
								[]string{"key1"},
							),
							),
						),
					),
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "PROJECT_GRANT_OWNER",
					},
				},
			},
			args: args{
				member: &AddProjectGrantMember{
					ProjectID: "project1",
					GrantID:   "projectgrant1",
					UserID:    "user1",
					Roles:     []string{"PROJECT_GRANT_OWNER"},
				},
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
				zitadelRoles:    tt.fields.zitadelRoles,
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.AddProjectGrantMember(context.Background(), tt.args.member)
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

func TestCommandSide_ChangeProjectGrantMember(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		zitadelRoles    []authz.RoleMapping
		checkPermission domain.PermissionCheck
	}
	type args struct {
		member *ChangeProjectGrantMember
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
			name: "invalid member, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				member: &ChangeProjectGrantMember{
					ProjectID: "project1",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid roles, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				member: &ChangeProjectGrantMember{
					ProjectID: "project1",
					GrantID:   "projectgrant1",
					UserID:    "user1",
					Roles:     []string{"PROJECT_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "member not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"org2",
								[]string{"rol1", "role2"},
							),
						),
					),
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "PROJECT_GRANT_OWNER",
					},
				},
			},
			args: args{
				member: &ChangeProjectGrantMember{
					ProjectID: "project1",
					GrantID:   "projectgrant1",
					UserID:    "user1",
					Roles:     []string{"PROJECT_GRANT_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "member not changed, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"org2",
								[]string{"rol1", "role2"},
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectGrantMemberAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"user1",
								"projectgrant1",
								[]string{"PROJECT_GRANT_OWNER"}...,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "PROJECT_GRANT_OWNER",
					},
				},
			},
			args: args{
				member: &ChangeProjectGrantMember{
					ProjectID: "project1",
					GrantID:   "projectgrant1",
					UserID:    "user1",
					Roles:     []string{"PROJECT_GRANT_OWNER"},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "member change, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"org2",
								[]string{"rol1", "role2"},
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectGrantMemberAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"user1",
								"projectgrant1",
								[]string{"PROJECT_GRANT_OWNER"}...,
							),
						),
					),
					expectPush(
						project.NewProjectGrantMemberChangedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"user1",
							"projectgrant1",
							[]string{"PROJECT_GRANT_OWNER", "PROJECT_GRANT_VIEWER"}...,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "PROJECT_GRANT_OWNER",
					},
					{
						Role: "PROJECT_GRANT_VIEWER",
					},
				},
			},
			args: args{
				member: &ChangeProjectGrantMember{
					ProjectID: "project1",
					GrantID:   "projectgrant1",
					UserID:    "user1",
					Roles:     []string{"PROJECT_GRANT_OWNER", "PROJECT_GRANT_VIEWER"},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "member change, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"org2",
								[]string{"rol1", "role2"},
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectGrantMemberAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"user1",
								"projectgrant1",
								[]string{"PROJECT_GRANT_OWNER"}...,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "PROJECT_GRANT_OWNER",
					},
					{
						Role: "PROJECT_GRANT_VIEWER",
					},
				},
			},
			args: args{
				member: &ChangeProjectGrantMember{
					ProjectID: "project1",
					GrantID:   "projectgrant1",
					UserID:    "user1",
					Roles:     []string{"PROJECT_GRANT_OWNER", "PROJECT_GRANT_VIEWER"},
				},
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
				zitadelRoles:    tt.fields.zitadelRoles,
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.ChangeProjectGrantMember(context.Background(), tt.args.member)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_RemoveProjectGrantMember(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx       context.Context
		projectID string
		grantID   string
		userID    string
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
			name: "invalid member projectid missing, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "",
				userID:    "user1",
				grantID:   "projectgrant1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid member userid missing, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
				userID:    "",
				grantID:   "projectgrant1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid member grantid missing, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
				userID:    "user1",
				grantID:   "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "member not existing, not found ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"org2",
								[]string{"rol1", "role2"},
							),
						),
					),
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
				userID:    "user1",
				grantID:   "projectgrant1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "member remove, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"org2",
								[]string{"rol1", "role2"},
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectGrantMemberAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"user1",
								"projectgrant1",
								[]string{"PROJECT_OWNER"}...,
							),
						),
					),
					expectPush(
						project.NewProjectGrantMemberRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"user1",
							"projectgrant1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
				userID:    "user1",
				grantID:   "projectgrant1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "member remove, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"org2",
								[]string{"rol1", "role2"},
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectGrantMemberAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"user1",
								"projectgrant1",
								[]string{"PROJECT_OWNER"}...,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
				userID:    "user1",
				grantID:   "projectgrant1",
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
			got, err := r.RemoveProjectGrantMember(tt.args.ctx, tt.args.projectID, tt.args.userID, tt.args.grantID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}
