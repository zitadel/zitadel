//go:build integration

package admin_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/text"
	"github.com/zitadel/zitadel/pkg/grpc/user"
)

func TestServer_Restrictions_AllowedLanguages(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	var (
		defaultAndAllowedLanguage = language.German
		supportedLanguagesStr     = []string{language.German.String(), language.English.String(), language.Japanese.String()}
		disallowedLanguage        = language.Spanish
		unsupportedLanguage       = language.Afrikaans
	)

	domain, _, _, iamOwnerCtx := Tester.UseIsolatedInstance(t, ctx, SystemCTX)
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
			awaitDiscoveryEndpoint(ttt, domain, supportedLanguagesStr, nil)
		})
	})
	t.Run("restricting the default language fails", func(tt *testing.T) {
		_, err := Tester.Client.Admin.SetRestrictions(iamOwnerCtx, &admin.SetRestrictionsRequest{AllowedLanguages: &admin.SelectLanguages{List: []string{defaultAndAllowedLanguage.String()}}})
		expectStatus, ok := status.FromError(err)
		require.True(tt, ok)
		require.Equal(tt, codes.FailedPrecondition, expectStatus.Code())
	})
	t.Run("not defining any restrictions throws an error", func(tt *testing.T) {
		_, err := Tester.Client.Admin.SetRestrictions(iamOwnerCtx, &admin.SetRestrictionsRequest{})
		expectStatus, ok := status.FromError(err)
		require.True(tt, ok)
		require.Equal(tt, codes.InvalidArgument, expectStatus.Code())
	})
	t.Run("setting the default language works", func(tt *testing.T) {
		setAndAwaitDefaultLanguage(iamOwnerCtx, tt, defaultAndAllowedLanguage)
	})
	t.Run("restricting allowed languages works", func(tt *testing.T) {
		setAndAwaitAllowedLanguages(iamOwnerCtx, tt, []string{defaultAndAllowedLanguage.String()})
	})
	t.Run("GetAllowedLanguage returns only the allowed languages", func(tt *testing.T) {
		expectContains, expectNotContains := []string{defaultAndAllowedLanguage.String()}, []string{disallowedLanguage.String()}
		adminResp, err := Tester.Client.Admin.GetAllowedLanguages(iamOwnerCtx, &admin.GetAllowedLanguagesRequest{})
		require.NoError(t, err)
		langs := adminResp.GetLanguages()
		assert.Condition(t, contains(langs, expectContains))
		assert.Condition(t, not(contains(langs, expectNotContains)))
	})
	t.Run("setting the default language to a disallowed language fails", func(tt *testing.T) {
		_, err := Tester.Client.Admin.SetDefaultLanguage(iamOwnerCtx, &admin.SetDefaultLanguageRequest{Language: disallowedLanguage.String()})
		expectStatus, ok := status.FromError(err)
		require.True(tt, ok)
		require.Equal(tt, codes.FailedPrecondition, expectStatus.Code())
	})
	t.Run("the list of supported languages includes the disallowed languages", func(tt *testing.T) {
		supported, err := Tester.Client.Admin.GetSupportedLanguages(iamOwnerCtx, &admin.GetSupportedLanguagesRequest{})
		require.NoError(tt, err)
		require.Condition(tt, contains(supported.GetLanguages(), supportedLanguagesStr))
	})
	t.Run("the disallowed language is not listed in the discovery endpoint", func(tt *testing.T) {
		awaitDiscoveryEndpoint(tt, domain, []string{defaultAndAllowedLanguage.String()}, []string{disallowedLanguage.String()})
	})
	t.Run("the login ui is rendered in the default language", func(tt *testing.T) {
		awaitLoginUILanguage(tt, domain, disallowedLanguage, defaultAndAllowedLanguage, "Allgemeine Geschäftsbedingungen und Datenschutz")
	})
	t.Run("preferred languages are not restricted by the supported languages", func(tt *testing.T) {
		tt.Run("change user profile", func(ttt *testing.T) {
			resp, err := Tester.Client.Mgmt.ListUsers(iamOwnerCtx, &management.ListUsersRequest{Queries: []*user.SearchQuery{{Query: &user.SearchQuery_UserNameQuery{UserNameQuery: &user.UserNameQuery{
				UserName: "zitadel-admin@zitadel.localhost"}},
			}}})
			require.NoError(ttt, err)
			require.Len(ttt, resp.GetResult(), 1)
			humanAdmin := resp.GetResult()[0]
			profile := humanAdmin.GetHuman().GetProfile()
			require.NotEqual(ttt, unsupportedLanguage.String(), profile.GetPreferredLanguage())
			_, updateErr := Tester.Client.Mgmt.UpdateHumanProfile(iamOwnerCtx, &management.UpdateHumanProfileRequest{
				PreferredLanguage: unsupportedLanguage.String(),
				UserId:            humanAdmin.GetId(),
				FirstName:         profile.GetFirstName(),
				LastName:          profile.GetLastName(),
				NickName:          profile.GetNickName(),
				DisplayName:       profile.GetDisplayName(),
				Gender:            profile.GetGender(),
			})
			require.NoError(ttt, updateErr)
		})
	})
	t.Run("custom texts are only restricted by the supported languages", func(tt *testing.T) {
		_, err := Tester.Client.Admin.SetCustomLoginText(iamOwnerCtx, &admin.SetCustomLoginTextsRequest{
			Language: disallowedLanguage.String(),
			EmailVerificationText: &text.EmailVerificationScreenText{
				Description: "hodor",
			},
		})
		assert.NoError(tt, err)
		_, err = Tester.Client.Mgmt.SetCustomLoginText(iamOwnerCtx, &management.SetCustomLoginTextsRequest{
			Language: disallowedLanguage.String(),
			EmailVerificationText: &text.EmailVerificationScreenText{
				Description: "hodor",
			},
		})
		assert.NoError(tt, err)
		_, err = Tester.Client.Mgmt.SetCustomInitMessageText(iamOwnerCtx, &management.SetCustomInitMessageTextRequest{
			Language: disallowedLanguage.String(),
			Text:     "hodor",
		})
		assert.NoError(tt, err)
		_, err = Tester.Client.Admin.SetDefaultInitMessageText(iamOwnerCtx, &admin.SetDefaultInitMessageTextRequest{
			Language: disallowedLanguage.String(),
			Text:     "hodor",
		})
		assert.NoError(tt, err)
	})
	t.Run("allowing all languages works", func(tt *testing.T) {
		tt.Run("restricting allowed languages works", func(ttt *testing.T) {
			setAndAwaitAllowedLanguages(iamOwnerCtx, ttt, make([]string, 0))
		})
	})

	t.Run("allowing the language makes it usable again", func(tt *testing.T) {
		tt.Run("the previously disallowed language is listed in the discovery endpoint again", func(ttt *testing.T) {
			awaitDiscoveryEndpoint(ttt, domain, []string{disallowedLanguage.String()}, nil)
		})
		tt.Run("the login ui is rendered in the previously disallowed language", func(ttt *testing.T) {
			awaitLoginUILanguage(ttt, domain, disallowedLanguage, disallowedLanguage, "Términos y condiciones")
		})
	})
}

