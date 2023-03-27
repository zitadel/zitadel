package actions

import (
	"context"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/sirupsen/logrus"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/emitters/execution"
)

var (
	logstoreService *logstore.Service
	_               console.Printer = (*logger)(nil)
)

func SetLogstoreService(svc *logstore.Service) {
	logstoreService = svc
}

type logger struct {
	ctx        context.Context
	started    time.Time
	instanceID string
}

// newLogger returns a *logger instance that should only be used for a single action run.
// The first log call sets the started field for subsequent log calls
func newLogger(ctx context.Context, instanceID string) *logger {
	return &logger{
		ctx:        ctx,
		started:    time.Time{},
		instanceID: instanceID,
	}
}

func (l *logger) Log(msg string) {
	l.log(msg, logrus.InfoLevel, false)
}

func (l *logger) Warn(msg string) {
	l.log(msg, logrus.WarnLevel, false)
}

func (l *logger) Error(msg string) {
	l.log(msg, logrus.ErrorLevel, false)
}

func (l *logger) log(msg string, level logrus.Level, last bool) {
	ts := time.Now()
	if l.started.IsZero() {
		l.started = ts
	}

	record := &execution.Record{
		LogDate:    ts,
		InstanceID: l.instanceID,
		Message:    msg,
		LogLevel:   level,
	}

	if last {
		record.Took = ts.Sub(l.started)
	}

	logstoreService.Handle(l.ctx, record)
}

func withLogger(ctx context.Context) Option {
	instance := authz.GetInstance(ctx)
	instanceID := instance.InstanceID()
	return func(c *runConfig) {
		c.logger = newLogger(ctx, instanceID)
		c.instanceID = instanceID
		c.modules["zitadel/log"] = func(runtime *goja.Runtime, module *goja.Object) {
			console.RequireWithPrinter(c.logger)(runtime, module)
		}
	}
}
