package command

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddProject(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		idGenerator     id.Generator
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx     context.Context
		project *AddProject
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
			name: "invalid project, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				project: &AddProject{
					ObjectRoot: models.ObjectRoot{ResourceOwner: "org1"},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project, resourceowner empty",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				project: &AddProject{
					ObjectRoot:             models.ObjectRoot{ResourceOwner: ""},
					Name:                   "project",
					ProjectRoleAssertion:   true,
					ProjectRoleCheck:       true,
					HasProjectCheck:        true,
					PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				project: &AddProject{
					ObjectRoot:             models.ObjectRoot{AggregateID: "project1", ResourceOwner: "org1"},
					Name:                   "project",
					ProjectRoleAssertion:   true,
					ProjectRoleCheck:       true,
					HasProjectCheck:        true,
					PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
				},
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "project, already exists",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "project1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				project: &AddProject{
					ObjectRoot:             models.ObjectRoot{ResourceOwner: "org1"},
					Name:                   "project",
					ProjectRoleAssertion:   true,
					ProjectRoleCheck:       true,
					HasProjectCheck:        true,
					PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
				},
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "project, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						project.NewProjectAddedEvent(
							context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"project", true, true, true,
							domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "project1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				project: &AddProject{
					ObjectRoot:             models.ObjectRoot{ResourceOwner: "org1"},
					Name:                   "project",
					ProjectRoleAssertion:   true,
					ProjectRoleCheck:       true,
					HasProjectCheck:        true,
					PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "project, with id, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						project.NewProjectAddedEvent(
							context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"project", true, true, true,
							domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				project: &AddProject{
					ObjectRoot:             models.ObjectRoot{AggregateID: "project1", ResourceOwner: "org1"},
					Name:                   "project",
					ProjectRoleAssertion:   true,
					ProjectRoleCheck:       true,
					HasProjectCheck:        true,
					PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
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
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				idGenerator:     tt.fields.idGenerator,
				checkPermission: tt.fields.checkPermission,
			}
			c.setMilestonesCompletedForTest("instanceID")
			got, err := c.AddProject(tt.args.ctx, tt.args.project)
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

