package logstore

import (
	"context"
	"database/sql"
	"math"
	"time"

	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

type reportFunc func(ctx context.Context, q *query.Quota, used uint64) (doLimit bool, err error)

type UsageQuerier interface {
	LogEmitter
	QuotaUnit() quota.Unit
	QueryUsage(ctx context.Context, instanceId string, start, end time.Time) (uint64, error)
}

type Service struct {
	usageQuerier     UsageQuerier
	dbClient         *sql.DB
	report           reportFunc
	enabledSinks     []*emitter
	sinkEnabled      bool
	reportingEnabled bool
}

func New(usageQuerierSink *emitter, dbClient *sql.DB, reportFunc reportFunc, additionalSink ...*emitter) *Service {

	svc := &Service{
		reportingEnabled: usageQuerierSink.enabled,
		usageQuerier:     usageQuerierSink.emitter.(UsageQuerier),
		dbClient:         dbClient,
		report:           reportFunc,
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
		if err := sink.Emit(ctx, record.RedactSecrets()); err != nil {
			return err
		}
	}
	return nil
}

// Limit TODO: Cache things in-memory here?
func (s *Service) Limit(ctx context.Context, instanceID string) (bool, *uint64, error) {
	if !s.reportingEnabled || instanceID == "" {
		return false, nil, nil
	}

	quota, err := query.GetInstanceQuota(ctx, s.dbClient, instanceID, s.usageQuerier.QuotaUnit())
	if err != nil {
		return false, nil, err
	}

	usage, err := s.usageQuerier.QueryUsage(ctx, instanceID, quota.PeriodStart, quota.PeriodEnd)
	if err != nil {
		return false, nil, err
	}

	remaining := uint64(math.Max(0, float64(quota.Amount)-float64(usage)))

	doLimit, err := s.report(ctx, quota, usage)
	return doLimit, &remaining, err
}
