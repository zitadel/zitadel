package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddProjectGrant(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		idGenerator     id.Generator
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx          context.Context
		projectGrant *AddProjectGrant
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
			name: "invalid usergrant, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &AddProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID: "grant1",
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
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &AddProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:      "grant1",
					GrantedOrgID: "grantedorg1",
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project not existing in org, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "otherorg").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1").Aggregate,
								"granted org",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "otherorg").Aggregate,
								"key1",
								"key",
								"",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &AddProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:      "grant1",
					GrantedOrgID: "grantedorg1",
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "granted org not existing, precondition error",
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
				projectGrant: &AddProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:      "grant1",
					GrantedOrgID: "grantedorg1",
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project roles not existing, precondition error",
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
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1").Aggregate,
								"granted org",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &AddProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:      "grant1",
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1"},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "grant for project, same resourceowner",
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
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"granted org",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "projectgrant1"),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &AddProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantedOrgID: "org1",
					RoleKeys:     []string{"key1"},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "grant for project, ok",
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
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1").Aggregate,
								"granted org",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"",
							),
						),
					),
					expectPush(
						project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "projectgrant1"),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &AddProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1"},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "grant for project, id, ok",
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
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1").Aggregate,
								"granted org",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"",
							),
						),
					),
					expectPush(
						project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &AddProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:      "projectgrant1",
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1"},
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
				idGenerator:     tt.fields.idGenerator,
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.AddProjectGrant(tt.args.ctx, tt.args.projectGrant)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.NotEmpty(t, got.ID)
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ChangeProjectGrant(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx                 context.Context
		projectGrant        *ChangeProjectGrant
		cascadeUserGrantIDs []string
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
			name: "invalid projectgrant, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &ChangeProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "projectgrant not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &ChangeProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:  "projectgrant1",
					RoleKeys: []string{"key1"},
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "project not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &ChangeProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID: "projectgrant1",
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project not existing in org, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "otherorg").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1").Aggregate,
								"granted org",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &ChangeProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:  "projectgrant1",
					RoleKeys: []string{"key1"},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "granted org not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
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
				projectGrant: &ChangeProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID: "projectgrant1",
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project roles not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1").Aggregate,
								"granted org",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &ChangeProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:  "projectgrant1",
					RoleKeys: []string{"key1"},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "projectgrant not changed, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1").Aggregate,
								"granted org",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &ChangeProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:  "projectgrant1",
					RoleKeys: []string{"key1"},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant only added roles, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1").Aggregate,
								"granted org",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key2",
								"key2",
								"",
							),
						),
					),
					expectPush(
						project.NewGrantChangedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							[]string{"key1", "key2"},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &ChangeProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:  "projectgrant1",
					RoleKeys: []string{"key1", "key2"},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant only added roles, grantedOrgID, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1").Aggregate,
								"granted org",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key2",
								"key2",
								"",
							),
						),
					),
					expectPush(
						project.NewGrantChangedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							[]string{"key1", "key2"},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &ChangeProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantedOrgID: "grantedorg1",
					RoleKeys:     []string{"key1", "key2"},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant remove roles, usergrant not found, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1", "key2"},
						)),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1").Aggregate,
								"granted org",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key2",
								"key2",
								"",
							),
						),
					),
					expectFilter(),
					expectPush(
						project.NewGrantChangedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							[]string{"key1"},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &ChangeProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:  "projectgrant1",
					RoleKeys: []string{"key1"},
				},
				cascadeUserGrantIDs: []string{"usergrant1"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant remove roles, usergrant not found, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1", "key2"},
						)),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("grantedorg1").Aggregate,
								"granted org",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key2",
								"key2",
								"",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"projectgrant1",
								[]string{"key1", "key2"}),
						),
					),
					expectPush(
						project.NewGrantChangedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							[]string{"key1"},
						),
						usergrant.NewUserGrantCascadeChangedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							[]string{"key1"},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				projectGrant: &ChangeProjectGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					GrantID:  "projectgrant1",
					RoleKeys: []string{"key1"},
				},
				cascadeUserGrantIDs: []string{"usergrant1"},
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
			got, err := r.ChangeProjectGrant(tt.args.ctx, tt.args.projectGrant, tt.args.cascadeUserGrantIDs...)
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

