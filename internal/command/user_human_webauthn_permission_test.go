package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// victimHumanAddedEvent is a human in org2, targeted by a caller acting in org1.
func victimHumanAddedEvent() *user.HumanAddedEvent {
	return user.NewHumanAddedEvent(context.Background(),
		&user.NewAggregate("user1", "org2").Aggregate,
		"username",
		"firstname",
		"lastname",
		"nickname",
		"displayname",
		language.German,
		domain.GenderUnspecified,
		"email@test.ch",
		true,
	)
}

// TestCommands_HumanPasswordlessInitCode_permission guards cross-organization issuance of v1
// passwordless init codes. The deprecated management RPCs pass the *caller's* organization as
// resource owner, which only scopes the read: the init code write model then matches no events
// for a user of another organization and keeps the caller's organization, while the pushed event
// is re-owned to the target's organization by the eventstore.
//
// Authorizing against that write model would therefore authorize against the attacker's own
// organization. The assertion on the organization handed to checkPermission is what this test
// exists for - asserting only that a denial happened would pass even against the wrong org.
func TestCommands_HumanPasswordlessInitCode_permission(t *testing.T) {
	tests := []struct {
		name string
		// wantPermission mirrors the auth annotation of the RPC calling the command; a change
		// here means roles gained or lost access.
		wantPermission string
		call           func(*Commands) error
	}{
		{
			name:           "HumanAddPasswordlessInitCode",
			wantPermission: domain.PermissionUserCredentialWrite,
			call: func(c *Commands) error {
				_, err := c.HumanAddPasswordlessInitCode(context.Background(), "user1", "org1", nil)
				return err
			},
		},
		{
			name:           "HumanSendPasswordlessInitCode",
			wantPermission: domain.PermissionUserWrite,
			call: func(c *Commands) error {
				_, err := c.HumanSendPasswordlessInitCode(context.Background(), "user1", "org1", nil)
				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPermission, gotOrgID, gotResourceID string
			c := &Commands{
				checkPermission: func(_ context.Context, permission, orgID, resourceID string) error {
					gotPermission, gotOrgID, gotResourceID = permission, orgID, resourceID
					return zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied")
				},
				// Only the resource owner resolution is expected: a denied caller gets no code
				// generated and no event pushed.
				eventstore:  expectEventstore(expectFilter(eventFromEventPusher(victimHumanAddedEvent())))(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t),
			}

			err := tt.call(c)
			assert.True(t, zerrors.IsPermissionDenied(err), "want permission denied, got %v", err)
			assert.Equal(t, tt.wantPermission, gotPermission)
			assert.Equal(t, "org2", gotOrgID, "must authorize against the target user's org, not the caller-supplied one")
			assert.Equal(t, "user1", gotResourceID)
		})
	}

	// internal/api/grpc/auth/passwordless.go passes ctxData.UserID, so self management has to
	// short-circuit before any permission lookup. Filtering is failed on purpose to end the call
	// right after the check, which spares this test mocking the whole push.
	t.Run("self management needs no permission", func(t *testing.T) {
		ctx := authz.SetCtxData(context.Background(), authz.CtxData{UserID: "user1"})
		filterErr := zerrors.ThrowInternal(nil, "TEST-Ei0oh", "filter failed")
		checked := false
		c := &Commands{
			checkPermission: func(context.Context, string, string, string) error {
				checked = true
				return zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied")
			},
			eventstore:  expectEventstore(expectFilterError(filterErr))(t),
			idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "code1"),
		}

		_, err := c.HumanAddPasswordlessInitCode(ctx, "user1", "org1", nil)
		assert.False(t, checked, "self management must not require a permission")
		assert.True(t, zerrors.IsInternal(err), "want the call to proceed past the permission check, got %v", err)
	})
}
