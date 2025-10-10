package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
		want    *AddUsersToGroupResponse
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
			name: "all users failed to be added (not found), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					),
					expectFilter(), // to get the user write model for user1
					expectFilter(), // to get the user write model for user2
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			want: &AddUsersToGroupResponse{
				FailedUserIDs: []string{"user1", "user2"},
			},
		},
		{
			name: "some users failed to be added (user in a different organization), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					),
					expectFilter( // to get the user write model for user1
						eventFromEventPusher(
							addNewUserEvent("user1", "org1"),
						),
					),
					expectFilter(), // to get the group user write model for user1
					expectFilter( // to get the user write model for user2
						eventFromEventPusher(
							addNewUserEvent("user2", "org2"),
						),
					),
					expectPush(
						group.NewGroupUserAddedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							"user1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			want: &AddUsersToGroupResponse{
				FailedUserIDs: []string{"user2"},
				ObjectDetails: &domain.ObjectDetails{
					ID:            "group1",
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "some users failed to be added (user in a different organization), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					),
					expectFilter( // to get the user write model for user1
						eventFromEventPusher(
							addNewUserEvent("user1", "org1"),
						),
					),
					expectFilter(), // to get the group user write model for user1
					expectFilter( // to get the user write model for user2
						eventFromEventPusher(
							addNewUserEvent("user2", "org2"),
						),
					),
					expectPush(
						group.NewGroupUserAddedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							"user1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			want: &AddUsersToGroupResponse{
				FailedUserIDs: []string{"user2"},
				ObjectDetails: &domain.ObjectDetails{
					ID:            "group1",
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "some users failed to be added (GroupUserWriteModel filter error), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					),
					expectFilter( // to get the user write model for user1
						eventFromEventPusher(
							addNewUserEvent("user1", "org1"),
						),
					),
					expectFilter(), // to get the group user write model for user1
					expectFilter( // to get the user write model for user2
						eventFromEventPusher(
							addNewUserEvent("user2", "org1"),
						),
					),
					expectFilterError(filterErr), // failed to get the group user write model for user2
					expectPush(
						group.NewGroupUserAddedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							"user1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			want: &AddUsersToGroupResponse{
				FailedUserIDs: []string{"user2"},
				ObjectDetails: &domain.ObjectDetails{
					ID:            "group1",
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "some users already exist in the group, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					),
					expectFilter( // to get the user write model for user1
						eventFromEventPusher(
							addNewUserEvent("user1", "org1"),
						),
					),
					expectFilter( // to get the group user write model for user1
						eventFromEventPusher(
							group.NewGroupUserAddedEvent(context.Background(),
								&group.NewAggregate("group1", "org1").Aggregate,
								"user1",
							),
						),
					),
					expectFilter( // to get the user write model for user2
						eventFromEventPusher(
							addNewUserEvent("user2", "org1"),
						),
					),
					expectFilter(), // to get the group user write model for user2
					expectPush(
						group.NewGroupUserAddedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							"user2",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			want: &AddUsersToGroupResponse{
				ObjectDetails: &domain.ObjectDetails{
					ID:            "group1",
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "failed to push events, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					),
					expectFilter(
						eventFromEventPusher(
							addNewUserEvent("user1", "org1"),
						),
					), // to get the user write model for user1
					expectFilter(), // to get the group user write model for user1
					expectPushFailed(
						pushErr,
						group.NewGroupUserAddedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							"user1",
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
					expectFilter(
						eventFromEventPusher(
							addNewGroupEvent("group1", "org1"),
						),
					),
					expectFilter( // to get the user write model for user1
						eventFromEventPusher(
							addNewUserEvent("user1", "org1"),
						),
					),
					expectFilter(), // to get the group user write model for user1
					expectFilter( // to get the user write model for user2
						eventFromEventPusher(
							addNewUserEvent("user2", "org1"),
						),
					),
					expectFilter(), // to get the group user write model for user2
					expectPush(
						group.NewGroupUserAddedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							"user1",
						),
						group.NewGroupUserAddedEvent(context.Background(),
							&group.NewAggregate("group1", "org1").Aggregate,
							"user2",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				groupID: "group1",
				userIDs: []string{"user1", "user2"},
			},
			want: &AddUsersToGroupResponse{
				ObjectDetails: &domain.ObjectDetails{
					ID:            "group1",
					ResourceOwner: "org1",
				},
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
			assert.Equal(t, tt.want, got)
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
