package logstore

import (
	"context"
	"math"
	"time"

	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

type UsageQuerier interface {
	LogEmitter
	QuotaUnit() quota.Unit
	QueryUsage(ctx context.Context, instanceId string, start, end time.Time) (uint64, error)
}

type UsageReporter interface {
	GetQuota(ctx context.Context, instanceID string, unit quota.Unit) (*query.Quota, error)
	// TODO: Determining doLimit is the services responsibility
	Report(ctx context.Context, q *query.Quota, used uint64) (err error)
}

type Service struct {
	usageQuerier     UsageQuerier
	usageReporter    UsageReporter
	enabledSinks     []*emitter
	sinkEnabled      bool
	reportingEnabled bool
}

func New(usageReporter UsageReporter, usageQuerierSink *emitter, additionalSink ...*emitter) *Service {

	svc := &Service{
		reportingEnabled: usageQuerierSink.enabled,
		usageQuerier:     usageQuerierSink.emitter.(UsageQuerier),
		usageReporter:    usageReporter,
	}

	for _, s := range append([]*emitter{usageQuerierSink}, additionalSink...) {
		if s.enabled {
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
		if err := sink.Emit(ctx, record.Redact()); err != nil {
			return err
		}
	}
	return nil
}

// Limit TODO: Cache things in-memory here?
func (s *Service) Limit(ctx context.Context, instanceID string) (*uint64, error) {
	if !s.reportingEnabled || instanceID == "" {
		return nil, nil
	}

	quota, err := s.usageReporter.GetQuota(ctx, instanceID, s.usageQuerier.QuotaUnit())
	if err != nil {
		return nil, err
	}

	usage, err := s.usageQuerier.QueryUsage(ctx, instanceID, quota.PeriodStart, quota.PeriodEnd)
	if err != nil {
		return nil, err
	}

	var remaining *uint64
	if quota.Limit {
		r := uint64(math.Max(0, float64(quota.Amount)-float64(usage)))
		remaining = &r
	}
	return remaining, s.usageReporter.Report(ctx, quota, usage)
}
