//go:build integration

package org_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

func TestServer_DeactivateOrganization_RejectsCallerToken(t *testing.T) {
	orgResp := Instance.CreateOrganization(CTX, integration.OrganizationName(), integration.Email())
	machineUser := Instance.CreateUserTypeMachine(CTX, orgResp.GetOrganizationId())
	pat := Instance.CreatePersonalAccessToken(CTX, machineUser.GetId())
	userCtx := integration.WithAuthorizationToken(CTX, pat.GetToken())

	_, err := Instance.Client.Auth.GetMyUser(userCtx, &auth.GetMyUserRequest{})
	require.NoError(t, err)

	_, err = Client.DeactivateOrganization(CTX, &org.DeactivateOrganizationRequest{
		OrganizationId: orgResp.GetOrganizationId(),
	})
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		_, getErr := Instance.Client.Auth.GetMyUser(userCtx, &auth.GetMyUserRequest{})
		assert.Error(collect, getErr)
	}, retryDuration, tick, "timeout waiting for deactivated org to reject caller token")

	// Instance admin can still target the deactivated org.
	_, err = Client.ListOrganizations(CTX, &org.ListOrganizationsRequest{
		Queries: []*org.SearchQuery{
			{
				Query: &org.SearchQuery_IdQuery{
					IdQuery: &org.OrganizationIDQuery{
						Id: orgResp.GetOrganizationId(),
					},
				},
			},
		},
	})
	require.NoError(t, err)
}
