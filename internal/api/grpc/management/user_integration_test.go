//go:build integration

package management_test

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/user"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	CTX    context.Context
	Tester *integration.Tester
	Client management.ManagementServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(3 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		CTX, _ = Tester.WithSystemAuthorization(ctx, integration.OrgOwner), errCtx
		Client = Tester.Client.Mgmt
		return m.Run()
	}())
}

func TestImport_and_Get(t *testing.T) {
	const N = 100
	var misses int

	for i := 0; i < N; i++ {
		firstName := strconv.Itoa(i)
		t.Run(firstName, func(t *testing.T) {
			// create unique names.
			lastName := strconv.FormatInt(time.Now().Unix(), 10)
			userName := strings.Join([]string{firstName, lastName}, "_")
			email := strings.Join([]string{userName, "zitadel.com"}, "@")

			res, err := Client.ImportHumanUser(CTX, &management.ImportHumanUserRequest{
				UserName: userName,
				Profile: &management.ImportHumanUserRequest_Profile{
					FirstName:         firstName,
					LastName:          lastName,
					PreferredLanguage: language.Afrikaans.String(),
					Gender:            user.Gender_GENDER_DIVERSE,
				},
				Email: &management.ImportHumanUserRequest_Email{
					Email:           email,
					IsEmailVerified: true,
				},
			})
			require.NoError(t, err)

			_, err = Client.GetUserByID(CTX, &management.GetUserByIDRequest{Id: res.GetUserId()})

			if s, ok := status.FromError(err); ok {
				if s == nil {
					return
				}
				if s.Code() == codes.NotFound {
					t.Log(s)
					misses++
					return
				}
			}
			require.NoError(t, err) // catch and fail on any other error
		})
	}
	assert.Zerof(t, misses, "Not Found errors %d out of %d", misses, N)
}
