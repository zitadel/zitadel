//go:build integration

package integration_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/scim"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestDeleteUser_errors(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		orgID       string
		errorStatus int
	}{
		{
			name:        "not authenticated",
			ctx:         context.Background(),
			errorStatus: http.StatusUnauthorized,
		},
		{
			name:        "no permissions",
			ctx:         Instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			errorStatus: http.StatusNotFound,
		},
		{
			name:        "unknown user id",
			errorStatus: http.StatusNotFound,
		},
		{
			name:        "another org",
			orgID:       SecondaryOrganization.OrganizationId,
			errorStatus: http.StatusNotFound,
		},
		{
			name:        "another org with permissions",
			orgID:       SecondaryOrganization.OrganizationId,
			ctx:         Instance.WithAuthorization(CTX, integration.UserTypeIAMOwner),
			errorStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if ctx == nil {
				ctx = CTX
			}

			orgID := tt.orgID
			if orgID == "" {
				orgID = Instance.DefaultOrg.Id
			}
			err := Instance.Client.SCIM.Users.Delete(ctx, orgID, "1")

			statusCode := tt.errorStatus
			if statusCode == 0 {
				statusCode = http.StatusBadRequest
			}

			scim.RequireScimError(t, statusCode, err)
		})
	}
}

func TestDeleteUser_ensureReallyDeleted(t *testing.T) {
	// create user and dependencies
	createUserResp := Instance.CreateHumanUser(CTX)
	proj := Instance.CreateProject(CTX, t, Instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

	Instance.CreateProjectUserGrant(t, CTX, Instance.DefaultOrg.GetId(), proj.Id, createUserResp.UserId)

	// delete user via scim
	err := Instance.Client.SCIM.Users.Delete(CTX, Instance.DefaultOrg.Id, createUserResp.UserId)
	assert.NoError(t, err)

	// ensure it is really deleted => try to delete again => should 404
	err = Instance.Client.SCIM.Users.Delete(CTX, Instance.DefaultOrg.Id, createUserResp.UserId)
	scim.RequireScimError(t, http.StatusNotFound, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		// try to get user via api => should 404
		_, err = Instance.Client.UserV2.GetUserByID(CTX, &user.GetUserByIDRequest{UserId: createUserResp.UserId})
		integration.AssertGrpcStatus(tt, codes.NotFound, err)
	}, retryDuration, tick)
}
