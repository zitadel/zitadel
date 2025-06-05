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
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddIAMMember(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		zitadelRoles []authz.RoleMapping
	}
	type args struct {
		instanceID string
		userID     string
		roles      []string
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
			args: args{},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
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
				userID: "user1",
				roles:  []string{"IAM_OWNER"},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
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
						Role: "IAM_OWNER",
					},
				},
			},
			args: args{
				instanceID: "",
				userID:     "user1",
				roles:      []string{"IAM_OWNER"},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
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
							instance.NewMemberAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
				instanceID: "",
				userID:     "user1",
				roles:      []string{"IAM_OWNER"},
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
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
					expectPushFailed(zerrors.ThrowAlreadyExists(nil, "ERROR", "internal"),
						instance.NewMemberAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"user1",
							[]string{"IAM_OWNER"}...,
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
				instanceID: "INSTANCE",
				userID:     "user1",
				roles:      []string{"IAM_OWNER"},
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "member add, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
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
						instance.NewMemberAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"user1",
							[]string{"IAM_OWNER"}...,
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
				instanceID: "INSTANCE",
				userID:     "user1",
				roles:      []string{"IAM_OWNER"},
			},
			res: res{
				want: &domain.Member{
					ObjectRoot: models.ObjectRoot{
						InstanceID:    "INSTANCE",
						ResourceOwner: "INSTANCE",
						AggregateID:   "INSTANCE",
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
			got, err := r.AddInstanceMember(context.Background(), tt.args.instanceID, tt.args.userID, tt.args.roles...)
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
		ctx        context.Context
		instanceID string
		member     *domain.Member
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "member not changed, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewMemberAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "member change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewMemberAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"user1",
								[]string{"IAM_OWNER"}...,
							),
						),
					),
					expectPush(
						instance.NewMemberChangedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"user1",
							[]string{"IAM_OWNER", "IAM_OWNER_VIEWER"}...,
						),
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
						ResourceOwner: "INSTANCE",
						AggregateID:   "INSTANCE",
						InstanceID:    "INSTANCE",
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
			got, err := r.ChangeInstanceMember(tt.args.ctx, tt.args.member)
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
		ctx        context.Context
		instanceID string
		userID     string
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
				err: zerrors.IsErrorInvalidArgument,
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
				ctx:    context.Background(),
				userID: "user1",
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
							instance.NewMemberAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"user1",
								[]string{"IAM_OWNER"}...,
							),
						),
					),
					expectPush(
						instance.NewMemberRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"user1",
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveInstanceMember(tt.args.ctx, tt.args.userID)
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
