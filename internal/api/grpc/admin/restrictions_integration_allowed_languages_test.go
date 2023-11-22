//go:build integration

package admin_test

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/text"
	"github.com/zitadel/zitadel/pkg/grpc/user"
	"golang.org/x/text/language"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestServer_Restrictions_AllowedLanguages(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var (
		defaultLanguage       = language.German
		allowedLanguage       = defaultLanguage
		supportedLanguagesStr = []string{language.German.String(), language.English.String(), language.Japanese.String()}
		disallowedLanguage    = language.Japanese
		someSupportedLanguage = language.Spanish
	)

	domain, _, iamOwnerCtx := Tester.UseIsolatedInstance(ctx, SystemCTX)
	t.Run("assumed defaults are correct", func(tt *testing.T) {
		tt.Run("languages are not restricted by default", func(ttt *testing.T) {
			restrictions, err := Tester.Client.Admin.GetRestrictions(iamOwnerCtx, &admin.GetRestrictionsRequest{})
			require.NoError(ttt, err)
			require.Len(ttt, restrictions.AllowedLanguages, 0)
		})
		tt.Run("default language is English by default", func(ttt *testing.T) {
			defaultLang, err := Tester.Client.Admin.GetDefaultLanguage(iamOwnerCtx, &admin.GetDefaultLanguageRequest{})
			require.NoError(ttt, err)
			require.Equal(ttt, language.Make(defaultLang.Language), language.English)
		})
		tt.Run("the discovery endpoint returns all supported languages", func(ttt *testing.T) {
			checkDiscoveryEndpoint(ttt, domain, supportedLanguagesStr, nil)
		})
	})
	t.Run("setting disallowed language as default fails", func(tt *testing.T) {
		_, err := Tester.Client.Admin.SetRestrictions(iamOwnerCtx, &admin.SetRestrictionsRequest{AllowedLanguages: &admin.SelectLanguages{List: []string{allowedLanguage.String()}}})
		require.Error(tt, err)
	})
	var importedUser *management.ImportHumanUserResponse
	t.Run("preferred languages that are disallowed are set to the instances default language", func(tt *testing.T) {
		_, err := Tester.Client.Admin.SetDefaultLanguage(iamOwnerCtx, &admin.SetDefaultLanguageRequest{Language: defaultLanguage.String()})
		require.NoError(tt, err)
		importedUser, err = importUser(iamOwnerCtx, disallowedLanguage)
		require.NoError(tt, err)
		awaitPreferredLanguage(iamOwnerCtx, tt, importedUser.GetUserId(), disallowedLanguage)
		_, err = Tester.Client.Admin.SetRestrictions(iamOwnerCtx, &admin.SetRestrictionsRequest{AllowedLanguages: &admin.SelectLanguages{List: []string{allowedLanguage.String()}}})
		require.NoError(tt, err)
		awaitPreferredLanguage(iamOwnerCtx, tt, importedUser.GetUserId(), defaultLanguage)
	})
	t.Run("using a disallowed language is limited", func(tt *testing.T) {
		assertIfLanguageIsUsable(iamOwnerCtx, tt, domain, supportedLanguagesStr, allowedLanguage, disallowedLanguage, importedUser.GetUserId(), false)
		tt.Run("the list of supported languages includes the disallowed languages", func(ttt *testing.T) {
			supported, err := Tester.Client.Admin.GetSupportedLanguages(iamOwnerCtx, &admin.GetSupportedLanguagesRequest{})
			require.NoError(ttt, err)
			require.Condition(ttt, contains(supported.GetLanguages(), supportedLanguagesStr))
		})
		tt.Run("setting a message text for a disallowed language still works, so they can be set before a language is allowed", func(ttt *testing.T) {
			_, err := Tester.Client.Admin.SetCustomLoginText(iamOwnerCtx, &admin.SetCustomLoginTextsRequest{
				Language: disallowedLanguage.String(),
				EmailVerificationText: &text.EmailVerificationScreenText{
					Description: "hodor",
				},
			})
			assert.NoError(ttt, err)
			_, err = Tester.Client.Mgmt.SetCustomLoginText(iamOwnerCtx, &management.SetCustomLoginTextsRequest{
				Language: disallowedLanguage.String(),
				EmailVerificationText: &text.EmailVerificationScreenText{
					Description: "hodor",
				},
			})
			assert.NoError(ttt, err)
			_, err = Tester.Client.Mgmt.SetCustomInitMessageText(iamOwnerCtx, &management.SetCustomInitMessageTextRequest{
				Language: disallowedLanguage.String(),
				Text:     "hodor",
			})
			assert.NoError(ttt, err)
			_, err = Tester.Client.Admin.SetDefaultInitMessageText(iamOwnerCtx, &admin.SetDefaultInitMessageTextRequest{
				Language: disallowedLanguage.String(),
				Text:     "hodor",
			})
			assert.NoError(ttt, err)
		})
	})
	t.Run("allowing the language makes it usable again", func(tt *testing.T) {
		setAndAwaitAllowedLanguages(iamOwnerCtx, tt, []string{allowedLanguage.String(), disallowedLanguage.String()})
		assertIfLanguageIsUsable(iamOwnerCtx, tt, domain, supportedLanguagesStr, allowedLanguage, someSupportedLanguage, importedUser.GetUserId(), false)
		assertIfLanguageIsUsable(iamOwnerCtx, tt, domain, supportedLanguagesStr, allowedLanguage, disallowedLanguage, importedUser.GetUserId(), true)
	})
	t.Run("setting the languages to the empty list allows all supported languages", func(tt *testing.T) {
		setAndAwaitAllowedLanguages(iamOwnerCtx, tt, []string{})
		assertIfLanguageIsUsable(iamOwnerCtx, tt, domain, supportedLanguagesStr, allowedLanguage, someSupportedLanguage, importedUser.GetUserId(), true)
		checkDiscoveryEndpoint(tt, domain, supportedLanguagesStr, nil)
	})
}

