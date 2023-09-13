package oidc

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2beta"
)

func Test_authRequestToPb(t *testing.T) {
	now := time.Now()
	arg := &query.AuthRequest{
		ID:           "authID",
		CreationDate: now,
		ClientID:     "clientID",
		Scope:        []string{"a", "b", "c"},
		RedirectURI:  "callbackURI",
		Prompt: []domain.Prompt{
			domain.PromptUnspecified,
			domain.PromptNone,
			domain.PromptLogin,
			domain.PromptConsent,
			domain.PromptSelectAccount,
			domain.PromptCreate,
			999,
		},
		UiLocales:  []string{"en", "fi"},
		LoginHint:  gu.Ptr("foo@bar.com"),
		MaxAge:     gu.Ptr(time.Minute),
		HintUserID: gu.Ptr("userID"),
	}
	want := &oidc_pb.AuthRequest{
		Id:           "authID",
		CreationDate: timestamppb.New(now),
		ClientId:     "clientID",
		RedirectUri:  "callbackURI",
		Prompt: []oidc_pb.Prompt{
			oidc_pb.Prompt_PROMPT_UNSPECIFIED,
			oidc_pb.Prompt_PROMPT_NONE,
			oidc_pb.Prompt_PROMPT_LOGIN,
			oidc_pb.Prompt_PROMPT_CONSENT,
			oidc_pb.Prompt_PROMPT_SELECT_ACCOUNT,
			oidc_pb.Prompt_PROMPT_CREATE,
			oidc_pb.Prompt_PROMPT_UNSPECIFIED,
		},
		UiLocales:  []string{"en", "fi"},
		Scope:      []string{"a", "b", "c"},
		LoginHint:  gu.Ptr("foo@bar.com"),
		MaxAge:     durationpb.New(time.Minute),
		HintUserId: gu.Ptr("userID"),
	}
	got := authRequestToPb(arg)
	if !proto.Equal(want, got) {
		t.Errorf("authRequestToPb() =\n%v\nwant\n%v\n", got, want)
	}
}

func Test_errorReasonToOIDC(t *testing.T) {
	tests := []struct {
		reason oidc_pb.ErrorReason
		want   string
	}{
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_UNSPECIFIED,
			want:   "server_error",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_INVALID_REQUEST,
			want:   "invalid_request",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_UNAUTHORIZED_CLIENT,
			want:   "unauthorized_client",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_ACCESS_DENIED,
			want:   "access_denied",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_UNSUPPORTED_RESPONSE_TYPE,
			want:   "unsupported_response_type",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_INVALID_SCOPE,
			want:   "invalid_scope",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_SERVER_ERROR,
			want:   "server_error",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_TEMPORARY_UNAVAILABLE,
			want:   "temporarily_unavailable",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_INTERACTION_REQUIRED,
			want:   "interaction_required",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_LOGIN_REQUIRED,
			want:   "login_required",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_ACCOUNT_SELECTION_REQUIRED,
			want:   "account_selection_required",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_CONSENT_REQUIRED,
			want:   "consent_required",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_INVALID_REQUEST_URI,
			want:   "invalid_request_uri",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_INVALID_REQUEST_OBJECT,
			want:   "invalid_request_object",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_REQUEST_NOT_SUPPORTED,
			want:   "request_not_supported",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_REQUEST_URI_NOT_SUPPORTED,
			want:   "request_uri_not_supported",
		},
		{
			reason: oidc_pb.ErrorReason_ERROR_REASON_REGISTRATION_NOT_SUPPORTED,
			want:   "registration_not_supported",
		},
		{
			reason: 99999,
			want:   "server_error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.reason.String(), func(t *testing.T) {
			got := errorReasonToOIDC(tt.reason)
			assert.Equal(t, tt.want, got)
		})
	}
}
