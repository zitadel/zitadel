package command

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/member"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

func TestCommandSide_AddIAMMember(t *testing.T) {
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
				ctx:    context.Background(),
				member: &domain.Member{},
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
					UserID: "user1",
					Roles:  []string{"IAM_OWNER"},
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
					UserID: "user1",
					Roles:  []string{"IAM_OWNER"},
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
							iam.NewMemberAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"user1",
							),
						),
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "IAM_OWNER",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					UserID: "user1",
					Roles:  []string{"IAM_OWNER"},
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
							eventFromEventPusher(iam.NewMemberAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"user1",
								[]string{"IAM_OWNER"}...,
							)),
						},
						nil,
						uniqueConstraintsFromEventConstraint(member.NewAddMemberUniqueConstraint("IAM", "user1")),
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "IAM_OWNER",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					UserID: "user1",
					Roles:  []string{"IAM_OWNER"},
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
							eventFromEventPusher(iam.NewMemberAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"user1",
								[]string{"IAM_OWNER"}...,
							)),
						},
						nil,
						uniqueConstraintsFromEventConstraint(member.NewAddMemberUniqueConstraint("IAM", "user1")),
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "IAM_OWNER",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					UserID: "user1",
					Roles:  []string{"IAM_OWNER"},
				},
			},
			res: res{
				want: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "IAM",
						AggregateID:   "IAM",
					},
					UserID: "user1",
					Roles:  []string{"IAM_OWNER"},
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
			got, err := r.AddIAMMember(tt.args.ctx, tt.args.member)
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

func TestCommandSide_ChangeIAMMember(t *testing.T) {
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
				ctx:    context.Background(),
				member: &domain.Member{},
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
					UserID: "user1",
					Roles:  []string{"IAM_OWNER"},
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
						Role: "IAM_OWNER",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					UserID: "user1",
					Roles:  []string{"IAM_OWNER"},
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
							iam.NewMemberAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"user1",
								[]string{"IAM_OWNER"}...,
							),
						),
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleIAMOwner,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					UserID: "user1",
					Roles:  []string{"IAM_OWNER"},
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
							iam.NewMemberAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"user1",
								[]string{"IAM_OWNER"}...,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(iam.NewMemberChangedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"user1",
								[]string{"IAM_OWNER", "IAM_OWNER_VIEWER"}...,
							)),
						},
						nil,
					),
				),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "IAM_OWNER",
					},
					{
						Role: "IAM_OWNER_VIEWER",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				member: &domain.Member{
					UserID: "user1",
					Roles:  []string{"IAM_OWNER", "IAM_OWNER_VIEWER"},
				},
			},
			res: res{
				want: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						ResourceOwner: "IAM",
						AggregateID:   "IAM",
					},
					UserID: "user1",
					Roles:  []string{"IAM_OWNER", "IAM_OWNER_VIEWER"},
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
			got, err := r.ChangeIAMMember(tt.args.ctx, tt.args.member)
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

func TestCommandSide_RemoveIAMMember(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		userID string
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
			name: "invalid member userid missing, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "",
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
				ctx:    context.Background(),
				userID: "user1",
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
							iam.NewMemberAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"user1",
								[]string{"IAM_OWNER"}...,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(iam.NewMemberRemovedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"user1",
							)),
						},
						nil,
						uniqueConstraintsFromEventConstraint(member.NewRemoveMemberUniqueConstraint("IAM", "user1")),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveIAMMember(tt.args.ctx, tt.args.userID)
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