func setAndAwaitAllowedLanguages(ctx context.Context, t *testing.T, selectLanguages []string) {
	_, err := Tester.Client.Admin.SetRestrictions(ctx, &admin.SetRestrictionsRequest{AllowedLanguages: &admin.SelectLanguages{List: selectLanguages}})
	assert.NoError(t, err)
	awaitCtx, awaitCancel := context.WithTimeout(ctx, 100*time.Second)
	defer awaitCancel()
	await(t, awaitCtx, func() bool {
		restrictions, getErr := Tester.Client.Admin.GetRestrictions(awaitCtx, &admin.GetRestrictionsRequest{})
		expectLanguages := selectLanguages
		if len(selectLanguages) == 0 {
			expectLanguages = nil
		}
		return assert.NoError(NoopAssertionT, getErr) &&
			assert.Equal(NoopAssertionT, expectLanguages, restrictions.GetAllowedLanguages())
	})
}

func assertIfLanguageIsUsable(ctx context.Context, t *testing.T, domain string, supportedLanguagesStr []string, allowedLanguage, disallowedLanguage language.Tag, existingUserID string, usable bool) {
	var (
		usableStr               = "usable"
		checkAPIErr             = require.NoError
		allowListMustContain    = supportedLanguagesStr
		allowListMustNotContain []string
	)

	if !usable {
		usableStr = "unusable"
		checkAPIErr = require.Error
		allowListMustContain = []string{allowedLanguage.String()}
		allowListMustNotContain = []string{disallowedLanguage.String()}
	}

	t.Run("the disallowed language is "+usableStr, func(tt *testing.T) {
		tt.Run("as an imported users preferred language", func(ttt *testing.T) {
			_, err := importUser(ctx, disallowedLanguage)
			checkAPIErr(ttt, err)
		})
		tt.Run("as preferred languang when updating a users profile", func(ttt *testing.T) {
			_, err := Tester.Client.Mgmt.UpdateHumanProfile(ctx, &management.UpdateHumanProfileRequest{
				UserId:            existingUserID,
				FirstName:         "hodor",
				LastName:          "hodor",
				NickName:          integration.RandString(5),
				DisplayName:       "hodor",
				PreferredLanguage: disallowedLanguage.String(),
				Gender:            user.Gender_GENDER_MALE,
			})
			checkAPIErr(ttt, err)
		})
		tt.Run("as ui locale according to the oidc discovery endpoint", func(ttt *testing.T) {
			checkDiscoveryEndpoint(ttt, domain, allowListMustContain, allowListMustNotContain)
		})
		tt.Run("for consumers of the GetAllowedLanguages endpoints", func(ttt *testing.T) {
			adminAllowed, err := Tester.Client.Admin.GetAllowedLanguages(ctx, &admin.GetAllowedLanguagesRequest{})
			require.NoError(ttt, err)
			assertAllowList(adminAllowed.Languages, allowListMustContain, allowListMustNotContain)
			mgmtAllowed, err := Tester.Client.Mgmt.GetAllowedLanguages(ctx, &management.GetAllowedLanguagesRequest{})
			require.NoError(ttt, err)
			assertAllowList(mgmtAllowed.Languages, allowListMustContain, allowListMustNotContain)
			authAllowed, err := Tester.Client.Auth.GetAllowedLanguages(ctx, &auth.GetAllowedLanguagesRequest{})
			require.NoError(ttt, err)
			assertAllowList(authAllowed.Languages, allowListMustContain, allowListMustNotContain)
		})
	})
}

