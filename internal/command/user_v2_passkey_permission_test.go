package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// TestCommands_addUserPasskeyCode_permission guards the cross-organization passkey enrollment
// code issuance. The v2 gRPC handlers pass an empty resourceOwner and the RPC's auth annotation
// carries no org_field, so the API interceptor only verifies the caller's permission in the
// request-header org. The authorization therefore has to happen here, against the target user's
// actual resource owner resolved from the write model - not against anything the caller supplied.
func TestCommands_addUserPasskeyCode_permission(t *testing.T) {
	alg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))

	// The victim lives in org2. The caller supplies no resource owner at all, exactly like
	// internal/api/grpc/user/v2/passkey.go does.
	victimAgg := &user.NewAggregate("user1", "org2").Aggregate
	victimAdded := eventFromEventPusher(user.NewHumanAddedEvent(context.Background(),
		victimAgg,
		"username",
		"firstname",
		"lastname",
		"nickname",
		"displayname",
		language.German,
		domain.GenderUnspecified,
		"email@test.ch",
		true,
	))

	const urlTmpl = "https://example.com/passkey/register?userID={{.UserID}}&orgID={{.OrgID}}&codeID={{.CodeID}}&code={{.Code}}"

	tests := []struct {
		name string
		call func(*Commands) error
	}{
		{
			name: "AddUserPasskeyCode",
			call: func(c *Commands) error {
				_, err := c.AddUserPasskeyCode(context.Background(), "user1", "", alg)
				return err
			},
		},
		{
			name: "AddUserPasskeyCodeURLTemplate",
			call: func(c *Commands) error {
				_, err := c.AddUserPasskeyCodeURLTemplate(context.Background(), "user1", "", alg, urlTmpl)
				return err
			},
		},
		{
			name: "AddUserPasskeyCodeReturn",
			call: func(c *Commands) error {
				_, err := c.AddUserPasskeyCodeReturn(context.Background(), "user1", "", alg)
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
				// No push is expected: an unauthorized caller must not get a code issued.
				// Neither ID generation nor code generation runs before the check, so both
				// mocks are declared without any expected call.
				eventstore:       expectEventstore(expectFilter(victimAdded))(t),
				idGenerator:      id_mock.NewIDGeneratorExpectIDs(t),
				newEncryptedCode: newEncryptedCodeNoCall(t),
				loginPaths:       expectLoginPathsNoCall(t),
			}

			err := tt.call(c)
			assert.True(t, zerrors.IsPermissionDenied(err), "want permission denied, got %v", err)
			assert.Equal(t, domain.PermissionUserPasskeyWrite, gotPermission)
			assert.Equal(t, "org2", gotOrgID, "must authorize against the target user's org, not the caller-supplied one")
			assert.Equal(t, "user1", gotResourceID)
		})
	}
}
