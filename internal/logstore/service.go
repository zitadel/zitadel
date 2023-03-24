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
	GetDueQuotaNotifications(ctx context.Context, config *quota.AddedEvent, periodStart time.Time, used uint64) ([]*quota.NotifiedEvent, error)
}

type UsageQuerier interface {
	LogEmitter
	QuotaUnit() quota.Unit
	QueryUsage(ctx context.Context, instanceId string, start time.Time) (uint64, error)
}

type UsageReporter interface {
	Report(ctx context.Context, notifications []*quota.NotifiedEvent) (err error)
}

type UsageReporterFunc func(context.Context, []*quota.NotifiedEvent) (err error)

func (u UsageReporterFunc) Report(ctx context.Context, notifications []*quota.NotifiedEvent) (err error) {
	return u(ctx, notifications)
}

type Service struct {
	usageQuerier     UsageQuerier
	quotaQuerier     QuotaQuerier
	usageReporter    UsageReporter
	enabledSinks     []*emitter
	sinkEnabled      bool
	reportingEnabled bool
}

func New(quotaQuerier QuotaQuerier, usageReporter UsageReporter, usageQuerierSink *emitter, additionalSink ...*emitter) *Service {
	var usageQuerier UsageQuerier
	if usageQuerierSink != nil {
		usageQuerier = usageQuerierSink.emitter.(UsageQuerier)
	}

	svc := &Service{
		reportingEnabled: usageQuerierSink != nil && usageQuerierSink.enabled,
		usageQuerier:     usageQuerier,
		quotaQuerier:     quotaQuerier,
		usageReporter:    usageReporter,
	}

	for _, s := range append([]*emitter{usageQuerierSink}, additionalSink...) {
		if s != nil && s.enabled {
			svc.enabledSinks = append(svc.enabledSinks, s)
		}
	}

	svc.sinkEnabled = len(svc.enabledSinks) > 0

	return svc
}

func (s *Service) Enabled() bool {
	return s.sinkEnabled
}

func (s *Service) Handle(ctx context.Context, record LogRecord) {
	for _, sink := range s.enabledSinks {
		logging.OnError(sink.Emit(ctx, record.Normalize())).WithField("record", record).Warn("failed to emit log record")
	}
}

func (s *Service) Limit(ctx context.Context, instanceID string) *uint64 {
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

func (s *Service) handleThresholds(ctx context.Context, quota *quota.AddedEvent, periodStart time.Time, usage uint64) {
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
