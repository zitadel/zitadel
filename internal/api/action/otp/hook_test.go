package otp

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestPreOTPSMSCodeContext_GetHTTPRequestBody(t *testing.T) {
	ctx := &PreOTPSMSCodeContext{
		FunctionName:         "preotpsmscode",
		RecipientPhoneNumber: "+15555551234",
		GeneratorConfig: &PublicGeneratorConfig{
			Length:        6,
			Expiry:        Duration(5 * time.Minute),
			IncludeDigits: true,
		},
	}
	body := ctx.GetHTTPRequestBody()
	assert.True(t, json.Valid(body))

	var decoded map[string]any
	assert.NoError(t, json.Unmarshal(body, &decoded))
	assert.Equal(t, "preotpsmscode", decoded["function"])
	assert.Equal(t, "+15555551234", decoded["recipient_phone_number"])
}

func TestPreOTPSMSCodeContext_SetHTTPResponseBody(t *testing.T) {
	tests := []struct {
		name         string
		resp         string
		wantErr      error
		wantCode     *string
		wantGenerate *GenerationOverrides
		wantExpiry   *Duration
	}{
		{
			name:       "code only",
			resp:       `{"code":"A7F2B9","expiry":"5m0s"}`,
			wantCode:   gu.Ptr("A7F2B9"),
			wantExpiry: gu.Ptr(Duration(5 * time.Minute)),
		},
		{
			name: "generate only",
			resp: `{"generate":{"length":4,"include_digits":true}}`,
			wantGenerate: &GenerationOverrides{
				Length:        gu.Ptr(uint32(4)),
				IncludeDigits: gu.Ptr(true),
			},
		},
		{
			name:    "code and generate both set rejected",
			resp:    `{"code":"123","generate":{"length":4}}`,
			wantErr: zerrors.ThrowPreconditionFailed(nil, "ACTION-k3j9z", "Errors.Execution.Invalid"),
		},
		{
			name:    "invalid json rejected",
			resp:    `not-json`,
			wantErr: zerrors.ThrowPreconditionFailed(nil, "ACTION-p7q2w", "Errors.Execution.ResponseIsNotValidJSON"),
		},
		{
			name: "empty response uses defaults",
			resp: `{}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &PreOTPSMSCodeContext{}
			err := ctx.SetHTTPResponseBody([]byte(tt.resp))
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr != nil {
				return
			}
			assert.Equal(t, tt.wantCode, ctx.Response.Code)
			assert.Equal(t, tt.wantGenerate, ctx.Response.Generate)
			assert.Equal(t, tt.wantExpiry, ctx.Response.Expiry)
		})
	}
}

func TestPreOTPSMSCodeContext_GetContent(t *testing.T) {
	ctx := &PreOTPSMSCodeContext{
		Response: &PreOTPSMSCodeResponse{Code: gu.Ptr("abc")},
	}
	got, ok := ctx.GetContent().(*PreOTPSMSCodeResponse)
	assert.True(t, ok)
	assert.Equal(t, "abc", *got.Code)
}

func TestPreOTPEmailCodeContext_GetHTTPRequestBody(t *testing.T) {
	ctx := &PreOTPEmailCodeContext{
		FunctionName:          "preotpemailcode",
		RecipientEmailAddress: "alice@example.com",
	}
	body := ctx.GetHTTPRequestBody()
	assert.True(t, json.Valid(body))

	var decoded map[string]any
	assert.NoError(t, json.Unmarshal(body, &decoded))
	assert.Equal(t, "preotpemailcode", decoded["function"])
	assert.Equal(t, "alice@example.com", decoded["recipient_email_address"])
}

func TestPreOTPEmailCodeContext_SetHTTPResponseBody(t *testing.T) {
	tests := []struct {
		name    string
		resp    string
		wantErr error
	}{
		{"code only", `{"code":"123456"}`, nil},
		{"generate only", `{"generate":{"length":8}}`, nil},
		{"empty", `{}`, nil},
		{
			"both set rejected",
			`{"code":"123","generate":{"length":8}}`,
			zerrors.ThrowPreconditionFailed(nil, "ACTION-r5t8b", "Errors.Execution.Invalid"),
		},
		{
			"invalid json rejected",
			`{`,
			zerrors.ThrowPreconditionFailed(nil, "ACTION-m8x4n", "Errors.Execution.ResponseIsNotValidJSON"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &PreOTPEmailCodeContext{}
			err := ctx.SetHTTPResponseBody([]byte(tt.resp))
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
