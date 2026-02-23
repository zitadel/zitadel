//go:build integration

package events_test

// The tests of this file fail because we would need to mock an ldap, OIDC and Oauth server.

// import (
// 	"net"
// 	"net/url"
// 	"testing"
// 	"time"

// 	"github.com/brianvoe/gofakeit/v6"
// 	"github.com/muhlemmer/gu"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	durationpb "google.golang.org/protobuf/types/known/durationpb"

// 	"github.com/zitadel/zitadel/backend/v3/domain"
// 	"github.com/zitadel/zitadel/backend/v3/storage/database"
// 	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
// 	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
// 	"github.com/zitadel/zitadel/internal/integration"
// 	"github.com/zitadel/zitadel/internal/integration/sink"
// 	"github.com/zitadel/zitadel/pkg/grpc/admin"
// 	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
// 	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
// )

// func TestServer_IDPIntentReduces(t *testing.T) {
// 	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Minute)

// 	idpIntentRepo := repository.IDPIntentRepository()
// 	instanceID := Instance.ID()

// 	t.Run("start LDAP ID provider intent should create intent and succeed", func(t *testing.T) {
// 		ldapProvider, err := AdminClient.AddLDAPProvider(IAMCTX, defaultInstanceLDAPRequest(gofakeit.Name()))
// 		require.NoError(t, err)

// 		req := &user.StartIdentityProviderIntentRequest{
// 			IdpId: ldapProvider.GetId(),
// 			Content: &user.StartIdentityProviderIntentRequest_Ldap{
// 				Ldap: &user.LDAPCredentials{
// 					Username: "some-user",
// 					Password: "some-pass",
// 				},
// 			},
// 		}
// 		res, err := UserClient.StartIdentityProviderIntent(IAMCTX, req)
// 		require.NoError(t, err)
// 		idpIntent, ok := res.GetNextStep().(*user.StartIdentityProviderIntentResponse_IdpIntent)
// 		require.True(t, ok)
// 		t.Cleanup(func() {
// 			idpIntentRepo.Delete(IAMCTX, pool, idpIntentRepo.PrimaryKeyCondition(instanceID, idpIntent.IdpIntent.GetIdpIntentId()))
// 			AdminClient.RemoveIDP(IAMCTX, &admin.RemoveIDPRequest{IdpId: ldapProvider.GetId()})
// 		})
// 		require.EventuallyWithT(t, func(collect *assert.CollectT) {
// 			retrievedRelationalIntent, err := idpIntentRepo.Get(
// 				IAMCTX,
// 				pool,
// 				database.WithCondition(database.And(
// 					idpIntentRepo.PrimaryKeyCondition(instanceID, idpIntent.IdpIntent.GetIdpIntentId()),
// 					idpIntentRepo.StateCondition(domain.IDPIntentStateSucceeded),
// 				)),
// 			)
// 			require.NoError(collect, err)
// 			assert.Equal(collect, domain.IDPIntentStateSucceeded, retrievedRelationalIntent.State)
// 			assert.Empty(collect, retrievedRelationalIntent.SuccessURL)
// 			assert.Empty(collect, retrievedRelationalIntent.FailureURL)
// 			assert.Equal(collect, ldapProvider.GetId(), retrievedRelationalIntent.IDPID)
// 			assert.Nil(collect, retrievedRelationalIntent.IDPArguments)
// 			assert.NotEmpty(collect, retrievedRelationalIntent.IDPUser)
// 			assert.NotEmpty(collect, retrievedRelationalIntent.IDPUserID)
// 			assert.NotEmpty(collect, retrievedRelationalIntent.EntryAttributes)
// 			assert.True(collect, retrievedRelationalIntent.ExpiresAt.IsZero())
// 		}, retryDuration, tick)
// 	})

// 	t.Run("oidc intent flow started with URLs should create intent in started state", func(t *testing.T) {
// 		idpCreateReq := defaultOIDCIDReq(gofakeit.Name(), true)
// 		idpCreateReq.Issuer = "http://localhost:8082"
// 		oidcIDP, err := AdminClient.AddOIDCIDP(IAMCTX, idpCreateReq)
// 		require.NoError(t, err)

