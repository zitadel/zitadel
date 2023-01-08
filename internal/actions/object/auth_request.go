package object

import (
	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/domain"
)

// AuthRequestField accepts the domain.AuthRequest by value, so its not mutated
func AuthRequestField(authRequest domain.AuthRequest) func(c *actions.FieldConfig) interface{} {
	return func(c *actions.FieldConfig) interface{} {
		return c.Runtime.ToValue(authRequest)
	}
}
