package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddInstanceMember(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		zitadelRoles    []authz.RoleMapping
		checkPermission domain.PermissionCheck
	}
	type args struct {
		member *AddInstanceMember
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
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				member: &AddInstanceMember{},
			},
			res: res{
				err: zerrors.IsInternal,
			},
		},
		{
			name: "invalid roles, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				member: &AddInstanceMember{
					InstanceID: "INSTANCE",
					UserID:     "user1",
					Roles:      []string{"IAM_OWNER"},
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
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "IAM_OWNER",
					},
				},
			},
			args: args{
				member: &AddInstanceMember{
					InstanceID: "INSTANCE",
					UserID:     "user1",
					Roles:      []string{"IAM_OWNER"},
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
							instance.NewMemberAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"user1",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "IAM_OWNER",
					},
				},
			},
			args: args{
				member: &AddInstanceMember{
					InstanceID: "INSTANCE",
					UserID:     "user1",
					Roles:      []string{"IAM_OWNER"},
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
						instance.NewMemberAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"user1",
							[]string{"IAM_OWNER"}...,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "IAM_OWNER",
					},
				},
			},
			args: args{
				member: &AddInstanceMember{
					InstanceID: "INSTANCE",
					UserID:     "user1",
					Roles:      []string{"IAM_OWNER"},
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
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "IAM_OWNER",
					},
				},
			},
			args: args{
				member: &AddInstanceMember{
					InstanceID: "INSTANCE",
					UserID:     "user1",
					Roles:      []string{"IAM_OWNER"},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "member add, no permission",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckNotAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "IAM_OWNER",
					},
				},
			},
			args: args{
				member: &AddInstanceMember{
					InstanceID: "INSTANCE",
					UserID:     "user1",
					Roles:      []string{"IAM_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				zitadelRoles:    tt.fields.zitadelRoles,
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.AddInstanceMember(context.Background(), tt.args.member)
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

func TestCommandSide_ChangeInstanceMember(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		zitadelRoles    []authz.RoleMapping
		checkPermission domain.PermissionCheck
	}
	type args struct {
		member *ChangeInstanceMember
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
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				member: &ChangeInstanceMember{},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid roles, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				member: &ChangeInstanceMember{
					InstanceID: "INSTANCE",
					UserID:     "user1",
					Roles:      []string{"IAM_OWNER"},
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
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: "IAM_OWNER",
					},
				},
			},
			args: args{
				member: &ChangeInstanceMember{
					InstanceID: "INSTANCE",
					UserID:     "user1",
					Roles:      []string{"IAM_OWNER"},
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
							instance.NewMemberAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"user1",
								[]string{"IAM_OWNER"}...,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleIAMOwner,
					},
				},
			},
			args: args{
				member: &ChangeInstanceMember{
					InstanceID: "INSTANCE",
					UserID:     "user1",
					Roles:      []string{"IAM_OWNER"},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "member change, ok",
			fields: fields{
				eventstore: expectEventstore(
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
				checkPermission: newMockPermissionCheckAllowed(),
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
				member: &ChangeInstanceMember{
					InstanceID: "INSTANCE",
					UserID:     "user1",
					Roles:      []string{"IAM_OWNER", "IAM_OWNER_VIEWER"},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "member change, ok",
			fields: fields{
				eventstore: expectEventstore(
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
				checkPermission: newMockPermissionCheckNotAllowed(),
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
				member: &ChangeInstanceMember{
					InstanceID: "INSTANCE",
					UserID:     "user1",
					Roles:      []string{"IAM_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				zitadelRoles:    tt.fields.zitadelRoles,
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.ChangeInstanceMember(context.Background(), tt.args.member)
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

func TestCommandSide_RemoveInstanceMember(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
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
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				instanceID: "INSTANCE",
				userID:     "",
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				instanceID: "INSTANCE",
				userID:     "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "member remove, ok",
			fields: fields{
				eventstore: expectEventstore(
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				instanceID: "INSTANCE",
				userID:     "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "member remove, no permission",
			fields: fields{
				eventstore: expectEventstore(
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
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				instanceID: "INSTANCE",
				userID:     "user1",
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.RemoveInstanceMember(context.Background(), tt.args.instanceID, tt.args.userID)
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
