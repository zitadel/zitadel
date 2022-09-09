package actions

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
)

func WithLogger(logger console.Printer) runOpt {
	return func(c *runConfig) {
		c.modules["zitadel/log"] = func(runtime *goja.Runtime, module *goja.Object) {
			console.RequireWithPrinter(logger)(runtime, module)
		}
	}
}