// 		req := &user.StartIdentityProviderIntentRequest{
// 			IdpId: oidcIDP.GetIdpId(),
// 			Content: &user.StartIdentityProviderIntentRequest_Urls{
// 				Urls: &user.RedirectURLs{
// 					SuccessUrl: "https://localhost:8081/success",
// 					FailureUrl: "https://localhost:8081/fail",
// 				},
// 			},
// 		}
// 		res, err := UserClient.StartIdentityProviderIntent(IAMCTX, req)
// 		require.NoError(t, err)
// 		authURL, ok := res.GetNextStep().(*user.StartIdentityProviderIntentResponse_AuthUrl)
// 		require.True(t, ok)
// 		parsedAuthURL, err := url.Parse(authURL.AuthUrl)
// 		require.NoError(t, err)
// 		idpIntentID := parsedAuthURL.Query().Get("state")
// 		t.Cleanup(func() {
// 			idpIntentRepo.Delete(IAMCTX, pool, idpIntentRepo.PrimaryKeyCondition(instanceID, idpIntentID))
// 			AdminClient.RemoveIDP(IAMCTX, &admin.RemoveIDPRequest{IdpId: oidcIDP.GetIdpId()})
// 		})
// 		require.EventuallyWithT(t, func(collect *assert.CollectT) {
// 			retrievedRelationalIntent, err := idpIntentRepo.Get(
// 				IAMCTX,
// 				pool,
// 				database.WithCondition(idpIntentRepo.PrimaryKeyCondition(instanceID, idpIntentID)),
// 			)
// 			require.NoError(collect, err)
// 			assert.Equal(collect, domain.IDPIntentStateStarted, retrievedRelationalIntent.State)
// 			assert.Equal(collect, "https://localhost:8081/success", retrievedRelationalIntent.SuccessURL.String())
// 			assert.Equal(collect, "https://localhost:8081/fail", retrievedRelationalIntent.FailureURL.String())
// 			assert.Empty(collect, retrievedRelationalIntent.IDPArguments)
// 		}, retryDuration, tick)

// 	})

// 	t.Run("when successful oidc intent should update intent state to succeeded", func(t *testing.T) {
// 		beforeCreate := time.Now()
// 		idpCreateReq := defaultOIDCIDReq(gofakeit.Name(), true)
// 		idpCreateReq.Issuer = "http://localhost:8082"
// 		oidcIDP, err := AdminClient.AddOIDCIDP(IAMCTX, idpCreateReq)
// 		require.NoError(t, err)

// 		expiryDate := time.Now().Add(time.Minute * 30)
// 		idpIntentID, _, _, _, err := sink.SuccessfulOIDCIntent(instanceID, oidcIDP.GetIdpId(), "some-idp-user", "some-user", expiryDate)
// 		require.NoError(t, err)
// 		afterSuccess := time.Now()
// 		t.Cleanup(func() {
// 			idpIntentRepo.Delete(IAMCTX, pool, idpIntentRepo.PrimaryKeyCondition(instanceID, idpIntentID))
// 			AdminClient.RemoveIDP(IAMCTX, &admin.RemoveIDPRequest{IdpId: oidcIDP.GetIdpId()})
// 		})
// 		require.EventuallyWithT(t, func(collect *assert.CollectT) {
// 			retrievedRelationalIntent, err := idpIntentRepo.Get(
// 				IAMCTX,
// 				pool,
// 				database.WithCondition(idpIntentRepo.PrimaryKeyCondition(instanceID, idpIntentID)),
// 			)
// 			require.NoError(collect, err)
// 			assert.Equal(collect, domain.IDPIntentStateSucceeded, retrievedRelationalIntent.State)
// 			assert.NotEmpty(collect, retrievedRelationalIntent.IDPUser)
// 			assert.Equal(collect, "some-idp-user", retrievedRelationalIntent.IDPUserID)
// 			assert.Equal(collect, "username", retrievedRelationalIntent.IDPUsername)
// 			assert.Equal(collect, "some-user", retrievedRelationalIntent.UserID)
// 			assert.NotEmpty(collect, retrievedRelationalIntent.IDPAccessToken)
// 			assert.Equal(collect, "idToken", retrievedRelationalIntent.IDPIDToken)
// 			require.NotNil(collect, retrievedRelationalIntent.SucceededAt)
// 			assert.WithinRange(collect, *retrievedRelationalIntent.SucceededAt, beforeCreate, afterSuccess)
// 			require.NotNil(collect, retrievedRelationalIntent.ExpiresAt)
// 			assert.WithinRange(collect, *retrievedRelationalIntent.ExpiresAt, expiryDate.Add(time.Millisecond*-5), expiryDate.Add(time.Second*1))
// 		}, retryDuration, tick)

