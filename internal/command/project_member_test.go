package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddProjectMember(t *testing.T) {
	type fields struct {
		eventstore   func(t *testing.T) *eventstore.Eventstore
		zitadelRoles []authz.RoleMapping
	}
	type args struct {
		member *AddProjectMember
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
			name: "invalid member, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				member: &AddProjectMember{
					ResourceOwner: "org1",
					ProjectID:     "project1",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid roles, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				member: &AddProjectMember{
					ResourceOwner: "org1",
					ProjectID:     "project1",
					UserID:        "user1",
					Roles:         []string{"PROJECT_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleProjectOwner,
					},
				},
			},
			args: args{
				member: &AddProjectMember{
					ResourceOwner: "org1",
					ProjectID:     "project1",
					UserID:        "user1",
					Roles:         []string{"PROJECT_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "member already exists, precondition error",
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
					),
					expectFilter(
						eventFromEventPusher(
							project.NewProjectMemberAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"user1",
							),
						),
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleProjectOwner,
					},
				},
			},
			args: args{
				member: &AddProjectMember{
					ResourceOwner: "org1",
					ProjectID:     "project1",
					UserID:        "user1",
					Roles:         []string{"PROJECT_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "member add uniqueconstraint err, already exists",
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
					),
					expectFilter(),
					expectPushFailed(zerrors.ThrowAlreadyExists(nil, "ERROR", "internal"),
						project.NewProjectMemberAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"user1",
							[]string{"PROJECT_OWNER"}...,
						),
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleProjectOwner,
					},
				},
			},
			args: args{
				member: &AddProjectMember{
					ResourceOwner: "org1",
					ProjectID:     "project1",
					UserID:        "user1",
					Roles:         []string{"PROJECT_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "member add, ok",
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
					),
					expectFilter(),
					expectPush(
						project.NewProjectMemberAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"user1",
							[]string{"PROJECT_OWNER"}...,
						),
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleProjectOwner,
					},
				},
			},
			args: args{
				member: &AddProjectMember{
					ResourceOwner: "org1",
					ProjectID:     "project1",
					UserID:        "user1",
					Roles:         []string{"PROJECT_OWNER"},
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
				eventstore:   tt.fields.eventstore(t),
				zitadelRoles: tt.fields.zitadelRoles,
			}
			got, err := r.AddProjectMember(context.Background(), tt.args.member)
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

func TestCommandSide_ChangeProjectMember(t *testing.T) {
	type fields struct {
		eventstore   func(t *testing.T) *eventstore.Eventstore
		zitadelRoles []authz.RoleMapping
	}
	type args struct {
		member *ChangeProjectMember
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
			name: "invalid member, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				member: &ChangeProjectMember{
					ResourceOwner: "org1",
					ProjectID:     "project1",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid roles, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				member: &ChangeProjectMember{
					ResourceOwner: "org1",
					ProjectID:     "project1",
					UserID:        "user1",
					Roles:         []string{"PROJECT_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "member not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleProjectOwner,
					},
				},
			},
			args: args{
				member: &ChangeProjectMember{
					ResourceOwner: "org1",
					ProjectID:     "project1",
					UserID:        "user1",
					Roles:         []string{"PROJECT_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "member not changed, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectMemberAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"user1",
								[]string{"PROJECT_OWNER"}...,
							),
						),
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleProjectOwner,
					},
				},
			},
			args: args{
				member: &ChangeProjectMember{
					ResourceOwner: "org1",
					ProjectID:     "project1",
					UserID:        "user1",
					Roles:         []string{"PROJECT_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "member change, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectMemberAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"user1",
								[]string{"PROJECT_OWNER"}...,
							),
						),
					),
					expectPush(
						project.NewProjectMemberChangedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"user1",
							[]string{"PROJECT_OWNER", "PROJECT_VIEWER"}...,
						),
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleProjectOwner,
					},
					{
						Role: "PROJECT_VIEWER",
					},
				},
			},
			args: args{
				member: &ChangeProjectMember{
					ResourceOwner: "org1",
					ProjectID:     "project1",
					UserID:        "user1",
					Roles:         []string{"PROJECT_OWNER", "PROJECT_VIEWER"},
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
				eventstore:   tt.fields.eventstore(t),
				zitadelRoles: tt.fields.zitadelRoles,
			}
			got, err := r.ChangeProjectMember(context.Background(), tt.args.member)
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

func TestCommandSide_RemoveProjectMember(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		projectID     string
		userID        string
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
			name: "invalid member projectid missing, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				projectID:     "",
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid member userid missing, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				projectID:     "project1",
				userID:        "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "member not existing, empty object details result",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				projectID:     "project1",
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{},
			},
		},
		{
			name: "member remove, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectMemberAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"user1",
								[]string{"PROJECT_OWNER"}...,
							),
						),
					),
					expectPush(
						project.NewProjectMemberRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"user1",
						),
					),
				),
			},
			args: args{
				projectID:     "project1",
				userID:        "user1",
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.RemoveProjectMember(context.Background(), tt.args.projectID, tt.args.userID, tt.args.resourceOwner)
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
