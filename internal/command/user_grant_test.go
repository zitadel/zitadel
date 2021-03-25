package command

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/id"
	id_mock "github.com/caos/zitadel/internal/id/mock"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/repository/usergrant"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

func TestCommandSide_AddUserGrant(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx           context.Context
		userGrant     *domain.UserGrant
		resourceOwner string
	}
	type res struct {
		want *domain.UserGrant
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid permissions, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				userGrant: &domain.UserGrant{
					UserID: "user1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPermissionDenied,
			},
		},
		{
			name: "invalid usergrant, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID: "user1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "user removed, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
						eventFromEventPusher(
							user.NewUserRemovedEvent(
								context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username1",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:    "user1",
					ProjectID: "project1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "project removed, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
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
				ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:    "user1",
					ProjectID: "project1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "project roles not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"roleKey"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "project grant not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"roleKey"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "project grant roles not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"rolekey1",
								"rolekey",
								"",
							),
						),
						eventFromEventPusher(
							project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"org2",
								nil,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"roleKey"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "usergrant for project, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"rolekey1",
								"rolekey",
								"",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"",
								[]string{"rolekey1"},
							)),
						},
						uniqueConstraintsFromEventConstraint(usergrant.NewAddUserGrantUniqueConstraint("org1", "user1", "project1", "")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "usergrant1"),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1"},
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1"},
					State:     domain.UserGrantStateActive,
				},
			},
		},
		{
			name: "usergrant for projectgrant, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"rolekey1",
								"rolekey",
								"",
							),
						),
						eventFromEventPusher(
							project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"org1",
								[]string{"rolekey1"},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"projectgrant1",
								[]string{"rolekey1"},
							)),
						},
						uniqueConstraintsFromEventConstraint(usergrant.NewAddUserGrantUniqueConstraint("org1", "user1", "project1", "projectgrant1")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "usergrant1"),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"rolekey1"},
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"rolekey1"},
					State:          domain.UserGrantStateActive,
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
			got, err := r.AddUserGrant(tt.args.ctx, tt.args.userGrant, tt.args.resourceOwner)
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

func TestCommandSide_ChangeUserGrant(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userGrant     *domain.UserGrant
		resourceOwner string
	}
	type res struct {
		want *domain.UserGrant
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid permissions, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				userGrant: &domain.UserGrant{
					UserID: "user1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPermissionDenied,
			},
		},
		{
			name: "invalid usergrant, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID: "user1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "usergrant not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "usergrant1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "usergrant not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "usergrant1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "usergrant roles not changed, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "usergrant1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "user removed, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
					),
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
						eventFromEventPusher(
							user.NewUserRemovedEvent(
								context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username1",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "usergrant1",
					},
					UserID:    "user1",
					ProjectID: "project1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "project removed, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
					),
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
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
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
				ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "usergrant1",
					},
					UserID:    "user1",
					ProjectID: "project1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "project roles not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
					),
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
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "usergrant1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"roleKey"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "project grant not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
					),
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
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "usergrant1",
					},
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"roleKey"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "project grant roles not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
					),
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
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"rolekey1",
								"rolekey",
								"",
							),
						),
						eventFromEventPusher(
							project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"org2",
								nil,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "usergrant1",
					},
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"roleKey"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "usergrant for project, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
					),
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
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"rolekey1",
								"rolekey",
								"",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"rolekey2",
								"rolekey 2",
								"",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(usergrant.NewUserGrantChangedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								[]string{"rolekey1", "rolekey2"},
							)),
						},
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "usergrant1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1", "rolekey2"},
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1", "rolekey2"},
					State:     domain.UserGrantStateActive,
				},
			},
		},
		{
			name: "usergrant for projectgrant, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"projectgrant1", []string{"rolekey1"}),
						),
					),
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
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectname1",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"rolekey1",
								"rolekey",
								"",
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"rolekey2",
								"rolekey2",
								"",
							),
						),
						eventFromEventPusher(
							project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"org1",
								[]string{"rolekey1", "rolekey2"},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(usergrant.NewUserGrantChangedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								[]string{"rolekey1", "rolekey2"},
							)),
						},
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "usergrant1",
					},
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"rolekey1", "rolekey2"},
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"rolekey1", "rolekey2"},
					State:          domain.UserGrantStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeUserGrant(tt.args.ctx, tt.args.userGrant, tt.args.resourceOwner)
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

