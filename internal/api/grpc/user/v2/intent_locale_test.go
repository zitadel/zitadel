package user

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/idp/providers/google"
	oidcprovider "github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

func Test_idpUserUpdates_preferredLanguage(t *testing.T) {
	for _, tt := range []struct {
		name   string
		locale *oidc.Locale
		want   *string
	}{
		{name: "missing locale preserves existing language"},
		{name: "undefined locale preserves existing language", locale: oidc.NewLocale(language.Und)},
		{name: "valid locale updates language", locale: oidc.NewLocale(language.English), want: gu.Ptr("en")},
		{name: "regional locale updates language", locale: oidc.NewLocale(language.AmericanEnglish), want: gu.Ptr("en-US")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			idpUser := &google.User{User: oidcprovider.NewUser(&oidc.UserInfo{
				UserInfoProfile: oidc.UserInfoProfile{GivenName: "Updated", FamilyName: "User", Locale: tt.locale},
			})}

			t.Run("UpdateUser", func(t *testing.T) {
				request := idpUserToUpdateUser("user-id", idpUser).UpdateUser
				profile := request.GetHuman().GetProfile()
				require.NotNil(t, profile)
				assert.Equal(t, tt.want, profile.PreferredLanguage)
				assert.Equal(t, "Updated", profile.GetGivenName())
				commandProfile := SetHumanProfileToProfile(profile)
				assertIdpLanguagePatch(t, commandProfile, tt.want)
				if tt.want == nil {
					// A nil patch field leaves the stored profile language unchanged.
					assert.Nil(t, commandProfile.PreferredLanguage)
				} else {
					require.NotNil(t, commandProfile.PreferredLanguage)
					assert.Equal(t, *tt.want, commandProfile.PreferredLanguage.String())
				}
			})
			t.Run("UpdateHumanUser", func(t *testing.T) {
				request := idpUserToUpdateHumanUser("user-id", idpUser)
				require.NotNil(t, request.Profile)
				assert.Equal(t, tt.want, request.Profile.PreferredLanguage)
				assert.Equal(t, "Updated", request.Profile.GivenName)
				command, err := updateHumanUserRequestToChangeHuman(request)
				require.NoError(t, err)
				assertIdpLanguagePatch(t, command.Profile, tt.want)
				if tt.want == nil {
					assert.Nil(t, command.Profile.PreferredLanguage)
				} else {
					require.NotNil(t, command.Profile.PreferredLanguage)
					assert.Equal(t, *tt.want, command.Profile.PreferredLanguage.String())
				}
			})
		})
	}
}

// Exercise the domain event/reducer as well as the transport conversion: an
// absent IdP locale must leave an existing Italian preference untouched while
// other profile attributes are still synchronized.
func assertIdpLanguagePatch(t *testing.T, profile *command.Profile, want *string) {
	t.Helper()
	wm := command.NewUserHumanWriteModel("user-id", "org-id", true, false, false, false, false, false, false)
	wm.PreferredLanguage = language.Italian
	event, err := wm.NewProfileChangedEvent(context.Background(), profile.FirstName, profile.LastName,
		profile.NickName, profile.DisplayName, profile.PreferredLanguage, profile.Gender)
	require.NoError(t, err)
	require.NotNil(t, event)
	wm.AppendEvents(event)
	require.NoError(t, wm.Reduce())
	assert.Equal(t, "Updated", wm.FirstName)
	expected := "it"
	if want != nil {
		expected = *want
	}
	assert.Equal(t, expected, wm.PreferredLanguage.String())
}
