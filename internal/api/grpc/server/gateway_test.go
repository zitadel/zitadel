package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorHandler_PreservesUnauthenticatedStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/auth/v1/users/me", nil)
	resp := httptest.NewRecorder()

	errorHandler(context.Background(), runtime.NewServeMux(), jsonMarshaler, resp, req, status.Error(codes.Unauthenticated, "auth header missing"))

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), `"code":16`)
	assert.Contains(t, resp.Body.String(), `"message":"auth header missing"`)
}
