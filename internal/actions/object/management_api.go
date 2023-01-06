package object

import (
	"github.com/dop251/goja"
	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/api/grpc/management"
)

func ManagementAPIField(server *management.Server) func(c *actions.FieldConfig) goja.Value {
	return func(c *actions.FieldConfig) goja.Value {
		return c.Runtime.ToValue(server)
	}
}
