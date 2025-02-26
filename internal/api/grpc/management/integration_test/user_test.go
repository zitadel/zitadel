//go:build integration

package management_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/user"
)

// TestImport_and_Get reproduces https://github.com/zitadel/zitadel/issues/5808
// which led to consistency issues due the call timestamp not being
// updated after a bulk Trigger.
// This test Imports a user and directly tries to Get it, 100 times in a loop.
// When the bug still existed, some (between 1 to 7 out of 100)
// Get calls would return a Not Found error.
func TestImport_and_Get(t *testing.T) {
	const N = 10

	for i := 0; i < N; i++ {
		firstName := strconv.Itoa(i)
		t.Run(firstName, func(tt *testing.T) {
			// create unique names.
			lastName := strconv.FormatInt(time.Now().Unix(), 10)
			userName := strings.Join([]string{firstName, lastName}, "_")
			email := strings.Join([]string{userName, "example.com"}, "@")

			res, err := Client.ImportHumanUser(OrgCTX, &management.ImportHumanUserRequest{
				UserName: userName,
				Profile: &management.ImportHumanUserRequest_Profile{
					FirstName:         firstName,
					LastName:          lastName,
					PreferredLanguage: language.Japanese.String(),
					Gender:            user.Gender_GENDER_DIVERSE,
				},
				Email: &management.ImportHumanUserRequest_Email{
					Email:           email,
					IsEmailVerified: true,
				},
			})
			require.NoError(tt, err)

			_, err = Client.GetUserByID(OrgCTX, &management.GetUserByIDRequest{Id: res.GetUserId()})

			s, ok := status.FromError(err)
			if ok && s != nil && s.Code() == codes.NotFound {
				tt.Errorf("iteration %d: user with id %q not found", i, res.GetUserId())
			}
			require.NoError(tt, err) // catch and fail on any other error
		})
	}
}

func TestImport_UnparsablePreferredLanguage(t *testing.T) {
	random := integration.RandString(5)
	_, err := Client.ImportHumanUser(OrgCTX, &management.ImportHumanUserRequest{
		UserName: random,
		Profile: &management.ImportHumanUserRequest_Profile{
			FirstName:         random,
			LastName:          random,
			PreferredLanguage: "not valid",
			Gender:            user.Gender_GENDER_DIVERSE,
		},
		Email: &management.ImportHumanUserRequest_Email{
			Email:           random + "@example.com",
			IsEmailVerified: true,
		},
	})
	require.NoError(t, err)
}

func TestAdd_MachineUser(t *testing.T) {
	random := integration.RandString(5)
	res, err := Client.AddMachineUser(OrgCTX, &management.AddMachineUserRequest{
		UserName:        random,
		Name:            "testMachineName1",
		Description:     "testMachineDescription1",
		AccessTokenType: 0,
	})
	require.NoError(t, err)

	_, err = Client.GetUserByID(OrgCTX, &management.GetUserByIDRequest{Id: res.GetUserId()})
	require.NoError(t, err)
}

func TestAdd_MachineUserCustomID(t *testing.T) {
	id := integration.RandString(5)
	random := integration.RandString(5)

	res, err := Client.AddMachineUser(OrgCTX, &management.AddMachineUserRequest{
		UserId:          &id,
		UserName:        random,
		Name:            "testMachineName1",
		Description:     "testMachineDescription1",
		AccessTokenType: 0,
	})
	require.NoError(t, err)

	_, err = Client.GetUserByID(OrgCTX, &management.GetUserByIDRequest{Id: id})
	require.NoError(t, err)

	require.Equal(t, id, res.GetUserId())
}
