package object

import (
	"maps"
	"net/http"

	"github.com/zitadel/zitadel/internal/actions"
)

// HTTPRequestField accepts the http.Request by value, so it's not mutated
func HTTPRequestField(request *http.Request) func(c *actions.FieldConfig) any {
	return func(c *actions.FieldConfig) any {
		return c.Runtime.ToValue(&httpRequest{
			Method:        request.Method,
			Url:           request.URL.String(),
			Proto:         request.Proto,
			ContentLength: request.ContentLength,
			Host:          request.Host,
			Form:          copyMap(request.Form),
			PostForm:      copyMap(request.PostForm),
			RemoteAddr:    request.RemoteAddr,
			Headers:       copyMap(request.Header),
		})
	}
}

type httpRequest struct {
	Method        string
	Url           string
	Proto         string
	ContentLength int64
	Host          string
	Form          map[string][]string
	PostForm      map[string][]string
	RemoteAddr    string
	Headers       map[string][]string
}

func copyMap(src map[string][]string) map[string][]string {
	dst := make(map[string][]string)
	maps.Copy(dst, src)
	return dst
}