// 	})

// 	t.Run("when successful SAML intent should update intent state to succeeded", func(t *testing.T) {
// 		beforeCreate := time.Now()
// 		samlIDP, err := AdminClient.AddSAMLProvider(IAMCTX, defaultInstanceSAMLRequest(gofakeit.Name(), gu.Ptr(true)))
// 		require.NoError(t, err)

// 		expiryDate := time.Now().Add(time.Minute * 30)
// 		idpIntentID, _, _, _, err := sink.SuccessfulSAMLIntent(instanceID, samlIDP.GetId(), "some-idp-user", "some-user", "", expiryDate)
// 		require.NoError(t, err)
// 		afterSuccess := time.Now()
// 		t.Cleanup(func() {
// 			idpIntentRepo.Delete(IAMCTX, pool, idpIntentRepo.PrimaryKeyCondition(instanceID, idpIntentID))
// 			AdminClient.RemoveIDP(IAMCTX, &admin.RemoveIDPRequest{IdpId: samlIDP.GetId()})
// 		})
// 		require.EventuallyWithT(t, func(collect *assert.CollectT) {
// 			retrievedRelationalIntent, err := idpIntentRepo.Get(
// 				IAMCTX,
// 				pool,
// 				database.WithCondition(idpIntentRepo.PrimaryKeyCondition(instanceID, idpIntentID)),
// 			)
// 			require.NoError(collect, err)
// 			assert.Equal(collect, domain.IDPIntentStateSucceeded, retrievedRelationalIntent.State)
// 			assert.NotEmpty(collect, retrievedRelationalIntent.IDPUser)
// 			assert.Equal(collect, "some-idp-user", retrievedRelationalIntent.IDPUserID)
// 			assert.Empty(collect, retrievedRelationalIntent.IDPUsername)
// 			assert.Equal(collect, "some-user", retrievedRelationalIntent.UserID)
// 			assert.NotEmpty(collect, retrievedRelationalIntent.Assertion)
// 			require.NotNil(collect, retrievedRelationalIntent.SucceededAt)
// 			assert.WithinRange(collect, *retrievedRelationalIntent.SucceededAt, beforeCreate, afterSuccess)
// 			require.NotNil(collect, retrievedRelationalIntent.ExpiresAt)
// 			assert.WithinRange(collect, *retrievedRelationalIntent.ExpiresAt, expiryDate.Add(time.Millisecond*-5), expiryDate.Add(time.Second*1))

// 		}, retryDuration, tick)
// 	})

// 	t.Run("when emitting SAMLRequest event should update intent with request ID", func(t *testing.T) {
// 		samlIDP, err := AdminClient.AddSAMLProvider(IAMCTX, defaultInstanceSAMLRequest(gofakeit.Name(), gu.Ptr(true)))
// 		require.NoError(t, err)

// 		expiry := time.Now().Add(time.Hour * 1)
// 		req := &user.StartIdentityProviderIntentRequest{
// 			IdpId: samlIDP.GetId(),
// 			Content: &user.StartIdentityProviderIntentRequest_Urls{
// 				Urls: &user.RedirectURLs{
// 					SuccessUrl: "https://localhost:8081/success",
// 					FailureUrl: "https://localhost:8081/fail",
// 				},
// 			},
// 		}
// 		res, err := UserClient.StartIdentityProviderIntent(IAMCTX, req)
// 		require.NoError(t, err)
// 		formData, ok := res.GetNextStep().(*user.StartIdentityProviderIntentResponse_FormData)
// 		require.True(t, ok)

// 		idpIntentID, ok := formData.FormData.GetFields()["RelayState"]
// 		require.True(t, ok)

