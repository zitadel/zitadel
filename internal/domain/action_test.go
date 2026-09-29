package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActionFunction_LocalizationKey(t *testing.T) {
	tests := []struct {
		fn   ActionFunction
		want string
	}{
		{ActionFunctionPreUserinfo, "preuserinfo"},
		{ActionFunctionPreAccessToken, "preaccesstoken"},
		{ActionFunctionPreSAMLResponse, "presamlresponse"},
		{ActionFunctionPreOTPSMSCode, "preotpsmscode"},
		{ActionFunctionPreOTPEmailCode, "preotpemailcode"},
		{ActionFunctionUnspecified, "unspecified"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.fn.LocalizationKey())
		})
	}
}

func TestActionFunctionExists(t *testing.T) {
	exists := ActionFunctionExists()
	for _, name := range []string{"preuserinfo", "preaccesstoken", "presamlresponse", "preotpsmscode", "preotpemailcode"} {
		t.Run(name, func(t *testing.T) {
			assert.True(t, exists(name))
		})
	}
	assert.False(t, exists("unspecified"))
	assert.False(t, exists("bogus"))
}
