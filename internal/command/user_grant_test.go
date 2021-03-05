package command

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
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
		//{
		//	name: "invalid permissions, error",
		//	fields: fields{
		//		eventstore: eventstoreExpect(
		//			t,
		//		),
		//	},
		//	args: args{
		//		ctx: context.Background(),
		//		userGrant: &domain.UserGrant{
		//			UserID: "user1",
		//		},
		//		resourceOwner: "org1",
		//	},
		//	res: res{
		//		err: caos_errs.IsPermissionDenied,
		//	},
		//},
		//{
		//	name: "invalid usergrant, error",
		//	fields: fields{
		//		eventstore: eventstoreExpect(
		//			t,
		//		),
		//	},
		//	args: args{
		//		ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
		//		userGrant: &domain.UserGrant{
		//			UserID: "user1",
		//		},
		//		resourceOwner: "org1",
		//	},
		//	res: res{
		//		err: caos_errs.IsErrorInvalidArgument,
		//	},
		//},
		//{
		//	name: "user removed, precondition error",
		//	fields: fields{
		//		eventstore: eventstoreExpect(
		//			t,
		//			expectFilter(
		//				eventFromEventPusher(
		//					user.NewHumanAddedEvent(context.Background(),
		//						&user.NewAggregate("user1", "org1").Aggregate,
		//						"username1",
		//						"firstname1",
		//						"lastname1",
		//						"nickname1",
		//						"displayname1",
		//						language.German,
		//						domain.GenderMale,
		//						"email1",
		//						true,
		//					),
		//				),
		//				eventFromEventPusher(
		//					user.NewUserRemovedEvent(
		//						context.Background(),
		//						&user.NewAggregate("user1", "org1").Aggregate,
		//						"username1",
		//						true,
		//					),
		//				),
		//			),
		//		),
		//	},
		//	args: args{
		//		ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
		//		userGrant: &domain.UserGrant{
		//			UserID: "user1",
		//			ProjectID: "project1",
		//		},
		//		resourceOwner: "org1",
		//	},
		//	res: res{
		//		err: caos_errs.IsPreconditionFailed,
		//	},
		//},
		//{
		//	name: "project removed, precondition error",
		//	fields: fields{
		//		eventstore: eventstoreExpect(
		//			t,
		//			expectFilter(
		//				eventFromEventPusher(
		//					user.NewHumanAddedEvent(context.Background(),
		//						&user.NewAggregate("user1", "org1").Aggregate,
		//						"username1",
		//						"firstname1",
		//						"lastname1",
		//						"nickname1",
		//						"displayname1",
		//						language.German,
		//						domain.GenderMale,
		//						"email1",
		//						true,
		//					),
		//				),
		//				eventFromEventPusher(
		//					project.NewProjectAddedEvent(context.Background(),
		//						&project.NewAggregate("project1", "org1").Aggregate,
		//						"projectname1",
		//					),
		//				),
		//				eventFromEventPusher(
		//					project.NewProjectRemovedEvent(context.Background(),
		//						&project.NewAggregate("project1", "org1").Aggregate,
		//						"projectname1",
		//					),
		//				),
		//			),
		//		),
		//	},
		//	args: args{
		//		ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
		//		userGrant: &domain.UserGrant{
		//			UserID: "user1",
		//			ProjectID: "project1",
		//		},
		//		resourceOwner: "org1",
		//	},
		//	res: res{
		//		err: caos_errs.IsPreconditionFailed,
		//	},
		//},
		//{
		//	name: "project roles not existing, precondition error",
		//	fields: fields{
		//		eventstore: eventstoreExpect(
		//			t,
		//			expectFilter(
		//				eventFromEventPusher(
		//					user.NewHumanAddedEvent(context.Background(),
		//						&user.NewAggregate("user1", "org1").Aggregate,
		//						"username1",
		//						"firstname1",
		//						"lastname1",
		//						"nickname1",
		//						"displayname1",
		//						language.German,
		//						domain.GenderMale,
		//						"email1",
		//						true,
		//					),
		//				),
		//				eventFromEventPusher(
		//					project.NewProjectAddedEvent(context.Background(),
		//						&project.NewAggregate("project1", "org1").Aggregate,
		//						"projectname1",
		//					),
		//				),
		//			),
		//		),
		//	},
		//	args: args{
		//		ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
		//		userGrant: &domain.UserGrant{
		//			UserID: "user1",
		//			ProjectID: "project1",
		//			RoleKeys: []string{"roleKey"},
		//		},
		//		resourceOwner: "org1",
		//	},
		//	res: res{
		//		err: caos_errs.IsPreconditionFailed,
		//	},
		//},
		//{
		//	name: "project grant not existing, precondition error",
		//	fields: fields{
		//		eventstore: eventstoreExpect(
		//			t,
		//			expectFilter(
		//				eventFromEventPusher(
		//					user.NewHumanAddedEvent(context.Background(),
		//						&user.NewAggregate("user1", "org1").Aggregate,
		//						"username1",
		//						"firstname1",
		//						"lastname1",
		//						"nickname1",
		//						"displayname1",
		//						language.German,
		//						domain.GenderMale,
		//						"email1",
		//						true,
		//					),
		//				),
		//				eventFromEventPusher(
		//					project.NewProjectAddedEvent(context.Background(),
		//						&project.NewAggregate("project1", "org1").Aggregate,
		//						"projectname1",
		//					),
		//				),
		//			),
		//		),
		//	},
		//	args: args{
		//		ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
		//		userGrant: &domain.UserGrant{
		//			UserID: "user1",
		//			ProjectID: "project1",
		//			ProjectGrantID: "projectgrant1",
		//			RoleKeys: []string{"roleKey"},
		//		},
		//		resourceOwner: "org1",
		//	},
		//	res: res{
		//		err: caos_errs.IsPreconditionFailed,
		//	},
		//},
		//{
		//	name: "project grant roles not existing, precondition error",
		//	fields: fields{
		//		eventstore: eventstoreExpect(
		//			t,
		//			expectFilter(
		//				eventFromEventPusher(
		//					user.NewHumanAddedEvent(context.Background(),
		//						&user.NewAggregate("user1", "org1").Aggregate,
		//						"username1",
		//						"firstname1",
		//						"lastname1",
		//						"nickname1",
		//						"displayname1",
		//						language.German,
		//						domain.GenderMale,
		//						"email1",
		//						true,
		//					),
		//				),
		//				eventFromEventPusher(
		//					project.NewProjectAddedEvent(context.Background(),
		//						&project.NewAggregate("project1", "org1").Aggregate,
		//						"projectname1",
		//					),
		//				),
		//				eventFromEventPusher(
		//					project.NewRoleAddedEvent(context.Background(),
		//						&project.NewAggregate("project1", "org1").Aggregate,
		//						"rolekey1",
		//						"rolekey",
		//						"",
		//						"project1",
		//					),
		//				),
		//				eventFromEventPusher(
		//					project.NewGrantAddedEvent(context.Background(),
		//						&project.NewAggregate("project1", "org1").Aggregate,
		//						"projectgrant1",
		//						"org2",
		//						"project1",
		//						nil,
		//						),
		//				),
		//			),
		//		),
		//	},
		//	args: args{
		//		ctx: authz.NewMockContextWithPermissions("org", "user", []string{domain.RoleProjectOwner}),
		//		userGrant: &domain.UserGrant{
		//			UserID: "user1",
		//			ProjectID: "project1",
		//			ProjectGrantID: "projectgrant1",
		//			RoleKeys: []string{"roleKey"},
		//		},
		//		resourceOwner: "org1",
		//	},
		//	res: res{
		//		err: caos_errs.IsPreconditionFailed,
		//	},
		//},
		//{
		//	name: "usergrant for project, ok",
		//	fields: fields{
		//		eventstore: eventstoreExpect(
		//			t,
		//			expectFilter(
		//				eventFromEventPusher(
		//					user.NewHumanAddedEvent(context.Background(),
		//						&user.NewAggregate("user1", "org1").Aggregate,
		//						"username1",
		//						"firstname1",
		//						"lastname1",
		//						"nickname1",
		//						"displayname1",
		//						language.German,
		//						domain.GenderMale,
		//						"email1",
		//						true,
		//					),
		//				),
		//				eventFromEventPusher(
		//					project.NewProjectAddedEvent(context.Background(),
		//						&project.NewAggregate("project1", "org1").Aggregate,
		//						"projectname1",
		//					),
		//				),
		//				eventFromEventPusher(
		//					project.NewRoleAddedEvent(context.Background(),
		//						&project.NewAggregate("project1", "org1").Aggregate,
		//						"rolekey1",
		//						"rolekey",
		//						"",
		//						"project1",
		//					),
		//				),
		//			),
		//			expectPush(
		//				[]*repository.Event{
		//					eventFromEventPusher(usergrant.NewUserGrantAddedEvent(context.Background(),
		//						&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
		//						"user1",
		//						"project1",
		//					"",
		//					[]string{"rolekey1"},
		//					)),
		//				},
		//				uniqueConstraintsFromEventConstraint(usergrant.NewAddUserGrantUniqueConstraint("org1", "user1", "project1", "")),
		//			),
		//		),
		//		idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "usergrant1"),
		//	},
		//	args: args{
		//		ctx: authz.NewMockContextWithPermissions("", "", []string{domain.RoleProjectOwner}),
		//		userGrant: &domain.UserGrant{
		//			UserID: "user1",
		//			ProjectID: "project1",
		//			RoleKeys: []string{"rolekey1"},
		//		},
		//		resourceOwner: "org1",
		//	},
		//	res: res{
		//		want: &domain.UserGrant{
		//			ObjectRoot: models.ObjectRoot{
		//				AggregateID:   "usergrant1",
		//				ResourceOwner: "org1",
		//			},
		//			UserID:          "user1",
		//			ProjectID: "project1",
		//			RoleKeys: []string{"rolekey1"},
		//			State:         domain.UserGrantStateActive,
		//		},
		//	},
		//},
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
								"project1",
							),
						),
						eventFromEventPusher(
							project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"projectgrant1",
								"org1",
								"project1",
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
								"",
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