// 		t.Cleanup(func() {
// 			idpIntentRepo.Delete(IAMCTX, pool, idpIntentRepo.PrimaryKeyCondition(instanceID, idpIntentID))
// 			AdminClient.RemoveIDP(IAMCTX, &admin.RemoveIDPRequest{IdpId: samlIDP.GetId()})
// 		})

// 		idpIntentID, token, _, _, err := sink.SuccessfulSAMLIntent(instanceID, samlIDP.GetId(), "idp-user-id", "some-user", idpIntentID, expiry)
// 		require.NoError(t, err)

// 		_, err = UserClient.RetrieveIdentityProviderIntent(IAMCTX, &user.RetrieveIdentityProviderIntentRequest{
// 			IdpIntentId:    idpIntentID,
// 			IdpIntentToken: token,
// 		})
// 		require.NoError(t, err)
// 		require.EventuallyWithT(t, func(collect *assert.CollectT) {
// 			retrievedRelationalIntent, err := idpIntentRepo.Get(
// 				IAMCTX,
// 				pool,
// 				database.WithCondition(database.And(
// 					idpIntentRepo.PrimaryKeyCondition(instanceID, idpIntentID),
// 				)),
// 			)
// 			require.NoError(collect, err)
// 			assert.NotEmpty(collect, retrievedRelationalIntent.RequestID)
// 		}, retryDuration, tick)

// 	})

// 	t.Run("when LDAP intent flow fails should update intent state and reason", func(t *testing.T) {
// 		beforeCreate := time.Now()
// 		ldapProvider, err := AdminClient.AddLDAPProvider(IAMCTX, defaultInstanceLDAPRequest(gofakeit.Name()))
// 		require.NoError(t, err)

// 		req := &user.StartIdentityProviderIntentRequest{
// 			IdpId: ldapProvider.GetId(),
// 			Content: &user.StartIdentityProviderIntentRequest_Ldap{
// 				Ldap: &user.LDAPCredentials{
// 					Username: ldap.FailingUser,
// 					Password: ldap.FailingPsw,
// 				},
// 			},
// 		}
// 		_, err = UserClient.StartIdentityProviderIntent(IAMCTX, req)
// 		require.Error(t, err)
// 		afterFail := time.Now()

// 		var idpIntentID string
// 		t.Cleanup(func() {
// 			idpIntentRepo.Delete(IAMCTX, pool, idpIntentRepo.PrimaryKeyCondition(instanceID, idpIntentID))
// 			AdminClient.RemoveIDP(IAMCTX, &admin.RemoveIDPRequest{IdpId: ldapProvider.GetId()})
// 		})

// 		require.EventuallyWithT(t, func(collect *assert.CollectT) {
// 			retrievedRelationalIntent, err := idpIntentRepo.Get(
// 				IAMCTX,
// 				pool,
// 				database.WithCondition(database.And(
// 					idpIntentRepo.InstanceIDCondition(instanceID),
// 					idpIntentRepo.StateCondition(domain.IDPIntentStateFailed),
// 				)),
// 				database.WithOrderByDescending(idpIntentRepo.CreatedAtColumn()),
// 			)
// 			require.NoError(collect, err)
// 			assert.Equal(collect, domain.IDPIntentStateFailed, retrievedRelationalIntent.State)
// 			assert.NotEmpty(collect, retrievedRelationalIntent.FailReason)
// 			require.NotNil(collect, retrievedRelationalIntent.FailedAt)
// 			assert.WithinRange(collect, *retrievedRelationalIntent.FailedAt, beforeCreate, afterFail)
// 			idpIntentID = retrievedRelationalIntent.ID
// 		}, retryDuration, tick)
// 	})

// 	t.Run("when session update checks intent should delete intent", func(t *testing.T) {
// 		// ==========
// 		//   GIVEN
// 		// ==========