func TestCommandSide_DeactivateProjectGrant(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
		projectID     string
		grantID       string
		grantedOrgID  string
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
			name: "missing projectid, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing grantid, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, precondition failed error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "projectgrant not existing, precondition error",
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
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "projectgrant already deactivated, ok",
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
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
						eventFromEventPusher(project.NewGrantDeactivateEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
						)),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant deactivate, ok",
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
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectPush(
						project.NewGrantDeactivateEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant deactivate, grantedOrgID, ok",
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
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectPush(
						project.NewGrantDeactivateEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantedOrgID:  "grantedorg1",
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
			got, err := r.DeactivateProjectGrant(tt.args.ctx, tt.args.projectID, tt.args.grantID, tt.args.grantedOrgID, tt.args.resourceOwner)
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

func TestCommandSide_ReactivateProjectGrant(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
		projectID     string
		grantID       string
		grantedOrgID  string
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
			name: "missing projectid, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing grantid, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, precondition failed error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "projectgrant not existing, precondition error",
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
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "projectgrant not inactive, precondition error",
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
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant reactivate, ok",
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
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
						eventFromEventPusher(project.NewGrantDeactivateEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
						)),
					),
					expectPush(
						project.NewGrantReactivatedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant reactivate, grantedOrgID, ok",
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
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
						eventFromEventPusher(project.NewGrantDeactivateEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
						)),
					),
					expectPush(
						project.NewGrantReactivatedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantedOrgID:  "grantedorg1",
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
			got, err := r.ReactivateProjectGrant(tt.args.ctx, tt.args.projectID, tt.args.grantID, tt.args.grantedOrgID, tt.args.resourceOwner)
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

func TestCommandSide_RemoveProjectGrant(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx                 context.Context
		projectID           string
		grantID             string
		resourceOwner       string
		cascadeUserGrantIDs []string
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
			name: "missing projectid, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing grantid, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project already removed, precondition failed error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
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
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "projectgrant not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "projectgrant remove, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectPush(
						project.NewGrantRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant remove, cascading usergrant not found, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(),
					expectPush(
						project.NewGrantRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:                 context.Background(),
				projectID:           "project1",
				grantID:             "projectgrant1",
				resourceOwner:       "org1",
				cascadeUserGrantIDs: []string{"usergrant1"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant remove with cascading usergrants, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(
						eventFromEventPusher(usergrant.NewUserGrantAddedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							"user1",
							"project1",
							"projectgrant1",
							[]string{"key1"}))),
					expectPush(
						project.NewGrantRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
						),
						usergrant.NewUserGrantCascadeRemovedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							"user1",
							"project1",
							"projectgrant1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:                 context.Background(),
				projectID:           "project1",
				grantID:             "projectgrant1",
				resourceOwner:       "org1",
				cascadeUserGrantIDs: []string{"usergrant1"},
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
			got, err := r.RemoveProjectGrant(tt.args.ctx, tt.args.projectID, tt.args.grantID, tt.args.resourceOwner, tt.args.cascadeUserGrantIDs...)
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

func TestCommandSide_DeleteProjectGrant(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx                 context.Context
		projectID           string
		grantID             string
		grantedOrgID        string
		resourceOwner       string
		cascadeUserGrantIDs []string
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
			name: "missing projectid, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing grantid, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project already removed, precondition failed error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
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
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant remove, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectPush(
						project.NewGrantRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantID:       "projectgrant1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant remove, grantedOrgID, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectPush(
						project.NewGrantRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				grantedOrgID:  "grantedorg1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant remove, cascading usergrant not found, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(),
					expectPush(
						project.NewGrantRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:                 context.Background(),
				projectID:           "project1",
				grantID:             "projectgrant1",
				resourceOwner:       "org1",
				cascadeUserGrantIDs: []string{"usergrant1"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "projectgrant remove with cascading usergrants, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
							[]string{"key1"},
						)),
					),
					expectFilter(
						eventFromEventPusher(usergrant.NewUserGrantAddedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							"user1",
							"project1",
							"projectgrant1",
							[]string{"key1"}))),
					expectPush(
						project.NewGrantRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"projectgrant1",
							"grantedorg1",
						),
						usergrant.NewUserGrantCascadeRemovedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							"user1",
							"project1",
							"projectgrant1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:                 context.Background(),
				projectID:           "project1",
				grantID:             "projectgrant1",
				resourceOwner:       "org1",
				cascadeUserGrantIDs: []string{"usergrant1"},
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
			got, err := r.DeleteProjectGrant(tt.args.ctx, tt.args.projectID, tt.args.grantID, tt.args.grantedOrgID, tt.args.resourceOwner, tt.args.cascadeUserGrantIDs...)
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
