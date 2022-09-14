package actions

import (
	"github.com/zitadel/logging"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
)

var ServerLog *logrus

type logrus struct{}

func (*logrus) Log(s string) {
	logging.WithFields("message", s).Info("log from action")
}
func (*logrus) Warn(s string) {
	logging.WithFields("message", s).Info("log from action")
}
func (*logrus) Error(s string) {
	logging.WithFields("message", s).Info("log from action")
}

func WithLogger(logger console.Printer) Option {
	return func(c *runConfig) {
		c.modules["zitadel/log"] = func(runtime *goja.Runtime, module *goja.Object) {
			console.RequireWithPrinter(logger)(runtime, module)
		}
	}
}
