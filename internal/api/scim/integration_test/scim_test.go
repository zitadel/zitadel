//go:build integration

package integration_test

import (
	"context"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

var (
	Instance              *integration.Instance
	SecondaryOrganization *org.AddOrganizationResponse
	CTX                   context.Context

	// remove comments in the json, as the default golang json unmarshaler cannot handle them
	// some test files (e.g. bulk, patch) are much easier to maintain with comments
	removeCommentsRegex = regexp.MustCompile("(?s)//.*?\n|/\\*.*?\\*/")
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)

		CTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeOrgOwner)

		iamOwnerCtx := Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
		SecondaryOrganization = Instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

		return m.Run()
	}())
}

func removeComments(json []byte) []byte {
	return removeCommentsRegex.ReplaceAll(json, nil)
}
