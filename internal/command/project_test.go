package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddProject(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id_generator.Generator
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
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "org with project owner, error already exists",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectPushFailed(zerrors.ThrowAlreadyExists(nil, "ERROR", "internl"),
						project.NewProjectAddedEvent(
							context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"project", true, true, true,
							domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
						),
						project.NewProjectMemberAddedEvent(
							context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"user1",
							[]string{domain.RoleProjectOwner}...,
						),
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
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "org with project owner, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectPush(
						project.NewProjectAddedEvent(
							context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"project", true, true, true,
							domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
						),
						project.NewProjectMemberAddedEvent(
							context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"user1",
							[]string{domain.RoleProjectOwner}...,
						),
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
				eventstore: tt.fields.eventstore,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsNotFound,
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
				err: zerrors.IsPreconditionFailed,
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
				eventstore: tt.fields.eventstore,
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsNotFound,
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
				err: zerrors.IsNotFound,
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
				err: zerrors.IsPreconditionFailed,
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
						project.NewProjectDeactivatedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate),
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
				eventstore: tt.fields.eventstore,
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsNotFound,
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
				err: zerrors.IsNotFound,
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
				err: zerrors.IsPreconditionFailed,
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
						project.NewProjectReactivatedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate),
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
				eventstore: tt.fields.eventstore,
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsNotFound,
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
				err: zerrors.IsNotFound,
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
						project.NewProjectRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"project",
							nil),
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
						project.NewProjectRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"project",
							[]*eventstore.UniqueConstraint{
								project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test.com/saml/metadata"),
							},
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
				eventstore: tt.fields.eventstore,
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
