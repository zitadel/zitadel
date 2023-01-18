package object

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/actions"
)

// HTTPRequestField accepts the http.Request by value, so it's not mutated
func HTTPRequestField(httpRequest http.Request) func(c *actions.FieldConfig) interface{} {
	return func(c *actions.FieldConfig) interface{} {
		return c.Runtime.ToValue(httpRequest)
	}
}