func TestCommandSide_DeactivateUserGrant(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userGrantID   string
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
			name: "invalid usergrantID, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid resourceOwner, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:         context.Background(),
				userGrantID: "usergrant1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "usergrant not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "usergrant removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
						eventFromEventPusher(
							usergrant.NewUserGrantRemovedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								""),
						),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "no permissions, permisison denied error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPermissionDenied,
			},
		},
		{
			name: "already deactivated, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
						eventFromEventPusher(
							usergrant.NewUserGrantDeactivatedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "deactivated, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								usergrant.NewUserGrantDeactivatedEvent(context.Background(),
									&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
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
			got, err := r.DeactivateUserGrant(tt.args.ctx, tt.args.userGrantID, tt.args.resourceOwner)
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

func TestCommandSide_ReactivateUserGrant(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userGrantID   string
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
			name: "invalid usergrantID, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid resourceOwner, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:         context.Background(),
				userGrantID: "usergrant1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "usergrant not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "usergrant removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
						eventFromEventPusher(
							usergrant.NewUserGrantRemovedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								""),
						),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "no permissions, permisison denied error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
						eventFromEventPusher(
							usergrant.NewUserGrantDeactivatedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPermissionDenied,
			},
		},
		{
			name: "already active, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "reactivated, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
						eventFromEventPusher(
							usergrant.NewUserGrantDeactivatedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								usergrant.NewUserGrantReactivatedEvent(context.Background(),
									&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
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
			got, err := r.ReactivateUserGrant(tt.args.ctx, tt.args.userGrantID, tt.args.resourceOwner)
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

func TestCommandSide_RemoveUserGrant(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userGrantID   string
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
			name: "invalid usergrantID, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "usergrant not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "usergrant removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
						eventFromEventPusher(
							usergrant.NewUserGrantRemovedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								""),
						),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "no permissions, permisison denied error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
						eventFromEventPusher(
							usergrant.NewUserGrantDeactivatedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPermissionDenied,
			},
		},
		{
			name: "remove usergrant project, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								usergrant.NewUserGrantRemovedEvent(context.Background(),
									&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
									"user1",
									"project1",
									"",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(usergrant.NewRemoveUserGrantUniqueConstraint("org1", "user1", "project1", "")),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "remove usergrant projectgrant, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"projectgrant1", []string{"rolekey1"}),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								usergrant.NewUserGrantRemovedEvent(context.Background(),
									&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
									"user1",
									"project1",
									"projectgrant1",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(usergrant.NewRemoveUserGrantUniqueConstraint("org1", "user1", "project1", "projectgrant1")),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
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
			got, err := r.RemoveUserGrant(tt.args.ctx, tt.args.userGrantID, tt.args.resourceOwner)
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

func TestCommandSide_BulkRemoveUserGrant(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userGrantIDs  []string
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
			name: "empty usergrantid list, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "usergrant not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrantIDs:  []string{"usergrant1", "usergrant2"},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "usergrant removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
						eventFromEventPusher(
							usergrant.NewUserGrantRemovedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								""),
						),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrantIDs:  []string{"usergrant1", "usergrant2"},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "no permissions, permisison denied error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userGrantIDs:  []string{"usergrant1", "usergrant2"},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPermissionDenied,
			},
		},
		{
			name: "remove usergrants project, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
					),
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant2", "org1").Aggregate,
								"user2",
								"project2",
								"", []string{"rolekey1"}),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								usergrant.NewUserGrantRemovedEvent(context.Background(),
									&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
									"user1",
									"project1",
									"",
								),
							),
							eventFromEventPusher(
								usergrant.NewUserGrantRemovedEvent(context.Background(),
									&usergrant.NewAggregate("usergrant2", "org1").Aggregate,
									"user2",
									"project2",
									"",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(usergrant.NewRemoveUserGrantUniqueConstraint("org1", "user1", "project1", "")),
						uniqueConstraintsFromEventConstraint(usergrant.NewRemoveUserGrantUniqueConstraint("org1", "user2", "project2", "")),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrantIDs:  []string{"usergrant1", "usergrant2"},
				resourceOwner: "org1",
			},
			res: res{},
		},
		{
			name: "remove usergrants projectgrant, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"projectgrant1", []string{"rolekey1"}),
						),
					),
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant2", "org1").Aggregate,
								"user2",
								"project2",
								"projectgrant2", []string{"rolekey1"}),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								usergrant.NewUserGrantRemovedEvent(context.Background(),
									&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
									"user1",
									"project1",
									"projectgrant1",
								),
							),
							eventFromEventPusher(
								usergrant.NewUserGrantRemovedEvent(context.Background(),
									&usergrant.NewAggregate("usergrant2", "org1").Aggregate,
									"user2",
									"project2",
									"projectgrant2",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(usergrant.NewRemoveUserGrantUniqueConstraint("org1", "user1", "project1", "projectgrant1")),
						uniqueConstraintsFromEventConstraint(usergrant.NewRemoveUserGrantUniqueConstraint("org1", "user2", "project2", "projectgrant2")),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
				userGrantIDs:  []string{"usergrant1", "usergrant2"},
				resourceOwner: "org1",
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			err := r.BulkRemoveUserGrant(tt.args.ctx, tt.args.userGrantIDs, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
