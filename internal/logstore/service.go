package logstore

import (
	"context"
	"math"
	"time"

	"github.com/zitadel/zitadel/internal/repository/quota"
)

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

func (s *Service) Handle(ctx context.Context, record LogRecord) error {
	for _, sink := range s.enabledSinks {
		if err := sink.Emit(ctx, record.Normalize()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Limit(ctx context.Context, instanceID string) (*uint64, error) {
	if !s.reportingEnabled || instanceID == "" {
		return nil, nil
	}

	quota, periodStart, err := s.quotaQuerier.GetCurrentQuotaPeriod(ctx, instanceID, s.usageQuerier.QuotaUnit())
	if err != nil || quota == nil {
		return nil, err
	}

	usage, err := s.usageQuerier.QueryUsage(ctx, instanceID, periodStart)
	if err != nil {
		return nil, err
	}

	var remaining *uint64
	if quota.Limit {
		r := uint64(math.Max(0, float64(quota.Amount)-float64(usage)))
		remaining = &r
	}

	notifications, err := s.quotaQuerier.GetDueQuotaNotifications(ctx, quota, periodStart, usage)
	if err != nil {
		return remaining, err
	}

	return remaining, s.usageReporter.Report(ctx, notifications)
}
