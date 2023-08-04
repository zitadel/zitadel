package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/member"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func TestCommandSide_AddProject(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx           context.Context
		project       *domain.Project
		resourceOwner string
		ownerID       string
	}
	type res struct {
		want *domain.Project
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				project:       &domain.Project{},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "org with project owner, error already exists",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectPushFailed(errors.ThrowAlreadyExists(nil, "ERROR", "internl"),
						[]*repository.Event{
							eventFromEventPusher(project.NewProjectAddedEvent(
								context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
							),
							),
							eventFromEventPusher(project.NewProjectMemberAddedEvent(
								context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"user1",
								[]string{domain.RoleProjectOwner}...,
							),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewAddProjectNameUniqueConstraint("project", "org1")),
						uniqueConstraintsFromEventConstraint(member.NewAddMemberUniqueConstraint("project1", "user1")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "project1"),
			},
			args: args{
				ctx: context.Background(),
				project: &domain.Project{
					Name:                   "project",
					ProjectRoleAssertion:   true,
					ProjectRoleCheck:       true,
					HasProjectCheck:        true,
					PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
				},
				resourceOwner: "org1",
				ownerID:       "user1",
			},
			res: res{
				err: errors.IsErrorAlreadyExists,
			},
		},
		{
			name: "org with project owner, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewProjectAddedEvent(
								context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
							),
							),
							eventFromEventPusher(project.NewProjectMemberAddedEvent(
								context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"user1",
								[]string{domain.RoleProjectOwner}...,
							),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewAddProjectNameUniqueConstraint("project", "org1")),
						uniqueConstraintsFromEventConstraint(member.NewAddMemberUniqueConstraint("project1", "user1")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "project1"),
			},
			args: args{
				ctx: context.Background(),
				project: &domain.Project{
					Name:                   "project",
					ProjectRoleAssertion:   true,
					ProjectRoleCheck:       true,
					HasProjectCheck:        true,
					PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
				},
				resourceOwner: "org1",
				ownerID:       "user1",
			},
			res: res{
				want: &domain.Project{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
						AggregateID:   "project1",
					},
					Name:                   "project",
					ProjectRoleAssertion:   true,
					ProjectRoleCheck:       true,
					HasProjectCheck:        true,
					PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				Eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			got, err := r.AddProject(tt.args.ctx, tt.args.project, tt.args.resourceOwner, tt.args.ownerID)
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

func TestCommandSide_ChangeProject(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		project       *domain.Project
		resourceOwner string
	}
	type res struct {
		want *domain.Project
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				project: &domain.Project{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid project empty aggregateid, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				project: &domain.Project{
					Name: "project",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				project: &domain.Project{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Name: "project change",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsNotFound,
			},
		},
		{
			name: "project removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
			},
			args: args{
				ctx: context.Background(),
				project: &domain.Project{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Name: "project change",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				project: &domain.Project{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Name:                   "project",
					ProjectRoleAssertion:   true,
					ProjectRoleCheck:       true,
					HasProjectCheck:        true,
					PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			name: "project change with name and unique constraints, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newProjectChangedEvent(context.Background(),
									"project1",
									"org1",
									"project",
									"project-new",
									false,
									false,
									false,
									domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewRemoveProjectNameUniqueConstraint("project", "org1")),
						uniqueConstraintsFromEventConstraint(project.NewAddProjectNameUniqueConstraint("project-new", "org1")),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				project: &domain.Project{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Name:                   "project-new",
					ProjectRoleAssertion:   false,
					ProjectRoleCheck:       false,
					HasProjectCheck:        false,
					PrivateLabelingSetting: domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Project{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					Name:                   "project-new",
					ProjectRoleAssertion:   false,
					ProjectRoleCheck:       false,
					HasProjectCheck:        false,
					PrivateLabelingSetting: domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
				},
			},
		},
		{
			name: "project change without name and unique constraints, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newProjectChangedEvent(context.Background(),
									"project1",
									"org1",
									"project",
									"",
									false,
									false,
									false,
									domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				project: &domain.Project{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Name:                   "project",
					ProjectRoleAssertion:   false,
					ProjectRoleCheck:       false,
					HasProjectCheck:        false,
					PrivateLabelingSetting: domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Project{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					Name:                   "project",
					ProjectRoleAssertion:   false,
					ProjectRoleCheck:       false,
					HasProjectCheck:        false,
					PrivateLabelingSetting: domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				Eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeProject(tt.args.ctx, tt.args.project, tt.args.resourceOwner)
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

func TestCommandSide_DeactivateProject(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "",
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid resourceowner, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsNotFound,
			},
		},
		{
			name: "project removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsNotFound,
			},
		},
		{
			name: "project already inactive, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			name: "project deactivate, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								project.NewProjectDeactivatedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate),
							),
						},
					),
				),
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
				Eventstore: tt.fields.eventstore,
			}
			got, err := r.DeactivateProject(tt.args.ctx, tt.args.projectID, tt.args.resourceOwner)
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

