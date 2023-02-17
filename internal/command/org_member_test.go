package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/member"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestAddMember(t *testing.T) {
	type args struct {
		a            *org.Aggregate
		userID       string
		roles        []string
		zitadelRoles []authz.RoleMapping
		filter       preparation.FilterToQueryReducer
	}

	ctx := context.Background()
	agg := org.NewAggregate("test")

	tests := []struct {
		name string
		args args
		want Want
	}{
		{
			name: "no user id",
			args: args{
				a:      agg,
				userID: "",
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "ORG-4Mlfs", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "no roles",
			args: args{
				a:      agg,
				userID: "12342",
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "V2-PfYhb", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "TODO: invalid roles",
			args: args{
				a:      agg,
				userID: "123",
				roles:  []string{"ORG_OWNER"},
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "Org-4N8es", ""),
			},
		},
		{
			name: "user not exists",
			args: args{
				a:      agg,
				userID: "userID",
				roles:  []string{"ORG_OWNER"},
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "ORG_OWNER",
					},
				},
				filter: NewMultiFilter().Append(
					func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).Filter(),
			},
			want: Want{
				CreateErr: errors.ThrowPreconditionFailed(nil, "ORG-GoXOn", "Errors.User.NotFound"),
			},
		},
		{
			name: "already member",
			args: args{
				a:      agg,
				userID: "userID",
				roles:  []string{"ORG_OWNER"},
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "ORG_OWNER",
					},
				},
				filter: NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							user.NewMachineAddedEvent(
								ctx,
								&user.NewAggregate("id", "ro").Aggregate,
								"userName",
								"name",
								"description",
								true,
								domain.OIDCTokenTypeBearer,
							),
						}, nil
					}).
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							org.NewMemberAddedEvent(
								ctx,
								&org.NewAggregate("id").Aggregate,
								"userID",
							),
						}, nil
					}).
					Filter(),
			},
			want: Want{
				CreateErr: errors.ThrowAlreadyExists(nil, "ORG-poWwe", "Errors.Org.Member.AlreadyExists"),
			},
		},
		{
			name: "correct",
			args: args{
				a:      agg,
				userID: "userID",
				roles:  []string{"ORG_OWNER"},
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "ORG_OWNER",
					},
				},
				filter: NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							user.NewMachineAddedEvent(
								ctx,
								&user.NewAggregate("id", "ro").Aggregate,
								"userName",
								"name",
								"description",
								true,
								domain.OIDCTokenTypeBearer,
							),
						}, nil
					}).
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).
					Filter(),
			},
			want: Want{
				Commands: []eventstore.Command{
					org.NewMemberAddedEvent(ctx, &agg.Aggregate, "userID", "ORG_OWNER"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertValidation(t, context.Background(), (&Commands{zitadelRoles: tt.args.zitadelRoles}).AddOrgMemberCommand(tt.args.a, tt.args.userID, tt.args.roles...), tt.args.filter, tt.want)
		})
	}
}

func TestIsMember(t *testing.T) {
	type args struct {
		filter preparation.FilterToQueryReducer
		orgID  string
		userID string
	}
	tests := []struct {
		name       string
		args       args
		wantExists bool
		wantErr    bool
	}{
		{
			name: "no events",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{}, nil
				},
				orgID:  "orgID",
				userID: "userID",
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name: "member added",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						org.NewMemberAddedEvent(
							context.Background(),
							&org.NewAggregate("orgID").Aggregate,
							"userID",
						),
					}, nil
				},
				orgID:  "orgID",
				userID: "userID",
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name: "member removed",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						org.NewMemberAddedEvent(
							context.Background(),
							&org.NewAggregate("orgID").Aggregate,
							"userID",
						),
						org.NewMemberRemovedEvent(
							context.Background(),
							&org.NewAggregate("orgID").Aggregate,
							"userID",
						),
					}, nil
				},
				orgID:  "orgID",
				userID: "userID",
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name: "member cascade removed",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						org.NewMemberAddedEvent(
							context.Background(),
							&org.NewAggregate("orgID").Aggregate,
							"userID",
						),
						org.NewMemberCascadeRemovedEvent(
							context.Background(),
							&org.NewAggregate("orgID").Aggregate,
							"userID",
						),
					}, nil
				},
				orgID:  "orgID",
				userID: "userID",
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name: "error durring filter",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return nil, errors.ThrowInternal(nil, "PROJE-Op26p", "Errors.Internal")
				},
				orgID:  "orgID",
				userID: "userID",
			},
			wantExists: false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExists, err := IsOrgMember(context.Background(), tt.args.filter, tt.args.orgID, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExistsUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExists != tt.wantExists {
				t.Errorf("ExistsUser() = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}

func TestCommandSide_AddOrgMember(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		zitadelRoles []authz.RoleMapping
	}
	type args struct {
		ctx    context.Context
		userID string
		orgID  string
		roles  []string
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
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
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
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				roles:  []string{"ORG_OWNER"},
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
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
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				roles:  []string{domain.RoleOrgOwner},
			},
			res: res{
				err: errors.IsPreconditionFailed,
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
								&org.NewAggregate("org1").Aggregate,
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
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				roles:  []string{"ORG_OWNER"},
			},
			res: res{
				err: errors.IsErrorAlreadyExists,
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
					expectPushFailed(errors.ThrowAlreadyExists(nil, "ERROR", "internal"),
						[]*repository.Event{
							eventFromEventPusher(org.NewMemberAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"user1",
								[]string{"ORG_OWNER"}...,
							)),
						},
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
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				roles:  []string{"ORG_OWNER"},
			},
			res: res{
				err: errors.IsErrorAlreadyExists,
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
								&org.NewAggregate("org1").Aggregate,
								"user1",
								[]string{"ORG_OWNER"}...,
							)),
						},
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
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				roles:  []string{"ORG_OWNER"},
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
			got, err := r.AddOrgMember(tt.args.ctx, tt.args.orgID, tt.args.userID, tt.args.roles...)
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
				err: errors.IsErrorInvalidArgument,
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
				err: errors.IsErrorInvalidArgument,
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
				err: errors.IsNotFound,
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
								&org.NewAggregate("org1").Aggregate,
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
				err: errors.IsPreconditionFailed,
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
								&org.NewAggregate("org1").Aggregate,
								"user1",
								[]string{"ORG_OWNER"}...,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewMemberChangedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"user1",
								[]string{"ORG_OWNER", "ORG_OWNER_VIEWER"}...,
							)),
						},
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
				err: errors.IsErrorInvalidArgument,
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
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "member not existing, empty object details result",
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
				want: &domain.ObjectDetails{},
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
