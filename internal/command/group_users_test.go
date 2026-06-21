package command

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_AddUsersToGroup(t *testing.T) {
	t.Parallel()
	filterErr := errors.New("filter error")
	pushErr := errors.New("push error")

	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		groupID string
		userIDs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr func(error) bool
	}{
		{
			name: "failed to get group write model, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(filterErr),
				),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, filterErr)
			},
		},
		{
			name: "group not found, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			wantErr: zerrors.IsPreconditionFailed,
		},
		{
			name: "missing permission, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					)),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			wantErr: zerrors.IsPermissionDenied,
		},
		{
			name: "user not found, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					),
					expectFilter(), // to get the user write model for user1
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			wantErr: zerrors.IsPreconditionFailed,
		},
		{
			name: "multiple users not found, all missing users reported, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					),
					expectFilter( // users existence check: only user2 exists
						eventFromEventPusher(
							addNewUserEvent("user2", "org1"),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2", "user3"},
			},
			wantErr: func(err error) bool {
				return zerrors.IsPreconditionFailed(err) &&
					strings.Contains(err.Error(), "user1") &&
					strings.Contains(err.Error(), "user3") &&
					!strings.Contains(err.Error(), "user2")
			},
		},
		{
			// expectFilterScoped asserts the users-existence filter is scoped
			// to the group's resource owner. Without that scope, a user in
			// another organization would be considered to exist, letting
			// AddUsersToGroup silently admit cross-org members. The mock
			// returns no events here to simulate the eventstore-filtered
			// view from org1's perspective when user1 only exists in org2.
			name: "user from different org treated as not found, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					),
					expectFilterScoped(t, "org1"),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1"},
			},
			wantErr: func(err error) bool {
				return zerrors.IsPreconditionFailed(err) &&
					strings.Contains(err.Error(), "user1")
			},
		},
		{
			// Positive scoping case: the users-existence filter must be
			// scoped to the group's resource owner even when the users do
			// exist. Together with the case above, this locks down the
			// scoping contract at the eventstore mock layer — replacing the
			// previous Test_usersExistenceWriteModel_QueryFiltersByResourceOwner
			// query-shape unit test with a behavior-level assertion.
			name: "users-existence filter is scoped to group's resource owner, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(addNewGroupEvent("group1", "org1")),
					),
					expectFilterScoped(t, "org1",
						eventFromEventPusher(addNewUserEvent("user1", "org1")),
					),
					expectPush(
						newGroupUserAddedEvent("group1", "org1", "user1"),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1"},
			},
			want: &domain.ObjectDetails{
				ID:            "group1",
				ResourceOwner: "org1",
			},
		},
		{
			name: "user removed from organization, not found, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					),
					expectFilter( // users existence check: user1 was removed
						eventFromEventPusher(
							addNewUserEvent("user1", "org1"),
						),
						eventFromEventPusher(
							user.NewUserRemovedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								nil,
								false,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1"},
			},
			wantErr: zerrors.IsPreconditionFailed,
		},
		{
			name: "some users already exist in the group, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter( // to get the group write model
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
						eventFromEventPusher(
							newGroupUserAddedEvent("group1", "org1", "user1"),
						),
					),
					expectFilter( // to get the user write model to check if user1 exists
						eventFromEventPusher(
							addNewUserEvent("user2", "org1"),
						),
					),
					expectPush(
						newGroupUserAddedEvent("group1", "org1", "user2"),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			want: &domain.ObjectDetails{
				ID:            "group1",
				ResourceOwner: "org1",
			},
		},
		{
			name: "all users already exist in the group, no events pushed, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter( // to get the group write model
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
						eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "user1")),
						eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "user2")),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			want: &domain.ObjectDetails{
				ID:            "group1",
				ResourceOwner: "org1",
			},
		},
		{
			name: "failed to push events, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter( // to get the group write model
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					),
					expectFilter( // to get the user write model for user1
						eventFromEventPusher(
							addNewUserEvent("user1", "org1"),
						),
					),
					expectPushFailed(
						pushErr,
						newGroupUserAddedEvent("group1", "org1", "user1"),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1"},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, pushErr)
			},
		},
		{
			name: "all users added, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter( // to get the group write model
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					),
					expectFilter( // users existence check
						eventFromEventPusher(
							addNewUserEvent("user1", "org1"),
						),
						eventFromEventPusher(
							addNewUserEvent("user2", "org1"),
						),
					),
					expectPush(
						newGroupUserAddedEvent("group1", "org1", "user1"),
						newGroupUserAddedEvent("group1", "org1", "user2"),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			want: &domain.ObjectDetails{
				ID:            "group1",
				ResourceOwner: "org1",
			},
		},
		{
			name: "all users added (with duplicate users), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter( // to get the group write model
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					),
					expectFilter( // users existence check
						eventFromEventPusher(
							addNewUserEvent("user1", "org1"),
						),
						eventFromEventPusher(
							addNewUserEvent("user2", "org1"),
						),
					),
					expectPush(
						newGroupUserAddedEvent("group1", "org1", "user1"),
						newGroupUserAddedEvent("group1", "org1", "user2"),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2", "user2", "user1"},
			},
			want: &domain.ObjectDetails{
				ID:            "group1",
				ResourceOwner: "org1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := c.AddUsersToGroup(context.Background(), tt.args.groupID, tt.args.userIDs)
			if tt.wantErr != nil {
				require.True(t, tt.wantErr(err))
				return
			}
			assertObjectDetails(t, tt.want, got)
		})
	}

}