func importUser(ctx context.Context, preferredLanguage language.Tag) (*management.ImportHumanUserResponse, error) {
	random := integration.RandString(5)
	return Tester.Client.Mgmt.ImportHumanUser(ctx, &management.ImportHumanUserRequest{
		UserName: "integration-test-user_" + random,
		Profile: &management.ImportHumanUserRequest_Profile{
			FirstName:         "hodor",
			LastName:          "hodor",
			NickName:          "hodor",
			PreferredLanguage: preferredLanguage.String(),
		},
		Email: &management.ImportHumanUserRequest_Email{
			Email:           random + "@hodor.hodor",
			IsEmailVerified: true,
		},
	})
}

func awaitPreferredLanguage(ctx context.Context, t *testing.T, userID string, preferredLanguage language.Tag) {
	awaitCtx, awaitCancel := context.WithTimeout(ctx, 10*time.Second)
	defer awaitCancel()
	await(t, awaitCtx, func() bool {
		resp, getErr := Tester.Client.Mgmt.GetHumanProfile(awaitCtx, &management.GetHumanProfileRequest{UserId: userID})
		return assert.NoError(NoopAssertionT, getErr) &&
			assert.Equal(NoopAssertionT, preferredLanguage, language.Make(resp.GetProfile().GetPreferredLanguage()))
	})

}

func checkDiscoveryEndpoint(t *testing.T, domain string, containsUILocales, notContainsUILocales []string) {
	endpoint, err := url.Parse("http://" + domain + ":8080/.well-known/openid-configuration")
	require.NoError(t, err)
	resp, err := http.Get(endpoint.String())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	doc := struct {
		UILocalesSupported []string `json:"ui_locales_supported"`
	}{}
	require.NoError(t, json.Unmarshal(body, &doc))
	assertAllowList(doc.UILocalesSupported, containsUILocales, notContainsUILocales)
}

func assertAllowList(allowList, mustContain, mustNotContain []string) {
	if mustContain != nil {
		assert.Condition(NoopAssertionT, contains(allowList, mustContain))
	}
	if mustNotContain != nil {
		assert.Condition(NoopAssertionT, not(contains(allowList, mustNotContain)))
	}
}

// We would love to use assert.Contains here, but it doesn't work with slices of strings
func contains(container []string, subset []string) assert.Comparison {
	return func() bool {
		if subset == nil {
			return true
		}
		for _, str := range subset {
			var found bool
			for _, containerStr := range container {
				if str == containerStr {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
}

func not(cmp assert.Comparison) assert.Comparison {
	return func() bool {
		return !cmp()
	}
}
