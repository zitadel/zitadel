package middleware

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestRecoverHandler(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantStatus int
		wantBody   string
	}{
		{
			name: "no panic",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, "ok") //nolint:errcheck
			},
			wantStatus: http.StatusOK,
			wantBody:   "ok",
		},
		{
			name: "panic with error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic(errors.New("oops"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "ID=RECOVER Message=Errors.Internal Parent=(oops)",
		},
		{
			name: "panic with string",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic("something went wrong")
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "ID=RECOVER Message=Errors.Internal Parent=(something went wrong)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RecoverHandler(writeResponse)(tt.handler)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			handler.ServeHTTP(w, r)

			res := w.Result()
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.wantStatus, res.StatusCode)
			assert.Equal(t, tt.wantBody, string(body))
		})
	}
}

func TestFallbackRecoverHandler(t *testing.T) {
	handler := FallbackRecoverHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("unexpected error")
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, "Internal Server Error\n", string(body))
}

func TestRecoverHandlerWithError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	// Panic on nil handler.
	err := RecoverHandlerWithError(nil)(w, r)
	assert.ErrorIs(t, err, zerrors.ThrowInternal(nil, zerrors.IDRecover, "Errors.Internal"))
}

func writeResponse(w http.ResponseWriter, _ *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, err.Error()) //nolint:errcheck
}