func TestCommandSide_ReactivateProject(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "",
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid resourceowner, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsNotFound,
			},
		},
		{
			name: "project removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsNotFound,
			},
		},
		{
			name: "project not inactive, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			name: "project reactivate, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
						[]*repository.Event{
							eventFromEventPusher(
								project.NewProjectReactivatedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate),
							),
						},
					),
				),
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
				Eventstore: tt.fields.eventstore,
			}
			got, err := r.ReactivateProject(tt.args.ctx, tt.args.projectID, tt.args.resourceOwner)
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

func TestCommandSide_RemoveProject(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "",
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid resourceowner, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsNotFound,
			},
		},
		{
			name: "project removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsNotFound,
			},
		},
		{
			name: "project remove, without entityConstraints, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
						[]*repository.Event{
							eventFromEventPusher(
								project.NewProjectRemovedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"project",
									nil),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewRemoveProjectNameUniqueConstraint("project", "org1")),
					),
				),
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
				eventstore: eventstoreExpect(
					t,
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
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								project.NewProjectRemovedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"project",
									[]*eventstore.EventUniqueConstraint{
										project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test.com/saml/metadata"),
									}),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewRemoveProjectNameUniqueConstraint("project", "org1")),
						uniqueConstraintsFromEventConstraint(project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test.com/saml/metadata")),
					),
				),
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
				eventstore: eventstoreExpect(
					t,
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
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								project.NewProjectRemovedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"project",
									[]*eventstore.EventUniqueConstraint{
										project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test1.com/saml/metadata"),
										project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test2.com/saml/metadata"),
										project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test3.com/saml/metadata"),
									}),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewRemoveProjectNameUniqueConstraint("project", "org1")),
						uniqueConstraintsFromEventConstraint(project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test1.com/saml/metadata")),
						uniqueConstraintsFromEventConstraint(project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test2.com/saml/metadata")),
						uniqueConstraintsFromEventConstraint(project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test3.com/saml/metadata")),
					),
				),
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
				Eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveProject(tt.args.ctx, tt.args.projectID, tt.args.resourceOwner)
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
	event, _ := project.NewProjectChangeEvent(ctx,
		&project.NewAggregate(projectID, resourceOwner).Aggregate,
		oldName,
		changes,
	)
	return event
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
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-C01yo", "Errors.Invalid.Argument"),
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
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-AO52V", "Errors.Invalid.Argument"),
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
				ValidationErr: errors.ThrowPreconditionFailed(nil, "PROJE-hzxwo", "Errors.Invalid.Argument"),
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
					project.NewProjectMemberAddedEvent(ctx, &agg.Aggregate,
						"CAOS AG",
						domain.RoleProjectOwner),
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

// func TestExistsProject(t *testing.T) {
// 	type args struct {
// 		filter        preparation.FilterToQueryReducer
// 		id            string
// 		resourceOwner string
// 	}
// 	tests := []struct {
// 		name       string
// 		args       args
// 		wantExists bool
// 		wantErr    bool
// 	}{
// 		{
// 			name: "no events",
// 			args: args{
// 				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
// 					return []eventstore.Event{}, nil
// 				},
// 				id:            "id",
// 				resourceOwner: "ro",
// 			},
// 			wantExists: false,
// 			wantErr:    false,
// 		},
// 		{
// 			name: "project added",
// 			args: args{
// 				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
// 					return []eventstore.Event{
// 						project.NewProjectAddedEvent(
// 							context.Background(),
// 							&project.NewAggregate("id", "ro").Aggregate,
// 							"name",
// 							false,
// 							false,
// 							false,
// 							domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
// 						),
// 					}, nil
// 				},
// 				id:            "id",
// 				resourceOwner: "ro",
// 			},
// 			wantExists: true,
// 			wantErr:    false,
// 		},
// 		{
// 			name: "project removed",
// 			args: args{
// 				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
// 					return []eventstore.Event{
// 						project.NewProjectAddedEvent(
// 							context.Background(),
// 							&project.NewAggregate("id", "ro").Aggregate,
// 							"name",
// 							false,
// 							false,
// 							false,
// 							domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
// 						),
// 						project.NewProjectRemovedEvent(
// 							context.Background(),
// 							&project.NewAggregate("id", "ro").Aggregate,
// 							"name",
// 						),
// 					}, nil
// 				},
// 				id:            "id",
// 				resourceOwner: "ro",
// 			},
// 			wantExists: false,
// 			wantErr:    false,
// 		},
// 		{
// 			name: "error durring filter",
// 			args: args{
// 				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
// 					return nil, errors.ThrowInternal(nil, "PROJE-Op26p", "Errors.Internal")
// 				},
// 				id:            "id",
// 				resourceOwner: "ro",
// 			},
// 			wantExists: false,
// 			wantErr:    true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotExists, err := projectWriteModel(context.Background(), tt.args.filter, tt.args.id, tt.args.resourceOwner)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ExistsUser() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if gotExists != tt.wantExists {
// 				t.Errorf("ExistsUser() = %v, want %v", gotExists, tt.wantExists)
// 			}
// 		})
// 	}
// }
