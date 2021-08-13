package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/id"
	id_mock "github.com/caos/zitadel/internal/id/mock"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/member"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/stretchr/testify/assert"
	"testing"
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
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "org with project owner, error already exists",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewGlobalOrgSetEventEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"globalorg",
							),
						),
					),
					expectPushFailed(caos_errs.ThrowAlreadyExists(nil, "ERROR", "internl"),
						[]*repository.Event{
							eventFromEventPusher(project.NewProjectAddedEvent(
								context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
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
					Name: "project",
				},
				resourceOwner: "org1",
				ownerID:       "user1",
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "global org with project owner global, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewGlobalOrgSetEventEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"globalorg",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewProjectAddedEvent(
								context.Background(),
								&project.NewAggregate("project1", "globalorg").Aggregate,
								"project", true, true, true,
							),
							),
							eventFromEventPusher(project.NewProjectMemberAddedEvent(
								context.Background(),
								&project.NewAggregate("project1", "globalorg").Aggregate,
								"user1",
								[]string{domain.RoleProjectOwnerGlobal}...,
							),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewAddProjectNameUniqueConstraint("project", "globalorg")),
						uniqueConstraintsFromEventConstraint(member.NewAddMemberUniqueConstraint("project1", "user1")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "project1"),
			},
			args: args{
				ctx: context.Background(),
				project: &domain.Project{
					Name: "project",
				},
				resourceOwner: "globalorg",
				ownerID:       "user1",
			},
			res: res{
				want: &domain.Project{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "globalorg",
						AggregateID:   "project1",
					},
					Name:                 "project",
					ProjectRoleAssertion: true,
					ProjectRoleCheck:     true,
					HasProjectCheck:      true,
				},
			},
		},
		{
			name: "org with project owner, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewGlobalOrgSetEventEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"globalorg",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewProjectAddedEvent(
								context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
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
					Name: "project",
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
					Name:                 "project",
					ProjectRoleAssertion: true,
					ProjectRoleCheck:     true,
					HasProjectCheck:      true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:  tt.fields.eventstore,
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsNotFound,
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
								"project", true, true, true),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project"),
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
				err: caos_errs.IsNotFound,
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
								"project", true, true, true),
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
					Name:                 "project",
					ProjectRoleAssertion: true,
					ProjectRoleCheck:     true,
					HasProjectCheck:      true,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
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
								"project", true, true, true),
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
									true,
									true,
									true),
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
					Name:                 "project-new",
					ProjectRoleAssertion: true,
					ProjectRoleCheck:     true,
					HasProjectCheck:      true,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Project{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					Name:                 "project-new",
					ProjectRoleAssertion: true,
					ProjectRoleCheck:     true,
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
								"project", true, true, true),
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
									true,
									true,
									true),
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
					Name:                 "project",
					ProjectRoleAssertion: true,
					ProjectRoleCheck:     true,
					HasProjectCheck:      true,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Project{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					Name:                 "project",
					ProjectRoleAssertion: true,
					ProjectRoleCheck:     true,
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsNotFound,
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
								"project", true, true, true),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project"),
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
				err: caos_errs.IsNotFound,
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
								"project", true, true, true),
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
				err: caos_errs.IsPreconditionFailed,
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
								"project", true, true, true),
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsNotFound,
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
								"project", true, true, true),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project"),
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
				err: caos_errs.IsNotFound,
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
								"project", true, true, true),
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
				err: caos_errs.IsPreconditionFailed,
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
								"project", true, true, true),
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsNotFound,
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
								"project", true, true, true),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project"),
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
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "project remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								project.NewProjectRemovedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"project"),
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

func newProjectChangedEvent(ctx context.Context, projectID, resourceOwner, oldName, newName string, roleAssertion, roleCheck, hasProjectCheck bool) *project.ProjectChangeEvent {
	changes := []project.ProjectChanges{
		project.ChangeProjectRoleAssertion(roleAssertion),
		project.ChangeProjectRoleCheck(roleCheck),
		project.ChangeHasProjectCheck(hasProjectCheck),
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
