package connect_middleware

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func emptyMockHandler(resp connect.AnyResponse) func(_ context.Context, _ connect.AnyRequest) (connect.AnyResponse, error) {
	return func(_ context.Context, _ connect.AnyRequest) (connect.AnyResponse, error) {
		return resp, nil
	}
}

func errorMockHandler(_ context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
	return nil, zerrors.ThrowInternal(nil, "test", "error")
}

type mockReq struct {
	connect.Request[struct{}]

	procedure string
	header    http.Header
}

func (m *mockReq) Spec() connect.Spec {
	return connect.Spec{
		Procedure: m.procedure,
	}
}

func (m *mockReq) Header() http.Header {
	if m.header == nil {
		m.header = make(http.Header)
	}
	return m.header
}
