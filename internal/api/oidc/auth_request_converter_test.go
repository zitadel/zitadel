package oidc

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/domain"
)

func TestResponseModeToBusiness(t *testing.T) {
	type args struct {
		responseMode oidc.ResponseMode
	}
	tests := []struct {
		name string
		args args
		want domain.OIDCResponseMode
	}{
		{
			name: "empty",
			args: args{""},
			want: domain.OIDCResponseModeUnspecified,
		},
		{
			name: "invalid",
			args: args{"foo"},
			want: domain.OIDCResponseModeUnspecified,
		},
		{
			name: "query",
			args: args{oidc.ResponseModeQuery},
			want: domain.OIDCResponseModeQuery,
		},
		{
			name: "fragment",
			args: args{oidc.ResponseModeFragment},
			want: domain.OIDCResponseModeFragment,
		},
		{
			name: "post_form",
			args: args{oidc.ResponseModeFormPost},
			want: domain.OIDCResponseModeFormPost,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResponseModeToBusiness(tt.args.responseMode)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResponseModeToOIDC(t *testing.T) {
	type args struct {
		responseMode domain.OIDCResponseMode
	}
	tests := []struct {
		name string
		args args
		want oidc.ResponseMode
	}{
		{
			name: "unspecified",
			args: args{domain.OIDCResponseModeUnspecified},
			want: "",
		},
		{
			name: "invalid",
			args: args{99},
			want: "",
		},
		{
			name: "query",
			args: args{domain.OIDCResponseModeQuery},
			want: oidc.ResponseModeQuery,
		},
		{
			name: "fragment",
			args: args{domain.OIDCResponseModeFragment},
			want: oidc.ResponseModeFragment,
		},
		{
			name: "form_post",
			args: args{domain.OIDCResponseModeFormPost},
			want: oidc.ResponseModeFormPost,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResponseModeToOIDC(tt.args.responseMode)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPromptToBusiness(t *testing.T) {
	type args struct {
		oidcPrompt []string
	}
	tests := []struct {
		name string
		args args
		want []domain.Prompt
	}{
		{
			name: "unspecified",
			args: args{nil},
			want: []domain.Prompt{},
		},
		{
			name: "invalid",
			args: args{[]string{"non_existing_prompt"}},
			want: []domain.Prompt{},
		},
		{
			name: "prompt_none",
			args: args{[]string{oidc.PromptNone}},
			want: []domain.Prompt{domain.PromptNone},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PromptToBusiness(tt.args.oidcPrompt)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAuthRequest_LogValue(t *testing.T) {
	tests := []struct {
		name     string
		req      *AuthRequest
		expected slog.Value
	}{
		{
			name:     "nil receiver",
			req:      nil,
			expected: slog.Value{},
		},
		{
			name:     "nil inner request",
			req:      &AuthRequest{},
			expected: slog.Value{},
		},
		{
			name: "populated auth request",
			req: &AuthRequest{
				AuthRequest: &domain.AuthRequest{
					ID:         "auth-123",
					InstanceID: "inst-456",
					ApplicationID:    "client-789",
					CallbackURI: "https://example.com/callback",
					Request: &domain.AuthRequestOIDC{
						Scopes:      []string{"openid", "profile"},
					},
				},
			},
			expected: slog.GroupValue(
				slog.String("id", "auth-123"),
				slog.String("client_id", "client-789"),
				slog.String("instance_id", "inst-456"),
				slog.Any("scopes", []string{"openid", "profile"}),
				slog.String("redirect_uri", "https://example.com/callback"),
			),
		},
		{
			name: "populated auth request with nil inner request",
			req: &AuthRequest{
				AuthRequest: &domain.AuthRequest{
					ID:            "auth-123",
					InstanceID:    "inst-456",
					ApplicationID: "client-789",
					CallbackURI:   "https://example.com/callback",
					Request:       nil,
				},
			},
			expected: slog.GroupValue(
				slog.String("id", "auth-123"),
				slog.String("client_id", "client-789"),
				slog.String("instance_id", "inst-456"),
				slog.Any("scopes", []string(nil)),
				slog.String("redirect_uri", "https://example.com/callback"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.LogValue()
			assert.Equal(t, tt.expected, got)
		})
	}
}