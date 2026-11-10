package command

import (
	"context"
	"errors"
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
			name: "some users already exist in the group, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter( // to get the group write model
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
						eventFromEventPusher(
							addNewGroupUsersAddedEvent("group1", "org1", []string{"user1"}),
						),
					),
					expectFilter( // to get the user write model to check if user1 exists
						eventFromEventPusher(
							addNewUserEvent("user2", "org1"),
						),
					),
					expectPush(
						group.NewGroupUsersAddedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							[]string{"user2"},
						),
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
						eventFromEventPusher(
							addNewGroupUsersAddedEvent("group1", "org1", []string{"user1", "user2"}),
						),
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
						group.NewGroupUsersAddedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							[]string{"user1"},
						),
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
					expectFilter( // to get the user write model for user1
						eventFromEventPusher(
							addNewUserEvent("user1", "org1"),
						),
					),
					expectFilter( // to get the user write model for user2
						eventFromEventPusher(
							addNewUserEvent("user2", "org1"),
						),
					),
					expectPush(
						group.NewGroupUsersAddedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							[]string{"user1", "user2"},
						),
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
					expectFilter( // to get the user write model for user1
						eventFromEventPusher(
							addNewUserEvent("user1", "org1"),
						),
					),
					expectFilter( // to get the user write model for user2
						eventFromEventPusher(
							addNewUserEvent("user2", "org1"),
						),
					),
					expectPush(
						group.NewGroupUsersAddedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							[]string{"user1", "user2"},
						),
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
						eventFromEventPusher(
							addNewGroupUsersAddedEvent("group1", "org1", []string{"user1", "user2"}),
						),
					),
					expectPush(
						group.NewGroupUsersRemovedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							[]string{"user1", "user2"},
						),
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
						eventFromEventPusher(
							addNewGroupUsersAddedEvent("group1", "org1", []string{"user1", "user2"}),
						),
					),
					expectPush(
						group.NewGroupUsersRemovedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							[]string{"user1", "user2"},
						),
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
						eventFromEventPusher(
							addNewGroupUsersAddedEvent("group1", "org1", []string{"user1", "user2"}),
						),
					),
					expectPush(
						group.NewGroupUsersRemovedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							[]string{"user2"},
						),
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
						eventFromEventPusher(
							addNewGroupUsersAddedEvent("group1", "org1", []string{"user1", "user2"}),
						),
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
						eventFromEventPusher(
							addNewGroupUsersAddedEvent("group1", "org1", []string{"user1", "user2"}),
						),
					),
					expectPushFailed(
						pushErr,
						group.NewGroupUsersRemovedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							[]string{"user1"},
						),
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

func addNewGroupUsersAddedEvent(groupID, orgID string, userIds []string) *group.GroupUsersAddedEvent {
	return group.NewGroupUsersAddedEvent(context.Background(),
		&group.NewAggregate(groupID, orgID).Aggregate,
		userIds,
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
