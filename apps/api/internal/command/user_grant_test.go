package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	errMockedPermissionCheck   = errors.New("mocked permission check error")
	isMockedPermissionCheckErr = func(err error) bool {
		return errors.Is(err, errMockedPermissionCheck)
	}
	succeedingUserGrantPermissionCheck = func(_, _ string) PermissionCheck {
		return func(_, _ string) error { return nil }
	}
	failingUserGrantPermissionCheck = func(_, _ string) PermissionCheck {
		return func(_, _ string) error { return errMockedPermissionCheck }
	}
)

func TestCommandSide_AddUserGrant(t *testing.T) {
	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator func(t *testing.T) id.Generator
	}
	type args struct {
		ctx       context.Context
		userGrant *domain.UserGrant
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID: "user1",
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user removed, precondition error",
			fields: fields{
				eventstore: expectEventstore(
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
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project removed, precondition error",
			fields: fields{
				eventstore: expectEventstore(
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
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project on other org, precondition error",
			fields: fields{
				eventstore: expectEventstore(
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
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org2",
					},
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
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project grant not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
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
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project grant roles not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
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
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project grant on other org, precondition error",
			fields: fields{
				eventstore: expectEventstore(
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
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org2",
					},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "usergrant for project, ok",
			fields: fields{
				eventstore: expectEventstore(
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
				idGenerator: func(t *testing.T) id.Generator {
					return id_mock.NewIDGeneratorExpectIDs(t, "usergrant1")
				},
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1"},
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
				},
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
			name: "usergrant without resource owner on project, ok",
			fields: fields{
				eventstore: expectEventstore(
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
				idGenerator: func(t *testing.T) id.Generator {
					return id_mock.NewIDGeneratorExpectIDs(t, "usergrant1")
				},
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1"},
				},
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
				eventstore: expectEventstore(
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
								&project.NewAggregate("project1", "org2").Aggregate,
								"projectname1", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
						eventFromEventPusher(
							project.NewRoleAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org2").Aggregate,
								"rolekey1",
								"rolekey",
								"",
							),
						),
						eventFromEventPusher(
							project.NewGrantAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org2").Aggregate,
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
				idGenerator: func(t *testing.T) id.Generator {
					return id_mock.NewIDGeneratorExpectIDs(t, "usergrant1")
				},
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"rolekey1"},
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
				},
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
		{
			name: "usergrant for granted resource owner, ok",
			fields: fields{
				eventstore: expectEventstore(
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
				idGenerator: func(t *testing.T) id.Generator {
					return id_mock.NewIDGeneratorExpectIDs(t, "usergrant1")
				},
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1"},
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
				},
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
	t.Run("without permission check", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := &Commands{
					eventstore: tt.fields.eventstore(t),
				}
				if tt.fields.idGenerator != nil {
					r.idGenerator = tt.fields.idGenerator(t)
				}
				got, err := r.AddUserGrant(tt.args.ctx, tt.args.userGrant, nil)
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
	})
	t.Run("with succeeding permission check", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := &Commands{
					eventstore: tt.fields.eventstore(t),
				}
				if tt.fields.idGenerator != nil {
					r.idGenerator = tt.fields.idGenerator(t)
				}
				// we use an empty context and only rely on the permission check implementation
				got, err := r.AddUserGrant(context.Background(), tt.args.userGrant, succeedingUserGrantPermissionCheck)
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
	})
	t.Run("with failing permission check", func(t *testing.T) {
		r := &Commands{
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
			),
		}
		// we use an empty context and only rely on the permission check implementation
		_, err := r.AddUserGrant(context.Background(), &domain.UserGrant{
			UserID:    "user1",
			ProjectID: "project1",
			RoleKeys:  []string{"rolekey1"},
			ObjectRoot: models.ObjectRoot{
				ResourceOwner: "org1",
			},
		}, failingUserGrantPermissionCheck)
		assert.ErrorIs(t, err, errMockedPermissionCheck)
	})
}

func TestCommandSide_ChangeUserGrant(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx             context.Context
		userGrant       *domain.UserGrant
		permissionCheck UserGrantPermissionCheck
		cascade         bool
		ignoreUnchanged bool
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "org", "user", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					UserID: "user1",
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
					},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid permissions, error",
			fields: fields{
				eventstore: expectEventstore(
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
				),
			},
			args: args{
				ctx: context.Background(),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID: "user1",
				},
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "usergrant not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1"},
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "usergrant not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1"},
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "usergrant roles not changed, precondition error",
			fields: fields{
				eventstore: expectEventstore(
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
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1"},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "user removed, precondition error",
			fields: fields{
				eventstore: expectEventstore(
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
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:    "user1",
					ProjectID: "project1",
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project removed, precondition error",
			fields: fields{
				eventstore: expectEventstore(
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
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:    "user1",
					ProjectID: "project1",
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
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"roleKey"},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project grant not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
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
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"roleKey"},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project grant roles not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
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
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"roleKey"},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "usergrant for project, ok",
			fields: fields{
				eventstore: expectEventstore(
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
							"user1",
							[]string{"rolekey1", "rolekey2"},
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1", "rolekey2"},
				},
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
			name: "usergrant for project cascade, ok",
			fields: fields{
				eventstore: expectEventstore(
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
						usergrant.NewUserGrantCascadeChangedEvent(context.Background(),
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
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					RoleKeys: []string{"rolekey1", "rolekey2"},
				},
				cascade: true,
			},
			res: res{
				want: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					State:     domain.UserGrantStateActive,
					RoleKeys:  []string{"rolekey1", "rolekey2"},
				},
			},
		},
		{
			name: "usergrant for projectgrant, ok",
			fields: fields{
				eventstore: expectEventstore(
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
							"user1",
							[]string{"rolekey1", "rolekey2"},
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:         "user1",
					ProjectID:      "project1",
					ProjectGrantID: "projectgrant1",
					RoleKeys:       []string{"rolekey1", "rolekey2"},
				},
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
		{
			name: "usergrant for project without resource owner, ok",
			fields: fields{
				eventstore: expectEventstore(
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
							"user1",
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
			name: "usergrant for project with passed succeeding permission check, ok",
			fields: fields{
				eventstore: expectEventstore(
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
							"user1",
							[]string{"rolekey1", "rolekey2"},
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1", "rolekey2"},
				},
				permissionCheck: succeedingUserGrantPermissionCheck,
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
			name: "usergrant for project with passed failing permission check, error",
			fields: fields{
				eventstore: expectEventstore(
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
				),
			},
			args: args{
				ctx: context.Background(),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1", "rolekey2"},
				},
				permissionCheck: failingUserGrantPermissionCheck,
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
				err: isMockedPermissionCheckErr,
			},
		},
		{
			name: "usergrant roles not changed, ignore unchanged, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1", "rolekey2"}),
						),
					),
				),
			},
			args: args{
				ctx: authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrant: &domain.UserGrant{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "usergrant1",
						ResourceOwner: "org1",
					},
					UserID:    "user1",
					ProjectID: "project1",
					RoleKeys:  []string{"rolekey1", "rolekey2"},
				},
				ignoreUnchanged: true,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.ChangeUserGrant(tt.args.ctx, tt.args.userGrant, tt.args.cascade, tt.args.ignoreUnchanged, tt.args.permissionCheck)
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
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userGrantID   string
		resourceOwner string
		check         UserGrantPermissionCheck
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "not provided resourceOwner, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org").Aggregate,
								"user1",
								"project1",
								"", []string{"rolekey1"}),
						),
					),
					expectPush(
						usergrant.NewUserGrantDeactivatedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:         authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrantID: "usergrant1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org",
				},
			},
		},
		{
			name: "usergrant not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
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
				eventstore: expectEventstore(
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
			name: "no permissions, permission denied error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
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
			name: "already deactivated, ignore, ok",
			fields: fields{
				eventstore: expectEventstore(
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
				want: &domain.ObjectDetails{
					ResourceOwner: "org",
				},
			},
		},
		{
			name: "deactivated, ok",
			fields: fields{
				eventstore: expectEventstore(
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
		{
			name: "with passed succeeding permission check, ok",
			fields: fields{
				eventstore: expectEventstore(
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
				ctx:           context.Background(),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
				check:         succeedingUserGrantPermissionCheck,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "with passed failing permission check, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
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
				check:         failingUserGrantPermissionCheck,
			},
			res: res{
				err: isMockedPermissionCheckErr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.DeactivateUserGrant(tt.args.ctx, tt.args.userGrantID, tt.args.resourceOwner, tt.args.check)
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

func TestCommandSide_ReactivateUserGrant(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userGrantID   string
		resourceOwner string
		check         UserGrantPermissionCheck
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "not provided resourceOwner, ok",
			fields: fields{
				eventstore: expectEventstore(
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
					expectPush(
						usergrant.NewUserGrantReactivatedEvent(context.Background(),
							&usergrant.NewAggregate("usergrant1", "org").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:         authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrantID: "usergrant1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org",
				},
			},
		},
		{
			name: "usergrant not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
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
				eventstore: expectEventstore(
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
			name: "no permissions, permission denied error",
			fields: fields{
				eventstore: expectEventstore(
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
			name: "already active, ignore, ok",
			fields: fields{
				eventstore: expectEventstore(
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
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org",
				},
			},
		},
		{
			name: "reactivated, ok",
			fields: fields{
				eventstore: expectEventstore(
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
		{
			name: "with passed succeeding permission check, ok",
			fields: fields{
				eventstore: expectEventstore(
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
				ctx:           context.Background(),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
				check:         succeedingUserGrantPermissionCheck,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "with passed failing permission check, error",
			fields: fields{
				eventstore: expectEventstore(
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
				),
			},
			args: args{
				ctx:           context.Background(),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
				check:         failingUserGrantPermissionCheck,
			},
			res: res{
				err: isMockedPermissionCheckErr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.ReactivateUserGrant(tt.args.ctx, tt.args.userGrantID, tt.args.resourceOwner, tt.args.check)
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

func TestCommandSide_RemoveUserGrant(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx            context.Context
		userGrantID    string
		resourceOwner  string
		ignoreNotFound bool
		check          UserGrantPermissionCheck
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "usergrant not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
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
			name: "usergrant not existing, ignore, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:            authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrantID:    "usergrant1",
				resourceOwner:  "org1",
				ignoreNotFound: true,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "usergrant removed, not found error",
			fields: fields{
				eventstore: expectEventstore(
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
			name: "no permissions, permission denied error",
			fields: fields{
				eventstore: expectEventstore(
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
				eventstore: expectEventstore(
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
			name: "not provided resourceOwner, ok",
			fields: fields{
				eventstore: expectEventstore(
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
				ctx:         authz.NewMockContextWithPermissions("", "", "", []string{domain.RoleProjectOwner}),
				userGrantID: "usergrant1",
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
				eventstore: expectEventstore(
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
		{
			name: "with passed succeeding permission check, ok",
			fields: fields{
				eventstore: expectEventstore(
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
				ctx:           context.Background(),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
				check:         succeedingUserGrantPermissionCheck,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "with passed failing permission check, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							usergrant.NewUserGrantAddedEvent(context.Background(),
								&usergrant.NewAggregate("usergrant1", "org1").Aggregate,
								"user1",
								"project1",
								"projectgrant1", []string{"rolekey1"}),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userGrantID:   "usergrant1",
				resourceOwner: "org1",
				check:         failingUserGrantPermissionCheck,
			},
			res: res{
				err: isMockedPermissionCheckErr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.RemoveUserGrant(tt.args.ctx, tt.args.userGrantID, tt.args.resourceOwner, tt.args.ignoreNotFound, tt.args.check)
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
				eventstore: eventstoreExpect(t),
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
