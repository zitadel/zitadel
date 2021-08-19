package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/repository/usergrant"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandSide_AddProjectRole(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		role          *domain.ProjectRole
		resourceOwner string
	}
	type res struct {
		want *domain.ProjectRole
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
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
							),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				role: &domain.ProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Key: "key1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid role, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				role: &domain.ProjectRole{
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
			name: "role key already exists, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
							),
						),
					),
					expectPushFailed(caos_errs.ThrowAlreadyExists(nil, "id", "internal"),
						[]*repository.Event{
							eventFromEventPusher(project.NewRoleAddedEvent(
								context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"group",
							),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewAddProjectRoleUniqueConstraint("key1", "project1")),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				role: &domain.ProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Key:         "key1",
					DisplayName: "key",
					Group:       "group",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "add role,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewRoleAddedEvent(
								context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"group",
							),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewAddProjectRoleUniqueConstraint("key1", "project1")),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				role: &domain.ProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Key:         "key1",
					DisplayName: "key",
					Group:       "group",
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					Key:         "key1",
					DisplayName: "key",
					Group:       "group",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddProjectRole(tt.args.ctx, tt.args.role, tt.args.resourceOwner)
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

func TestCommandSide_BulkAddProjectRole(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		roles         []*domain.ProjectRole
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
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
							),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				roles: []*domain.ProjectRole{
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
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid role, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				roles: []*domain.ProjectRole{
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
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "role key already exists, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
							),
						),
					),
					expectPushFailed(caos_errs.ThrowAlreadyExists(nil, "id", "internal"),
						[]*repository.Event{
							eventFromEventPusher(project.NewRoleAddedEvent(
								context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"group",
							),
							),
							eventFromEventPusher(project.NewRoleAddedEvent(
								context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key2",
								"key2",
								"group",
							),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewAddProjectRoleUniqueConstraint("key1", "project1")),
						uniqueConstraintsFromEventConstraint(project.NewAddProjectRoleUniqueConstraint("key2", "project1")),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				roles: []*domain.ProjectRole{
					{
						Key:         "key1",
						DisplayName: "key",
						Group:       "group",
					},
					{
						Key:         "key2",
						DisplayName: "key2",
						Group:       "group",
					},
				},
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "add roles,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewRoleAddedEvent(
								context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key1",
								"key",
								"group",
							),
							),
							eventFromEventPusher(project.NewRoleAddedEvent(
								context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"key2",
								"key2",
								"group",
							),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewAddProjectRoleUniqueConstraint("key1", "project1")),
						uniqueConstraintsFromEventConstraint(project.NewAddProjectRoleUniqueConstraint("key2", "project1")),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				roles: []*domain.ProjectRole{
					{
						Key:         "key1",
						DisplayName: "key",
						Group:       "group",
					},
					{
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
				eventstore: tt.fields.eventstore,
			}
			got, err := r.BulkAddProjectRole(tt.args.ctx, tt.args.projectID, tt.args.resourceOwner, tt.args.roles)
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

func TestCommandSide_ChangeProjectRole(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		role          *domain.ProjectRole
		resourceOwner string
	}
	type res struct {
		want *domain.ProjectRole
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				role: &domain.ProjectRole{
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
			name: "project not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
							),
						),
						eventFromEventPusher(
							project.NewProjectRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				role: &domain.ProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Key: "key1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "role removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
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
			},
			args: args{
				ctx: context.Background(),
				role: &domain.ProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Key:         "key1",
					DisplayName: "key",
					Group:       "group",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "role not changed, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
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
			},
			args: args{
				ctx: context.Background(),
				role: &domain.ProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Key:         "key1",
					DisplayName: "key",
					Group:       "group",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "role changed, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1", true, true, true,
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
						[]*repository.Event{
							eventFromEventPusher(
								newRoleChangedEvent(context.Background(), "project1", "org1", "key1", "keychanged", "groupchanged"),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				role: &domain.ProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					Key:         "key1",
					DisplayName: "keychanged",
					Group:       "groupchanged",
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ProjectRole{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					Key:         "key1",
					DisplayName: "keychanged",
					Group:       "groupchanged",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeProjectRole(tt.args.ctx, tt.args.role, tt.args.resourceOwner)
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

func TestCommandSide_RemoveProjectRole(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				key:           "key1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid key, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				key:           "",
				projectID:     "project1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "role not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				key:           "key",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "role removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				key:           "key",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "role removed, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
						[]*repository.Event{
							eventFromEventPusher(
								project.NewRoleRemovedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"key1",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewRemoveProjectRoleUniqueConstraint("key1", "project1")),
					),
				),
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
				eventstore: eventstoreExpect(
					t,
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
						[]*repository.Event{
							eventFromEventPusher(
								project.NewRoleRemovedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"key1",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewRemoveProjectRoleUniqueConstraint("key1", "project1")),
					),
				),
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
				eventstore: eventstoreExpect(
					t,
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
						[]*repository.Event{
							eventFromEventPusher(
								project.NewRoleRemovedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"key1",
								),
							),
							eventFromEventPusher(
								project.NewGrantCascadeChangedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"projectgrant1",
									[]string{},
								),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewRemoveProjectRoleUniqueConstraint("key1", "project1")),
					),
				),
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
				eventstore: eventstoreExpect(
					t,
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
						[]*repository.Event{
							eventFromEventPusher(
								project.NewRoleRemovedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"key1",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewRemoveProjectRoleUniqueConstraint("key1", "project1")),
					),
				),
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
				eventstore: eventstoreExpect(
					t,
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
						[]*repository.Event{
							eventFromEventPusher(
								project.NewRoleRemovedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"key1",
								),
							),
							eventFromEventPusher(
								usergrant.NewUserGrantCascadeChangedEvent(context.Background(),
									&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
									[]string{},
								),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewRemoveProjectRoleUniqueConstraint("key1", "project1")),
					),
				),
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
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveProjectRole(tt.args.ctx, tt.args.projectID, tt.args.key, tt.args.resourceOwner, tt.args.cascadingProjectGrantIDs, tt.args.cascadingUserGrantIDs...)
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
