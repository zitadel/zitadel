//go:build integration

package events_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-resty/resty/v2"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/headzoo/surf.v1"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user"
	v2_user "github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_TestHumanUserReduces(t *testing.T) {
	instanceID := Instance.ID()
	orgID := Instance.DefaultOrg.Id

	userRepo := repository.UserRepository()

	t.Run("test human user added reduced", func(t *testing.T) {
		humanUserRequest := &management.AddHumanUserRequest{
			UserName: gofakeit.Username(),
			Profile: &management.AddHumanUserRequest_Profile{
				FirstName:         "first",
				LastName:          "last",
				NickName:          "nick",
				DisplayName:       "display",
				PreferredLanguage: "en",
				Gender:            user.Gender_GENDER_MALE,
			},
			Email: &management.AddHumanUserRequest_Email{
				Email:           gofakeit.Email(),
				IsEmailVerified: true,
			},
			Phone: &management.AddHumanUserRequest_Phone{
				Phone:           "+" + gofakeit.Phone(),
				IsPhoneVerified: true,
			},
			InitialPassword: "Password1!",
		}

		before := time.Now()
		resp, err := MgmtClient.AddHumanUser(IAMCTX, humanUserRequest)
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			// event user.added
			// event user.human.added
			// domain.User
			assert.Equal(t, instanceID, user.InstanceID)
			assert.Equal(t, orgID, user.OrganizationID)
			assert.Equal(t, resp.UserId, user.ID)
			assert.Equal(t, humanUserRequest.UserName, user.Username)
			assert.Equal(t, domain.UserStateActive, user.State)
			assert.WithinRange(t, user.UpdatedAt, before, after)
			assert.WithinRange(t, user.CreatedAt, before, after)
			// Email
			assert.Equal(t, humanUserRequest.Email.Email, user.Human.Email.Address)
			assert.NotZero(t, user.Human.Email.VerifiedAt)
			assert.Nil(t, user.Human.Email.PendingVerification)
			// Phone
			assert.Equal(t, humanUserRequest.Phone.Phone, user.Human.Phone.Number)
			assert.NotZero(t, user.Human.Phone.VerifiedAt)
			assert.Nil(t, user.Human.Phone.PendingVerification)
			// Human
			assert.Equal(t, humanUserRequest.Profile.FirstName, user.Human.FirstName)
			assert.Equal(t, humanUserRequest.Profile.LastName, user.Human.LastName)
			assert.Equal(t, humanUserRequest.Profile.NickName, user.Human.Nickname)
			assert.Equal(t, humanUserRequest.Profile.DisplayName, user.Human.DisplayName)
			assert.Equal(t, humanUserRequest.Profile.PreferredLanguage, user.Human.PreferredLanguage)
			assert.Equal(t, uint8(humanUserRequest.Profile.Gender), user.Human.Gender)
		}, retryDuration, tick)
	})

	t.Run("test human user register reduced", func(t *testing.T) {
		token := integration.SystemToken
		client := resty.New()

		bow := surf.NewBrowser()
		err := bow.Open("http://localhost:8080" + "/ui/login/register/org")
		require.NoError(t, err)
		require.Equal(t, 200, bow.StatusCode())

		csfr, err := bow.Forms()[1].Value("gorilla.csrf.Token")
		require.NoError(t, err)

		before := time.Now()
		client.SetCookieJar(bow.CookieJar())
		firstName := gofakeit.Name()
		lastName := gofakeit.Name()
		email := gofakeit.Email()
		out, err := client.R().SetAuthToken(token).
			SetFormData(map[string]string{
				"gorilla.csrf.Token": csfr,
				"orgname":            gofakeit.Name(),
				"firstname":          firstName,
				"lastname":           lastName,
				// "email":                          "@zitadel.localhost",
				"email":                          email,
				"register-password":              "Password1!",
				"register-password-confirmation": "Password1!",
			}).
			Post("http://localhost:8080" + "/ui/login/register/org")

		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())
		after := time.Now()

		ctx := t.Context()
		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		assert.NoError(t, err)
		instanceID := instance.ID

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.Human().FirstNameCondition(database.TextOperationEqual, firstName),
					userRepo.Human().LastNameCondition(database.TextOperationEqual, lastName),
				)),
			)
			require.NoError(t, err)

			// domain.User
			assert.Equal(t, instanceID, user.InstanceID)
			assert.Equal(t, email, user.Username)
			assert.Equal(t, domain.UserStateActive, user.State)
			assert.WithinRange(t, user.CreatedAt, before, after)
			// Email
			assert.Equal(t, email, user.Human.Email.Address)
			assert.NotZero(t, user.Human.Email.VerifiedAt)
			assert.Nil(t, user.Human.Email.PendingVerification)
			// Phone
			assert.Equal(t, email, user.Human.Phone.Number)
			assert.NotZero(t, user.Human.Phone.VerifiedAt)
			assert.Nil(t, user.Human.Phone.PendingVerification)
			// Human
			assert.Equal(t, firstName, user.Human.FirstName)
			assert.Equal(t, lastName, user.Human.LastName)
		}, retryDuration, tick)
	})

	// TODO
	// t.Run("test human user added init reduced", func(t *testing.T) {
	// 	humanUserRequest := &management.AddHumanUserRequest{
	// UserName: gofakeit.Username(),
	// 		Profile: &management.AddHumanUserRequest_Profile{
	// 			FirstName:         "first",
	// 			LastName:          "last",
	// 			NickName:          "nick",
	// 			DisplayName:       "display",
	// 			PreferredLanguage: "en",
	// 			Gender:            user.Gender_GENDER_MALE,
	// 		},
	// 		Email: &management.AddHumanUserRequest_Email{
	// 			Email:           gofakeit.Email(),
	// 			IsEmailVerified: true,
	// 		},
	// 		Phone: &management.AddHumanUserRequest_Phone{
	// 			Phone:           "+" + gofakeit.Phone(),
	// 			IsPhoneVerified: true,
	// 		},
	// 		InitialPassword: "Password1!",
	// 	}

	// 	resp, err := MgmtClient.AddHumanUser(IAMCTX, humanUserRequest)
	// 	require.NoError(t, err)

	// 	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
	// 	assert.EventuallyWithT(t, func(t *assert.CollectT) {
	// 		user, err := userRepo.Get(
	// 			CTX,
	// 			pool,
	// 			database.WithCondition(database.And(
	// 				userRepo.InstanceIDCondition(instanceID),
	// 				userRepo.OrganizationIDCondition(orgID),
	// 				userRepo.IDCondition(resp.UserId),
	// 			)),
	// 		)
	// 		require.NoError(t, err)

	// 		// event user.human.added.initialization.code.added
	// 	}, retryDuration, tick)
	// })

	// TODO
	// t.Run("test human user added init check succeeded reduced", func(t *testing.T) {
	// 	humanUserRequest := &management.AddHumanUserRequest{
	// UserName: gofakeit.Username(),
	// 		Profile: &management.AddHumanUserRequest_Profile{
	// 			FirstName:         "first",
	// 			LastName:          "last",
	// 			NickName:          "nick",
	// 			DisplayName:       "display",
	// 			PreferredLanguage: "en",
	// 			Gender:            user.Gender_GENDER_MALE,
	// 		},
	// 		Email: &management.AddHumanUserRequest_Email{
	// 			Email:           gofakeit.Email(),
	// 			IsEmailVerified: true,
	// 		},
	// 		Phone: &management.AddHumanUserRequest_Phone{
	// 			Phone:           "+" + gofakeit.Phone(),
	// 			IsPhoneVerified: true,
	// 		},
	// 		InitialPassword: "Password1!",
	// 	}

	// 	resp, err := MgmtClient.AddHumanUser(IAMCTX, humanUserRequest)
	// 	require.NoError(t, err)

	// 	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
	// 	assert.EventuallyWithT(t, func(t *assert.CollectT) {
	// 		user, err := userRepo.Get(
	// 			CTX,
	// 			pool,
	// 			database.WithCondition(database.And(
	// 				userRepo.InstanceIDCondition(instanceID),
	// 				userRepo.OrganizationIDCondition(orgID),
	// 				userRepo.IDCondition(resp.UserId),
	// 			)),
	// 		)
	// 		require.NoError(t, err)

	// 		// event user.human.added.initialization.check.succeeded
	// 	}, retryDuration, tick)
	// })

	t.Run("test human user locked reduced", func(t *testing.T) {
		humanUserRequest := &management.AddHumanUserRequest{
			UserName: gofakeit.Username(),
			Profile: &management.AddHumanUserRequest_Profile{
				FirstName:         "first",
				LastName:          "last",
				NickName:          "nick",
				DisplayName:       "display",
				PreferredLanguage: "en",
				Gender:            user.Gender_GENDER_MALE,
			},
			Email: &management.AddHumanUserRequest_Email{
				Email:           gofakeit.Email(),
				IsEmailVerified: true,
			},
			Phone: &management.AddHumanUserRequest_Phone{
				Phone:           "+" + gofakeit.Phone(),
				IsPhoneVerified: true,
			},
			InitialPassword: "Password1!",
		}

		resp, err := MgmtClient.AddHumanUser(IAMCTX, humanUserRequest)
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			assert.Equal(t, domain.UserStateActive, user.State)
		}, retryDuration, tick)

		before := time.Now()
		_, err = MgmtClient.LockUser(IAMCTX, &management.LockUserRequest{
			Id: resp.UserId,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			// event user.locked
			assert.Equal(t, domain.UserStateLocked, user.State)
			assert.WithinRange(t, user.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test human user unlock reduced", func(t *testing.T) {
		humanUserRequest := &management.AddHumanUserRequest{
			UserName: gofakeit.Username(),
			Profile: &management.AddHumanUserRequest_Profile{
				FirstName:         "first",
				LastName:          "last",
				NickName:          "nick",
				DisplayName:       "display",
				PreferredLanguage: "en",
				Gender:            user.Gender_GENDER_MALE,
			},
			Email: &management.AddHumanUserRequest_Email{
				Email:           gofakeit.Email(),
				IsEmailVerified: true,
			},
			Phone: &management.AddHumanUserRequest_Phone{
				Phone:           "+" + gofakeit.Phone(),
				IsPhoneVerified: true,
			},
			InitialPassword: "Password1!",
		}

		resp, err := MgmtClient.AddHumanUser(IAMCTX, humanUserRequest)
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			assert.Equal(t, domain.UserStateActive, user.State)
		}, retryDuration, tick)

		_, err = MgmtClient.LockUser(IAMCTX, &management.LockUserRequest{
			Id: resp.UserId,
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			assert.Equal(t, domain.UserStateLocked, user.State)
		}, retryDuration, tick)

		before := time.Now()
		_, err = MgmtClient.UnlockUser(IAMCTX, &management.UnlockUserRequest{
			Id: resp.UserId,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			// event user.unlocked
			assert.Equal(t, domain.UserStateActive, user.State)
			assert.WithinRange(t, user.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test human user deactivate reduced", func(t *testing.T) {
		humanUserRequest := &management.AddHumanUserRequest{
			UserName: gofakeit.Username(),
			Profile: &management.AddHumanUserRequest_Profile{
				FirstName:         "first",
				LastName:          "last",
				NickName:          "nick",
				DisplayName:       "display",
				PreferredLanguage: "en",
				Gender:            user.Gender_GENDER_MALE,
			},
			Email: &management.AddHumanUserRequest_Email{
				Email:           gofakeit.Email(),
				IsEmailVerified: true,
			},
			Phone: &management.AddHumanUserRequest_Phone{
				Phone:           "+" + gofakeit.Phone(),
				IsPhoneVerified: true,
			},
			InitialPassword: "Password1!",
		}

		resp, err := MgmtClient.AddHumanUser(IAMCTX, humanUserRequest)
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			assert.Equal(t, domain.UserStateActive, user.State)
		}, retryDuration, tick)

		before := time.Now()
		_, err = MgmtClient.DeactivateUser(IAMCTX, &management.DeactivateUserRequest{
			Id: resp.UserId,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			// event user.deactivated
			assert.Equal(t, domain.UserStateInactive, user.State)
			assert.WithinRange(t, user.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test human user reeactivate reduced", func(t *testing.T) {
		humanUserRequest := &management.AddHumanUserRequest{
			UserName: gofakeit.Username(),
			Profile: &management.AddHumanUserRequest_Profile{
				FirstName:         "first",
				LastName:          "last",
				NickName:          "nick",
				DisplayName:       "display",
				PreferredLanguage: "en",
				Gender:            user.Gender_GENDER_MALE,
			},
			Email: &management.AddHumanUserRequest_Email{
				Email:           gofakeit.Email(),
				IsEmailVerified: true,
			},
			Phone: &management.AddHumanUserRequest_Phone{
				Phone:           "+" + gofakeit.Phone(),
				IsPhoneVerified: true,
			},
			InitialPassword: "Password1!",
		}

		resp, err := MgmtClient.AddHumanUser(IAMCTX, humanUserRequest)
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			assert.Equal(t, domain.UserStateActive, user.State)
		}, retryDuration, tick)

		_, err = MgmtClient.DeactivateUser(IAMCTX, &management.DeactivateUserRequest{
			Id: resp.UserId,
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			assert.Equal(t, domain.UserStateInactive, user.State)
		}, retryDuration, tick)

		before := time.Now()
		_, err = MgmtClient.ReactivateUser(IAMCTX, &management.ReactivateUserRequest{
			Id: resp.UserId,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			// event user.reactivated
			assert.Equal(t, domain.UserStateActive, user.State)
			assert.WithinRange(t, user.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test human user removed reduced", func(t *testing.T) {
		humanUserRequest := &management.AddHumanUserRequest{
			UserName: gofakeit.Username(),
			Profile: &management.AddHumanUserRequest_Profile{
				FirstName:         "first",
				LastName:          "last",
				NickName:          "nick",
				DisplayName:       "display",
				PreferredLanguage: "en",
				Gender:            user.Gender_GENDER_MALE,
			},
			Email: &management.AddHumanUserRequest_Email{
				Email:           gofakeit.Email(),
				IsEmailVerified: true,
			},
			Phone: &management.AddHumanUserRequest_Phone{
				Phone:           "+" + gofakeit.Phone(),
				IsPhoneVerified: true,
			},
			InitialPassword: "Password1!",
		}

		resp, err := MgmtClient.AddHumanUser(IAMCTX, humanUserRequest)
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			assert.NotNil(t, user)
		}, retryDuration, tick)

		_, err = UserClient.DeleteUser(IAMCTX, &v2_user.DeleteUserRequest{
			UserId: resp.UserId,
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)

			// event user.removed
			require.ErrorIs(t, &database.NoRowFoundError{}, err)
		}, retryDuration, tick)
	})

	t.Run("test human user change name reduced", func(t *testing.T) {
		humanUserRequest := &management.AddHumanUserRequest{
			UserName: gofakeit.Username(),
			Profile: &management.AddHumanUserRequest_Profile{
				FirstName:         "first",
				LastName:          "last",
				NickName:          "nick",
				DisplayName:       "display",
				PreferredLanguage: "en",
				Gender:            user.Gender_GENDER_MALE,
			},
			Email: &management.AddHumanUserRequest_Email{
				Email:           gofakeit.Email(),
				IsEmailVerified: true,
			},
			Phone: &management.AddHumanUserRequest_Phone{
				Phone:           "+" + gofakeit.Phone(),
				IsPhoneVerified: true,
			},
			InitialPassword: "Password1!",
		}

		resp, err := MgmtClient.AddHumanUser(IAMCTX, humanUserRequest)
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			assert.NotNil(t, user)
		}, retryDuration, tick)

		username := gofakeit.Username()
		_, err = UserClient.UpdateHumanUser(IAMCTX, &v2_user.UpdateHumanUserRequest{
			UserId:   resp.UserId,
			Username: &username,
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			// event user.username.changed
			assert.Equal(t, username, user.Username)
		}, retryDuration, tick)
	})

	// TODO
	// t.Run("test user domain claimed reduced", func(t *testing.T) {
	// 	org, err := OrgClient.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
	// 		Name: gofakeit.Name(),
	// 	})
	// 	require.NoError(t, err)
	// 	orgID := org.Id

	// 	// change custom domain policy to make validation mandatory
	// 	_, err = AdminClient.AddCustomDomainPolicy(IAMCTX, &admin.AddCustomDomainPolicyRequest{
	// 		OrgId:              orgID,
	// 		ValidateOrgDomains: true,
	// 	})
	// 	require.NoError(t, err)

	// 	// add user
	// 	userName := gofakeit.Name()
	// 	userID := gofakeit.Name()
	// 	UserClient.AddHumanUser(IAMCTX,
	// 		&v2beta_user.AddHumanUserRequest{
	// 			Username: &userName,
	// 			UserId:   &userID,
	// 			Organization: &object.Organization{
	// 				Org: &object.Organization_OrgId{
	// 					OrgId: orgID,
	// 				},
	// 			},
	// 			Profile: &v2beta_user.SetHumanProfile{
	// 				GivenName:         "Donald",
	// 				FamilyName:        "Duck",
	// 				NickName:          gu.Ptr("Dukkie"),
	// 				DisplayName:       gu.Ptr("Donald Duck"),
	// 				PreferredLanguage: gu.Ptr("en"),
	// 				Gender:            v2beta_user.Gender_GENDER_DIVERSE.Enum(),
	// 			},
	// 			Email: &v2beta_user.SetHumanEmail{
	// 				Email: gofakeit.Email(),
	// 			},
	// 			Phone: &v2beta_user.SetHumanPhone{},
	// 			Metadata: []*v2beta_user.SetMetadataEntry{
	// 				{
	// 					Key:   "somekey",
	// 					Value: []byte("somevalue"),
	// 				},
	// 			},
	// 			PasswordType: &v2beta_user.AddHumanUserRequest_Password{
	// 				Password: &v2beta_user.Password{
	// 					Password:       "DifficultPW666!",
	// 					ChangeRequired: true,
	// 				},
	// 			},
	// 		})
	// 	// humanUserRequest := &management.AddHumanUserRequest{
	// 	// 	UserName: userName,
	// 	// 	Profile: &management.AddHumanUserRequest_Profile{
	// 	// 		FirstName:         "first",
	// 	// 		LastName:          "last",
	// 	// 		NickName:          "nick",
	// 	// 		DisplayName:       "display",
	// 	// 		PreferredLanguage: "en",
	// 	// 		Gender:            user.Gender_GENDER_MALE,
	// 	// 	},
	// 	// 	Email: &management.AddHumanUserRequest_Email{
	// 	// 		Email:           gofakeit.Email(),
	// 	// 		IsEmailVerified: true,
	// 	// 	},
	// 	// 	Phone: &management.AddHumanUserRequest_Phone{
	// 	// 		Phone:           "+" + gofakeit.Phone(),
	// 	// 		IsPhoneVerified: true,
	// 	// 	},
	// 	// 	InitialPassword: "Password1!",
	// 	// }
	// 	// resp, err := MgmtClient.AddHumanUser(IAMCTX, humanUserRequest)
	// 	require.NoError(t, err)

	// 	// check user exists
	// 	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
	// 	assert.EventuallyWithT(t, func(t *assert.CollectT) {
	// 		user, err := userRepo.Get(
	// 			CTX,
	// 			pool,
	// 			database.WithCondition(database.And(
	// 				userRepo.InstanceIDCondition(instanceID),
	// 				userRepo.OrganizationIDCondition(orgID),
	// 				userRepo.IDCondition(userID),
	// 			)),
	// 		)
	// 		require.NoError(t, err)

	// 		assert.NotNil(t, user)
	// 	}, retryDuration, tick)

	// 	domain := "localhost:8383"

	// 	// add organization domain
	// 	_, err = OrgClient.AddOrganizationDomain(IAMCTX, &v2beta_org.AddOrganizationDomainRequest{
	// 		OrganizationId: orgID,
	// 		Domain:         domain,
	// 	})
	// 	require.NoError(t, err)

	// 	// request validation for domain
	// 	r, err := OrgClient.GenerateOrganizationDomainValidation(IAMCTX, &v2beta_org.GenerateOrganizationDomainValidationRequest{
	// 		OrganizationId: orgID,
	// 		Domain:         domain,
	// 		Type:           v2beta_org.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP,
	// 	})
	// 	require.NoError(t, err)

	// 	token := r.Token

	// 	// crate http server for domain validation
	// 	http.HandleFunc("/.well-known/zitadel-challenge/"+token+".txt", func(w http.ResponseWriter, req *http.Request) {
	// 		w.Header().Set("Content-Type", "text/plain")
	// 		w.Write([]byte(token))
	// 	})
	// 	// NOTE Apple introduced a policy change where TLS server certificates issued after July 1, 2019, must have a validity period of 825 days or fewer.
	// 	// to regenerate; openssl req -new -x509 -key server.key -out server.crt -days 800 -subj "/CN=localhost" -addext "subjectAltName = DNS:localhost,IP:127.0.0.1
	// 	cert, err := tls.X509KeyPair(serverCrt, serverKey)
	// 	require.NoError(t, err)
	// 	tlsConfig := &tls.Config{
	// 		Certificates: []tls.Certificate{cert},
	// 		MinVersion:   tls.VersionTLS12, // Recommended for security
	// 	}
	// 	server := &http.Server{
	// 		Addr:      ":8383",
	// 		Handler:   http.DefaultServeMux,
	// 		TLSConfig: tlsConfig,
	// 	}
	// 	go func() {
	// 		err := server.ListenAndServeTLS("", "")
	// 		require.NoError(t, err)
	// 	}()
	// 	defer server.Close()

	// 	// validate domain
	// 	_, err = OrgClient.VerifyOrganizationDomain(IAMCTX, &v2beta_org.VerifyOrganizationDomainRequest{
	// 		OrganizationId: orgID,
	// 		Domain:         domain,
	// 	})
	// 	require.NoError(t, err)

	// 	// check usernme now includes the domain
	// 	retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
	// 	assert.EventuallyWithT(t, func(t *assert.CollectT) {
	// 		user, err := userRepo.Get(
	// 			CTX,
	// 			pool,
	// 			database.WithCondition(database.And(
	// 				userRepo.InstanceIDCondition(instanceID),
	// 				userRepo.OrganizationIDCondition(orgID),
	// 				userRepo.IDCondition(userID),
	// 			)),
	// 		)
	// 		require.NoError(t, err)

	// 		// event user.domain.claimed
	// 		assert.Equal(t, userName+domain, user.Username)
	// 	}, retryDuration, tick)
	// })

	t.Run("test human user profile change reduced", func(t *testing.T) {
		humanUserRequest := &management.AddHumanUserRequest{
			UserName: gofakeit.Username(),
			Profile: &management.AddHumanUserRequest_Profile{
				FirstName:         "first",
				LastName:          "last",
				NickName:          "nick",
				DisplayName:       "display",
				PreferredLanguage: "en",
				Gender:            user.Gender_GENDER_MALE,
			},
			Email: &management.AddHumanUserRequest_Email{
				Email:           gofakeit.Email(),
				IsEmailVerified: true,
			},
			Phone: &management.AddHumanUserRequest_Phone{
				Phone:           "+" + gofakeit.Phone(),
				IsPhoneVerified: true,
			},
			InitialPassword: "Password1!",
		}

		resp, err := MgmtClient.AddHumanUser(IAMCTX, humanUserRequest)
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			assert.NotNil(t, user)
		}, retryDuration, tick)

		profile := &v2_user.SetHumanProfile{
			GivenName:         gofakeit.FirstName(),
			FamilyName:        gofakeit.LastName(),
			NickName:          gu.Ptr(gofakeit.Color()),
			DisplayName:       gu.Ptr(gofakeit.Name()),
			PreferredLanguage: gu.Ptr("fr"),
			Gender:            gu.Ptr(v2_user.Gender_GENDER_FEMALE),
		}
		_, err = UserClient.UpdateHumanUser(IAMCTX, &v2_user.UpdateHumanUserRequest{
			UserId:  resp.UserId,
			Profile: profile,
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			// event user.human.profile.changed
			assert.Equal(t, profile.GivenName, user.Human.FirstName)
			assert.Equal(t, profile.FamilyName, user.Human.LastName)
			assert.Equal(t, *profile.NickName, user.Human.Nickname)
			assert.Equal(t, *profile.DisplayName, user.Human.DisplayName)
			assert.Equal(t, *profile.PreferredLanguage, user.Human.PreferredLanguage)
			assert.Equal(t, uint8(*profile.Gender), user.Human.Gender)
		}, retryDuration, tick)
	})

	// TODO fix this: the issue is that a user called tester needs to be created for this test to work on the default instance and default organization
	// I added a delete user command (Over REST because there is no gRPC client (as far as I can tell) for the default instance)
	// the issue is if the user does not exist it will throw and error when trying to delete it
	t.Run("test human user avatar added reduced", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)
		orgRepo := repository.OrganizationRepository()
		org_, err := orgRepo.Get(ctx, pool, database.WithCondition(database.And(orgRepo.InstanceIDCondition(instanceID), orgRepo.NameCondition(database.TextOperationEqual, "ZITADEL"))))
		require.NoError(t, err)
		orgID := org_.ID

		// delete user tester
		deleteUserReq := v2_user.DeleteUserRequest{
			UserId: "tester",
		}
		deleteUserReqJSON, err := json.Marshal(&deleteUserReq)
		require.NoError(t, err)
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgID).
			SetBody(deleteUserReqJSON).
			Delete("http://localhost:8080" + "/v2beta/users/human")
		require.NoError(t, err)
		// require.Equal(t, 200, out.StatusCode())
		// require.NoError(t, err)

		// create user tester
		createUserReq := &v2_user.AddHumanUserRequest{
			UserId: gu.Ptr("tester"), // YOLO
			Organization: &object.Organization{
				Org: &object.Organization_OrgId{
					OrgId: orgID,
				},
			},
			Profile: &v2_user.SetHumanProfile{
				GivenName:         "Donald",
				FamilyName:        "Duck",
				NickName:          gu.Ptr("Dukkie"),
				DisplayName:       gu.Ptr("Donald Duck"),
				PreferredLanguage: gu.Ptr("en"),
				Gender:            v2_user.Gender_GENDER_DIVERSE.Enum(),
			},
			Email: &v2_user.SetHumanEmail{
				Email: gofakeit.Email(),
			},
			Phone: &v2_user.SetHumanPhone{},
			Metadata: []*v2_user.SetMetadataEntry{
				{
					Key:   "somekey",
					Value: []byte("somevalue"),
				},
			},
			PasswordType: &v2_user.AddHumanUserRequest_Password{
				Password: &v2_user.Password{
					Password:       "DifficultPW666!",
					ChangeRequired: true,
				},
			},
		}
		createUserReqJSON, err := json.Marshal(createUserReq)
		require.NoError(t, err)
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgID).
			SetBody(createUserReqJSON).
			Post("http://localhost:8080" + "/v2beta/users/human")
		require.NoError(t, err)
		// require.Equal(t, 201, out.StatusCode())

		require.NoError(t, err)

		// client = resty.New()
		// POST avatar
		before := time.Now()
		out, err = client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			SetHeader("x-zitadel-orgid", orgID).
			Post("http://localhost:8080" + "/assets/v1" + "/users/me/avatar")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition("tester"),
				)),
			)
			require.NoError(t, err)

			// event user.human.avatar.added
			avatarUrl := "users/tester/avatar?"
			assert.Contains(t, user.Human.AvatarKey, avatarUrl)
			assert.WithinRange(t, user.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	// TODO fix this: the issue is that a user called tester needs to be created for this test to work on the default instance and default organization
	// I added a delete user command (Over REST because there is no gRPC client (as far as I can tell) for the default instance)
	// the issue is if the user does not exist it will throw and error when trying to delete it
	t.Run("test human user avatar removed reduced", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)
		orgRepo := repository.OrganizationRepository()
		org_, err := orgRepo.Get(ctx, pool, database.WithCondition(database.And(orgRepo.InstanceIDCondition(instanceID), orgRepo.NameCondition(database.TextOperationEqual, "ZITADEL"))))
		require.NoError(t, err)
		orgID := org_.ID

		// delete user tester
		deleteUserReq := v2_user.DeleteUserRequest{
			UserId: "tester",
		}
		deleteUserReqJSON, err := json.Marshal(&deleteUserReq)
		require.NoError(t, err)
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgID).
			SetBody(deleteUserReqJSON).
			Delete("http://localhost:8080" + "/v2beta/users/human")
		require.NoError(t, err)
		// require.Equal(t, 200, out.StatusCode())
		// require.NoError(t, err)

		// create user tester
		createUserReq := &v2_user.AddHumanUserRequest{
			UserId: gu.Ptr("tester"), // YOLO
			Organization: &object.Organization{
				Org: &object.Organization_OrgId{
					OrgId: orgID,
				},
			},
			Profile: &v2_user.SetHumanProfile{
				GivenName:         "Donald",
				FamilyName:        "Duck",
				NickName:          gu.Ptr("Dukkie"),
				DisplayName:       gu.Ptr("Donald Duck"),
				PreferredLanguage: gu.Ptr("en"),
				Gender:            v2_user.Gender_GENDER_DIVERSE.Enum(),
			},
			Email: &v2_user.SetHumanEmail{
				Email: gofakeit.Email(),
			},
			Phone: &v2_user.SetHumanPhone{},
			Metadata: []*v2_user.SetMetadataEntry{
				{
					Key:   "somekey",
					Value: []byte("somevalue"),
				},
			},
			PasswordType: &v2_user.AddHumanUserRequest_Password{
				Password: &v2_user.Password{
					Password:       "DifficultPW666!",
					ChangeRequired: true,
				},
			},
		}
		createUserReqJSON, err := json.Marshal(createUserReq)
		require.NoError(t, err)
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgID).
			SetBody(createUserReqJSON).
			Post("http://localhost:8080" + "/v2beta/users/human")
		require.NoError(t, err)
		// require.Equal(t, 201, out.StatusCode())

		require.NoError(t, err)

		// POST avatar
		out, err = client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			SetHeader("x-zitadel-orgid", orgID).
			Post("http://localhost:8080" + "/assets/v1" + "/users/me/avatar")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition("tester"),
				)),
			)
			require.NoError(t, err)

			require.NotEmpty(t, user.Human.AvatarKey)

			avatarUrl := "users/tester/avatar?"
			assert.Equal(t, avatarUrl, string(user.Human.AvatarKey)[:len(avatarUrl)])
		}, retryDuration, tick)

		// delete avatar
		before := time.Now()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgID).
			Delete("http://localhost:8080" + "/auth/v1" + "/users/me/avatar")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())
		after := time.Now()

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition("tester"),
				)),
			)
			require.NoError(t, err)

			assert.Nil(t, user.Human.AvatarKey)
			assert.WithinRange(t, user.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	password := gofakeit.Password(true, true, true, true, true, 10)
	t.Run("test human user password change reduced", func(t *testing.T) {
		createUserReq := &v2_user.AddHumanUserRequest{
			// UserId: gu.Ptr("tester"),
			Organization: &object.Organization{
				Org: &object.Organization_OrgId{
					OrgId: orgID,
				},
			},
			Profile: &v2_user.SetHumanProfile{
				GivenName:         "Donald",
				FamilyName:        "Duck",
				NickName:          gu.Ptr("Dukkie"),
				DisplayName:       gu.Ptr("Donald Duck"),
				PreferredLanguage: gu.Ptr("en"),
				Gender:            v2_user.Gender_GENDER_DIVERSE.Enum(),
			},
			Email: &v2_user.SetHumanEmail{
				Email: gofakeit.Email(),
			},
			Phone: &v2_user.SetHumanPhone{},
			Metadata: []*v2_user.SetMetadataEntry{
				{
					Key:   "somekey",
					Value: []byte("somevalue"),
				},
			},
			PasswordType: &v2_user.AddHumanUserRequest_Password{
				Password: &v2_user.Password{
					Password:       password,
					ChangeRequired: true,
				},
			},
		}
		res, err := UserClient.AddHumanUser(IAMCTX, createUserReq)
		require.NoError(t, err)

		userID := res.UserId

		// make sure user exists
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(userID),
				)),
			)
			require.NoError(t, err)

			assert.NotNil(t, user)
		}, retryDuration, tick)

		// change user password
		before := time.Now()
		_, err = UserClient.UpdateHumanUser(IAMCTX, &v2_user.UpdateHumanUserRequest{
			UserId: userID,
			Password: &v2_user.SetPassword{
				PasswordType: &v2_user.SetPassword_Password{
					Password: &v2_user.Password{
						Password:       gofakeit.Password(true, true, true, true, true, 10),
						ChangeRequired: true,
					},
				},
				Verification: &v2_user.SetPassword_CurrentPassword{
					CurrentPassword: password,
				},
			},
		})
		require.NoError(t, err)
		after := time.Now()

		// check password updated
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			humanUser, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(userID),
				)),
			)
			require.NoError(t, err)

			assert.Equal(t, true, humanUser.Human.Password.IsChangeRequired)
			assert.WithinRange(t, humanUser.Human.Password.ChangedAt, before, after)
		}, retryDuration, tick)
	})
}

