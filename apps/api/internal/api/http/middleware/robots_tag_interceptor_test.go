package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RobotsTagInterceptor(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	handler := RobotsTagHandler(http.HandlerFunc(testHandler))
	handler.ServeHTTP(recorder, req)

	res := recorder.Result()
	exp := res.Header.Get("X-Robots-Tag")
	assert.Equal(t, "none", exp)

	defer res.Body.Close()
}
