package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddUserGrant(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id_generator.Generator
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
			name: "invalid usergrant, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID: "user1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
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
								nil,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:    "user1",
					ProjectID: "project1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
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
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:    "user1",
					ProjectID: "project1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project on other org, precondition error",
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
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "org2", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:    "user1",
					ProjectID: "project1",
				},
				resourceOwner: "org2",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
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
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"roleKey"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
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
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"roleKey"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
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
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
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
				ctx: authz.NewMockContextWithPermissions("", "org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"roleKey"},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project grant on other org, precondition error",
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
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
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
								"org3",
								[]string{"rolekey1"},
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "org2", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"rolekey1"},
				},
				resourceOwner: "org2",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
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
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
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
						usergrant.NewUserGrantAddedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							"user1",
							"project1",
							"",
							[]string{"rolekey1"},
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "usergrant1"),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
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
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
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
						usergrant.NewUserGrantAddedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							"user1",
							"project1",
							"projectgrant1",
							[]string{"rolekey1"},
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "usergrant1"),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
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
				eventstore: tt.fields.eventstore,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
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
			name: "invalid usergrant, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID: "user1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid permissions, error",
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
				ctx: context.Background(),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "usergrant1",
					},
					UserID: "user1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPermissionDenied,
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
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
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
				err: zerrors.IsNotFound,
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
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
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
				err: zerrors.IsNotFound,
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
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
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
				err: zerrors.IsPreconditionFailed,
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
								nil,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "org", "user", []string{domain.RoleProjectOwner}),
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
				err: zerrors.IsPreconditionFailed,
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
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "org", "user", []string{domain.RoleProjectOwner}),
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
				err: zerrors.IsPreconditionFailed,
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
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "org", "user", []string{domain.RoleProjectOwner}),
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
				err: zerrors.IsPreconditionFailed,
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
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "org", "user", []string{domain.RoleProjectOwner}),
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
				err: zerrors.IsPreconditionFailed,
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
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
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
				ctx: authz.NewMockContextWithPermissions("", "org", "user", []string{domain.RoleProjectOwner}),
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
				err: zerrors.IsPreconditionFailed,
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
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
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
						usergrant.NewUserGrantChangedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							[]string{"rolekey1", "rolekey2"},
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
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
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
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
						usergrant.NewUserGrantChangedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							[]string{"rolekey1", "rolekey2"},
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsErrorInvalidArgument,
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
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
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
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
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
				err: zerrors.IsPermissionDenied,
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
				err: zerrors.IsPreconditionFailed,
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
						usergrant.NewUserGrantDeactivatedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsErrorInvalidArgument,
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
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
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
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
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
				err: zerrors.IsPermissionDenied,
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
				err: zerrors.IsPreconditionFailed,
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
						usergrant.NewUserGrantReactivatedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
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
				err: zerrors.IsErrorInvalidArgument,
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
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
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
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
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
				err: zerrors.IsPermissionDenied,
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
						usergrant.NewUserGrantRemovedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							"user1",
							"project1",
							"",
						),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
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
						usergrant.NewUserGrantRemovedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							"user1",
							"project1",
							"projectgrant1",
						),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
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
				err: zerrors.IsErrorInvalidArgument,
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
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrantIDs:  []string{"usergrant1", "usergrant2"},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
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
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrantIDs:  []string{"usergrant1", "usergrant2"},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
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
				err: zerrors.IsPermissionDenied,
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
						usergrant.NewUserGrantRemovedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							"user1",
							"project1",
							"",
						),
						usergrant.NewUserGrantRemovedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant2", "org1").Aggregate,
							"user2",
							"project2",
							"",
						),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
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
						usergrant.NewUserGrantRemovedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
							"user1",
							"project1",
							"projectgrant1",
						),
						usergrant.NewUserGrantRemovedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant2", "org1").Aggregate,
							"user2",
							"project2",
							"projectgrant2",
						),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
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
