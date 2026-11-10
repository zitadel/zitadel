//go:build integration

package management_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/management"
)

// TestPasswordlessRegistration_CrossOrg simulates issuing a passwordless enrollment link for a
// user of a foreign organization through the deprecated v1 management API.
//
// The caller passes only its own organization, which the interceptor is happy with because that
// is where the caller holds the permission. The write lands on the target's aggregate regardless,
// so the denial has to come from the command, authorizing against the target's real owner.
func TestPasswordlessRegistration_CrossOrg(t *testing.T) {
	orgB := Instance.CreateOrganization(IAMOwnerCTX, integration.OrganizationName(), integration.Email())
	victimID := Instance.CreateHumanUserVerified(
		IAMOwnerCTX, orgB.GetOrganizationId(), integration.Email(), integration.Phone(),
	).GetUserId()

	// Either code is acceptable: the default permission backend reports the caller's missing
	// membership in the target organization as not found, PermissionCheckV2 as permission denied.
	deniedCodes := []codes.Code{codes.NotFound, codes.PermissionDenied}

	t.Run("org owner of another org is denied", func(t *testing.T) {
		// ReturnCode equivalent: on success the link with the plaintext code is in the response,
		// so asserting it stays empty proves nothing leaked.
		got, err := Client.AddPasswordlessRegistration(OrgCTX, &management.AddPasswordlessRegistrationRequest{
			UserId: victimID,
		})
		require.Error(t, err)
		assert.Contains(t, deniedCodes, status.Code(err))
		assert.Empty(t, got.GetLink())

		_, err = Client.SendPasswordlessRegistration(OrgCTX, &management.SendPasswordlessRegistrationRequest{
			UserId: victimID,
		})
		require.Error(t, err)
		assert.Contains(t, deniedCodes, status.Code(err))
	})

	t.Run("org owner of the same org is allowed", func(t *testing.T) {
		ownUserID := Instance.CreateHumanUserVerified(
			IAMOwnerCTX, Instance.DefaultOrg.GetId(), integration.Email(), integration.Phone(),
		).GetUserId()

		got, err := Client.AddPasswordlessRegistration(OrgCTX, &management.AddPasswordlessRegistrationRequest{
			UserId: ownUserID,
		})
		require.NoError(t, err)
		assert.NotEmpty(t, got.GetLink())

		_, err = Client.SendPasswordlessRegistration(OrgCTX, &management.SendPasswordlessRegistrationRequest{
			UserId: ownUserID,
		})
		require.NoError(t, err)
	})

	// SendPasswordlessRegistration is annotated with user.write, which IAM_USER_MANAGER holds
	// instance wide. Gating the command on user.credential.write instead - the annotation of the
	// sibling AddPasswordlessRegistration - would lock this role out of the endpoint.
	t.Run("instance wide user manager may still send", func(t *testing.T) {
		_, pat, err := Instance.CreateMachineUserPATWithMembership(IAMOwnerCTX, "IAM_USER_MANAGER")
		require.NoError(t, err)
		managerCTX := integration.WithAuthorizationToken(CTX, pat)

		_, err = Client.SendPasswordlessRegistration(managerCTX, &management.SendPasswordlessRegistrationRequest{
			UserId: victimID,
		})
		require.NoError(t, err)
	})
}
