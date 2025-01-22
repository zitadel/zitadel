//go:build integration

package integration_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/scim"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
	"google.golang.org/grpc/codes"
	"net/http"
	"testing"
	"time"
)

func TestDeleteUser_errors(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if ctx == nil {
				ctx = CTX
			}

			err := Instance.Client.SCIM.Users.Delete(ctx, Instance.DefaultOrg.Id, "1")

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
	proj, err := Instance.CreateProject(CTX)
	require.NoError(t, err)

	Instance.CreateProjectUserGrant(t, CTX, proj.Id, createUserResp.UserId)

	// delete user via scim
	err = Instance.Client.SCIM.Users.Delete(CTX, Instance.DefaultOrg.Id, createUserResp.UserId)
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

func TestDeleteUser_anotherOrg(t *testing.T) {
	createUserResp := Instance.CreateHumanUser(CTX)
	org := Instance.CreateOrganization(Instance.WithAuthorization(CTX, integration.UserTypeIAMOwner), gofakeit.Name(), gofakeit.Email())
	err := Instance.Client.SCIM.Users.Delete(CTX, org.OrganizationId, createUserResp.UserId)
	scim.RequireScimError(t, http.StatusNotFound, err)
}
