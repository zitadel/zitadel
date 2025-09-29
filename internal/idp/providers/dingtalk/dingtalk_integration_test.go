package dingtalk
package dingtalk_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/zitadel/internal/idp/providers/dingtalk"
)

func TestDingTalkProvider_New(t *testing.T) {
	tests := []struct {
		name        string
		clientID    string
		secret      string
		callbackURL string
		scopes      []string
		wantErr     bool
	}{
		{
			name:        "valid provider creation",
			clientID:    "test-client-id",
			secret:      "test-secret",
			callbackURL: "https://example.com/callback",
			scopes:      []string{"openid", "profile"},
			wantErr:     false,
		},
		{
			name:        "empty client ID",
			clientID:    "",
			secret:      "test-secret",
			callbackURL: "https://example.com/callback",
			scopes:      []string{"openid", "profile"},
			wantErr:     false, // OAuth provider allows empty client ID
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := dingtalk.New(tt.clientID, tt.secret, tt.callbackURL, tt.scopes)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, provider)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
				assert.Equal(t, "DingTalk", provider.Name())
			}
		})
	}
}

func TestDingTalkUser_Methods(t *testing.T) {
	user := &dingtalk.User{
		UnionID:   "test_union_123",
		Nick:      "张三",
		Email:     "zhangsan@example.com",
		Mobile:    "13800138000",
		AvatarURL: "https://avatar.example.com/user.jpg",
	}

	tests := []struct {
		name     string
		method   func() interface{}
		expected interface{}
	}{
		{"GetID", func() interface{} { return user.GetID() }, "test_union_123"},
		{"GetDisplayName", func() interface{} { return user.GetDisplayName() }, "张三"},
		{"GetNickname", func() interface{} { return user.GetNickname() }, "张三"},
		{"GetPreferredUsername", func() interface{} { return user.GetPreferredUsername() }, "张三"},
		{"GetEmail", func() interface{} { return string(user.GetEmail()) }, "zhangsan@example.com"},
		{"GetPhone", func() interface{} { return string(user.GetPhone()) }, "13800138000"},
		{"GetAvatarURL", func() interface{} { return user.GetAvatarURL() }, "https://avatar.example.com/user.jpg"},
		{"IsEmailVerified", func() interface{} { return user.IsEmailVerified() }, true},
		{"IsPhoneVerified", func() interface{} { return user.IsPhoneVerified() }, true},
		{"GetFirstName", func() interface{} { return user.GetFirstName() }, ""},
		{"GetLastName", func() interface{} { return user.GetLastName() }, ""},
		{"GetProfile", func() interface{} { return user.GetProfile() }, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.method()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDingTalkUser_PhoneVerification(t *testing.T) {
	tests := []struct {
		name           string
		mobile         string
		expectedVerified bool
	}{
		{"with mobile", "13800138000", true},
		{"without mobile", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &dingtalk.User{
				UnionID: "test_union_123",
				Mobile:  tt.mobile,
			}
			assert.Equal(t, tt.expectedVerified, user.IsPhoneVerified())
		})
	}
}

func TestDingTalkUser_Language(t *testing.T) {
	user := &dingtalk.User{}
	lang := user.GetPreferredLanguage()
	assert.Equal(t, "zh", lang.String())
}