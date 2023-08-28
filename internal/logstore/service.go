package logstore

import (
	"context"
	"math"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

const handleThresholdTimeout = time.Minute

type QuotaQuerier interface {
	GetCurrentQuotaPeriod(ctx context.Context, instanceID string, unit quota.Unit) (config *quota.AddedEvent, periodStart time.Time, err error)
	GetDueQuotaNotifications(ctx context.Context, config *quota.AddedEvent, periodStart time.Time, used uint64) ([]*quota.NotificationDueEvent, error)
}

type UsageQuerier[T LogRecord[T]] interface {
	LogEmitter[T]
	QuotaUnit() quota.Unit
	QueryUsage(ctx context.Context, instanceId string, start time.Time) (uint64, error)
}

type UsageReporter interface {
	Report(ctx context.Context, notifications []*quota.NotificationDueEvent) (err error)
}

type UsageReporterFunc func(context.Context, []*quota.NotificationDueEvent) (err error)

func (u UsageReporterFunc) Report(ctx context.Context, notifications []*quota.NotificationDueEvent) (err error) {
	return u(ctx, notifications)
}

type Service[T LogRecord[T]] struct {
	usageQuerier     UsageQuerier[T]
	quotaQuerier     QuotaQuerier
	usageReporter    UsageReporter
	enabledSinks     []*emitter[T]
	sinkEnabled      bool
	reportingEnabled bool
}

func New[T LogRecord[T]](quotaQuerier QuotaQuerier, usageReporter UsageReporter, usageQuerierSink *emitter[T], additionalSink ...*emitter[T]) *Service[T] {
	var usageQuerier UsageQuerier[T]
	if usageQuerierSink != nil {
		usageQuerier = usageQuerierSink.emitter.(UsageQuerier[T])
	}

	svc := &Service[T]{
		reportingEnabled: usageQuerierSink != nil && usageQuerierSink.enabled,
		usageQuerier:     usageQuerier,
		quotaQuerier:     quotaQuerier,
		usageReporter:    usageReporter,
	}

	for _, s := range append([]*emitter[T]{usageQuerierSink}, additionalSink...) {
		if s != nil && s.enabled {
			svc.enabledSinks = append(svc.enabledSinks, s)
		}
	}

	svc.sinkEnabled = len(svc.enabledSinks) > 0

	return svc
}

func (s *Service[T]) Enabled() bool {
	return s.sinkEnabled
}

func (s *Service[T]) Handle(ctx context.Context, record T) {
	for _, sink := range s.enabledSinks {
		logging.OnError(sink.Emit(ctx, record.Normalize())).WithField("record", record).Warn("failed to emit log record")
	}
}

func (s *Service[T]) Limit(ctx context.Context, instanceID string) *uint64 {
	var err error
	defer func() {
		logging.OnError(err).Warn("failed to check is usage should be limited")
	}()

	if !s.reportingEnabled || instanceID == "" {
		return nil
	}

	quota, periodStart, err := s.quotaQuerier.GetCurrentQuotaPeriod(ctx, instanceID, s.usageQuerier.QuotaUnit())
	if err != nil || quota == nil {
		return nil
	}

	usage, err := s.usageQuerier.QueryUsage(ctx, instanceID, periodStart)
	if err != nil {
		return nil
	}

	go s.handleThresholds(ctx, quota, periodStart, usage)

	var remaining *uint64
	if quota.Limit {
		r := uint64(math.Max(0, float64(quota.Amount)-float64(usage)))
		remaining = &r
	}
	return remaining
}

func (s *Service[T]) handleThresholds(ctx context.Context, quota *quota.AddedEvent, periodStart time.Time, usage uint64) {
	var err error
	defer func() {
		logging.OnError(err).Warn("handling quota thresholds failed")
	}()

	detatchedCtx, cancel := context.WithTimeout(authz.Detach(ctx), handleThresholdTimeout)
	defer cancel()

	notifications, err := s.quotaQuerier.GetDueQuotaNotifications(detatchedCtx, quota, periodStart, usage)
	if err != nil || len(notifications) == 0 {
		return
	}

	err = s.usageReporter.Report(detatchedCtx, notifications)
}