func TestCommandSide_ChangeProject(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
		project       *ChangeProject
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
			name: "invalid project, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				project: &ChangeProject{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Name: gu.Ptr(""),
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid project empty aggregateid, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				project: &ChangeProject{
					Name: gu.Ptr("project"),
				},
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				project: &ChangeProject{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Name: gu.Ptr("project change"),
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "project removed, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project",
								nil),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				project: &ChangeProject{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Name: gu.Ptr("project change"),
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				project: &ChangeProject{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Name:                   gu.Ptr("project"),
					ProjectRoleAssertion:   gu.Ptr(true),
					ProjectRoleCheck:       gu.Ptr(true),
					HasProjectCheck:        gu.Ptr(true),
					PrivateLabelingSetting: gu.Ptr(domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "no changes, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx: context.Background(),
				project: &ChangeProject{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Name:                   gu.Ptr("project"),
					ProjectRoleAssertion:   gu.Ptr(true),
					ProjectRoleCheck:       gu.Ptr(true),
					HasProjectCheck:        gu.Ptr(true),
					PrivateLabelingSetting: gu.Ptr(domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "project change with name and unique constraints, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
					expectPush(
						newProjectChangedEvent(context.Background(),
							"project1",
							"org1",
							"project",
							"project-new",
							false,
							false,
							false,
							domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				project: &ChangeProject{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Name:                   gu.Ptr("project-new"),
					ProjectRoleAssertion:   gu.Ptr(false),
					ProjectRoleCheck:       gu.Ptr(false),
					HasProjectCheck:        gu.Ptr(false),
					PrivateLabelingSetting: gu.Ptr(domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy),
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "project change without name and unique constraints, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
					expectPush(
						newProjectChangedEvent(context.Background(),
							"project1",
							"org1",
							"",
							"",
							false,
							false,
							false,
							domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx: context.Background(),
				project: &ChangeProject{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Name:                   gu.Ptr("project"),
					ProjectRoleAssertion:   gu.Ptr(false),
					ProjectRoleCheck:       gu.Ptr(false),
					HasProjectCheck:        gu.Ptr(false),
					PrivateLabelingSetting: gu.Ptr(domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy),
				},
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
			got, err := r.ChangeProject(tt.args.ctx, tt.args.project)
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

func TestCommandSide_DeactivateProject(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
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
			name: "invalid project id, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "",
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "project removed, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project",
								nil),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "project already inactive, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
						eventFromEventPusher(
							project.NewProjectDeactivatedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project deactivate,no resourceOwner, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
					expectPush(
						project.NewProjectDeactivatedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "project deactivate, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "project deactivate, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
					expectPush(
						project.NewProjectDeactivatedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
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
			got, err := r.DeactivateProject(tt.args.ctx, tt.args.projectID, tt.args.resourceOwner)
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

func TestCommandSide_ReactivateProject(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
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
			name: "invalid project id, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "",
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "project removed, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project",
								nil),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "project not inactive, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project reactivate, no resourceOwner, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
						eventFromEventPusher(
							project.NewProjectDeactivatedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate),
						),
					),
					expectPush(
						project.NewProjectReactivatedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "project reactivate, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
						eventFromEventPusher(
							project.NewProjectDeactivatedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "project reactivate, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
						eventFromEventPusher(
							project.NewProjectDeactivatedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate),
						),
					),
					expectPush(
						project.NewProjectReactivatedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
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
			got, err := r.ReactivateProject(tt.args.ctx, tt.args.projectID, tt.args.resourceOwner)
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

func TestCommandSide_RemoveProject(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
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
			name: "invalid project id, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid resourceowner, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "",
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "project removed, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project",
								nil),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "project remove, without entityConstraints, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
					// no saml application events
					expectFilter(),
					expectPush(
						project.NewProjectRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"project",
							nil),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "project remove, with entityConstraints, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
						)),
						eventFromEventPusher(
							project.NewSAMLConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"https://test.com/saml/metadata",
								[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
								"http://localhost:8080/saml/metadata",
								domain.LoginVersionUnspecified,
								"",
							),
						),
					),
					expectPush(
						project.NewProjectRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"project",
							[]*eventstore.UniqueConstraint{
								project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test.com/saml/metadata"),
							},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "project remove, with multiple entityConstraints, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
						)),
						eventFromEventPusher(
							project.NewSAMLConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"https://test1.com/saml/metadata",
								[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
								"",
								domain.LoginVersionUnspecified,
								"",
							),
						),
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app2",
							"app",
						)),
						eventFromEventPusher(
							project.NewSAMLConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app2",
								"https://test2.com/saml/metadata",
								[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
								"",
								domain.LoginVersionUnspecified,
								"",
							),
						),
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app3",
							"app",
						)),
						eventFromEventPusher(
							project.NewSAMLConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app3",
								"https://test3.com/saml/metadata",
								[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
								"",
								domain.LoginVersionUnspecified,
								"",
							),
						),
					),
					expectPush(
						project.NewProjectRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"project",
							[]*eventstore.UniqueConstraint{
								project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test1.com/saml/metadata"),
								project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test2.com/saml/metadata"),
								project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test3.com/saml/metadata"),
							},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
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
			got, err := r.RemoveProject(tt.args.ctx, tt.args.projectID, tt.args.resourceOwner)
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

func TestCommandSide_DeleteProject(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
		projectID     string
		resourceOwner string
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
			name: "invalid project id, invalid error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "",
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: nil,
			},
		},
		{
			name: "project removed, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project",
								nil),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: nil,
			},
		}, {
			name: "project remove, no resourceOwner, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
					// no saml application events
					expectFilter(),
					expectPush(
						project.NewProjectRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"project",
							nil),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "",
			},
			res: res{
				err: nil,
			},
		},
		{
			name: "project remove, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "",
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "project remove, without entityConstraints, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
					// no saml application events
					expectFilter(),
					expectPush(
						project.NewProjectRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"project",
							nil),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: nil,
			},
		},
		{
			name: "project remove, with entityConstraints, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
						)),
						eventFromEventPusher(
							project.NewSAMLConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"https://test.com/saml/metadata",
								[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
								"http://localhost:8080/saml/metadata",
								domain.LoginVersionUnspecified,
								"",
							),
						),
					),
					expectPush(
						project.NewProjectRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"project",
							[]*eventstore.UniqueConstraint{
								project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test.com/saml/metadata"),
							},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: nil,
			},
		},
		{
			name: "project remove, with multiple entityConstraints, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
						)),
						eventFromEventPusher(
							project.NewSAMLConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"https://test1.com/saml/metadata",
								[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
								"",
								domain.LoginVersionUnspecified,
								"",
							),
						),
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app2",
							"app",
						)),
						eventFromEventPusher(
							project.NewSAMLConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app2",
								"https://test2.com/saml/metadata",
								[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
								"",
								domain.LoginVersionUnspecified,
								"",
							),
						),
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app3",
							"app",
						)),
						eventFromEventPusher(
							project.NewSAMLConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app3",
								"https://test3.com/saml/metadata",
								[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
								"",
								domain.LoginVersionUnspecified,
								"",
							),
						),
					),
					expectPush(
						project.NewProjectRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"project",
							[]*eventstore.UniqueConstraint{
								project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test1.com/saml/metadata"),
								project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test2.com/saml/metadata"),
								project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test3.com/saml/metadata"),
							},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			_, err := r.DeleteProject(tt.args.ctx, tt.args.projectID, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func newProjectChangedEvent(ctx context.Context, projectID, resourceOwner, oldName, newName string, roleAssertion, roleCheck, hasProjectCheck bool, privateLabelingSetting domain.PrivateLabelingSetting) *project.ProjectChangeEvent {
	changes := []project.ProjectChanges{
		project.ChangeProjectRoleAssertion(roleAssertion),
		project.ChangeProjectRoleCheck(roleCheck),
		project.ChangeHasProjectCheck(hasProjectCheck),
		project.ChangePrivateLabelingSetting(privateLabelingSetting),
	}
	if newName != "" {
		changes = append(changes, project.ChangeName(newName))
	}
	return project.NewProjectChangeEvent(ctx,
		&project.NewAggregate(projectID, resourceOwner).Aggregate,
		oldName,
		changes,
	)
}

func TestAddProject(t *testing.T) {
	type args struct {
		a                      *project.Aggregate
		name                   string
		owner                  string
		privateLabelingSetting domain.PrivateLabelingSetting
	}

	ctx := context.Background()
	agg := project.NewAggregate("test", "test")

	tests := []struct {
		name string
		args args
		want Want
	}{
		{
			name: "invalid name",
			args: args{
				a:                      agg,
				name:                   "",
				owner:                  "owner",
				privateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
			},
			want: Want{
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "PROJE-C01yo", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid private labeling setting",
			args: args{
				a:                      agg,
				name:                   "name",
				owner:                  "owner",
				privateLabelingSetting: -1,
			},
			want: Want{
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "PROJE-AO52V", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid owner",
			args: args{
				a:                      agg,
				name:                   "name",
				owner:                  "",
				privateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
			},
			want: Want{
				ValidationErr: zerrors.ThrowPreconditionFailed(nil, "PROJE-hzxwo", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "correct",
			args: args{
				a:                      agg,
				name:                   "ZITADEL",
				owner:                  "CAOS AG",
				privateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
			},
			want: Want{
				Commands: []eventstore.Command{
					project.NewProjectAddedEvent(ctx, &agg.Aggregate,
						"ZITADEL",
						false,
						false,
						false,
						domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertValidation(t, context.Background(), AddProjectCommand(tt.args.a, tt.args.name, tt.args.owner, false, false, false, tt.args.privateLabelingSetting), nil, tt.want)
		})
	}
}
