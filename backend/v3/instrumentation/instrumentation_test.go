package instrumentation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestFilter(t *testing.T) {
	filter := RequestFilter("/foo", "/bar")
	var tests = []struct {
		name string
		r    *http.Request
		want bool
	}{
		{
			name: "foo subpath false",
			r:    httptest.NewRequest("POST", "/foo/some", nil),
			want: false,
		},
		{
			name: "foo exact false",
			r:    httptest.NewRequest("POST", "/foo", nil),
			want: false,
		},
		{
			name: "other path true",
			r:    httptest.NewRequest("POST", "/other", nil),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filter(tt.r)
			assert.Equal(t, tt.want, got)
		})
	}
}