func addNewGroupEvent(groupID, orgID string) *group.GroupAddedEvent {
	return group.NewGroupAddedEvent(context.Background(),
		&group.NewAggregate(groupID, orgID).Aggregate,
		groupID,
		"group description",
	)
}

func TestCommands_RemoveUsersFromGroup(t *testing.T) {
	t.Parallel()
	filterErr := errors.New("filter error")
	pushErr := errors.New("push error")

	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		groupID string
		userIDs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr func(error) bool
	}{
		{
			name: "failed to get group write model, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(filterErr),
				),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, filterErr)
			},
		},
		{
			name: "group not found, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			wantErr: zerrors.IsPreconditionFailed,
		},
		{
			name: "missing permission, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					)),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			wantErr: zerrors.IsPermissionDenied,
		},
		{
			name: "remove users, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter( // to get the group write model
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
						eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "user1")),
						eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "user2")),
					),
					expectPush(
						newGroupUserRemovedEvent("group1", "org1", "user1"),
						newGroupUserRemovedEvent("group1", "org1", "user2"),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			want: &domain.ObjectDetails{
				ResourceOwner: "org1",
				ID:            "group1",
			},
		},
		{
			name: "remove users (with duplicates), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter( // to get the group write model
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
						eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "user1")),
						eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "user2")),
					),
					expectPush(
						newGroupUserRemovedEvent("group1", "org1", "user1"),
						newGroupUserRemovedEvent("group1", "org1", "user2"),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2", "user1", "user2", "user1", "user2"},
			},
			want: &domain.ObjectDetails{
				ResourceOwner: "org1",
				ID:            "group1",
			},
		},
		{
			name: "some users not in group, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter( // to get the group write model
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
						eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "user1")),
						eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "user2")),
					),
					expectPush(
						newGroupUserRemovedEvent("group1", "org1", "user2"),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user2", "user3"},
			},
			want: &domain.ObjectDetails{
				ResourceOwner: "org1",
				ID:            "group1",
			},
		},
		{
			name: "all users not in group, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter( // to get the group write model
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
						eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "user1")),
						eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "user2")),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user3", "user4"},
			},
			want: &domain.ObjectDetails{
				ResourceOwner: "org1",
				ID:            "group1",
			},
		},
		{
			name: "failed to push events, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter( // to get the group write model
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
						eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "user1")),
						eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "user2")),
					),
					expectPushFailed(
						pushErr,
						newGroupUserRemovedEvent("group1", "org1", "user1"),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1"},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, pushErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := c.RemoveUsersFromGroup(context.Background(), tt.args.groupID, tt.args.userIDs)
			if tt.wantErr != nil {
				require.True(t, tt.wantErr(err))
				return
			}
			assertObjectDetails(t, tt.want, got)
		})
	}
}

func TestCommands_removeUserFromGroups(t *testing.T) {
	t.Parallel()

	const (
		userID = "user1"
		orgID  = "org1"
	)

	tests := []struct {
		name     string
		groupIDs []string
	}{
		{name: "no groups, no events", groupIDs: nil},
		{name: "single group, single event", groupIDs: []string{"group1"}},
		{name: "multiple groups, one event per group, in order", groupIDs: []string{"group1", "group2", "group3"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Commands{}
			events, err := c.removeUserFromGroups(context.Background(), userID, tt.groupIDs, orgID)
			require.NoError(t, err)
			require.Len(t, events, len(tt.groupIDs))

			for i, ev := range events {
				removed, ok := ev.(*group.GroupUserRemovedEvent)
				require.Truef(t, ok, "event %d is %T, want *group.GroupUserRemovedEvent", i, ev)
				require.Equal(t, group.GroupUserRemovedEventType, removed.EventType)
				require.Equal(t, tt.groupIDs[i], removed.Aggregate().ID)
				require.Equal(t, orgID, removed.Aggregate().ResourceOwner)
				require.Equal(t, userID, removed.UserID)
			}
		})
	}
}

func newGroupUserAddedEvent(groupID, orgID, userID string) *group.GroupUserAddedEvent {
	return group.NewGroupUserAddedEvent(context.Background(),
		&group.NewAggregate(groupID, orgID).Aggregate,
		userID,
	)
}

func newGroupUserRemovedEvent(groupID, orgID, userID string) *group.GroupUserRemovedEvent {
	return group.NewGroupUserRemovedEvent(context.Background(),
		&group.NewAggregate(groupID, orgID).Aggregate,
		userID,
	)
}

func addNewUserEvent(userID, orgID string) *user.HumanAddedEvent {
	return user.NewHumanAddedEvent(context.Background(),
		&user.NewAggregate(userID, orgID).Aggregate,
		"username",
		"firstname",
		"lastname",
		"nickname",
		"displayname",
		AllowedLanguage,
		domain.GenderUnspecified,
		"email@test.ch",
		true,
	)
}
