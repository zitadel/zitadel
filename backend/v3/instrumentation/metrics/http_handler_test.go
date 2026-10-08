package metrics

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestURI(t *testing.T) {
	tests := []struct {
		name       string
		requestURI string
		// prepare returns the context the request is served with.
		prepare func(context.Context) context.Context
		want    string
	}{
		{
			name:       "no pattern support, path is recorded",
			requestURI: "/oauth/v2/token",
			prepare:    func(ctx context.Context) context.Context { return ctx },
			want:       "/oauth/v2/token",
		},
		{
			name:       "no pattern support, query is stripped",
			requestURI: "/oauth/v2/authorize?client_id=1234&state=abcd",
			prepare:    func(ctx context.Context) context.Context { return ctx },
			want:       "/oauth/v2/authorize",
		},
		{
			name:       "pattern supported but not set, path is recorded",
			requestURI: "/oauth/v2/token",
			prepare:    WithRequestURIPattern,
			want:       "/oauth/v2/token",
		},
		{
			name:       "pattern set, pattern is recorded instead of the object id",
			requestURI: "/v2/sessions/385063742120926058",
			prepare: func(ctx context.Context) context.Context {
				ctx = WithRequestURIPattern(ctx)
				SetRequestURIPattern(ctx, "/v2/sessions/{session_id}")
				return ctx
			},
			want: "/v2/sessions/{session_id}",
		},
		{
			name:       "unknown path, no path is recorded",
			requestURI: "/junk",
			prepare: func(ctx context.Context) context.Context {
				ctx = WithRequestURIPattern(ctx)
				SetRequestURIPattern(ctx, UnknownPath)
				return ctx
			},
			want: UnknownPath,
		},
		{
			name:       "pattern set without context support, path is recorded",
			requestURI: "/v2/sessions/385063742120926058",
			prepare: func(ctx context.Context) context.Context {
				SetRequestURIPattern(ctx, "/v2/sessions/{session_id}")
				return ctx
			},
			want: "/v2/sessions/385063742120926058",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.requestURI, nil)
			r = r.WithContext(tt.prepare(r.Context()))

			assert.Equal(t, tt.want, RequestURI(r))
		})
	}
}

// TestSetRequestURIPattern_derivedContext ensures the pattern set by a handler is visible
// to the middleware that seeded the context, even though the handler received a derived
// context. Routers such as the grpc-gateway only ever report the pattern on their own
// context, so a value copied by [context.WithValue] would not reach the caller.
func TestSetRequestURIPattern_derivedContext(t *testing.T) {
	r := httptest.NewRequest(http.MethodDelete, "/v2/sessions/385063742120926058", nil)
	r = r.WithContext(WithRequestURIPattern(r.Context()))

	type unrelatedKey struct{}
	derived := context.WithValue(r.Context(), unrelatedKey{}, "value")
	SetRequestURIPattern(derived, "/v2/sessions/{session_id}")

	assert.Equal(t, "/v2/sessions/{session_id}", RequestURI(r))
}
