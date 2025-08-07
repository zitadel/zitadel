package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddProjectRole(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx  context.Context
		role *AddProjectRole
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
			name: "project not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
								nil,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				role: &AddProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Key: "key1",
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid role, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				role: &AddProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "role key already exists, already exists error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectPushFailed(zerrors.ThrowAlreadyExists(nil, "id", "internal"),
						project.NewRoleAddedEvent(
							context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"key1",
							"key",
							"group",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				role: &AddProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Key:         "key1",
					DisplayName: "key",
					Group:       "group",
				},
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "add role,ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectPush(
						project.NewRoleAddedEvent(
							context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"key1",
							"key",
							"group",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				role: &AddProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Key:         "key1",
					DisplayName: "key",
					Group:       "group",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "add role, resourceowner, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectPush(
						project.NewRoleAddedEvent(
							context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"key1",
							"key",
							"group",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				role: &AddProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					Key:         "key1",
					DisplayName: "key",
					Group:       "group",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
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
			got, err := r.AddProjectRole(tt.args.ctx, tt.args.role)
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

func TestCommandSide_BulkAddProjectRole(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
		roles         []*AddProjectRole
		projectID     string
		resourceOwner string
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
			name: "project not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
								nil,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				roles: []*AddProjectRole{
					{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "project1",
						},
						Key: "key1",
					},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid role, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				roles: []*AddProjectRole{
					{
						ObjectRoot: models.ObjectRoot{},
					},
					{
						ObjectRoot: models.ObjectRoot{},
					},
				},
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "role key already exists, already exists error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectPushFailed(zerrors.ThrowAlreadyExists(nil, "id", "internal"),
						project.NewRoleAddedEvent(
							context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"key1",
							"key",
							"group",
						),
						project.NewRoleAddedEvent(
							context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"key2",
							"key2",
							"group",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				roles: []*AddProjectRole{
					{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "project1",
						},
						Key:         "key1",
						DisplayName: "key",
						Group:       "group",
					},
					{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "project1",
						},
						Key:         "key2",
						DisplayName: "key2",
						Group:       "group",
					},
				},
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "add roles,ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectPush(
						project.NewRoleAddedEvent(
							context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"key1",
							"key",
							"group",
						),
						project.NewRoleAddedEvent(
							context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"key2",
							"key2",
							"group",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				roles: []*AddProjectRole{
					{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "project1",
						},
						Key:         "key1",
						DisplayName: "key",
						Group:       "group",
					},
					{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "project1",
						},
						Key:         "key2",
						DisplayName: "key2",
						Group:       "group",
					},
				},
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
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
			got, err := r.BulkAddProjectRole(tt.args.ctx, tt.args.projectID, tt.args.resourceOwner, tt.args.roles)
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

func TestCommandSide_ChangeProjectRole(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx  context.Context
		role *ChangeProjectRole
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
			name: "invalid role, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				role: &ChangeProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
								nil,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				role: &ChangeProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Key: "key1",
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "role removed, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"group",
							),
						),
						eventFromEventPusher(
							project.NewRoleRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				role: &ChangeProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Key:         "key1",
					DisplayName: "key",
					Group:       "group",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "role not changed, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"group",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				role: &ChangeProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Key:         "key1",
					DisplayName: "key",
					Group:       "group",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "role changed, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"group",
							),
						),
					),
					expectPush(
						newRoleChangedEvent(context.Background(), "project1", "org1", "key1", "keychanged", "groupchanged"),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				role: &ChangeProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Key:         "key1",
					DisplayName: "keychanged",
					Group:       "groupchanged",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
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
			got, err := r.ChangeProjectRole(tt.args.ctx, tt.args.role)
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

func TestCommandSide_RemoveProjectRole(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx                      context.Context
		projectID                string
		key                      string
		resourceOwner            string
		cascadingProjectGrantIDs []string
		cascadingUserGrantIDs    []string
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
			name: "invalid projectid, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				key:           "key1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid key, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				key:           "",
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "role not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				key:           "key",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "role removed, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"ggroup",
							),
						),
						eventFromEventPusher(
							project.NewRoleRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				key:           "key",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "role removed, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"group",
							),
						),
					),
					expectPush(
						project.NewRoleRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"key1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				key:           "key1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "role removed with cascadingProjectGrantids, grant not found, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"group",
							),
						),
					),
					expectFilter(),
					expectPush(
						project.NewRoleRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"key1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:                      context.Background(),
				projectID:                "project1",
				key:                      "key1",
				resourceOwner:            "org1",
				cascadingProjectGrantIDs: []string{"projectgrant1"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "role removed with cascadingProjectGrantids, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"group",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"org2",
								[]string{"key1"},
							),
						),
					),
					expectPush(
						project.NewRoleRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"key1",
						),
						project.NewGrantCascadeChangedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							[]string{},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:                      context.Background(),
				projectID:                "project1",
				key:                      "key1",
				resourceOwner:            "org1",
				cascadingProjectGrantIDs: []string{"projectgrant1"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "role removed with cascadingUserGrantIDs, grant not found, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"group",
							),
						),
					),
					expectFilter(),
					expectPush(
						project.NewRoleRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"key1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:                   context.Background(),
				projectID:             "project1",
				key:                   "key1",
				resourceOwner:         "org1",
				cascadingUserGrantIDs: []string{"usergrant1"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "role removed with cascadingUserGrantIDs, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"group",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"",
								[]string{"key1"})),
					),
					expectPush(
						project.NewRoleRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"key1",
						),
						usergrant.NewUserGrantCascadeChangedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							[]string{},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:                   context.Background(),
				projectID:             "project1",
				key:                   "key1",
				resourceOwner:         "org1",
				cascadingUserGrantIDs: []string{"usergrant1"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
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
			got, err := r.RemoveProjectRole(tt.args.ctx, tt.args.projectID, tt.args.key, tt.args.resourceOwner, tt.args.cascadingProjectGrantIDs, tt.args.cascadingUserGrantIDs...)
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

func newRoleChangedEvent(ctx context.Context, projectID, resourceOwner, key, displayName, group string) *project.RoleChangedEvent {
	event, _ := project.NewRoleChangedEvent(ctx,
		&project.NewAggregate(projectID, resourceOwner).Aggregate,
		key,
		[]project.RoleChanges{
			project.ChangeDisplayName(displayName),
			project.ChangeGroup(group),
		},
	)
	return event
}