// 		// Create a user
// 		email := integration.Email()
// 		testUser := Instance.CreateUserTypeHuman(IAMCTX, email)
// 		_, err := UserClient.VerifyEmail(IAMCTX, &user.VerifyEmailRequest{
// 			UserId:           testUser.GetId(),
// 			VerificationCode: testUser.GetEmailCode(),
// 		})
// 		require.NoError(t, err)
// 		_, err = UserClient.SetPassword(CTX, &user.SetPasswordRequest{
// 			UserId: testUser.GetId(),
// 			NewPassword: &user.Password{
// 				Password: integration.UserPassword,
// 			},
// 			Verification: nil,
// 		})
// 		require.NoError(t, err)

// 		// Create an ID Provider
// 		idpID := Instance.AddGenericOAuthProvider(IAMCTX, integration.IDPName()).GetId()

// 		// Link the user with the provider
// 		idpUserID := integration.ID()
// 		_, err = UserClient.AddIDPLink(IAMCTX, &user.AddIDPLinkRequest{
// 			UserId: testUser.GetId(),
// 			IdpLink: &user.IDPLink{
// 				IdpId:    idpID,
// 				UserId:   idpUserID,
// 				UserName: integration.Username(),
// 			},
// 		})
// 		require.NoError(t, err)

// 		// Create an IDP intent
// 		_, err = url.Parse(Instance.CreateIntent(CTX, idpID).GetAuthUrl())
// 		require.NoError(t, err)
// 		expiry := time.Now().Add(1 * time.Hour)

// 		// Fake a successful intent flow
// 		intentID, intentToken, _, _, err := sink.SuccessfulOAuthIntent(Instance.ID(), idpID, idpUserID, testUser.GetId(), expiry)
// 		require.NoError(t, err)

// 		// Create a session for the user
// 		createdSessionID := createSessionForUser(t)
// 		require.NoError(t, err)

// 		// ==========
// 		//    TEST
// 		// ==========

// 		// Update the session with user, password and IPDIntent checks
// 		_, err = SessionClient.SetSession(IAMCTX, &session.SetSessionRequest{
// 			SessionId: createdSessionID,
// 			Checks: &session.Checks{
// 				User: &session.CheckUser{
// 					Search: &session.CheckUser_LoginName{
// 						LoginName: email,
// 					},
// 				},
// 				Password: &session.CheckPassword{
// 					Password: integration.UserPassword,
// 				},
// 				IdpIntent: &session.CheckIDPIntent{
// 					IdpIntentId:    intentID,
// 					IdpIntentToken: intentToken,
// 				},
// 			},
// 		})
// 		require.NoError(t, err)

// 		// ==========
// 		//   VERIFY
// 		// ==========
// 		t.Cleanup(func() {
// 			idpIntentRepo.Delete(IAMCTX, pool, idpIntentRepo.PrimaryKeyCondition(instanceID, intentID))
// 			AdminClient.RemoveIDP(IAMCTX, &admin.RemoveIDPRequest{IdpId: idpID})
// 			SessionClient.DeleteSession(IAMCTX, &session.DeleteSessionRequest{SessionId: createdSessionID})
// 		})

// 		require.EventuallyWithT(t, func(collect *assert.CollectT) {
// 			_, err := idpIntentRepo.Get(
// 				IAMCTX,
// 				pool,
// 				database.WithCondition(idpIntentRepo.PrimaryKeyCondition(instanceID, intentID)),
// 			)
// 			assert.Error(collect, err)
// 		}, retryDuration, tick)
// 	})
// }

// func createSessionForUser(t *testing.T) (sessionID string) {
// 	t.Helper()
// 	lifetime := time.Hour
// 	createdSession, err := SessionClient.CreateSession(CTX, &session.CreateSessionRequest{
// 		Checks:     nil,
// 		Metadata:   nil,
// 		Challenges: nil,
// 		UserAgent: &session.UserAgent{
// 			FingerprintId: gu.Ptr(integration.ID()),
// 			Ip:            gu.Ptr(net.IPv4(127, 0, 0, 1).String()),
// 			Description:   gu.Ptr("description"),
// 			Header: map[string]*session.UserAgent_HeaderValues{
// 				"User-Agent": {Values: []string{"ZITADEL-Integration-Test"}},
// 			},
// 		},
// 		Lifetime: durationpb.New(lifetime),
// 	})
// 	require.NoError(t, err)
// 	return createdSession.GetSessionId()
// }
