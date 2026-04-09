package connect_middleware

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_setInstance_errorCodes(t *testing.T) {
	i18n.SupportLanguages(language.English)
	translator := i18n.NewZitadelTranslator(language.English)

	cases := []struct {
		name     string
		err      error
		wantCode connect.Code
	}{
		{
			name:     "not found from verifier propagates as NotFound",
			err:      zerrors.ThrowNotFound(nil, "TEST-001", "Errors.Instance.NotFound"),
			wantCode: connect.CodeNotFound,
		},
		{
			name:     "internal error from verifier propagates as Internal",
			err:      zerrors.ThrowInternal(errors.New("FATAL: the database system is shutting down (SQLSTATE 57P03)"), "TEST-002", "Errors.Internal"),
			wantCode: connect.CodeInternal,
		},
		{
			name:     "unavailable error from verifier propagates as Internal",
			err:      zerrors.ThrowUnavailable(nil, "TEST-003", "Errors.Unavailable"),
			wantCode: connect.CodeInternal,
		},
	}

	for _, tc := range cases {
		verifier := &mockInstanceVerifier{err: tc.err}

		t.Run("byRequestedHost/"+tc.name, func(t *testing.T) {
			ctx := http_util.WithDomainContext(context.Background(), &http_util.DomainCtx{InstanceHost: "host"})
			_, err := setInstance(ctx, &mockReq[struct{}]{}, nil, verifier, "", translator)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			var connectErr *connect.Error
			if !errors.As(err, &connectErr) {
				t.Fatalf("expected *connect.Error, got %T", err)
			}
			if got := connectErr.Code(); got != tc.wantCode {
				t.Errorf("got code %v, want %v", got, tc.wantCode)
			}
		})
	}
}

type mockInstanceVerifier struct {
	err error
}

func (m *mockInstanceVerifier) InstanceByHost(_ context.Context, _, _ string) (authz.Instance, error) {
	return nil, m.err
}

func (m *mockInstanceVerifier) InstanceByID(_ context.Context, _ string) (authz.Instance, error) {
	return nil, m.err
}
