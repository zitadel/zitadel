//go:build integration

package saml_test

import (
	"testing"

	"github.com/crewjam/saml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
)

// TestServer_CreateResponse_LoginClientHeaderGrantsNoPermission ensures that the
// x-zitadel-login-client header of the (unauthenticated) SAML request does not
// grant the named user any permission on the SAML request.
func TestServer_CreateResponse_LoginClientHeaderGrantsNoPermission(t *testing.T) {
	idpMetadata, err := Instance.GetSAMLIDPMetadata()
	require.NoError(t, err)
	acsPost := idpMetadata.IDPSSODescriptors[0].SingleSignOnServices[1]
	_, _, sp := createSAMLApplication(CTX, t, idpMetadata, saml.HTTPPostBinding, false, false)

	caller := Instance.Users.Get(integration.UserTypeNoPermission)
	callerCtx := Instance.WithAuthorization(CTX, integration.UserTypeNoPermission)
	victimID := Instance.CreateHumanUser(CTX).GetUserId()

	t.Run("create response", func(t *testing.T) {
		// a session with only a user check, as created when the login name is entered in the login UI
		req := createSessionAndSmlRequestForCallback(LoginCTX, t, sp, caller.ID, acsPost, victimID, saml.HTTPPostBinding)

		got, err := Client.CreateResponse(callerCtx, req)
		if err == nil {
			t.Fatalf("SAML request completed by a caller without session.link (response for victim %q, post: %t)", victimID, got.GetPost().GetSamlResponse() != "")
		}
		assertDenied(t, err)
	})

	t.Run("login client with permission", func(t *testing.T) {
		req := createSessionAndSmlRequestForCallback(LoginCTX, t, sp, Instance.Users.Get(integration.UserTypeLogin).ID, acsPost, victimID, saml.HTTPPostBinding)

		_, err := Client.CreateResponse(LoginCTX, req)
		require.NoError(t, err)
	})
}

// assertDenied checks the call was rejected by the permission check.
// A caller without any membership gets NotFound ("membership not found"), one with memberships PermissionDenied.
func assertDenied(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	st := status.Convert(err)
	switch st.Code() {
	case codes.PermissionDenied:
	case codes.NotFound:
		// a caller without any membership is rejected by the permission check itself;
		// assert the message so the test cannot pass on a genuinely missing request.
		assert.Contains(t, st.Message(), "membership not found", err)
	default:
		t.Errorf("expected PermissionDenied or NotFound, got %s: %v", st.Code(), err)
	}
}
