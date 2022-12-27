package actions

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"

	"github.com/zitadel/logging"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/sirupsen/logrus"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/emitters/execution"
)

var (
	logstoreService *logstore.Service
)

func SetLogstoreService(svc *logstore.Service) {
	logstoreService = svc
}

var _ console.Printer = (*logger)(nil)

type logger struct {
	ctx                             context.Context
	started                         time.Time
	instanceID, projectID, actionID string
	metadata                        map[string]interface{}
}

// newLogger returns a *logger instance that should only be used for a single action run.
// The first log call sets the started field for subsequent log calls
func newLogger(ctx context.Context, actionID string, metadata map[string]interface{}) *logger {
	instance := authz.GetInstance(ctx)

	return &logger{
		ctx:        ctx,
		started:    time.Time{},
		instanceID: instance.InstanceID(),
		projectID:  instance.ProjectID(),
		actionID:   actionID,
		metadata:   metadata,
	}
}

func (l *logger) Log(msg string) {
	l.log(msg, logrus.InfoLevel)
}

func (l *logger) Warn(msg string) {
	l.log(msg, logrus.WarnLevel)
}

func (l *logger) Error(msg string) {
	l.log(msg, logrus.ErrorLevel)
}

func (l *logger) log(msg string, level logrus.Level) {

	ts := time.Now()
	if l.started.IsZero() {
		l.started = ts
	}

	record := &execution.Record{
		Timestamp:  ts,
		Started:    l.started,
		InstanceID: l.instanceID,
		ProjectID:  l.projectID,
		ActionID:   l.actionID,
		Message:    msg,
		LogLevel:   level,
		Metadata:   l.metadata,
	}
	if err := logstoreService.Handle(context.TODO() /* TODO: context */, record); err != nil {
		logging.WithError(err).WithField("record", record).Errorf("handling execution log failed")
	}
}

func WithLogger(ctx context.Context, actionID string, metadata map[string]interface{}) Option {
	return func(c *runConfig) {
		c.logger = newLogger(ctx, actionID, metadata)
		c.modules["zitadel/log"] = func(runtime *goja.Runtime, module *goja.Object) {
			console.RequireWithPrinter(c.logger)(runtime, module)
		}
	}
}