func TestServer_TesMachinetUserReduces(t *testing.T) {
	instanceID := Instance.ID()
	orgID := Instance.DefaultOrg.Id

	userRepo := repository.UserRepository()

	t.Run("test machine user add reduced", func(t *testing.T) {
		machineUserRequest := &management.AddMachineUserRequest{
			UserId:          gu.Ptr(gofakeit.Name()),
			UserName:        gofakeit.Name(),
			Name:            gofakeit.Name(),
			Description:     gofakeit.Blurb(),
			AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER,
		}

		before := time.Now()
		resp, err := MgmtClient.AddMachineUser(IAMCTX, machineUserRequest)
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			// event user.machine.added
			// user
			assert.Equal(t, instanceID, user.InstanceID)
			assert.Equal(t, orgID, user.OrganizationID)
			assert.Equal(t, *machineUserRequest.UserId, user.ID)
			assert.Equal(t, machineUserRequest.UserName, user.Username)
			assert.Equal(t, domain.UserStateActive, user.State)
			assert.WithinRange(t, user.UpdatedAt, before, after)
			assert.WithinRange(t, user.CreatedAt, before, after)
			// machine
			assert.Equal(t, machineUserRequest.Name, user.Machine.Name)
			assert.Equal(t, machineUserRequest.Description, user.Machine.Description)
			assert.Equal(t, domain.AccessTokenTypeBearer, user.Machine.AccessTokenType)
		}, retryDuration, tick)
	})

	t.Run("test machine user change reduced", func(t *testing.T) {
		userID := gofakeit.Name()
		machineUserRequest := &management.AddMachineUserRequest{
			UserId:          &userID,
			UserName:        gofakeit.Name(),
			Name:            gofakeit.Name(),
			Description:     gofakeit.Blurb(),
			AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER,
		}

		_, err := MgmtClient.AddMachineUser(IAMCTX, machineUserRequest)
		require.NoError(t, err)

		updateMachineUserReq := &management.UpdateMachineRequest{
			UserId:          userID,
			Name:            gofakeit.Name(),
			Description:     gofakeit.Blurb(),
			AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
		}

		before := time.Now()
		_, err = MgmtClient.UpdateMachine(IAMCTX, updateMachineUserReq)
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.Get(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID),
					userRepo.OrganizationIDCondition(orgID),
					userRepo.IDCondition(userID),
				)),
			)
			require.NoError(t, err)

			// event user.machine.changed
			// user
			assert.WithinRange(t, user.UpdatedAt, before, after)
			// machine
			assert.Equal(t, updateMachineUserReq.Name, user.Machine.Name)
			assert.Equal(t, machineUserRequest.Description, user.Machine.Description)
			assert.Equal(t, domain.AccessTokenTypeJWT, user.Machine.AccessTokenType)
		}, retryDuration, tick)
	})
}
