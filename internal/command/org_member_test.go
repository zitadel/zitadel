package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/member"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestCommandSide_AddOrgMember(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		zitadelRoles []authz.RoleMapping
	}
	type args struct {
		ctx    context.Context
		member *domain.Member
	}
	type res struct {
		want *domain.Member
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					UserID: "user1",
					Roles:  []string{domain.RoleOrgOwner},
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid roles, error",
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
					),
				),
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					UserID: "user1",
					Roles:  []string{"ORG_OWNER"},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "member already exists, precondition error",
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
					),
					expectFilter(
						eventFromEventPusher(
							org.NewMemberAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"user1",
							),
						),
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleOrgOwner,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					UserID: "user1",
					Roles:  []string{"ORG_OWNER"},
				},
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "member add uniqueconstraint err, already exists",
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
					),
					expectFilter(),
					expectPushFailed(caos_errs.ThrowAlreadyExists(nil, "ERROR", "internal"),
						[]*repository.Event{
							eventFromEventPusher(org.NewMemberAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"user1",
								[]string{"ORG_OWNER"}...,
							)),
						},
						nil,
						uniqueConstraintsFromEventConstraint(member.NewAddMemberUniqueConstraint("org1", "user1")),
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleOrgOwner,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					UserID: "user1",
					Roles:  []string{"ORG_OWNER"},
				},
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "member add, ok",
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
					),
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewMemberAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"user1",
								[]string{"ORG_OWNER"}...,
							)),
						},
						nil,
						uniqueConstraintsFromEventConstraint(member.NewAddMemberUniqueConstraint("org1", "user1")),
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleOrgOwner,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					UserID: "user1",
					Roles:  []string{"ORG_OWNER"},
				},
			},
			res: res{
				want: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
						AggregateID:   "org1",
					},
					UserID: "user1",
					Roles:  []string{domain.RoleOrgOwner},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:   tt.fields.eventstore,
				zitadelRoles: tt.fields.zitadelRoles,
			}
			got, err := r.AddOrgMember(tt.args.ctx, tt.args.member)
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

func TestCommandSide_ChangeOrgMember(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		zitadelRoles []authz.RoleMapping
	}
	type args struct {
		ctx    context.Context
		member *domain.Member
	}
	type res struct {
		want *domain.Member
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid roles, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					UserID: "user1",
					Roles:  []string{"PROJECT_OWNER"},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "member not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleOrgOwner,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					UserID: "user1",
					Roles:  []string{"ORG_OWNER"},
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "member not changed, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewMemberAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"user1",
								[]string{"ORG_OWNER"}...,
							),
						),
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleOrgOwner,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					UserID: "user1",
					Roles:  []string{"ORG_OWNER"},
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewMemberAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"user1",
								[]string{"ORG_OWNER"}...,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewMemberChangedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"user1",
								[]string{"ORG_OWNER", "ORG_OWNER_VIEWER"}...,
							)),
						},
						nil,
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "ORG_OWNER",
					},
					{
						Role: "ORG_OWNER_VIEWER",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					UserID: "user1",
					Roles:  []string{"ORG_OWNER", "ORG_OWNER_VIEWER"},
				},
			},
			res: res{
				want: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "org1",
						AggregateID:   "org1",
					},
					UserID: "user1",
					Roles:  []string{"ORG_OWNER", "ORG_OWNER_VIEWER"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:   tt.fields.eventstore,
				zitadelRoles: tt.fields.zitadelRoles,
			}
			got, err := r.ChangeOrgMember(tt.args.ctx, tt.args.member)
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

func TestCommandSide_RemoveOrgMember(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "",
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid member userid missing, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				userID:        "",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "member not existing, nil result",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				want: nil,
			},
		},
		{
			name: "member remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
						[]*repository.Event{
							eventFromEventPusher(project.NewProjectMemberRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"user1",
							)),
						},
						nil,
						uniqueConstraintsFromEventConstraint(member.NewRemoveMemberUniqueConstraint("project1", "user1")),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
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
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveProjectMember(tt.args.ctx, tt.args.projectID, tt.args.userID, tt.args.resourceOwner)
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