func setAndAwaitAllowedLanguages(ctx context.Context, t *testing.T, selectLanguages []string) {
	_, err := Tester.Client.Admin.SetRestrictions(ctx, &admin.SetRestrictionsRequest{AllowedLanguages: &admin.SelectLanguages{List: selectLanguages}})
	require.NoError(t, err)
	awaitCtx, awaitCancel := context.WithTimeout(ctx, 10*time.Second)
	defer awaitCancel()
	await(t, awaitCtx, func(tt *assert.CollectT) {
		restrictions, getErr := Tester.Client.Admin.GetRestrictions(awaitCtx, &admin.GetRestrictionsRequest{})
		expectLanguages := selectLanguages
		if len(selectLanguages) == 0 {
			expectLanguages = nil
		}
		assert.NoError(tt, getErr)
		assert.Equal(tt, expectLanguages, restrictions.GetAllowedLanguages())
	})
}

func setAndAwaitDefaultLanguage(ctx context.Context, t *testing.T, lang language.Tag) {
	_, err := Tester.Client.Admin.SetDefaultLanguage(ctx, &admin.SetDefaultLanguageRequest{Language: lang.String()})
	require.NoError(t, err)
	awaitCtx, awaitCancel := context.WithTimeout(ctx, 10*time.Second)
	defer awaitCancel()
	await(t, awaitCtx, func(tt *assert.CollectT) {
		defaultLang, getErr := Tester.Client.Admin.GetDefaultLanguage(awaitCtx, &admin.GetDefaultLanguageRequest{})
		assert.NoError(tt, getErr)
		assert.Equal(tt, lang.String(), defaultLang.GetLanguage())
	})
}

func awaitDiscoveryEndpoint(t *testing.T, domain string, containsUILocales, notContainsUILocales []string) {
	awaitCtx, awaitCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer awaitCancel()
	await(t, awaitCtx, func(tt *assert.CollectT) {
		req, err := http.NewRequestWithContext(awaitCtx, http.MethodGet, "http://"+domain+":8080/.well-known/openid-configuration", nil)
		require.NoError(tt, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(tt, err)
		require.Equal(tt, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		defer func() {
			require.NoError(tt, resp.Body.Close())
		}()
		require.NoError(tt, err)
		doc := struct {
			UILocalesSupported []string `json:"ui_locales_supported"`
		}{}
		require.NoError(tt, json.Unmarshal(body, &doc))
		if containsUILocales != nil {
			assert.Condition(tt, contains(doc.UILocalesSupported, containsUILocales))
		}
		if notContainsUILocales != nil {
			assert.Condition(tt, not(contains(doc.UILocalesSupported, notContainsUILocales)))
		}
	})
}

func awaitLoginUILanguage(t *testing.T, domain string, acceptLanguage language.Tag, expectLang language.Tag, containsText string) {
	awaitCtx, awaitCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer awaitCancel()
	await(t, awaitCtx, func(tt *assert.CollectT) {
		req, err := http.NewRequestWithContext(awaitCtx, http.MethodGet, "http://"+domain+":8080/ui/login/register", nil)
		req.Header.Set("Accept-Language", acceptLanguage.String())
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		defer func() {
			require.NoError(t, resp.Body.Close())
		}()
		require.NoError(t, err)
		assert.Containsf(t, string(body), containsText, "login ui language is in "+expectLang.String())
	})
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
