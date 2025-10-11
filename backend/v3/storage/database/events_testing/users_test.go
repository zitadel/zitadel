//go:build integration

package events_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/user"
	v2beta_user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
	// user "github.com/zitadel/zitadel/pkg/grpc/user"
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
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			// event user.added
			// event user.human.added
			// domain.User
			assert.Equal(t, instanceID, user.InstanceID)
			assert.Equal(t, orgID, user.OrgID)
			assert.Equal(t, resp.UserId, user.ID)
			assert.Equal(t, humanUserRequest.UserName, user.Username)
			assert.Equal(t, domain.UserStateActive, user.State)
			// TODO
			// assert.Equal(t, true, user.UsernameOrgUnique)
			assert.WithinRange(t, user.UpdatedAt, before, after)
			assert.WithinRange(t, user.CreatedAt, before, after)
			// Email
			assert.Equal(t, domain.ContactTypeEmail, *user.HumanEmailContact.Type)
			assert.Equal(t, humanUserRequest.Email.Email, *user.HumanEmailContact.Value)
			// TODO
			// assert.Equal(t, true, *user.HumanEmailContact.IsVerified)
			assert.Nil(t, user.HumanEmailContact.UnverifiedValue)
			// Phone
			assert.Equal(t, domain.ContactTypePhone, *user.HumanPhoneContact.Type)
			assert.Equal(t, humanUserRequest.Phone.Phone, *user.HumanPhoneContact.Value)
			// TODO
			// assert.Equal(t, true, *user.HumanPhoneContact.IsVerified)
			assert.Nil(t, user.HumanPhoneContact.UnverifiedValue)
			// Human
			assert.Equal(t, humanUserRequest.Profile.FirstName, user.FirstName)
			assert.Equal(t, humanUserRequest.Profile.LastName, user.LastName)
			assert.Equal(t, humanUserRequest.Profile.NickName, user.NickName)
			assert.Equal(t, humanUserRequest.Profile.DisplayName, user.DisplayName)
			assert.Equal(t, humanUserRequest.Profile.PreferredLanguage, user.PreferredLanguage)
			assert.Equal(t, uint8(humanUserRequest.Profile.Gender), user.Gender)
		}, retryDuration, tick)
	})

	// TODO
	// // ~/go/src/zitadel/internal/api/ui/login/register_org_handler.go
	// t.Run("test human user register reduced", func(t *testing.T) {
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

	// 	userAddRequest := &v2User.AddHumanUserRequest{
	// 		UserId:   gu.Ptr("userID"),
	// 		Username: gu.Ptr(gofakeit.Username()),
	// 		Organization: &v2Object.Organization{
	// 			Org: &v2Object.Organization_OrgId{
	// 				OrgId: orgID,
	// 			},
	// 		},
	// 		Profile: &v2User.SetHumanProfile{
	// 			GivenName:         gofakeit.FirstName(),
	// 			FamilyName:        gofakeit.LastName(),
	// 			NickName:          gu.Ptr(gofakeit.Username()),
	// 			DisplayName:       gu.Ptr(gofakeit.Name()),
	// 			PreferredLanguage: gu.Ptr("en"),
	// 			Gender:            gu.Ptr(v2User.Gender_GENDER_MALE),
	// 		},
	// 		Email: &v2User.SetHumanEmail{
	// 			Email: gofakeit.Email(),
	// 			Verification: &v2User.SetHumanEmail_IsVerified{
	// 				IsVerified: true,
	// 			},
	// 		},
	// 		Phone: &v2User.SetHumanPhone{
	// 			Phone: "+" + gofakeit.Phone(),
	// 			Verification: &v2User.SetHumanPhone_IsVerified{
	// 				IsVerified: true,
	// 			},
	// 		},
	// 	}

	// 	// resp, err := MgmtClient.AddHumanUser(IAMCTX, humanUserRequest)
	// 	before := time.Now()
	// 	resp, err := UserClient.AddHumanUser(CTX, userAddRequest)
	// 	fmt.Printf("[DEBUGPRINT] [users_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> respo = %+v\n", resp)
	// 	fmt.Printf("[DEBUGPRINT] [users_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> err = %+v\n", err)
	// 	return
	// 	require.NoError(t, err)
	// 	after := time.Now()

	// 	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
	// 	assert.EventuallyWithT(t, func(t *assert.CollectT) {
	// 		user, err := userRepo.GetHuman(
	// 			CTX,
	// 			pool,
	// 			database.WithCondition(database.And(
	// 				userRepo.Human().InstanceIDCondition(instanceID),
	// 				userRepo.Human().OrgIDCondition(orgID),
	// 				userRepo.Human().IDCondition(resp.UserId),
	// 			)),
	// 		)
	// 		require.NoError(t, err)

	// 		// event user.added
	// 		// event user.human.added
	// 		// domain.User
	// 		assert.Equal(t, instanceID, user.InstanceID)
	// 		assert.Equal(t, orgID, user.OrgID)
	// 		assert.Equal(t, resp.UserId, user.ID)
	// 		assert.Equal(t, humanUserRequest.UserName, user.Username)
	// 		assert.Equal(t, domain.UserStateActive, user.State)
	// 		// TODO
	// 		// assert.Equal(t, true, user.UsernameOrgUnique)
	// 		assert.WithinRange(t, user.UpdatedAt, before, after)
	// 		assert.WithinRange(t, user.CreatedAt, before, after)
	// 		// Email
	// 		assert.Equal(t, domain.ContactTypeEmail, *user.HumanEmailContact.Type)
	// 		assert.Equal(t, humanUserRequest.Email.Email, *user.HumanEmailContact.Value)
	// 		// TODO
	// 		// assert.Equal(t, true, *user.HumanEmailContact.IsVerified)
	// 		assert.Nil(t, user.HumanEmailContact.UnverifiedValue)
	// 		// Phone
	// 		assert.Equal(t, domain.ContactTypePhone, *user.HumanPhoneContact.Type)
	// 		assert.Equal(t, humanUserRequest.Phone.Phone, *user.HumanPhoneContact.Value)
	// 		// TODO
	// 		// assert.Equal(t, true, *user.HumanPhoneContact.IsVerified)
	// 		assert.Nil(t, user.HumanPhoneContact.UnverifiedValue)
	// 		// Human
	// 		assert.Equal(t, humanUserRequest.Profile.FirstName, user.FirstName)
	// 		assert.Equal(t, humanUserRequest.Profile.LastName, user.LastName)
	// 		assert.Equal(t, humanUserRequest.Profile.NickName, user.NickName)
	// 		assert.Equal(t, humanUserRequest.Profile.DisplayName, user.DisplayName)
	// 		assert.Equal(t, humanUserRequest.Profile.PreferredLanguage, user.PreferredLanguage)
	// 		assert.Equal(t, uint8(humanUserRequest.Profile.Gender), user.Gender)
	// 	}, retryDuration, tick)
	// })

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
	// 		user, err := userRepo.GetHuman(
	// 			CTX,
	// 			pool,
	// 			database.WithCondition(database.And(
	// 				userRepo.Human().InstanceIDCondition(instanceID),
	// 				userRepo.Human().OrgIDCondition(orgID),
	// 				userRepo.Human().IDCondition(resp.UserId),
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
	// 		user, err := userRepo.GetHuman(
	// 			CTX,
	// 			pool,
	// 			database.WithCondition(database.And(
	// 				userRepo.Human().InstanceIDCondition(instanceID),
	// 				userRepo.Human().OrgIDCondition(orgID),
	// 				userRepo.Human().IDCondition(resp.UserId),
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
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
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
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
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
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
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
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
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
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
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
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
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
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
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
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
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
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
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
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
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
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			assert.NotNil(t, user)
		}, retryDuration, tick)

		_, err = UserClient.DeleteUser(IAMCTX, &v2beta_user.DeleteUserRequest{
			UserId: resp.UserId,
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
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
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			assert.NotNil(t, user)
		}, retryDuration, tick)

		username := gofakeit.Username()
		_, err = UserClient.UpdateHumanUser(IAMCTX, &v2beta_user.UpdateHumanUserRequest{
			UserId:   resp.UserId,
			Username: &username,
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			// event user.username.changed
			assert.Equal(t, username, user.Username)
		}, retryDuration, tick)
	})

	// TODO
	// t.Run("test user domain claimed reduced", func(t *testing.T) {
	// 	org, err := OrgClient.CreateOrganization(CTX, &v2beta_org.CreateOrganizationRequest{
	// 		Name: gofakeit.Name(),
	// 	})
	// 	require.NoError(t, err)
	// 	orgID := org.Id

	// 	_, err = AdminClient.AddCustomDomainPolicy(IAMCTX, &admin.AddCustomDomainPolicyRequest{
	// 		OrgId:              orgID,
	// 		ValidateOrgDomains: true,
	// 	})
	// 	require.NoError(t, err)

	// 	humanUserRequest := &management.AddHumanUserRequest{
	// 		UserName: gofakeit.Username(),
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
	// 		user, err := userRepo.GetHuman(
	// 			CTX,
	// 			pool,
	// 			database.WithCondition(database.And(
	// 				userRepo.Human().InstanceIDCondition(instanceID),
	// 				userRepo.Human().OrgIDCondition(orgID),
	// 				userRepo.Human().IDCondition(resp.UserId),
	// 			)),
	// 		)
	// 		require.NoError(t, err)

	// 		assert.NotNil(t, user)
	// 	}, retryDuration, tick)

	// 	username := gofakeit.Username()
	// 	// _, err = UserClient.UpdateHumanUser(IAMCTX, &v2beta_user.UpdateHumanUserRequest{
	// 	_, err = OrgClient.AddOrganizationDomain(IAMCTX, &v2beta_org.AddOrganizationDomainRequest{
	// 		OrganizationId: orgID,
	// 		Domain:         "www.example.com",
	// 	})
	// 	require.NoError(t, err)

	// 	r, err := OrgClient.GenerateOrganizationDomainValidation(IAMCTX, &v2beta_org.GenerateOrganizationDomainValidationRequest{
	// 		OrganizationId: orgID,
	// 		Domain:         "www.example.com",
	// 		Type:           v2beta_org.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP,
	// 	})
	// 	require.NoError(t, err)
	// 	fmt.Printf("[DEBUGPRINT] [users_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> r = %+v\n", r)

	// 	_, err = OrgClient.VerifyOrganizationDomain(IAMCTX, &v2beta_org.VerifyOrganizationDomainRequest{
	// 		OrganizationId: orgID,
	// 		Domain:         "www.example.com",
	// 	})
	// 	require.NoError(t, err)

	// 	retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
	// 	assert.EventuallyWithT(t, func(t *assert.CollectT) {
	// 		user, err := userRepo.GetHuman(
	// 			CTX,
	// 			pool,
	// 			database.WithCondition(database.And(
	// 				userRepo.Human().InstanceIDCondition(instanceID),
	// 				userRepo.Human().OrgIDCondition(orgID),
	// 				userRepo.Human().IDCondition(resp.UserId),
	// 			)),
	// 		)
	// 		require.NoError(t, err)

	// 		// event user.domain.claimed
	// 		assert.Equal(t, username, user.Username)
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
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			assert.NotNil(t, user)
		}, retryDuration, tick)

		profile := &v2beta_user.SetHumanProfile{
			GivenName:         gofakeit.FirstName(),
			FamilyName:        gofakeit.LastName(),
			NickName:          gu.Ptr(gofakeit.Color()),
			DisplayName:       gu.Ptr(gofakeit.Name()),
			PreferredLanguage: gu.Ptr("fr"),
			Gender:            gu.Ptr(v2beta_user.Gender_GENDER_FEMALE),
		}
		_, err = UserClient.UpdateHumanUser(IAMCTX, &v2beta_user.UpdateHumanUserRequest{
			UserId:  resp.UserId,
			Profile: profile,
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			// event user.human.profile.changed
			assert.Equal(t, profile.GivenName, user.FirstName)
			assert.Equal(t, profile.FamilyName, user.LastName)
			assert.Equal(t, *profile.NickName, user.NickName)
			assert.Equal(t, *profile.DisplayName, user.DisplayName)
			assert.Equal(t, *profile.PreferredLanguage, user.PreferredLanguage)
			assert.Equal(t, uint8(*profile.Gender), user.Gender)
		}, retryDuration, tick)
	})

	// TODO
	// t.Run("test human user avatar added reduced", func(t *testing.T) {
	// 	humanUserRequest := &management.AddHumanUserRequest{
	// 		UserName: gofakeit.Username(),
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
	// 		user, err := userRepo.GetHuman(
	// 			CTX,
	// 			pool,
	// 			database.WithCondition(database.And(
	// 				userRepo.Human().InstanceIDCondition(instanceID),
	// 				userRepo.Human().OrgIDCondition(orgID),
	// 				userRepo.Human().IDCondition(resp.UserId),
	// 			)),
	// 		)
	// 		require.NoError(t, err)

	// 		assert.NotNil(t, user)
	// 	}, retryDuration, tick)

	// 	profile := &v2beta_user.SetHumanProfile{
	// 		GivenName:         gofakeit.FirstName(),
	// 		FamilyName:        gofakeit.LastName(),
	// 		NickName:          gu.Ptr(gofakeit.Color()),
	// 		DisplayName:       gu.Ptr(gofakeit.Name()),
	// 		PreferredLanguage: gu.Ptr("fr"),
	// 		Gender:            gu.Ptr(v2beta_user.Gender_GENDER_FEMALE),
	// 	}
	// 	// _, err = UserClient.UpdateHumanUser(IAMCTX, &v2beta_user.UpdateHumanUserRequest{
	// 	_, err = UserClient.UpdateHumanUser(IAMCTX, &v2beta_user.UpdateHumanUserRequest{
	// 		UserId:  resp.UserId,
	// 		Profile: profile,
	// 	})
	// 	require.NoError(t, err)

	// 	retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
	// 	assert.EventuallyWithT(t, func(t *assert.CollectT) {
	// 		user, err := userRepo.GetHuman(
	// 			CTX,
	// 			pool,
	// 			database.WithCondition(database.And(
	// 				userRepo.Human().InstanceIDCondition(instanceID),
	// 				userRepo.Human().OrgIDCondition(orgID),
	// 				userRepo.Human().IDCondition(resp.UserId),
	// 			)),
	// 		)
	// 		require.NoError(t, err)

	// 		// event user.human.profile.changed
	// 		assert.Equal(t, profile.GivenName, user.FirstName)
	// 		assert.Equal(t, profile.FamilyName, user.LastName)
	// 		assert.Equal(t, *profile.NickName, user.NickName)
	// 		assert.Equal(t, *profile.DisplayName, user.DisplayName)
	// 		assert.Equal(t, *profile.PreferredLanguage, user.PreferredLanguage)
	// 		assert.Equal(t, uint8(*profile.Gender), user.Gender)
	// 	}, retryDuration, tick)
	// })

	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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

		resp, err := MgmtClient.AddMachineUser(IAMCTX, machineUserRequest)
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			assert.NotNil(t, user)
		}, retryDuration, tick)

		profile := &v2beta_user.SetHumanProfile{
			GivenName:         gofakeit.FirstName(),
			FamilyName:        gofakeit.LastName(),
			NickName:          gu.Ptr(gofakeit.Color()),
			DisplayName:       gu.Ptr(gofakeit.Name()),
			PreferredLanguage: gu.Ptr("fr"),
			Gender:            gu.Ptr(v2beta_user.Gender_GENDER_FEMALE),
		}
		_, err = UserClient.UpdateHumanUser(IAMCTX, &v2beta_user.UpdateHumanUserRequest{
			UserId:  resp.UserId,
			Profile: profile,
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			user, err := userRepo.GetHuman(
				CTX,
				pool,
				database.WithCondition(database.And(
					userRepo.Human().InstanceIDCondition(instanceID),
					userRepo.Human().OrgIDCondition(orgID),
					userRepo.Human().IDCondition(resp.UserId),
				)),
			)
			require.NoError(t, err)

			// event user.human.profile.changed
			assert.Equal(t, profile.GivenName, user.FirstName)
			assert.Equal(t, profile.FamilyName, user.LastName)
			assert.Equal(t, *profile.NickName, user.NickName)
			assert.Equal(t, *profile.DisplayName, user.DisplayName)
			assert.Equal(t, *profile.PreferredLanguage, user.PreferredLanguage)
			assert.Equal(t, uint8(*profile.Gender), user.Gender)
		}, retryDuration, tick)
	})
}
