package robots_txt

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RobotsTxt(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	recorder := httptest.NewRecorder()

	handler, err := Start()
	handler.ServeHTTP(recorder, req)
	assert.Equal(t, nil, err)

	res := recorder.Result()
	body, err := io.ReadAll(res.Body)
	assert.Equal(t, nil, err)

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "User-agent: *\nDisallow: /\n", string(body))

	defer res.Body.Close()
}
