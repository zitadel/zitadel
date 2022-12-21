package actions

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/logstore"
)

var (
	ServerLog   *logrus
	logstoreSvc *logstore.Service
)

func SetLogstoreService(svc *logstore.Service) {
	logstoreSvc = svc
}

type logrus struct{}

func (*logrus) Log(s string) {
	logging.WithFields("message", s).Info("log from action")
}
func (*logrus) Warn(s string) {
	logging.WithFields("message", s).Info("warn from action")
}
func (*logrus) Error(s string) {
	logging.WithFields("message", s).Info("error from action")
}

func WithLogger(logger console.Printer) Option {
	return func(c *runConfig) {
		c.modules["zitadel/log"] = func(runtime *goja.Runtime, module *goja.Object) {
			console.RequireWithPrinter(logger)(runtime, module)
		}
	}
}
