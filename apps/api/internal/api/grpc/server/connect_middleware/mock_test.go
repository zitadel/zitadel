package connect_middleware

import (
	"context"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func emptyMockHandler(resp connect.AnyResponse, expectedCtxData authz.CtxData) func(*testing.T) connect.UnaryFunc {
	return func(t *testing.T) connect.UnaryFunc {
		return func(ctx context.Context, _ connect.AnyRequest) (connect.AnyResponse, error) {
			assert.Equal(t, expectedCtxData, authz.GetCtxData(ctx))
			return resp, nil
		}
	}
}

func errorMockHandler() func(*testing.T) connect.UnaryFunc {
	return func(t *testing.T) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			return nil, zerrors.ThrowInternal(nil, "test", "error")
		}
	}
}

type mockReq[t any] struct {
	connect.Request[t]

	procedure string
	header    http.Header
}

func (m *mockReq[T]) Spec() connect.Spec {
	return connect.Spec{
		Procedure: m.procedure,
	}
}

func (m *mockReq[T]) Header() http.Header {
	if m.header == nil {
		m.header = make(http.Header)
	}
	return m.header
}
