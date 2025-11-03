package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestAddMember(t *testing.T) {
	type args struct {
		member       *AddOrgMember
		zitadelRoles []authz.RoleMapping
		filter       preparation.FilterToQueryReducer
	}

	ctx := context.Background()

	tests := []struct {
		name string
		args args
		want Want
	}{
		{
			name: "no user id",
			args: args{
				member: &AddOrgMember{
					OrgID:  "test",
					UserID: "",
				},
			},
			want: Want{
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "ORG-4Mlfs", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "no roles",
			args: args{
				member: &AddOrgMember{
					OrgID:  "test",
					UserID: "12342",
				},
			},
			want: Want{
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "ORG-4Mlfs", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "TODO: invalid roles",
			args: args{
				member: &AddOrgMember{
					OrgID:  "test",
					UserID: "12342",
					Roles:  []string{"ORG_OWNER"},
				},
			},
			want: Want{
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "Org-4N8es", ""),
			},
		},
		{
			name: "user not exists",
			args: args{
				member: &AddOrgMember{
					OrgID:  "test",
					UserID: "userID",
					Roles:  []string{"ORG_OWNER"},
				},
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
				CreateErr: zerrors.ThrowPreconditionFailed(nil, "ORG-GoXOn", "Errors.User.NotFound"),
			},
		},
		{
			name: "already member",
			args: args{
				member: &AddOrgMember{
					OrgID:  "test",
					UserID: "userID",
					Roles:  []string{"ORG_OWNER"},
				},
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
				CreateErr: zerrors.ThrowAlreadyExists(nil, "ORG-poWwe", "Errors.Org.Member.AlreadyExists"),
			},
		},
		{
			name: "correct",
			args: args{
				member: &AddOrgMember{
					OrgID:  "test",
					UserID: "userID",
					Roles:  []string{"ORG_OWNER"},
				},
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
					org.NewMemberAddedEvent(ctx, &org.NewAggregate("test").Aggregate, "userID", "ORG_OWNER"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertValidation(t, context.Background(), (&Commands{zitadelRoles: tt.args.zitadelRoles}).AddOrgMemberCommand(tt.args.member), tt.args.filter, tt.want)
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
					return nil, zerrors.ThrowInternal(nil, "PROJE-Op26p", "Errors.Internal")
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
		checkPermission domain.PermissionCheck
		eventstore      func(t *testing.T) *eventstore.Eventstore
		zitadelRoles    []authz.RoleMapping
	}
	type args struct {
		member *AddOrgMember
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
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				member: &AddOrgMember{
					OrgID: "org1",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid roles, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				member: &AddOrgMember{
					OrgID:  "org1",
					UserID: "user1",
					Roles:  []string{"ORG_OWNER"},
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
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org",
							),
						),
					),
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleOrgOwner,
					},
				},
			},
			args: args{
				member: &AddOrgMember{
					OrgID:  "org1",
					UserID: "user1",
					Roles:  []string{"ORG_OWNER"},
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
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org",
							),
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
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleOrgOwner,
					},
				},
			},
			args: args{
				member: &AddOrgMember{
					OrgID:  "org1",
					UserID: "user1",
					Roles:  []string{"ORG_OWNER"},
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
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org",
							),
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
					),
					expectFilter(),
					expectPushFailed(zerrors.ThrowAlreadyExists(nil, "ERROR", "internal"),
						org.NewMemberAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"user1",
							[]string{"ORG_OWNER"}...,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleOrgOwner,
					},
				},
			},
			args: args{
				member: &AddOrgMember{
					OrgID:  "org1",
					UserID: "user1",
					Roles:  []string{"ORG_OWNER"},
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
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org",
							),
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
					),
					expectFilter(),
					expectPush(
						org.NewMemberAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"user1",
							[]string{"ORG_OWNER"}...,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleOrgOwner,
					},
				},
			},
			args: args{
				member: &AddOrgMember{
					OrgID:  "org1",
					UserID: "user1",
					Roles:  []string{"ORG_OWNER"},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "member add, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"org",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleOrgOwner,
					},
				},
			},
			args: args{
				member: &AddOrgMember{
					OrgID:  "org1",
					UserID: "user1",
					Roles:  []string{"ORG_OWNER"},
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
			got, err := r.AddOrgMember(context.Background(), tt.args.member)
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

func TestCommandSide_ChangeOrgMember(t *testing.T) {
	type fields struct {
		checkPermission domain.PermissionCheck
		eventstore      func(t *testing.T) *eventstore.Eventstore
		zitadelRoles    []authz.RoleMapping
	}
	type args struct {
		member *ChangeOrgMember
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
				member: &ChangeOrgMember{
					OrgID: "org1",
				},
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
				member: &ChangeOrgMember{
					OrgID:  "org1",
					UserID: "user1",
					Roles:  []string{"PROJECT_OWNER"},
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
						Role: domain.RoleOrgOwner,
					},
				},
			},
			args: args{
				member: &ChangeOrgMember{
					OrgID:  "org1",
					UserID: "user1",
					Roles:  []string{"ORG_OWNER"},
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "member not changed, no change",
			fields: fields{
				eventstore: expectEventstore(
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
				checkPermission: newMockPermissionCheckAllowed(),
				zitadelRoles: []authz.RoleMapping{
					{
						Role: domain.RoleOrgOwner,
					},
				},
			},
			args: args{
				member: &ChangeOrgMember{
					OrgID:  "org1",
					UserID: "user1",
					Roles:  []string{"ORG_OWNER"},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "member change, ok",
			fields: fields{
				eventstore: expectEventstore(
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
						org.NewMemberChangedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"user1",
							[]string{"ORG_OWNER", "ORG_OWNER_VIEWER"}...,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
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
				member: &ChangeOrgMember{
					OrgID:  "org1",
					UserID: "user1",
					Roles:  []string{"ORG_OWNER", "ORG_OWNER_VIEWER"},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "member change, no permission",
			fields: fields{
				eventstore: expectEventstore(
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
				checkPermission: newMockPermissionCheckNotAllowed(),
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
				member: &ChangeOrgMember{
					OrgID:  "org1",
					UserID: "user1",
					Roles:  []string{"ORG_OWNER", "ORG_OWNER_VIEWER"},
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
			got, err := r.ChangeOrgMember(context.Background(), tt.args.member)
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

func TestCommandSide_RemoveOrgMember(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx    context.Context
		orgID  string
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
			name: "invalid member orgID missing, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "",
				userID: "user1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid member userid missing, error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "",
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
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "member remove, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewMemberAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"user1",
								[]string{"PROJECT_OWNER"}...,
							),
						),
					),
					expectPush(
						org.NewMemberRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"user1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "member remove, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewMemberAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"user1",
								[]string{"PROJECT_OWNER"}...,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
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
			got, err := r.RemoveOrgMember(tt.args.ctx, tt.args.orgID, tt.args.userID)
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
